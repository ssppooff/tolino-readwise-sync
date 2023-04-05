package tolino

import (
	"errors"
	"fmt"
	"testing"
)

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

const noType = `
Bookmark on page 25: "PROGRAMME ONE

Relaxation should take place in a quiet, softly lit room that is pleasantly warm, airy but free from draughts. Later you may be able to relax even with slight discomfort or distraction, "
Added on 01/09/202 | 01:51

-----------------------------------

Relaxation should take place in a quiet, softly lit room that is pleasantly warm, airy but free from draughts. Later you may be able to relax even with slight discomfort or distraction, "
Added on 01/09/202 | 01:51

-----------------------------------

Why Work Sucks and How to Fix It: The Results-Only Revolution (Ressler, Cali)
Highlight on page 24: " But mostly we learn about what is normal at work by experiencing it. One of the lessons work teaches us right away—whether we’re working at a restaurant or doing grunt work in an office or mowing lawns for our neighbors—is that there is the job you do and the job you appear to be doing."
Added on 12/21/2022 | 12:24

-----------------------------------

`

const wrongHighlight = `
Bookmark on page 25: "PROGRAMME ONE

Relaxation should take place in a quiet, softly lit room that is pleasantly warm, airy but free from draughts. Later you may be able to relax even with slight discomfort or distraction, "
Added on 01/09/202 | 01:51

-----------------------------------

The 80/20 Principle (Koch, Richard)
Note on page 20: So do keep frequently used files on your desktop? ^^
"Incidentally, Zipf also provided a scientific justification for the messy desk by justifying clutter with another law: frequency of use draws near to us things that are frequently used. Intelligent secretaries have long known that files in frequent use should not be filed!"

-----------------------------------

Why Work Sucks and How to Fix It: The Results-Only Revolution (Ressler, Cali)
Highlight on page 24: " But mostly we learn about what is normal at work by experiencing it. One of the lessons work teaches us right away—whether we’re working at a restaurant or doing grunt work in an office or mowing lawns for our neighbors—is that there is the job you do and the job you appear to be doing."
Added on 12/21/2022 | 12:24

-----------------------------------

`

func TestExtractEntries(t *testing.T) {
	testCases := map[string]struct {
		input string
		num   int
	}{
		// "number of entries": {input: multipleHighlights, num: 2},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got, err := ExtractEntries(tC.input); len(got) != tC.num {
				if err != nil {
					t.Fatalf("%q: got and error, didn't want one: %v", name, err)
				}
				t.Errorf("%v, want %v", len(got), tC.num)
			}
		})
	}

	errorCases := map[string]struct {
		input  string
		errStr string
	}{
		"not enough entries": {input: wrongHighlight, errStr: ErrNotEnoughEntries},
		"can't extract type": {input: noType, errStr: ErrTypeExtraction},
	}
	for name, eC := range errorCases {
		t.Run(name, func(t *testing.T) {
			_, err := ExtractEntries(eC.input)
			if err == nil {
				t.Fatalf("%q: didn't get an error, wanted one", name)
			}

			toErr := err.(TolinoError)
			if toErr.err != eC.errStr {
				fmt.Println(errors.Unwrap(err))
				t.Errorf("%q: got %q, wanted %q", name, toErr, eC.errStr)
			}
		})
	}
}
