// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rtree

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// LeafO implements ROOT TLeafO
type LeafO struct {
	rvers int16
	tleaf
	ptr *bool
	sli *[]bool
	min bool
	max bool
}

func newLeafO(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafO {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafO{
		rvers: rvers.LeafO,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafO) Class() string {
	return "TLeafO"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafO) Minimum() bool {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafO) Maximum() bool {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafO) Kind() reflect.Kind {
	return reflect.Bool
}

// Type returns the leaf's type.
func (leaf *LeafO) Type() reflect.Type {
	var v bool
	return reflect.TypeOf(v)
}

func (leaf *LeafO) TypeName() string {
	return "bool"
}

func (leaf *LeafO) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteBool(leaf.min)
	w.WriteBool(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafO) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafO {
		panic(fmt.Errorf("rtree: invalid TLeafO version=%d > %d", hdr.Vers, rvers.LeafO))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadBool()
	leaf.max = r.ReadBool()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafO) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadBool()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeBool(*leaf.sli, nn)
			r.ReadArrayBool(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeBool(*leaf.sli, nn)
			r.ReadArrayBool(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafO) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 1
	arr := (*[0]bool)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafO) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]bool:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *bool:
		leaf.ptr = v
	case *[]bool:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]bool, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafO) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteBool(*leaf.ptr)
		nbytes += leaf.tleaf.etype
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayBool((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayBool((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafO{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafO", f)
}

var (
	_ root.Object        = (*LeafO)(nil)
	_ root.Named         = (*LeafO)(nil)
	_ Leaf               = (*LeafO)(nil)
	_ rbytes.Marshaler   = (*LeafO)(nil)
	_ rbytes.Unmarshaler = (*LeafO)(nil)
)

// LeafB implements ROOT TLeafB
type LeafB struct {
	rvers int16
	tleaf
	ptr *int8
	sli *[]int8
	min int8
	max int8
}

func newLeafB(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafB {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafB{
		rvers: rvers.LeafB,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafB) Class() string {
	return "TLeafB"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafB) Minimum() int8 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafB) Maximum() int8 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafB) Kind() reflect.Kind {
	if leaf.IsUnsigned() {
		return reflect.Uint8
	}
	return reflect.Int8
}

// Type returns the leaf's type.
func (leaf *LeafB) Type() reflect.Type {
	if leaf.IsUnsigned() {
		var v uint8
		return reflect.TypeOf(v)
	}
	var v int8
	return reflect.TypeOf(v)
}

// ivalue returns the first leaf value as int
func (leaf *LeafB) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *LeafB) imax() int {
	return int(leaf.max)
}

func (leaf *LeafB) TypeName() string {
	if leaf.IsUnsigned() {
		return "uint8"
	}
	return "int8"
}

func (leaf *LeafB) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI8(leaf.min)
	w.WriteI8(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafB) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafB {
		panic(fmt.Errorf("rtree: invalid TLeafB version=%d > %d", hdr.Vers, rvers.LeafB))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI8()
	leaf.max = r.ReadI8()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafB) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadI8()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeI8(*leaf.sli, nn)
			r.ReadArrayI8(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeI8(*leaf.sli, nn)
			r.ReadArrayI8(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafB) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 1
	arr := (*[0]int8)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafB) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]int8:
			return leaf.setAddress(sli)
		case *[]uint8:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int8:
		leaf.ptr = v
	case *[]int8:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]int8, 0)
		}
	case *uint8:
		leaf.ptr = (*int8)(unsafe.Pointer(v))
	case *[]uint8:
		leaf.sli = (*[]int8)(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]int8, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafB) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteI8(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayI8((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayI8((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafB{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafB", f)
}

var (
	_ root.Object        = (*LeafB)(nil)
	_ root.Named         = (*LeafB)(nil)
	_ Leaf               = (*LeafB)(nil)
	_ rbytes.Marshaler   = (*LeafB)(nil)
	_ rbytes.Unmarshaler = (*LeafB)(nil)
)

// LeafS implements ROOT TLeafS
type LeafS struct {
	rvers int16
	tleaf
	ptr *int16
	sli *[]int16
	min int16
	max int16
}

func newLeafS(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafS {
	const etype = 2
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafS{
		rvers: rvers.LeafS,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafS) Class() string {
	return "TLeafS"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafS) Minimum() int16 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafS) Maximum() int16 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafS) Kind() reflect.Kind {
	if leaf.IsUnsigned() {
		return reflect.Uint16
	}
	return reflect.Int16
}

// Type returns the leaf's type.
func (leaf *LeafS) Type() reflect.Type {
	if leaf.IsUnsigned() {
		var v uint16
		return reflect.TypeOf(v)
	}
	var v int16
	return reflect.TypeOf(v)
}

// ivalue returns the first leaf value as int
func (leaf *LeafS) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *LeafS) imax() int {
	return int(leaf.max)
}

func (leaf *LeafS) TypeName() string {
	if leaf.IsUnsigned() {
		return "uint16"
	}
	return "int16"
}

func (leaf *LeafS) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI16(leaf.min)
	w.WriteI16(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafS) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafS {
		panic(fmt.Errorf("rtree: invalid TLeafS version=%d > %d", hdr.Vers, rvers.LeafS))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI16()
	leaf.max = r.ReadI16()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafS) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadI16()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeI16(*leaf.sli, nn)
			r.ReadArrayI16(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeI16(*leaf.sli, nn)
			r.ReadArrayI16(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafS) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 2
	arr := (*[0]int16)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafS) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]int16:
			return leaf.setAddress(sli)
		case *[]uint16:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int16:
		leaf.ptr = v
	case *[]int16:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]int16, 0)
		}
	case *uint16:
		leaf.ptr = (*int16)(unsafe.Pointer(v))
	case *[]uint16:
		leaf.sli = (*[]int16)(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]int16, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafS) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteI16(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayI16((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayI16((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafS{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafS", f)
}

var (
	_ root.Object        = (*LeafS)(nil)
	_ root.Named         = (*LeafS)(nil)
	_ Leaf               = (*LeafS)(nil)
	_ rbytes.Marshaler   = (*LeafS)(nil)
	_ rbytes.Unmarshaler = (*LeafS)(nil)
)

// LeafI implements ROOT TLeafI
type LeafI struct {
	rvers int16
	tleaf
	ptr *int32
	sli *[]int32
	min int32
	max int32
}

func newLeafI(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafI {
	const etype = 4
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafI{
		rvers: rvers.LeafI,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafI) Class() string {
	return "TLeafI"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafI) Minimum() int32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafI) Maximum() int32 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafI) Kind() reflect.Kind {
	if leaf.IsUnsigned() {
		return reflect.Uint32
	}
	return reflect.Int32
}

// Type returns the leaf's type.
func (leaf *LeafI) Type() reflect.Type {
	if leaf.IsUnsigned() {
		var v uint32
		return reflect.TypeOf(v)
	}
	var v int32
	return reflect.TypeOf(v)
}

// ivalue returns the first leaf value as int
func (leaf *LeafI) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *LeafI) imax() int {
	return int(leaf.max)
}

func (leaf *LeafI) TypeName() string {
	if leaf.IsUnsigned() {
		return "uint32"
	}
	return "int32"
}

func (leaf *LeafI) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI32(leaf.min)
	w.WriteI32(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafI) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafI {
		panic(fmt.Errorf("rtree: invalid TLeafI version=%d > %d", hdr.Vers, rvers.LeafI))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafI) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadI32()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeI32(*leaf.sli, nn)
			r.ReadArrayI32(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeI32(*leaf.sli, nn)
			r.ReadArrayI32(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafI) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]int32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafI) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]int32:
			return leaf.setAddress(sli)
		case *[]uint32:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int32:
		leaf.ptr = v
	case *[]int32:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]int32, 0)
		}
	case *uint32:
		leaf.ptr = (*int32)(unsafe.Pointer(v))
	case *[]uint32:
		leaf.sli = (*[]int32)(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]int32, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafI) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteI32(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayI32((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayI32((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafI{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafI", f)
}

var (
	_ root.Object        = (*LeafI)(nil)
	_ root.Named         = (*LeafI)(nil)
	_ Leaf               = (*LeafI)(nil)
	_ rbytes.Marshaler   = (*LeafI)(nil)
	_ rbytes.Unmarshaler = (*LeafI)(nil)
)

// LeafL implements ROOT TLeafL
type LeafL struct {
	rvers int16
	tleaf
	ptr *int64
	sli *[]int64
	min int64
	max int64
}

func newLeafL(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafL {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafL{
		rvers: rvers.LeafL,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafL) Class() string {
	return "TLeafL"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafL) Minimum() int64 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafL) Maximum() int64 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafL) Kind() reflect.Kind {
	if leaf.IsUnsigned() {
		return reflect.Uint64
	}
	return reflect.Int64
}

// Type returns the leaf's type.
func (leaf *LeafL) Type() reflect.Type {
	if leaf.IsUnsigned() {
		var v uint64
		return reflect.TypeOf(v)
	}
	var v int64
	return reflect.TypeOf(v)
}

// ivalue returns the first leaf value as int
func (leaf *LeafL) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *LeafL) imax() int {
	return int(leaf.max)
}

func (leaf *LeafL) TypeName() string {
	if leaf.IsUnsigned() {
		return "uint64"
	}
	return "int64"
}

func (leaf *LeafL) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI64(leaf.min)
	w.WriteI64(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafL) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafL {
		panic(fmt.Errorf("rtree: invalid TLeafL version=%d > %d", hdr.Vers, rvers.LeafL))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI64()
	leaf.max = r.ReadI64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafL) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadI64()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeI64(*leaf.sli, nn)
			r.ReadArrayI64(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeI64(*leaf.sli, nn)
			r.ReadArrayI64(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafL) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]int64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafL) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]int64:
			return leaf.setAddress(sli)
		case *[]uint64:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int64:
		leaf.ptr = v
	case *[]int64:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]int64, 0)
		}
	case *uint64:
		leaf.ptr = (*int64)(unsafe.Pointer(v))
	case *[]uint64:
		leaf.sli = (*[]int64)(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]int64, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafL) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteI64(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayI64((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayI64((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafL{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafL", f)
}

var (
	_ root.Object        = (*LeafL)(nil)
	_ root.Named         = (*LeafL)(nil)
	_ Leaf               = (*LeafL)(nil)
	_ rbytes.Marshaler   = (*LeafL)(nil)
	_ rbytes.Unmarshaler = (*LeafL)(nil)
)

// LeafG implements ROOT TLeafG
type LeafG struct {
	rvers int16
	tleaf
	ptr *int64
	sli *[]int64
	min int64
	max int64
}

func newLeafG(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafG {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafG{
		rvers: rvers.LeafG,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafG) Class() string {
	return "TLeafG"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafG) Minimum() int64 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafG) Maximum() int64 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafG) Kind() reflect.Kind {
	if leaf.IsUnsigned() {
		return reflect.Uint64
	}
	return reflect.Int64
}

// Type returns the leaf's type.
func (leaf *LeafG) Type() reflect.Type {
	if leaf.IsUnsigned() {
		var v uint64
		return reflect.TypeOf(v)
	}
	var v int64
	return reflect.TypeOf(v)
}

// ivalue returns the first leaf value as int
func (leaf *LeafG) ivalue() int {
	return int(*leaf.ptr)
}

// imax returns the leaf maximum value as int
func (leaf *LeafG) imax() int {
	return int(leaf.max)
}

func (leaf *LeafG) TypeName() string {
	if leaf.IsUnsigned() {
		return "uint64"
	}
	return "int64"
}

func (leaf *LeafG) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI64(leaf.min)
	w.WriteI64(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafG) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafG {
		panic(fmt.Errorf("rtree: invalid TLeafG version=%d > %d", hdr.Vers, rvers.LeafG))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI64()
	leaf.max = r.ReadI64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafG) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadI64()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeI64(*leaf.sli, nn)
			r.ReadArrayI64(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeI64(*leaf.sli, nn)
			r.ReadArrayI64(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafG) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]int64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafG) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]int64:
			return leaf.setAddress(sli)
		case *[]uint64:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int64:
		leaf.ptr = v
	case *[]int64:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]int64, 0)
		}
	case *uint64:
		leaf.ptr = (*int64)(unsafe.Pointer(v))
	case *[]uint64:
		leaf.sli = (*[]int64)(unsafe.Pointer(v))
		if *v == nil {
			*leaf.sli = make([]int64, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafG) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteI64(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayI64((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayI64((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafG{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafG", f)
}

var (
	_ root.Object        = (*LeafG)(nil)
	_ root.Named         = (*LeafG)(nil)
	_ Leaf               = (*LeafG)(nil)
	_ rbytes.Marshaler   = (*LeafG)(nil)
	_ rbytes.Unmarshaler = (*LeafG)(nil)
)

// LeafF implements ROOT TLeafF
type LeafF struct {
	rvers int16
	tleaf
	ptr *float32
	sli *[]float32
	min float32
	max float32
}

func newLeafF(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafF {
	const etype = 4
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafF{
		rvers: rvers.LeafF,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafF) Class() string {
	return "TLeafF"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafF) Minimum() float32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafF) Maximum() float32 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafF) Kind() reflect.Kind {
	return reflect.Float32
}

// Type returns the leaf's type.
func (leaf *LeafF) Type() reflect.Type {
	var v float32
	return reflect.TypeOf(v)
}

func (leaf *LeafF) TypeName() string {
	return "float32"
}

func (leaf *LeafF) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteF32(leaf.min)
	w.WriteF32(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafF) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafF {
		panic(fmt.Errorf("rtree: invalid TLeafF version=%d > %d", hdr.Vers, rvers.LeafF))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadF32()
	leaf.max = r.ReadF32()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafF) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadF32()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeF32(*leaf.sli, nn)
			r.ReadArrayF32(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeF32(*leaf.sli, nn)
			r.ReadArrayF32(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafF) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]float32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafF) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]float32:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *float32:
		leaf.ptr = v
	case *[]float32:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]float32, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafF) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteF32(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayF32((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayF32((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafF{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafF", f)
}

var (
	_ root.Object        = (*LeafF)(nil)
	_ root.Named         = (*LeafF)(nil)
	_ Leaf               = (*LeafF)(nil)
	_ rbytes.Marshaler   = (*LeafF)(nil)
	_ rbytes.Unmarshaler = (*LeafF)(nil)
)

// LeafD implements ROOT TLeafD
type LeafD struct {
	rvers int16
	tleaf
	ptr *float64
	sli *[]float64
	min float64
	max float64
}

func newLeafD(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafD {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafD{
		rvers: rvers.LeafD,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafD) Class() string {
	return "TLeafD"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafD) Minimum() float64 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafD) Maximum() float64 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafD) Kind() reflect.Kind {
	return reflect.Float64
}

// Type returns the leaf's type.
func (leaf *LeafD) Type() reflect.Type {
	var v float64
	return reflect.TypeOf(v)
}

func (leaf *LeafD) TypeName() string {
	return "float64"
}

func (leaf *LeafD) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteF64(leaf.min)
	w.WriteF64(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafD) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafD {
		panic(fmt.Errorf("rtree: invalid TLeafD version=%d > %d", hdr.Vers, rvers.LeafD))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadF64()
	leaf.max = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafD) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadF64()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeF64(*leaf.sli, nn)
			r.ReadArrayF64(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeF64(*leaf.sli, nn)
			r.ReadArrayF64(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafD) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]float64)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafD) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]float64:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *float64:
		leaf.ptr = v
	case *[]float64:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]float64, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafD) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteF64(*leaf.ptr)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayF64((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayF64((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafD{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafD", f)
}

var (
	_ root.Object        = (*LeafD)(nil)
	_ root.Named         = (*LeafD)(nil)
	_ Leaf               = (*LeafD)(nil)
	_ rbytes.Marshaler   = (*LeafD)(nil)
	_ rbytes.Unmarshaler = (*LeafD)(nil)
)

// LeafF16 implements ROOT TLeafF16
type LeafF16 struct {
	rvers int16
	tleaf
	ptr *root.Float16
	sli *[]root.Float16
	min root.Float16
	max root.Float16
	elm rbytes.StreamerElement
}

func newLeafF16(b Branch, name string, shape []int, unsigned bool, count Leaf, elm rbytes.StreamerElement) *LeafF16 {
	const etype = 4
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafF16{
		rvers: rvers.LeafF16,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
		elm:   elm,
	}
}

// Class returns the ROOT class name.
func (leaf *LeafF16) Class() string {
	return "TLeafF16"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafF16) Minimum() root.Float16 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafF16) Maximum() root.Float16 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafF16) Kind() reflect.Kind {
	return reflect.Float32
}

// Type returns the leaf's type.
func (leaf *LeafF16) Type() reflect.Type {
	var v root.Float16
	return reflect.TypeOf(v)
}

func (leaf *LeafF16) TypeName() string {
	return "root.Float16"
}

func (leaf *LeafF16) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteF16(leaf.min, leaf.elm)
	w.WriteF16(leaf.max, leaf.elm)

	return w.SetHeader(hdr)
}

func (leaf *LeafF16) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafF16 {
		panic(fmt.Errorf("rtree: invalid TLeafF16 version=%d > %d", hdr.Vers, rvers.LeafF16))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadF16(leaf.elm)
	leaf.max = r.ReadF16(leaf.elm)

	if strings.Contains(leaf.Title(), "[") {
		elm := rdict.Element{
			Name:   *rbase.NewNamed(fmt.Sprintf("%s_Element", leaf.Name()), leaf.Title()),
			Offset: 0,
			Type:   rmeta.Float16,
		}.New()
		leaf.elm = &elm
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafF16) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadF16(leaf.elm)
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeF16(*leaf.sli, nn)
			r.ReadArrayF16(*leaf.sli, leaf.elm)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeF16(*leaf.sli, nn)
			r.ReadArrayF16(*leaf.sli, leaf.elm)
		}
	}
	return r.Err()
}

func (leaf *LeafF16) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 4
	arr := (*[0]root.Float16)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafF16) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]root.Float16:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *root.Float16:
		leaf.ptr = v
	case *[]root.Float16:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]root.Float16, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafF16) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteF16(*leaf.ptr, leaf.elm)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayF16((*leaf.sli)[:end], leaf.elm)
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayF16((*leaf.sli)[:leaf.tleaf.len], leaf.elm)
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafF16{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafF16", f)
}

