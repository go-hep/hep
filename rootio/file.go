// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"

	"github.com/pkg/errors"
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
	uuid        [18]byte

	dir    tdirectoryFile // root directory of this file
	siKey  Key
	sinfos []StreamerInfo

	spans freeList // list of free spans on file
}

// Open opens the named ROOT file for reading. If successful, methods on the
// returned file can be used for reading; the associated file descriptor
// has mode os.O_RDONLY.
func Open(path string) (*File, error) {
	fd, err := openFile(path)
	if err != nil {
		return nil, fmt.Errorf("rootio: unable to open %q (%q)", path, err.Error())
	}

	f := &File{
		r:      fd,
		seeker: fd,
		closer: fd,
		id:     path,
	}
	f.dir.dir.file = f

	err = f.readHeader()
	if err != nil {
		return nil, fmt.Errorf("rootio: failed to read header %q: %v", path, err)
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
	f.dir.dir.file = f

	err := f.readHeader()
	if err != nil {
		return nil, fmt.Errorf("rootio: failed to read header: %v", err)
	}
	f.id = f.dir.dir.Name()

	return f, nil
}

// Create creates the named ROOT file for writing.
func Create(name string) (*File, error) {
	fd, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("rootio: unable to create %q (%q)", name, err.Error())
	}

	f := &File{
		w:           fd,
		seeker:      fd,
		closer:      fd,
		id:          name,
		version:     rootVersion,
		begin:       kBEGIN,
		end:         kBEGIN,
		units:       4,
		compression: 1,
		sinfos:      []StreamerInfo{},
	}
	f.dir = *newDirectoryFile(name, f)
	f.spans.add(kBEGIN, kStartBigFile)

	// write directory info
	namelen := f.dir.dir.named.sizeof()
	nbytes := namelen + f.dir.sizeof()
	key := createKey(f.dir.Name(), f.dir.Title(), "TFile", nbytes, f)
	f.nbytesname = key.keylen + namelen
	f.dir.dir.nbytesname = key.keylen + namelen
	f.dir.dir.seekdir = key.seekkey
	f.seekfree = 0
	f.nbytesfree = 0

	err = f.writeHeader()
	if err != nil {
		_ = fd.Close()
		_ = os.RemoveAll(name)
		return nil, fmt.Errorf("rootio: failed to write header %q: %v", name, err)
	}
	buf := NewWBuffer(make([]byte, key.bytes), nil, 0, f)
	_, err = f.dir.MarshalROOT(buf)
	if err != nil {
		return nil, errors.Wrapf(err, "rootio: failed to write header")
	}
	key.bytes = int32(len(buf.buffer()))
	key.buf = buf.buffer()

	_, err = key.writeFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "rootio: failed to write key header")
	}

	return f, nil
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
	return nil, errors.Errorf("rootio: underlying file w/o os.FileInfo")
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

	buf := make([]byte, 64)
	if _, err := f.ReadAt(buf, 0); err != nil {
		return err
	}

	r := NewRBuffer(buf, nil, 0, nil)

	// Header

	var magic [4]byte
	if _, err := io.ReadFull(r.r, magic[:]); err != nil || string(magic[:]) != "root" {
		if err != nil {
			return fmt.Errorf("rootio: failed to read ROOT file magic header: %v", err)
		}
		return fmt.Errorf("rootio: %q is not a root file", f.id)
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

	if _, err := io.ReadFull(r.r, f.uuid[:]); err != nil || r.Err() != nil {
		if err != nil {
			return fmt.Errorf("rootio: failed to read ROOT's UUID file: %v", err)
		}
		return r.Err()
	}

	var err error

	err = f.dir.readDirInfo()
	if err != nil {
		return fmt.Errorf("rootio: failed to read ROOT directory infos: %v", err)
	}

	if f.seekfree > 0 {
		err = f.readFreeSegments()
		if err != nil {
			return errors.Wrapf(err, "rootio: failed to read ROOT file free segments")
		}
	}

	if f.seekinfo > 0 {
		err = f.readStreamerInfo()
		if err != nil {
			return fmt.Errorf("rootio: failed to read ROOT streamer infos: %v", err)
		}
	}

	err = f.dir.readKeys()
	if err != nil {
		return fmt.Errorf("rootio: failed to read ROOT file keys: %v", err)
	}

	return nil
}

func (f *File) writeHeader() error {
	var (
		err   error
		span  = f.spans.last()
		nfree = int32(len(f.spans))
	)

	if span != nil {
		f.end = span.first
	}

	buf := NewWBuffer(make([]byte, f.begin), nil, 0, f)
	buf.write([]byte("root"))

	version := f.version
	if version < 1000000 && f.end > kStartBigFile {
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
	buf.write(f.uuid[:])

	_, _ = f.w.WriteAt(make([]byte, f.begin), 0)
	_, err = f.w.WriteAt(buf.buffer(), 0)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not write file header")
	}

	if w, ok := f.w.(syncer); ok {
		err = w.Sync()
	}

	return err
}

func (f *File) Tell() int64 {
	where, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		panic(err)
	}
	return where
}

