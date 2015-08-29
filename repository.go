package main

import (
	"crypto/md5"
	"errors"
	"fmt"

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
		for _, bucket := range []string{"Content", "Nodes", "Timeline", "Spam", "Quarantine", "Tags", "Sums"} {
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
		contentId := b.Get([]byte(sum(cir.Content)))
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
func (r *Repository) ContentSummer(cir ContentIngestionResult) ContentIngestionResult {
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Sums"))
		err := b.Put([]byte(sum(cir.Content)), cir.Content.Id.Bytes())
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
	cir := ContentIngestionResult{Content: content, Successful: true}

	cir = r.EnsurePosterId(r.EnsureId(cir))
	// // 3. sum it
	if r.PredictNotMalicious(r.PredictNotSpam(r.EnsureNotAlreadyPresent(cir))).Successful {
		cir = r.Timeline(r.Tag(r.Persist(cir)))
	}
	return cir
}

func (r *Repository) Timeline(cir ContentIngestionResult) ContentIngestionResult {
	return cir
}

func sum(c Content) string {
	hash := fmt.Sprintf("%x", md5.Sum(c.AsJson()))
	return hash
}
