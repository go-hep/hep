// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
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
	case io.SeekStart:
		r.c = int(offset)
	case io.SeekCurrent:
		r.c += int(offset)
	case io.SeekEnd:
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
	sictx  StreamerInfoContext
}

func NewRBuffer(data []byte, refs map[int64]interface{}, offset uint32, ctx StreamerInfoContext) *RBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}

	return &RBuffer{
		r:      &rbuff{p: data, c: 0},
		refs:   refs,
		offset: offset,
		sictx:  ctx,
	}
}

func (r *RBuffer) StreamerInfo(name string) (StreamerInfo, error) {
	if r.sictx == nil {
		return nil, fmt.Errorf("rootio: no streamers")
	}
	return r.sictx.StreamerInfo(name)
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

func (r *RBuffer) ReadString() string {
	if r.err != nil {
		return ""
	}

	u8 := r.ReadU8()
	n := int(u8)
	if u8 == 255 {
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

func (r *RBuffer) ReadBool() bool {
	if r.err != nil {
		return false
	}

	v := r.ReadI8()
	if v != 0 {
		return true
	}
	return false
}

func (r *RBuffer) ReadI8() int8 {
	if r.err != nil {
		return 0
	}

	var v byte
	v, r.err = r.r.ReadByte()
	if r.err != nil {
		return 0
	}
	return int8(v)
}

func (r *RBuffer) ReadI16() int16 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 2
	v := int16(binary.BigEndian.Uint16(r.r.p[beg:r.r.c]))
	return v
}

func (r *RBuffer) ReadI32() int32 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 4
	v := int32(binary.BigEndian.Uint32(r.r.p[beg:r.r.c]))
	return v
}

func (r *RBuffer) ReadI64() int64 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 8
	v := int64(binary.BigEndian.Uint64(r.r.p[beg:r.r.c]))
	return v
}

func (r *RBuffer) ReadU8() uint8 {
	if r.err != nil {
		return 0
	}

	var v byte
	v, r.err = r.r.ReadByte()
	if r.err != nil {
		return 0
	}
	return uint8(v)
}

func (r *RBuffer) ReadU16() uint16 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 2
	v := binary.BigEndian.Uint16(r.r.p[beg:r.r.c])
	return v
}

func (r *RBuffer) ReadU32() uint32 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 4
	v := binary.BigEndian.Uint32(r.r.p[beg:r.r.c])
	return v
}

func (r *RBuffer) ReadU64() uint64 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 8
	v := binary.BigEndian.Uint64(r.r.p[beg:r.r.c])
	return v
}

func (r *RBuffer) ReadF32() float32 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 4
	v := binary.BigEndian.Uint32(r.r.p[beg:r.r.c])
	return math.Float32frombits(v)
}

func (r *RBuffer) ReadF64() float64 {
	if r.err != nil {
		return 0
	}

	beg := r.r.c
	r.r.c += 8
	v := binary.BigEndian.Uint64(r.r.p[beg:r.r.c])
	return math.Float64frombits(v)
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

func (r *RBuffer) ReadFastArrayBool(n int) []bool {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]bool, n)
	for i := range arr {
		arr[i] = r.ReadBool()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayI8(n int) []int8 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int8, n)
	for i := range arr {
		arr[i] = r.ReadI8()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayI16(n int) []int16 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int16, n)
	for i := range arr {
		arr[i] = r.ReadI16()
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

func (r *RBuffer) ReadFastArrayI64(n int) []int64 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]int64, n)
	for i := range arr {
		arr[i] = r.ReadI64()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayU8(n int) []uint8 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]uint8, n)
	for i := range arr {
		arr[i] = r.ReadU8()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayU16(n int) []uint16 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]uint16, n)
	for i := range arr {
		arr[i] = r.ReadU16()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayU32(n int) []uint32 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]uint32, n)
	for i := range arr {
		arr[i] = r.ReadU32()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayU64(n int) []uint64 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]uint64, n)
	for i := range arr {
		arr[i] = r.ReadU64()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayF32(n int) []float32 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]float32, n)
	for i := range arr {
		arr[i] = r.ReadF32()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayF64(n int) []float64 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]float64, n)
	for i := range arr {
		arr[i] = r.ReadF64()
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayString(n int) []string {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]string, n)
	for i := range arr {
		arr[i] = r.ReadString()
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
	if (int64(bcnt) & ^kByteCountMask) != 0 {
		n = int32(int64(bcnt) & ^kByteCountMask)
	} else {
		r.err = fmt.Errorf("rootio.ReadVersion: too old file")
		return
	}

	vers = int16(r.ReadU16())
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
	vers := r.ReadI16()
	if vers&kByteCountVMask != 0 {
		_, r.err = r.r.Seek(4, io.SeekCurrent)
		if r.err != nil {
			return
		}
	}
	_ = r.ReadU32() // fUniqueID
	fbits := r.ReadU32() | kIsOnHeap

	if fbits&kIsReferenced != 0 {
		_, r.err = r.r.Seek(2, io.SeekCurrent)
		if r.err != nil {
			return
		}
	}
	return
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
		bcnt  = r.ReadU32()
	)

	if int64(bcnt)&kByteCountMask == 0 || int64(bcnt) == kNewClassTag {
		tag = bcnt
		bcnt = 0
	} else {
		vers = 1
		start = r.Pos()
		tag = r.ReadU32()
	}

	tag64 := int64(tag)
	switch {
	case tag64&kClassMask == 0:
		if tag64 == 0 {
			return nil
		}
		// FIXME(sbinet): tag==1 means "self". not implemented yet.
		if tag == 1 {
			r.err = fmt.Errorf("rootio: tag == 1 means 'self'. not implemented yet")
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
		cname := r.ReadCString(80)
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
