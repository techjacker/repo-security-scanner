package main

import (
	"fmt"
	"net/http"
)

//go:generate gojson -name githubResponseFull -input test/fixtures/github_diff_response.json -o github.go

type GithubResponse struct {
	Body struct {
		Commits []struct {
			Added    []string      `json:"added"`
			ID       string        `json:"id"`
			Modified []interface{} `json:"modified"`
			// Removed  []interface{} `json:"removed"`
			// "https://github.com/ukhomeoffice-bot-test/testgithubintegration/commit/47797c0123bc0f5adfcae3d3467a2ed12e72b2cb",
			// curl -s "https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits/47797c0123bc0f5adfcae3d3467a2ed12e72b2cb" -H "Accept: application/vnd.github.VERSION.diff"
			// URL string `json:"url"`
		} `json:"commits"`
	} `json:"body"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Email interface{} `json:"email"`
			Name  string      `json:"name"`
		} `json:"owner"`
	} `json:"repository"`
	// Installation struct {
	// 	ID int64 `json:"id"`
	// } `json:"installation"`
	Headers struct {
		X_github_event string `json:"x-github-event"`
	} `json:"headers"`
}

//   .then(() => req.body.commits)
//       sha: commit.id,
// const installationId = req.body.installation.id
// const repo = req.body.repository.name
// const owner = req.body.repository.owner.name

func getGithubDiff(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &http.Response{}, fmt.Errorf("Could not get diff: %s\n", err)
	}
	req.Header.Set("Accept", "application/vnd.github.VERSION.diff")
	return http.DefaultClient.Do(req)
}
