// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

// recordHeader describes the on-disk record (header part)
type recordHeader struct {
	Len uint32
	Typ uint32
}

// recordData describes the on-disk record (payload part)
type recordData struct {
	Options uint32
	DataLen uint32 // length of compressed record data
	UCmpLen uint32 // length of uncompressed record data
	NameLen uint32 // length of record name
}

// Record manages blocks of data
type Record struct {
	name    string           // record name
	unpack  bool             // whether to unpack incoming records
	options uint32           // options (flag word)
	blocks  map[string]Block // connected blocks
}

// Name returns the name of this record
func (rec *Record) Name() string {
	return rec.name
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
	return rec.options&optCompress != 0
}

// SetCompress sets or resets the compression flag
func (rec *Record) SetCompress(compress bool) {
	rec.options &= optNotCompress
	if compress {
		rec.options |= optCompress
	}
}

// Options returns the options of this record.
func (rec *Record) Options() uint32 {
	return rec.options
}

// Connect connects a Block to this Record (for reading or writing)
func (rec *Record) Connect(name string, ptr interface{}) error {
	var err error
	_, dup := rec.blocks[name]
	if dup {
		//return fmt.Errorf("sio.Record: Block name [%s] already connected", name)
		//return ErrBlockConnected
	}
	var block Block
	switch ptr := ptr.(type) {
	case Block:
		block = ptr
	case Codec:
		rt := reflect.TypeOf(ptr)
		block = &userBlock{
			blk:     ptr,
			version: 0,
			name:    rt.Name(),
		}

	default:
		rt := reflect.TypeOf(ptr)
		if rt.Kind() != reflect.Ptr {
			return fmt.Errorf("sio: Connect needs a pointer to a block of data")
		}
		block = &genericBlock{
			rt:      rt,
			rv:      reflect.ValueOf(ptr),
			version: 0,
			name:    rt.Name(),
		}
	}
	rec.blocks[name] = block
	return err
}

// read reads a record
func (rec *Record) read(r *reader) error {
	var err error
	// fmt.Printf("::: reading record [%s]... [%d]\n", rec.name, buf.Len())
	type fixlink struct {
		link Linker
		vers uint32
	}
	var linkers []fixlink
	// loop until data has been depleted
	for r.Len() > 0 {
		// read block header
		var hdr blockHeader
		err = bread(r, &hdr)
		if err != nil {
			return err
		}
		if hdr.Typ != blkMarker {
			// fmt.Printf("*** err record[%s]: noblockmarker\n", rec.name)
			return ErrRecordNoBlockMarker
		}

		var data blockData
		err = bread(r, &data)
		if err != nil {
			return err
		}
		r.ver = data.Version

		var cbuf bytes.Buffer
		nlen := align4U32(data.NameLen)
		n, err := io.CopyN(&cbuf, r, int64(nlen))
		if err != nil {
			// fmt.Printf(">>> err:%v\n", err)
			return err
		}
		if n != int64(nlen) {
			return fmt.Errorf("sio: read too few bytes (got=%d. expected=%d)", n, nlen)
		}
		name := string(cbuf.Bytes()[:data.NameLen])
		blk, ok := rec.blocks[name]
		if ok {
			// fmt.Printf("### %q\n", string(buf.Bytes()))
			err = blk.UnmarshalSio(r)
			if err != nil {
				// fmt.Printf("*** error unmarshaling record=%q block=%q: %v\n", rec.name, name, err)
				return err
			}
			// fmt.Printf(">>> read record=%q block=%q (buf=%d)\n", rec.name, name, buf.Len())
			if ublk, ok := blk.(*userBlock); ok {
				if link, ok := ublk.blk.(Linker); ok {
					linkers = append(linkers, fixlink{link, data.Version})
				}
			}
		}

		// check whether there is still something to be read.
		// if there is, check whether there is a block-marker
		if r.Len() > 0 {
			next := bytes.Index(r.Bytes(), blkMarkerBeg)
			if next > 0 {
				pos := next - 4 // sizeof mark-block
				r.Next(pos)     // drain the buffer until next block
			} else {
				// drain the whole buffer
				r.Next(r.Len())
			}
		}
	}
	r.relocate()
	for _, fix := range linkers {
		err = fix.link.LinkSio(fix.vers)
		if err != nil {
			return err
		}
	}

	//fmt.Printf("::: reading record [%s]... [done]\n", rec.name)
	return err
}

func (rec *Record) write(w *writer) error {
	var err error
	for k, blk := range rec.blocks {

		bhdr := blockHeader{
			Typ: blkMarker,
		}

		bdata := blockData{
			Version: blk.Version(),
			NameLen: uint32(len(k)),
		}

		wblk := newWriterFrom(w)
		wblk.ver = bdata.Version

		err = blk.MarshalSio(wblk)
		if err != nil {
			return err
		}

		bhdr.Len = uint32(unsafe.Sizeof(bhdr)) +
			uint32(unsafe.Sizeof(bdata)) +
			align4U32(bdata.NameLen) + uint32(wblk.Len())

		// fmt.Printf("blockHeader: %v\n", bhdr)
		// fmt.Printf("blockData:   %v (%s)\n", bdata, k)

		err = bwrite(w, &bhdr)
		if err != nil {
			return err
		}

		err = bwrite(w, &bdata)
		if err != nil {
			return err
		}

		_, err = w.Write([]byte(k))
		if err != nil {
			return err
		}
		padlen := align4U32(bdata.NameLen) - bdata.NameLen
		if padlen > 0 {
			_, err = w.Write(make([]byte, int(padlen)))
			if err != nil {
				return err
			}
		}

		_, err := io.Copy(w, wblk.buf)
		if err != nil {
			return err
		}
		w.ids = wblk.ids
		w.tag = wblk.tag
		w.ptr = wblk.ptr
	}
	return err
}
