// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"compress/flate"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/root"
)

var (
	ErrReadOnly = errors.New("riofs: file read-only")
)

type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

type Writer interface {
	io.Writer
	io.WriterAt
	io.Seeker
	io.Closer
}

type syncer interface {
	// Sync commits the current contents of the file to stable storage.
	Sync() error
}

type stater interface {
	// Stat returns a FileInfo describing the file.
	Stat() (os.FileInfo, error)
}

// FileOption configures internal states of a ROOT file.
type FileOption func(f *File) error

// A ROOT file is a suite of consecutive data records (TKey's) with
// the following format (see also the TKey class). If the key is
// located past the 32 bit file limit (> 2 GB) then some fields will
// be 8 instead of 4 bytes:
//    1->4            Nbytes    = Length of compressed object (in bytes)
//    5->6            Version   = TKey version identifier
//    7->10           ObjLen    = Length of uncompressed object
//    11->14          Datime    = Date and time when object was written to file
//    15->16          KeyLen    = Length of the key structure (in bytes)
//    17->18          Cycle     = Cycle of key
//    19->22 [19->26] SeekKey   = Pointer to record itself (consistency check)
//    23->26 [27->34] SeekPdir  = Pointer to directory header
//    27->27 [35->35] lname     = Number of bytes in the class name
//    28->.. [36->..] ClassName = Object Class Name
//    ..->..          lname     = Number of bytes in the object name
//    ..->..          Name      = lName bytes with the name of the object
//    ..->..          lTitle    = Number of bytes in the object title
//    ..->..          Title     = Title of the object
//    ----->          DATA      = Data bytes associated to the object
//
// The first data record starts at byte fBEGIN (currently set to kBEGIN).
// Bytes 1->kBEGIN contain the file description, when fVersion >= 1000000
// it is a large file (> 2 GB) and the offsets will be 8 bytes long and
// fUnits will be set to 8:
//    1->4            "root"      = Root file identifier
//    5->8            fVersion    = File format version
//    9->12           fBEGIN      = Pointer to first data record
//    13->16 [13->20] fEND        = Pointer to first free word at the EOF
//    17->20 [21->28] fSeekFree   = Pointer to FREE data record
//    21->24 [29->32] fNbytesFree = Number of bytes in FREE data record
//    25->28 [33->36] nfree       = Number of free data records
//    29->32 [37->40] fNbytesName = Number of bytes in TNamed at creation time
//    33->33 [41->41] fUnits      = Number of bytes for file pointers
//    34->37 [42->45] fCompress   = Compression level and algorithm
//    38->41 [46->53] fSeekInfo   = Pointer to TStreamerInfo record
//    42->45 [54->57] fNbytesInfo = Number of bytes in TStreamerInfo record
//    46->63 [58->75] fUUID       = Universal Unique ID
type File struct {
	r      Reader
	w      Writer
	seeker io.Seeker
	closer io.Closer

	id string //non-root, identifies filename, etc.

	version int32
	begin   int64

	// Remainder of record is variable length, 4 or 8 bytes per pointer
	end         int64
	seekfree    int64 // first available record
	nbytesfree  int32 // total bytes available
	nfree       int32 // total free bytes
	nbytesname  int32 // number of bytes in TNamed at creation time
	units       byte
	compression int32
	seekinfo    int64 // pointer to TStreamerInfo
	nbytesinfo  int32 // sizeof(TStreamerInfo)
	uuid        rbase.UUID

	dir    tdirectoryFile // root directory of this file
	siKey  Key
	sinfos []rbytes.StreamerInfo
	simap  map[rbytes.StreamerInfo]struct{} // local set of streamers, when writing

	spans freeList // list of free spans on file
}

// Open opens the named ROOT file for reading. If successful, methods on the
// returned file can be used for reading; the associated file descriptor
// has mode os.O_RDONLY.
func Open(path string) (*File, error) {
	fd, err := openFile(path)
	if err != nil {
		return nil, errors.Errorf("riofs: unable to open %q (%q)", path, err.Error())
	}

	f := &File{
		r:      fd,
		seeker: fd,
		closer: fd,
		id:     path,
	}
	f.dir.file = f

	err = f.readHeader()
	if err != nil {
		return nil, errors.Errorf("riofs: failed to read header %q: %v", path, err)
	}

	return f, nil
}

