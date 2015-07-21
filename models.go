package main

import (
	"net/url"
	"time"

	"github.com/satori/go.uuid"
)

// Id is a concatentation of SubmittedAt and ContentId
type Content struct {
	Id          string    `json:"id"`
	ContentId   uuid.UUID `json:"content_id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Url         url.URL   `json:"string"`
	PosterId    uuid.UUID `json:"poster_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	Tags        []string  `json:"tags"`
}

// Id is a concatenation of ContentId and VoteId
type Vote struct {
	Type      string    `json:"type"`
	Id        string    `json:"id"`
	Positive  bool      `json:"positive"`
	PosterId  uuid.UUID `json:"poster_id"`
	ContentId uuid.UUID `json:"content_id"`
	VoteId    uuid.UUID `json:"vote_id"`
}

// Id is a concatenation of ContentId and CommentId
type Comment struct {
	Id        string    `json:"id"`
	ContentId uuid.UUID `json:"content_id"`
	CommentId uuid.UUID `json:"comment_id"`
	Body      string    `json:"body"`
	PosterId  uuid.UUID `json:"poster_id"`
	ParentId  string    `json:"response_to"`
}

// Not sure this type is necessary, given that time is a component of content id
type Timeline struct {
	PostedAt  time.Time `json:"posted_at"`
	ContentId uuid.UUID `json:"content_id"`
}

type Node struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"string"`
	Score float64   `json:"score"`
}

type ApiRoot struct {
	Meta `json:"meta"`
}

type Meta struct {
	Name      string `json:"name"`
	Licensing string `json:"licensing"`
}
