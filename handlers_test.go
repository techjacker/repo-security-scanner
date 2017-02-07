package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

func getRules(t *testing.T) *[]diffence.Rule {
	// get rules
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		t.Fatal(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}
	return rules
}

// TODO:
// 1. make table driven
// 2. mock github API response
func TestGithubHandler(t *testing.T) {

	const (
		testPath = "/github"
		expected = "Push contains no offenses"
	)

	router := httprouter.New()
	router.POST(testPath, GithubHandler(&diffValidator{getRules(t)}))

	// no query params -> pass 'nil' as the third parameter
	params := getFixture("test/fixtures/github_event_push.json")
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
