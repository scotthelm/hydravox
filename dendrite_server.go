package main

import (
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

type DendriteServer struct {
	Port   string
	DBFile string
	DB     *bolt.DB
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