// NewReader creates a new ROOT file reader.
func NewReader(r Reader) (*File, error) {
	f := &File{
		r:      r,
		seeker: r,
		closer: r,
	}
	f.dir.file = f

	err := f.readHeader()
	if err != nil {
		return nil, errors.Errorf("riofs: failed to read header: %v", err)
	}
	f.id = f.dir.Name()

	return f, nil
}

// Create creates the named ROOT file for writing.
func Create(name string, opts ...FileOption) (*File, error) {
	fd, err := os.Create(name)
	if err != nil {
		return nil, errors.Errorf("riofs: unable to create %q (%q)", name, err.Error())
	}

	f := &File{
		w:           fd,
		seeker:      fd,
		closer:      fd,
		id:          name,
		version:     root.Version,
		begin:       kBEGIN,
		end:         kBEGIN,
		units:       4,
		compression: 1,
		sinfos:      nil,
		simap:       make(map[rbytes.StreamerInfo]struct{}),
	}
	f.dir = *newDirectoryFile(name, "", f, nil)
	f.dir.dir.named.SetTitle("")
	f.spans.add(f.begin, kStartBigFile)

	f.setCompression(kZLIB, flate.BestCompression)

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		err := opt(f)
		if err != nil {
			return nil, errors.Wrapf(err, "riofs: could not apply option to ROOT file")
		}
	}

	// write directory info
	namelen := f.dir.dir.named.Sizeof()
	objlen := namelen + int32(f.dir.recordSize(f.version))
	key := newKey(&f.dir, f.dir.Name(), f.dir.Title(), "TFile", objlen, f)
	f.nbytesname = key.keylen + namelen
	f.dir.nbytesname = key.keylen + namelen
	f.dir.seekdir = key.seekkey
	f.seekfree = 0
	f.nbytesfree = 0

	err = f.writeHeader()
	if err != nil {
		_ = fd.Close()
		_ = os.RemoveAll(name)
		return nil, errors.Errorf("riofs: failed to write header %q: %v", name, err)
	}

	buf := rbytes.NewWBuffer(make([]byte, objlen), nil, 0, f)
	buf.WriteString(f.id)
	buf.WriteString(f.Title())

	_, err = f.dir.MarshalROOT(buf)
	if err != nil {
		return nil, errors.Wrapf(err, "riofs: failed to write header")
	}
	key.buf = buf.Bytes()

	_, err = key.writeFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "riofs: failed to write key header")
	}

	return f, nil
}

func (f *File) setEnd(pos int64) error {
	f.end = pos
	if f.spans.Len() == 0 {
		return errors.Errorf("riofs: empty free segment list")
	}
	blk := f.spans.last()
	if blk == nil {
		return errors.Errorf("riofs: last free segment is nil")
	}

	if blk.last != kStartBigFile {
		return errors.Errorf("riofs: last free segment is not the file ending")
	}

	blk.first = pos
	return nil
}

// Stat returns the os.FileInfo structure describing this file.
func (f *File) Stat() (os.FileInfo, error) {
	if f.r != nil {
		if st, ok := f.r.(stater); ok {
			return st.Stat()
		}
	}
	if f.w != nil {
		if st, ok := f.w.(stater); ok {
			return st.Stat()
		}
	}
	return nil, errors.Errorf("riofs: underlying file w/o os.FileInfo")
}

// Read implements io.Reader
func (f *File) Read(p []byte) (int, error) {
	return f.r.Read(p)
}

// ReadAt implements io.ReaderAt
func (f *File) ReadAt(p []byte, off int64) (int, error) {
	return f.r.ReadAt(p, off)
}

// Seek implements io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	return f.seeker.Seek(offset, whence)
}

// Version returns the ROOT version this file was created with.
func (f *File) Version() int {
	return int(f.version)
}

