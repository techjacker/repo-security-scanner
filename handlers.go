package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	headerGithubEvt = "x-github-event"
	msgHealthOk     = "ok"
	msgBadRequest   = "bad request"
	msgIgnore       = "Not a push evt; ignoring"
	msgSuccess      = "Push contains no offenses"
	msgFail         = "TODO: email list of recipients to be notified of violations"
)

// GithubHandler is a github integration HTTP handler
func GithubHandler(dc DiffChecker, dg DiffGetterHTTP) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		if r.Header.Get(headerGithubEvt) != "push" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(msgIgnore))
			return
		}

		// decode github push event payload
		gitRes, err := NewGithubResponse(r.Body)
		if err != nil {
			http.Error(w, msgBadRequest, http.StatusBadRequest)
			return
		}

		// analyse all pushed commits
		results := []DiffStats{}
		allPassed := true
		for _, commit := range gitRes.Commits {
			resp, err := dg.Get(gitRes.getDiffURL(commit.ID))
			if err != nil {
				http.Error(w, msgBadRequest, http.StatusBadRequest)
				return
			}
			diffResults, err := dc.Check(resp.Body)
			stats := statDiffResults(diffResults, commit.ID)
			if stats.Pass != true {
				allPassed = false
			}
			results = append(results, stats)
		}

		// TODO: Notify recipients if fails checks
		// stringify results vs pass to logger
		if allPassed == true {
			fmt.Fprintf(w, "%s\n", msgSuccess)
		} else {
			fmt.Fprintf(w, "%s\n", msgFail)
		}
	}
}

// HealthHandler is an endpoint for healthchecks
func HealthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msgHealthOk))
}
