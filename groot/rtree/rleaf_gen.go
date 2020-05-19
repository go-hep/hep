// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rtree

import (
	"reflect"
	"unsafe"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

// rleafValBool implements rleaf for ROOT TLeafO
type rleafValBool struct {
	base *LeafO
	v    *bool
}

func newRLeafBool(leaf *LeafO, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]bool)
		if *slice == nil {
			*slice = make([]bool, 0, rleafDefaultSliceCap)
		}
		return &rleafSliBool{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrBool{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayBool(rvar.Value)).Elem().Interface().([]bool),
		}

	default:
		return &rleafValBool{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*bool),
		}
	}
}

func (leaf *rleafValBool) Leaf() Leaf { return leaf.base }

func (leaf *rleafValBool) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValBool) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadBool()
	return r.Err()
}

var (
	_ rleaf = (*rleafValBool)(nil)
)

// rleafArrBool implements rleaf for ROOT TLeafO
type rleafArrBool struct {
	base *LeafO
	v    []bool
}

func (leaf *rleafArrBool) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrBool) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayBool(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 1
	arr := (*[0]bool)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrBool) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayBool(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrBool)(nil)
)

// rleafSliBool implements rleaf for ROOT TLeafO
type rleafSliBool struct {
	base *LeafO
	n    func() int
	v    *[]bool
}

func (leaf *rleafSliBool) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliBool) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliBool) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeBool(*leaf.v, n)
	r.ReadArrayBool(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliBool)(nil)
)

// rleafValI8 implements rleaf for ROOT TLeafB
type rleafValI8 struct {
	base *LeafB
	v    *int8
}

func newRLeafI8(leaf *LeafB, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]int8)
		if *slice == nil {
			*slice = make([]int8, 0, rleafDefaultSliceCap)
		}
		return &rleafSliI8{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrI8{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayI8(rvar.Value)).Elem().Interface().([]int8),
		}

	default:
		return &rleafValI8{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*int8),
		}
	}
}

func (leaf *rleafValI8) Leaf() Leaf { return leaf.base }

func (leaf *rleafValI8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValI8) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValI8) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadI8()
	return r.Err()
}

var (
	_ rleaf = (*rleafValI8)(nil)
)

// rleafArrI8 implements rleaf for ROOT TLeafB
type rleafArrI8 struct {
	base *LeafB
	v    []int8
}

func (leaf *rleafArrI8) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrI8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayI8(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 1
	arr := (*[0]int8)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrI8) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayI8(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrI8)(nil)
)

// rleafSliI8 implements rleaf for ROOT TLeafB
type rleafSliI8 struct {
	base *LeafB
	n    func() int
	v    *[]int8
}

func (leaf *rleafSliI8) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliI8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliI8) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeI8(*leaf.v, n)
	r.ReadArrayI8(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliI8)(nil)
)

// rleafValI16 implements rleaf for ROOT TLeafS
type rleafValI16 struct {
	base *LeafS
	v    *int16
}

func newRLeafI16(leaf *LeafS, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]int16)
		if *slice == nil {
			*slice = make([]int16, 0, rleafDefaultSliceCap)
		}
		return &rleafSliI16{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrI16{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayI16(rvar.Value)).Elem().Interface().([]int16),
		}

	default:
		return &rleafValI16{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*int16),
		}
	}
}

func (leaf *rleafValI16) Leaf() Leaf { return leaf.base }

func (leaf *rleafValI16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValI16) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValI16) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadI16()
	return r.Err()
}

var (
	_ rleaf = (*rleafValI16)(nil)
)

// rleafArrI16 implements rleaf for ROOT TLeafS
type rleafArrI16 struct {
	base *LeafS
	v    []int16
}

func (leaf *rleafArrI16) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrI16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayI16(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 2
	arr := (*[0]int16)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrI16) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayI16(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrI16)(nil)
)

// rleafSliI16 implements rleaf for ROOT TLeafS
type rleafSliI16 struct {
	base *LeafS
	n    func() int
	v    *[]int16
}

func (leaf *rleafSliI16) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliI16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliI16) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeI16(*leaf.v, n)
	r.ReadArrayI16(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliI16)(nil)
)

