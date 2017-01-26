// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"reflect"
	"time"
)

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

	classname string // object class name
	name      string // name of the object
	title     string // title of the object

	read bool
	data []byte
}

func (k *Key) Class() string {
	return k.classname
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Title() string {
	return k.title
}

// Bytes returns the buffer of bytes corresponding to the Key's value
func (k *Key) Bytes() ([]byte, error) {
	if !k.read {
		data, err := k.load()
		if err != nil {
			return nil, err
		}
		k.data = data
		k.read = true
	}
	return k.data, nil
}

// Note: this contains ZL[src][dst] where src and dst are 3 bytes each.
// Won't bother with this for the moment, since we can cross-check against
// objlen.
const rootHDRSIZE = 9

func (k *Key) load() ([]byte, error) {
	var buf bytes.Buffer
	if k.isCompressed() {
		start := k.seekkey + int64(k.keylen) + rootHDRSIZE
		r := io.NewSectionReader(k.f, start, int64(k.bytes)-int64(k.keylen))
		rc, err := zlib.NewReader(r)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.bytes))
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Value returns the data corresponding to the Key's value
func (k *Key) Value() interface{} {
	var v interface{}

	factory := Factory.Get(k.Class())
	if factory == nil {
		panic(fmt.Errorf("key[%v]: no factory for type [%s]\n", k.Name(), k.Class()))
	}

	vv := factory()
	if vv, ok := vv.Interface().(ROOTUnmarshaler); ok {
		data, err := k.Bytes()
		if err != nil {
			panic(fmt.Errorf("key[%v]: %v", k.Name(), err))
		}
		buf := bytes.NewBuffer(data)
		err = vv.UnmarshalROOT(buf)
		if err != nil {
			panic(err)
		}
		v = vv
	} else {
		panic(fmt.Errorf(
			"key[%v]: type [%s] does not satisfy the ROOTUnmarshaler interface",
			k.Name(), k.Class(),
		))
	}

	// FIXME: hack.
	if vv, ok := v.(*Tree); ok {
		vv.f = k.f
	}
	return v
}

/*
// Return Basket data associated with this key, if there is any
func (k *Key) AsBasket() *Basket {
	if k.classname != "TBasket" {
		panic("rootio.Key: Key is not a basket!")
	}
	b := &Basket{}
	_, err := k.f.Seek(int64(k.pdat), io.SeekStart)
	if err != nil {
		panic(fmt.Errorf("rootio.Key: %v", err))
	}

	dec := rootDecoder{r: k.f}
	err = dec.readBin(b)
	if err != nil {
		panic(fmt.Errorf("rootio.Key: %v", err))
	}
	return b
}
*/

func (k *Key) isCompressed() bool {
	return k.objlen != k.bytes-k.keylen
}

func (k *Key) Read() error {
	var err error
	f := k.f

	key_offset := f.Tell()
	myprintf("Key::Read -- @%v\n", key_offset)
	dec, err := newDecoderFromReader(f, 8)
	if err != nil {
		return err
	}

	dec.readInt32(&k.bytes)
	if dec.err != nil {
		return dec.err
	}

	myprintf("Key::Read -- @%v => %v\n", key_offset, k.bytes)

	if k.bytes < 0 {
		if k.classname == "[GAP]" {
			_, err = k.f.Seek(int64(-k.bytes)-4, io.SeekCurrent)
			if err != nil {
				return err
			}
			return k.Read()
		} else {
			return fmt.Errorf("rootio.Key: invalid bytes size [%v]", k.bytes)
		}
	}

	data := make([]byte, int(k.bytes))
	_, err = f.ReadAt(data, key_offset)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	err = k.UnmarshalROOT(buf)
	if err != nil {
		return err
	}

	if k.seekkey == 0 {
		k.seekkey = key_offset
	}

	if k.seekkey != key_offset {
		err = fmt.Errorf(
			"rootio.Key: Consistency failure: key offset %v doesn't match actual offset %v",
			k.seekkey, key_offset,
		)
		//fmt.Printf("*** %v\n", err)
		//panic(err)
	}

	_, err = k.f.Seek(int64(key_offset+int64(k.bytes)), io.SeekStart)
	return err
}

// UnmarshalROOT decodes the content of data into the Key
func (k *Key) UnmarshalROOT(data *bytes.Buffer) error {
	dec := newDecoder(data)
	myprintf("--- key ---\n")

	dec.readInt32(&k.bytes)
	myprintf("key-nbytes:  %v\n", k.bytes)

	if k.bytes < 0 {
		myprintf("Jumping gap: %v\n", k.bytes)
		k.classname = "[GAP]"
		return nil
	}

	dec.readInt16(&k.version)
	dec.readInt32(&k.objlen)
	var datetime uint32
	dec.readBin(&datetime)
	k.datetime = datime2time(datetime)
	dec.readInt16(&k.keylen)
	dec.readInt16(&k.cycle)

	if k.version > 1000 {
		dec.readInt64(&k.seekkey)
		dec.readInt64(&k.seekpdir)
	} else {
		dec.readInt32(&k.seekkey)
		dec.readInt32(&k.seekpdir)
	}

	dec.readString(&k.classname)
	dec.readString(&k.name)
	dec.readString(&k.title)

	myprintf("key-version: %v\n", k.version)
	myprintf("key-objlen:  %v\n", k.objlen)
	myprintf("key-cdate:   %v\n", k.datetime)
	myprintf("key-keylen:  %v\n", k.keylen)
	myprintf("key-cycle:   %v\n", k.cycle)
	myprintf("key-seekkey: %v\n", k.seekkey)
	myprintf("key-seekpdir:%v\n", k.seekpdir)
	myprintf("key-compress: %v %v %v %v %v\n", k.isCompressed(), k.objlen, k.bytes-k.keylen, k.bytes, k.keylen)
	myprintf("key-class: [%v]\n", k.classname)
	myprintf("key-name:  [%v]\n", k.name)
	myprintf("key-title: [%v]\n", k.title)

	//k.pdat = data

	return dec.err
}

func init() {
	f := func() reflect.Value {
		o := &Key{}
		return reflect.ValueOf(o)
	}
	Factory.add("TKey", f)
	Factory.add("*rootio.Key", f)
}

// testing interfaces
var _ Object = (*Key)(nil)
var _ ROOTUnmarshaler = (*Key)(nil)
