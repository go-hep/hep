// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
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

	var buf *bytes.Buffer
	buf = bytes.NewBuffer(data[f.nbytesname:])
	if err := dir.UnmarshalROOT(buf); err != nil {
		return err
	}

	nk := 4 // Key::fNumberOfBytes
	buf = bytes.NewBuffer(data[nk:])
	dec := newDecoder(buf)
	var keyversion int16
	dec.readBin(&keyversion)
	if dec.err != nil {
		return dec.err
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

	buf = bytes.NewBuffer(data[nk:])
	dec = newDecoder(buf)
	classname := ""
	dec.readString(&classname)
	myprintf("class: [%v]\n", classname)

	cname := ""
	dec.readString(&cname)
	myprintf("cname: [%v]\n", cname)
	dir.named.name = cname

	title := ""
	dec.readString(&title)
	myprintf("title: [%v]\n", title)
	dir.named.title = title

	if dir.nbytesname < 10 || dir.nbytesname > 1000 {
		return fmt.Errorf("rootio: can't read directory info.")
	}

	return dec.err
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

	dec := newDecoder(bytes.NewBuffer(data))

	var nkeys int32
	dec.readInt32(&nkeys)
	for i := 0; i < int(nkeys); i++ {
		err = dir.readKey()
		if err != nil {
			return err
		}
	}
	return dec.err
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

func (dir *directory) UnmarshalROOT(data *bytes.Buffer) error {
	dec := newDecoder(data)

	var version int16
	dec.readBin(&version)
	myprintf("dir-version: %v\n", version)

	var ctime uint32
	dec.readBin(&ctime)
	dir.ctime = datime2time(ctime)
	myprintf("dir-ctime: %v\n", dir.ctime)

	var mtime uint32
	dec.readBin(&mtime)
	dir.mtime = datime2time(mtime)
	myprintf("dir-mtime: %v\n", dir.mtime)

	dec.readInt32(&dir.nbyteskeys)
	dec.readInt32(&dir.nbytesname)
	readptr := dec.readInt64
	if version <= 1000 {
		readptr = dec.readInt32
	}
	readptr(&dir.seekdir)
	readptr(&dir.seekparent)
	readptr(&dir.seekkeys)
	return dec.err
}

var _ Object = (*directory)(nil)
var _ Named = (*directory)(nil)
var _ Directory = (*directory)(nil)
var _ ROOTUnmarshaler = (*directory)(nil)
