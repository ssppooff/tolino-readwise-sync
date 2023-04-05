package tolino

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	ErrTypeExtraction   = "couldn't extract type"
	ErrNotEnoughEntries = "not enough entries found"
)

type TolinoError struct {
	err   string
	entry string
}

func (e TolinoError) Error() string {
	return e.err
}

type Entry struct {
	title     string
	author    string
	entryType string
	note      string
	page      int
	highlight string
	changed   bool
	date      time.Time
}

type Book struct {
	title  string
	author string
}

func (e Entry) GetBook() Book {
	return Book{e.title, e.author}
}

func ExtractEntries(input string) (entries []Entry, err error) {
	const delim = "\n\n-----------------------------------\n\n"
	tmp := strings.Split(input, delim)
	tmp = tmp[:len(tmp)-1]

	for _, t := range tmp {
		entryType, err := extractType(t)
		if err != nil {
			return nil, TolinoError{ErrTypeExtraction, err.(TolinoError).entry}
		}

		if entryType == "Note" || entryType == "Highlight" {
			entry, err := extractNote(t)
			if err != nil {
				return nil, err
			}

			entries = append(entries, entry)
		}
	}

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

	pageNum, err := strconv.Atoi(tmp[4])
	if err != nil {
		return
	}

	var hasChanged bool
	if tmp[7] == "Changed" {
		hasChanged = true
	}

	const timestampLayout = "01/02/2006 | 15:04"
	timestamp, err := time.Parse(timestampLayout, tmp[8])
	if err != nil {
		return
	}

	entry = Entry{
		title:     tmp[1],
		author:    tmp[2],
		entryType: tmp[3],
		page:      pageNum,
		note:      strings.TrimSpace(tmp[5]),
		highlight: tmp[6],
		changed:   hasChanged,
		date:      timestamp,
	}

	return
}

func ExtractBooks(entries []Entry) []Book {
	return mapF(entries, func(e Entry) Book { return e.GetBook() })
}

func mapF(entries []Entry, fn func(Entry) Book) []Book {
	var res = make([]Book, len(entries))
	for i, entry := range entries {
		res[i] = fn(entry)
	}
	return res
}
