package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

const (
	sharedSecret   = "blah"
	testPath       = "/auth"
	testHandlerMsg = "test handler message body"
)

func TestAuthMiddleware(t *testing.T) {
	type args struct {
		headerGithubEvt string
		headerGithubMAC string
		body            string
	}
	// authTestHandler returns a http.Handler for testing http middleware
	authTestHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testHandlerMsg))
	})

	tests := []struct {
		name           string
		args           args
		wantStatusCode int
		wantResBody    string
	}{
		{
			name: "Not a push event",
			args: args{
				headerGithubEvt: "not-push",
				headerGithubMAC: "sha1=033fa0645f29a6f8a4decd8e8aee8a9e341151cc",
				body:            "hello world",
			},
			wantStatusCode: http.StatusOK,
			wantResBody:    msgIgnore,
		},
		{
			name: "Push event with valid signature",
			args: args{
				headerGithubEvt: "push",
				headerGithubMAC: "sha1=033fa0645f29a6f8a4decd8e8aee8a9e341151cc",
				body:            "hello world",
			},
			wantStatusCode: http.StatusOK,
			wantResBody:    testHandlerMsg,
		},
		{
			name: "Push event with invalid signature",
			args: args{
				headerGithubEvt: "push",
				headerGithubMAC: "sha1=0000000000000000000000000000000000000000",
				body:            "hello world",
			},
			wantStatusCode: http.StatusUnauthorized,
			wantResBody:    http.StatusText(http.StatusUnauthorized),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			router := httprouter.New()
			auth := AuthMiddleware(GithubAuthenticator{[]byte(sharedSecret)})
			router.Handler("POST", testPath, Adapt(authTestHandler, auth))

			r, err := http.NewRequest("POST", testPath, strings.NewReader(tt.args.body))
			r.Header.Set(headerGithubEvt, tt.args.headerGithubEvt)
			r.Header.Set(headerGithubMAC, tt.args.headerGithubMAC)
			if err != nil {
				t.Error(err)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, r)

			if w.Code != tt.wantStatusCode {
				t.Fatalf("TestAuthHandler returned unexpected status code: got %v want %v",
					w.Code, tt.wantStatusCode)
			}
			if strings.EqualFold(strings.TrimSpace(w.Body.String()), tt.wantResBody) != true {
				t.Fatalf("TestAuthHandler returned unexpected body: got %v want %v",
					w.Body.String(), tt.wantResBody)
			}
		})
	}

}
