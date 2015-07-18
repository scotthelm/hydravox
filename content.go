package main

import (
	"net/url"
	"time"

	"github.com/satori/go.uuid"
)

type Content struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Url         url.URL   `json:"string"`
	PosterId    uuid.UUID `json:"poster_id"`
	SubmittedAt time.Time `json:"submitted_at"`
	Votes       []Vote    `json:"votes"`
	Comments    []Comment `json:"comments"`
	Tags        []string  `json:"tags"`
}

type Vote struct {
	Type     string    `json:"type"`
	Id       uuid.UUID `json:"id"`
	Positive bool      `json:"positive"`
	PosterId uuid.UUID `json:"poster_id"`
}

type Comment struct {
	ContentId uuid.UUID `json:"content_id"`
	Id        uuid.UUID `json:"id"`
	Body      string    `json:"body"`
	PosterId  uuid.UUID `json:"poster_id"`
	ParentId  uuid.UUID `json:"response_to"`
}

type Timeline struct {
	PostedAt  time.Time `json:"posted_at"`
	ContentId uuid.UUID
}
