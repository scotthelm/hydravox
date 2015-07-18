package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boltdb/bolt"
)

func ApiRootHandler(res http.ResponseWriter, req *http.Request) {
	r := ApiRoot{Meta: Meta{Name: "dendrite api", Licensing: "Creative Commons Attribution Share-Alike"}}
	json.NewEncoder(res).Encode(r)
}

func ContentCreateHandler(res http.ResponseWriter, req *http.Request) {
	server.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	res.Write([]byte("{\"id\" : \"1\"}"))
}
