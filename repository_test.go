package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

// repository should have a db
func TestRepositoryDatabase(t *testing.T) {
	_ = Repository{server.DB, make([]*bolt.Bucket, 3)}
}

// repository should have a set of buckets
func TestRepositoryBuckets(t *testing.T) {
	r := Repository{server.DB, make([]*bolt.Bucket, 3)}
	r.InitializeBuckets()
	fmt.Println(r)
}

// repository should be able to create content
func TestRepositoryCreateContent(t *testing.T) {
	r := Repository{server.DB, make([]*bolt.Bucket, 3)}
	_, err := r.CreateContent(
		Content{
			Id:          uuid.NewV4(),
			Title:       "Test",
			Body:        "This is a test.",
			PosterId:    server.NodeId,
			SubmittedAt: time.Now(),
			Tags:        []string{"Test"},
		})
	if err != nil {
		t.Error("Error creating content", err)
	}
}
