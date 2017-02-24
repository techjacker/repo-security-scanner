package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

const (
	headerGithubMAC     = "X-Hub-Signature"
	headerGithubEvt     = "X-Github-Event"
	githubWebhookSecret = "GITHUB_WEBHOOKSECRET"
	gitrobRules         = "rules/gitrob.json"
	serverPort          = 8080
)

func getRules(rulesPath string) *[]diffence.Rule {
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.LoadRulesJSON(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		panic(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}
	return rules
}

func getEnvVar(name string) []byte {
	secretStr := os.Getenv(name)
	if len(secretStr) < 1 {
		panic(fmt.Sprintf("Env var:%s not set in environment", name))
	}
	return []byte(secretStr)
}

func main() {
	router := httprouter.New()
	router.Handler("GET", "/healthz", http.HandlerFunc(HealthHandler))
	router.Handler("POST", "/github", Adapt(
		GithubHandler(
			diffence.DiffChecker{getRules(gitrobRules)},
			diffGetterGithub{},
		),
		AuthMiddleware(GithubAuthenticator{getEnvVar(githubWebhookSecret)}),
	))

	fmt.Printf("Server listening on port: %d", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