func (f *File) readHeader() error {

	buf := make([]byte, 64+12) // 64: small file + extra space for big file
	if _, err := f.ReadAt(buf, 0); err != nil {
		return err
	}

	r := rbytes.NewRBuffer(buf, nil, 0, nil)

	// Header

	var magic [4]byte
	if _, err := io.ReadFull(r, magic[:]); err != nil || string(magic[:]) != "root" {
		if err != nil {
			return errors.Errorf("riofs: failed to read ROOT file magic header: %v", err)
		}
		return errors.Errorf("riofs: %q is not a root file", f.id)
	}

	f.version = r.ReadI32()
	f.begin = int64(r.ReadI32())
	if f.version < 1000000 { // small file
		f.end = int64(r.ReadI32())
		f.seekfree = int64(r.ReadI32())
		f.nbytesfree = r.ReadI32()
		f.nfree = r.ReadI32()
		f.nbytesname = r.ReadI32()
		f.units = r.ReadU8()
		f.compression = r.ReadI32()
		f.seekinfo = int64(r.ReadI32())
		f.nbytesinfo = r.ReadI32()
	} else { // large files
		f.end = r.ReadI64()
		f.seekfree = r.ReadI64()
		f.nbytesfree = r.ReadI32()
		f.nfree = r.ReadI32()
		f.nbytesname = r.ReadI32()
		f.units = r.ReadU8()
		f.compression = r.ReadI32()
		f.seekinfo = r.ReadI64()
		f.nbytesinfo = r.ReadI32()
	}
	f.version %= 1000000

	if _, err := io.ReadFull(r, f.uuid[:]); err != nil || r.Err() != nil {
		if err != nil {
			return errors.Errorf("riofs: failed to read ROOT's UUID file: %v", err)
		}
		return r.Err()
	}

	var err error

	err = f.dir.readDirInfo()
	if err != nil {
		return errors.Errorf("riofs: failed to read ROOT directory infos: %v", err)
	}

	if f.seekfree > 0 {
		err = f.readFreeSegments()
		if err != nil {
			return errors.Wrapf(err, "riofs: failed to read ROOT file free segments")
		}
	}

	if f.seekinfo > 0 {
		err = f.readStreamerInfo()
		if err != nil {
			return errors.Errorf("riofs: failed to read ROOT streamer infos: %v", err)
		}
	}

	err = f.dir.readKeys()
	if err != nil {
		return errors.Errorf("riofs: failed to read ROOT file keys: %v", err)
	}

	return nil
}

func (f *File) writeHeader() error {
	var (
		err   error
		nfree = int32(len(f.spans))
	)

	buf := rbytes.NewWBuffer(make([]byte, f.begin), nil, 0, f)
	buf.Write([]byte("root"))

	version := f.version
	if version < 1000000 && (f.end > kStartBigFile ||
		f.seekfree > kStartBigFile ||
		f.seekinfo > kStartBigFile) {
		version += 1000000
		f.units = 8
	}
	buf.WriteI32(version)
	buf.WriteI32(int32(f.begin))
	switch {
	case version < 1000000:
		buf.WriteI32(int32(f.end))
		buf.WriteI32(int32(f.seekfree))
		buf.WriteI32(f.nbytesfree)
		buf.WriteI32(nfree)
		buf.WriteI32(f.nbytesname)
		buf.WriteU8(f.units)
		buf.WriteI32(f.compression)
		buf.WriteI32(int32(f.seekinfo))
		buf.WriteI32(f.nbytesinfo)
	default:
		buf.WriteI64(f.end)
		buf.WriteI64(f.seekfree)
		buf.WriteI32(f.nbytesfree)
		buf.WriteI32(nfree)
		buf.WriteI32(f.nbytesname)
		buf.WriteU8(f.units)
		buf.WriteI32(f.compression)
		buf.WriteI64(f.seekinfo)
		buf.WriteI32(f.nbytesinfo)
	}

	_, err = f.uuid.MarshalROOT(buf)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write UUID's file header")
	}

	_, _ = f.w.WriteAt(make([]byte, f.begin), 0)
	_, err = f.w.WriteAt(buf.Bytes(), 0)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write file header")
	}

	if w, ok := f.w.(syncer); ok {
		err = w.Sync()
	}

	return err
}

// Close closes the File, rendering it unusable for I/O.
// It returns an error, if any.
func (f *File) Close() error {
	if f.closer == nil {
		return nil
	}

	var err error

	err = f.dir.close()
	if err != nil {
		return err
	}

	if f.w != nil {
		err = f.writeStreamerInfo()
		if err != nil {
			return err
		}

		err = f.writeFreeSegments()
		if err != nil {
			return err
		}

		err = f.writeHeader()
		if err != nil {
			return err
		}
	}

	for i := range f.dir.keys {
		k := &f.dir.keys[i]
		k.f = nil
	}
	f.dir.keys = nil
	f.dir.file = nil

	err = f.closer.Close()
	f.closer = nil
	return err
}

// Keys returns the list of keys this File contains
func (f *File) Keys() []Key {
	return f.dir.Keys()
}

func (f *File) Name() string {
	return f.dir.Name()
}

func (f *File) Title() string {
	return f.dir.Title()
}

