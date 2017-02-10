package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	msgHealthOk   = "ok"
	msgBadRequest = "bad request"
	msgSuccess    = "Push contains no offenses"
	msgFail       = "TODO: email list of recipients to be notified of violations"
)

// GithubHandler is a github integration HTTP handler
func GithubHandler(dv DiffValidator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		// decode github push event payload
		gitRes, err := NewGithubResponse(r.Body)
		if err != nil {
			http.Error(w, msgBadRequest, http.StatusBadRequest)
			return
		}

		// analyse all pushed commits
		results := []Stats{}
		allPassed := true
		for _, commit := range gitRes.Body.Commits {
			ruleResults, err := dv.Check(gitRes.getDiffURL(commit.ID))
			if err != nil {
				http.Error(w, msgBadRequest, http.StatusBadRequest)
				return
			}
			stats := dv.Stat(ruleResults, commit.ID)
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

func HealthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msgHealthOk))
}
