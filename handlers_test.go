package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/techjacker/diffence"
)

type testDiffGetter struct {
	fixt string
}

func (t testDiffGetter) Get(_ string) (*http.Response, error) {
	return &http.Response{
		Body: ioutil.NopCloser(getFixture(t.fixt)),
	}, nil
}

type testLogger struct {
	log *logrus.Logger
}

func (l testLogger) Log(v ...interface{}) {
	return
}

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
			name: "Incorrect JSON payload returns 4xx status code",
			args: args{
				githubPayloadPath: "test/fixtures/nonsense.json",
				rulesPath:         gitrobRules,
				diffPath:          "",
			},
			wantStatusCode: http.StatusBadRequest,
			wantResBody:    msgBadRequest,
		},
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
			name: "1 offense in diff",
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
			router.Handler("POST", testPath, GithubHandler(
				diffence.DiffChecker{Rules: getTestRules(t, tt.args.rulesPath)},
				testDiffGetter{tt.args.diffPath},
				testLogger{},
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
			wantResBody:    msgHealthOk,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := httprouter.New()
			router.Handler("GET", testPath, http.HandlerFunc(HealthHandler))

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