// rleafValI32 implements rleaf for ROOT TLeafI
type rleafValI32 struct {
	base *LeafI
	v    *int32
}

func newRLeafI32(leaf *LeafI, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]int32)
		if *slice == nil {
			*slice = make([]int32, 0, rleafDefaultSliceCap)
		}
		return &rleafSliI32{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrI32{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayI32(rvar.Value)).Elem().Interface().([]int32),
		}

	default:
		return &rleafValI32{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*int32),
		}
	}
}

func (leaf *rleafValI32) Leaf() Leaf { return leaf.base }

func (leaf *rleafValI32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValI32) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValI32) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadI32()
	return r.Err()
}

var (
	_ rleaf = (*rleafValI32)(nil)
)

// rleafArrI32 implements rleaf for ROOT TLeafI
type rleafArrI32 struct {
	base *LeafI
	v    []int32
}

func (leaf *rleafArrI32) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrI32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayI32(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]int32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrI32) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayI32(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrI32)(nil)
)

// rleafSliI32 implements rleaf for ROOT TLeafI
type rleafSliI32 struct {
	base *LeafI
	n    func() int
	v    *[]int32
}

func (leaf *rleafSliI32) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliI32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliI32) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeI32(*leaf.v, n)
	r.ReadArrayI32(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliI32)(nil)
)

// rleafValI64 implements rleaf for ROOT TLeafL
type rleafValI64 struct {
	base *LeafL
	v    *int64
}

func newRLeafI64(leaf *LeafL, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]int64)
		if *slice == nil {
			*slice = make([]int64, 0, rleafDefaultSliceCap)
		}
		return &rleafSliI64{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrI64{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayI64(rvar.Value)).Elem().Interface().([]int64),
		}

	default:
		return &rleafValI64{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*int64),
		}
	}
}

func (leaf *rleafValI64) Leaf() Leaf { return leaf.base }

func (leaf *rleafValI64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValI64) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValI64) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadI64()
	return r.Err()
}

var (
	_ rleaf = (*rleafValI64)(nil)
)

// rleafArrI64 implements rleaf for ROOT TLeafL
type rleafArrI64 struct {
	base *LeafL
	v    []int64
}

func (leaf *rleafArrI64) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrI64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayI64(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]int64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrI64) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayI64(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrI64)(nil)
)

// rleafSliI64 implements rleaf for ROOT TLeafL
type rleafSliI64 struct {
	base *LeafL
	n    func() int
	v    *[]int64
}

func (leaf *rleafSliI64) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliI64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliI64) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeI64(*leaf.v, n)
	r.ReadArrayI64(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliI64)(nil)
)

// rleafValU8 implements rleaf for ROOT TLeafB
type rleafValU8 struct {
	base *LeafB
	v    *uint8
}

func newRLeafU8(leaf *LeafB, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]uint8)
		if *slice == nil {
			*slice = make([]uint8, 0, rleafDefaultSliceCap)
		}
		return &rleafSliU8{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrU8{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayU8(rvar.Value)).Elem().Interface().([]uint8),
		}

	default:
		return &rleafValU8{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*uint8),
		}
	}
}

func (leaf *rleafValU8) Leaf() Leaf { return leaf.base }

func (leaf *rleafValU8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValU8) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValU8) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadU8()
	return r.Err()
}

var (
	_ rleaf = (*rleafValU8)(nil)
)

// rleafArrU8 implements rleaf for ROOT TLeafB
type rleafArrU8 struct {
	base *LeafB
	v    []uint8
}

func (leaf *rleafArrU8) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrU8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayU8(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 1
	arr := (*[0]uint8)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrU8) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayU8(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrU8)(nil)
)

// rleafSliU8 implements rleaf for ROOT TLeafB
type rleafSliU8 struct {
	base *LeafB
	n    func() int
	v    *[]uint8
}

func (leaf *rleafSliU8) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliU8) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliU8) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeU8(*leaf.v, n)
	r.ReadArrayU8(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliU8)(nil)
)

// rleafValU16 implements rleaf for ROOT TLeafS
type rleafValU16 struct {
	base *LeafS
	v    *uint16
}

