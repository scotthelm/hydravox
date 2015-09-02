package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
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
		} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
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
