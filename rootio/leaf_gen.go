// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
	"unsafe"
)

// LeafO implements ROOT TLeafO
type LeafO struct {
	tleaf
	ptr *bool
	arr *[]bool
	//val []bool
	min bool
	max bool
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
func (*LeafO) Kind() reflect.Kind {
	return reflect.Bool
}

// Type returns the leaf's type.
func (*LeafO) Type() reflect.Type {
	var v bool
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafO) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafO) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

func (leaf *LeafO) TypeName() string {
	return "bool"
}

func (leaf *LeafO) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadBool(&leaf.min)
	r.ReadBool(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafO")
	return r.Err()
}

func (leaf *LeafO) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadBool(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayBool(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayBool(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafO) resize(n int) {
	if leaf.arr == nil {
		arr := make([]bool, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]bool, n)
}

func (leaf *LeafO) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]bool:
		leaf.arr = v
	case *bool:
		leaf.ptr = v
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafO) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *bool:
		*v = *leaf.ptr
	case *[]bool:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]bool, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []bool:
		copy(v, *leaf.arr)

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafO{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafO", f)
	Factory.add("*rootio.LeafO", f)
}

var _ Object = (*LeafO)(nil)
var _ Named = (*LeafO)(nil)
var _ Leaf = (*LeafO)(nil)
var _ ROOTUnmarshaler = (*LeafO)(nil)

// LeafB implements ROOT TLeafB
type LeafB struct {
	tleaf
	ptr *int8
	arr *[]int8
	//val []int8
	min int8
	max int8
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
func (*LeafB) Kind() reflect.Kind {
	return reflect.Int8
}

// Type returns the leaf's type.
func (*LeafB) Type() reflect.Type {
	var v int8
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafB) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafB) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

// ivalue returns the first leaf value as int
func (leaf *LeafB) ivalue() int {
	if leaf.ptr != nil {
		return int(*leaf.ptr)
	}
	return int((*leaf.arr)[0])
}

// imax returns the leaf maximum value as int
func (leaf *LeafB) imax() int {
	return int(leaf.max)
}

func (leaf *LeafB) TypeName() string {
	return "int8"
}

func (leaf *LeafB) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadI8(&leaf.min)
	r.ReadI8(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafB")
	return r.Err()
}

