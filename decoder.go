// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package rootio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
)

type decoder struct {
	buf *bytes.Buffer
	len int64
	err error

	refs map[int64]interface{}
}

func newDecoder(buf *bytes.Buffer) *decoder {
	dec := &decoder{
		buf:  buf,
		len:  int64(buf.Len()),
		refs: make(map[int64]interface{}),
	}
	return dec
}

func newDecoderFromBytes(data []byte) *decoder {
	buf := make([]byte, len(data))
	copy(buf, data)
	return newDecoder(bytes.NewBuffer(buf))
}

func newDecoderFromReader(r io.Reader, size int) (*decoder, error) {
	data := make([]byte, size)
	n, err := r.Read(data)
	if err != nil {
		return nil, err
	}
	if n != size {
		return nil, fmt.Errorf("rootio.Decoder: read too few bytes [%v]. requested [%v]", n, size)
	}
	return newDecoder(bytes.NewBuffer(data)), nil
}

func (dec *decoder) Clone() *decoder {
	o := newDecoderFromBytes(dec.buf.Bytes())
	o.len = dec.len
	o.err = dec.err
	return o
}

func (dec *decoder) Pos() int64 {
	return dec.len - int64(dec.buf.Len())
}

func (dec *decoder) Len() int64 {
	return int64(dec.buf.Len())
}

func (dec *decoder) readString(s *string) {
	if dec.err != nil {
		return
	}

	var length byte
	var buf [256]byte

	dec.readBin(&length)
	if length != 0 {
		dec.readBin(buf[:length])
		if dec.err != nil {
			return
		}
		*s = string(buf[:length])
	}
}

func (dec *decoder) readCString(n int, s *string) {
	if dec.err != nil {
		return
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		var b byte
		dec.readBin(&b)
		if b == 0 {
			buf = buf[:i]
			break
		}
		buf[i] = b
	}
	*s = string(buf)
}

func (dec *decoder) readBin(v interface{}) {
	if dec.err != nil {
		return
	}
	dec.err = binary.Read(dec.buf, binary.BigEndian, v)
}

func (dec *decoder) readInt16(v interface{}) {
	if dec.err != nil {
		return
	}
	var d int16
	dec.readBin(&d)
	if dec.err != nil {
		return
	}

	switch uv := v.(type) {
	case *int16:
		*uv = int16(d)
	case *int32:
		*uv = int32(d)
	case *int64:
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
}

func (dec *decoder) readInt32(v interface{}) {
	if dec.err != nil {
		return
	}
	switch uv := v.(type) {
	case *int32:
		dec.readBin(v)
	case *int64:
		var d int32
		dec.readBin(&d)
		if dec.err != nil {
			return
		}
		*uv = int64(d)
	default:
		panic(fmt.Errorf("rootio: unknown type %T", v))
	}
}

func (dec *decoder) readInt64(v interface{}) {
	if dec.err != nil {
		return
	}
	switch uv := v.(type) {
	case *int64:
		var d int64
		dec.readBin(&d)
		if dec.err != nil {
			return
		}
		*uv = int64(d)
	default:
		panic(fmt.Errorf("rootio: unknown type %T", v))
	}
}

func (dec *decoder) readVersion() (version int16, position, bytecount int32) {
	if dec.err != nil {
		return
	}

	start := dec.Pos()

	var bcnt uint32
	dec.readBin(&bcnt)
	if dec.err != nil {
		return
	}
	myprintf("readVersion - bytecount=%v\n", bcnt)
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		bytecount = int32(int64(bcnt) & ^kByteCountMask)
	} else {
		dec.err = fmt.Errorf("rootio.readVersion: too old file")
		return
	}

	var vers uint16
	dec.readBin(&vers)
	if dec.err != nil {
		return
	}
	version = int16(vers)

	position = int32(start)
	myprintf("readVersion => [%v] [%v] [%v]\n", position, version, bytecount)
	return version, position, bytecount
}

func (dec *decoder) readClass(name *string, count *int32, isref *bool) {
	if dec.err != nil {
		return
	}

	var tag uint32
	dec.readBin(&tag)
	if dec.err != nil {
		return
	}
	myprintf("::readClass. first int: [%v]\n", tag)
	switch {
	case tag == kNullTag:
		*isref = false
		return

	case (tag & kByteCountMask) != 0:
		// bufvers = 1
		classtag := ""
		dec.readString(&classtag)
		if dec.err != nil {
			return
		}
		if classtag == "" {
			dec.err = fmt.Errorf("rootio.readClass: empty class tag")
			return
		}
		*name = classtag
		*count = int32(int64(tag) & ^kByteCountMask)
		*isref = false
	default:
		*count = int32(tag)
		*isref = true
	}
	return
}

