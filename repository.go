package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

type Repository struct {
	DB     *bolt.DB
	Config Config
}

type ContentIngestionResult struct {
	Content    Content
	Errors     []error
	Task       string
	Successful bool
}

type ContentActor func(ContentIngestionResult) ContentIngestionResult

func (r *Repository) InitializeBuckets() {
	r.DB.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{
			"Content",
			"Nodes",
			"Timeline",
			"Spam",
			"Quarantine",
			"Tags",
			"Sums",
			"Votes",
			"Comments",
		} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
}

func (r *Repository) CheckSumId(cir ContentIngestionResult) ContentIngestionResult {
	_ = r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sums"))
		sum := sum(contentOnly(cir.Content))
		log.Printf("%s", b)
		log.Printf("%s", sum)
		contentId := b.Get([]byte(sum))
		if contentId != nil {
			uid, err := uuid.FromBytes(contentId)
			if err != nil {
				panic(err)
			}

			cir.Content.Id = uid
		}
		return nil
	})

	return cir
}
func (r *Repository) EnsureNotAlreadyPresent(cir ContentIngestionResult) ContentIngestionResult {
	_ = r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sums"))
		contentId := b.Get([]byte(sum(contentOnly(cir.Content))))
		cir.Successful = (contentId == nil)
		return nil
	})

	return cir
}

func (r *Repository) PredictNotSpam(cir ContentIngestionResult) ContentIngestionResult {
	//this is where the bayesian filter will go - some decisions need to be made
	return cir
}
func (r *Repository) PredictNotMalicious(cir ContentIngestionResult) ContentIngestionResult {
	// this is where we will compare the content to a list of known bad actors
	// and remove javascript
	return cir
}
func (r *Repository) EnsureId(cir ContentIngestionResult) ContentIngestionResult {
	if cir.Content.Id == uuid.Nil {
		cir.Content.Id = uuid.NewV4()
	}
	return cir
}

func (r *Repository) EnsurePosterId(cir ContentIngestionResult) ContentIngestionResult {
	if cir.Content.PosterId == uuid.Nil {
		cir.Content.PosterId = r.Config.NodeId
	}
	return cir
}
func (r *Repository) Sum(cir ContentIngestionResult) ContentIngestionResult {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sums"))
		err := b.Put([]byte(sum(contentOnly(cir.Content))), cir.Content.Id.Bytes())
		cir.Successful = (err == nil)
		return err
	})

	if err != nil {
		cir.Errors = append(cir.Errors, err)
	}
	return cir
}

func (r *Repository) Persist(cir ContentIngestionResult) ContentIngestionResult {
	if cir.Successful {
		err := r.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Content"))
			err := b.Put(cir.Content.Id.Bytes(), []byte(cir.Content.AsJson()))
			return err
		})

		if err != nil {
			cir.Errors = append(cir.Errors, err)
		}
	} else {
		cir.Errors = append(cir.Errors, errors.New("did not persist due to prior failure"))
	}
	return cir
}

func (r *Repository) Tag(cir ContentIngestionResult) ContentIngestionResult {
	if cir.Successful {
		_ = r.DB.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Tags"))
			for _, t := range cir.Content.Tags {
				tagBucket, bucketErr := b.CreateBucketIfNotExists([]byte(t))
				if bucketErr == nil {
					err := tagBucket.Put(cir.Content.Id.Bytes(), []byte(""))
					if err != nil {
						cir.Errors = append(cir.Errors, err)
						cir.Successful = false
					}
				} else {
					cir.Errors = append(cir.Errors, bucketErr)
					cir.Successful = false
				}
			}
			return nil
		})
	} else {
		cir.Errors = append(cir.Errors, errors.New("did not persist due to prior failure"))
	}
	return cir
}

func (r *Repository) CreateContent(content Content) ContentIngestionResult {
	content.Votes = nil
	content.Comments = nil
	cir := ContentIngestionResult{Content: content, Successful: true}
	cir = r.CheckSumId(cir)
	if r.EnsureNotAlreadyPresent(cir).Successful {
		cir = r.EnsurePosterId(r.EnsureId(cir))
		// // 3. sum it
		if r.PredictNotMalicious(r.PredictNotSpam(cir)).Successful {
			cir = r.Timeline(r.Tag(r.Persist(r.Sum(cir))))
		}
	}
	return cir
}

func (r *Repository) Timeline(cir ContentIngestionResult) ContentIngestionResult {
	arry, _ := json.Marshal([]string{cir.Content.Id.String()})
	timeKey := []byte(time.Now().Format(time.RFC3339Nano))
	var err error
	_ = r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Timeline"))
		// see if the id exists
		contentIds := b.Get(timeKey)
		if contentIds == nil {
			b.Put(timeKey, arry)
		} else {
			var ids []string
			json.Unmarshal(contentIds, &ids)
			ids = append(ids, cir.Content.Id.String())
			idsArry, _ := json.Marshal([]string{cir.Content.Id.String()})
			err = b.Put(timeKey, idsArry)
		}
		return err
	})
	return cir
}

func sum(c Content) []byte {
	arry := md5.Sum(c.AsJson())
	return arry[:]
}
func contentOnly(c Content) Content {
	return Content{
		Title: c.Title,
		Body:  c.Body,
		Url:   c.Url,
	}
}
func (r *Repository) GetContent(id string) Content {
	c := Content{}
	_ = r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Content"))
		uid, _ := uuid.FromString(id)
		json.Unmarshal(b.Get(uid.Bytes()), &c)
		return nil
	})
	return c
}

func (r *Repository) GetContentFull(id string) Content {
	content := Content{}
	_ = r.DB.View(func(tx *bolt.Tx) error {
		cb := tx.Bucket([]byte("Content"))
		uid, _ := uuid.FromString(id)
		json.Unmarshal(cb.Get(uid.Bytes()), &content)
		r.GetVotes(&content, tx)
		r.GetComments(&content, tx)
		content.Score = content.GetScore()
		return nil

	})
	return content
}

func (r *Repository) GetVotes(content *Content, tx *bolt.Tx) {
	c := tx.Bucket([]byte("Votes")).Cursor()
	prefix := content.Id.Bytes()
	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		vote := Vote{}
		_ = json.Unmarshal(v, &vote)
		_ = append(content.Votes, vote)
	}
}

func (r *Repository) GetComments(content *Content, tx *bolt.Tx) {
	c := tx.Bucket([]byte("Comments")).Cursor()
	prefix := content.Id.Bytes()
	for k, v := c.Seek(prefix); bytes.HasPrefix(k, prefix); k, v = c.Next() {
		comment := Comment{}
		_ = json.Unmarshal(v, &comment)
		_ = append(content.Comments, comment)
	}
}

func (r *Repository) CreateVote(vote Vote) (Vote, error) {
	vote.VoteId = uuid.NewV4()
	vote.Id = fmt.Sprintf("%s:%s", vote.ContentId.String(), vote.VoteId.String())
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Votes"))
		j, _ := json.Marshal(vote)
		err := b.Put([]byte(vote.Id), []byte(j))
		return err
	})

	return vote, err
}
