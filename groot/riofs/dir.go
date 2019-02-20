// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"reflect"
	"time"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type tdirectory struct {
	rvers  int16
	named  rbase.Named // name+title of this directory
	parent Directory
	objs   []root.Object
	uuid   rbase.UUID
}

func (dir *tdirectory) RVersion() int16 { return dir.rvers }

func (dir *tdirectory) Class() string {
	return "TDirectory"
}

func (dir *tdirectory) Name() string {
	return dir.named.Name()
}

func (dir *tdirectory) Title() string {
	return dir.named.Title()
}

func (dir *tdirectory) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(dir.RVersion())

	dir.named.MarshalROOT(w)
	switch dir.parent {
	case nil:
		w.WriteObjectAny((*rbase.Object)(nil))
	default:
		w.WriteObjectAny(dir.parent.(root.Object))
	}
	// FIXME(sbinet): stream list
	dir.uuid.MarshalROOT(w)

	return w.SetByteCount(pos, dir.Class())
}

func (dir *tdirectory) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(dir.Class())

	dir.rvers = vers

	dir.named.UnmarshalROOT(r)
	obj := r.ReadObjectAny()
	if obj != nil {
		dir.parent = obj.(Directory)
	}
	// FIXME(sbinet): stream list
	dir.uuid.UnmarshalROOT(r)

	r.CheckByteCount(pos, bcnt, start, dir.Class())
	return r.Err()
}

type tdirectoryFile struct {
	dir tdirectory

	ctime      time.Time // time of directory's creation
	mtime      time.Time // time of directory's last modification
	nbyteskeys int32     // number of bytes for the keys
	nbytesname int32     // number of bytes in TNamed at creation time
	seekdir    int64     // location of directory on file
	seekparent int64     // location of parent directory on file
	seekkeys   int64     // location of Keys record on file

	classname string

	file *File // pointer to current file in memory
	keys []Key
}

func newDirectoryFile(name string, f *File, parent *tdirectoryFile) *tdirectoryFile {
	now := nowUTC()
	dir := &tdirectoryFile{
		dir: tdirectory{
			rvers: rvers.DirectoryFile,
			named: *rbase.NewNamed(name, name),
		},
		ctime: now,
		mtime: now,
		file:  f,
	}
	if parent != nil {
		dir.dir.parent = parent
	}
	return dir
}

func (dir *tdirectoryFile) isBigFile() bool {
	return dir.dir.rvers > 1000
}

// recordSize returns the size of the directory header in bytes
func (dir *tdirectoryFile) recordSize(version int32) int64 {
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

func (dir *tdirectoryFile) readDirInfo() error {
	f := dir.file
	nbytes := int64(f.nbytesname) + dir.recordSize(f.version)

	if nbytes+f.begin > f.end {
		return fmt.Errorf(
			"riofs: file [%v] has an incorrect header length [%v] or incorrect end of file length [%v]",
			f.id,
			f.begin+nbytes,
			f.end,
		)
	}

	data := make([]byte, int(nbytes))
	if _, err := f.ReadAt(data, f.begin); err != nil {
		return err
	}

	r := rbytes.NewRBuffer(data[f.nbytesname:], nil, 0, nil)
	if err := dir.UnmarshalROOT(r); err != nil {
		return err
	}

	nk := 4 // Key::fNumberOfBytes
	r = rbytes.NewRBuffer(data[nk:], nil, 0, nil)
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

	r = rbytes.NewRBuffer(data[nk:], nil, 0, nil)
	dir.classname = r.ReadString()

	dir.dir.named.SetName(r.ReadString())
	dir.dir.named.SetTitle(r.ReadString())

	if dir.nbytesname < 10 || dir.nbytesname > 1000 {
		return fmt.Errorf("riofs: can't read directory info.")
	}

	return r.Err()
}

func (dir *tdirectoryFile) readKeys() error {
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
	err = hdr.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, dir))
	if err != nil {
		return err
	}

	buf = make([]byte, hdr.objlen)
	_, err = dir.file.ReadAt(buf, dir.seekkeys+int64(hdr.keylen))
	if err != nil {
		return err
	}

	r := rbytes.NewRBuffer(buf, nil, 0, dir)
	nkeys := r.ReadI32()
	if r.Err() != nil {
		return r.Err()
	}
	dir.keys = make([]Key, int(nkeys))
	for i := range dir.keys {
		k := &dir.keys[i]
		k.f = dir.file
		err := k.UnmarshalROOT(r)
		if err != nil {
			return err
		}
		// support old ROOT versions.
		if k.class == "TDirectory" {
			k.class = "TDirectoryFile"
		}
	}
	return nil
}