func (f *File) Class() string {
	return "TFile"
}

// readStreamerInfo reads the list of StreamerInfo from this file
func (f *File) readStreamerInfo() error {
	if f.seekinfo <= 0 || f.seekinfo >= f.end {
		return errors.Errorf("riofs: invalid pointer to StreamerInfo (pos=%v end=%v)", f.seekinfo, f.end)

	}
	buf := make([]byte, int(f.nbytesinfo))
	nbytes, err := f.ReadAt(buf, f.seekinfo)
	if err != nil {
		return err
	}
	if nbytes != int(f.nbytesinfo) {
		return errors.Errorf("riofs: requested [%v] bytes. read [%v] bytes from file", f.nbytesinfo, nbytes)
	}

	err = f.siKey.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, nil))
	f.siKey.f = f
	if err != nil {
		return err
	}

	objs := f.siKey.Value().(root.List)
	f.sinfos = make([]rbytes.StreamerInfo, 0, objs.Len())
	for i := 0; i < objs.Len(); i++ {
		obj, ok := objs.At(i).(rbytes.StreamerInfo)
		if !ok {
			continue
		}
		f.sinfos = append(f.sinfos, obj)
		rdict.StreamerInfos.Add(obj)
	}
	return nil
}

// writeStreamerInfo rites the list of StreamerInfos used in this file.
func (f *File) writeStreamerInfo() error {
	if f.w == nil {
		return nil
	}

	var (
		err    error
		sinfos = rcont.NewList("", nil)
		rules  = rcont.NewList("listOfRules", nil)
	)

	err = f.findDepStreamers()
	if err != nil {
		return errors.Wrap(err, "riofs: could not find dependent streamers")
	}

	for _, si := range f.sinfos {
		sinfos.Append(si)
	}

	if rules.Len() > 0 {
		sinfos.Append(rules)
	}

	if f.seekinfo != 0 {
		f.markFree(f.seekinfo, f.seekinfo+int64(f.nbytesinfo)-1)
	}

	key := newKey(&f.dir, "StreamerInfo", sinfos.Title(), sinfos.Class(), 0, f)
	offset := uint32(key.keylen)
	buf := rbytes.NewWBuffer(nil, nil, offset, f)
	_, err = sinfos.MarshalROOT(buf)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write StreamerInfo list")
	}

	key = newKey(&f.dir, "StreamerInfo", sinfos.Title(), sinfos.Class(), int32(len(buf.Bytes())), f)
	key.buf = buf.Bytes()
	f.seekinfo = key.seekkey
	f.nbytesinfo = key.nbytes

	_, err = key.writeFile(f)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write StreamerInfo list key")
	}

	return nil
}

// findDepStreamers finds all the needed streamers for proper persistency.
func (f *File) findDepStreamers() error {
	type depsType struct {
		name string
		vers int
	}

	var (
		deps []depsType
		err  error
	)

	for _, si := range f.sinfos {
		err = rdict.Visit(rdict.StreamerInfos, si, func(depth int, se rbytes.StreamerElement) error {
			switch se := se.(type) {
			case *rdict.StreamerBase:
				deps = append(deps, depsType{se.Name(), se.Base()})
			case *rdict.StreamerObject, *rdict.StreamerObjectAny:
				deps = append(deps, depsType{se.TypeName(), -1})
			case *rdict.StreamerObjectPointer, *rdict.StreamerObjectAnyPointer:
				deps = append(deps, depsType{strings.TrimRight(se.TypeName(), "*"), -1})
			case *rdict.StreamerString, *rdict.StreamerSTLstring:
				deps = append(deps, depsType{se.TypeName(), -1})

			case *rdict.StreamerSTL:
				deps = append(deps, depsType{se.ElemTypeName(), -1})
			}
			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "riofs: could not visit all dependent streamers for %#v", si)
		}
	}

	for _, dep := range deps {
		if isCoreType(dep.name) || isCxxBuiltin(dep.name) {
			continue
		}
		sub, err := rdict.StreamerInfos.StreamerInfo(dep.name, dep.vers)
		if err != nil {
			return errors.Wrapf(err, "riofs: could not find streamer for %q and version=%d", dep.name, dep.vers)
		}
		f.addStreamer(sub)
	}

	return nil
}

