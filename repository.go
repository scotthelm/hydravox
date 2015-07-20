package main

import (
	"fmt"

	"github.com/boltdb/bolt"
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
	return content, nil
}
