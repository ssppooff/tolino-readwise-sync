package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/readwise"
	"github.com/ssppooff/tolino-readwise-sync/tolino"
	"github.com/ssppooff/tolino-readwise-sync/utils"
)

const appIdentifier = "tolino-sync"

/*
1. parse Tolino file
2. Filter for only new highlights
3. check API token
4. Transform new Tolino Highlights into compatible Readwise highlights, add "Tolino" as source
5. Upload all highlights to Readwise
*/
func main() {
	tolino_file := "path to tolino notes.txt file"
	file, _ := os.Open(tolino_file)
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	entries, _ := tolino.ExtractEntries(string(bytes))
	entries, _ = utils.Filter(entries, func(te tolino.Entry) bool { return te.Changed == false })

	filename := "path_to_token_file"
	token, _ := readToken(filename)
	ok, _ := readwise.CheckAPItoken(token, readwise.AuthURL)
	if !ok {
		return
	}

	// TODO refactor: don't give JSON to readwise package
	err := readwise.CreateHighlights(Convert(entries), readwise.HighlightsURL, token)
	if err != nil {
		return
	}
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

// Sets the 'location_type' to 'page', 'category' to 'books', and 'source_type' to the content of the constant 'appIdentifier'
func Convert(tes []tolino.Entry) []readwise.HlCreate {
	hls := []readwise.HlCreate{}

	for _, te := range tes {
		var tmp = readwise.HlCreate{
			Text:           te.Highlight,
			Note:           te.Note,
			Title:          &te.Title,
			Author:         &te.Author,
			Location:       &te.Page,
			Location_type:  utils.Ptr("page"),
			Source_type:    utils.Ptr(appIdentifier),
			Category:       utils.Ptr("books"),
			Highlighted_at: utils.Ptr(te.Date.Format(time.RFC3339)),
		}
		hls = append(hls, tmp)
	}

	return hls
}
