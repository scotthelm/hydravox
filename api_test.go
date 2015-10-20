package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

// var server *Server
var router *mux.Router

func TestMain(m *testing.M) {
	server = new(Server).Configure("config_test.json")
	db, err := bolt.Open(server.DBFile, 0660, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("the error is: ", err)
	}
	server.DB = db
	router = server.NewRouter()
	os.Exit(m.Run())
}

func TestRootHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("unable to create GET / request")
	}
	body, err := routeTest(res, req, t)
	if !strings.Contains(string(body), "html") {
		t.Error("expected body to contain 'html', got", string(body))
	}
}

func TestApiRootHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api", nil)
	if err != nil {
		t.Fatal("unable to create GET /api request")
	}
	body, err := routeTest(res, req, t)
	if !strings.Contains(string(body), "api") {
		t.Error("expected body to contain 'api', got", string(body))
	}
}

func TestCreateContent(t *testing.T) {
	url, err := url.Parse("https://google.com")
	content := Content{Title: "Title test api call", Body: "This is a test", Url: *url}
	json, err := json.Marshal(content)
	reader := bytes.NewReader([]byte(json))
	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/api/content", reader)
	if err != nil {
		t.Fatal("unable to create POST /api/content request")
	}
	body, err := routeTest(res, req, t)
	if !strings.Contains(string(body), "id") {
		t.Error("expected body to contain 'id', got", string(body))
	}
}

func TestCreateVote(t *testing.T) {
	_, result := createTestContent()
	voteId := uuid.NewV4()
	v := Vote{
		Type:      "Content",
		ContentId: result.Content.Id,
		VoteId:    voteId,
		PosterId:  server.Config.NodeId,
		Positive:  true,
		Id:        fmt.Sprintf("%s:%s", result.Content.Id.String(), voteId.String()),
	}
	fmt.Println(v)
}

func routeTest(res *httptest.ResponseRecorder, req *http.Request, t *testing.T) ([]byte, error) {
	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Error("expected 200, got", res.Code)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("unable to read body")
	}
	return body, err
}

func TestScore(t *testing.T) {
	c := new(Content)
	c.Id = uuid.NewV4()
	subAt := time.Now()
	subAt = subAt.Add(-time.Hour * 1)
	c.SubmittedAt = subAt
	c.addVotes(200, 100)
	if c.GetScore() != 70 {
		t.Error("Expected 70, got ", c.GetScore())
	}
}

func TestUpVotes(t *testing.T) {
	c := new(Content)
	c.Id = uuid.NewV4()
	c.addVotes(200, 100)
	if c.UpVotes() != 200 {
		t.Error("Expected 200, got ", c.UpVotes())
	}
}

func TestDownVotes(t *testing.T) {
	c := new(Content)
	c.Id = uuid.NewV4()
	c.addVotes(200, 100)
	if c.DownVotes() != 100 {
		t.Error("Expected 100, got ", c.DownVotes())
	}
}

func (c *Content) addVotes(positive int, negative int) {
	votes := make([]Vote, positive+negative)
	for v := 0; v < positive+negative; v++ {
		votes[v].Positive = v > negative-1
	}
	c.Votes = votes
}
