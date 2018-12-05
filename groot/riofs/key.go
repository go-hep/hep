// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/pkg/errors"
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

	nbytes   int32     // number of bytes for the compressed object+key
	rvers    int16     // version of the Key struct
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
}

func newKey(name, title, class string, nbytes int32, f *File) Key {
	k := Key{
		f:        f,
		rvers:    rvers.Key,
		objlen:   nbytes,
		datetime: nowUTC(),
		cycle:    1,
		class:    class,
		name:     name,
		title:    title,
	}
	k.keylen = k.sizeof()
	k.nbytes = nbytes + k.keylen

	return k
}

// createKey creates a new key of the specified size.
func createKey(name, title, class string, nbytes int32, f *File) Key {
	k := newKey(name, title, class, nbytes, f)
	if f.end > kStartBigFile {
		k.rvers += 1000
	}

	nsize := nbytes + k.keylen
	err := k.adjust(nsize)
	if err != nil {
		panic(err)
	}

	k.seekpdir = f.dir.seekdir
	return k
}

func newKeyFrom(obj root.Object, wbuf *rbytes.WBuffer) (Key, error) {
	if wbuf == nil {
		wbuf = rbytes.NewWBuffer(nil, nil, 0, nil)
	}
	beg := int(wbuf.Pos())
	n, err := obj.(rbytes.Marshaler).MarshalROOT(wbuf)
	if err != nil {
		return Key{}, err
	}
	end := beg + n
	data := wbuf.Bytes()[beg:end]

	name := ""
	title := ""
	if obj, ok := obj.(root.Named); ok {
		name = obj.Name()
		title = obj.Title()
	}

	k := Key{
		rvers:    rvers.Key,
		objlen:   int32(n),
		datetime: nowUTC(),
		class:    obj.Class(),
		name:     name,
		title:    title,
		buf:      data,
		obj:      obj,
		otyp:     reflect.TypeOf(obj),
	}
	return k, nil
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

func (k *Key) ObjLen() int32 { return k.objlen }
func (k *Key) KeyLen() int32 { return k.keylen }

func (k *Key) SetFile(f *File)      { k.f = f }
func (k *Key) SetBuffer(buf []byte) { k.buf = buf }

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
		panic(fmt.Errorf("error loading payload for %q: %+v", k.Name(), err))
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
		return nil, errors.Wrapf(err, "riofs: could not load key payload")
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
		return nil, errors.Wrapf(err, "riofs: could not unmarshal key payload")
	}

	if vv, ok := obj.(SetFiler); ok {
		vv.SetFile(k.f)
	}
	if dir, ok := obj.(*tdirectoryFile); ok {
		dir.file = k.f
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
		err := decompress(sr, buf)
		if err != nil {
			return nil, errors.Wrapf(err, "riofs: could not decompress key payload")
		}
		return buf, nil
	}
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.nbytes))
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (k *Key) store() error {
	if k.buf != nil {
		return nil
	}

	k.keylen = k.sizeof()

	buf := rbytes.NewWBuffer(make([]byte, k.nbytes), nil, uint32(k.keylen), k.f)
	_, err := k.obj.(rbytes.Marshaler).MarshalROOT(buf)
	if err != nil {
		return err
	}
	k.objlen = int32(len(buf.Bytes()))
	k.buf, err = compress(k.f.compression, buf.Bytes())
	if err != nil {
		return err
	}
	nbytes := int32(len(k.buf))
	k.nbytes = k.keylen + nbytes

	if k.seekkey <= 0 {
		// find a place on file where to store that key.
		err = k.adjust(k.nbytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *Key) adjust(nsize int32) error {
	best := k.f.spans.best(int64(nsize))
	if best == nil {
		return errors.Errorf("riofs: could not find a suitable free segment")
	}

	k.seekkey = best.first

	switch {
	case k.seekkey >= k.f.end:
		// segment at the end of the file.
		k.f.end = k.seekkey + int64(nsize)
		best.first = k.f.end
		if k.f.end > best.last {
			best.last += 1000000000
		}
		k.left = -1
	default:
		k.left = int32(best.last - k.seekkey - int64(nsize) + 1)
	}

	k.nbytes = nsize

	switch {
	case k.left == 0:
		// key's payload fills exactly a deleted gap.
		panic("not implemented -- k.left==0")

	case k.left > 0:
		// key's payload placed in a deleted gap larger than strictly needed.
		panic("not implemented -- k.left >0")
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
	nbytes := int32(22)
	if k.isBigFile() {
		nbytes += 8
	}
	nbytes += datimeSizeof()
	nbytes += tstringSizeof(k.class)
	nbytes += tstringSizeof(k.name)
	nbytes += tstringSizeof(k.title)
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