func newRLeafU16(leaf *LeafS, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]uint16)
		if *slice == nil {
			*slice = make([]uint16, 0, rleafDefaultSliceCap)
		}
		return &rleafSliU16{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrU16{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayU16(rvar.Value)).Elem().Interface().([]uint16),
		}

	default:
		return &rleafValU16{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*uint16),
		}
	}
}

func (leaf *rleafValU16) Leaf() Leaf { return leaf.base }

func (leaf *rleafValU16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValU16) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValU16) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadU16()
	return r.Err()
}

var (
	_ rleaf = (*rleafValU16)(nil)
)

// rleafArrU16 implements rleaf for ROOT TLeafS
type rleafArrU16 struct {
	base *LeafS
	v    []uint16
}

func (leaf *rleafArrU16) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrU16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayU16(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 2
	arr := (*[0]uint16)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrU16) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayU16(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrU16)(nil)
)

// rleafSliU16 implements rleaf for ROOT TLeafS
type rleafSliU16 struct {
	base *LeafS
	n    func() int
	v    *[]uint16
}

func (leaf *rleafSliU16) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliU16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliU16) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeU16(*leaf.v, n)
	r.ReadArrayU16(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliU16)(nil)
)

// rleafValU32 implements rleaf for ROOT TLeafI
type rleafValU32 struct {
	base *LeafI
	v    *uint32
}

func newRLeafU32(leaf *LeafI, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]uint32)
		if *slice == nil {
			*slice = make([]uint32, 0, rleafDefaultSliceCap)
		}
		return &rleafSliU32{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrU32{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayU32(rvar.Value)).Elem().Interface().([]uint32),
		}

	default:
		return &rleafValU32{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*uint32),
		}
	}
}

func (leaf *rleafValU32) Leaf() Leaf { return leaf.base }

func (leaf *rleafValU32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValU32) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValU32) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadU32()
	return r.Err()
}

var (
	_ rleaf = (*rleafValU32)(nil)
)

// rleafArrU32 implements rleaf for ROOT TLeafI
type rleafArrU32 struct {
	base *LeafI
	v    []uint32
}

func (leaf *rleafArrU32) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrU32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayU32(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]uint32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrU32) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayU32(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrU32)(nil)
)

// rleafSliU32 implements rleaf for ROOT TLeafI
type rleafSliU32 struct {
	base *LeafI
	n    func() int
	v    *[]uint32
}

func (leaf *rleafSliU32) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliU32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliU32) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeU32(*leaf.v, n)
	r.ReadArrayU32(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliU32)(nil)
)

// rleafValU64 implements rleaf for ROOT TLeafL
type rleafValU64 struct {
	base *LeafL
	v    *uint64
}

func newRLeafU64(leaf *LeafL, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]uint64)
		if *slice == nil {
			*slice = make([]uint64, 0, rleafDefaultSliceCap)
		}
		return &rleafSliU64{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrU64{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayU64(rvar.Value)).Elem().Interface().([]uint64),
		}

	default:
		return &rleafValU64{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*uint64),
		}
	}
}

func (leaf *rleafValU64) Leaf() Leaf { return leaf.base }

func (leaf *rleafValU64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValU64) ivalue() int { return int(*leaf.v) }

func (leaf *rleafValU64) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadU64()
	return r.Err()
}

var (
	_ rleaf = (*rleafValU64)(nil)
)

// rleafArrU64 implements rleaf for ROOT TLeafL
type rleafArrU64 struct {
	base *LeafL
	v    []uint64
}

func (leaf *rleafArrU64) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrU64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayU64(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]uint64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrU64) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayU64(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrU64)(nil)
)

// rleafSliU64 implements rleaf for ROOT TLeafL
type rleafSliU64 struct {
	base *LeafL
	n    func() int
	v    *[]uint64
}

func (leaf *rleafSliU64) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliU64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliU64) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeU64(*leaf.v, n)
	r.ReadArrayU64(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliU64)(nil)
)

// rleafValF32 implements rleaf for ROOT TLeafF
type rleafValF32 struct {
	base *LeafF
	v    *float32
}