func (leaf *LeafB) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadI8(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayI8(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayI8(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafB) resize(n int) {
	if leaf.arr == nil {
		arr := make([]int8, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]int8, n)
}

func (leaf *LeafB) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]int8:
		leaf.arr = v
	case *int8:
		leaf.ptr = v
	case *[]uint8:
		leaf.arr = (*[]int8)(unsafe.Pointer(v))
	case *uint8:
		leaf.ptr = (*int8)(unsafe.Pointer(v))
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafB) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int8:
		*v = *leaf.ptr
	case *[]int8:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]int8, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []int8:
		copy(v, *leaf.arr)

	case *uint8:
		*v = uint8(*leaf.ptr)
	case *[]uint8:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]uint8, len(*leaf.arr))
		}
		for i, u := range *leaf.arr {
			(*v)[i] = uint8(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint8:
		for i := range v {
			v[i] = uint8(((*leaf.arr)[i]))
		}

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafB{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafB", f)
	Factory.add("*rootio.LeafB", f)
}

var _ Object = (*LeafB)(nil)
var _ Named = (*LeafB)(nil)
var _ Leaf = (*LeafB)(nil)
var _ ROOTUnmarshaler = (*LeafB)(nil)

// LeafS implements ROOT TLeafS
type LeafS struct {
	tleaf
	ptr *int16
	arr *[]int16
	//val []int16
	min int16
	max int16
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
func (*LeafS) Kind() reflect.Kind {
	return reflect.Int16
}

// Type returns the leaf's type.
func (*LeafS) Type() reflect.Type {
	var v int16
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafS) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafS) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

// ivalue returns the first leaf value as int
func (leaf *LeafS) ivalue() int {
	if leaf.ptr != nil {
		return int(*leaf.ptr)
	}
	return int((*leaf.arr)[0])
}

// imax returns the leaf maximum value as int
func (leaf *LeafS) imax() int {
	return int(leaf.max)
}

func (leaf *LeafS) TypeName() string {
	return "int16"
}

func (leaf *LeafS) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadI16(&leaf.min)
	r.ReadI16(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafS")
	return r.Err()
}

func (leaf *LeafS) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadI16(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayI16(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayI16(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafS) resize(n int) {
	if leaf.arr == nil {
		arr := make([]int16, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]int16, n)
}

func (leaf *LeafS) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]int16:
		leaf.arr = v
	case *int16:
		leaf.ptr = v
	case *[]uint16:
		leaf.arr = (*[]int16)(unsafe.Pointer(v))
	case *uint16:
		leaf.ptr = (*int16)(unsafe.Pointer(v))
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafS) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int16:
		*v = *leaf.ptr
	case *[]int16:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]int16, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []int16:
		copy(v, *leaf.arr)

	case *uint16:
		*v = uint16(*leaf.ptr)
	case *[]uint16:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]uint16, len(*leaf.arr))
		}
		for i, u := range *leaf.arr {
			(*v)[i] = uint16(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint16:
		for i := range v {
			v[i] = uint16(((*leaf.arr)[i]))
		}

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafS{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafS", f)
	Factory.add("*rootio.LeafS", f)
}

var _ Object = (*LeafS)(nil)
var _ Named = (*LeafS)(nil)
var _ Leaf = (*LeafS)(nil)
var _ ROOTUnmarshaler = (*LeafS)(nil)

// LeafI implements ROOT TLeafI
type LeafI struct {
	tleaf
	ptr *int32
	arr *[]int32
	//val []int32
	min int32
	max int32
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
func (*LeafI) Kind() reflect.Kind {
	return reflect.Int32
}

// Type returns the leaf's type.
func (*LeafI) Type() reflect.Type {
	var v int32
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafI) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafI) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

// ivalue returns the first leaf value as int
func (leaf *LeafI) ivalue() int {
	if leaf.ptr != nil {
		return int(*leaf.ptr)
	}
	return int((*leaf.arr)[0])
}

// imax returns the leaf maximum value as int
func (leaf *LeafI) imax() int {
	return int(leaf.max)
}

func (leaf *LeafI) TypeName() string {
	return "int32"
}

func (leaf *LeafI) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadI32(&leaf.min)
	r.ReadI32(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafI")
	return r.Err()
}

func (leaf *LeafI) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadI32(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayI32(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayI32(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafI) resize(n int) {
	if leaf.arr == nil {
		arr := make([]int32, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]int32, n)
}

func (leaf *LeafI) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]int32:
		leaf.arr = v
	case *int32:
		leaf.ptr = v
	case *[]uint32:
		leaf.arr = (*[]int32)(unsafe.Pointer(v))
	case *uint32:
		leaf.ptr = (*int32)(unsafe.Pointer(v))
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafI) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int32:
		*v = *leaf.ptr
	case *[]int32:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]int32, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []int32:
		copy(v, *leaf.arr)

	case *uint32:
		*v = uint32(*leaf.ptr)
	case *[]uint32:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]uint32, len(*leaf.arr))
		}
		for i, u := range *leaf.arr {
			(*v)[i] = uint32(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint32:
		for i := range v {
			v[i] = uint32(((*leaf.arr)[i]))
		}

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafI{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafI", f)
	Factory.add("*rootio.LeafI", f)
}

var _ Object = (*LeafI)(nil)
var _ Named = (*LeafI)(nil)
var _ Leaf = (*LeafI)(nil)
var _ ROOTUnmarshaler = (*LeafI)(nil)

// LeafL implements ROOT TLeafL
type LeafL struct {
	tleaf
	ptr *int64
	arr *[]int64
	//val []int64
	min int64
	max int64
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
func (*LeafL) Kind() reflect.Kind {
	return reflect.Int64
}

// Type returns the leaf's type.
func (*LeafL) Type() reflect.Type {
	var v int64
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafL) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafL) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

// ivalue returns the first leaf value as int
func (leaf *LeafL) ivalue() int {
	if leaf.ptr != nil {
		return int(*leaf.ptr)
	}
	return int((*leaf.arr)[0])
}

// imax returns the leaf maximum value as int
func (leaf *LeafL) imax() int {
	return int(leaf.max)
}

func (leaf *LeafL) TypeName() string {
	return "int64"
}

func (leaf *LeafL) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadI64(&leaf.min)
	r.ReadI64(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafL")
	return r.Err()
}

func (leaf *LeafL) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadI64(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayI64(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayI64(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafL) resize(n int) {
	if leaf.arr == nil {
		arr := make([]int64, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]int64, n)
}

func (leaf *LeafL) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]int64:
		leaf.arr = v
	case *int64:
		leaf.ptr = v
	case *[]uint64:
		leaf.arr = (*[]int64)(unsafe.Pointer(v))
	case *uint64:
		leaf.ptr = (*int64)(unsafe.Pointer(v))
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafL) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int64:
		*v = *leaf.ptr
	case *[]int64:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]int64, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []int64:
		copy(v, *leaf.arr)

	case *uint64:
		*v = uint64(*leaf.ptr)
	case *[]uint64:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]uint64, len(*leaf.arr))
		}
		for i, u := range *leaf.arr {
			(*v)[i] = uint64(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint64:
		for i := range v {
			v[i] = uint64(((*leaf.arr)[i]))
		}

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafL{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafL", f)
	Factory.add("*rootio.LeafL", f)
}

var _ Object = (*LeafL)(nil)
var _ Named = (*LeafL)(nil)
var _ Leaf = (*LeafL)(nil)
var _ ROOTUnmarshaler = (*LeafL)(nil)

// LeafF implements ROOT TLeafF
type LeafF struct {
	tleaf
	ptr *float32
	arr *[]float32
	//val []float32
	min float32
	max float32
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
func (*LeafF) Kind() reflect.Kind {
	return reflect.Float32
}

// Type returns the leaf's type.
func (*LeafF) Type() reflect.Type {
	var v float32
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafF) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafF) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

func (leaf *LeafF) TypeName() string {
	return "float32"
}

func (leaf *LeafF) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadF32(&leaf.min)
	r.ReadF32(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafF")
	return r.Err()
}

func (leaf *LeafF) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadF32(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayF32(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayF32(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafF) resize(n int) {
	if leaf.arr == nil {
		arr := make([]float32, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]float32, n)
}

func (leaf *LeafF) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]float32:
		leaf.arr = v
	case *float32:
		leaf.ptr = v
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafF) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *float32:
		*v = *leaf.ptr
	case *[]float32:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]float32, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []float32:
		copy(v, *leaf.arr)

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafF{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafF", f)
	Factory.add("*rootio.LeafF", f)
}

var _ Object = (*LeafF)(nil)
var _ Named = (*LeafF)(nil)
var _ Leaf = (*LeafF)(nil)
var _ ROOTUnmarshaler = (*LeafF)(nil)

// LeafD implements ROOT TLeafD
type LeafD struct {
	tleaf
	ptr *float64
	arr *[]float64
	//val []float64
	min float64
	max float64
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
func (*LeafD) Kind() reflect.Kind {
	return reflect.Float64
}

// Type returns the leaf's type.
func (*LeafD) Type() reflect.Type {
	var v float64
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafD) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafD) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

func (leaf *LeafD) TypeName() string {
	return "float64"
}

func (leaf *LeafD) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadF64(&leaf.min)
	r.ReadF64(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafD")
	return r.Err()
}

func (leaf *LeafD) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadF64(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayF64(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayF64(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafD) resize(n int) {
	if leaf.arr == nil {
		arr := make([]float64, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]float64, n)
}

func (leaf *LeafD) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]float64:
		leaf.arr = v
	case *float64:
		leaf.ptr = v
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafD) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *float64:
		*v = *leaf.ptr
	case *[]float64:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]float64, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []float64:
		copy(v, *leaf.arr)

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafD{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafD", f)
	Factory.add("*rootio.LeafD", f)
}

