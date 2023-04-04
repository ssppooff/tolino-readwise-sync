package tolino

import (
	"testing"
)

const highlightTest = `The 80/20 Principle (Koch, Richard)
Highlight on page 126: "The 80/20 Principle treats time as a friend, not an enemy. Time gone is not time lost. Time will always come round again. This is why there are seven days in a week, twelve months in a year, why the seasons come round again. Insight and value are likely to come from placing ourselves in a comfortable, relaxed, and collaborative position toward time. It is our use of time, and not time itself, that is the enemy.
			• The 80/20 Principle says that we should act less. Action drives out thought. It is because we have so much time that we squander it. The most productive time on a project is usually the last 20 percent, simply because the work has to be completed before a deadline. Productivity on most projects could be doubled simply by halving the amount of time for their completion. This is not evidence that time is in short supply."
Added on 07/18/2022 | 13:00

-----------------------------------`

const noteTest = `The 80/20 Principle (Koch, Richard)
Note on page 44: That would be fixing your emotions!
"exercise control over our lives with the least possible effort"
Added on 07/19/2022 | 20:59

-----------------------------------`

func TestFoo(t *testing.T) {
	tests := []struct {
		name  string
		entry string
	}{
		{name: "highlight", entry: highlightTest},
		{name: "note", entry: noteTest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Foo(tt.entry)
		})
	}
}

func Test_extractType(t *testing.T) {
	const noteType = `The 80/20 Principle (Koch, Richard)
Note`
	const highlightType = `The 80/20 Principle (Koch, Richard)
Highlight`
	tests := []struct {
		name  string
		entry string
		want  string
	}{
		{name: "note type", entry: noteType, want: "Note"},
		{name: "note type", entry: highlightType, want: "Highlight"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractType(tt.entry); got != tt.want {
				t.Errorf("extractType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractPage(t *testing.T) {
	const noteType = `44: That would be fixing your emotions!
	"exercise control over our lives with the least possible effort"
	Added on`
	const highlightType = `126: "The 80/20 Principle treats time as a friend, not an enemy. Time gone is not time lost. Time will always come round again. This is why there are seven days in a week, twelve months in a year, why the seasons come round again. Insight and value are likely to come from placing ourselves in a comfortable, relaxed, and collaborative position toward time. It is our use of time, and not time itself, that is the enemy.
	• The 80/20 Principle says that we should act less. Action drives out thought. It is because we have so much time that we squander it. The most productive time on a project is usually the last 20 percent, simply because the work has to be completed before a deadline. Productivity on most projects could be doubled simply by halving the amount of time for their completion. This is not evidence that time is in short supply."
Added on`

	tests := []struct {
		name  string
		entry string
		want  int
	}{
		{name: "note type", entry: noteType, want: 44},
		{name: "note type", entry: highlightType, want: 126},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractPage(tt.entry); got != tt.want {
				t.Errorf("extractType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractNote(t *testing.T) {
	const noteType = `44: That would be fixing your emotions!
	"exercise control over our lives with the least possible effort"
	Added on`
	const highlightType = `126: "The 80/20 Principle treats time as a friend, not an enemy. Time gone is not time lost. Time will always come round again. This is why there are seven days in a week, twelve months in a year, why the seasons come round again. Insight and value are likely to come from placing ourselves in a comfortable, relaxed, and collaborative position toward time. It is our use of time, and not time itself, that is the enemy.
	• The 80/20 Principle says that we should act less. Action drives out thought. It is because we have so much time that we squander it. The most productive time on a project is usually the last 20 percent, simply because the work has to be completed before a deadline. Productivity on most projects could be doubled simply by halving the amount of time for their completion. This is not evidence that time is in short supply."
Added on`

	const shortNote = `44: That would!
	"exercise control !"
	Added on`

	tests := []struct {
		name      string
		entry     string
		page      string
		note      string
		highlight string
	}{
		// {name: "note type", entry: noteType, want: ""},
		// {name: "note type", entry: highlightType, want: ""},
		{name: "note type", entry: shortNote, page: "44", note: "That would!", highlight: "exercise control !"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotNote, gotHL := extractNote(tt.entry)
			if gotPage != tt.page {
				t.Errorf("extractNote() = %v, want %v", gotPage, tt.page)
			}
			if gotNote != tt.note {
				t.Errorf("extractNote() = %q, want %q", gotNote, tt.note)
			}
			if gotHL != tt.highlight {
				t.Errorf("extractNote() = %q, want %q", gotHL, tt.highlight)
			}
		})
	}
}
