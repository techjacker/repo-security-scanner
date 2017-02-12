package main

import (
	"fmt"
	"net/http"
)

// DiffGetterHTTP gets diff text for a commit
type DiffGetterHTTP interface {
	Get(string) (*http.Response, error)
}

type diffGetterGithub struct{}

func (g diffGetterGithub) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("could not get diff: %s", err)
	}
	req.Header.Set("Accept", "application/vnd.github.VERSION.diff")
	return http.DefaultClient.Do(req)
}
