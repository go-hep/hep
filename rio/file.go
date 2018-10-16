// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"io"
	"sort"
)

// File random-read-access to a rio stream
type File struct {
	r    io.ReadSeeker
	meta Metadata
}

// Open creates a new read-only File.
func Open(r io.ReadSeeker) (*File, error) {

	f := &File{
		r: r,
	}

	// a rio stream starts with rio magic
	hdr := [4]byte{}
	_, err := f.r.Read(hdr[:])
	if err != nil {
		return nil, errorf("rio: error reading magic-header: %v", err)
	}
	if hdr != rioMagic {
		return nil, errorf("rio: not a rio-stream. magic-header=%q. want=%q",
			string(hdr[:]),
			string(rioMagic[:]),
		)
	}

	// a seek-able rio streams sports a rioFooter at the end.
	_, err = f.r.Seek(-int64(ftrSize), io.SeekEnd)
	if err != nil {
		return nil, errorf("rio: error seeking footer (err=%v)", err)
	}

	// {
	// 	fmt.Printf("==== tail ==== (%d)\n", ftrSize)
	// 	buf := new(bytes.Buffer)
	// 	io.Copy(buf, f.r)
	// 	fmt.Printf("buf: %v\n", buf.Bytes())
	// 	pos, err = f.r.Seek(-int64(ftrSize), 2)
	// 	fmt.Printf("=== [tail] ===\n")
	// }

	var ftr rioFooter
	err = ftr.RioUnmarshal(f.r)
	if err != nil {
		return nil, err
	}

	_, err = f.r.Seek(ftr.Meta, io.SeekStart)
	if err != nil {
		return nil, errorf("rio: error seeking metadata (err=%v)", err)
	}

	rec := newRecord(MetaRecord, 0)
	rec.unpack = true

	err = rec.readRecord(f.r)
	if err != nil {
		return nil, err
	}

	err = rec.Block(MetaRecord).Read(&f.meta)
	if err != nil {
		return nil, err
	}

	return f, err
}

// Keys returns the list of record names.
func (f *File) Keys() []RecordDesc {
	keys := make([]RecordDesc, 0, len(f.meta.Records))
	for _, rec := range f.meta.Records {
		keys = append(keys, rec)
	}
	sort.Sort(recordsByName(keys))
	return keys
}

// Get reads the value `name` into `ptr`
func (f *File) Get(name string, ptr interface{}) error {
	offsets, ok := f.meta.Offsets[name]
	if !ok {
		return errorf("rio: no record [%s]", name)
	}

	if len(offsets) > 1 {
		return errorf("rio: multi-record streams unsupported")
	}

	offset := offsets[0]
	_, err := f.r.Seek(offset.Pos, 0)
	if err != nil {
		return err
	}

	rec := newRecord(name, 0)
	rec.unpack = true

	err = rec.readRecord(f.r)
	if err != nil {
		return err
	}

	err = rec.Block(name).Read(ptr)
	if err != nil {
		return err
	}

	return err
}

// Has returns whether a record `name` exists in this file.
func (f *File) Has(name string) bool {
	_, ok := f.meta.Offsets[name]
	return ok
}

// Close closes access to the rio-file.
// It does not (and can not) close the underlying reader.
func (f *File) Close() error {
	return nil
}
