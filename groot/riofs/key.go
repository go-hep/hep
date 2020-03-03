// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"io"
	"reflect"
	"strings"
	"time"

	"go-hep.org/x/hep/groot/internal/rcompress"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// noKeyError is the error returned when a riofs.Key could not be found.
type noKeyError struct {
	key string
	obj root.Named
}

func (err noKeyError) Error() string {
	return fmt.Sprintf("riofs: %s: could not find key %q", err.obj.Name(), err.key)
}

// keyTypeError is the error returned when a riofs.Key was found but the associated
// value is not of the expected type.
type keyTypeError struct {
	key   string
	class string
}

func (err keyTypeError) Error() string {
	return fmt.Sprintf("riofs: inconsistent value type for key %q (type=%s)", err.key, err.class)
}

// Key is a key (a label) in a ROOT file
//
//  The Key class includes functions to book space on a file,
//   to create I/O buffers, to fill these buffers
//   to compress/uncompress data buffers.
//
//  Before saving (making persistent) an object on a file, a key must
//  be created. The key structure contains all the information to
//  uniquely identify a persistent object on a file.
//  The Key class is used by ROOT:
//    - to write an object in the Current Directory
//    - to write a new ntuple buffer
type Key struct {
	f *File // underlying file

	rvers    int16     // version of the Key struct
	nbytes   int32     // number of bytes for the compressed object+key
	objlen   int32     // length of uncompressed object
	datetime time.Time // Date/Time when the object was written
	keylen   int32     // number of bytes for the Key struct
	cycle    int16     // cycle number of the object

	// address of the object on file (points to Key.bytes)
	// this is a redundant information used to cross-check
	// the data base integrity
	seekkey  int64
	seekpdir int64 // pointer to the directory supporting this object

	class string // object class name
	name  string // name of the object
	title string // title of the object

	left int32 // number of bytes left in current segment.

	buf []byte      // buffer of the Key's value
	obj root.Object // Key's value

	otyp reflect.Type // Go type of the Key's payload.

	parent Directory // directory holding this key
}

func newKey(dir *tdirectoryFile, name, title, class string, objlen int32, f *File) Key {
	k := Key{
		f:        f,
		rvers:    rvers.Key,
		objlen:   objlen,
		datetime: nowUTC(),
		cycle:    1,
		class:    class,
		name:     name,
		title:    title,
		seekpdir: f.begin, // FIXME(sbinet): see https://sft.its.cern.ch/jira/browse/ROOT-10352
		parent:   dir,
	}
	k.keylen = k.sizeof()
	// FIXME(sbinet): this assumes the key-payload isn't compressed.
	// if the key's payload is actually compressed, we introduce a hole
	// with the f.setEnd call below.
	k.nbytes = k.objlen + k.keylen
	if objlen > 0 {
		k.seekkey = f.end
		err := f.setEnd(k.seekkey + int64(k.nbytes))
		if err != nil {
			panic(err)
		}
	}

	if f.end > kStartBigFile {
		k.rvers += 1000
	}

	if dir != nil {
		k.seekpdir = dir.seekdir
	}

	return k
}

// NewKey creates a new key from the provided serialized object buffer.
// NewKey puts the key and its payload at the end of the provided file f.
// Depending on the file configuration, NewKey may compress the provided object buffer.
func NewKey(dir Directory, name, title, class string, cycle int16, obj []byte, f *File) (Key, error) {
	var d *tdirectoryFile
	if dir != nil {
		d = dir.(*tdirectoryFile)
	}
	return newKeyFromBuf(d, name, title, class, cycle, obj, f)
}

func newKeyFrom(dir *tdirectoryFile, name, title, class string, obj root.Object, f *File) (Key, error) {
	var err error
	if dir == nil {
		dir = &f.dir
	}

	keylen := keylenFor(name, title, class, dir)

	buf := rbytes.NewWBuffer(nil, nil, uint32(keylen), dir.file)
	switch obj := obj.(type) {
	case rbytes.Marshaler:
		_, err = obj.MarshalROOT(buf)
		if err != nil {
			return Key{}, fmt.Errorf("riofs: could not marshal object %T for key=%q: %w", obj, name, err)
		}
	default:
		return Key{}, fmt.Errorf("riofs: object %T can not be ROOT serialized", obj)
	}

	objlen := int32(len(buf.Bytes()))
	k := Key{
		f:        f,
		nbytes:   keylen + objlen,
		rvers:    rvers.Key,
		keylen:   keylen,
		objlen:   objlen,
		datetime: nowUTC(),
		cycle:    1,
		class:    class,
		name:     name,
		title:    title,
		seekkey:  f.end,
		seekpdir: dir.seekdir,
		obj:      obj,
		otyp:     reflect.TypeOf(obj),
		parent:   dir,
	}
	if f.end > kStartBigFile {
		k.rvers += 1000
	}

	k.buf, err = rcompress.Compress(nil, buf.Bytes(), k.f.compression)
	if err != nil {
		return k, fmt.Errorf("riofs: could not compress object %T for key %q: %w", obj, name, err)
	}
	k.nbytes = k.keylen + int32(len(k.buf))

	err = f.setEnd(k.seekkey + int64(k.nbytes))
	if err != nil {
		return k, fmt.Errorf("riofs: could not update ROOT file end: %w", err)
	}

	return k, nil
}

