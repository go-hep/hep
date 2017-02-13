// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"fmt"
	"reflect"
)

type tbranch struct {
	named          tnamed
	attfill        attfill
	compress       int      // compression level and algorithm
	basketSize     int      // initial size of Basket buffer
	entryOffsetLen int      // initial length of entryOffset table in the basket buffers
	writeBasket    int      // last basket number written
	entryNumber    int64    // current entry number (last one filled in this branch)
	offset         int      // offset of this branch
	maxBaskets     int      // maximum number of baskets so far
	splitLevel     int      // branch split level
	entries        int64    // number of entries
	firstEntry     int64    // number of the first entry in this branch
	totBytes       int64    // total number of bytes in all leaves before compression
	zipBytes       int64    // total number of bytes in all leaves after compression
	branches       []Branch // list of branches of this branch
	leaves         []Leaf   // list of leaves of this branch
	baskets        []Basket // list of baskets of this branch

	basketBytes []int32 // length of baskets on file
	basketEntry []int64 // table of first entry in each basket
	basketSeek  []int64 // addresses of baskets on file

	fname string // named of file where buffers are stored (empty if in same file as Tree header)

	readbasket  int     // current basket number when reading
	readentry   int64   // current entry number when reading
	firstbasket int64   // first entry in the current basket
	nextbasket  int64   // next entry that will reaquire us to go to the next basket
	basket      *Basket // pointer to the current basket

	tree   Tree        // tree header
	mother *tbranch    // top-level parent branch in the tree
	parent *tbranch    // parent branch
	dir    *tdirectory // directory where this branch's buffers are stored
}

func (b *tbranch) Name() string {
	return b.named.Name()
}

func (b *tbranch) Title() string {
	return b.named.Title()
}

func (b *tbranch) Class() string {
	return "TBranch"
}

func (b *tbranch) SetTree(t Tree) {
	b.tree = t
}

func (b *tbranch) Branches() []Branch {
	return b.branches
}

func (b *tbranch) Leaves() []Leaf {
	return b.leaves
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranch) UnmarshalROOT(r *RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()

	b.tree = nil
	b.basket = nil
	b.firstbasket = -1
	b.nextbasket = -1

	if vers < 12 {
		panic(fmt.Errorf("rootio: too old TBanch version (%d<12)", vers))
	}

	if err := b.named.UnmarshalROOT(r); err != nil {
		return err
	}

	if err := b.attfill.UnmarshalROOT(r); err != nil {
		return err
	}

	b.compress = int(r.ReadI32())
	b.basketSize = int(r.ReadI32())
	b.entryOffsetLen = int(r.ReadI32())
	b.writeBasket = int(r.ReadI32())
	b.entryNumber = r.ReadI64()
	b.offset = int(r.ReadI32())
	b.maxBaskets = int(r.ReadI32())
	b.splitLevel = int(r.ReadI32())
	b.entries = r.ReadI64()
	b.firstEntry = r.ReadI64()
	b.totBytes = r.ReadI64()
	b.zipBytes = r.ReadI64()

	{
		var branches objarray
		if err := branches.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.branches = make([]Branch, branches.last+1)
		for i := range b.branches {
			br := branches.At(i).(Branch)
			b.branches[i] = br
		}
	}

	{
		var leaves objarray
		if err := leaves.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.leaves = make([]Leaf, leaves.last+1)
		for i := range b.leaves {
			leaf := leaves.At(i).(Leaf)
			b.leaves[i] = leaf
		}
	}
	{
		var baskets objarray
		if err := baskets.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
		b.baskets = make([]Basket, baskets.last+1)
		for i := range b.baskets {
			bk := baskets.At(i)
			// FIXME(sbinet) check why some are nil
			if bk != nil {
				b.baskets[i] = *(bk.(*Basket))
			}
		}
	}

	b.basketBytes = nil
	b.basketEntry = nil
	b.basketSeek = nil

	/*isArray*/ _ = r.ReadI8()
	b.basketBytes = r.ReadFastArrayI32(b.maxBaskets)

	/*isArray*/ _ = r.ReadI8()
	b.basketEntry = r.ReadFastArrayI64(b.maxBaskets)

	/*isArray*/ _ = r.ReadI8()
	b.basketSeek = r.ReadFastArrayI64(b.maxBaskets)

	b.fname = r.ReadString()

	r.CheckByteCount(pos, bcnt, beg, "TBranch")

	if b.splitLevel == 0 && len(b.branches) > 0 {
		b.splitLevel = 1
	}

	return r.Err()
}

