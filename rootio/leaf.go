// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

type tleaf struct {
	named    tnamed
	len      int
	etype    int
	offset   int
	hasrange bool
	unsigned bool
	count    leafCount
	branch   Branch
}

// Name returns the name of the instance
func (leaf *tleaf) Name() string {
	return leaf.named.Name()
}

// Title returns the title of the instance
func (leaf *tleaf) Title() string {
	return leaf.named.Title()
}

func (leaf *tleaf) Class() string {
	return "TLeaf"
}

func (leaf *tleaf) ArrayDim() int {
	panic("not implemented")
}

func (leaf *tleaf) setBranch(b Branch) {
	leaf.branch = b
}

func (leaf *tleaf) Branch() Branch {
	return leaf.branch
}

func (leaf *tleaf) HasRange() bool {
	return leaf.hasrange
}

func (leaf *tleaf) IsUnsigned() bool {
	return leaf.unsigned
}

func (leaf *tleaf) LeafCount() Leaf {
	return leaf.count
}

func (leaf *tleaf) Len() int {
	if leaf.count != nil {
		// variable length array
		n := leaf.count.ivalue()
		max := leaf.count.imax()
		if n > max {
			n = max
		}
		return leaf.len * n
	}
	return leaf.len
}

func (leaf *tleaf) LenType() int {
	return leaf.etype
}

func (leaf *tleaf) MaxIndex() []int {
	panic("not implemented")
}

func (leaf *tleaf) Offset() int {
	return leaf.offset
}

func (leaf *tleaf) Value(i int) interface{} {
	panic("not implemented")
}

func (leaf *tleaf) value() interface{} {
	panic("not implemented")
}

func (leaf *tleaf) readBasket(r *RBuffer) error {
	panic("not implemented")
}

func (leaf *tleaf) scan(r *RBuffer, ptr interface{}) error {
	panic("not implemented")
}

func (leaf *tleaf) TypeName() string {
	panic("not implemented")
}

func (leaf *tleaf) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := leaf.named.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.len = int(r.ReadI32())
	leaf.etype = int(r.ReadI32())
	leaf.offset = int(r.ReadI32())
	leaf.hasrange = r.ReadBool()
	leaf.unsigned = r.ReadBool()

	leaf.count = nil
	ptr := r.ReadObjectAny()
	if ptr != nil {
		leaf.count = ptr.(leafCount)
	}

	r.CheckByteCount(pos, bcnt, start, "TLeaf")
	if leaf.len == 0 {
		leaf.len = 1
	}

	return r.Err()
}

// tleafElement is a Leaf for a general object derived from Object.
type tleafElement struct {
	tleaf
	id    int32 // element serial number in fInfo
	ltype int32 // leaf type
}

func (leaf *tleafElement) Class() string {
	return "TLeafElement"
}

func (leaf *tleafElement) ivalue() int {
	panic("not implemented")
}

func (leaf *tleafElement) imax() int {
	panic("not implemented")
}

func (leaf *tleafElement) TypeName() string {
	panic("not implemented")
}

func (leaf *tleafElement) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}
	beg := r.Pos()

	_ /*vers*/, pos, bcnt := r.ReadVersion()

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.id = r.ReadI32()
	leaf.ltype = r.ReadI32()

	r.CheckByteCount(pos, bcnt, beg, "TLeafElement")
	return r.err
}

func (leaf *tleafElement) readBasket(r *RBuffer) error {
	panic("not implemented")
}

func (leaf *tleafElement) scan(r *RBuffer, ptr interface{}) error {
	panic("not implemented")
}

func init() {
	{
		f := func() reflect.Value {
			o := &tleaf{}
			return reflect.ValueOf(o)
		}
		Factory.add("TLeaf", f)
		Factory.add("*rootio.tleaf", f)
	}
	{
		f := func() reflect.Value {
			o := &tleafElement{}
			return reflect.ValueOf(o)
		}
		Factory.add("TLeafElement", f)
		Factory.add("*rootio.tleafElement", f)
	}
}

var _ Object = (*tleaf)(nil)
var _ Named = (*tleaf)(nil)
var _ Leaf = (*tleaf)(nil)
var _ ROOTUnmarshaler = (*tleaf)(nil)

var _ Object = (*tleafElement)(nil)
var _ Named = (*tleafElement)(nil)
var _ Leaf = (*tleafElement)(nil)
var _ ROOTUnmarshaler = (*tleafElement)(nil)
