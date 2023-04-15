package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	//"github.com/ssppooff/tolino-readwise-sync/readwise"
	"github.com/ssppooff/tolino-readwise-sync/tolino"
)

func main() {}

const appIdentifier = "tolino-sync"

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

type HlCreate struct {
	Text           string
	Note           string
	Title, Author  string
	Location       int    // page number
	Location_type  string // always page
	Source_type    string // app identifier
	Category       string // always books
	Highlighted_at string // ISO 8601
}

// Convert converts multiple Tolino entries into a JSON string suitable for sending to API call Highlight CREATE
//
// Sets the 'location_type' to 'page', 'category' to 'books', and 'source_type' to the content of the constant 'appIdentifier'
func Convert(tes []tolino.Entry) (string, error) {
	type payload struct{ Highlights []HlCreate }
	var p = payload{}

	for _, te := range tes {
		var tmp = HlCreate{
			Text:           te.Highlight,
			Note:           te.Note,
			Title:          te.Title,
			Author:         te.Author,
			Location:       te.Page,
			Location_type:  "page",
			Source_type:    appIdentifier,
			Category:       "books",
			Highlighted_at: te.Date.Format(time.RFC3339),
		}
		p.Highlights = append(p.Highlights, tmp)
	}

	b, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