// Close closes the File, rendering it unusable for I/O.
// It returns an error, if any.
func (f *File) Close() error {
	if f.closer == nil {
		return nil
	}

	var err error

	if f.w != nil {
		err = f.writeStreamerInfo()
		if err != nil {
			return err
		}
	}

	err = f.dir.Close()
	if err != nil {
		return err
	}

	if f.w != nil {
		err = f.writeFreeSegments()
		if err != nil {
			return err
		}

		err = f.writeHeader()
		if err != nil {
			return err
		}
	}

	for _, k := range f.dir.dir.keys {
		k.f = nil
	}
	f.dir.dir.keys = nil
	f.dir.dir.file = nil

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
		return fmt.Errorf("rootio: invalid pointer to StreamerInfo (pos=%v end=%v)", f.seekinfo, f.end)

	}
	buf := make([]byte, int(f.nbytesinfo))
	nbytes, err := f.ReadAt(buf, f.seekinfo)
	if err != nil {
		return err
	}
	if nbytes != int(f.nbytesinfo) {
		return fmt.Errorf("rootio: requested [%v] bytes. read [%v] bytes from file", f.nbytesinfo, nbytes)
	}

	err = f.siKey.UnmarshalROOT(NewRBuffer(buf, nil, 0, nil))
	f.siKey.f = f
	if err != nil {
		return err
	}

	objs := f.siKey.Value().(List)
	f.sinfos = make([]StreamerInfo, 0, objs.Len())
	for i := 0; i < objs.Len(); i++ {
		obj, ok := objs.At(i).(StreamerInfo)
		if !ok {
			continue
		}
		f.sinfos = append(f.sinfos, obj)
		streamers.add(obj)
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
		sinfos = newList("")
		rules  = newList("listOfRules")
	)

	for _, si := range f.sinfos {
		sinfos.objs = append(sinfos.objs, si)
	}

	if len(rules.objs) > 0 {
		sinfos.objs = append(sinfos.objs, rules)
	}

	if f.seekinfo != 0 {
		f.markFree(f.seekinfo, f.seekinfo+int64(f.nbytesinfo)-1)
	}

	buf := NewWBuffer(nil, nil, 0, f)
	_, err = sinfos.MarshalROOT(buf)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not write StreamerInfo list")
	}

	key := createKey("StreamerInfo", sinfos.Title(), sinfos.Class(), int32(len(buf.buffer())), f)
	key.buf = buf.buffer()
	f.seekinfo = key.seekkey
	f.nbytesinfo = key.bytes

	_, err = key.writeFile(f)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not write StreamerInfo list key")
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
	buf := NewWBuffer(make([]byte, 4), nil, 0, f)
	buf.writeI32(-int32(nbytes))
	if end == f.end-1 {
		f.end = span.first
	}
	_, err := f.w.WriteAt(buf.buffer(), span.first)
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
		return errors.Errorf("rootio: requested [%v] bytes, read [%v] bytes from file", f.nbytesfree, nbytes)
	}

	var key = Key{f: f}
	err = key.UnmarshalROOT(NewRBuffer(buf, nil, 0, nil))
	if err != nil {
		panic(err)
		return err
	}
	buf, err = key.Bytes()
	if err != nil {
		return errors.Wrapf(err, "rootio: could not read key payload")
	}
	rbuf := NewRBuffer(buf, nil, 0, nil)
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
		key := createKey(f.Name(), f.Title(), "TFile", nbytes, f)
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
	buf := NewWBuffer(make([]byte, nbytes), nil, 0, f)
	for _, span := range f.spans {
		_, err := span.MarshalROOT(buf)
		if err != nil {
			return errors.Wrapf(err, "rootio: could not marshal free-block")
		}
	}
	if abytes := buf.Pos(); abytes != int64(nbytes) {
		switch {
		case abytes < int64(nbytes):
			// most likely one of the 'free' segments was used
			// to store this key.
			// we thus have one less free-block to store than planned.
			copy(buf.buffer()[abytes:], make([]byte, int64(nbytes)-abytes))
		default:
			panic("rootio: free block list larger than expected")
		}
	}

	f.nbytesfree = key.bytes
	f.seekfree = key.seekkey
	key.buf = buf.buffer()
	_, err = key.writeFile(f)
	if err != nil {
		return errors.Wrapf(err, "rootio: could not write free-block list")
	}
	return nil
}

// StreamerInfos returns the list of StreamerInfos of this file.
func (f *File) StreamerInfos() []StreamerInfo {
	return f.sinfos
}

// StreamerInfo returns the StreamerInfo with name of this file and an error if any.
func (f *File) StreamerInfo(name string) (StreamerInfo, error) {
	if len(f.sinfos) == 0 {
		return nil, fmt.Errorf("rootio: no streamer for %q (no streamerinfo list)", name)
	}

	for _, si := range f.sinfos {
		if si.Name() == name {
			return si, nil
		}
	}

	// no streamer for "name" in that file.
	// try whether "name" isn't actually std::vector<T> and a streamer
	// for T is in that file.
	o := reStdVector.FindStringSubmatch(name)
	if o != nil {
		si := stdvecSIFrom(name, o[1], f)
		if si != nil {
			f.sinfos = append(f.sinfos, si)
			streamers.add(si)
			return si, nil
		}
	}

	return nil, fmt.Errorf("rootio: no streamer for %q", name)
}