func newKeyFromBuf(dir *tdirectoryFile, name, title, class string, cycle int16, buf []byte, f *File) (Key, error) {
	var err error
	if dir == nil {
		dir = &f.dir
	}

	keylen := keylenFor(name, title, class, dir)
	objlen := int32(len(buf))
	k := Key{
		f:        f,
		nbytes:   keylen + objlen,
		rvers:    rvers.Key,
		keylen:   keylen,
		objlen:   objlen,
		datetime: nowUTC(),
		cycle:    cycle,
		class:    class,
		name:     name,
		title:    title,
		seekkey:  f.end,
		seekpdir: dir.seekdir,
		parent:   dir,
	}
	if f.end > kStartBigFile {
		k.rvers += 1000
	}

	k.buf, err = rcompress.Compress(nil, buf, k.f.compression)
	if err != nil {
		return k, fmt.Errorf("riofs: could not compress object %s for key %q: %w", class, name, err)
	}
	k.nbytes = k.keylen + int32(len(k.buf))

	err = f.setEnd(k.seekkey + int64(k.nbytes))
	if err != nil {
		return k, fmt.Errorf("riofs: could not update ROOT file end: %w", err)
	}

	return k, nil
}

// NewKeyForBasketInternal creates a new empty key.
// This is needed for Tree/Branch/Basket persistency.
//
// DO NOT USE.
func NewKeyForBasketInternal(dir Directory, name, title, class string, cycle int16) Key {
	var (
		f = fileOf(dir)
		d *tdirectoryFile
	)
	switch v := dir.(type) {
	case *File:
		d = &v.dir
	case *tdirectoryFile:
		d = v
	default:
		panic(fmt.Errorf("riofs: invalid directory type %T", dir))
	}

	k := Key{
		f:        f,
		rvers:    rvers.Key,
		cycle:    cycle,
		datetime: nowUTC(),
		class:    class,
		name:     name,
		title:    title,
		seekpdir: f.begin, // FIXME(sbinet): see https://sft.its.cern.ch/jira/browse/ROOT-10352
		parent:   dir,
	}
	k.keylen = k.sizeof()
	k.nbytes = k.keylen
	if f.end > kStartBigFile {
		k.rvers += 1000
	}

	if d != nil {
		k.seekpdir = d.seekdir
	}

	return k
}

// KeyFromDir creates a new empty key (with no associated payload object)
// with provided name and title, and the expected object type name.
// The key will be held by the provided directory.
func KeyFromDir(dir Directory, name, title, class string) Key {
	f := fileOf(dir)
	var k Key
	switch v := dir.(type) {
	case *File:
		k = newKey(&v.dir, name, title, class, 0, f)
	case *tdirectoryFile:
		k = newKey(v, name, title, class, 0, f)
	default:
		panic(fmt.Errorf("riofs: invalid directory type %T", dir))
	}
	return k
}

func (k *Key) RVersion() int16 { return k.rvers }

func (*Key) Class() string {
	return "TKey"
}

