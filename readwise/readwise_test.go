package readwise

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/utils"
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

func TestGetPageParams(t *testing.T) {
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
			url := r.URL
			q := url.Query()
			if reflect.DeepEqual(q["source"], []string{"appID"}) {
				w.Header().Set("Content-Type", "application/json")
				const jsonPayload = `{"count":1, "next": null, "previous": null, "results":[{"id":100, "title":"Random Title", "author": "John Doe", "category": "books", "source": "appID", "num_highlights": 68, "last_highlight_at": "2020-10-01T17:47:31.234826Z", "asin": "B0046LU7H0"}]}`
				w.Write([]byte(jsonPayload))
			}
		}
	}))
	defer ts.Close()

	t.Run("unexpected response status code", func(t *testing.T) {
		params := url.Values{}
		err := GetPageParams(&Page[Highlight]{}, ts.URL, "wrongToken", &params)
		if err == nil {
			t.Errorf("Wanted an error, didn't get one!")
		}
	})

	t.Run("get highlight, no parameters", func(t *testing.T) {
		want := Page[Highlight]{Count: 1, Results: []Highlight{{
			ID:   100,
			Text: "random text",
		}}}

		var page Page[Highlight]
		params := url.Values{}
		err := GetPageParams(&page, ts.URL, "validHighlightToken", &params)
		checkNoError(t, err)
		if !reflect.DeepEqual(page, want) {
			t.Errorf("wrong page,\ngot   : %#v,\nwanted: %#v", page, want)
		}
	})

	t.Run("get book, parameter source", func(t *testing.T) {
		hlTime, err := time.Parse(time.RFC3339, "2020-10-01T17:47:31.234826Z")
		if err != nil {
			t.Fatalf("couldn't parse time string: %#v", err)
		}

		want := Page[Book]{Count: 1, Results: []Book{{
			ID:                100,
			Title:             "Random Title",
			Author:            "John Doe",
			Category:          "books",
			Source:            "appID",
			Num_highlights:    68,
			Last_highlight_at: hlTime,
			ASIN:              "B0046LU7H0",
		}}}

		var page Page[Book]
		params := url.Values{"source": {"appID"}}
		err = GetPageParams(&page, ts.URL, "validBookToken", &params)
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

func TestCreateHighlight(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ah := r.Header["Authorization"]
		if len(ah) != 1 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if ah[0] == "Token wrongToken" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if ah[0] == "Token validToken" {
			body, _ := io.ReadAll(r.Body)
			defer r.Body.Close()

			const jsonPayload = `{"Highlights":[{"Text":"Some Text","Note":"","Title":"Title1","Author":"John Doe","Location":null,"Location_type":null,"Source_type":"app_ID1","Category":"books","Highlighted_at":null},{"Text":"Some more text","Note":"some note","Title":"Title2","Author":"John Smith","Location":null,"Location_type":null,"Source_type":"app_ID2","Category":"books","Highlighted_at":null}]}`
			if string(body) != jsonPayload {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			// API call returns list of created/updated books/articles/podcasts
			w.Header().Set("Content-Type", "application/json")
			const respJSON = `[{"title": "Title1", "author": "John Doe", "source": "app_ID1", "category": "books"}, {"title": "Title2", "author": "John Smith", "source": "app_ID2", "category": "books"}]`
			w.Write([]byte(respJSON))
		}
	}))
	defer ts.Close()

	hls := []HlCreate{
		{
			Text:        "Some Text",
			Note:        "",
			Title:       utils.Ptr("Title1"),
			Author:      utils.Ptr("John Doe"),
			Category:    utils.Ptr("books"),
			Source_type: utils.Ptr("app_ID1")},
		{
			Text:        "Some more text",
			Note:        "some note",
			Title:       utils.Ptr("Title2"),
			Author:      utils.Ptr("John Smith"),
			Category:    utils.Ptr("books"),
			Source_type: utils.Ptr("app_ID2")}}

	t.Run("positive answer", func(t *testing.T) {
		token := "validToken"

		// Both highlights together should create 2 new books: Title1, by John Doe & Title2, by John Smith
		err := CreateHighlights(hls, ts.URL, token)
		checkNoError(t, err)
	})
}
