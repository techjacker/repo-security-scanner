package main

import (
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

const (
	rulesPath = "rules/gitrob.json"
)

func GithubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get rules
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read rule file: %s\n", err), 500)
		return
	}

	// decode github push event payload
	gitRes, err := NewGithubResponse(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return
	}

	// analyse all pushed commits
	results := []RuleResults{}
	allPassed := true
	for _, commit := range gitRes.Body.Commits {
		ruleResultsUnanalysed, err := runRules(
			rules,
			gitRes.getDiffURL(commit.ID),
		)
		if err != nil {
			http.Error(w, fmt.Sprintf("%s", err), 500)
			return
		}
		ruleResultsAnalysed := checkRuleResults(
			ruleResultsUnanalysed,
			commit.ID,
		)
		if ruleResultsAnalysed.Pass != true {
			allPassed = false
		}
		results = append(results, ruleResultsAnalysed)
	}

	// TODO: Notify recipients if fails checks
	// stringify results vs pass to logger
	if allPassed == true {
		fmt.Fprintf(w, "Push contains no offenses\n\n")
	} else {
		fmt.Fprintf(w, "TODO: email list of recipients to be notified of violations\n\n")
	}

}
