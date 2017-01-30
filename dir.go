// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"time"
)

type directory struct {
	ctime      time.Time // time of directory's creation
	mtime      time.Time // time of directory's last modification
	nbyteskeys int32     // number of bytes for the keys
	nbytesname int32     // number of bytes in TNamed at creation time
	seekdir    int64     // location of directory on file
	seekparent int64     // location of parent directory on file
	seekkeys   int64     // location of Keys record on file

	named named // name+title of this directory
	file  *File // pointer to current file in memory
	keys  []Key
}

// recordSize returns the size of the directory header in bytes
func (dir *directory) recordSize(version int32) int64 {
	var nbytes int64
	nbytes += 2 // fVersion
	nbytes += 4 // ctime
	nbytes += 4 // mtime
	nbytes += 4 // nbyteskeys
	nbytes += 4 // nbytesname
	if version >= 40000 {
		// assume that the file may be above 2 Gbytes if file version is > 4
		nbytes += 8 // seekdir
		nbytes += 8 // seekparent
		nbytes += 8 // seekkeys
	} else {
		nbytes += 4 // seekdir
		nbytes += 4 // seekparent
		nbytes += 4 // seekkeys
	}
	return nbytes
}

func (dir *directory) readDirInfo() error {
	f := dir.file
	nbytes := int64(f.nbytesname) + dir.recordSize(f.version)

	if nbytes+f.begin > f.end {
		return fmt.Errorf(
			"rootio: file [%v] has an incorrect header length [%v] or incorrect end of file length [%v]",
			f.id,
			f.begin+nbytes,
			f.end,
		)
	}

	data := make([]byte, int(nbytes))
	if _, err := f.ReadAt(data, f.begin); err != nil {
		return err
	}

	r := NewRBuffer(data[f.nbytesname:], nil, 0)
	if err := dir.UnmarshalROOT(r); err != nil {
		return err
	}

	nk := 4 // Key::fNumberOfBytes
	r = NewRBuffer(data[nk:], nil, 0)
	keyversion := r.ReadI16()
	if r.Err() != nil {
		return r.Err()
	}

	if keyversion > 1000 {
		// large files
		nk += 2     // Key::fVersion
		nk += 2 * 4 // Key::fObjectSize, Date
		nk += 2 * 2 // Key::fKeyLength, fCycle
		nk += 2 * 8 // Key::fSeekKey, fSeekParentDirectory
	} else {
		nk += 2     // Key::fVersion
		nk += 2 * 4 // Key::fObjectSize, Date
		nk += 2 * 2 // Key::fKeyLength, fCycle
		nk += 2 * 4 // Key::fSeekKey, fSeekParentDirectory
	}

	r = NewRBuffer(data[nk:], nil, 0)
	classname := r.ReadString()

	dir.named.name = r.ReadString()
	dir.named.title = r.ReadString()

	myprintf("class: [%v]\n", classname)
	myprintf("cname: [%v]\n", dir.named.name)
	myprintf("title: [%v]\n", dir.named.title)

	if dir.nbytesname < 10 || dir.nbytesname > 1000 {
		return fmt.Errorf("rootio: can't read directory info.")
	}

	return r.Err()
}

func (dir *directory) readKeys() error {
	var err error
	if dir.seekkeys <= 0 {
		return nil
	}

	_, err = dir.file.Seek(dir.seekkeys, io.SeekStart)
	if err != nil {
		return err
	}

	hdr := Key{f: dir.file}
	err = hdr.Read()
	if err != nil {
		return err
	}
	//myprintf("==> hdr: %#v\n", hdr)

	_, err = dir.file.Seek(dir.seekkeys+int64(hdr.keylen), io.SeekStart)
	if err != nil {
		return err
	}
	data := make([]byte, 4)
	_, err = dir.file.Read(data)
	if err != nil {
		return err
	}

	r := NewRBuffer(data, nil, 0)

	nkeys := r.ReadI32()
	for i := 0; i < int(nkeys); i++ {
		err = dir.readKey()
		if err != nil {
			return err
		}
	}
	return r.Err()
}

// readKey reads a key and appends it to dir.keys
func (dir *directory) readKey() error {
	dir.keys = append(dir.keys, Key{f: dir.file})
	key := &(dir.keys[len(dir.keys)-1])
	return key.Read()
}

func (dir *directory) Class() string {
	return "TDirectory"
}

func (dir *directory) Name() string {
	return dir.named.Name()
}

func (dir *directory) Title() string {
	return dir.named.Title()
}

// Get returns the object identified by namecycle
//   namecycle has the format name;cycle
//   name  = * is illegal, cycle = * is illegal
//   cycle = "" or cycle = 9999 ==> apply to a memory object
//
//   examples:
//     foo   : get object named foo in memory
//             if object is not in memory, try with highest cycle from file
//     foo;1 : get cycle 1 of foo on file
func (dir *directory) Get(namecycle string) (Object, bool) {
	name, cycle := decodeNameCycle(namecycle)
	for _, k := range dir.keys {
		if k.Name() == name {
			if cycle != 9999 {
				if k.cycle == cycle {
					return &k, true
				} else {
					return nil, false
				}
			}
			return &k, true
		}
	}
	return nil, false
}

func (dir *directory) UnmarshalROOT(r *RBuffer) error {
	var (
		version = r.ReadI16()
		ctime   = r.ReadU32()
		mtime   = r.ReadU32()
	)
	myprintf("dir-version: %v\n", version)
	myprintf("dir-ctime: %v\n", dir.ctime)
	myprintf("dir-mtime: %v\n", dir.mtime)

	dir.mtime = datime2time(mtime)
	dir.ctime = datime2time(ctime)

	dir.nbyteskeys = r.ReadI32()
	dir.nbytesname = r.ReadI32()

	readptr := r.ReadI64
	if version <= 1000 {
		readptr = func() int64 { return int64(r.ReadI32()) }
	}
	dir.seekdir = readptr()
	dir.seekparent = readptr()
	dir.seekkeys = readptr()
	return r.Err()
}

var _ Object = (*directory)(nil)
var _ Named = (*directory)(nil)
var _ Directory = (*directory)(nil)
var _ ROOTUnmarshaler = (*directory)(nil)
