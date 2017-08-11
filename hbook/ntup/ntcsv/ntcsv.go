// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ntcsv provides a convenient access to CSV files as n-tuple data.
//
// Example:
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
package ntcsv // import "go-hep.org/x/hep/hbook/ntup/ntcsv"

import (
	"os"

	"go-hep.org/x/hep/csvutil/csvdriver"
	"go-hep.org/x/hep/hbook/ntup"
)

// Open opens a CSV file in read-only mode and returns a n-tuple
// connected to that.
func Open(name string, opts ...Option) (*ntup.Ntuple, error) {
	c := conn{
		c: csvdriver.Conn{
			File:    name,
			Mode:    os.O_RDONLY,
			Perm:    0,
			Comma:   ',',
			Comment: '#',
		},
		header: false,
	}

	for _, opt := range opts {
		opt(&c)
	}

	db, err := c.c.Open()
	if err != nil {
		return nil, err
	}

	return ntup.Open(db, "csv")
}

type conn struct {
	c      csvdriver.Conn
	header bool // whether the CSV file has a header
}

// Option configures the underlying sql.DB connection to the n-tuple.
type Option func(c *conn)

// Comma configures the n-tuple to use v as the comma delimiter between columns.
func Comma(v rune) Option {
	return func(c *conn) {
		c.c.Comma = v
	}
}

// Comment configures the n-tuple to use v as the comment character
// for start of line.
func Comment(v rune) Option {
	return func(c *conn) {
		c.c.Comment = v
	}
}
