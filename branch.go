// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"reflect"
)

type Branch struct {
	f *File

	named      named
	autodelete bool
	branches   []Branch // list of branches of this branch
	leaves     []Leaf   // list of leaves of this branch
	baskets    []Basket // list of baskets of this branch

	readbasket  int     // current basket number when reading
	readentry   int64   // current entry number when reading
	firstbasket int64   // first entry in the current basket
	nextbasket  int64   // next entry that will reaquire us to go to the next basket
	currbasket  *Basket // pointer to the current basket

	tree   *Tree      // tree header
	mother *Branch    // top-level parent branch in the tree
	parent *Branch    // parent branch
	dir    *directory // directory where this branch's buffers are stored
}

func (b *Branch) Name() string {
	return b.named.Name()
}

func (b *Branch) Title() string {
	return b.named.Title()
}

func (b *Branch) Class() string {
	return "TBranch"
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (b *Branch) UnmarshalROOT(data *bytes.Buffer) error {
	var err error
	panic("not implemented")
	return err
}

func init() {
	f := func() reflect.Value {
		o := &Branch{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBranch"] = f
	Factory.db["*rootio.Branch"] = f
}

var _ Object = (*Branch)(nil)
var _ ROOTUnmarshaler = (*Branch)(nil)
