package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	serverPort = 8080
)

func main() {
	router := httprouter.New()
	router.POST("/github", GithubHandler)
	fmt.Printf("Server listening on port: %d", serverPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", serverPort), router))
}