var (
	_ root.Object        = (*LeafF16)(nil)
	_ root.Named         = (*LeafF16)(nil)
	_ Leaf               = (*LeafF16)(nil)
	_ rbytes.Marshaler   = (*LeafF16)(nil)
	_ rbytes.Unmarshaler = (*LeafF16)(nil)
)

// LeafD32 implements ROOT TLeafD32
type LeafD32 struct {
	rvers int16
	tleaf
	ptr *root.Double32
	sli *[]root.Double32
	min root.Double32
	max root.Double32
	elm rbytes.StreamerElement
}

func newLeafD32(b Branch, name string, shape []int, unsigned bool, count Leaf, elm rbytes.StreamerElement) *LeafD32 {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafD32{
		rvers: rvers.LeafD32,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
		elm:   elm,
	}
}

// Class returns the ROOT class name.
func (leaf *LeafD32) Class() string {
	return "TLeafD32"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafD32) Minimum() root.Double32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafD32) Maximum() root.Double32 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafD32) Kind() reflect.Kind {
	return reflect.Float64
}

// Type returns the leaf's type.
func (leaf *LeafD32) Type() reflect.Type {
	var v root.Double32
	return reflect.TypeOf(v)
}

func (leaf *LeafD32) TypeName() string {
	return "root.Double32"
}

