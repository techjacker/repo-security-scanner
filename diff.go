package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/techjacker/diffence"
)

// DiffChecker checks diffs for rule violations
type DiffChecker interface {
	Check(io.ReadCloser) (diffence.Results, error)
}

type diffChecker struct {
	rules *[]diffence.Rule
}

func (d diffChecker) Check(in io.ReadCloser) (diffence.Results, error) {
	defer in.Close()
	return diffence.CheckDiffs(in, d.rules)
}

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

// DiffStats analyses a series of diff rule results
type DiffStats struct {
	Pass bool
	Info string
}

func statDiffResults(ruleResults diffence.Results, commitID string) DiffStats {

	dirty := false
	info := fmt.Sprintf("Commit ID: %s\n", commitID)

	for k, v := range ruleResults {
		if len(v) > 0 {
			dirty = true
			info += fmt.Sprintf("File %s violates %d rules:\n", k, len(v))
			for _, r := range v {
				info += fmt.Sprintf("\n%s\n", r.String())
			}
		}
	}

	return DiffStats{
		Pass: !dirty,
		Info: info,
	}
}
