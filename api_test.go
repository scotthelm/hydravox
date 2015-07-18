package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
)

var server *DendriteServer
var router *mux.Router

func TestMain(m *testing.M) {
	db, err := bolt.Open("test2.db", 0660, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		fmt.Println("the error is: ", err)
	}
	fmt.Println("thist is db", db)
	server = new(DendriteServer)
	server.Port = ":7778"
	server.DBFile = "test2.db"
	server.DB = db
	router = server.NewRouter()
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	os.Exit(m.Run())
}

func TestRootHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("unable to create GET / request")
	}
	body, err := routeTest(res, req, t)
	fmt.Println(string(body))
	if !strings.Contains(string(body), "api") {
		t.Error("expected body to contain 'api', got", string(body))
	}
}

func TestCreateContent(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/content", nil)
	if err != nil {
		t.Fatal("unable to create POST /content request")
	}
	body, err := routeTest(res, req, t)
	if !strings.Contains(string(body), "id") {
		t.Error("expected body to contain 'id', got", string(body))
	}
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
