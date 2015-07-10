package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootHandler(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("unable to create GET / request")
	}
	ds := DendriteServer{":7777"}
	router := ds.NewRouter()
	router.ServeHTTP(res, req)
	if res.Code != 200 {
		t.Error("expected 200, got", res.Code)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("unable to read body")
	}
	fmt.Println(string(body))
	if !strings.Contains(string(body), "api") {
		t.Error("expected body to contain 'api', got", string(body))
	}
}
