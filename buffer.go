// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

func (k *Key) DecodeVector(in *bytes.Buffer, dst interface{}) (int, error) {
	// Discard three int16s (like 40 00 00 0e 00 09)
	x := in.Next(6)
	_ = x // sometimes we want to look at this.

	var n int32
	err := binary.Read(in, binary.BigEndian, &n)
	if err != nil {
		return -1, err
	}

	err = binary.Read(in, binary.BigEndian, reflect.ValueOf(dst).Slice(0, int(n)).Interface())
	if err != nil {
		return -1, err
	}
	return int(n), nil
}

// RBuffer is a read-only ROOT buffer for streaming.
type RBuffer struct {
	r    *bytes.Reader
	err  error
	klen uint32
	refs map[int64]interface{}
}

func NewRBuffer(data []byte, refs map[int64]interface{}, klen uint32) *RBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}

	return &RBuffer{
		r:    bytes.NewReader(data),
		refs: refs,
		klen: klen,
	}
}

func (r *RBuffer) Pos() int64 {
	pos, _ := r.r.Seek(0, io.SeekCurrent)
	return pos
}

func (r *RBuffer) Len() int64 {
	return int64(r.r.Len())
}

func (r *RBuffer) Err() error {
	return r.err
}

func (r *RBuffer) read(data []byte) {
	if r.err != nil {
		return
	}
	_, r.err = io.ReadFull(r.r, data)
}

func (r *RBuffer) bytes() []byte {
	pos := r.Pos()
	out := make([]byte, int(r.Len()))
	io.ReadFull(r.r, out)
	r.r.Seek(pos, io.SeekStart)
	return out
}

func (r *RBuffer) ReadString() string {
	if r.err != nil {
		return ""
	}

	n := int(r.ReadU8())
	if n == 255 {
		// large string
		n = int(r.ReadU32())
	}
	if n == 0 {
		return ""
	}
	v := r.ReadU8()
	if v == 0 {
		return ""
	}
	buf := make([]byte, n)
	buf[0] = v
	if n != 0 {
		r.read(buf[1:])
		if r.err != nil {
			return ""
		}
		return string(buf)
	}
	return ""
}

func (r *RBuffer) ReadCString(n int) string {
	if r.err != nil {
		return ""
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		r.read(buf[i : i+1])
		if buf[i] == 0 {
			buf = buf[:i]
			break
		}
	}
	return string(buf)
}

func (r *RBuffer) ReadI8() int8 {
	if r.err != nil {
		return 0
	}

	var buf [1]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return int8(buf[0])
}

func (r *RBuffer) ReadI16() int16 {
	if r.err != nil {
		return 0
	}

	var buf [2]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return int16(binary.BigEndian.Uint16(buf[:]))
}

func (r *RBuffer) ReadI32() int32 {
	if r.err != nil {
		return 0
	}

	var buf [4]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return int32(binary.BigEndian.Uint32(buf[:]))
}

func (r *RBuffer) ReadI64() int64 {
	if r.err != nil {
		return 0
	}

	var buf [8]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return int64(binary.BigEndian.Uint64(buf[:]))
}

func (r *RBuffer) ReadU8() uint8 {
	if r.err != nil {
		return 0
	}

	var buf [1]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return uint8(buf[0])
}

func (r *RBuffer) ReadU16() uint16 {
	if r.err != nil {
		return 0
	}

	var buf [2]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return binary.BigEndian.Uint16(buf[:])
}

func (r *RBuffer) ReadU32() uint32 {
	if r.err != nil {
		return 0
	}

	var buf [4]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return binary.BigEndian.Uint32(buf[:])
}

func (r *RBuffer) ReadU64() uint64 {
	if r.err != nil {
		return 0
	}

	var buf [8]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return binary.BigEndian.Uint64(buf[:])
}

func (r *RBuffer) ReadF32() float32 {
	if r.err != nil {
		return 0
	}

	var buf [4]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return math.Float32frombits(binary.BigEndian.Uint32(buf[:]))
}

func (r *RBuffer) ReadF64() float64 {
	if r.err != nil {
		return 0
	}

	var buf [8]byte
	_, r.err = io.ReadFull(r.r, buf[:])
	if r.err != nil {
		return 0
	}
	return math.Float64frombits(binary.BigEndian.Uint64(buf[:]))
}

func (r *RBuffer) ReadStaticArrayI32() []int32 {
	if r.err != nil {
		return nil
	}

	n := int(r.ReadI32())
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int32, n)
	for i := range arr {
		arr[i] = r.ReadI32()
	}

	if r.err != nil {
		return nil
	}

	return arr
}

func (r *RBuffer) ReadFastArrayI32(n int) []int32 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int32, n)
	for i := range arr {
		arr[i] = r.ReadI32()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadVersion() (vers int16, pos, n int32) {
	if r.err != nil {
		return
	}

	pos = int32(r.Pos())

	bcnt := r.ReadU32()
	myprintf("readVersion - bytecount=%v\n", bcnt)
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		n = int32(int64(bcnt) & ^kByteCountMask)
	} else {
		r.err = fmt.Errorf("rootio.ReadVersion: too old file")
		return
	}

	vers = int16(r.ReadU16())
	myprintf("readVersion => [%v] [%v] [%v]\n", pos, vers, n)
	return vers, pos, n
}

