// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"io"
	"io/ioutil"
)

// Record manages and describes blocks of data
type Record struct {
	unpack bool           // whether to unpack incoming/outcoming records
	blocks []Block        // connected blocks
	bmap   map[string]int // connected blocks

	w *Writer
	r *Reader

	raw rioRecord
}

func newRecord(name string, options Options) *Record {
	rec := Record{
		unpack: false,
		blocks: make([]Block, 0, 2),
		bmap:   make(map[string]int, 2),
		raw: rioRecord{
			Header: rioHeader{
				Len:   0,
				Frame: recFrame,
			},
			Options: options,
			Name:    name,
		},
	}
	return &rec
}

// Connect connects a Block to this Record (for reading or writing)
func (rec *Record) Connect(name string, ptr interface{}) error {
	_, dup := rec.bmap[name]
	if dup {
		return errorf("rio: block [%s] already connected to record [%s]", name, rec.Name())
	}

	version := Version(0)
	switch t := ptr.(type) {
	case RioStreamer:
		version = t.RioVersion()
	}

	rec.bmap[name] = len(rec.blocks)
	rec.blocks = append(
		rec.blocks,
		newBlock(name, version),
	)

	return nil
}

// Block returns the block named name for reading or writing
// Block returns nil if the block doesn't exist
func (rec *Record) Block(name string) *Block {
	i, ok := rec.bmap[name]
	if !ok {
		return nil
	}
	block := &rec.blocks[i]
	return block
}

// Write writes data to the Writer, in the rio format
func (rec *Record) Write() error {
	var err error
	xbuf := new(bytes.Buffer) // FIXME(sbinet): use a sync.Pool

	for i := range rec.blocks {
		block := &rec.blocks[i]
		err = block.raw.RioEncode(xbuf)
		if err != nil {
			return errorf("rio: error writing block #%d (%s): %v", i, block.Name(), err)
		}
	}

	// FIXME(sbinet): handle compression
	cbuf := xbuf

	rec.raw.Header.Len = uint64(cbuf.Len())
	rec.raw.CLen = uint64(cbuf.Len())
	rec.raw.XLen = uint64(xbuf.Len())

	buf := new(bytes.Buffer)
	err = rec.raw.RioEncode(buf)
	if err != nil {
		return err
	}

	_, err = rec.w.w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = rec.w.w.Write(cbuf.Bytes())
	if err != nil {
		return err
	}

	n := rioAlignU64(rec.raw.Header.Len)
	if n != rec.raw.Header.Len {
		_, err = rec.w.w.Write(make([]byte, int(n-rec.raw.Header.Len)))
	}

	return err
}

// Read reads data from the Reader, in the rio format
func (rec *Record) Read() error {
	err := rec.raw.RioDecode(rec.r.r)
	if err != nil {
		return err
	}

	clen := int64(rioAlignU64(rec.raw.CLen))
	if !rec.unpack {
		switch r := rec.r.r.(type) {
		case io.Seeker:
			_, err = r.Seek(clen, 0)
		default:
			_, err = io.CopyN(ioutil.Discard, r, clen)
		}
		return err
	}

	r := &io.LimitedReader{
		R: rec.r.r,
		N: clen,
	}
	for i := range rec.blocks {
		block := &rec.blocks[i]
		err = block.raw.RioDecode(r)
		if err != nil {
			return err
		}
	}

	if r.N > 0 {
		return errorf("rio: record read too few bytes (want=%d. got=%d)", clen, clen-r.N)
	}
	return err
}

// Name returns the name of this record
func (rec *Record) Name() string {
	return rec.raw.Name
}

// Unpack returns whether to unpack incoming records
func (rec *Record) Unpack() bool {
	return rec.unpack
}

// SetUnpack sets whether to unpack incoming records
func (rec *Record) SetUnpack(unpack bool) {
	rec.unpack = unpack
}

// Compress returns the compression flag
func (rec *Record) Compress() bool {
	return rec.raw.Options&gCompress != 0
}

// SetCompress sets or resets the compression flag
func (rec *Record) SetCompress(compress bool) {
	rec.raw.Options &= gNotCompress
	if compress {
		rec.raw.Options |= gCompress
	}
}

// Options returns the options of this record.
func (rec *Record) Options() Options {
	return rec.raw.Options
}

// EOF
