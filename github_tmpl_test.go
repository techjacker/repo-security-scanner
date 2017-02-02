package main

import (
	"io"
	"testing"
)

func Test_decodeGithubJSON(t *testing.T) {
	type args struct {
		body io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *GithubResponse
		wantErr bool
	}{
		{
			name:    "Decodes github diff JSON response into struct",
			args:    args{getFixture("test/fixtures/github_diff_response.json")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeGithubJSON(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("decodeGithubJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			equals(t, got.Body.Commits[0].ID, "47797c0123bc0f5adfcae3d3467a2ed12e72b2cb")
			equals(t, got.Body.Commits[0].Added, []string{"ba.txt"})
			equals(t, got.Body.Repository.Name, "testgithubintegration")
			equals(t, got.Body.Repository.Owner.Name, "ukhomeoffice-bot-test")
			equals(t, got.Body.Repository.Owner.Email, nil)
			equals(t, got.Headers.X_github_event, "push")

			equals(t, got.getDiffURLStem(), "https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits")
		})
	}
}
