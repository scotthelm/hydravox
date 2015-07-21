package main

import (
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

type Repository struct {
	DB      *bolt.DB
	Buckets []*bolt.Bucket
}

func (r *Repository) InitializeBuckets() {
	r.DB.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{"Content", "Nodes", "Timeline"} {
			_, err := tx.CreateBucketIfNotExists([]byte(bucket))
			if err != nil {
				return fmt.Errorf("create bucket: %s", err)
			}
		}
		return nil
	})
}

func (r *Repository) CreateContent(content Content) (Content, error) {
	id := uuid.NewV4()
	content.Id = fmt.Sprintf("%s|%s", content.SubmittedAt, id)
	content.ContentId = id
	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Content"))
		value, err := json.Marshal(content)
		if err != nil {
			return err
		}
		err = b.Put([]byte(content.Id), []byte(value))
		return err
	})
	return content, err
}
