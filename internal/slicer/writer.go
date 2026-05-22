package slicer

import (
	"bufio"
	"fmt"
	"io"
)

// WriteOptions controls how results are written.
type WriteOptions struct {
	ShowCount bool
}

// Write writes the sliced lines to w, optionally appending a summary line.
func Write(w io.Writer, result *Result, opts WriteOptions) error {
	bw := bufio.NewWriter(w)

	for _, line := range result.Lines {
		if _, err := fmt.Fprintln(bw, line); err != nil {
			return err
		}
	}

	if opts.ShowCount {
		if _, err := fmt.Fprintf(bw, "# %d line(s) matched\n", result.Count); err != nil {
			return err
		}
	}

	return bw.Flush()
}
