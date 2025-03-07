// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"fmt"
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
		return nil, fmt.Errorf("rio: error reading magic-header: %w", err)
	}
	if hdr != rioMagic {
		return nil, fmt.Errorf("rio: not a rio-stream. magic-header=%q. want=%q",
			string(hdr[:]),
			string(rioMagic[:]),
		)
	}

	// a seek-able rio streams sports a rioFooter at the end.
	_, err = f.r.Seek(-int64(ftrSize), io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("rio: error seeking footer: %w", err)
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
		return nil, fmt.Errorf("rio: error seeking metadata: %w", err)
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
	keys := make([]RecordDesc, len(f.meta.Records))
	copy(keys, f.meta.Records)
	sort.Sort(recordsByName(keys))
	return keys
}

// Get reads the value `name` into `ptr`
func (f *File) Get(name string, ptr any) error {
	offsets, ok := f.meta.Offsets[name]
	if !ok {
		return fmt.Errorf("rio: no record [%s]", name)
	}

	if len(offsets) > 1 {
		return fmt.Errorf("rio: multi-record streams unsupported")
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
