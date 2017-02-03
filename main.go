package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
	df "github.com/techjacker/diffence"
)

const (
	rulesPath  = "rules/gitrob.json"
	serverPort = 8080
)

func main() {
	router := httprouter.New()
	router.POST("/github", GithubHandler)
	fmt.Printf("Server listening on port: %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router))
}

func GithubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// get rules
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := df.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read rule file: %s\n", err), 500)
		return
	}

	// decode github push event payload
	gitRes, err := decodeGithubJSON(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return
	}

	// analyse rules
	ruleResultsUnanalysed, err := runRules(
		rules,
		gitRes.getDiffURL(gitRes.Body.Commits[0].ID),
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return
	}

	// TODO: Notify recipients if fails checks
	ruleResults := checkRuleResults(ruleResultsUnanalysed, gitRes.Body.Commits[0].ID)
	if ruleResults.Pass == true {
		fmt.Fprintf(w, "Diff contains no offenses\n\n")
	} else {
		fmt.Fprintf(w, "TODO: email list of recipients to be notified of violations\n\n")
	}

}
