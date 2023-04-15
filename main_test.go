package main

import (
	// "encoding/json"
	// "encoding/json"
	// "fmt"
	"testing"
	"time"

	"github.com/ssppooff/tolino-readwise-sync/tolino"
)

func TestConvert(t *testing.T) {
	tmpTime, _ := time.ParseInLocation(time.DateTime, "2022-09-01 21:36:00", time.Local)
	var te1 = tolino.Entry{Title: "Random Title", Author: "Random Author", Page: 33, Note: "Some Note", Highlight: "High Text", Changed: false, Date: tmpTime}
	var te2 = tolino.Entry{Title: "Another Random Title", Author: "Random Author", Page: 73, Note: "Some Note", Highlight: "High Text", Changed: false, Date: tmpTime.Add(time.Hour * 1)}

	t.Run("single entry", func(t *testing.T) {
		got := Convert([]tolino.Entry{te1})
		want := `{"Highlights":[{"Text":"High Text","Note":"Some Note","Title":"Random Title","Author":"Random Author","Location":33,"Location_type":"page","Source_type":"tolino-sync","Category":"books","Highlighted_at":"2022-09-01T21:36:00+02:00"}]}`

		if got != want {
			t.Errorf("\ngot : %q,\nwant: %q", got, want)
		}
	})

	t.Run("multiple entries", func(t *testing.T) {
		got := Convert([]tolino.Entry{te1, te2})
		want := `{"Highlights":[{"Text":"High Text","Note":"Some Note","Title":"Random Title","Author":"Random Author","Location":33,"Location_type":"page","Source_type":"tolino-sync","Category":"books","Highlighted_at":"2022-09-01T21:36:00+02:00"},{"Text":"High Text","Note":"Some Note","Title":"Another Random Title","Author":"Random Author","Location":73,"Location_type":"page","Source_type":"tolino-sync","Category":"books","Highlighted_at":"2022-09-01T22:36:00+02:00"}]}`

		if got != want {
			t.Errorf("\ngot : %q,\nwant: %q", got, want)
		}
	})
}
