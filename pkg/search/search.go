package search

import (
	"context"
	"strings"
)

// Result describes result of search
type Result struct {
	// phrase you're looking for
	Phrase string
	// the whole line with "Phrase"
	Line string
	// line number of phrase
	LineNum int64
	// column number of phrase
	ColNum int64
}

func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	//wg := sync.WaitGroup{}

	return ch
}

func Search(line string, phrase string, index int64, first bool) (res []Result) {
	for strings.Contains(line, phrase) {
		idx := strings.Index(line, phrase)
		res = append(res, Result{
			Phrase:  phrase,
			Line:    line,
			LineNum: index,
			ColNum:  int64(idx),
		})
		if first == true {
			return res
		}
		line = strings.Replace(line, phrase, string('#')+phrase[1:], 1)
	}
	return res
}
