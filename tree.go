// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"encoding/hex"
	"fmt"
	"reflect"
)

// A Tree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type Tree struct {
	f *File // underlying file

	named named

	entries  int64 // Number of entries
	totbytes int64 // Total number of bytes in all branches before compression
	zipbytes int64 // Total number of bytes in all branches after  compression

	branches []Object // list of branches
	leaves   []Object // direct pointers to individual branch leaves
}

func (tree *Tree) Class() string {
	return "TTree" //tree.classname
}

func (tree *Tree) Name() string {
	return tree.named.Name()
}

func (tree *Tree) Title() string {
	return tree.named.Title()
}

func (tree *Tree) Entries() int64 {
	return tree.entries
}

func (tree *Tree) TotBytes() int64 {
	return tree.totbytes
}

func (tree *Tree) ZipBytes() int64 {
	return tree.zipbytes
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (tree *Tree) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	myprintf(">>> => [%v] [%v] [%v]\n", pos, vers, bcnt)
	fmt.Printf("--- Tree vers=%d pos=%d count=%d\n", vers, pos, bcnt)
	{
		buf := r.bytes()
		if len(buf) > 256 {
			buf = buf[:256]
		}
		fmt.Printf("--- hex ---\n%s\n", string(hex.Dump(buf)))
	}

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

	// FIXME: hack. where do these 18 bytes come from ?
	// var trash [18]byte
	// err = dec.readBin(&trash)
	// if err != nil {
	// 	return err
	// }

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
		if nclus == 0 {
			_ = r.ReadI8() // fClusterRangeEnd
			_ = r.ReadI8() // fClusterSize
		} else {
			panic("not implemented")
		}
	}

	var branches objarray
	if err := branches.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.branches = branches.arr[:branches.last+1]

	var leaves objarray
	if err := leaves.UnmarshalROOT(r); err != nil {
		return err
	}
	tree.leaves = leaves.arr[:leaves.last+1]

	r.CheckByteCount(pos, bcnt, beg, "TTree")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &Tree{}
		return reflect.ValueOf(o)
	}
	Factory.add("TTree", f)
	Factory.add("*rootio.Tree", f)
}

var _ Object = (*Tree)(nil)
var _ Named = (*Tree)(nil)
var _ ROOTUnmarshaler = (*Tree)(nil)
