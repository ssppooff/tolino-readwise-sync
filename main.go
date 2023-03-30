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

	auth, err := CeckAPItoken(token, authURL)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Println(auth)

}

// GET Request to https://readwise.io/api/v2/auth/ with header: key "Authorization", value "Token XXX"
func CeckAPItoken(token string, url string) (bool, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return false, errors.Join(errors.New("checkAPItoken: couldn't create HTTP request"), err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, errors.Join(errors.New("checkAPItoken: couldn't send request"), err)
	}

	if resp.StatusCode != 204 {
		return false, nil
	} else {
		return true, nil
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
