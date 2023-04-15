package main

import (
	"testing"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/readwise"
	"github.com/ssppooff/tolino-readwise-sync/tolino"
	"github.com/ssppooff/tolino-readwise-sync/utils"
)

func TestConvert(t *testing.T) {
	tmpTime, _ := time.ParseInLocation(time.DateTime, "2022-09-01 21:36:00", time.Local)
	var te = tolino.Entry{
		Title:     "Random Title",
		Author:    "Random Author",
		Page:      33,
		Note:      "Some Note",
		Highlight: "High Text",
		Changed:   false,
		Date:      tmpTime,
	}

	want := readwise.HlCreate{
		Text:           "High text",
		Note:           "Some Note",
		Title:          utils.Ptr("Random Title"),
		Author:         utils.Ptr("Random Author"),
		Location:       utils.Ptr(33),
		Location_type:  utils.Ptr("page"),
		Source_type:    utils.Ptr(appIdentifier),
		Category:       utils.Ptr("books"),
		Highlighted_at: utils.Ptr(tmpTime.Format(time.RFC3339)),
	}
	got := Convert(te)

	if checkHlCreate(t, got, want) {
		t.Errorf("\ngot : %v,\nwant: %v", got, want)
	}
}

func checkHlCreate(t *testing.T, got, want readwise.HlCreate) bool {
	t.Helper()
	return got.Text == want.Text &&
		got.Note == want.Note &&
		*got.Title == *want.Title &&
		*got.Author == *want.Author &&
		*got.Location == *want.Location &&
		*got.Source_type == appIdentifier &&
		*got.Highlighted_at == *want.Highlighted_at
}