func (dec *decoder) readClassTag(classtag *string) {
	if dec.err != nil {
		return
	}

	var tag uint32
	dec.readBin(&tag)
	if dec.err != nil {
		return
	}

	tagNewClass := tag == kNewClassTag
	tagClassMask := (int64(tag) & (^int64(kClassMask))) != 0

	if tagNewClass {
		dec.readString(classtag)
		if dec.err != nil {
			return
		}
	} else if tagClassMask {
		panic("not implemented")
	} else {
		panic(fmt.Errorf("rootio.readClassTag: unknown class-tag [%v]", tag))
	}

	return
}

func (dec *decoder) checkByteCount(pos, count int32, start int64, class string) {
	if dec.err != nil {
		return
	}

	if count == 0 {
		return
	}

	lenbuf := int64(pos) + int64(count) + 4
	diff := dec.Pos() - start
	if diff == lenbuf {
		return
	}
	dec.err = fmt.Errorf(
		"**error** [%v] diff=%v len=%v (pos=%v, count=%v, start=%v)",
		class, diff, lenbuf, pos, count, start,
	)
	panic(dec.err)
	return
}

/*
func (dec *decoder) readObject(o *Object) {
	if dec.err != nil {
		return
	}

	//start := dec.Pos()
	//orig := dec.Clone()

	var class string
	var count int32
	var isref bool
	dec.readClass(&class, &count, &isref)
	return
}
*/

func (dec *decoder) skipObject() {
	if dec.err != nil {
		return
	}
	_, dec.err = io.CopyN(ioutil.Discard, dec.buf, 10)
}

func (dec *decoder) readObject(name string) Object {
	if dec.err != nil {
		return nil
	}

	fct := Factory.get(name)
	obj := fct().Interface().(Object)

	dec.err = obj.(ROOTUnmarshaler).UnmarshalROOT(NewRBuffer(dec.buf.Bytes(), nil))
	return obj
}

func (dec *decoder) readObjectRef() Object {
	if dec.err != nil {
		return nil
	}

	fmt.Printf("--- readobjref ---\n")

	var (
		objStartPos = dec.Pos()
		tag         uint32
		vers        uint32
		startPos    int64
		bcnt        int32
	)
	dec.readBin(&bcnt)
	fmt.Printf("bcnt=%v\n", bcnt)

	if bcnt&kByteCountMask == 0 || int64(bcnt) == kNewClassTag {
		tag = uint32(bcnt)
		bcnt = 0
	} else {
		vers = 1
		startPos = dec.Pos()
		dec.readBin(&tag)
	}

	tag64 := int64(tag)

	fmt.Printf("objP= %v\n", objStartPos)
	fmt.Printf("vers= %v\n", vers)
	fmt.Printf("tag=  %v | class-mask=%v | new-class=%v\n", tag, int64(tag)&kClassMask == 0, int64(tag) == kNewClassTag)
	fmt.Printf("bcnt= %v\n", bcnt)
	fmt.Printf("spos= %v\n", startPos)

	if tag64&kClassMask == 0 {
		fmt.Printf("--> kClassMask\n")
		switch tag64 {
		case 0:
			return nil
		case 1:
			// FIXME(sbinet): tag==1 means "self", but we don't currently have self available
			panic("rootio: tag==1 'self' not implemented")
			return nil
		}

		obj := dec.refs[tag64]
		if obj == nil {
			panic(fmt.Errorf("rootio: invalid object ref [%d]", tag64))
		}
		return obj.(Object)
	}

	if tag64 == kNewClassTag {
		fmt.Printf("--> kNewClassTag\n")
		var cname string
		dec.readCString(80, &cname)
		fmt.Printf("--> class-name: %q\n", cname)

		fct := Factory.get(cname)

		if vers > 0 {
			dec.refs[startPos+kMapOffset] = fct
		} else {
			dec.refs[int64(len(dec.refs)+1)] = fct
		}

		obj := fct().Interface().(Object)

		if vers > 0 {
			dec.refs[objStartPos+kMapOffset] = obj
		} else {
			dec.refs[int64(len(dec.refs)+1)] = obj
		}
		dec.err = obj.(ROOTUnmarshaler).UnmarshalROOT(NewRBuffer(dec.buf.Bytes(), nil))
		if dec.err != nil {
			return nil
		}
		return obj

	} else {
		fmt.Printf("--> ^kClassMask\n")
		tag64 &= ^kClassMask
		fmt.Printf("--> tag64=%v\n", tag64)
	}

	fmt.Printf("tag=  %v | class-mask=%v | new-class=%v\n", tag64, tag64&kClassMask == 0, tag64 == kNewClassTag)
	return nil
}
