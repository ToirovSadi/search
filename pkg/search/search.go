package search

import (
	"context"
	"io/ioutil"
	"strings"
	"sync"
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

//All find all phrases from files and return them
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	wg := sync.WaitGroup{}
	ch := make(chan []Result)
	_, cancel := context.WithCancel(ctx)
	for i := 0; i < len(files); i++ {
		file := files[i]
		wg.Add(1)
		go func(file string, ch chan []Result) {
			defer wg.Done()
			words := search(file, phrase, false)
			if len(words) > 0 {
				ch <- words
			}
		}(file, ch)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()
	cancel()
	return ch
}

func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	wg := sync.WaitGroup{}
	ch := make(chan Result)
	findResult := Result{}
	_, cancel := context.WithCancel(ctx)
	for i := 0; i < len(files); i++ {
		file := files[i]
		wg.Add(1)
		go func(file string, ch chan Result) {
			defer wg.Done()
			words := search(file, phrase, true)
			if len(words) > 0 {
				findResult = words[0]
			}
		}(file, ch)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if findResult != (Result{}) {
			ch <- findResult
		}
		close(ch)
		cancel()
	}()
	wg.Wait()
	return ch
}

func search(fileName string, phrase string, first bool) (res []Result) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(data), "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		givenLine := line
		for strings.Contains(line, phrase) {
			idx := strings.Index(line, phrase)
			res = append(res, Result{
				Phrase:  phrase,
				Line:    givenLine,
				LineNum: int64(i + 1),
				ColNum:  int64(idx + 1),
			})
			if first == true {
				return res
			}
			line = strings.Replace(line, phrase, string('#')+phrase[1:], 1)
		}
	}
	return res
}
