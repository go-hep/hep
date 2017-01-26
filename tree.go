// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
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
func (tree *Tree) UnmarshalROOT(data *bytes.Buffer) error {
	dec := newDecoder(data)

	myprintf(">>>>>>>>>>>>>> Tree.unmarshal...\n")
	vers, pos, bcnt := dec.readVersion()
	myprintf(">>> => [%v] [%v] [%v]\n", pos, vers, bcnt)

	for _, a := range []ROOTUnmarshaler{
		&tree.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		err := a.UnmarshalROOT(data)
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

	//fmt.Printf("### data = %v\n", dec.data.Bytes()[:64])
	dec.readBin(&tree.entries)
	dec.readBin(&tree.totbytes)
	dec.readBin(&tree.zipbytes)
	if vers >= 18 {
		var flushedbytes int64
		dec.readInt64(&flushedbytes)
	}

	// dummy values
	var (
		f64 float64
		i32 int32
		i64 int64
	)

	dec.readBin(&f64)   // fWeight
	dec.readInt32(&i32) // fTimerInterval
	dec.readInt32(&i32) // fScanField
	dec.readInt32(&i32) // fUpdate

	if vers >= 18 {
		dec.readInt32(&i32) // fDefaultEntryOffsetLen
	}

	dec.readInt64(&i64) // fMaxEntries
	dec.readInt64(&i64) // fMaxEntryLoop
	dec.readInt64(&i64) // fMaxVirtualSize
	dec.readInt64(&i64) // fAutoSave

	if vers >= 18 {
		dec.readInt64(&i64) // fAutoFlush
	}

	dec.readInt64(&i64) // fEstimate

	return dec.err
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
