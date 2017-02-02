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
	router.GET("/", Index)
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

	// get full diff from github API
	diffURL := fmt.Sprintf(
		"%s/%s",
		gitRes.getDiffURLStem(),
		// TODO: loop through all the commits
		gitRes.Body.Commits[0].ID,
	)
	resp, err := getGithubDiff(diffURL)
	defer resp.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("%s", err), 500)
		return
	}

	// check body of diff
	res, err := df.CheckDiffs(resp.Body, rules)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading diff: %s\n", err), 500)
		return
	}

	// print errors
	// -> save to string when plugin in email
	dirty := false
	for k, v := range res {
		if len(v) > 0 {
			dirty = true
			fmt.Fprintf(w, "File %s violates %d rules:\n", k, len(v))
			for _, r := range v {
				fmt.Fprintf(w, "\n%s\n", r.String())
			}
		}
	}

	// Notify recipients if fails checks
	if dirty == false {
		fmt.Fprintf(w, "Diff contains no offenses\n\n")
	} else {
		fmt.Fprintf(w, "TODO: email list of recipients to be notified of violations\n\n")
	}
}
