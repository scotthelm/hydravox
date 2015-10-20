package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/boltdb/bolt"
)

// repository should have a db
func TestRepositoryDatabase(t *testing.T) {
	_ = Repository{server.DB, server.Config}
}

// repository should have a set of buckets
func TestRepositoryBuckets(t *testing.T) {
	r := Repository{server.DB, server.Config}
	r.InitializeBuckets()
}

// repository should be able to create content
func TestRepositoryCreateContent(t *testing.T) {
	r := Repository{server.DB, server.Config}
	result := r.CreateContent(
		Content{
			Title:       "Test",
			Body:        "This is a test.",
			PosterId:    server.NodeId,
			SubmittedAt: time.Now(),
			Tags:        []string{"Test"},
		})
	if result.Successful == false {
		t.Error("Error creating content", result.Errors)
	}
	checkResults(r, []byte("Content"))
	checkResults(r, []byte("Sums"))
	checkResults(r, []byte("Timeline"))
	checkResults(r, []byte("Tags"))
}

func TestRepoGetContent(t *testing.T) {
	r, result := createTestContent()
	fmt.Println("***************************", result.Content.Id)
	content := r.GetContentFull(result.Content.Id.String())
	fmt.Println("%%%%%%%%%%%%%%%%%%%%%%%%%%%", content.Id)
	if content.Id != result.Content.Id {
		t.Error("Error : content get not correct.")
	}
}

func createTestContent() (Repository, ContentIngestionResult) {
	r := Repository{server.DB, server.Config}
	result := r.CreateContent(
		Content{
			Title:       "Test",
			Body:        "This is a test.",
			PosterId:    server.NodeId,
			SubmittedAt: time.Now(),
			Tags:        []string{"Test"},
		})
	return r, result
}

func checkResults(r Repository, bucketName []byte) {
	_ = r.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("value for %s:%s = %s\n----------------\n", bucketName, k, v)
			return nil
		})
		return nil
	})
}
