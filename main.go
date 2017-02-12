package main

import (
	"fmt"
	"net/http"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

func getRules(rulesPath string) *[]diffence.Rule {
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		panic(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}
	return rules
}

const (
	gitrobRules = "rules/gitrob.json"
	serverPort  = 8080
)

func main() {
	router := httprouter.New()
	router.GET("/healthz", HealthHandler)
	router.POST("/github", GithubHandler(
		diffChecker{getRules(gitrobRules)},
		diffGetterGithub{},
	))
	fmt.Printf("Server listening on port: %d", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