func (k *Key) ClassName() string {
	return k.class
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Title() string {
	return k.title
}

func (k *Key) Cycle() int {
	return int(k.cycle)
}

func (k *Key) Nbytes() int32  { return k.nbytes }
func (k *Key) ObjLen() int32  { return k.objlen }
func (k *Key) KeyLen() int32  { return k.keylen }
func (k *Key) SeekKey() int64 { return k.seekkey }
func (k *Key) Buffer() []byte { return k.buf }

func (k *Key) SetFile(f *File)      { k.f = f }
func (k *Key) SetBuffer(buf []byte) { k.buf = buf; k.objlen = int32(len(buf)) }

// ObjectType returns the Key's payload type.
//
// ObjectType returns nil if the Key's payload type is not known
// to the registry of groot.
func (k *Key) ObjectType() reflect.Type {
	if k.otyp != nil {
		return k.otyp
	}
	if !rtypes.Factory.HasKey(k.class) {
		return nil
	}
	k.otyp = rtypes.Factory.Get(k.class)().Type()
	return k.otyp
}

// Value returns the data corresponding to the Key's value
func (k *Key) Value() interface{} {
	v, err := k.Object()
	if err != nil {
		panic(fmt.Errorf("error loading payload for %q: %w", k.Name(), err))
	}
	return v
}

// Object returns the (ROOT) object corresponding to the Key's value.
func (k *Key) Object() (root.Object, error) {
	if k.obj != nil {
		return k.obj, nil
	}

	buf, err := k.Bytes()
	if err != nil {
		return nil, fmt.Errorf("riofs: could not load key payload: %w", err)
	}

	fct := rtypes.Factory.Get(k.class)
	if fct == nil {
		return nil, fmt.Errorf("riofs: no registered factory for class %q (key=%q)", k.class, k.Name())
	}

	v := fct()
	obj, ok := v.Interface().(root.Object)
	if !ok {
		return nil, fmt.Errorf("riofs: class %q does not implement root.Object (key=%q)", k.class, k.Name())
	}

	vv, ok := obj.(rbytes.Unmarshaler)
	if !ok {
		return nil, fmt.Errorf("riofs: class %q does not implement rbytes.Unmarshaler (key=%q)", k.class, k.Name())
	}

	err = vv.UnmarshalROOT(rbytes.NewRBuffer(buf, nil, uint32(k.keylen), k.f))
	if err != nil {
		return nil, fmt.Errorf("riofs: could not unmarshal key payload: %w", err)
	}

	if vv, ok := obj.(SetFiler); ok {
		vv.SetFile(k.f)
	}
	if dir, ok := obj.(*tdirectoryFile); ok {
		dir.file = k.f
		dir.dir.parent = k.parent
		dir.dir.named.SetName(k.Name())
		dir.dir.named.SetTitle(k.Name())
		dir.classname = k.class
		err = dir.readKeys()
		if err != nil {
			return nil, err
		}
	}

	k.obj = obj
	return obj, nil
}

// Bytes returns the buffer of bytes corresponding to the Key's value
func (k *Key) Bytes() ([]byte, error) {
	data, err := k.load(nil)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (k *Key) Load(buf []byte) ([]byte, error) {
	return k.load(buf)
}

func (k *Key) load(buf []byte) ([]byte, error) {
	if len(buf) < int(k.objlen) {
		buf = make([]byte, k.objlen)
	}
	if len(k.buf) > 0 {
		copy(buf, k.buf)
		return buf, nil
	}
	if k.isCompressed() {
		start := k.seekkey + int64(k.keylen)
		sr := io.NewSectionReader(k.f, start, int64(k.nbytes)-int64(k.keylen))
		err := rcompress.Decompress(buf, sr)
		if err != nil {
			return nil, fmt.Errorf("riofs: could not decompress key payload: %w", err)
		}
		return buf, nil
	}
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.nbytes))
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, fmt.Errorf("riofs: could not read key payload: %w", err)
	}
	return buf, nil
}

func (k *Key) store() error {
	if k.buf != nil {
		return nil
	}

	k.keylen = k.sizeof()

	buf := rbytes.NewWBuffer(make([]byte, k.objlen), nil, uint32(k.keylen), k.f)
	_, err := k.obj.(rbytes.Marshaler).MarshalROOT(buf)
	if err != nil {
		return err
	}
	k.objlen = int32(len(buf.Bytes()))
	k.buf, err = rcompress.Compress(nil, buf.Bytes(), k.f.compression)
	if err != nil {
		return err
	}
	nbytes := int32(len(k.buf))
	k.nbytes = k.keylen + nbytes

	if k.seekkey <= 0 {
		panic("impossible: seekkey <= 0")
	}

	return nil
}

func (k *Key) isCompressed() bool {
	return k.objlen != k.nbytes-k.keylen
}

func (k *Key) isBigFile() bool {
	return k.rvers > 1000
}

// sizeof returns the size in bytes of the key header structure.
func (k *Key) sizeof() int32 {
	return keylenFor(k.name, k.title, k.class, &k.f.dir)
}

func keylenFor(name, title, class string, dir *tdirectoryFile) int32 {
	nbytes := int32(22)
	if dir.isBigFile() {
		nbytes += 8
	}
	nbytes += datimeSizeof()
	nbytes += tstringSizeof(class)
	nbytes += tstringSizeof(name)
	nbytes += tstringSizeof(title)
	if class == "TBasket" {
		nbytes += 2 // version
		nbytes += 4 // bufsize
		nbytes += 4 // nevsize
		nbytes += 4 // nevbuf
		nbytes += 4 // last
		nbytes += 1 // flag
	}
	return nbytes
}