var _ Object = (*LeafD)(nil)
var _ Named = (*LeafD)(nil)
var _ Leaf = (*LeafD)(nil)
var _ ROOTUnmarshaler = (*LeafD)(nil)

// LeafC implements ROOT TLeafC
type LeafC struct {
	tleaf
	ptr *string
	arr *[]string
	//val []string
	min int32
	max int32
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
func (*LeafC) Kind() reflect.Kind {
	return reflect.String
}

// Type returns the leaf's type.
func (*LeafC) Type() reflect.Type {
	var v string
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafC) Value(i int) interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return (*leaf.arr)[i]
}

// value returns the leaf value.
func (leaf *LeafC) value() interface{} {
	if leaf.ptr != nil {
		return *leaf.ptr
	}
	return *leaf.arr
}

func (leaf *LeafC) TypeName() string {
	return "string"
}

func (leaf *LeafC) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	_, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	r.ReadI32(&leaf.min)
	r.ReadI32(&leaf.max)

	r.CheckByteCount(pos, bcnt, start, "TLeafC")
	return r.Err()
}

func (leaf *LeafC) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && leaf.ptr != nil {
		r.ReadString(leaf.ptr)
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.resize(leaf.tleaf.len * n)
			r.ReadFastArrayString(*leaf.arr)
		} else {
			leaf.resize(leaf.tleaf.len)
			r.ReadFastArrayString(*leaf.arr)
		}
	}
	return r.err
}

func (leaf *LeafC) resize(n int) {
	if leaf.arr == nil {
		arr := make([]string, n)
		leaf.arr = &arr
		return
	}
	if len(*leaf.arr) > n {
		*leaf.arr = (*leaf.arr)[:n]
		return
	}
	*leaf.arr = make([]string, n)
}

func (leaf *LeafC) setAddress(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return errorf("rootio: invalid kind (%T). want a pointer", ptr)
	}
	if rv.Elem().Kind() == reflect.Array {
		panic("not implemented")
	}
	switch v := ptr.(type) {
	case *[]string:
		leaf.arr = v
	case *string:
		leaf.ptr = v
	default:
		panic(errorf("rootio: invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}
	return nil
}

func (leaf *LeafC) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *string:
		*v = *leaf.ptr
	case *[]string:
		if len(*v) < len(*leaf.arr) || *v == nil {
			*v = make([]string, len(*leaf.arr))
		}
		copy(*v, *leaf.arr)
		*v = (*v)[:leaf.count.ivalue()]
	case []string:
		copy(v, *leaf.arr)

	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &LeafC{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeafC", f)
	Factory.add("*rootio.LeafC", f)
}

var _ Object = (*LeafC)(nil)
var _ Named = (*LeafC)(nil)
var _ Leaf = (*LeafC)(nil)
var _ ROOTUnmarshaler = (*LeafC)(nil)
