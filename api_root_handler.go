package main

import (
	"encoding/json"
	"net/http"
)

func ApiRootHandler(res http.ResponseWriter, req *http.Request) {
	r := ApiRoot{Meta: Meta{Name: "dendrite api", Licensing: "Creative Commons Attribution Share-Alike"}}
	json.NewEncoder(res).Encode(r)
}
