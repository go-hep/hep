// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"time"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
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

	if _, err := dir.named.MarshalROOT(w); err != nil {
		return 0, w.Err()
	}
	if err := w.WriteObjectAny((*rbase.Object)(nil)); err != nil {
		return 0, w.Err()
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
	dirs []*tdirectoryFile
}

func newDirectoryFile(name, title string, f *File, parent *tdirectoryFile) *tdirectoryFile {
	now := nowUTC()
	if title == "" {
		title = name
	}
	dir := &tdirectoryFile{
		dir: tdirectory{
			rvers: rvers.DirectoryFile,
			named: *rbase.NewNamed(name, title),
		},
		ctime: now,
		mtime: now,
		file:  f,
	}
	if parent == nil {
		return dir
	}

	dir.dir.parent = parent
	dir.seekparent = parent.seekdir

	objlen := int32(dir.recordSize(f.version))
	key := newKey(parent, name, title, "TDirectory", objlen, f)
	dir.nbytesname = key.keylen
	dir.seekdir = key.seekkey

	buf := rbytes.NewWBuffer(make([]byte, objlen), nil, 0, f)
	buf.WriteString(f.id)
	buf.WriteString(f.Title())
	// dir-marshal
	_, err := dir.MarshalROOT(buf)
	if err != nil {
		panic(xerrors.Errorf("riofs: failed to write header: %w", err))
	}
	key.buf = buf.Bytes()
	key.obj = dir

	parent.keys = append(parent.keys, key)

	// key-write-file
	_, err = key.writeFile(f)
	if err != nil {
		panic(xerrors.Errorf("riofs: failed to write key header: %w", err))
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
	nbytes += int64(dir.dir.uuid.Sizeof())

	return nbytes
}

func (dir *tdirectoryFile) readDirInfo() error {
	f := dir.file
	nbytes := int64(f.nbytesname) + dir.recordSize(f.version)

	if nbytes+f.begin > f.end {
		return xerrors.Errorf(
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
		return xerrors.Errorf("riofs: can't read directory info")
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
		k.parent = dir
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

func (dir *tdirectoryFile) close() error {
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
	err := dir.saveSelf()
	if err != nil {
		return err
	}

	for _, sub := range dir.dirs {
		err = sub.save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (dir *tdirectoryFile) saveSelf() (err error) {
	err = dir.writeKeys()
	if err != nil {
		return err
	}

	err = dir.writeHeader()
	if err != nil {
		return err
	}

	return nil
}

// writeDirHeader overwrites the Directory header record.
func (dir *tdirectoryFile) writeHeader() error {
	var (
		err error
	)
	dir.mtime = nowUTC()

	nbytes := int32(dir.recordSize(dir.file.version))
	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	_, err = dir.MarshalROOT(buf)
	if err != nil {
		return xerrors.Errorf("riofs: could not marshal dir-info: %w", err)
	}

	_, err = dir.file.w.WriteAt(buf.Bytes(), dir.seekdir+int64(dir.nbytesname))
	if err != nil {
		return xerrors.Errorf("riofs: could not write dir-info to file: %w", err)
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
	var keys []*Key
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
			keys = append(keys, k)
		}
	}
	var key *Key
	switch len(keys) {
	case 0:
		return nil, noKeyError{key: namecycle, obj: dir}
	case 1:
		key = keys[0]
	default:
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].Cycle() < keys[j].Cycle()
		})
		key = keys[len(keys)-1]
	}

	obj, err := key.Object()
	if err != nil {
		return nil, err
	}

	if obj != nil {
		switch obj := obj.(type) {
		case *tdirectoryFile:
			obj.dir.parent = dir
			if obj.dir.Name() == "" {
				obj.dir.named.SetName(name)
			}
			if obj.Title() == "" {
				obj.dir.named.SetTitle(name)
			}
		}
	}
	return obj, nil
}

func (dir *tdirectoryFile) Put(name string, obj root.Object) error {
	if dir.file.w == nil {
		return xerrors.Errorf("could not put %q into directory %q: %w", name, dir.dir.Name(), ErrReadOnly)
	}

	if strings.Contains(name, "/") {
		return xerrors.Errorf("riofs: invalid path name %q (contains a '/')", name)
	}

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
		return xerrors.Errorf("riofs: empty key name")
	}

	// FIXME(sbinet): implement a fast look-up ?
	for i := range dir.keys {
		key := &dir.keys[i]
		if key.name != name {
			continue
		}
		if key.ClassName() != obj.Class() {
			return keyTypeError{key: name, class: key.ClassName()}
		}
		if key.cycle > cycle {
			cycle = key.cycle
		}
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
				return xerrors.Errorf("riofs: could not generate streamer for key %q and type %T: %w", name, obj, err)
			}
			si, err = dir.StreamerInfo(cxx, -1)
		}
		if err != nil {
			return xerrors.Errorf("riofs: could not find streamer for %T: %w", obj, err)
		}
		dir.addStreamer(si)
	}

	key, err := newKeyFrom(dir, name, title, rdict.GoName2Cxx(typename), obj, dir.file)
	if err != nil {
		return xerrors.Errorf("riofs: could not create key %q for object %T: %w", name, obj, err)
	}
	key.cycle = cycle
	_, err = key.writeFile(dir.file)
	if err != nil {
		return xerrors.Errorf("riofs: could not write key %q to file: %w", name, err)
	}

	dir.keys = append(dir.keys, key)

	return nil
}

