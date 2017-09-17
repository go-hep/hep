// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sort"
)

type rbuff struct {
	p []byte // buffer of data to read from
	c int    // current position in buffer of data
}

func (r *rbuff) Read(p []byte) (int, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	n := copy(p, r.p[r.c:])
	r.c += n
	return n, nil
}

func (r *rbuff) ReadByte() (byte, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	v := r.p[r.c]
	r.c++
	return v, nil
}

func (r *rbuff) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case ioSeekStart:
		r.c = int(offset)
	case ioSeekCurrent:
		r.c += int(offset)
	case ioSeekEnd:
		r.c = len(r.p) - int(offset)
	default:
		return 0, errors.New("rootio.rbuff.Seek: invalid whence")
	}
	if r.c < 0 {
		return 0, errors.New("rootio.rbuff.Seek: negative position")
	}
	return int64(r.c), nil
}

// RBuffer is a read-only ROOT buffer for streaming.
type RBuffer struct {
	r      *rbuff
	err    error
	offset uint32
	refs   map[int64]interface{}
}

func NewRBuffer(data []byte, refs map[int64]interface{}, offset uint32) *RBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}

	return &RBuffer{
		r:      &rbuff{p: data, c: 0},
		refs:   refs,
		offset: offset,
	}
}

func (r *RBuffer) Pos() int64 {
	return int64(r.r.c) + int64(r.offset)
}

func (r *RBuffer) setPos(pos int64) error {
	pos -= int64(r.offset)
	r.r.c = int(pos)
	return nil
}

func (r *RBuffer) Len() int64 {
	return int64(len(r.r.p) - r.r.c)
}

func (r *RBuffer) Err() error {
	return r.err
}

func (r *RBuffer) read(data []byte) {
	if r.err != nil {
		return
	}
	n := copy(data, r.r.p[r.r.c:])
	r.r.c += n
}

func (r *RBuffer) bytes() []byte {
	return r.r.p[r.r.c:]
}

func (r *RBuffer) dumpRefs() {
	fmt.Printf("--- refs ---\n")
	ids := make([]int64, 0, len(r.refs))
	for k := range r.refs {
		ids = append(ids, k)
	}
	sort.Sort(int64Slice(ids))
	for _, id := range ids {
		fmt.Printf(" id=%4d -> %v\n", id, r.refs[id])
	}
}

type int64Slice []int64

func (p int64Slice) Len() int           { return len(p) }
func (p int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (r *RBuffer) dumpHex(n int) {
	buf := r.bytes()
	if len(buf) > n {
		buf = buf[:n]
	}
	fmt.Printf("--- hex --- (pos=%d len=%d end=%d)\n%s\n", r.Pos(), n, r.Len(), string(hex.Dump(buf)))
}

func (r *RBuffer) ReadString(s *string) {
	if r.err != nil {
		return
	}

	var u8 uint8
	r.ReadU8(&u8)
	n := int(u8)
	if u8 == 255 {
		// large string
		var u32 uint32
		r.ReadU32(&u32)
		n = int(u32)
	}
	if n == 0 {
		*s = ""
		return
	}
	r.ReadU8(&u8)
	if u8 == 0 {
		*s = ""
		return
	}
	buf := make([]byte, n)
	buf[0] = u8
	if n != 0 {
		r.read(buf[1:])
		if r.err != nil {
			*s = ""
			return
		}
		*s = string(buf)
		return
	}
	*s = ""
	return
}

func (r *RBuffer) ReadCString(n int, s *string) {
	if r.err != nil {
		return
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		r.read(buf[i : i+1])
		if buf[i] == 0 {
			buf = buf[:i]
			break
		}
	}
	*s = string(buf)
}

func (r *RBuffer) ReadBool(v *bool) {
	if r.err != nil {
		return
	}

	var i8 int8
	r.ReadI8(&i8)
	if i8 != 0 {
		*v = true
		return
	}
	*v = false
}

func (r *RBuffer) ReadStaticArrayI32() []int32 {
	if r.err != nil {
		return nil
	}

	var n int32
	r.ReadI32(&n)
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int32, n)
	for i := range arr {
		r.ReadI32(&arr[i])
	}

	if r.err != nil {
		return nil
	}

	return arr
}

func (r *RBuffer) ReadFastArrayBool(v []bool) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}

	for i := range v {
		r.ReadBool(&v[i])
	}
}

