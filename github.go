package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// NewGithubResponse marshalls a github event JSON payload into a struct
func NewGithubResponse(body io.Reader) (*GithubResponse, error) {
	var gitJSON GithubResponse
	dec := json.NewDecoder(body)
	err := dec.Decode(&gitJSON)
	return &gitJSON, err
}

// GithubResponse is for marshalling Github push event JSON payloads
type GithubResponse struct {
	Body struct {
		Commits []struct {
			Added []string `json:"added"`
			ID    string   `json:"id"`
		} `json:"commits"`
		Repository struct {
			Name  string `json:"name"`
			Owner struct {
				Email interface{} `json:"email"`
				Name  string      `json:"name"`
			} `json:"owner"`
		} `json:"repository"`
	} `json:"body"`
	// Installation struct {
	// 	ID int64 `json:"id"`
	// } `json:"installation"`
	Headers struct {
		XGithubEvent string `json:"x-github-event"`
	} `json:"headers"`
}

// getDiffURL returns the URL of the Github Diff API endpoint for a particular commit
func (g *GithubResponse) getDiffURL(commitID string) string {
	return fmt.Sprintf(
		"%s/%s",
		g.getDiffURLStem(),
		commitID,
	)
}

// getDiffURLStem/${CommitID}
func (g *GithubResponse) getDiffURLStem() string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits",
		g.Body.Repository.Owner.Name,
		g.Body.Repository.Name,
	)
}