// Keys returns the list of keys being held by this directory.
func (dir *tdirectoryFile) Keys() []Key {
	return dir.keys
}

// Mkdir creates a new subdirectory
func (dir *tdirectoryFile) Mkdir(name string) (Directory, error) {
	if _, err := dir.Get(name); err == nil {
		return nil, xerrors.Errorf("riofs: %q already exists", name)
	}

	if strings.Contains(name, "/") {
		return nil, xerrors.Errorf("riofs: invalid directory name %q (contains a '/')", name)
	}

	sub := newDirectoryFile(name, "", dir.file, dir)
	dir.dirs = append(dir.dirs, sub)

	return sub, nil
}

// Parent returns the directory holding this directory.
// Parent returns nil if this is the top-level directory.
func (dir *tdirectoryFile) Parent() Directory {
	if dir.dir.parent == nil {
		return dir.file
	}
	return dir.dir.parent
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

	version := dir.RVersion()
	if dir.isBigFile() && version < 1000 {
		version += 1000
	}
	w.WriteI16(version)
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

	switch version := version % 1000; {
	case version == 2:
		_ = dir.dir.uuid.UnmarshalROOTv1(r)
	case version > 2:
		_ = dir.dir.uuid.UnmarshalROOT(r)
	}

	return r.Err()
}

// StreamerInfo returns the StreamerInfo with name of this directory, or nil otherwise.
// If version is negative, the latest version should be returned.
func (dir *tdirectoryFile) StreamerInfo(name string, version int) (rbytes.StreamerInfo, error) {
	if dir.file == nil {
		return nil, xerrors.Errorf("riofs: no streamers")
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
		nbytes += key.keylen
	}

	hdr := newKey(dir, dir.Name(), dir.Title(), "TDirectory", nbytes, dir.file)

	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.WriteI32(int32(len(dir.Keys())))
	for _, k := range dir.Keys() {
		_, err = k.MarshalROOT(buf)
		if err != nil {
			return xerrors.Errorf("riofs: could not write key: %w", err)
		}
	}
	hdr.buf = buf.Bytes()

	dir.seekkeys = hdr.seekkey
	dir.nbyteskeys = hdr.nbytes

	_, err = hdr.writeFile(dir.file)
	if err != nil {
		return xerrors.Errorf("riofs: could not write header key: %w", err)
	}
	return nil
}

// writeDirHeader overwrites the Directory header record.
func (dir *tdirectoryFile) writeDirHeader() error {
	var (
		err error
	)
	dir.mtime = nowUTC()

	nbytes := int32(dir.recordSize(dir.file.version)) + int32(dir.file.nbytesname)
	key := newKey(dir, dir.Name(), dir.Title(), "TFile", nbytes, dir.file)
	key.seekkey = dir.seekdir
	key.seekpdir = dir.file.begin

	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, nil)
	buf.WriteString(dir.Name())
	buf.WriteString(dir.Title())
	_, err = dir.MarshalROOT(buf)
	if err != nil {
		return xerrors.Errorf("riofs: could not marshal dir-info: %w", err)
	}

	key.buf = buf.Bytes()
	_, err = key.writeFile(dir.file)
	if err != nil {
		return xerrors.Errorf("riofs: could not write dir-info to file: %w", err)
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

func (dir *tdirectoryFile) records(w io.Writer, indent int) error {
	hdr := strings.Repeat("  ", indent)
	fmt.Fprintf(w, "%s=== dir %q @%d ===\n", hdr, dir.Name(), dir.seekdir)
	parent := "<nil>"
	if dir.dir.parent != nil {
		parent = fmt.Sprintf("@%d", dir.dir.parent.(*tdirectoryFile).seekdir)
	}
	fmt.Fprintf(w, "%sparent:      %s\n", hdr, parent)
	fmt.Fprintf(w, "%snbytes-keys: %d\n", hdr, dir.nbyteskeys)
	fmt.Fprintf(w, "%snbytes-name: %d\n", hdr, dir.nbytesname)
	fmt.Fprintf(w, "%sseek-dir:    %d\n", hdr, dir.seekdir)
	fmt.Fprintf(w, "%sseek-parent: %d\n", hdr, dir.seekparent)
	fmt.Fprintf(w, "%sseek-keys:   %d\n", hdr, dir.seekkeys)
	fmt.Fprintf(w, "%sclass:       %q\n", hdr, dir.classname)
	fmt.Fprintf(w, "%skeys:        %d\n", hdr, len(dir.keys))
	for i := range dir.keys {
		k := &dir.keys[i]
		fmt.Fprintf(w, "%skey[%d]: %q\n", hdr+" ", i, k.Name())
		err := k.records(w, indent+1)
		if err != nil {
			return xerrors.Errorf("could not inspect key %q: %w", k.Name(), err)
		}
	}

	return nil
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
			o := newDirectoryFile("", "", nil, nil)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TDirectoryFile", f)
	}
}

// coreTypes is the set of types that do not need a streamer info.
var coreTypes = map[string]struct{}{
	"TObject":        {},
	"TFile":          {},
	"TDirectoryFile": {},
	"TKey":           {},
	"TString":        {},

	"TDatime":       {},
	"TVirtualIndex": {},
	"TBasket":       {},
}

func isCoreType(typename string) bool {
	_, ok := coreTypes[typename]
	return ok
}

func isCxxBuiltin(typename string) bool {
	_, ok := rmeta.CxxBuiltins[typename]
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
