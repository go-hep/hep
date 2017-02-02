// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type tleaf struct {
	named    tnamed
	len      int
	etype    int
	offset   int
	hasrange bool
	unsigned bool
	count    Leaf
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

func (leaf *tleaf) SetBranch(b Branch) {
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

func (leaf *tleaf) Value(int) interface{} {
	panic("not implemented")
}

func (leaf *tleaf) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("tleaf: %v %v %v\n", vers, pos, bcnt)

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
		leaf.count = ptr.(Leaf)
	}

	r.CheckByteCount(pos, bcnt, start, "TLeaf")
	if leaf.len == 0 {
		leaf.len = 1
	}

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tleaf{}
		return reflect.ValueOf(o)
	}
	Factory.add("TLeaf", f)
	Factory.add("*rootio.tleaf", f)
}

var _ Object = (*tleaf)(nil)
var _ Named = (*tleaf)(nil)
var _ Leaf = (*tleaf)(nil)
var _ ROOTUnmarshaler = (*tleaf)(nil)
