package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{"Index", "GET", "/", Index},
	Route{"ApiIndex", "GET", "/api", ApiRootHandler},
	Route{"ContentCreate", "POST", "/api/content", ContentCreateHandler},
	Route{"CreateComment", "POST", "/api/content/{content_id}/comments", CommentCreateHandler},
	Route{"CreateVote", "POST", "/api/content/{content_id}/votes", VoteCreateHandler},
	Route{"GetContent", "GET", "/api/content/{content_id}", ContentGetHandler},
}
