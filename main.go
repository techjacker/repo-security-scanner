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

func getRules(rulesPath string) *[]diffence.Rule {
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.LoadRulesJSON(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		panic(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}
	return rules
}

func main() {
	router := httprouter.New()
	router.GET("/healthz", HealthHandler)
	router.POST("/github", GithubHandler(
		diffence.DiffChecker{getRules(gitrobRules)},
		diffGetterGithub{},
	))
	fmt.Printf("Server listening on port: %d", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