func (r *RBuffer) ReadFastArrayString(v []string) {
	if r.err != nil {
		return
	}
	if n := len(v); n == 0 || int64(n) > r.Len() {
		return
	}

	for i := range v {
		r.ReadString(&v[i])
	}
}

func (r *RBuffer) ReadVersion() (vers int16, pos, n int32) {
	if r.err != nil {
		return
	}

	pos = int32(r.Pos())

	var bcnt uint32
	r.ReadU32(&bcnt)
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		n = int32(int64(bcnt) & ^kByteCountMask)
	} else {
		r.err = fmt.Errorf("rootio.ReadVersion: too old file")
		return
	}

	r.ReadI16(&vers)
	return vers, pos, n
}

func (r *RBuffer) SkipVersion(class string) {
	if r.err != nil {
		return
	}

	var version int16
	r.ReadI16(&version)

	if int64(version)&kByteCountVMask != 0 {
		var i16 int16
		r.ReadI16(&i16)
		r.ReadI16(&i16)
	}

	if class != "" && version <= 1 {
		panic("not implemented")
	}
}

func (r *RBuffer) chk(pos, count int32) bool {
	if count <= 0 {
		return true
	}

	var (
		want = int64(pos) + int64(count) + 4
		got  = r.Pos()
	)

	return got == want
}

func (r *RBuffer) CheckByteCount(pos, count int32, start int64, class string) {
	if r.err != nil {
		return
	}

	if count <= 0 {
		return
	}

	var (
		want = int64(pos) + int64(count) + 4
		got  = r.Pos()
	)

	switch {
	case got == want:
		return

	case got > want:
		r.err = fmt.Errorf("rootio.CheckByteCount: read too many bytes. got=%d, want=%d (pos=%d count=%d start=%d) [class=%q]",
			got, want, pos, count, start, class,
		)
		return

	case got < want:
		r.err = fmt.Errorf("rootio.CheckByteCount: read too few bytes. got=%d, want=%d (pos=%d count=%d start=%d) [class=%q]",
			got, want, pos, count, start, class,
		)
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
	r.r.Seek(10, ioSeekCurrent)
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

	beg := r.Pos()
	var (
		tag   uint32
		vers  int32
		start int64
		bcnt  uint32
	)
	r.ReadU32(&bcnt)

	if int64(bcnt)&kByteCountMask == 0 || int64(bcnt) == kNewClassTag {
		tag = bcnt
		bcnt = 0
	} else {
		vers = 1
		start = r.Pos()
		r.ReadU32(&tag)
	}

	tag64 := int64(tag)
	switch {
	case tag64&kClassMask == 0:
		if tag64 == 0 {
			return nil
		}
		// FIXME(sbinet): tag==1 means "self". not implemented yet.
		if tag == 1 {
			return nil
		}

		o, ok := r.refs[tag64]
		if !ok {
			r.setPos(beg + int64(bcnt) + 4)
			// r.err = fmt.Errorf("rootio: invalid tag [%v] found", tag64)
			return nil
		}
		obj, ok = o.(Object)
		if !ok {
			r.err = fmt.Errorf("rootio: invalid tag [%v] found (not a rootio.Object)", tag64)
			return nil
		}
		return obj

	case tag64 == kNewClassTag:
		var cname string
		r.ReadCString(80, &cname)
		fct := Factory.get(cname)

		if vers > 0 {
			r.refs[start+kMapOffset] = fct
		} else {
			r.refs[int64(len(r.refs))+1] = fct
		}

		obj = fct().Interface().(Object)
		if err := obj.(ROOTUnmarshaler).UnmarshalROOT(r); err != nil {
			r.err = err
			return nil
		}

		if vers > 0 {
			r.refs[beg+kMapOffset] = obj
		} else {
			r.refs[int64(len(r.refs))+1] = obj
		}
		return obj

	default:
		ref := tag64 & ^kClassMask
		cls, ok := r.refs[ref]
		if !ok {
			r.err = fmt.Errorf("rootio: invalid class-tag reference [%v] found", ref)
			return nil
		}

		fct, ok := cls.(FactoryFct)
		if !ok {
			r.err = fmt.Errorf("rootio: invalid class-tag reference [%v] found (not a rootio.FactoryFct)", ref)
			return nil
		}

		obj = fct().Interface().(Object)
		if vers > 0 {
			r.refs[beg+kMapOffset] = obj
		} else {
			r.refs[int64(len(r.refs))+1] = obj
		}

		if err := obj.(ROOTUnmarshaler).UnmarshalROOT(r); err != nil {
			r.err = err
			return nil
		}
		return obj
	}
}
