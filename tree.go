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
	var err error
	dec := newDecoder(data)

	myprintf(">>>>>>>>>>>>>> Tree.unmarshal...\n")
	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	}
	myprintf(">>> => [%v] [%v] [%v]\n", pos, vers, bcnt)

	for _, a := range []ROOTUnmarshaler{
		&tree.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		err = a.UnmarshalROOT(data)
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
	err = dec.readBin(&tree.entries)
	if err != nil {
		return err
	}

	err = dec.readBin(&tree.totbytes)
	if err != nil {
		return err
	}

	err = dec.readBin(&tree.zipbytes)
	if err != nil {
		return err
	}

	if vers >= 18 {
		var flushedbytes int64
		err = dec.readInt64(&flushedbytes)
		if err != nil {
			return err
		}
	}
	var dummy_f64 float64
	var dummy_i32 int32
	var dummy_i64 int64

	err = dec.readBin(&dummy_f64) // fWeight
	if err != nil {
		return err
	}

	err = dec.readInt32(&dummy_i32) // fTimerInterval
	if err != nil {
		return err
	}
	err = dec.readInt32(&dummy_i32) // fScanField
	if err != nil {
		return err
	}
	err = dec.readInt32(&dummy_i32) // fUpdate
	if err != nil {
		return err
	}

	if vers >= 18 {
		err = dec.readInt32(&dummy_i32) // fDefaultEntryOffsetLen
		if err != nil {
			return err
		}
	}

	err = dec.readInt64(&dummy_i64) // fMaxEntries
	if err != nil {
		return err
	}

	err = dec.readInt64(&dummy_i64) // fMaxEntryLoop
	if err != nil {
		return err
	}

	err = dec.readInt64(&dummy_i64) // fMaxVirtualSize
	if err != nil {
		return err
	}

	err = dec.readInt64(&dummy_i64) // fAutoSave
	if err != nil {
		return err
	}

	if vers >= 18 {
		err = dec.readInt64(&dummy_i64) // fAutoFlush
		if err != nil {
			return err
		}
	}

	err = dec.readInt64(&dummy_i64) // fEstimate
	if err != nil {
		return err
	}

	return err
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
var _ ROOTUnmarshaler = (*Tree)(nil)
