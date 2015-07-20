package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GeertJohan/go.rice"
	"github.com/boltdb/bolt"
)

var server *DendriteServer

func main() {
	server = new(DendriteServer).Configure("config.json")

	f, err := os.OpenFile(server.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open logfile")
	}

	defer f.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, f))

	db, err := bolt.Open(server.DBFile, 0660, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("the error is: ", err)
	}
	server.DB = db
	router := server.NewRouter()
	router.PathPrefix("/templates/").Handler(templatesHandler())

	log.Fatal(http.ListenAndServe(server.WebPort, router))
}

func templatesHandler() http.Handler {
	assets := rice.MustFindBox("templates")
	return http.StripPrefix("/templates/", http.FileServer(assets.HTTPBox()))
}