func (f *File) addStreamer(streamer StreamerInfo) {
	var (
		idx  = -1
		name = streamer.Name()
	)

	for i, si := range f.sinfos {
		if si.Name() == name {
			idx = i
			break
		}
	}
	if idx != -1 {
		return
	}

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
func (f *File) Get(namecycle string) (Object, error) {
	return f.dir.Get(namecycle)
}

// Put puts the object v under the key with the given name.
func (f *File) Put(name string, v Object) error {
	return f.dir.Put(name, v)
}

var (
	_ Object              = (*File)(nil)
	_ Named               = (*File)(nil)
	_ Directory           = (*File)(nil)
	_ StreamerInfoContext = (*File)(nil)
	_ streamerInfoStore   = (*File)(nil)
)

type freeSegment struct {
	first int64 // first free word of segment
	last  int64 // last free word of segment
}

func (freeSegment) Class() string {
	return "TFree"
}

func (seg freeSegment) free() int64 {
	return seg.last - seg.first + 1
}

func (seg freeSegment) sizeof() int32 {
	if seg.last > kStartBigFile {
		return 18
	}
	return 10
}

func (seg freeSegment) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()

	w.w.grow(int(seg.sizeof()))

	vers := int16(1)
	if seg.last > kStartBigFile {
		vers += 1000
	}
	w.writeI16(vers)
	switch {
	case vers > 1000:
		w.writeI64(seg.first)
		w.writeI64(seg.last)
	default:
		w.writeI32(int32(seg.first))
		w.writeI32(int32(seg.last))
	}

	end := w.Pos()
	return int(end - pos), w.err
}

func (seg *freeSegment) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	vers := r.ReadI16()
	switch {
	case vers > 1000:
		seg.first = r.ReadI64()
		seg.last = r.ReadI64()
	default:
		seg.first = int64(r.ReadI32())
		seg.last = int64(r.ReadI32())
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &freeSegment{}
		return reflect.ValueOf(o)
	}
	Factory.add("TFree", f)
	Factory.add("*rootio.freeSegment", f)
}

var (
	_ Object          = (*freeSegment)(nil)
	_ ROOTMarshaler   = (*freeSegment)(nil)
	_ ROOTUnmarshaler = (*freeSegment)(nil)
)

// freeList describes the list of free segments on a ROOT file.
//
// Each ROOT file has a linked list of free segments.
// Each free segment is described by its first and last addresses.
// When an object is written to a file, a new Key is created. The first free
// segment big enough to accomodate the object is used.
//
// If the object size has a length corresponding to the size of the free segment,
// the free segment is deleted from the list of free segments.
// When an object is deleted from a file, a new freeList object is generated.
// If the deleted object is contiguous to an already deleted object, the free
// segments are merged in one single segment.
type freeList []freeSegment

func (p freeList) Len() int      { return len(p) }
func (p freeList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p freeList) Less(i, j int) bool {
	pi := p[i]
	pj := p[j]
	if pi.first < pj.first {
		return true
	}
	if pi.first == pj.first {
		return pi.last < pj.last
	}
	return false
}

func (fl *freeList) add(first, last int64) *freeSegment {
	elmt := freeSegment{first, last}
	*fl = append(*fl, elmt)
	sort.Sort(*fl)
	fl.consolidate()
	return fl.find(elmt)
}

func (fl *freeList) find(elmt freeSegment) *freeSegment {
	// FIXME(sbinet): use sort.Search
	for i := range *fl {
		cur := &(*fl)[i]
		if elmt.last < cur.first || cur.last < elmt.first {
			continue
		}
		if cur.first <= elmt.first && elmt.first <= cur.last &&
			elmt.last <= cur.last {
			return cur
		}
	}
	return nil
}

func (fl *freeList) consolidate() {
	for i := len(*fl) - 1; i >= 1; i-- {
		cur := &(*fl)[i]
		prev := &(*fl)[i-1]
		if prev.last+1 < cur.first {
			continue
		}
		if cur.last >= prev.last {
			prev.last = cur.last
		}
		fl.remove(i)
	}
}

func (fl *freeList) remove(i int) {
	list := *fl
	*fl = append(list[:i], list[i+1:]...)
}

// best returns the best free segment where to store nbytes.
func (fl freeList) best(nbytes int64) *freeSegment {
	var best *freeSegment

	if len(fl) == 0 {
		return best
	}

	for i, cur := range fl {
		nleft := cur.free()
		if nleft == nbytes {
			// exact match.
			return &fl[i]
		}
		if nleft >= nbytes+4 && best == nil {
			best = &fl[i]
		}
	}

	if best != nil {
		return best
	}

	// try big file
	best = &fl[len(fl)-1]
	best.last += 1000000000
	return best
}

func (fl freeList) last() *freeSegment {
	if len(fl) == 0 {
		return nil
	}
	return &fl[len(fl)-1]
}
