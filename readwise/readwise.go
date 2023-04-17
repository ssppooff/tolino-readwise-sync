package readwise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/utils"
	"golang.org/x/exp/slices"
)

const (
	AuthURL       = "https://readwise.io/api/v2/auth/"
	HighlightsURL = "https://readwise.io/api/v2/highlights/"
	BooksURL      = "https://readwise.io/api/v2/books/"
)

func main() {}

// GET Request to https://readwise.io/api/v2/auth/ with header: key "Authorization", value "Token XXX"
func CheckAPItoken(token string, url string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, errors.Join(errors.New("checkAPItoken: couldn't create HTTP request"), err)
	}

	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.Join(errors.New("checkAPItoken: couldn't send request"), err)
	}

	if resp.StatusCode != 204 {
		return false, nil
	}

	return true, nil
}

type Highlight struct {
	ID, Bookd_ID         int
	Text, Note, Location string
	Tags                 []Tag
}

type Tag struct {
	ID   int
	Name string
}

type Book struct {
	ID                int
	Title, Author     string
	Source            string
	ASIN              string
	Num_highlights    int
	Updated           time.Time
	Last_highlight_at time.Time
	Document_note     string
	Tags              []Tag
	Category          string
}

type Page[E Highlight | Tag | Book] struct {
	Count          int64
	Next, Previous string
	Results        []E
}

// Gotcha: parameters will be modified inside function
func GetAll[E Highlight | Book](sl *[]E, apiURL, token string, parameters *url.Values) error {
	var page = Page[E]{}
	var pageCount = 1
	err := GetPageParams(&page, apiURL, token, parameters)
	if err != nil {
		return errors.Join(errors.New("error while fetching pages"), err)
	}

	(*sl) = append((*sl), page.Results...)
	for page.Next != "" {
		pageCount += 1
		(*parameters).Set("page", strconv.Itoa(pageCount))
		page = Page[E]{}
		err = GetPageParams(&page, apiURL, token, parameters)
		if err != nil {
			return errors.Join(errors.New("error while fetching pages"), err)
		}
		(*sl) = append((*sl), page.Results...)
	}

	return nil
}

func GetPageParams[E Highlight | Book](page *Page[E], url, token string, parameters *url.Values) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return errors.Join(errors.New("couldn't create HTTP GET request"), err)
	}

	req.URL.RawQuery = (*parameters).Encode()
	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Join(errors.New("couldn't send GET request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Join(errors.New("couldn't get page"), err)
	}

	rh := resp.Header["Content-Type"]
	if len(rh) != 1 || rh[0] != "application/json" {
		return fmt.Errorf("something wrong with response header: %#v", rh)
	}

	err = json.NewDecoder(resp.Body).Decode(page)
	if err != nil {
		return errors.Join(errors.New("couldn't decode response body:"), err)
	}
	return nil
}

func setAuthHeader(token string, r *http.Request) *http.Request {
	r.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	return r
}

// Values not required by API call are pointer to base type
type HlCreate struct {
	Text           string  `json:"text"`
	Note           string  `json:"note,omitempty"`
	Title          *string `json:"title,omitempty"`
	Author         *string `json:"author,omitempty"`
	Location       *int    `json:"location,omitempty"`
	Location_type  *string `json:"location_type,omitempty"`
	Source_type    *string `json:"source_type,omitempty"` // app identifier
	Category       *string `json:"category,omitempty"`
	Highlighted_at *string `json:"highlighted_at,omitempty"` // ISO 8601
}

func CreateHighlights(hls []HlCreate, url, token string) (modBooks []string, err error) {
	var payload = struct {
		Highlights []HlCreate `json:"highlights"`
	}{Highlights: hls}

	body, err := json.Marshal(payload)
	if err != nil {
		err = errors.Join(errors.New("error while converting highlights to JSON"))
		return
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		err = errors.Join(errors.New("couldn't create HTTP POST request"), err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Join(errors.New("couldn't send POST request"), err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Join(fmt.Errorf("received status code: %d", resp.StatusCode))
		return
	}

	ct := resp.Header["Content-Type"]
	if len(ct) != 1 || !strings.Contains(ct[0], "application/json") {
		err = fmt.Errorf("something wrong with response header or content-type: %#v", ct)
		return
	}

	got := []Book{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		err = errors.Join(errors.New("couldn't decode response body"), err)
		return
	}

	modBooks, ok := checkResponseBody(hls, got)
	if !ok {
		err = errors.New("response does not correspond to what is expected")
		return
	}
	return
}

func checkResponseBody(hls []HlCreate, books []Book) ([]string, bool) {
	bksSent := []Book{}
	for _, hl := range hls {
		bksSent = append(bksSent, Book{
			Author:   *hl.Author,
			Title:    *hl.Title,
			Source:   *hl.Source_type,
			Category: *hl.Category,
		})
	}

	bksSent = slices.CompactFunc(bksSent, func(t, o Book) bool {
		return t.Title == o.Title &&
			t.Author == o.Author &&
			t.Source == o.Source &&
			t.Category == o.Category
	})

	if len(bksSent) == len(books) {
		modBooks := utils.Map(books, func(b Book) string { return b.Title })
		return modBooks, true
	}
	return []string{}, false
}
