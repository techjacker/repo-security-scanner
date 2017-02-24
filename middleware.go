package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"net/http"
)

// Adapter defines a middleware handler
type Adapter func(http.Handler) http.Handler

// Adapt chains a series of middleware handlers
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

// Middleware wraps a function and returns an http.Handler fn type
type Middleware func(http.Handler) http.Handler

// Authenticator is an interface for authorizing streams
type Authenticator interface {
	CheckMAC([]byte, []byte) (bool, error)
}

// GithubAuthenticator authorizes a payload with a shared secret
type GithubAuthenticator struct {
	secret []byte
}

// CheckMAC checks text matches a HMAC signature
func (g GithubAuthenticator) CheckMAC(body, expectedMAC []byte) (bool, error) {
	if len(g.secret) < 1 {
		return false, errors.New("secret not set")
	}
	mac := hmac.New(sha1.New, g.secret)
	mac.Write(body)
	actualMAC := []byte(hex.EncodeToString(mac.Sum(nil)))
	return hmac.Equal(expectedMAC, append([]byte("sha1="), actualMAC...)), nil
}

// AuthMiddleware authenticates a request body
func AuthMiddleware(ag Authenticator) Adapter {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(headerGithubEvt) != "push" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(msgIgnore))
				return
			}
			buf, err := ioutil.ReadAll(r.Body)
			rdr1 := ioutil.NopCloser(bytes.NewBuffer(buf))
			r.Body = rdr1
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return
			}
			authorized, err := ag.CheckMAC(
				buf,
				[]byte(r.Header.Get(headerGithubMAC)),
			)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
				return
			}
			if authorized != true {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
