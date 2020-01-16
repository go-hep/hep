// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"sort"

	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
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
		return 0, xerrors.Errorf("rbytes: invalid whence")
	}
	if r.c < 0 {
		return 0, xerrors.Errorf("rbytes: negative position")
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

// StreamerInfo returns the named StreamerInfo.
// If version is negative, the latest version should be returned.
func (r *RBuffer) StreamerInfo(name string, version int) (StreamerInfo, error) {
	if r.sictx == nil {
		return nil, xerrors.Errorf("rbytes: no streamers")
	}
	return r.sictx.StreamerInfo(name, version)
}

func (r *RBuffer) Pos() int64 {
	return int64(r.r.c) + int64(r.offset)
}

func (r *RBuffer) SetPos(pos int64) error { return r.setPos(pos) }
func (r *RBuffer) setPos(pos int64) error {
	pos -= int64(r.offset)
	r.r.c = int(pos)
	return nil
}

func (r *RBuffer) Len() int64 {
	return int64(len(r.r.p) - r.r.c)
}

func (r *RBuffer) Err() error       { return r.err }
func (r *RBuffer) SetErr(err error) { r.err = err }

func (r *RBuffer) read(data []byte) {
	if r.err != nil {
		return
	}
	n := copy(data, r.r.p[r.r.c:])
	r.r.c += n
}

func (r *RBuffer) Read(p []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	n, err := r.r.Read(p)
	r.err = err
	return n, r.err
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

func (r *RBuffer) DumpHex(n int) {
	buf := r.bytes()
	if len(buf) > n {
		buf = buf[:n]
	}
	fmt.Printf("--- hex --- (pos=%d len=%d end=%d)\n%s\n", r.Pos(), n, r.Len(), string(hex.Dump(buf)))
}

func (r *RBuffer) ReadSTLString() string {
	if r.Err() != nil {
		return ""
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion("string") // FIXME(sbinet): streamline with RStreamROOT
	if vers != rvers.StreamerInfo {
		r.SetErr(xerrors.Errorf("rbytes: invalid version for std::string. got=%v, want=%v", vers, rvers.StreamerInfo))
		return ""
	}

	o := r.ReadString()
	r.CheckByteCount(pos, bcnt, start, "string")

	return o
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

func (r *RBuffer) ReadF16(elm StreamerElement) root.Float16 {
	switch {
	case elm != nil && elm.Factor() != 0:
		return r.readWithFactorF16(elm.Factor(), elm.XMin())
	default:
		var nbits uint32
		if elm != nil {
			nbits = uint32(elm.XMin())
		}
		if nbits == 0 {
			nbits = 12
		}
		return r.readWithNbitsF16(nbits)
	}
}

func (r *RBuffer) readWithFactorF16(f, xmin float64) root.Float16 {
	v := float64(r.ReadU32())
	return root.Float16(v/f + xmin)
}

func (r *RBuffer) readWithNbitsF16(nbits uint32) root.Float16 {
	var (
		exp = uint32(r.ReadU8())
		man = uint32(r.ReadU16())
		val = uint32(exp)
	)
	val <<= 23
	val |= (man & ((1 << (nbits + 1)) - 1)) << (23 - nbits)

	f := math.Float32frombits(val)
	if (1 << (nbits + 1) & man) != 0 {
		f = -f
	}

	return root.Float16(f)
}

func (r *RBuffer) ReadD32(elm StreamerElement) root.Double32 {
	switch {
	case elm != nil && elm.Factor() != 0:
		return r.readWithFactorD32(elm.Factor(), elm.XMin())
	default:
		var nbits uint32
		if elm != nil {
			nbits = uint32(elm.XMin())
		}
		if nbits == 0 {
			f32 := r.ReadF32()
			return root.Double32(f32)
		}
		return r.readWithNbitsD32(nbits)
	}
}

func (r *RBuffer) readWithFactorD32(f, xmin float64) root.Double32 {
	v := float64(r.ReadU32())
	return root.Double32(v/f + xmin)
}

func (r *RBuffer) readWithNbitsD32(nbits uint32) root.Double32 {
	var (
		exp = uint32(r.ReadU8())
		man = uint32(r.ReadU16())
		val = uint32(exp)
	)
	val <<= 23
	val |= (man & ((1 << (nbits + 1)) - 1)) << (23 - nbits)

	f := math.Float32frombits(val)
	if (1 << (nbits + 1) & man) != 0 {
		f = -f
	}

	return root.Double32(f)
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

func (r *RBuffer) ReadFastArrayF16(n int, elm StreamerElement) []root.Float16 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]root.Float16, n)
	for i := range arr {
		arr[i] = r.ReadF16(elm)
	}

	if r.err != nil {
		return nil
	}
	return arr
}

func (r *RBuffer) ReadFastArrayD32(n int, elm StreamerElement) []root.Double32 {
	if r.err != nil {
		return nil
	}
	if n <= 0 || int64(n) > r.Len() {
		return nil
	}

	arr := make([]root.Double32, n)
	for i := range arr {
		arr[i] = r.ReadD32(elm)
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

func (r *RBuffer) ReadVersion(class string) (vers int16, pos, n int32) {
	if r.err != nil {
		return
	}

	pos = int32(r.Pos())

	bcnt := r.ReadU32()
	if (int64(bcnt) & kByteCountMask) != 0 {
		n = int32(int64(bcnt) & ^kByteCountMask)
		vers = int16(r.ReadU16())
	} else {
		// no byte count. rewind and read version
		r.SetPos(int64(pos))
		vers = int16(r.ReadU16())
	}

	if vers <= 0 {
		if class != "" && r.sictx != nil {
			si, err := r.sictx.StreamerInfo(class, -1)
			if err == nil && si.ClassVersion() != int(vers) {
				chksum := r.ReadU32()
				if si.CheckSum() == int(chksum) {
					vers = int16(si.ClassVersion())
				}
			}
		}
	}

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
		r.err = xerrors.Errorf("rbytes: read too many bytes. got=%d, want=%d (pos=%d count=%d start=%d) [class=%q]",
			got, want, pos, count, start, class,
		)
		return

	case got < want:
		r.err = xerrors.Errorf("rbytes: read too few bytes. got=%d, want=%d (pos=%d count=%d start=%d) [class=%q]",
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

func (r *RBuffer) ReadObject(class string) root.Object {
	if r.err != nil {
		return nil
	}

	fct := rtypes.Factory.Get(class)
	obj := fct().Interface().(root.Object)
	r.err = obj.(Unmarshaler).UnmarshalROOT(r)
	return obj
}

func (r *RBuffer) ReadObjectAny() (obj root.Object) {
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
			r.err = xerrors.Errorf("rbytes: tag == 1 means 'self'. not implemented yet")
			return nil
		}

		o, ok := r.refs[tag64]
		if !ok {
			r.setPos(beg + int64(bcnt) + 4)
			// r.err = xerrors.Errorf("rbytes: invalid tag [%v] found", tag64)
			return nil
		}
		obj, ok = o.(root.Object)
		if !ok {
			r.err = xerrors.Errorf("rbytes: invalid tag [%v] found (not a root.Object)", tag64)
			return nil
		}
		return obj

	case tag64 == kNewClassTag:
		cname := r.ReadCString(80)
		fct := rtypes.Factory.Get(cname)

		if vers > 0 {
			r.refs[start+kMapOffset] = fct
		} else {
			r.refs[int64(len(r.refs))+1] = fct
		}

		obj = fct().Interface().(root.Object)
		if err := obj.(Unmarshaler).UnmarshalROOT(r); err != nil {
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
			r.err = xerrors.Errorf("rbytes: invalid class-tag reference [%v] found", ref)
			return nil
		}

		fct, ok := cls.(rtypes.FactoryFct)
		if !ok {
			r.err = xerrors.Errorf("rbytes: invalid class-tag reference [%v] found (not a rypes.FactoryFct: %T)", ref, cls)
			return nil
		}

		obj = fct().Interface().(root.Object)
		if vers > 0 {
			r.refs[beg+kMapOffset] = obj
		} else {
			r.refs[int64(len(r.refs))+1] = obj
		}

		if err := obj.(Unmarshaler).UnmarshalROOT(r); err != nil {
			r.err = err
			return nil
		}
		return obj
	}
}

var (
	_ StreamerInfoContext = (*RBuffer)(nil)
)
