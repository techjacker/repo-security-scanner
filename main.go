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

const rulesPath = "rules/gitrob.json"
const diffURL = "https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits/f591c33a1b9500d0721b6664cfb6033d47a00793"
const serverPort = 8080

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)
	router.POST("/github", GithubHandler)
	fmt.Printf("Server listening on port: %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router))
}

// curl -X POST http://localhost:8080/github
func GithubHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := df.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		http.Error(w, fmt.Sprintf("Cannot read rule file: %s\n", err), 500)

		return
	}

	req, err := http.NewRequest("GET", diffURL, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get diff: %s\n", err), 500)
	}
	req.Header.Set("Accept", "application/vnd.github.VERSION.diff")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Could not get diff content from github API: %s\n", err), 500)
	}
	defer resp.Body.Close()
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
