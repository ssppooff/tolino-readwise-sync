package readwise

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestCheckAPItoken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header["Authorization"]
		if len(ah) == 1 && ah[0] == "Token validToken" {
			w.WriteHeader(204)
		}
	}))
	defer ts.Close()

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

func TestGetPage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header["Authorization"]
		if len(ah) != 1 {
			return
		}

		if ah[0] == "Token wrongToken" {
			w.WriteHeader(210)
		}

		if ah[0] == "Token validHighlightToken" {
			w.Header().Set("Content-Type", "application/json")
			const jsonPayload = `{"count":1, "next": null, "previous": null, "results":[{"id":100,"text":"random text"}]}`
			w.Write([]byte(jsonPayload))
		}

		if ah[0] == "Token validBookToken" {
			w.Header().Set("Content-Type", "application/json")
			const jsonPayload = `{"count":1, "next": null, "previous": null, "results":[{"id":100, "title":"Random Title", "author": "John Doe", "num_highlights": 68, "last_highlight_at": "2020-10-01T17:47:31.234826Z", "asin": "B0046LU7H0"}]}`
			w.Write([]byte(jsonPayload))
		}
	}))
	defer ts.Close()

	t.Run("unexpected response status code", func(t *testing.T) {
		err := GetPage(&Page[Highlight]{}, ts.URL, "wrongToken")
		if err == nil {
			t.Errorf("Wanted an error, didn't get one!")
		}
	})

	t.Run("get highlight", func(t *testing.T) {
		want := Page[Highlight]{Count: 1, Results: []Highlight{{ID: 100, Text: "random text"}}}

		var page Page[Highlight]
		err := GetPage(&page, ts.URL, "validHighlightToken")
		checkNoError(t, err)
		if !reflect.DeepEqual(page, want) {
			t.Errorf("wrong page,\ngot   : %#v,\nwanted: %#v", page, want)
		}
	})

	t.Run("get book", func(t *testing.T) {
		hlTime, err := time.Parse(time.RFC3339, "2020-10-01T17:47:31.234826Z")
		if err != nil {
			t.Fatalf("couldn't parse time string: %#v", err)
		}

		want := Page[Book]{Count: 1, Results: []Book{{ID: 100,
			Title:             "Random Title",
			Author:            "John Doe",
			Num_highlights:    68,
			Last_highlight_at: hlTime,
			ASIN:              "B0046LU7H0",
		}}}

		var page Page[Book]
		err = GetPage(&page, ts.URL, "validBookToken")
		checkNoError(t, err)
		if !reflect.DeepEqual(page, want) {
			t.Errorf("wrong page,\ngot   : %#v,\nwanted: %#v", page, want)
		}
	})
}

func checkNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Got an error, didn't want one: %v", err)
	}
}
