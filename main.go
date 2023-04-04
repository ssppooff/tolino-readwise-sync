package main

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const authURL = "https://readwise.io/api/v2/auth/"
const highlightsURL = "https://readwise.io/api/v2/highlights/"

func main() {

	token, err := readToken("token")
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}

	// auth, err := CheckAPItoken(token, authURL)
	// if err != nil {
	// 	fmt.Printf("Error: %v", err)
	// 	return
	// }
	// fmt.Println(auth)

	// _, err := decodeJSONpayload("highlights_JSON.json")
	resp, err := GetHighlights(highlightsURL, token)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Println(resp)
}

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
	body, err := io.ReadAll(respBody)
	if err != nil {
		return err
	}

	_, err = f.Write(body)
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}
	return nil
}

/*
	{
		"id": 59758950,
		"text": "The fundamental belief of metaphysicians is THE BELIEF IN ANTITHESES OF VALUES.",
		"note": "",
		"location": 9,
		"location_type": "order",
		"highlighted_at": null,
		"url": null,
		"color": "",
		"updated": "2020-10-01T12:58:44.716235Z",
		"book_id": 2608248,
		"tags": [
			{
					"id": 123456,
					"name": "philosophy"
			},
		]
	}
*/
type Highlight struct {
	ID, Bookd_ID         int
	Text, Note, Location string
	Tags                 []Tag
}

type Tag struct {
	ID   int
	Name string
}

type List[T any] struct {
	Count          int64
	Next, Previous string
	Results        []T
}

func decodeJSONpayload(filename string) (string, error) {
	payload, err := os.ReadFile(filename)
	if err != nil {
		return "", errors.Join(fmt.Errorf("decode JSON: couldn't open file %q", filename), err)
	}

	var list List[Highlight]
	// b2 := []byte(`{"count": 1912, "next": "url", "previous": null}`)
	err = json.Unmarshal(payload, &list)

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

func GetHighlights(url, token string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Join(errors.New("GetHighlights: couldn't create HTTP request"), err)
	}

	setAuthHeader(token, req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Join(errors.New("GetHighlights: couldn't send request"), err)
	}

	if resp.StatusCode != 200 {
		return nil, errors.Join(errors.New("GetHighlights: couldn't get highlights"), err)
	} else {
		if rh := resp.Header["Content-Type"]; len(rh) == 1 && rh[0] == "application/json" {
			// body, err := io.ReadAll(resp.Body)
			// writeJSONpayload(resp.Body)
			decoder := json.NewDecoder(resp.Body)

			var list List[Highlight]
			err := decoder.Decode(&list)
			if err != nil {
				return nil, errors.Join(errors.New("GetHighlights: couldn't decode response body"), err)
			}
			defer resp.Body.Close()
			return nil, nil
		}
		return resp, errors.New("Something wrong with the response")
	}
}


func readToken(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", errors.Join(fmt.Errorf("readToken: couldn't open file %q", filename), err)
	}

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
