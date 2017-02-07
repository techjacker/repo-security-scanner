package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	rulesPath = "rules/gitrob.json"
)

// GithubHandler is a github integration HTTP handler
func GithubHandler(dv DiffValidator) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		// func GithubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		// decode github push event payload
		gitRes, err := NewGithubResponse(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
			return
		}

		// analyse all pushed commits
		results := []Stats{}
		allPassed := true
		for _, commit := range gitRes.Body.Commits {
			ruleResults, err := dv.Check(gitRes.getDiffURL(commit.ID))
			if err != nil {
				http.Error(w, fmt.Sprintf("%s", err), 500)
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
			fmt.Fprintf(w, "Push contains no offenses\n\n")
		} else {
			fmt.Fprintf(w, "TODO: email list of recipients to be notified of violations\n\n")
		}
	}
}