func (leaf *LeafD32) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteD32(leaf.min, leaf.elm)
	w.WriteD32(leaf.max, leaf.elm)

	return w.SetHeader(hdr)
}

func (leaf *LeafD32) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafD32 {
		panic(fmt.Errorf("rtree: invalid TLeafD32 version=%d > %d", hdr.Vers, rvers.LeafD32))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadD32(leaf.elm)
	leaf.max = r.ReadD32(leaf.elm)

	if strings.Contains(leaf.Title(), "[") {
		elm := rdict.Element{
			Name:   *rbase.NewNamed(fmt.Sprintf("%s_Element", leaf.Name()), leaf.Title()),
			Offset: 0,
			Type:   rmeta.Double32,
		}.New()
		leaf.elm = &elm
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafD32) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadD32(leaf.elm)
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeD32(*leaf.sli, nn)
			r.ReadArrayD32(*leaf.sli, leaf.elm)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeD32(*leaf.sli, nn)
			r.ReadArrayD32(*leaf.sli, leaf.elm)
		}
	}
	return r.Err()
}

func (leaf *LeafD32) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 8
	arr := (*[0]root.Double32)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafD32) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]root.Double32:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *root.Double32:
		leaf.ptr = v
	case *[]root.Double32:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]root.Double32, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafD32) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteD32(*leaf.ptr, leaf.elm)
		nbytes += leaf.tleaf.etype
		if v := *leaf.ptr; v > leaf.max {
			leaf.max = v
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayD32((*leaf.sli)[:end], leaf.elm)
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayD32((*leaf.sli)[:leaf.tleaf.len], leaf.elm)
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafD32{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafD32", f)
}

var (
	_ root.Object        = (*LeafD32)(nil)
	_ root.Named         = (*LeafD32)(nil)
	_ Leaf               = (*LeafD32)(nil)
	_ rbytes.Marshaler   = (*LeafD32)(nil)
	_ rbytes.Unmarshaler = (*LeafD32)(nil)
)

// LeafC implements ROOT TLeafC
type LeafC struct {
	rvers int16
	tleaf
	ptr *string
	sli *[]string
	min int32
	max int32
}

func newLeafC(b Branch, name string, shape []int, unsigned bool, count Leaf) *LeafC {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafC{
		rvers: rvers.LeafC,
		tleaf: newLeaf(name, shape, etype, 0, false, unsigned, lcnt, b),
	}
}

// Class returns the ROOT class name.
func (leaf *LeafC) Class() string {
	return "TLeafC"
}

// Minimum returns the minimum value of the leaf.
func (leaf *LeafC) Minimum() int32 {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *LeafC) Maximum() int32 {
	return leaf.max
}

// Kind returns the leaf's kind.
func (leaf *LeafC) Kind() reflect.Kind {
	return reflect.String
}

// Type returns the leaf's type.
func (leaf *LeafC) Type() reflect.Type {
	var v string
	return reflect.TypeOf(v)
}

func (leaf *LeafC) TypeName() string {
	return "string"
}

func (leaf *LeafC) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(leaf.Class(), leaf.rvers)
	w.WriteObject(&leaf.tleaf)
	w.WriteI32(leaf.min)
	w.WriteI32(leaf.max)

	return w.SetHeader(hdr)
}

