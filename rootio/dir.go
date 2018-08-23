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

// FIXME(sbinet): reorganize tdirectory/tdirectoryFile fields
// to closer match that of ROOT's.

type tdirectory struct {
	rvers      int16
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
	uuid  tuuid
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

	// make sure we have a streamer for this type.
	if _, err := dir.StreamerInfo(obj.Class()); err != nil {
		_, err = streamerInfoFrom(obj, dir)
		if err != nil {
			return errors.Wrapf(err, "rootio: could not generate streamer for key")
		}
		_, err = dir.StreamerInfo(obj.Class())
		if err != nil {
			panic(err)
		}
	}

	dir.keys = append(dir.keys, Key{
		f:        dir.file,
		version:  4, // FIXME(sbinet): harmonize versions
		datetime: nowUTC(),
		cycle:    cycle,
		class:    obj.Class(),
		name:     name,
		title:    title,
		obj:      obj,
		seekpdir: dir.seekdir,
	})

	return nil
}

func (dir *tdirectory) Keys() []Key {
	return dir.keys
}

func (dir *tdirectory) isBigFile() bool {
	return dir.rvers > 1000
}

func (dir *tdirectory) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	beg := w.Pos()

	w.WriteI16(dir.rvers)
	w.WriteU32(time2datime(dir.ctime))
	w.WriteU32(time2datime(dir.mtime))
	w.WriteI32(dir.nbyteskeys)
	w.WriteI32(dir.nbytesname)

	switch {
	case dir.isBigFile():
		w.WriteI64(dir.seekdir)
		w.WriteI64(dir.seekparent)
		w.WriteI64(dir.seekkeys)
	default:
		w.WriteI32(int32(dir.seekdir))
		w.WriteI32(int32(dir.seekparent))
		w.WriteI32(int32(dir.seekkeys))
	}

	dir.uuid.MarshalROOT(w)

	end := w.Pos()

	return int(end - beg), w.err
}

func (dir *tdirectory) UnmarshalROOT(r *RBuffer) error {
	var (
		version = r.ReadI16()
		ctime   = r.ReadU32()
		mtime   = r.ReadU32()
	)

	dir.rvers = version
	dir.ctime = datime2time(ctime)
	dir.mtime = datime2time(mtime)

	dir.nbyteskeys = r.ReadI32()
	dir.nbytesname = r.ReadI32()

	switch {
	case dir.isBigFile():
		dir.seekdir = r.ReadI64()
		dir.seekparent = r.ReadI64()
		dir.seekkeys = r.ReadI64()
	default:
		dir.seekdir = int64(r.ReadI32())
		dir.seekparent = int64(r.ReadI32())
		dir.seekkeys = int64(r.ReadI32())
	}

	dir.uuid.UnmarshalROOT(r)

	return r.Err()
}

// StreamerInfo returns the StreamerInfo with name of this directory, or nil otherwise.
func (dir *tdirectory) StreamerInfo(name string) (StreamerInfo, error) {
	if dir.file == nil {
		return nil, fmt.Errorf("rootio: no streamers")
	}
	return dir.file.StreamerInfo(name)
}

func (dir *tdirectory) addStreamer(streamer StreamerInfo) {
	dir.file.addStreamer(streamer)
}

type tdirectoryFile struct {
	dir tdirectory
}

func newDirectoryFile(name string, f *File) *tdirectoryFile {
	now := nowUTC()
	return &tdirectoryFile{tdirectory{
		rvers: 5, // FIXME(sbinet)
		ctime: now,
		mtime: now,
		named: tnamed{name: name},
		file:  f,
	}}
}

func (dir *tdirectoryFile) readKeys() error {
	return dir.dir.readKeys()
}

func (dir *tdirectoryFile) readDirInfo() error {
	return dir.dir.readDirInfo()
}

func (dir *tdirectoryFile) Close() error {
	if dir.dir.file.w == nil {
		return nil
	}

	// FIXME(sbinet): ROOT applies this optimization. should we ?
	//	if len(dir.dir.keys) == 0 || dir.dir.seekdir == 0 {
	//		return nil
	//	}

	err := dir.save()
	if err != nil {
		return err
	}

	return nil
}

