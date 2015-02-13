// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"io"
)

// Scanner provides a convenient interface for reading records of a rio-stream.
type Scanner struct {
	r   *Reader
	err error   // first non-EOF error encountered while reading the rio-stream.
	rec *Record // last record encountered while reading the rio-stream.
}

// NewScanner returns a new Scanner to read from r.
func NewScanner(r *Reader) *Scanner {
	scan := &Scanner{
		r:   r,
		err: nil,
		rec: newRecord("<N/A>", 0),
	}
	scan.rec.unpack = false
	scan.rec.r = r
	return scan
}

// Scan scans the next Record until io.EOF
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	// FIXME(sbinet): what happens when different record-shapes and record.unpack==true?
	// should we allocate anew?
	err := s.rec.Read()
	if err != nil {
		s.err = err
		return false
	}

	return true
}

// Err returns the first non-EOF error encountered by the reader.
func (s *Scanner) Err() error {
	if s.err == io.EOF {
		return nil
	}
	return s.err
}

// Record returns the last Record read by the Scanner.
func (s *Scanner) Record() *Record {
	return s.rec
}
