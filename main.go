package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

const (
	serverPort = 8080
)

func main() {

	// get rules
	_, cmd, _, _ := runtime.Caller(0)
	rules, err := diffence.ReadRulesFromFile(path.Join(path.Dir(cmd), rulesPath))
	if err != nil {
		panic(fmt.Sprintf("Cannot read rule file: %s\n", err))
	}

	// run server
	router := httprouter.New()
	router.POST("/github", GithubHandler(diffValidator{rules}))
	fmt.Printf("Server listening on port: %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router))
}
