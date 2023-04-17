package readwise

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

func writeJSONpayload(respBody io.Reader) error {
	f, err := os.Create("highlights_JSON.json")
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := io.ReadAll(respBody)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	return nil
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
	Source, Updated   string
	ASIN              string
	Num_highlights    int
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

func decodeJSONpayload(filename string) (string, error) {
	payload, err := os.ReadFile(filename)
	if err != nil {
		return "", errors.Join(fmt.Errorf("decode JSON: couldn't open file %q", filename), err)
	}

	var page Page[Highlight]
	err = json.Unmarshal(payload, &page)

	if err != nil {
		return "", errors.Join(fmt.Errorf("decode JSON: couldn't unmarshal JSON"), err)
	}

	// var page map[string]interface{}
	// if err != nil {
	// 	return nil, errors.Join(errors.New("GetHighlights: couldn't decode response body"), err)
	// }
	// defer resp.Body.Close()
	// return nil, nil

	// 	return "", errors.Join(errors.New("readToken: error while scanning for token"), err)
	// }
	return "", nil
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
	Text           string
	Note           string
	Title, Author  *string
	Location       *int
	Location_type  *string
	Source_type    *string // app identifier
	Category       *string
	Highlighted_at *string // ISO 8601
}

func CreateHighlights(hls []HlCreate, url, token string) error {
	var payload = struct{ Highlights []HlCreate }{Highlights: hls}

	body, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(errors.New("error while converting highlights to JSON"))
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return errors.Join(errors.New("couldn't create HTTP POST request"), err)
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Join(errors.New("couldn't send POST request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Join(errors.New("received error after trying creating highlights on Readwise"), err)
	}

	ct := resp.Header["Content-Type"]
	if len(ct) != 1 || !strings.Contains(ct[0], "application/json") {
		err = fmt.Errorf("something wrong with response header or content-type: %#v", ct)
		return
	}

	want := extractBooksFromHighlights(hls)
	got := []Book{}
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		return errors.Join(errors.New("couldn't decode response body"), err)
	}

	if !reflect.DeepEqual(got, want) {
		return errors.New("response does not correspond to what is expected")
	}

	return nil
}

func extractBooksFromHighlights(hls []HlCreate) (books []Book) {
	for _, hl := range hls {
		books = append(books, Book{
			Author:   *hl.Author,
			Title:    *hl.Title,
			Source:   *hl.Source_type,
			Category: *hl.Category,
		})
	}
	return
}
