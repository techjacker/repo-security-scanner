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
	gitrobRules = "rules/gitrob.json"
	serverPort  = 8080
)

func main() {

	// get rules
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.ReadRulesFromFile(path.Join(path.Dir(cmd), gitrobRules))
	if err != nil {
		panic(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}

	// run server
	router := httprouter.New()
	router.GET("/healthz", HealthHandler)
	router.POST("/github", GithubHandler(diffValidator{
		rules:      rules,
		diffGetter: getGithubDiff,
	}))
	fmt.Printf("Server listening on port: %d", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
