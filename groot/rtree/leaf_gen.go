// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rtree

import (
	"reflect"
	"unsafe"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
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

func newLeafO(b Branch, name string, len int, unsigned bool, count Leaf) *LeafO {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafO{
		rvers: rvers.LeafO,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafO) Kind() reflect.Kind {
	return reflect.Bool
}

// Type returns the leaf's type.
func (leaf *LeafO) Type() reflect.Type {
	var v bool
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafO) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafO) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
}

func (leaf *LeafO) TypeName() string {
	return "bool"
}

func (leaf *LeafO) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteBool(leaf.min)
	w.WriteBool(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafO) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadBool()
	leaf.max = r.ReadBool()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayBool(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayBool(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafO) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *bool:
		*v = *leaf.ptr
	case *[]bool:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]bool, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []bool:
		copy(v, *leaf.sli)
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafO) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []bool:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *bool:
		leaf.ptr = v
	case *[]bool:
		leaf.sli = v
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayBool((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayBool((*leaf.sli)[:leaf.tleaf.len])
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

func newLeafB(b Branch, name string, len int, unsigned bool, count Leaf) *LeafB {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafB{
		rvers: rvers.LeafB,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafB) Kind() reflect.Kind {
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

// Value returns the leaf value at index i.
func (leaf *LeafB) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafB) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
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

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteI8(leaf.min)
	w.WriteI8(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafB) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadI8()
	leaf.max = r.ReadI8()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayI8(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayI8(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafB) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int8:
		*v = *leaf.ptr
	case *[]int8:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]int8, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []int8:
		copy(v, *leaf.sli)
	case *uint8:
		*v = uint8(*leaf.ptr)
	case *[]uint8:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]uint8, len(*leaf.sli))
		}
		for i, u := range *leaf.sli {
			(*v)[i] = uint8(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint8:
		for i := range v {
			v[i] = uint8((*leaf.sli)[i])
		}
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafB) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []int8:
			return leaf.setAddress(&sli)
		case []uint8:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int8:
		leaf.ptr = v
	case *[]int8:
		leaf.sli = v
	case *uint8:
		leaf.ptr = (*int8)(unsafe.Pointer(v))
	case *[]uint8:
		leaf.sli = (*[]int8)(unsafe.Pointer(v))
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayI8((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayI8((*leaf.sli)[:leaf.tleaf.len])
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

func newLeafS(b Branch, name string, len int, unsigned bool, count Leaf) *LeafS {
	const etype = 2
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafS{
		rvers: rvers.LeafS,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafS) Kind() reflect.Kind {
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

// Value returns the leaf value at index i.
func (leaf *LeafS) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafS) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
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

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteI16(leaf.min)
	w.WriteI16(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafS) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadI16()
	leaf.max = r.ReadI16()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayI16(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayI16(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafS) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int16:
		*v = *leaf.ptr
	case *[]int16:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]int16, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []int16:
		copy(v, *leaf.sli)
	case *uint16:
		*v = uint16(*leaf.ptr)
	case *[]uint16:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]uint16, len(*leaf.sli))
		}
		for i, u := range *leaf.sli {
			(*v)[i] = uint16(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint16:
		for i := range v {
			v[i] = uint16((*leaf.sli)[i])
		}
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafS) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []int16:
			return leaf.setAddress(&sli)
		case []uint16:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int16:
		leaf.ptr = v
	case *[]int16:
		leaf.sli = v
	case *uint16:
		leaf.ptr = (*int16)(unsafe.Pointer(v))
	case *[]uint16:
		leaf.sli = (*[]int16)(unsafe.Pointer(v))
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayI16((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayI16((*leaf.sli)[:leaf.tleaf.len])
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

func newLeafI(b Branch, name string, len int, unsigned bool, count Leaf) *LeafI {
	const etype = 4
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafI{
		rvers: rvers.LeafI,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafI) Kind() reflect.Kind {
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

// Value returns the leaf value at index i.
func (leaf *LeafI) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafI) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
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

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteI32(leaf.min)
	w.WriteI32(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafI) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayI32(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayI32(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafI) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int32:
		*v = *leaf.ptr
	case *[]int32:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]int32, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []int32:
		copy(v, *leaf.sli)
	case *uint32:
		*v = uint32(*leaf.ptr)
	case *[]uint32:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]uint32, len(*leaf.sli))
		}
		for i, u := range *leaf.sli {
			(*v)[i] = uint32(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint32:
		for i := range v {
			v[i] = uint32((*leaf.sli)[i])
		}
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafI) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []int32:
			return leaf.setAddress(&sli)
		case []uint32:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int32:
		leaf.ptr = v
	case *[]int32:
		leaf.sli = v
	case *uint32:
		leaf.ptr = (*int32)(unsafe.Pointer(v))
	case *[]uint32:
		leaf.sli = (*[]int32)(unsafe.Pointer(v))
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayI32((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayI32((*leaf.sli)[:leaf.tleaf.len])
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

func newLeafL(b Branch, name string, len int, unsigned bool, count Leaf) *LeafL {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafL{
		rvers: rvers.LeafL,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafL) Kind() reflect.Kind {
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

// Value returns the leaf value at index i.
func (leaf *LeafL) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafL) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
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

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteI64(leaf.min)
	w.WriteI64(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafL) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadI64()
	leaf.max = r.ReadI64()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayI64(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayI64(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafL) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *int64:
		*v = *leaf.ptr
	case *[]int64:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]int64, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []int64:
		copy(v, *leaf.sli)
	case *uint64:
		*v = uint64(*leaf.ptr)
	case *[]uint64:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]uint64, len(*leaf.sli))
		}
		for i, u := range *leaf.sli {
			(*v)[i] = uint64(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []uint64:
		for i := range v {
			v[i] = uint64((*leaf.sli)[i])
		}
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafL) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []int64:
			return leaf.setAddress(&sli)
		case []uint64:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *int64:
		leaf.ptr = v
	case *[]int64:
		leaf.sli = v
	case *uint64:
		leaf.ptr = (*int64)(unsafe.Pointer(v))
	case *[]uint64:
		leaf.sli = (*[]int64)(unsafe.Pointer(v))
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayI64((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayI64((*leaf.sli)[:leaf.tleaf.len])
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

// LeafF implements ROOT TLeafF
type LeafF struct {
	rvers int16
	tleaf
	ptr *float32
	sli *[]float32
	min float32
	max float32
}

func newLeafF(b Branch, name string, len int, unsigned bool, count Leaf) *LeafF {
	const etype = 4
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafF{
		rvers: rvers.LeafF,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafF) Kind() reflect.Kind {
	return reflect.Float32
}

// Type returns the leaf's type.
func (leaf *LeafF) Type() reflect.Type {
	var v float32
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafF) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafF) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
}

func (leaf *LeafF) TypeName() string {
	return "float32"
}

func (leaf *LeafF) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteF32(leaf.min)
	w.WriteF32(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafF) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadF32()
	leaf.max = r.ReadF32()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayF32(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayF32(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafF) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *float32:
		*v = *leaf.ptr
	case *[]float32:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]float32, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []float32:
		copy(v, *leaf.sli)
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafF) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []float32:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *float32:
		leaf.ptr = v
	case *[]float32:
		leaf.sli = v
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayF32((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayF32((*leaf.sli)[:leaf.tleaf.len])
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

func newLeafD(b Branch, name string, len int, unsigned bool, count Leaf) *LeafD {
	const etype = 8
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafD{
		rvers: rvers.LeafD,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafD) Kind() reflect.Kind {
	return reflect.Float64
}

// Type returns the leaf's type.
func (leaf *LeafD) Type() reflect.Type {
	var v float64
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafD) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafD) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
}

func (leaf *LeafD) TypeName() string {
	return "float64"
}

func (leaf *LeafD) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteF64(leaf.min)
	w.WriteF64(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafD) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadF64()
	leaf.max = r.ReadF64()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayF64(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayF64(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafD) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *float64:
		*v = *leaf.ptr
	case *[]float64:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]float64, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []float64:
		copy(v, *leaf.sli)
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafD) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []float64:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *float64:
		leaf.ptr = v
	case *[]float64:
		leaf.sli = v
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayF64((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayF64((*leaf.sli)[:leaf.tleaf.len])
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

// LeafC implements ROOT TLeafC
type LeafC struct {
	rvers int16
	tleaf
	ptr *string
	sli *[]string
	min int32
	max int32
}

func newLeafC(b Branch, name string, len int, unsigned bool, count Leaf) *LeafC {
	const etype = 1
	var lcnt leafCount
	if count != nil {
		lcnt = count.(leafCount)
	}
	return &LeafC{
		rvers: rvers.LeafC,
		tleaf: newLeaf(name, len, etype, 0, false, unsigned, lcnt, b),
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
func (*LeafC) Kind() reflect.Kind {
	return reflect.String
}

// Type returns the leaf's type.
func (leaf *LeafC) Type() reflect.Type {
	var v string
	return reflect.TypeOf(v)
}

// Value returns the leaf value at index i.
func (leaf *LeafC) Value(i int) interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return (*leaf.sli)[i]
	}
}

// value returns the leaf value.
func (leaf *LeafC) value() interface{} {
	switch {
	case leaf.ptr != nil:
		return *leaf.ptr
	default:
		return *leaf.sli
	}
}

func (leaf *LeafC) TypeName() string {
	return "string"
}

func (leaf *LeafC) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(leaf.rvers)
	leaf.tleaf.MarshalROOT(w)
	w.WriteI32(leaf.min)
	w.WriteI32(leaf.max)

	return w.SetByteCount(pos, leaf.Class())
}

func (leaf *LeafC) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(leaf.Class())
	leaf.rvers = vers

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		return err
	}

	leaf.min = r.ReadI32()
	leaf.max = r.ReadI32()

	r.CheckByteCount(pos, bcnt, start, leaf.Class())
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
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			*leaf.sli = r.ReadFastArrayString(leaf.tleaf.len * n)
		} else {
			copy(*leaf.sli, r.ReadFastArrayString(leaf.tleaf.len))
		}
	}
	return r.Err()
}

func (leaf *LeafC) scan(r *rbytes.RBuffer, ptr interface{}) error {
	if r.Err() != nil {
		return r.Err()
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *string:
		*v = *leaf.ptr
	case *[]string:
		if len(*v) < len(*leaf.sli) || *v == nil {
			*v = make([]string, len(*leaf.sli))
		}
		copy(*v, *leaf.sli)
		*v = (*v)[:leaf.count.ivalue()]
	case []string:
		copy(v, *leaf.sli)
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.Err()
}

func (leaf *LeafC) setAddress(ptr interface{}) error {
	if ptr == nil {
		return leaf.setAddress(newValue(leaf))
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		arr := reflect.ValueOf(ptr).Elem()
		switch sli := arr.Slice(0, rv.Len()).Interface().(type) {
		case []string:
			return leaf.setAddress(&sli)
		default:
			panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", ptr, leaf.Name(), leaf))
		}
	}

	switch v := ptr.(type) {
	case *string:
		leaf.ptr = v
	case *[]string:
		leaf.sli = v
	default:
		panic(xerrors.Errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
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
		w.WriteFastArrayString((*leaf.sli)[:end])
		nbytes += leaf.tleaf.etype * end
	default:
		w.WriteFastArrayString((*leaf.sli)[:leaf.tleaf.len])
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
