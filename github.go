package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
		X_github_event string `json:"x-github-event"`
	} `json:"headers"`
}

// GithubAPIEvent reveals information from a Github API event payload
type GithubAPIEvent interface {
	getDiffURL(string) string
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

func decodeGithubJSON(body io.Reader) (*GithubResponse, error) {
	var gitJSON GithubResponse
	dec := json.NewDecoder(body)
	err := dec.Decode(&gitJSON)
	return &gitJSON, err
}

func getGithubDiff(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Could not get diff: %s\n", err)
	}
	req.Header.Set("Accept", "application/vnd.github.VERSION.diff")
	return http.DefaultClient.Do(req)
}
