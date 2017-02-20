// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

// A ttree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type ttree struct {
	f *File // underlying file

	named tnamed

	entries  int64 // Number of entries
	totbytes int64 // Total number of bytes in all branches before compression
	zipbytes int64 // Total number of bytes in all branches after  compression

	branches []Branch // list of branches
	leaves   []Leaf   // direct pointers to individual branch leaves
}

func (tree *ttree) Class() string {
	return "TTree"
}

func (tree *ttree) Name() string {
	return tree.named.Name()
}

func (tree *ttree) Title() string {
	return tree.named.Title()
}

func (tree *ttree) Entries() int64 {
	return tree.entries
}

func (tree *ttree) TotBytes() int64 {
	return tree.totbytes
}

func (tree *ttree) ZipBytes() int64 {
	return tree.zipbytes
}

func (tree *ttree) Branches() []Branch {
	return tree.branches
}

func (tree *ttree) Branch(name string) Branch {
	for _, br := range tree.branches {
		if br.Name() == name {
			return br
		}
	}
	return nil
}

func (tree *ttree) Leaves() []Leaf {
	return tree.leaves
}

func (tree *ttree) SetFile(f *File) {
	tree.f = f
}

func (tree *ttree) getFile() *File {
	return tree.f
}

func (tree *ttree) loadEntry(entry int64) error {
	for _, b := range tree.branches {
		err := b.loadEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (tree *ttree) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()

	for _, a := range []ROOTUnmarshaler{
		&tree.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		err := a.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}

	if vers < 16 {
		return fmt.Errorf(
			"rootio.Tree: tree [%s] with version [%v] is not supported (too old)",
			tree.Name(),
			vers,
		)
	}

	tree.entries = r.ReadI64()
	tree.totbytes = r.ReadI64()
	tree.zipbytes = r.ReadI64()
	if vers >= 19 { // FIXME
		_ = r.ReadI64() // fSavedBytes
	}
	if vers >= 18 {
		_ = r.ReadI64() // flushed bytes
	}

	_ = r.ReadF64() // fWeight
	_ = r.ReadI32() // fTimerInterval
	_ = r.ReadI32() // fScanField
	_ = r.ReadI32() // fUpdate

	if vers >= 18 {
		_ = r.ReadI32() // fDefaultEntryOffsetLen
	}
	nclus := 0
	if vers >= 19 { // FIXME
		nclus = int(r.ReadI32()) // fNClusterRange
	}

	_ = r.ReadI64() // fMaxEntries
	_ = r.ReadI64() // fMaxEntryLoop
	_ = r.ReadI64() // fMaxVirtualSize
	_ = r.ReadI64() // fAutoSave

	if vers >= 18 {
		_ = r.ReadI64() // fAutoFlush
	}

	_ = r.ReadI64() // fEstimate

	if vers >= 19 { // FIXME
		_ = r.ReadI8()
		_ = r.ReadFastArrayI64(nclus) // fClusterRangeEnd
		_ = r.ReadI8()
		_ = r.ReadFastArrayI64(nclus) // fClusterSize
	}

	var branches objarray
	if err := branches.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.branches = make([]Branch, branches.last+1)
	for i := range tree.branches {
		tree.branches[i] = branches.At(i).(Branch)
		tree.branches[i].setTree(tree)
	}

	var leaves objarray
	if err := leaves.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.leaves = make([]Leaf, leaves.last+1)
	for i := range tree.leaves {
		tree.leaves[i] = leaves.At(i).(Leaf)
		// FIXME(sbinet)
		//tree.leaves[i].SetBranch(tree.branches[i])
	}

	for _ = range []string{
		"fAliases", "fIndexValues", "fIndex", "fTreeIndex", "fFriends",
		"fUserInfo", "fBranchRef",
	} {
		_ = r.ReadObjectAny()
	}

	r.CheckByteCount(pos, bcnt, beg, "TTree")
	return r.Err()
}

type tntuple struct {
	ttree
	nvars int
}

func (nt *tntuple) Class() string {
	return "TNtuple"
}

func (nt *tntuple) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	if err := nt.ttree.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	nt.nvars = int(r.ReadI32())

	r.CheckByteCount(pos, bcnt, beg, "TNtuple")
	return r.err
}

func init() {
	{
		f := func() reflect.Value {
			o := &ttree{}
			return reflect.ValueOf(o)
		}
		Factory.add("TTree", f)
		Factory.add("*rootio.ttree", f)
	}
	{
		f := func() reflect.Value {
			o := &tntuple{}
			return reflect.ValueOf(o)
		}
		Factory.add("TNtuple", f)
		Factory.add("*rootio.tntuple", f)
	}
}

var _ Object = (*ttree)(nil)
var _ Named = (*ttree)(nil)
var _ Tree = (*ttree)(nil)
var _ ROOTUnmarshaler = (*ttree)(nil)

var _ Object = (*tntuple)(nil)
var _ Named = (*tntuple)(nil)
var _ Tree = (*tntuple)(nil)
var _ ROOTUnmarshaler = (*tntuple)(nil)
