package tolino

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/exp/slices"
)

type Entry struct {
}

func extractEntries(input string) []string {
	const delim = "\n\n-----------------------------------\n\n"
	tmp := strings.Split(input, delim)
	tmp = tmp[:len(tmp)-1]

	ss := []string{}
	for _, t := range tmp {
		if entryType := extractType(t); entryType == "Note" || entryType == "Highlight" {
			ss = append(ss, t)
		}
	}

	return ss
}

func extractType(token string) string {
	pattern := regexp.MustCompile(`.*\n(\w+)\x{00A0}`)
	return pattern.FindStringSubmatch(token)[1]
}

func extractPage(token string) int {
	pattern := regexp.MustCompile(`^(\d+): `)
	page, err := strconv.Atoi(pattern.FindStringSubmatch(token)[1])
	if err != nil {
		return -1
	}
	return page
}

/*
const noteType = `44: That would be fixing your emotions!
"exercise control over our lives with the least possible effort"
Added on`

const highlightType = `126: "The 80/20 Principle treats time as a friend, not an enemy. Time gone is not time lost. Time will always come round again. This is why there are seven days in a week, twelve months in a year, why the seasons come round again. Insight and value are likely to come from placing ourselves in a comfortable, relaxed, and collaborative position toward time. It is our use of time, and not time itself, that is the enemy.
• The 80/20 Principle says that we should act less. Action drives out thought. It is because we have so much time that we squander it. The most productive time on a project is usually the last 20 percent, simply because the work has to be completed before a deadline. Productivity on most projects could be doubled simply by halving the amount of time for their completion. This is not evidence that time is in short supply."
Added on`
*/
func extractNote(token string) (page, note, highlight string) {
	// pattern := regexp.MustCompile(`(\d*): (.*) "(.*)"$`)
	// fmt.Printf("%q\n", token)
	// pattern := regexp.MustCompile(`(.*)\s*Added on$`)
	pattern := regexp.MustCompile(`(?s)^(\d+): (.*)"(.*)"\s+Added on$`)

	tmp := pattern.FindStringSubmatch(token)
	page, note, highlight = tmp[1], tmp[2], tmp[3]
	note = strings.TrimSpace(note)
	return
	// fmt.Println(len(tmp))
	// // fmt.Println(tmp)
	// fmt.Printf("0: %q\n", tmp[0])
	// fmt.Printf("1: %q\n", tmp[1])
	// fmt.Printf("2: %q\n", tmp[2])
	// fmt.Printf("2 trim: %q\n", strings.TrimSpace(tmp[2]))
	// fmt.Printf("3: %q\n", tmp[3])
	// fmt.Println()
	// return ""
}

func extractHighlight(token string) string {
	return ""
}

// fmt.Printf("entry: %q\n", extractType(tokens[id]))

func Foo(entry string) (string, error) {
	delim := "\u00a0" // U+00a0 (non-breaking space)
	var tokens []string

	tokens = strings.Split(entry, delim)
	// for i, s := range tokens {
	// fmt.Printf("%d: %q\n", i, s)
	// }

	// id := slices.Index(tokens, "on page") - 1

	// extractAuthor := func(token string) string {
	// 	pattern := regexp.MustCompile(`\((.*),\ (.*)\)`)
	// 	return strings.Join(pattern.FindStringSubmatch(token)[1:], " ")
	// }
	// fmt.Printf("author: %q\n", extractAuthor(tokens[id]))

	// extractTitle := func(token string) string {
	// 	pattern := regexp.MustCompile(`(.*) \(`)
	// 	return strings.Join(pattern.FindStringSubmatch(token)[1:], " ")
	// }
	// fmt.Printf("title: %q\n", extractTitle(tokens[id]))

	// a := strings.LastIndex(tokens[0], "\n")
	// entryType := tokens[0][a+1:]
	// title := strings.Split(tokens[0], " (")[0]
	// fmt.Println(entryType)
	// fmt.Println(title)

	var indeces = make(map[string]int)
	d := slices.Index(tokens, "on page")
	indeces["date"] = d
	// fmt.Println(indeces)
	indeces["page"] = slices.Index(tokens, "on page") + 1
	indeces[""] = slices.Index(tokens, "on page") + 1
	// for _, part := range strings.Split(entry, delim) {
	// tokens = append(tokens, strings.Split(part, "\n")...)
	// }
	// for i, s := range tokens {
	// 	fmt.Printf("%d: %q\n", i, s)
	// }

	return "", nil
}

func Baz() (string, error) {
	file, err := os.Open("../notes.txt")
	filename := "wrong file"
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