func (leaf *LeafC) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(leaf.Class())
	if hdr.Vers > rvers.LeafC {
		panic(fmt.Errorf("rtree: invalid TLeafC version=%d > %d", hdr.Vers, rvers.LeafC))
	}

	leaf.rvers = hdr.Vers

	r.ReadObject(&leaf.tleaf)

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckHeader(hdr)
	return r.Err()
}

func (leaf *LeafC) readFromBuffer(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	if leaf.count == nil && leaf.ptr != nil {
		*leaf.ptr = r.ReadString()
	} else {
		if leaf.count != nil {
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			nn := leaf.tleaf.len * n
			*leaf.sli = rbytes.ResizeStr(*leaf.sli, nn)
			r.ReadArrayString(*leaf.sli)
		} else {
			nn := leaf.tleaf.len
			*leaf.sli = rbytes.ResizeStr(*leaf.sli, nn)
			r.ReadArrayString(*leaf.sli)
		}
	}
	return r.Err()
}

func (leaf *LeafC) unsafeDecayArray(ptr interface{}) interface{} {
	rv := reflect.ValueOf(ptr).Elem()
	sz := rv.Type().Size() / 16
	arr := (*[0]string)(unsafe.Pointer(rv.UnsafeAddr()))
	sli := (*arr)[:]
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&sli))
	hdr.Len = int(sz)
	hdr.Cap = int(sz)
	return &sli
}

