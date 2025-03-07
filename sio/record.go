// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
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
	name    string         // record name
	unpack  bool           // whether to unpack incoming records
	options uint32         // options (flag word)
	bindex  map[string]int // index of connected blocks
	bnames  []string       // connected blocks names
	blocks  []Block        // connected blocks
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

// Disconnect disconnects all blocks previously connected to this
// Record (for reading or writing.)
func (rec *Record) Disconnect() {
	rec.bnames = rec.bnames[:0]
	rec.bindex = make(map[string]int)
	rec.blocks = rec.blocks[:0]
}

// Connect connects a Block to this Record (for reading or writing)
func (rec *Record) Connect(name string, ptr any) error {
	var err error
	iblk, ok := rec.bindex[name]
	if !ok {
		iblk = len(rec.blocks)
		rec.bnames = append(rec.bnames, name)
		rec.blocks = append(rec.blocks, nil)
		rec.bindex[name] = iblk
		//return fmt.Errorf("sio.Record: Block name [%s] already connected", name)
		//return ErrBlockConnected
	}
	var block Block
	switch ptr := ptr.(type) {
	case Block:
		block = ptr
	case Codec:
		rt := reflect.TypeOf(ptr)
		var vers uint32
		if ptr, ok := ptr.(Versioner); ok {
			vers = ptr.VersionSio()
		}
		block = &userBlock{
			blk:     ptr,
			version: vers,
			name:    rt.Name(),
		}

	default:
		rt := reflect.TypeOf(ptr)
		if rt.Kind() != reflect.Ptr {
			return fmt.Errorf("sio: Connect needs a pointer to a block of data")
		}
		var vers uint32
		if ptr, ok := ptr.(Versioner); ok {
			vers = ptr.VersionSio()
		}
		block = &genericBlock{
			rt:      rt,
			rv:      reflect.ValueOf(ptr),
			version: vers,
			name:    rt.Name(),
		}
	}
	rec.blocks[iblk] = block
	return err
}

// read reads a record
func (rec *Record) read(r *reader) error {
	var err error
	// fmt.Printf("::: reading record [%s]... [%d]\n", rec.name, r.Len())
	type fixlink struct {
		link Linker
		vers uint32
	}
	var linkers []fixlink
	// loop until data has been depleted
	for r.Len() > 0 {
		beg := r.Len()
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
			return ErrBlockShortRead
		}
		iblk, ok := rec.bindex[string(cbuf.Bytes()[:data.NameLen])]
		if ok {
			blk := rec.blocks[iblk]
			// fmt.Printf("### %q\n", buf.String())
			err = blk.UnmarshalSio(r)
			end := r.Len()
			if err != nil {
				// fmt.Printf("*** error unmarshaling record=%q block=%q: %v\n", rec.name, name, err)
				return err
			}
			if beg-end != int(hdr.Len) {
				/*
					if true {
						var typ any
						switch blk := blk.(type) {
						case *userBlock:
							typ = blk.blk
						case *genericBlock:
							typ = blk.rv.Interface()
						}
						log.Printf("record %q block %q (%T) (beg-end=%d-%d=%d != %d)", rec.Name(), name, typ, beg, end, beg-end, int(hdr.Len))
					} else {
				*/
				return ErrBlockShortRead
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

	// fmt.Printf("::: reading record [%s]... [done]\n", rec.name)
	return err
}

func (rec *Record) write(w *writer) error {
	var (
		err  error
		work = make([]byte, 16*1024*1024)
	)
	for i, k := range rec.bnames {
		blk := rec.blocks[i]
		bhdr := blockHeader{
			Typ: blkMarker,
		}

		bdata := blockData{
			Version: blk.VersionSio(),
			NameLen: uint32(len(k)),
		}

		wblk := newWriterFrom(w)
		wblk.ver = bdata.Version

		err = blk.MarshalSio(wblk)
		if err != nil {
			return err
		}

		bhdr.Len = uint32(blockHeaderSize) + uint32(blockDataSize) +
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

		_, err := io.CopyBuffer(w, wblk.buf, work)
		if err != nil {
			return err
		}
		w.ids = wblk.ids
		w.tag = wblk.tag
		w.ptr = wblk.ptr
	}
	return err
}
