package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/ssppooff/tolino-readwise-sync/readwise"
	"github.com/ssppooff/tolino-readwise-sync/tolino"
	"github.com/ssppooff/tolino-readwise-sync/utils"
)

const appIdentifier = "tolino-sync"

type args struct {
	TokenFile  string `arg:"-t,--,required" placeholder:"TOKEN_FILE" help:"path to file with API token"`
	TolinoFile string `arg:"-n,--,required" placeholder:"NOTES_FILE" help:"path to file with highlights & notes from Tolino"`
}

func (args) Description() string {
	return "Uploads all highlights and notes from TOLINO_FILE to Readwise\n"
}

func (args) Epilogue() string {
	return "For more information visit github.com/ssppooff/tolino-readwise-sync"
}

func main() {
	var args args
	arg.MustParse(&args)
	token, err := readToken(args.TokenFile)
	if err != nil {
		panic(err)
	}

	ok, err := readwise.CheckAPItoken(token, readwise.AuthURL)
	if err != nil {
		panic(err)
	}
	if !ok {
		panic("Token not valid")
	} else {
		fmt.Println("Token valid")
	}

	content, err := os.ReadFile(args.TolinoFile)
	if err != nil {
		panic(err)
	}

	entries, err := tolino.ExtractEntries(string(content))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Extracted %d highlights & notes\n", len(entries))

	modBooks, err := readwise.CreateHighlights(utils.Map(entries, Convert), readwise.HighlightsURL, token)
	if err != nil {
		fmt.Println()
		panic(err)
	}
	fmt.Println("Upload successful")
	fmt.Printf("Added or modified %d book(s): %s\n", len(modBooks), strings.Join(modBooks, ", "))
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
func Convert(te tolino.Entry) readwise.HlCreate {
	return readwise.HlCreate{
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
}
