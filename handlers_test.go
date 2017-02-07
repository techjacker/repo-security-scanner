package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

// TODO:
// 1. make table driven
// 2. mock github API response
func TestGithubHandler(t *testing.T) {

	const (
		testPath = "/github"
		expected = "Push contains no offenses"
	)

	router := httprouter.New()
	router.POST(testPath, GithubHandler)

	// no query params -> pass 'nil' as the third parameter
	params := getFixture("test/fixtures/github_diff_response.json")
	r, err := http.NewRequest("POST", testPath, params)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	// Check the status code is what we expect.
	if w.Code != http.StatusOK {
		t.FailNow()
	}

	if strings.EqualFold(strings.TrimSpace(w.Body.String()), expected) != true {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), expected)
	}
}
