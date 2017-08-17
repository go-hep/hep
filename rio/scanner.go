// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"fmt"
	"io"
	"io/ioutil"
)

// Selector selects Records based on their name
type Selector struct {
	Name   string // Record name
	Unpack bool   // Whether to unpack the Record
}

// Scanner provides a convenient interface for reading records of a rio-stream.
type Scanner struct {
	r   *Reader
	err error   // first non-EOF error encountered while reading the rio-stream.
	rec *Record // last record encountered while reading the rio-stream.

	filter map[string]Selector // records to read. if nil, return everything.
}

// NewScanner returns a new Scanner to read from r.
func NewScanner(r *Reader) *Scanner {
	scan := &Scanner{
		r:      r,
		err:    nil,
		rec:    newRecord("<N/A>", 0),
		filter: make(map[string]Selector, 0),
	}
	scan.rec.unpack = false
	scan.rec.r = r
	return scan
}

// Select sets the records selection function.
func (s *Scanner) Select(selectors []Selector) {
	s.filter = make(map[string]Selector, len(selectors))
	for _, sel := range selectors {
		s.filter[sel.Name] = sel
	}
}

// Scan scans the next Record until io.EOF
func (s *Scanner) Scan() bool {
	if s.err != nil {
		return false
	}

	for {
		var hdr rioHeader
		err := hdr.RioUnmarshal(s.r.r)
		if err != nil {
			s.err = err
			return false
		}

		switch hdr.Frame {
		case ftrFrame:
			ftr := rioFooter{Header: hdr}
			err = ftr.unmarshalData(s.r.r)
			if err != nil {
				s.err = err
				return false
			}
			continue
		case recFrame:
			s.rec.raw.Header = hdr
			err := s.rec.raw.unmarshalData(s.r.r)
			if err != nil {
				s.err = err
				return false
			}

			clen := int64(rioAlignU32(s.rec.raw.CLen))

			name := s.rec.Name()
			if len(s.filter) > 0 {
				_, ok := s.filter[name]
				if !ok {
					_, err = s.seek(clen, 0)
					if err != nil {
						s.err = err
						return false
					}
					continue
				}
			}

			s.rec.unpack = s.filter[name].Unpack

			switch s.rec.unpack {

			case true:
				err = s.rec.readBlocks(s.r.r)
				if err != nil {
					s.err = err
					return false
				}
				return true

			case false:
				_, err = s.seek(clen, 0)
				if err != nil {
					s.err = err
					return false
				}
				return true
			}

		default:
			panic(fmt.Errorf("unknown frame %v", hdr.Frame))
		}
	}
}

// seek sets the offset for the next Read or Write on file to offset,
// interpreted according to whence: 0 means relative to the origin of the
// file, 1 means relative to the current offset, and 2 means relative to
// the end. It returns the new offset and an error, if any.
func (s *Scanner) seek(offset int64, whence int) (ret int64, err error) {
	switch r := s.r.r.(type) {
	case io.Seeker:
		return r.Seek(offset, whence)
	default:
		if whence != 0 {
			panic("not implemented")
		}
		return io.CopyN(ioutil.Discard, r, offset)
	}
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
