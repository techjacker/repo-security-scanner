package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
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
func GithubHandler(dc diffence.Checker, dg DiffGetterHTTP) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		if r.Header.Get(headerGithubEvt) != "push" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(msgIgnore))
			return
		}
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
		// TODO: Notify recipients if fails checks
		// stringify results vs pass to logger
		if results.Matches() > 0 {
			fmt.Fprintf(w, "%s\n", msgFail)
		} else {
			fmt.Fprintf(w, "%s\n", msgSuccess)
		}
	}
}

// HealthHandler is an endpoint for healthchecks
func HealthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msgHealthOk))
}
