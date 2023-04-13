package readwise

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	authURL       = "https://readwise.io/api/v2/auth/"
	highlightsURL = "https://readwise.io/api/v2/highlights/"
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

func GetHighlights(url, token string) (Page[Highlight], error) {
	var list Page[Highlight]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return list, errors.Join(errors.New("GetHighlights: couldn't create HTTP request"), err)
	}

	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return list, errors.Join(errors.New("GetHighlights: couldn't send request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return list, errors.Join(errors.New("GetHighlights: couldn't get highlights"), err)
	}

	rh := resp.Header["Content-Type"]
	if len(rh) != 1 || rh[0] != "application/json" {
		return list, fmt.Errorf("something wrong with response header: %#v", rh)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&list)
	if err != nil {
		return list, errors.Join(errors.New("GetHighlights: couldn't decode response body"), err)
	}

	return list, nil
}

func GetBooks(url, token string) (Page[Book], error) {
	var page Page[Book]
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return page, errors.Join(errors.New("GetBooks: couldn't create HTTP request"), err)
	}

	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return page, errors.Join(errors.New("GetBooks: couldn't send request"), err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return page, errors.Join(errors.New("GetBooks: couldn't get books"), err)
	}

	rh := resp.Header["Content-Type"]
	if len(rh) != 1 || rh[0] != "application/json" {
		return page, fmt.Errorf("something wrong with response header: %#v", rh)
	}

	decoder := json.NewDecoder(resp.Body)

	err = decoder.Decode(&page)
	if err != nil {
		return page, errors.Join(errors.New("GetBooks: couldn't decode response body:"), err)
	}
	return page, nil
}

func readToken(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", errors.Join(fmt.Errorf("readToken: couldn't open file %q", filename), err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	token := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", errors.Join(errors.New("readToken: error while scanning for token"), err)
	}
	return token, nil
}

func setAuthHeader(token string, r *http.Request) *http.Request {
	r.Header.Set("Authorization", fmt.Sprintf("Token %s", token))
	return r
}