func (dir *tdirectoryFile) Close() error {
	if dir.file.w == nil {
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
	if dir.file.w == nil {
		return err
	}

	for i := range dir.keys {
		k := &dir.keys[i]
		err = k.store()
		if err != nil {
			return err
		}

		switch obj := k.obj.(type) {
		case *tdirectoryFile:
			err = obj.save()
			if err != nil {
				return err
			}
		}
	}

	err = dir.saveSelf()
	if err != nil {
		return errors.Wrapf(err, "riofs: could not save directory")
	}

	return nil
}

func (dir *tdirectoryFile) saveSelf() error {
	if dir.file.w == nil {
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

// Get returns the object identified by namecycle
//   namecycle has the format name;cycle
//   name  = * is illegal, cycle = * is illegal
//   cycle = "" or cycle = 9999 ==> apply to a memory object
//
//   examples:
//     foo   : get object named foo in memory
//             if object is not in memory, try with highest cycle from file
//     foo;1 : get cycle 1 of foo on file
func (dir *tdirectoryFile) Get(namecycle string) (root.Object, error) {
	name, cycle := decodeNameCycle(namecycle)
	for i := range dir.keys {
		k := &dir.keys[i]
		if k.Name() == name {
			if cycle != 9999 {
				if k.cycle == cycle {
					return k.Object()
				}
				continue
			}
			return k.Object()
		}
	}
	return nil, noKeyError{key: namecycle, obj: dir}
}

func (dir *tdirectoryFile) Put(name string, obj root.Object) error {
	var (
		cycle int16
		title = ""
	)
	if v, ok := obj.(root.Named); ok {
		if name == "" {
			name = v.Name()
		}
		title = v.Title()
	}
	if name == "" {
		return errors.Errorf("riofs: empty key name")
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

	typename := obj.Class()

	// make sure we have a streamer for this type.
	if !isCoreType(typename) {
		cxx := rdict.GoName2Cxx(typename)
		si, err := dir.StreamerInfo(cxx, -1)
		if err != nil {
			si, err = streamerInfoFrom(obj, dir)
			if err != nil {
				return errors.Wrapf(err, "riofs: could not generate streamer for key %q and type %T", name, obj)
			}
			si, err = dir.StreamerInfo(cxx, -1)
		}
		if err != nil {
			return errors.Wrapf(err, "riofs: could not find streamer for %T", obj)
		}
		dir.addStreamer(si)
	}

	dir.keys = append(dir.keys, Key{
		f:        dir.file,
		rvers:    rvers.Key,
		datetime: nowUTC(),
		cycle:    cycle,
		class:    rdict.GoName2Cxx(typename),
		name:     name,
		title:    title,
		obj:      obj,
		seekpdir: dir.seekdir,
	})

	return nil
}

// Keys returns the list of keys being held by this directory.
func (dir *tdirectoryFile) Keys() []Key {
	return dir.keys
}

// Mkdir creates a new subdirectory
func (dir *tdirectoryFile) Mkdir(name string) (Directory, error) {
	if _, err := dir.Get(name); err == nil {
		return nil, errors.Errorf("rootio: %q already exist", name)
	}

	sub := newDirectoryFile(name, dir.file, dir)
	err := dir.Put(name, sub)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (dir *tdirectoryFile) RVersion() int16 { return dir.dir.rvers }

func (dir *tdirectoryFile) Class() string {
	return "TDirectoryFile"
}

func (dir *tdirectoryFile) Name() string {
	return dir.dir.named.Name()
}

func (dir *tdirectoryFile) Title() string {
	return dir.dir.named.Title()
}

func (dir *tdirectoryFile) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	beg := w.Pos()

	w.WriteI16(dir.RVersion())
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

	dir.dir.uuid.MarshalROOT(w)

	end := w.Pos()

	return int(end - beg), w.Err()
}

func (dir *tdirectoryFile) UnmarshalROOT(r *rbytes.RBuffer) error {
	var (
		version = r.ReadI16()
		ctime   = r.ReadU32()
		mtime   = r.ReadU32()
	)

	dir.dir.rvers = version
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

	if r.Len() != 0 {
		dir.dir.uuid.UnmarshalROOT(r)
	}

	return r.Err()
}

// StreamerInfo returns the StreamerInfo with name of this directory, or nil otherwise.
// If version is negative, the latest version should be returned.
func (dir *tdirectoryFile) StreamerInfo(name string, version int) (rbytes.StreamerInfo, error) {
	if dir.file == nil {
		return nil, fmt.Errorf("riofs: no streamers")
	}
	return dir.file.StreamerInfo(name, version)
}

func (dir *tdirectoryFile) addStreamer(streamer rbytes.StreamerInfo) {
	dir.file.addStreamer(streamer)
}

// writeKeys writes the list of keys to the file.
// The list of keys is written out as a single data record.
func (dir *tdirectoryFile) writeKeys() error {
	var (
		err    error
		nbytes = int32(4) // space for n-keys
	)

	if dir.file.end > kStartBigFile {
		nbytes += 8
	}
	for i := range dir.Keys() {
		key := &dir.keys[i]
		nbytes += key.sizeof()
	}

	if dir.seekdir <= 0 {
		nbytes := dir.sizeof() + int32(dir.file.nbytesname)
		blk := dir.file.spans.best(int64(nbytes))
		dir.seekdir = blk.first
	}

	hdr := createKey(dir.Name(), dir.Title(), dir.Class(), nbytes, dir.file)

	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.WriteI32(int32(len(dir.Keys())))
	for _, k := range dir.Keys() {
		_, err = k.MarshalROOT(buf)
		if err != nil {
			return errors.Errorf("riofs: could not write key: %v", err)
		}
	}
	hdr.buf = buf.Bytes()

	dir.seekkeys = hdr.seekkey
	dir.nbyteskeys = hdr.nbytes

	for i := range dir.keys {
		k := &dir.keys[i]
		k.seekpdir = dir.seekdir

		k.buf = nil // force re-computation of serialized key
		err = k.store()
		if err != nil {
			return err
		}

		_, err = k.writeFile(dir.file)
		if err != nil {
			return errors.Wrapf(err, "riofs: could not write sub-key")
		}
	}

	_, err = hdr.writeFile(dir.file)
	if err != nil {
		return errors.Errorf("riofs: could not write header key: %v", err)
	}
	return nil
}

// writeDirHeader overwrites the Directory header record.
func (dir *tdirectoryFile) writeDirHeader() error {
	var (
		err error
	)
	dir.mtime = nowUTC()

	nbytes := dir.sizeof() + int32(dir.file.nbytesname)
	key := newKey(dir.Name(), dir.Title(), "TFile", nbytes, dir.file)
	key.seekkey = dir.file.begin
	key.seekpdir = dir.seekdir

	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.WriteString(dir.Name())
	buf.WriteString(dir.Title())
	_, err = dir.MarshalROOT(buf)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not marshal dir-info")
	}

	key.buf = buf.Bytes()
	_, err = key.writeFile(dir.file)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write dir-info to file")
	}

	return nil
}

func (dir *tdirectoryFile) sizeof() int32 {
	nbytes := int32(22)

	nbytes += datimeSizeof() // ctime
	nbytes += datimeSizeof() // mtime
	nbytes += dir.dir.uuid.Sizeof()
	if dir.file.version >= 40000 {
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
		rtypes.Factory.Add("TDirectory", f)
	}
	{
		f := func() reflect.Value {
			o := newDirectoryFile("", nil, nil)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TDirectoryFile", f)
	}
}

// coreTypes is the set of types that do not neet a streamer info.
var coreTypes = map[string]struct{}{
	"TObject":        {},
	"TFile":          {},
	"TDirectoryFile": {},
	"TKey":           {},
}

func isCoreType(typename string) bool {
	_, ok := coreTypes[typename]
	return ok
}

var (
	_ root.Object        = (*tdirectory)(nil)
	_ root.Named         = (*tdirectory)(nil)
	_ rbytes.Marshaler   = (*tdirectory)(nil)
	_ rbytes.Unmarshaler = (*tdirectory)(nil)

	_ root.Object                = (*tdirectoryFile)(nil)
	_ root.Named                 = (*tdirectoryFile)(nil)
	_ Directory                  = (*tdirectoryFile)(nil)
	_ rbytes.StreamerInfoContext = (*tdirectoryFile)(nil)
	_ streamerInfoStore          = (*tdirectoryFile)(nil)
	_ rbytes.Marshaler           = (*tdirectoryFile)(nil)
	_ rbytes.Unmarshaler         = (*tdirectoryFile)(nil)
)
