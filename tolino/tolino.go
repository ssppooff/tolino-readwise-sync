package tolino

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

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

func extractEntries(input string) (entries []Entry, err error) {
	const delim = "\n\n-----------------------------------\n\n"
	tmp := strings.Split(input, delim)
	tmp = tmp[:len(tmp)-1]

	for _, t := range tmp {
		entryType := extractType(t)
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

func extractType(token string) string {
	pattern := regexp.MustCompile(`.*\n(\w+)\x{00A0}`)
	return pattern.FindStringSubmatch(token)[1]
}

func extractNote(token string) (entry Entry, err error) {
	pattern := regexp.MustCompile(`(?s)^(?P<title>.*)\s\((?P<author>.*, .*)\)\n(?P<type>.*?)\x{00A0}.+?(?P<page>\d+): (?P<note>.*)(?:\n?")(?P<highlight>.*)"\n(?P<isChange>.+) on\x{00A0}(?P<timestamp>.+)`)

	tmp := pattern.FindStringSubmatch(token)
	if len(tmp) != 9 {
		err = errors.New("issue extracting note: not enough entries found")
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
