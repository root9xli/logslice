// Package timeparse provides utilities for parsing timestamps from log lines.
package timeparse

import (
	"errors"
	"time"
)

// Common log timestamp formats to attempt when parsing.
var knownFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05.000",
	"2006-01-02 15:04:05.000000",
	"02/Jan/2006:15:04:05 -0700", // Apache combined log
	"Jan _2 15:04:05",            // syslog
	"Jan  2 15:04:05",
}

// ErrNoTimestamp is returned when no recognisable timestamp is found in a line.
var ErrNoTimestamp = errors.New("no timestamp found in line")

// Parser holds configuration for timestamp extraction.
type Parser struct {
	formats []string
	location *time.Location
}

// NewParser creates a Parser that tries the given formats in order.
// If formats is nil the built-in knownFormats list is used.
func NewParser(formats []string, loc *time.Location) *Parser {
	if formats == nil {
		formats = knownFormats
	}
	if loc == nil {
		loc = time.UTC
	}
	return &Parser{formats: formats, location: loc}
}

// Extract attempts to parse a timestamp from the beginning of line.
// It returns the parsed time and the byte offset after the timestamp.
func (p *Parser) Extract(line string) (time.Time, int, error) {
	for _, fmt := range p.formats {
		if len(line) < len(fmt) {
			continue
		}
		// Try progressively longer prefixes up to twice the format length.
		maxLen := len(fmt) + 10
		if maxLen > len(line) {
			maxLen = len(line)
		}
		for end := len(fmt); end <= maxLen; end++ {
			t, err := time.ParseInLocation(fmt, line[:end], p.location)
			if err == nil {
				return t, end, nil
			}
		}
	}
	return time.Time{}, 0, ErrNoTimestamp
}

// InRange reports whether t falls within [from, to] (inclusive).
func InRange(t, from, to time.Time) bool {
	return !t.Before(from) && !t.After(to)
}
