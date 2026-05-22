// Package slicer provides the core log-slicing logic for logslice.
//
// It consumes an io.Reader (plain or gzip-wrapped via the reader package),
// filters log lines whose timestamps fall within a caller-supplied [from, to]
// window using the timeparse package, and writes the matching lines to any
// io.Writer.
//
// Basic usage:
//
//	p, _ := timeparse.NewParser("")
//	res, err := slicer.Slice(r, slicer.Options{
//		Parser: p,
//		From:   "2024-01-10T10:00:00Z",
//		To:     "2024-01-10T11:00:00Z",
//	})
//	if err != nil { ... }
//	slicer.Write(os.Stdout, res, slicer.WriteOptions{ShowCount: true})
package slicer
