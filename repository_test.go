package main

import (
	"fmt"
	"testing"
	"time"
)

// repository should have a db
func TestRepositoryDatabase(t *testing.T) {
	_ = Repository{server.DB, server.Config}
}

// repository should have a set of buckets
func TestRepositoryBuckets(t *testing.T) {
	r := Repository{server.DB, server.Config}
	r.InitializeBuckets()
	fmt.Println(r)
}

// repository should be able to create content
func TestRepositoryCreateContent(t *testing.T) {
	r := Repository{server.DB, server.Config}
	content, err := r.CreateContent(
		Content{
			Title:       "Test",
			Body:        "This is a test.",
			PosterId:    server.NodeId,
			SubmittedAt: time.Now(),
			Tags:        []string{"Test"},
		})
	if err != nil {
		t.Error("Error creating content", err)
	}
	fmt.Println("test content: ", content)
}
