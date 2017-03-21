package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/Sirupsen/logrus"
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

func getRequiredEnvVar(name string) []byte {
	val, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("Env var:%s not set in environment", name))
	}
	return []byte(val)
}

func getLogger() Log {
	esURL, ok := os.LookupEnv("ELASTICSEARCH_URL")
	if ok {
		fmt.Printf("Logging to elasticsearch: %s\n", esURL)
		logger, err := NewESLogger(esURL, "githubintegration")
		if err != nil {
			panic(err)
		}
		return logger
	}
	return Logger{
		log: logrus.New(),
	}
}

func main() {
	router := httprouter.New()
	router.Handler("GET", "/healthz", http.HandlerFunc(HealthHandler))
	router.Handler("POST", "/github", Adapt(
		GithubHandler(
			diffence.DiffChecker{Rules: getRules(gitrobRules)},
			diffGetterGithub{},
			getLogger(),
		),
		AuthMiddleware(GithubAuthenticator{getRequiredEnvVar(githubWebhookSecret)}),
	))

	fmt.Printf("Server listening on port: %d\n", serverPort)
	http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router)
}
