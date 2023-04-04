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

const multipleHighlights = `The Complete Relaxation Book (Hewitt, James)
Bookmark on page 25: "PROGRAMME ONE

Relaxation should take place in a quiet, softly lit room that is pleasantly warm, airy but free from draughts. Later you may be able to relax even with slight discomfort or distraction, "
Added on 01/09/202 | 01:51

-----------------------------------

The 80/20 Principle (Koch, Richard)
Note on page 20: So do keep frequently used files on your desktop? ^^
"Incidentally, Zipf also provided a scientific justification for the messy desk by justifying clutter with another law: frequency of use draws near to us things that are frequently used. Intelligent secretaries have long known that files in frequent use should not be filed!"
Changed on 07/16/2022 | 23:29

-----------------------------------

Why Work Sucks and How to Fix It: The Results-Only Revolution (Ressler, Cali)
Highlight on page 24: " But mostly we learn about what is normal at work by experiencing it. One of the lessons work teaches us right away—whether we’re working at a restaurant or doing grunt work in an office or mowing lawns for our neighbors—is that there is the job you do and the job you appear to be doing."
Added on 12/21/2022 | 12:24

-----------------------------------

`

func Test_extractEntries(t *testing.T) {
	testCases := map[string]struct {
		input string
		num   int
	}{
		"number of entries": {input: multipleHighlights, num: 2},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got, err := extractEntries(tC.input); len(got) != tC.num {
				if err != nil {
					t.Fatalf("%#v: got and error, didn't want one %v", tC.input, err)
				}
				t.Errorf("%v, want %v", len(got), tC.num)
			}
		})
	}
}

func Test_extractType(t *testing.T) {
	tests := []struct {
		entry string
		want  string
	}{
		{want: "Note", entry: "Title of book (author)\nNote\u00a0"},
		{want: "Highlight", entry: "Title of book (author)\nHighlight\u00a0"},
		{want: "Bookmark", entry: "Title of book (author)\nBookmark\u00a0"},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := extractType(tt.entry); got != tt.want {
				t.Errorf("extractType() = %v, want %v", got, tt.want)
			}
		})
	}
}
