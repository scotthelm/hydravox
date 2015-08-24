package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"

	"github.com/GeertJohan/go.rice"
)

func ApiRootHandler(res http.ResponseWriter, req *http.Request) {
	r := ApiRoot{Meta: Meta{Name: "hydravox api", Licensing: "Creative Commons Attribution Share-Alike"}}
	json.NewEncoder(res).Encode(r)
}

func ContentCreateHandler(res http.ResponseWriter, req *http.Request) {
	r := Repository{server.DB, server.Config}
	content := Content{}
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&content)

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	content, err = r.CreateContent(content)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(res).Encode(content)
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
