package main

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/satori/go.uuid"
)

const SCORING_FACTOR = float64(170001)

// Id is a concatentation of SubmittedAt and ContentId
type Content struct {
	Id            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Body          string    `json:"body"`
	Url           url.URL   `json:"string"`
	PosterId      uuid.UUID `json:"poster_id"`
	SubmittedAt   time.Time `json:"submitted_at"`
	Tags          []string  `json:"tags"`
	Votes         []Vote    `json:"votes"`
	Comments      []Comment `json:"comments"`
	SpamScore     float64   `json:"spam_score"`
	IsSpam        bool      `json:"is_spam"`
	IsQuarantined bool      `json:"is_quarantined"`
	Score         int
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

func (c *Content) GetScore() int {
	var ups float64
	var downs float64
	submitted_at := float64(c.SubmittedAt.Unix())
	now := float64(time.Now().Unix())
	diff := submitted_at + ((now - submitted_at) * SCORING_FACTOR)
	for _, v := range c.Votes {
		if v.Positive {
			ups = ups + float64(1)
		} else {
			downs = downs + float64(1)
		}
	}
	// this is a parabolic function approaching 0
	// (((ups - downs) * submitted_at) / ((ups - downs) * now)) * (ups - downs)
	return int(((ups - downs) * submitted_at) / ((ups - downs) * diff) * (ups - downs))
}

func (c *Content) UpVotes() int {
	voteCounter := 0
	for v := 0; v < len(c.Votes); v++ {
		if c.Votes[v].Positive {
			voteCounter++
		}
	}
	return voteCounter
}

func (c *Content) DownVotes() int {
	voteCounter := 0
	for v := 0; v < len(c.Votes); v++ {
		if !c.Votes[v].Positive {
			voteCounter++
		}
	}
	return voteCounter
}

func (c *Content) AsJson() []byte {
	value, _ := json.Marshal(c)
	return value
}

// Id is a concatenation of ContentId and CommentId
type Comment struct {
	Id        string    `json:"id"`
	ContentId uuid.UUID `json:"content_id"`
	CommentId uuid.UUID `json:"comment_id"`
	Body      string    `json:"body"`
	PosterId  uuid.UUID `json:"poster_id"`
	ParentId  string    `json:"response_to"`
	Comments  []Comment `json:"children"`
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
