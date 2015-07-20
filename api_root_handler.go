package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/GeertJohan/go.rice"
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

func Index(w http.ResponseWriter, r *http.Request) {
	tmpls, err := rice.FindBox("templates")
	if err != nil {
		log.Fatal(err)
	}
	t, err := tmpls.String("index.html")
	if err != nil {
		log.Panic(err)
	}
	tmpl, err := template.New("Index").Parse(t)
	if err != nil {
		log.Panic(err)
	}
	tmpl.Execute(w, nil)
}
