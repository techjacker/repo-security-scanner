package main

import (
	"fmt"
	"net/http"

	"github.com/techjacker/diffence"
)

const (
	msgHealthOk          = "ok"
	msgBadRequest        = "bad request"
	msgIgnore            = "Not a push evt; ignoring"
	msgNoViolationsFound = "Push contains no offenses"
	msgViolationFound    = "Push contains violations"
)

// GithubHandler is a github integration HTTP handler
func GithubHandler(dc diffence.Checker, dg DiffGetterHTTP, log Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// decode github push event payload
		gitRes := &GithubResponse{}
		err := DecodeJSON(r.Body, gitRes)
		if err != nil {
			http.Error(w, msgBadRequest, http.StatusBadRequest)
			return
		}
		// analyse all pushed commits
		results := diffence.Results{}
		for _, commit := range gitRes.Commits {
			resp, err := dg.Get(gitRes.getDiffURL(commit.ID))
			if err != nil {
				http.Error(w, msgBadRequest, http.StatusBadRequest)
				return
			}
			diffRes, err := dc.Check(resp.Body)
			resp.Body.Close()
			if err != nil {
				http.Error(w, msgBadRequest, http.StatusBadRequest)
				return
			}
			results = append(results, diffRes)
		}

		if results.Matches() < 1 {
			fmt.Fprintf(w, "%s\n", msgNoViolationsFound)
			return
		}

		for _, res := range results {
			if res.Matches() > 0 {
				log.Log(
					res.MatchedRules,
					gitRes.Repository.Owner.Name,
					gitRes.Repository.Name,
					gitRes.Compare,
				)
			}
		}
		fmt.Fprintf(w, "%s\n", msgViolationFound)
	})
}

// HealthHandler is an endpoint for healthchecks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msgHealthOk))
}
