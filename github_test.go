package main

import (
	"fmt"
	"io"
	"testing"
)

func Test_decodeGithubJSON(t *testing.T) {
	type args struct {
		body io.Reader
	}
	type want struct {
		CommitsID            string
		CommitsAdded         []string
		RepositoryName       string
		RepositoryOwnerName  string
		RepositoryOwnerEmail interface{}
		HeadersXGithubEvent  string
		GithubAPIDiffURL     string
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "Decodes github diff JSON response into struct",
			args: args{getFixture("test/fixtures/github_event_push.json")},
			want: want{
				CommitsID:            "47797c0123bc0f5adfcae3d3467a2ed12e72b2cb",
				CommitsAdded:         []string{"ba.txt"},
				RepositoryName:       "testgithubintegration",
				RepositoryOwnerName:  "ukhomeoffice-bot-test",
				RepositoryOwnerEmail: nil,
				HeadersXGithubEvent:  "push",
				GithubAPIDiffURL:     "https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGithubResponse(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeGithubJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			equals(t, got.Body.Commits[0].ID, tt.want.CommitsID)
			equals(t, got.Body.Commits[0].Added, tt.want.CommitsAdded)
			equals(t, got.Body.Repository.Name, tt.want.RepositoryName)
			equals(t, got.Body.Repository.Owner.Name, tt.want.RepositoryOwnerName)
			equals(t, got.Body.Repository.Owner.Email, tt.want.RepositoryOwnerEmail)
			equals(t, got.Headers.XGithubEvent, tt.want.HeadersXGithubEvent)
			equals(t, got.getDiffURLStem(), tt.want.GithubAPIDiffURL)
			equals(t, got.getDiffURL(tt.want.CommitsID), fmt.Sprintf("%s/%s", tt.want.GithubAPIDiffURL, tt.want.CommitsID))
		})
	}
}