func newRLeafF32(leaf *LeafF, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]float32)
		if *slice == nil {
			*slice = make([]float32, 0, rleafDefaultSliceCap)
		}
		return &rleafSliF32{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrF32{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayF32(rvar.Value)).Elem().Interface().([]float32),
		}

	default:
		return &rleafValF32{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*float32),
		}
	}
}

func (leaf *rleafValF32) Leaf() Leaf { return leaf.base }

func (leaf *rleafValF32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValF32) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadF32()
	return r.Err()
}

var (
	_ rleaf = (*rleafValF32)(nil)
)

// rleafArrF32 implements rleaf for ROOT TLeafF
type rleafArrF32 struct {
	base *LeafF
	v    []float32
}

func (leaf *rleafArrF32) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrF32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayF32(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]float32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrF32) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayF32(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrF32)(nil)
)

// rleafSliF32 implements rleaf for ROOT TLeafF
type rleafSliF32 struct {
	base *LeafF
	n    func() int
	v    *[]float32
}

func (leaf *rleafSliF32) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliF32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliF32) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeF32(*leaf.v, n)
	r.ReadArrayF32(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliF32)(nil)
)

// rleafValF64 implements rleaf for ROOT TLeafD
type rleafValF64 struct {
	base *LeafD
	v    *float64
}

func newRLeafF64(leaf *LeafD, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]float64)
		if *slice == nil {
			*slice = make([]float64, 0, rleafDefaultSliceCap)
		}
		return &rleafSliF64{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrF64{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayF64(rvar.Value)).Elem().Interface().([]float64),
		}

	default:
		return &rleafValF64{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*float64),
		}
	}
}

func (leaf *rleafValF64) Leaf() Leaf { return leaf.base }

func (leaf *rleafValF64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValF64) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadF64()
	return r.Err()
}

var (
	_ rleaf = (*rleafValF64)(nil)
)

// rleafArrF64 implements rleaf for ROOT TLeafD
type rleafArrF64 struct {
	base *LeafD
	v    []float64
}

func (leaf *rleafArrF64) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrF64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayF64(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]float64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrF64) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayF64(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrF64)(nil)
)

// rleafSliF64 implements rleaf for ROOT TLeafD
type rleafSliF64 struct {
	base *LeafD
	n    func() int
	v    *[]float64
}

func (leaf *rleafSliF64) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliF64) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliF64) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeF64(*leaf.v, n)
	r.ReadArrayF64(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliF64)(nil)
)

// rleafValD32 implements rleaf for ROOT TLeafD32
type rleafValD32 struct {
	base *LeafD32
	v    *root.Double32
	elm  rbytes.StreamerElement
}

func newRLeafD32(leaf *LeafD32, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]root.Double32)
		if *slice == nil {
			*slice = make([]root.Double32, 0, rleafDefaultSliceCap)
		}
		return &rleafSliD32{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrD32{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayD32(rvar.Value)).Elem().Interface().([]root.Double32),
		}

	default:
		return &rleafValD32{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*root.Double32),
		}
	}
}

func (leaf *rleafValD32) Leaf() Leaf { return leaf.base }

func (leaf *rleafValD32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValD32) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadD32(leaf.elm)
	return r.Err()
}

var (
	_ rleaf = (*rleafValD32)(nil)
)

// rleafArrD32 implements rleaf for ROOT TLeafD32
type rleafArrD32 struct {
	base *LeafD32
	v    []root.Double32
	elm  rbytes.StreamerElement
}

func (leaf *rleafArrD32) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrD32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayD32(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]root.Double32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrD32) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayD32(leaf.v, leaf.elm)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrD32)(nil)
)

// rleafSliD32 implements rleaf for ROOT TLeafD32
type rleafSliD32 struct {
	base *LeafD32
	n    func() int
	v    *[]root.Double32
	elm  rbytes.StreamerElement
}

func (leaf *rleafSliD32) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliD32) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliD32) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeD32(*leaf.v, n)
	r.ReadArrayD32(sli, leaf.elm)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliD32)(nil)
)