func (r *RBuffer) SkipVersion(class string) {
	if r.err != nil {
		return
	}

	version := r.ReadI16()

	if int64(version)&kByteCountVMask != 0 {
		_ = r.ReadI16()
		_ = r.ReadI16()
	}

	if class != "" && version <= 1 {
		panic("not implemented")
	}
}

func (r *RBuffer) CheckByteCount(pos, count int32, start int64, class string) {
	if r.err != nil {
		return
	}

	if count <= 0 {
		return
	}

	var (
		n    = int64(pos) + int64(count) + 4
		diff = r.Pos()
	)

	if diff == n {
		return
	}

	if diff != n {
		r.err = fmt.Errorf("rootio.CheckByteCount: len=%d, want=%d (pos=%d count=%d start=%d) [class=%q]",
			n, diff, pos, count, start, class,
		)
		fmt.Printf("*** err: %v\n", r.err)
		panic(r.err) // FIXME(sbinet)
		return
	}

	return
}

func (r *RBuffer) SkipObject() {
	if r.err != nil {
		return
	}
	//v, pos, n := r.ReadVersion()
	//fmt.Printf("--- skip-object: v=%d pos=%d n=%d\n", v, pos, n)
	r.r.Seek(10, io.SeekCurrent)
	//_, r.err = r.r.Seek(int64(n), io.SeekCurrent)
}

func (r *RBuffer) ReadObject(class string) Object {
	if r.err != nil {
		return nil
	}

	fct := Factory.get(class)
	obj := fct().Interface().(Object)
	r.err = obj.(ROOTUnmarshaler).UnmarshalROOT(r)
	return obj
}

func (r *RBuffer) ReadObjectAny() (obj Object) {
	if r.err != nil {
		return obj
	}

	start := r.Pos()

	name, count, isref := r.ReadClass()
	if isref {
		panic("rootio: not implemented")
	} else {
		switch name {
		case "":
			obj = nil
			r.r.Seek(start+int64(count)+4, io.SeekStart)
		default:
			fct := Factory.get(name)
			obj = fct().Interface().(Object)
			if err := obj.(ROOTUnmarshaler).UnmarshalROOT(r); err != nil {
				r.err = err
			}
			r.r.Seek(start+int64(count)+4, io.SeekStart)
		}
	}

	return obj
}

func (r *RBuffer) ReadClass() (name string, count uint32, isref bool) {
	if r.err != nil {
		return
	}

	i := r.ReadU32()
	switch {
	case i == kNullTag:
		fmt.Printf("+++ kNullTag\n")

	case i&kByteCountMask != 0:
		clstag := r.ReadClassTag()
		if clstag == "" {
			panic("rootio: empty class tag")
		}
		name = clstag
		count = uint32(int64(i) & ^kByteCountMask)
	default:
		count = uint32(i)
		isref = true
	}
	return name, count, isref
}

func (r *RBuffer) ReadClassTag() (clstag string) {
	if r.err != nil {
		return ""
	}

	tag := r.ReadU32()
	switch {
	case tag == kNewClassTag:
		clstag = r.ReadCString(80)

	case (int64(tag) & int64(kClassMask)) != 0:
		ref := uint32(int64(tag) & ^kClassMask)
		ref -= r.klen
		pos := r.Pos()
		defer r.r.Seek(pos, io.SeekStart)
		r.r.Seek(-(pos - int64(ref-kMapOffset)), io.SeekCurrent)
		clstag = r.ReadClassTag()
	default:
		panic(fmt.Errorf("rootio: unknown class-tag: %v\n", tag))
	}
	return clstag
}

/*
func (r *RBuffer) ReadObjectRef() Object {
	if r.err != nil {
		return nil
	}

	var (
		objStartPos = r.Pos()
		tag         uint32
		vers        uint32
		startPos    int64
		bcnt        = r.ReadI32()
	)

	if bcnt&kByteCountMask == 0 || int64(bcnt) == kNewClassTag {
		tag = uint32(bcnt)
		bcnt = 0
	} else {
		vers = 1
		startPos = r.Pos()
		tag = r.ReadU32()
	}

	tag64 := int64(tag)

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

		obj := r.refs[tag64]
		if obj == nil {
			panic(fmt.Errorf("rootio: invalid object ref [%d]", tag64))
		}
		return obj.(Object)
	}

	if tag64 == kNewClassTag {
		cname := r.ReadCString(80)

		fct := Factory.get(cname)

		if vers > 0 {
			r.refs[startPos+kMapOffset] = fct
		} else {
			r.refs[int64(len(r.refs)+1)] = fct
		}

		obj := fct().Interface().(Object)

		if vers > 0 {
			r.refs[objStartPos+kMapOffset] = obj
		} else {
			r.refs[int64(len(r.refs)+1)] = obj
		}
		if r.err != nil {
			return nil
		}

		r.err = obj.(ROOTUnmarshaler).UnmarshalROOT(r)
		if r.err != nil {
			return nil
		}
		return obj

	} else {
		tag64 &= ^kClassMask
	}

	return nil
}
*/