// MarshalROOT encodes the key to the provided buffer.
func (k *Key) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()

	w.WriteI32(k.nbytes)
	if k.nbytes < 0 {
		return int(w.Pos() - pos), nil
	}

	w.WriteI16(k.RVersion())
	w.WriteI32(k.objlen)
	w.WriteU32(time2datime(k.datetime))
	w.WriteI16(int16(k.keylen))
	w.WriteI16(k.cycle)
	switch {
	case k.isBigFile():
		w.WriteI64(k.seekkey)
		// FIXME(sbinet): handle PidOffsetShift and PidOffset that are stored in the 16 highest bits of seekpdir
		w.WriteI64(k.seekpdir)
	default:
		w.WriteI32(int32(k.seekkey))
		w.WriteI32(int32(k.seekpdir))
	}
	w.WriteString(k.class)
	w.WriteString(k.name)
	w.WriteString(k.title)

	return int(w.Pos() - pos), nil
}

// UnmarshalROOT decodes the content of data into the Key
func (k *Key) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	k.nbytes = r.ReadI32()
	if k.nbytes < 0 {
		k.class = "[GAP]"
		return nil
	}

	k.rvers = r.ReadI16()
	k.objlen = r.ReadI32()
	k.datetime = datime2time(r.ReadU32())
	k.keylen = int32(r.ReadI16())
	k.cycle = r.ReadI16()

	switch {
	case k.isBigFile():
		k.seekkey = r.ReadI64()
		// FIXME(sbinet): handle PidOffsetShift and PidOffset that are stored in the 16 highest bits of seekpdir
		k.seekpdir = r.ReadI64()
	default:
		k.seekkey = int64(r.ReadI32())
		k.seekpdir = int64(r.ReadI32())
	}

	k.class = r.ReadString()
	k.name = r.ReadString()
	k.title = r.ReadString()

	return r.Err()
}

// writeFile writes the key's payload to the file
func (k *Key) writeFile(f *File) (int, error) {
	if k.left > 0 {
		w := rbytes.NewWBuffer(nil, nil, 0, nil)
		w.WriteI32(int32(-k.left))
		k.buf = append(k.buf, w.Bytes()...)
	}

	buf := rbytes.NewWBuffer(make([]byte, k.nbytes), nil, 0, f)
	_, err := k.MarshalROOT(buf)
	if err != nil {
		return 0, err
	}

	n, err := f.w.WriteAt(buf.Bytes(), k.seekkey)
	if err != nil {
		return n, err
	}
	nn, err := f.w.WriteAt(k.buf, k.seekkey+int64(k.keylen))
	n += nn
	if err != nil {
		return n, err
	}

	k.buf = nil
	return n, nil
}

func (k *Key) records(w io.Writer, indent int) error {
	hdr := strings.Repeat("  ", indent)
	fmt.Fprintf(w, "%s=== key %q ===\n", hdr, k.Name())
	fmt.Fprintf(w, "%snbytes:    %d\n", hdr, k.nbytes)
	fmt.Fprintf(w, "%skeylen:    %d\n", hdr, k.keylen)
	fmt.Fprintf(w, "%sobjlen:    %d\n", hdr, k.objlen)
	fmt.Fprintf(w, "%scycle:     %d\n", hdr, k.cycle)
	fmt.Fprintf(w, "%sseek-key:  %d\n", hdr, k.seekkey)
	fmt.Fprintf(w, "%sseek-pdir: %d\n", hdr, k.seekpdir)
	fmt.Fprintf(w, "%sclass:     %q\n", hdr, k.class)
	parent := "<nil>"
	if k.parent != nil {
		parent = fmt.Sprintf("@%d", k.parent.(*tdirectoryFile).seekdir)
	}
	fmt.Fprintf(w, "%sparent:    %s\n", hdr, parent)

	switch k.class {
	case "TDirectory", "TDirectoryFile":
		obj, err := k.Object()
		if err != nil {
			return fmt.Errorf("could not load object of key %q: %w", k.Name(), err)
		}
		return obj.(*tdirectoryFile).records(w, indent+1)
	}
	return nil
}

func init() {
	f := func() reflect.Value {
		o := &Key{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TKey", f)
}

var (
	_ root.Object        = (*Key)(nil)
	_ root.Named         = (*Key)(nil)
	_ rbytes.Marshaler   = (*Key)(nil)
	_ rbytes.Unmarshaler = (*Key)(nil)
)
