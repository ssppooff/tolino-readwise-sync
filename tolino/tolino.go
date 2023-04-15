package tolino

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/utils"
)

const (
	ErrTypeExtraction   = TolinoErrStr("couldn't extract type")
	ErrNotEnoughEntries = TolinoErrStr("not enough entries found")
	ErrWrongTimeStamp   = TolinoErrStr("date and time in wrong format")
)

type TolinoErrStr string

type TolinoError struct {
	err   TolinoErrStr
	entry string
}

func (e TolinoError) Error() string {
	return string(e.err)
}

func (e TolinoError) Is(errStr TolinoErrStr) bool {
	return e.err == errStr
}

type Entry struct {
	Highlight string
	Note      string
	Title     string
	Author    string
	EntryType string // really needed?
	Page      int
	Changed   bool
	Date      time.Time
}

type Book struct {
	title  string
	author string
}

func (e Entry) GetBook() Book {
	return Book{e.Title, e.Author}
}

func (e Entry) isEmpty() bool {
	if e.EntryType == "" && e.Title == "" && e.Author == "" && e.Highlight == "" {
		return true
	}

	return false
}

func ExtractEntries(input string) (entries []Entry, err error) {
	const delim = "\n\n-----------------------------------\n\n"
	tmp := strings.Split(input, delim)
	tmp = tmp[:len(tmp)-1]

	entries, err = utils.MapWithErr(tmp, func(t string) (Entry, error) {
		entryType, err := extractType(t)
		if err != nil {
			return Entry{}, TolinoError{ErrTypeExtraction, err.(TolinoError).entry}
		}

		if entryType == "Note" || entryType == "Highlight" {
			entry, err := extractNote(t)
			if err != nil {
				return Entry{}, err
			}
			return entry, nil
		}

		return Entry{}, nil
	})

	_, entries = utils.Filter(entries, Entry.isEmpty)
	return
}

func extractType(token string) (string, error) {
	pattern := regexp.MustCompile(`.*\n(\w+)\x{00A0}`)

	matches := pattern.FindStringSubmatch(token)
	if len(matches) < 2 {
		return "", TolinoError{err: "couldn't extract entry type", entry: token}
	}

	return matches[1], nil
}

func extractNote(token string) (entry Entry, err error) {
	pattern := regexp.MustCompile(`(?s)^(?P<title>.*)\s\((?P<author>.*, .*)\)\n(?P<type>.*?)\x{00A0}.+?(?P<page>\d+): (?P<note>.*)(?:\n?")(?P<highlight>.*)"\n(?P<isChange>.+) on\x{00A0}(?P<timestamp>.+)`)

	tmp := pattern.FindStringSubmatch(token)
	if len(tmp) != 9 {
		err = TolinoError{err: ErrNotEnoughEntries, entry: token}
		return
	}

	// won't return error: if there is no page number, or in a wrong format (ie.,
	//   not regex `\d+`), FindStringSubmatch(token) will not find enough entries
	//   and will therefore return from this function with an error
	pageNum, _ := strconv.Atoi(tmp[4])

	var hasChanged bool
	if tmp[7] == "Changed" {
		hasChanged = true
	}

	const timestampLayout = "01/02/2006 | 15:04"
	timestamp, err := time.ParseInLocation(timestampLayout, tmp[8], time.Local)
	if err != nil {
		err = TolinoError{err: ErrWrongTimeStamp, entry: token}
		return
	}
	timestamp = timestamp.Local()

	entry = Entry{
		Title:     tmp[1],
		Author:    tmp[2],
		EntryType: tmp[3],
		Page:      pageNum,
		Note:      tmp[5],
		Highlight: tmp[6],
		Changed:   hasChanged,
		Date:      timestamp,
	}

	return
}

func ExtractBooks(entries []Entry) []Book {
	return utils.Map(entries, func(e Entry) Book { return e.GetBook() })
}
