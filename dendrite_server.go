package main

import (
	"github.com/gorilla/mux"
	"net/http"
)

type DendriteServer struct {
	Port string
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
