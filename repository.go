package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/satori/go.uuid"
)

type Repository struct {
	DB     *bolt.DB
	Config Config
}

type ContentInjestionResult struct {
	Content       Content
	Error         []error
	Task          string
	TaskCompleted bool
}

type contentActor func(*ContentInjestionResult) *ContentInjestionResult

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
	content.Id = uuid.NewV4()
	if content.PosterId == uuid.Nil {
		content.PosterId = r.Config.NodeId
	}

	// we need to break this apart soon, but for now:
	// 1. set id
	// 2. set posterId if not present
	// 3. sum it
	// 4. see if it needs to be quarantined
	// 5. if not, then see if it's spam
	// 6. persist the content
	// 7. tag it
	// 8. score it
	// 9. timeline it

	err := r.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Content"))
		value, err := json.Marshal(content)
		if err != nil {
			return err
		}
		err = b.Put([]byte(content.Id.String()), []byte(value))
		return err
	})
	return content, err
}

func merge(channels ...<-chan ContentInjestionResult) <-chan ContentInjestionResult {
	var wg sync.WaitGroup
	out := make(chan ContentInjestionResult)
	// copy from the incoming channels to the outgoing channel
	copier := func(c <-chan ContentInjestionResult) {
		for cir := range c {
			out <- cir
		}
		wg.Done()
	}
	wg.Add(len(channels))
	// for each incoming channel, get it in the copier
	for _, c := range channels {
		go copier(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
