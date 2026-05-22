package slicer

import (
	"bufio"
	"io"

	"github.com/user/logslice/internal/timeparse"
)

// Options configures the slicing behaviour.
type Options struct {
	Parser    *timeparse.Parser
	From      string
	To        string
	MaxLines  int // 0 means unlimited
}

// Result holds the output of a slice operation.
type Result struct {
	Lines []string
	Count int
}

// Slice reads lines from r and returns only those whose timestamp falls
// within [from, to]. Lines without a recognisable timestamp are skipped.
func Slice(r io.Reader, opts Options) (*Result, error) {
	result := &Result{}
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		t, ok := opts.Parser.Extract(line)
		if !ok {
			continue
		}

		in, err := opts.Parser.InRange(t, opts.From, opts.To)
		if err != nil {
			return nil, err
		}
		if !in {
			continue
		}

		result.Lines = append(result.Lines, line)
		result.Count++

		if opts.MaxLines > 0 && result.Count >= opts.MaxLines {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
