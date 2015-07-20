package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

type DendriteServer struct {
	WebPort          string
	NotificationPort string
	DBFile           string
	NodeId           uuid.UUID
	LogFile          string
	DB               *bolt.DB
	configuration    Config
}

func (ds *DendriteServer) NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, r := range routes {
		var handler http.Handler
		handler = r.HandlerFunc
		handler = Logger(Recoverer(handler), r.Name)

		router.Methods(r.Method).Path(r.Pattern).Name(r.Name).Handler(handler)
	}

	return router
}

func NewDendriteServer() *DendriteServer {
	return new(DendriteServer).Configure("config.json")
}

func (ds *DendriteServer) Configure(path string) *DendriteServer {
	config := ReadConfiguration(path)
	config = WriteConfigurationIfNeeded(config, path)
	ds.WebPort = config.WebPort
	ds.NotificationPort = config.NotificationPort
	ds.DBFile = config.DBPath
	ds.configuration = config
	ds.NodeId = config.NodeId
	ds.LogFile = config.LogPath
	return ds
}

func ReadConfiguration(path string) Config {
	file, err := os.OpenFile(path, os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		log.Fatal("Cannot open configuration")
	}
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Cannot parse configuration: ", err)
	}
	return config
}

func WriteConfigurationIfNeeded(config Config, path string) Config {
	empty_uuid := uuid.UUID{}
	if uuid.Equal(config.NodeId, empty_uuid) {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			log.Fatal("Cannot open configuration for writing: ", err)
		}
		config.NodeId = uuid.NewV4()
		encoder := json.NewEncoder(file)
		if err := encoder.Encode(&config); err != nil {
			log.Fatal("Cannot write NodeId to configuration:", err)
		}
	}
	return config
}
