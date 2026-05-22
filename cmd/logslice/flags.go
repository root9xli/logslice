package main

import (
	"errors"
	"flag"
	"fmt"
	"time"
)

const defaultTimeFormat = ""

func parseFlags() (*config, error) {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)

	var (
		input      = fs.String("f", "", "input log file (plain or .gz)")
		output     = fs.String("o", "", "output file (default: stdout)")
		fmt_       = fs.String("fmt", defaultTimeFormat, "timestamp format (Go layout or empty for auto-detect)")
		fromStr    = fs.String("from", "", "start time (RFC3339 or free-form)")
		toStr      = fs.String("to", "", "end time (RFC3339 or free-form)")
		maxLines   = fs.Int("n", 0, "maximum number of lines to output (0 = unlimited)")
		showCount  = fs.Bool("count", false, "print matched line count to stderr")
	)

	if err := fs.Parse(flag.Args()); err != nil {
		return nil, err
	}

	if *input == "" {
		return nil, errors.New("flag -f (input file) is required")
	}

	var from, to time.Time
	var err error

	if *fromStr != "" {
		from, err = parseTime(*fromStr)
		if err != nil {
			return nil, fmt.Errorf("invalid -from value: %w", err)
		}
	}

	if *toStr != "" {
		to, err = parseTime(*toStr)
		if err != nil {
			return nil, fmt.Errorf("invalid -to value: %w", err)
		}
	}

	if !from.IsZero() && !to.IsZero() && !to.After(from) {
		return nil, errors.New("-to must be after -from")
	}

	return &config{
		inputFile:  *input,
		outputFile: *output,
		timeFormat: *fmt_,
		from:       from,
		to:         to,
		maxLines:   *maxLines,
		showCount:  *showCount,
	}, nil
}

func parseTime(s string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as a timestamp", s)
}
