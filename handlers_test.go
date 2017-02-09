package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

func TestGithubHandler(t *testing.T) {
	const (
		testPath = "/github"
	)

	type args struct {
		githubPayloadPath string
		rulesPath         string
		diffPath          string
	}
	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResBody    string
	}{
		{
			name: "No offenses in diff",
			args: args{
				githubPayloadPath: "test/fixtures/github_event_push.json",
				rulesPath:         gitrobRules,
				diffPath:          "test/fixtures/no_offenses.diff",
			},
			wantStatusCode: http.StatusOK,
			wantResBody:    msgSuccess,
		},
		{
			name: "No offenses in diff",
			args: args{
				githubPayloadPath: "test/fixtures/github_event_push.json",
				rulesPath:         gitrobRules,
				diffPath:          "test/fixtures/offenses_x1.diff",
			},
			wantStatusCode: http.StatusOK,
			wantResBody:    msgFail,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := httprouter.New()
			router.POST(testPath, GithubHandler(
				NewTestDiffValidator(getDiffenceRules(t, tt.args.rulesPath), tt.args.diffPath),
			))

			params := getFixture(tt.args.githubPayloadPath)
			r, err := http.NewRequest("POST", testPath, params)
			if err != nil {
				t.Error(err)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Fatalf("Githubhandler returned unexpected status code: got %v want %v",
					w.Code, tt.wantStatusCode)
			}
			if strings.EqualFold(strings.TrimSpace(w.Body.String()), tt.wantResBody) != true {
				t.Fatalf("Githubhandler returned unexpected body: got %v want %v",
					w.Body.String(), tt.wantResBody)
			}
		})
	}

}

func TestHealthHandler(t *testing.T) {
	const (
		testPath = "/healthz"
	)

	tests := []struct {
		name           string
		wantStatusCode int
		wantResBody    string
	}{
		{
			name:           "Healthcheck handler",
			wantStatusCode: http.StatusOK,
			wantResBody:    healthMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := httprouter.New()
			router.GET(testPath, HealthHandler)

			r, err := http.NewRequest("GET", testPath, nil)
			if err != nil {
				t.Error(err)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Fatalf("Healthhandler returned unexpected status code: got %v want %v",
					w.Code, tt.wantStatusCode)
			}
			if strings.EqualFold(strings.TrimSpace(w.Body.String()), tt.wantResBody) != true {
				t.Fatalf("Healthhandler returned unexpected body: got %v want %v",
					w.Body.String(), tt.wantResBody)
			}
		})
	}

}

type TestDiffValidator struct {
	diffValidator
}

func NewTestDiffValidator(rules *[]diffence.Rule, fixt string) *TestDiffValidator {
	tdv := TestDiffValidator{
		diffValidator{
			rules:      rules,
			diffGetter: testHTTPGetterFactory(fixt),
		},
	}
	return &tdv
}

func testHTTPGetterFactory(fixt string) httpGetter {
	return func(url string) (*http.Response, error) {
		return &http.Response{
			Body: ioutil.NopCloser(getFixture(fixt)),
		}, nil
	}
}
