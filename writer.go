// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bufio"
	"compress/flate"
	"io"

	riobin "github.com/gonuts/binary"
)

type cwriter struct {
	w *bufio.Writer
	n int64
}

func (w *cwriter) Write(data []byte) (int, error) {
	n, err := w.w.Write(data)
	w.n += int64(n)
	return n, err
}

func (w *cwriter) Flush() error {
	return w.w.Flush()
}

// Writer is a rio write-only stream
type Writer struct {
	w *cwriter

	options Options
	version Version

	recs    map[string]*Record
	offsets map[string][]int64
	closed  bool
}

// NewWriter returns a new write-only rio stream
func NewWriter(w io.Writer) (*Writer, error) {
	ww := &cwriter{bufio.NewWriter(w), 0}
	// a rio stream starts with rio magic.
	_, err := ww.Write(rioMagic[:])
	if err != nil {
		return nil, err
	}

	return &Writer{
		w:       ww,
		options: NewOptions(CompressDefault, flate.DefaultCompression, 0),
		version: 1,
		recs:    make(map[string]*Record),
		offsets: make(map[string][]int64),
	}, nil
}

// SetCompressor enables compression and sets the compression method.
func (w *Writer) SetCompressor(compr CompressorKind, lvl int) error {
	var err error

	// FIXME(sbinet) handle codec (gob|cbor|xdr|riobin|...)
	codec := 0
	w.options = NewOptions(compr, lvl, codec)

	return err
}

// Record adds a Record to the list of records to write or
// returns the Record with that name.
func (w *Writer) Record(name string) *Record {
	rec, ok := w.recs[name]
	if !ok {
		rec = newRecord(name, w.options)
		rec.w = w
		w.recs[name] = rec
	}
	return rec
}

// Close finishes writing the rio write-only stream.
// It does not (and can not) close the underlying writer.
func (w *Writer) Close() error {
	if w.closed {
		return nil
	}
	w.closed = true
	pos := w.w.n
	var meta Metadata
	for _, rec := range w.recs {
		var blocks []struct{ Name, Type string }
		for _, blk := range rec.blocks {
			blocks = append(blocks, struct{ Name, Type string }{blk.Name(), nameFromType(blk.typ)})
		}
		meta.Records = append(meta.Records, struct {
			Name   string
			Blocks []struct{ Name, Type string }
		}{
			Name:   rec.Name(),
			Blocks: blocks,
		})
	}
	meta.Offsets = w.offsets

	err := w.WriteValue(rioMeta, &meta)
	if err != nil {
		return err
	}

	ftr := rioFooter{
		Header: rioHeader{
			Len:   uint32(ftrSize),
			Frame: ftrFrame,
		},
		Meta: pos,
	}
	err = ftr.RioMarshal(w.w)
	if err != nil {
		return err
	}
	return w.w.Flush()
}

// writeRecord writes all the record data
func (w *Writer) writeRecord(rec *Record, hdr, data []byte) error {
	var err error
	w.offsets[rec.Name()] = append(w.offsets[rec.Name()], w.w.n)

	_, err = w.w.Write(hdr)
	if err != nil {
		return err
	}

	_, err = w.w.Write(data)
	if err != nil {
		return err
	}

	n := rioAlignU32(rec.raw.Header.Len)
	if n != rec.raw.Header.Len {
		_, err = w.w.Write(make([]byte, int(n-rec.raw.Header.Len)))
	}

	return err
}

// WriteValue writes a value to the stream
func (w *Writer) WriteValue(name string, value interface{}) error {
	var err error

	rec := w.Record(name)
	err = rec.Connect(name, value)
	if err != nil {
		return err
	}

	blk := rec.Block(name)
	err = blk.Write(value)
	if err != nil {
		return err
	}

	err = rec.Write()
	if err != nil {
		return err
	}

	return err
}

// encoder manages the encoding of data values into rioRecords
type encoder struct {
	w io.Writer
}

func (enc *encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case Marshaler:
		return v.RioMarshal(enc.w)
	}

	e := riobin.NewEncoder(enc.w)
	e.Order = Endian
	return e.Encode(v)
}
