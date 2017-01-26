// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

type decoder struct {
	buf *bytes.Buffer
	len int64
	err error
}

func newDecoder(buf *bytes.Buffer) *decoder {
	dec := &decoder{
		buf: buf,
		len: int64(buf.Len()),
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

func (dec *decoder) readBin(v interface{}) {
	if dec.err != nil {
		return
	}
	dec.err = binary.Read(dec.buf, E, v)
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

	/*
	 */
	//FIXME: hack
	// var trash [8]byte
	// err = dec.readBin(&trash)
	// if err != nil {
	// 	return
	// }
	//fmt.Printf("## data = %#v\n", trash[:])

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
