// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bufio"
	"bytes"
	"io"

	riobin "github.com/gonuts/binary"
)

// Reader is a rio read-only stream
type Reader struct {
	r io.Reader

	rrr     *bytes.Buffer
	options Options
	version Version

	recs map[string]*Record // map of all connected records
}

// NewReader returns a new read-only rio stream
func NewReader(r io.Reader) (*Reader, error) {
	// a rio stream starts with rio magic
	hdr := [4]byte{}
	_, err := r.Read(hdr[:])
	if err != nil {
		return nil, errorf("rio: error reading magic-header: %v", err)
	}
	if hdr != rioMagic {
		return nil, errorf("rio: not a rio-stream. magic-header=%q. want=%q",
			string(hdr[:]),
			string(rioMagic[:]),
		)
	}

	buf := new(bytes.Buffer)
	r = io.TeeReader(r, buf)
	rr := bufio.NewReader(r) //Size(r, 4)
	return &Reader{
		r:       rr,
		options: 0,
		version: rioHdrVersion,
		recs:    make(map[string]*Record),
		rrr:     buf,
	}, nil
}

// Record adds a Record to the list of records to read or
// returns the Record with that name.
func (r *Reader) Record(name string) *Record {
	rec, ok := r.recs[name]
	if !ok {
		rec = newRecord(name, r.options)
		rec.r = r
		rec.unpack = true
		r.recs[name] = rec
	}
	return rec
}

// Records returns the list of connected Records
func (r *Reader) Records() []*Record {
	recs := make([]*Record, 0, len(r.recs))
	for _, rec := range r.recs {
		recs = append(recs, rec)
	}
	return recs
}

// Close finishes reading the rio read-only stream.
// It does not (and can not) close the underlying reader.
func (r *Reader) Close() error {
	var err error
	return err
}

// decoder manages the decoding of data values from rioRecords
type decoder struct {
	r io.Reader
}

func (dec *decoder) Decode(v interface{}) error {
	switch v := v.(type) {
	case RioDecoder:
		return v.RioDecode(dec.r)
	}

	d := riobin.NewDecoder(dec.r)
	d.Order = Endian
	return d.Decode(v)
}
