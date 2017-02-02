// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

// The TNamed class is the base class for all named ROOT classes
// A TNamed contains the essential elements (name, title)
// to identify a derived object in containers, directories and files.
// Most member functions defined in this base class are in general
// overridden by the derived classes.
type tnamed struct {
	name  string
	title string
}

// Name returns the name of the instance
func (n *tnamed) Name() string {
	return n.name
}

// Title returns the title of the instance
func (n *tnamed) Title() string {
	return n.title
}

func (n *tnamed) Class() string {
	return "TNamed"
}

func (n *tnamed) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("named: %v %v %v\n", vers, pos, bcnt)

	var (
		_    = r.ReadU32() // id
		bits = r.ReadU32() // bits
	)
	bits |= kIsOnHeap // by definition, de-serialized object is on heap
	if (bits & kIsReferenced) == 0 {
		_ = r.ReadU16()
	}

	n.name = r.ReadString()
	n.title = r.ReadString()

	r.CheckByteCount(pos, bcnt, start, "TNamed")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tnamed{}
		return reflect.ValueOf(o)
	}
	Factory.add("TNamed", f)
	Factory.add("*rootio.tnamed", f)
}

var _ Object = (*tnamed)(nil)
var _ Named = (*tnamed)(nil)
var _ ROOTUnmarshaler = (*tnamed)(nil)
