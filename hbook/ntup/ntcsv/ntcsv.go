// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntcsv provides a convenient access to CSV files as n-tuple data.
//
// Examples:
//
//  nt, err := ntcsv.Open("testdata/simple.csv")
//  if err != nil {
//      log.Fatal(err)
//  }
//  defer nt.DB().Close()
//
// or, with a different configuration for the comma/comment runes:
//
//  nt, err := ntcsv.Open("testdata/simple.csv", ntcsv.Comma(' '), ntcsv.Comment('#'))
//  if err != nil {
//      log.Fatal(err)
//  }
//  defer nt.DB().Close()
//
// Give our own names to the CSV columns (default: "var1", "var2", ...):
//
//  nt, err := ntcsv.Open("testdata/simple.csv", ntcsv.Columns("var1", "i64", "foo"))
//
// Take the names from the CSV header (note that the header *must* exist):
//
//  nt, err := ntcsv.Open("testdata/simple-with-header.csv", ntcsv.Header())
//
// Override the names from the CSV header with our own:
//
//  nt, err := ntcsv.Open("testdata/simple-with-header.csv", ntcsv.Header(), ntcsv.Columns("v1", "v2", "v3")
package ntcsv // import "go-hep.org/x/hep/hbook/ntup/ntcsv"

import (
	"fmt"

	"go-hep.org/x/hep/csvutil/csvdriver"
	"go-hep.org/x/hep/hbook/ntup"
)

// Open opens a CSV file in read-only mode and returns a n-tuple
// connected to that.
func Open(name string, opts ...Option) (*ntup.Ntuple, error) {
	c := csvdriver.Conn{File: name}
	for _, opt := range opts {
		opt(&c)
	}

	db, err := c.Open()
	if err != nil {
		return nil, err
	}

	nt, err := ntup.Open(db, "csv")
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("could not open n-tuple: %w", err)
	}

	return nt, nil
}

// Option configures the underlying sql.DB connection to the n-tuple.
type Option func(c *csvdriver.Conn)

// Comma configures the n-tuple to use v as the comma delimiter between columns.
func Comma(v rune) Option {
	return func(c *csvdriver.Conn) {
		c.Comma = v
	}
}

// Comment configures the n-tuple to use v as the comment character
// for start of line.
func Comment(v rune) Option {
	return func(c *csvdriver.Conn) {
		c.Comment = v
	}
}

// Header informs the n-tuple the CSV file has a header line.
func Header() Option {
	return func(c *csvdriver.Conn) {
		c.Header = true
	}
}

// Columns names the n-tuple columns with the given slice.
func Columns(names ...string) Option {
	return func(c *csvdriver.Conn) {
		if len(names) == 0 {
			return
		}
		c.Names = make([]string, len(names))
		copy(c.Names, names)
	}
}
