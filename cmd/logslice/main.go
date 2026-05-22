package main

import (
	"fmt"
	"os"
	"time"

	"github.com/user/logslice/internal/reader"
	"github.com/user/logslice/internal/slicer"
	"github.com/user/logslice/internal/timeparse"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := parseFlags()
	if err != nil {
		return err
	}

	r, err := reader.New(cfg.inputFile)
	if err != nil {
		return fmt.Errorf("opening input: %w", err)
	}
	defer r.Close()

	parser := timeparse.NewParser(cfg.timeFormat)

	var out *os.File
	if cfg.outputFile == "" {
		out = os.Stdout
	} else {
		out, err = os.Create(cfg.outputFile)
		if err != nil {
			return fmt.Errorf("creating output: %w", err)
		}
		defer out.Close()
	}

	result, err := slicer.Slice(r, parser, cfg.from, cfg.to, cfg.maxLines)
	if err != nil {
		return fmt.Errorf("slicing: %w", err)
	}

	return slicer.Write(out, result, cfg.showCount)
}

type config struct {
	inputFile  string
	outputFile string
	timeFormat string
	from       time.Time
	to         time.Time
	maxLines   int
	showCount  bool
}
