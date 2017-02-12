package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

// Valid exposes a method to validate a type
type Valid interface {
	OK() error
}

// DecodeJSON marshalls JSON into a struct
func DecodeJSON(r io.Reader, v interface{}) error {
	err := json.NewDecoder(r).Decode(v)
	if err != nil {
		return err
	}
	obj, ok := v.(Valid)
	if !ok {
		return nil // no OK method
	}
	err = obj.OK()
	if err != nil {
		return err
	}
	return nil
}

// GithubResponse is for marshalling Github push event JSON payloads
type GithubResponse struct {
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
}

// OK validates a marshalled struct
func (g *GithubResponse) OK() error {
	if len(g.Commits) < 1 {
		return errors.New("empty payload")
	}
	for _, c := range g.Commits {
		if len(c.ID) < 1 {
			return errors.New("missing commit ID")
		}
	}
	if len(g.Repository.Name) < 1 {
		return errors.New("missing repository name")
	}
	return nil
}

// getDiffURL returns the URL of the Github Diff API endpoint for a particular commit
func (g GithubResponse) getDiffURL(commitID string) string {
	return fmt.Sprintf(
		"%s/%s",
		g.getDiffURLStem(),
		commitID,
	)
}

func (g *GithubResponse) getDiffURLStem() string {
	return fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits",
		g.Repository.Owner.Name,
		g.Repository.Name,
	)
}