func (b *tbranch) loadEntry(ientry int64) error {
	var err error
	err = b.loadBasket(ientry)
	if err != nil {
		return err
	}

	err = b.basket.loadEntry(ientry)
	if err != nil {
		return err
	}

	return err
}

func (b *tbranch) loadBasket(entry int64) error {
	var err error

	ib := b.findBasketIndex(entry)
	if ib < 0 {
		return errorf("rootio: no basket for entry %d", entry)
	}
	if ib < len(b.baskets) {
		b.basket = &b.baskets[ib]
		return nil
	}

	buf := make([]byte, int(b.basketBytes[ib]))
	// FIXME(sbinet)
	f := b.tree.(*ttree).f
	_, err = f.ReadAt(buf, b.basketSeek[ib])
	if err != nil {
		return err
	}

	b.baskets = append(b.baskets, Basket{})
	b.basket = &b.baskets[len(b.baskets)-1]
	err = b.basket.UnmarshalROOT(NewRBuffer(buf, nil, 0))
	if err != nil {
		return err
	}
	b.basket.f = f

	buf, err = b.basket.Bytes()
	if err != nil {
		return err
	}
	b.basket.buf = bytes.NewReader(buf)

	return err
}

func (b *tbranch) findBasketIndex(entry int64) int {
	return 0
}

func (b *tbranch) scan(ptr interface{}) error {
	return b.basket.scan(ptr)
}

// tbranchElement is a Branch for objects.
type tbranchElement struct {
	tbranch
	class   string          // class name of referenced object
	parent  string          // name of parent class
	clones  string          // named of class in TClonesArray (if any)
	chksum  uint32          // checksum of class
	clsver  uint16          // version number of class
	id      int32           // element serial number in fInfo
	btype   int32           // branch type
	stype   int32           // branch streamer type
	max     int32           // maximum entries for a TClonesArray or variable array
	bcount1 *tbranchElement // pointer to primary branchcount branch
	bcount2 *tbranchElement // pointer to secondary branchcount branch
}

func (b *tbranchElement) Class() string {
	return "TBranchElement"
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *tbranchElement) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 9 {
		r.err = fmt.Errorf("rootio: TBranchElement version too old (%d < 9)", vers)
		return r.err
	}

	if err := b.tbranch.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	b.class = r.ReadString()
	b.parent = r.ReadString()
	b.clones = r.ReadString()
	b.chksum = r.ReadU32()
	if vers >= 10 {
		b.clsver = r.ReadU16()
	} else {
		b.clsver = uint16(r.ReadU32())
	}
	b.id = r.ReadI32()
	b.btype = r.ReadI32()
	b.stype = r.ReadI32()
	b.max = r.ReadI32()

	bcount1 := r.ReadObjectAny()
	if bcount1 != nil {
		b.bcount1 = bcount1.(*tbranchElement)
	}

	bcount2 := r.ReadObjectAny()
	if bcount2 != nil {
		b.bcount2 = bcount2.(*tbranchElement)
	}

	r.CheckByteCount(pos, bcnt, beg, "TBranchElement")
	return r.err
}

func init() {
	{
		f := func() reflect.Value {
			o := &tbranch{}
			return reflect.ValueOf(o)
		}
		Factory.add("TBranch", f)
		Factory.add("*rootio.tbranch", f)
	}
	{
		f := func() reflect.Value {
			o := &tbranchElement{}
			return reflect.ValueOf(o)
		}
		Factory.add("TBranchElement", f)
		Factory.add("*rootio.tbranchElement", f)
	}
}

var _ Object = (*tbranch)(nil)
var _ Named = (*tbranch)(nil)
var _ Branch = (*tbranch)(nil)
var _ ROOTUnmarshaler = (*tbranch)(nil)

var _ Object = (*tbranchElement)(nil)
var _ Named = (*tbranchElement)(nil)
var _ Branch = (*tbranchElement)(nil)
var _ ROOTUnmarshaler = (*tbranchElement)(nil)