func (dir *tdirectoryFile) save() error {
	var err error
	if dir.dir.file.w == nil {
		return err
	}

	for i := range dir.dir.keys {
		k := &dir.dir.keys[i]
		err = k.store()
		if err != nil {
			return err
		}
		_, err = k.writeFile(dir.dir.file)
		if err != nil {
			return errors.Wrapf(err, "rootio: could not write key for directory %q", dir.Name())
		}
	}

	err = dir.saveSelf()
	if err != nil {
		return errors.Wrapf(err, "rootio: could not save directory")
	}

	// FIXME(sbinet): recursively save sub-directories.

	return nil
}

func (dir *tdirectoryFile) saveSelf() error {
	if dir.dir.file.w == nil {
		return nil
	}

	var err error
	err = dir.writeKeys()
	if err != nil {
		return err
	}

	err = dir.writeDirHeader()
	if err != nil {
		return err
	}

	return nil
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

func (dir *tdirectoryFile) addStreamer(streamer StreamerInfo) {
	dir.dir.addStreamer(streamer)
}

func (dir *tdirectoryFile) MarshalROOT(w *WBuffer) (int, error) {
	return dir.dir.MarshalROOT(w)
}

func (dir *tdirectoryFile) UnmarshalROOT(r *RBuffer) error {
	return dir.dir.UnmarshalROOT(r)
}

// writeKeys writes the list of keys to the file.
// The list of keys is written out as a single data record.
func (dir *tdirectoryFile) writeKeys() error {
	var (
		err    error
		nbytes = int32(4) // space for n-keys
	)

	if dir.dir.file.end > kStartBigFile {
		nbytes += 8
	}
	for i := range dir.Keys() {
		key := &dir.dir.keys[i]
		nbytes += key.sizeof()
	}

	hdr := createKey(dir.Name(), dir.Title(), dir.Class(), nbytes, dir.dir.file)

	buf := NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.writeI32(int32(len(dir.Keys())))
	for _, k := range dir.Keys() {
		_, err = k.MarshalROOT(buf)
		if err != nil {
			return errors.Errorf("rootio: could not write key: %v", err)
		}
	}
	hdr.buf = buf.buffer()

	dir.dir.seekkeys = hdr.seekkey
	dir.dir.nbyteskeys = hdr.bytes

	_, err = hdr.writeFile(dir.dir.file)
	if err != nil {
		return errors.Errorf("rootio: could not write header key: %v", err)
	}
	return nil
}

// writeDirHeader overwrites the Directory header record.
func (dir *tdirectoryFile) writeDirHeader() error {
	var (
		err error
	)
	dir.dir.mtime = nowUTC()

	nbytes := dir.sizeof() + int32(dir.dir.file.nbytesname)
	key := newKey(dir.Name(), dir.Title(), "TFile", nbytes, dir.dir.file)
	key.seekkey = dir.dir.file.begin
	key.seekpdir = dir.dir.seekdir

	buf := NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.WriteString(dir.Name())
	buf.WriteString(dir.Title())
	_, err = dir.MarshalROOT(buf)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not marshal dir-info")
	}

	key.buf = buf.buffer()
	_, err = key.writeFile(dir.dir.file)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not write dir-info to file")
	}

	return nil
}

func (dir *tdirectoryFile) sizeof() int32 {
	nbytes := int32(22)

	nbytes += datimeSizeof() // ctime
	nbytes += datimeSizeof() // mtime
	nbytes += dir.dir.uuid.sizeof()
	if dir.dir.file.version >= 40000 {
		nbytes += 12 // files with >= 2Gb
	}
	return nbytes
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
			o := newDirectoryFile("", nil)
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
	_ streamerInfoStore   = (*tdirectory)(nil)
	_ ROOTMarshaler       = (*tdirectory)(nil)
	_ ROOTUnmarshaler     = (*tdirectory)(nil)

	_ Object              = (*tdirectoryFile)(nil)
	_ Named               = (*tdirectoryFile)(nil)
	_ Directory           = (*tdirectoryFile)(nil)
	_ StreamerInfoContext = (*tdirectoryFile)(nil)
	_ streamerInfoStore   = (*tdirectoryFile)(nil)
	_ ROOTMarshaler       = (*tdirectoryFile)(nil)
	_ ROOTUnmarshaler     = (*tdirectoryFile)(nil)
)
