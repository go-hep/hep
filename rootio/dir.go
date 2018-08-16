// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

type tdirectory struct {
	ctime      time.Time // time of directory's creation
	mtime      time.Time // time of directory's last modification
	nbyteskeys int32     // number of bytes for the keys
	nbytesname int32     // number of bytes in TNamed at creation time
	seekdir    int64     // location of directory on file
	seekparent int64     // location of parent directory on file
	seekkeys   int64     // location of Keys record on file

	classname string

	named tnamed // name+title of this directory
	file  *File  // pointer to current file in memory
	keys  []Key
}

// recordSize returns the size of the directory header in bytes
func (dir *tdirectory) recordSize(version int32) int64 {
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

func (dir *tdirectory) readDirInfo() error {
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

	r := NewRBuffer(data[f.nbytesname:], nil, 0, nil)
	if err := dir.UnmarshalROOT(r); err != nil {
		return err
	}

	nk := 4 // Key::fNumberOfBytes
	r = NewRBuffer(data[nk:], nil, 0, nil)
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

	r = NewRBuffer(data[nk:], nil, 0, nil)
	dir.classname = r.ReadString()

	dir.named.name = r.ReadString()
	dir.named.title = r.ReadString()

	if dir.nbytesname < 10 || dir.nbytesname > 1000 {
		return fmt.Errorf("rootio: can't read directory info.")
	}

	return r.Err()
}

func (dir *tdirectory) readKeys() error {
	var err error
	if dir.seekkeys <= 0 {
		return nil
	}

	buf := make([]byte, int(dir.nbyteskeys))
	_, err = dir.file.ReadAt(buf, dir.seekkeys)
	if err != nil {
		return err
	}

	hdr := Key{f: dir.file}
	err = hdr.UnmarshalROOT(NewRBuffer(buf, nil, 0, dir))
	if err != nil {
		return err
	}

	buf = make([]byte, hdr.objlen)
	_, err = dir.file.ReadAt(buf, dir.seekkeys+int64(hdr.keylen))
	if err != nil {
		return err
	}

	r := NewRBuffer(buf, nil, 0, dir)
	nkeys := r.ReadI32()
	if r.Err() != nil {
		return r.err
	}
	dir.keys = make([]Key, int(nkeys))
	for i := range dir.keys {
		k := &dir.keys[i]
		k.f = dir.file
		err := k.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dir *tdirectory) Class() string {
	return "TDirectory"
}

func (dir *tdirectory) Name() string {
	return dir.named.Name()
}

func (dir *tdirectory) Title() string {
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
func (dir *tdirectory) Get(namecycle string) (Object, error) {
	name, cycle := decodeNameCycle(namecycle)
	for i := range dir.keys {
		k := &dir.keys[i]
		if k.Name() == name {
			if cycle != 9999 {
				if k.cycle == cycle {
					return k.Value().(Object), nil
				}
				continue
			}
			return k.Value().(Object), nil
		}
	}
	return nil, noKeyError{key: namecycle, obj: dir}
}

func (dir *tdirectory) Put(name string, obj Object) error {
	var (
		cycle int16
		title = ""
	)
	if name == "" {
		if v, ok := obj.(Named); ok {
			name = v.Name()
			title = v.Title()
		}
	}
	if name == "" {
		return errors.Errorf("rootio: empty key name")
	}

	// FIXME(sbinet): implement a fast look-up ?
	for i := range dir.keys {
		key := &dir.keys[i]
		if key.name != name {
			continue
		}
		cycle = key.cycle
	}
	cycle++

	dir.keys = append(dir.keys, Key{
		f:     dir.file,
		cycle: cycle,
		class: obj.Class(),
		name:  name,
		title: title,
		obj:   obj,
	})
	panic("not implemented")
}

func (dir *tdirectory) Keys() []Key {
	return dir.keys
}

func (dir *tdirectory) UnmarshalROOT(r *RBuffer) error {
	var (
		version = r.ReadI16()
		ctime   = r.ReadU32()
		mtime   = r.ReadU32()
	)

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

// StreamerInfo returns the StreamerInfo with name of this directory, or nil otherwise.
func (dir *tdirectory) StreamerInfo(name string) (StreamerInfo, error) {
	if dir.file == nil {
		return nil, fmt.Errorf("rootio: no streamers")
	}
	return dir.file.StreamerInfo(name)
}

type tdirectoryFile struct {
	dir tdirectory
}

func (dir *tdirectoryFile) Get(namecycle string) (Object, error) {
	return dir.dir.Get(namecycle)
}

func (dir *tdirectoryFile) Put(name string, v Object) error {
	return dir.dir.Put(name, v)
}

func (dir *tdirectoryFile) Keys() []Key {
	return dir.dir.Keys()
}

func (dir *tdirectoryFile) Class() string {
	return "TDirectoryFile"
}

func (dir *tdirectoryFile) Name() string {
	return dir.dir.named.Name()
}

func (dir *tdirectoryFile) Title() string {
	return dir.dir.named.Title()
}

func (dir *tdirectoryFile) StreamerInfo(name string) (StreamerInfo, error) {
	return dir.dir.StreamerInfo(name)
}

func (dir *tdirectoryFile) UnmarshalROOT(r *RBuffer) error {
	err := dir.dir.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	{
		f := func() reflect.Value {
			o := &tdirectory{}
			return reflect.ValueOf(o)
		}
		Factory.add("TDirectory", f)
		Factory.add("*rootio.tdirectory", f)
	}
	{
		f := func() reflect.Value {
			o := &tdirectoryFile{}
			return reflect.ValueOf(o)
		}
		Factory.add("TDirectoryFile", f)
		Factory.add("*rootio.tdirectoryFile", f)
	}
}

var (
	_ Object              = (*tdirectory)(nil)
	_ Named               = (*tdirectory)(nil)
	_ Directory           = (*tdirectory)(nil)
	_ StreamerInfoContext = (*tdirectory)(nil)
	_ ROOTUnmarshaler     = (*tdirectory)(nil)

	_ Object              = (*tdirectoryFile)(nil)
	_ Named               = (*tdirectoryFile)(nil)
	_ Directory           = (*tdirectoryFile)(nil)
	_ StreamerInfoContext = (*tdirectoryFile)(nil)
	_ ROOTUnmarshaler     = (*tdirectoryFile)(nil)
)
