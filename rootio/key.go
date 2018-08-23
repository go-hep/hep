// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"reflect"
	"time"

	"github.com/pkg/errors"
)

// noKeyError is the error returned when a rootio.Key could not be found.
type noKeyError struct {
	key string
	obj Named
}

func (err noKeyError) Error() string {
	return fmt.Sprintf("rootio: %s: could not find key %q", err.obj.Name(), err.key)
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

	bytes    int32     // number of bytes for the compressed object+key
	version  int16     // version of the Key struct
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

	buf []byte // buffer of the Key's value
	obj Object // Key's value
}

func newKey(name, title, class string, nbytes int32, f *File) Key {
	k := Key{
		f:        f,
		version:  4, // FIXME(sbinet): harmonize versions
		objlen:   nbytes,
		datetime: nowUTC(),
		cycle:    1,
		class:    class,
		name:     name,
		title:    title,
	}
	k.keylen = k.sizeof()
	k.bytes = nbytes + k.keylen

	return k
}

// createKey creates a new key of the specified size.
func createKey(name, title, class string, nbytes int32, f *File) Key {
	k := newKey(name, title, class, nbytes, f)
	if f.end > kStartBigFile {
		k.version += 1000
	}

	nsize := nbytes + k.keylen
	err := k.adjust(nsize)
	if err != nil {
		panic(err)
	}

	k.seekpdir = f.dir.dir.seekdir
	return k
}

func newKeyFrom(obj Object, wbuf *WBuffer) (Key, error) {
	if wbuf == nil {
		wbuf = NewWBuffer(nil, nil, 0, nil)
	}
	beg := int(wbuf.Pos())
	n, err := obj.(ROOTMarshaler).MarshalROOT(wbuf)
	if err != nil {
		return Key{}, err
	}
	end := beg + n
	data := wbuf.buffer()[beg:end]

	name := ""
	title := ""
	if obj, ok := obj.(Named); ok {
		name = obj.Name()
		title = obj.Title()
	}

	k := Key{
		version:  4, // FIXME(sbinet): harmonize versions
		objlen:   int32(n),
		datetime: nowUTC(),
		class:    obj.Class(),
		name:     name,
		title:    title,
		buf:      data,
		obj:      obj,
	}
	return k, nil
}

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

// Value returns the data corresponding to the Key's value
func (k *Key) Value() interface{} {
	v, err := k.Object()
	if err != nil {
		panic(fmt.Errorf("error loading payload for %q: %v", k.Name(), err))
	}
	return v
}

// Object returns the (ROOT) object corresponding to the Key's value.
func (k *Key) Object() (Object, error) {
	if k.obj != nil {
		return k.obj, nil
	}

	buf, err := k.Bytes()
	if err != nil {
		return nil, err
	}

	fct := Factory.Get(k.class)
	if fct == nil {
		return nil, fmt.Errorf("rootio: no registered factory for class %q (key=%q)", k.class, k.Name())
	}

	v := fct()
	obj, ok := v.Interface().(Object)
	if !ok {
		return nil, fmt.Errorf("rootio: class %q does not implement rootio.Object (key=%q)", k.class, k.Name())
	}

	vv, ok := obj.(ROOTUnmarshaler)
	if !ok {
		return nil, fmt.Errorf("rootio: class %q does not implement rootio.ROOTUnmarshaler (key=%q)", k.class, k.Name())
	}

	err = vv.UnmarshalROOT(NewRBuffer(buf, nil, uint32(k.keylen), k.f))
	if err != nil {
		return nil, err
	}

	if vv, ok := obj.(SetFiler); ok {
		vv.SetFile(k.f)
	}
	if dir, ok := obj.(*tdirectory); ok {
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
		sr := io.NewSectionReader(k.f, start, int64(k.bytes)-int64(k.keylen))
		err := decompress(sr, buf)
		if err != nil {
			return nil, err
		}
		return buf, nil
	}
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.bytes))
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

	buf := NewWBuffer(make([]byte, k.bytes), nil, 0, k.f)
	_, err := k.obj.(ROOTMarshaler).MarshalROOT(buf)
	if err != nil {
		return err
	}
	k.buf = buf.buffer()         // FIXME(sbinet): handle compression
	k.objlen = int32(len(k.buf)) // FIXME(sbinet): handle compression
	nbytes := k.objlen
	k.keylen = k.sizeof()
	k.bytes = k.keylen + nbytes

	if k.seekkey == 0 {
		// find a place on file where to store that key.
		err = k.adjust(k.bytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k *Key) adjust(nsize int32) error {
	best := k.f.spans.best(int64(nsize))
	if best == nil {
		return errors.Errorf("rootio: could not find a suitable free segment")
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

	k.bytes = nsize

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
	return k.objlen != k.bytes-k.keylen
}

func (k *Key) isBigFile() bool {
	return k.version > 1000
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
func (k *Key) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()

	w.WriteI32(k.bytes)
	if k.bytes < 0 {
		return int(int64(w.w.c) - pos), nil
	}

	w.WriteI16(k.version)
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

	return int(int64(w.w.c) - pos), nil
}

// UnmarshalROOT decodes the content of data into the Key
func (k *Key) UnmarshalROOT(r *RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	k.bytes = r.ReadI32()
	if k.bytes < 0 {
		k.class = "[GAP]"
		return nil
	}

	k.version = r.ReadI16()
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
		w := NewWBuffer(nil, nil, 0, nil)
		w.WriteI32(int32(-k.left))
		k.buf = append(k.buf, w.buffer()...)
	}

	buf := NewWBuffer(make([]byte, k.bytes), nil, 0, f)
	_, err := k.MarshalROOT(buf)
	if err != nil {
		return 0, err
	}

	n, err := f.w.WriteAt(buf.buffer(), k.seekkey)
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
	Factory.add("TKey", f)
	Factory.add("*rootio.Key", f)
}

var (
	_ Object          = (*Key)(nil)
	_ Named           = (*Key)(nil)
	_ ROOTMarshaler   = (*Key)(nil)
	_ ROOTUnmarshaler = (*Key)(nil)
)
