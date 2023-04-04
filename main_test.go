package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckAPItoken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header["Authorization"]
		if len(ah) == 1 && ah[0] == "Token validToken" {
			w.WriteHeader(204)
		}
	}))

	t.Run("invalid token", func(t *testing.T) {
		token := "wrong token"

		got, err := CheckAPItoken(token, ts.URL)
		want := false

		checkNoError(t, err)
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})

	t.Run("valid token", func(t *testing.T) {
		token := "validToken"

		got, err := CheckAPItoken(token, ts.URL)
		want := true

		checkNoError(t, err)
		if got != want {
			t.Errorf("got %v, wanted %v", got, want)
		}
	})
}

func TestGetHighlights(t *testing.T) {
	tsWrongSC := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	t.Run("unexpected response status code", func(t *testing.T) {
		_, err := GetHighlights(tsWrongSC.URL, "")
		if err == nil {
			t.Errorf("Wanted an error, didn't get one!")
		}
	})

	t.Run("correct response status code", func(t *testing.T) {
		resp, err := GetHighlights(ts.URL, "")
		checkNoError(t, err)
		fmt.Println(resp)
	})

}

func checkNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Got an error, didn't want one: %v", err)
	}
}
