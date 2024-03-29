package tolino

import (
	"reflect"
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

const noPage = `

Why Work Sucks and How to Fix It: The Results-Only Revolution (Ressler, Cali)
Highlight on page 24b: " But mostly we learn about what is normal at work by experiencing it. One of the lessons work teaches us right away—whether we’re working at a restaurant or doing grunt work in an office or mowing lawns for our neighbors—is that there is the job you do and the job you appear to be doing."
Added on 12/21/2022 | 12:24

-----------------------------------

`

const wrongTime = `

Why Work Sucks and How to Fix It: The Results-Only Revolution (Ressler, Cali)
Highlight on page 24: " But mostly we learn about what is normal at work by experiencing it. One of the lessons work teaches us right away—whether we’re working at a restaurant or doing grunt work in an office or mowing lawns for our neighbors—is that there is the job you do and the job you appear to be doing."
Added on 12/21/2022 12:24

-----------------------------------

`

func TestExtractEntries(t *testing.T) {
	testCases := map[string]struct {
		input string
		num   int
	}{
		"number of entries": {input: multipleHighlights, num: 2},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got, err := ExtractEntries(tC.input); len(got) != tC.num {
				if err != nil {
					t.Fatalf("%q: got and error, didn't want one: %v", name, err)
				}

				t.Errorf("got %v, want %v", len(got), tC.num)
			}
		})
	}

	errorCases := map[string]struct {
		input  string
		errStr TolinoErrStr
	}{
		"not enough entries":     {input: wrongHighlight, errStr: ErrNotEnoughEntries},
		"can't extract type":     {input: noType, errStr: ErrTypeExtraction},
		"not a page number":      {input: noPage, errStr: ErrNotEnoughEntries},
		"wrong timestamp layout": {input: wrongTime, errStr: ErrWrongTimeStamp},
	}
	for name, eC := range errorCases {
		t.Run(name, func(t *testing.T) {
			_, err := ExtractEntries(eC.input)
			if err == nil {
				t.Fatal("didn't get an error, wanted one")
			}

			toErr := err.(TolinoError)
			if !toErr.Is(eC.errStr) {
				t.Errorf("got error %q, wanted error %q", toErr, eC.errStr)
			}
		})
	}
}

func TestEntry_GetBook(t *testing.T) {
	testEntries := map[string]struct {
		entry Entry
		book  Book
	}{
		"correct entry":  {Entry{Title: "Book title", Author: "John Doe"}, Book{"Book title", "John Doe"}},
		"missing title":  {Entry{Title: "", Author: "John Doe"}, Book{"", "John Doe"}},
		"missing author": {Entry{Title: "Book title", Author: ""}, Book{"Book title", ""}},
	}
	for name, tC := range testEntries {
		t.Run(name, func(t *testing.T) {
			got := tC.entry.GetBook()
			if got.author != tC.book.author {
				t.Errorf("%q: didn't get correct author: got %q, wanted %q", name, got.author, tC.book.author)
			}

			if got.title != tC.book.title {
				t.Errorf("%q: didn't get correct title: got %q, wanted %q", name, got.title, tC.book.title)
			}
		})
	}
}

func TestExtractBookList(t *testing.T) {
	testCases := map[string]struct {
		entries []Entry
		want    []Book
	}{
		"name": {
			entries: []Entry{{Title: "t1", Author: "a1"}, {Title: "t2", Author: "a2"}, {Title: "t3", Author: "a3"}, {Title: "t4", Author: "a4"}},
			want:    []Book{{"t1", "a1"}, {"t2", "a2"}, {"t3", "a3"}, {"t4", "a4"}},
		},
	}
	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := ExtractBooks(tC.entries); !reflect.DeepEqual(got, tC.want) {
				t.Errorf("ExtractBookList() = %v, want %v", got, tC.want)
			}
		})
	}
}

func TestTolinoError_Is(t *testing.T) {
	toErr := TolinoError{err: ErrNotEnoughEntries}

	testCases := map[string]struct {
		errStr TolinoErrStr
		want   bool
	}{
		"correct TolinoErrorStr":   {errStr: ErrNotEnoughEntries, want: true},
		"different TolinoErrorStr": {errStr: ErrTypeExtraction, want: false},
	}

	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := toErr.Is(tC.errStr); got != tC.want {
				t.Errorf("got %v, wanted %v", got, tC.want)
			}
		})
	}
}

func TestTolinoError_Error(t *testing.T) {
	toErr := TolinoError{err: ErrTypeExtraction}

	testCases := map[string]struct {
		errStr string
		want   bool
	}{
		"correct Tolino error":   {errStr: string(ErrTypeExtraction), want: true},
		"different Tolino error": {errStr: string(ErrNotEnoughEntries), want: false},
		"not an error string":    {errStr: "some other string", want: false},
	}

	for name, tC := range testCases {
		t.Run(name, func(t *testing.T) {
			if got := toErr.Error() == tC.errStr; got != tC.want {
				t.Errorf("got %v, wanted %v", got, tC.want)
			}
		})
	}
}