// markFree marks unused bytes on the file.
// it's the equivalent of slice[beg:end] = nil.
func (f *File) markFree(beg, end int64) {
	if len(f.spans) == 0 {
		return
	}

	span := f.spans.add(beg, end)
	if span == nil {
		return
	}
	nbytes := span.free()
	if nbytes > 2000000000 {
		nbytes = 2000000000
	}
	buf := rbytes.NewWBuffer(make([]byte, 4), nil, 0, f)
	buf.WriteI32(-int32(nbytes))
	if end == f.end-1 {
		f.end = span.first
	}
	_, err := f.w.WriteAt(buf.Bytes(), span.first)
	if err != nil {
		panic(err)
	}
}

func (f *File) readFreeSegments() error {
	var err error
	buf := make([]byte, f.nbytesfree)
	nbytes, err := f.ReadAt(buf, f.seekfree)
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return err
	}
	if nbytes != len(buf) {
		return errors.Errorf("riofs: requested [%v] bytes, read [%v] bytes from file", f.nbytesfree, nbytes)
	}

	var key = Key{f: f}
	err = key.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, nil))
	if err != nil {
		panic(err)
		return err
	}
	buf, err = key.Bytes()
	if err != nil {
		return errors.Wrapf(err, "riofs: could not read key payload")
	}
	rbuf := rbytes.NewRBuffer(buf, nil, 0, nil)
	for rbuf.Len() > 0 {
		var span freeSegment
		err = span.UnmarshalROOT(rbuf)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		f.spans = append(f.spans, span)
	}

	return err
}

func (f *File) writeFreeSegments() error {
	var err error

	if f.seekfree != 0 {
		f.markFree(f.seekfree, f.seekfree+int64(f.nbytesfree)-1)
	}

	key := func() *Key {
		var nbytes int32
		for _, span := range f.spans {
			nbytes += span.sizeof()
		}
		if nbytes == 0 {
			return nil
		}
		key := newKey(&f.dir, f.Name(), f.Title(), "TFile", nbytes, f)
		if key.seekkey == 0 {
			return nil
		}
		return &key
	}()

	if key == nil {
		return nil
	}

	isBigFile := f.end > kStartBigFile
	if !isBigFile && f.end > kStartBigFile {
		// the free block list is large enough to bring the file over the
		// 2Gb limit.
		// The references and offsets are now 64b, so we need to redo the
		// calculation since the list of free blocks will not fit in the
		// original size.
		panic("not implemented")
	}

	nbytes := key.objlen
	buf := rbytes.NewWBuffer(make([]byte, nbytes), nil, 0, f)
	for _, span := range f.spans {
		_, err := span.MarshalROOT(buf)
		if err != nil {
			return errors.Wrapf(err, "riofs: could not marshal free-block")
		}
	}
	if abytes := buf.Pos(); abytes != int64(nbytes) {
		switch {
		case abytes < int64(nbytes):
			// most likely one of the 'free' segments was used
			// to store this key.
			// we thus have one less free-block to store than planned.
			copy(buf.Bytes()[abytes:], make([]byte, int64(nbytes)-abytes))
		default:
			panic("riofs: free block list larger than expected")
		}
	}

	f.nbytesfree = key.nbytes
	f.seekfree = key.seekkey
	key.buf = buf.Bytes()
	_, err = key.writeFile(f)
	if err != nil {
		return errors.Wrapf(err, "riofs: could not write free-block list")
	}
	return nil
}

// StreamerInfos returns the list of StreamerInfos of this file.
func (f *File) StreamerInfos() []rbytes.StreamerInfo {
	return f.sinfos
}

// StreamerInfo returns the named StreamerInfo.
// If version is negative, the latest version should be returned.
func (f *File) StreamerInfo(name string, version int) (rbytes.StreamerInfo, error) {
	if len(f.sinfos) == 0 {
		return nil, errors.Errorf("riofs: no streamer for %q (no streamerinfo list)", name)
	}

	for _, si := range f.sinfos {
		if si.Name() == name {
			return si, nil
		}
		if _, ok := rdict.Typename(name, si.Title()); ok {
			return si, nil
		}
	}

	si, ok := rdict.StreamerInfos.Get(name, version)
	if ok {
		return si, nil
	}

	// no streamer for "name" in that file.
	// try whether "name" isn't actually std::vector<T> and a streamer
	// for T is in that file.
	o := reStdVector.FindStringSubmatch(name)
	if o != nil {
		si := stdvecSIFrom(name, o[1], f)
		if si != nil {
			f.sinfos = append(f.sinfos, si)
			rdict.StreamerInfos.Add(si)
			return si, nil
		}
	}

	return nil, errors.Errorf("riofs: no streamer for %q", name)
}