func (leaf *LeafC) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		sli := leaf.unsafeDecayArray(ptr)
		switch sli := sli.(type) {
		case *[]string:
			return leaf.setAddress(sli)
		default:
			panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *string:
		leaf.ptr = v
	case *[]string:
		leaf.sli = v
		if *v == nil {
			*leaf.sli = make([]string, 0)
		}
	default:
		panic(fmt.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafC) writeToBuffer(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	var nbytes int
	switch {
	case leaf.ptr != nil:
		w.WriteString(*leaf.ptr)
		sz := len(*leaf.ptr)
		nbytes += sz
		if v := int32(sz); v >= leaf.max {
			leaf.max = v + 1
		}
		if sz >= leaf.tleaf.len {
			leaf.tleaf.len = sz + 1
		}
	case leaf.count != nil:
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		end := leaf.tleaf.len * n
		w.WriteArrayString((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteArrayString((*leaf.sli)[:leaf.tleaf.len])
		nbytes += leaf.tleaf.etype * leaf.tleaf.len
	}

	return nbytes, w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &LeafC{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TLeafC", f)
}

var (
	_ root.Object        = (*LeafC)(nil)
	_ root.Named         = (*LeafC)(nil)
	_ Leaf               = (*LeafC)(nil)
	_ rbytes.Marshaler   = (*LeafC)(nil)
	_ rbytes.Unmarshaler = (*LeafC)(nil)
)