// rleafValF16 implements rleaf for ROOT TLeafF16
type rleafValF16 struct {
	base *LeafF16
	v    *root.Float16
	elm  rbytes.StreamerElement
}

func newRLeafF16(leaf *LeafF16, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]root.Float16)
		if *slice == nil {
			*slice = make([]root.Float16, 0, rleafDefaultSliceCap)
		}
		return &rleafSliF16{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrF16{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayF16(rvar.Value)).Elem().Interface().([]root.Float16),
		}

	default:
		return &rleafValF16{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*root.Float16),
		}
	}
}

func (leaf *rleafValF16) Leaf() Leaf { return leaf.base }

func (leaf *rleafValF16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValF16) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadF16(leaf.elm)
	return r.Err()
}

var (
	_ rleaf = (*rleafValF16)(nil)
)

// rleafArrF16 implements rleaf for ROOT TLeafF16
type rleafArrF16 struct {
	base *LeafF16
	v    []root.Float16
	elm  rbytes.StreamerElement
}

func (leaf *rleafArrF16) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrF16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayF16(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]root.Float16)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrF16) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayF16(leaf.v, leaf.elm)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrF16)(nil)
)

// rleafSliF16 implements rleaf for ROOT TLeafF16
type rleafSliF16 struct {
	base *LeafF16
	n    func() int
	v    *[]root.Float16
	elm  rbytes.StreamerElement
}

func (leaf *rleafSliF16) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliF16) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliF16) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeF16(*leaf.v, n)
	r.ReadArrayF16(sli, leaf.elm)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliF16)(nil)
)

// rleafValStr implements rleaf for ROOT TLeafC
type rleafValStr struct {
	base *LeafC
	v    *string
}

func newRLeafStr(leaf *LeafC, rvar ReadVar, rctx rleafCtx) rleaf {
	switch {
	case leaf.count != nil:
		slice := reflect.ValueOf(rvar.Value).Interface().(*[]string)
		if *slice == nil {
			*slice = make([]string, 0, rleafDefaultSliceCap)
		}
		return &rleafSliStr{
			base: leaf,
			n:    rctx.rcountFunc(leaf.count.Name()),
			v:    slice,
		}

	case leaf.len > 1:
		return &rleafArrStr{
			base: leaf,
			v:    reflect.ValueOf(unsafeDecayArrayStr(rvar.Value)).Elem().Interface().([]string),
		}

	default:
		return &rleafValStr{
			base: leaf,
			v:    reflect.ValueOf(rvar.Value).Interface().(*string),
		}
	}
}

func (leaf *rleafValStr) Leaf() Leaf { return leaf.base }

func (leaf *rleafValStr) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafValStr) readFromBuffer(r *rbytes.RBuffer) error {
	*leaf.v = r.ReadString()
	return r.Err()
}

var (
	_ rleaf = (*rleafValStr)(nil)
)

// rleafArrStr implements rleaf for ROOT TLeafC
type rleafArrStr struct {
	base *LeafC
	v    []string
}

func (leaf *rleafArrStr) Leaf() Leaf { return leaf.base }

func (leaf *rleafArrStr) Offset() int64 {
	return int64(leaf.base.Offset())
}

func unsafeDecayArrayStr(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 16
	arr := (*[0]string)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *rleafArrStr) readFromBuffer(r *rbytes.RBuffer) error {
	r.ReadArrayString(leaf.v)
	return r.Err()
}

var (
	_ rleaf = (*rleafArrStr)(nil)
)

// rleafSliStr implements rleaf for ROOT TLeafC
type rleafSliStr struct {
	base *LeafC
	n    func() int
	v    *[]string
}

func (leaf *rleafSliStr) Leaf() Leaf { return leaf.base }

func (leaf *rleafSliStr) Offset() int64 {
	return int64(leaf.base.Offset())
}

func (leaf *rleafSliStr) readFromBuffer(r *rbytes.RBuffer) error {
	n := leaf.base.tleaf.len * leaf.n()
	sli := rbytes.ResizeStr(*leaf.v, n)
	r.ReadArrayString(sli)
	*leaf.v = sli
	return r.Err()
}

var (
	_ rleaf = (*rleafSliStr)(nil)
)