func (f *File) addStreamer(streamer rbytes.StreamerInfo) {
	if isCoreType(streamer.Name()) {
		return
	}

	if _, dup := f.simap[streamer]; dup {
		return
	}

	f.simap[streamer] = struct{}{}
	f.sinfos = append(f.sinfos, streamer)
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
func (f *File) Get(namecycle string) (root.Object, error) {
	return f.dir.Get(namecycle)
}

// Put puts the object v under the key with the given name.
func (f *File) Put(name string, v root.Object) error {
	if f.w == nil {
		return errors.Wrapf(ErrReadOnly, "could not put %q into file %q", name, f.Name())
	}
	return f.dir.Put(name, v)
}

// Mkdir creates a new subdirectory
func (f *File) Mkdir(name string) (Directory, error) {
	if f.w == nil {
		return nil, errors.Wrapf(ErrReadOnly, "could not mkdir %q in file %q", name, f.Name())
	}
	return f.dir.Mkdir(name)
}

// Parent returns the directory holding this directory.
// Parent returns nil if this is the top-level directory.
func (*File) Parent() Directory { return nil }

// SegmentMap displays to w the file's segments map.
func (f *File) SegmentMap(w io.Writer) (err error) {
	const timefmt = "20060102/150405"
	var (
		idcur = f.begin
		sz    = int64(64)
		date  time.Time
	)

	ndigits := int(math.Log10(float64(f.end))) + 1
	for idcur < f.end {
		var (
			buf = make([]byte, sz)
			n   int
		)
		n, err = f.ReadAt(buf, idcur)
		switch err {
		case nil:
			// ok
		case io.EOF:
			if n <= 0 {
				return errors.Wrapf(err, "could not buffer at position %d", idcur)
			}
		default:
			return errors.Wrapf(err, "could not buffer at position %d", idcur)
		}

		var k Key
		err = k.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, 0, f))
		if err != nil {
			return errors.Wrapf(err, "could not unmarshal key at %d", idcur)
		}
		date = k.datetime
		cname := k.ClassName()
		if k.nbytes < 0 {
			cname = "=== [GAP] ==="
			fmt.Fprintf(w, "%s  At:%-*d  N=%-8d  %-14s\n", strings.Repeat("*", 15), ndigits+1, idcur, k.nbytes, cname)
			idcur += int64(-k.nbytes)
			continue
		}
		if k.nbytes == 0 {
			return errors.Errorf("invalid key %q", k.name)
		}

		switch idcur {
		case f.seekfree:
			cname = "FreeSegments"
		case f.seekinfo:
			cname = "StreamerInfo"
		case f.dir.seekkeys:
			cname = "KeysList"
		}

		switch {
		case k.isCompressed():
			cx := float64(k.ObjLen()+k.KeyLen()) / float64(k.nbytes)
			fmt.Fprintf(w, "%s  At:%-*d  N=%-8d  %-14s CX = %5.2f\n", date.Format(timefmt), ndigits+1, idcur, k.nbytes, cname, cx)
		default:
			fmt.Fprintf(w, "%s  At:%-*d  N=%-8d  %-14s\n", date.Format(timefmt), ndigits+1, idcur, k.nbytes, cname)
		}
		idcur += int64(k.nbytes)
	}

	fmt.Fprintf(w, "%s  At:%-*d  N=%-8d  %-14s\n", date.Format(timefmt), ndigits+1, idcur, 1, "END")
	return err
}

// Records writes the records structure of the ROOT file to w.
func (f *File) Records(w io.Writer) error {
	fmt.Fprintf(w, "=== file %q ===\n", f.id)
	fmt.Fprintf(w, "begin: %d\n", f.begin)
	fmt.Fprintf(w, "end:   %d\n", f.end)
	fmt.Fprintf(w, "seek-free: %d nbytes-free=%d nfree=%d\n", f.seekfree, f.nbytesfree, f.nfree)
	fmt.Fprintf(w, "seek-info: %d nbytes-info=%d\n", f.seekinfo, f.nbytesinfo)

	return f.dir.records(w, 0)
}

var (
	_ root.Object                = (*File)(nil)
	_ root.Named                 = (*File)(nil)
	_ Directory                  = (*File)(nil)
	_ rbytes.StreamerInfoContext = (*File)(nil)
	_ streamerInfoStore          = (*File)(nil)
)
