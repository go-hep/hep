// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"io"
	"io/ioutil"
	"reflect"

	"github.com/pkg/errors"
)

// Record manages and describes blocks of data
type Record struct {
	unpack bool           // whether to unpack incoming/outcoming records
	blocks []Block        // connected blocks
	bmap   map[string]int // connected blocks

	w *Writer
	r *Reader

	cw Compressor
	xr Decompressor

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
		return errors.Errorf("rio: block [%s] already connected to record [%s]", name, rec.Name())
	}

	version := Version(0)
	switch t := ptr.(type) {
	case Streamer:
		version = t.RioVersion()
	}

	rec.bmap[name] = len(rec.blocks)
	rec.blocks = append(
		rec.blocks,
		newBlock(name, version),
	)
	rec.blocks[rec.bmap[name]].typ = reflect.TypeOf(ptr)

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
		err = block.raw.RioMarshal(xbuf)
		if err != nil {
			return errors.Errorf("rio: error writing block #%d (%s): %v", i, block.Name(), err)
		}
	}

	xlen := xbuf.Len()

	var cbuf *bytes.Buffer
	switch {
	case rec.Compress():
		cbuf = new(bytes.Buffer)
		switch {
		case rec.cw == nil:
			compr := rec.raw.Options.CompressorKind()
			cw, err := compr.NewCompressor(cbuf, rec.raw.Options)
			if err != nil {
				return err
			}
			rec.cw = cw
		default:
			err = rec.cw.Reset(cbuf)
			if err != nil {
				return err
			}
		}
		_, err = io.CopyBuffer(rec.cw, xbuf, make([]byte, 16*1024*1024))
		if err != nil {
			return errors.Errorf("rio: error compressing blocks: %v", err)
		}
		err = rec.cw.Flush()
		if err != nil {
			return errors.Errorf("rio: error compressing blocks: %v", err)
		}

	default:
		cbuf = xbuf
	}

	clen := cbuf.Len()

	rec.raw.Header.Len = uint32(clen)
	rec.raw.CLen = uint32(clen)
	rec.raw.XLen = uint32(xlen)

	buf := new(bytes.Buffer)
	err = rec.raw.RioMarshal(buf)
	if err != nil {
		return err
	}

	err = rec.w.writeRecord(rec, buf.Bytes(), cbuf.Bytes())

	return err
}

// Read reads data from the Reader, in the rio format
func (rec *Record) Read() error {
	return rec.readRecord(rec.r.r)
}

// readRecord reads data from the Reader r, in the rio format
func (rec *Record) readRecord(r io.Reader) error {
	err := rec.raw.RioUnmarshal(r)
	if err != nil {
		return err
	}

	clen := int64(rioAlignU32(rec.raw.CLen))
	if !rec.unpack {
		switch r := r.(type) {
		case io.Seeker:
			_, err = r.Seek(clen, 0)
		default:
			_, err = io.CopyN(ioutil.Discard, r, clen)
		}
		return err
	}

	return rec.readBlocks(r)
}

// readBlocks reads the blocks data from the Reader
func (rec *Record) readBlocks(r io.Reader) error {
	var err error
	clen := int64(rioAlignU32(rec.raw.CLen))

	lr := &io.LimitedReader{
		R: r,
		N: clen,
	}

	// decompression
	switch {
	case rec.xr == nil:
		compr := rec.raw.Options.CompressorKind()
		xr, err := compr.NewDecompressor(lr)
		if err != nil {
			return err
		}
		rec.xr = xr
		lr = &io.LimitedReader{
			R: xr,
			N: int64(rec.raw.XLen),
		}

	default:
		err = rec.xr.Reset(lr)
		if err != nil {
			return err
		}
		lr = &io.LimitedReader{
			R: rec.xr,
			N: int64(rec.raw.XLen),
		}
	}

	for lr.N > 0 {
		blk := newBlock("", 0)
		err = blk.raw.RioUnmarshal(lr)
		if err == io.EOF {
			err = nil
			break
		}
		if err != nil {
			return err
		}
		n := blk.Name()
		if i, ok := rec.bmap[n]; ok {
			rec.blocks[i] = blk
		} else {
			rec.bmap[n] = len(rec.blocks)
			rec.blocks = append(rec.blocks, blk)
		}
	}

	if lr.N > 0 {
		return errors.Errorf("rio: record read too few bytes (want=%d. got=%d)", clen, clen-lr.N)
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
	return CompressorKind((rec.raw.Options&gMaskCompr)>>16) != CompressNone
}

// Options returns the options of this record.
func (rec *Record) Options() Options {
	return rec.raw.Options
}

// EOF
