// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type tobjString struct {
	rvers int16
	obj   tobject
	str   string
}

func (*tobjString) Class() string {
	return "TObjString"
}

func (obj *tobjString) Name() string {
	return obj.str
}

func (obj *tobjString) Title() string {
	return "Collectable string class"
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (obj *tobjString) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	obj.rvers = vers
	r.SkipObject()
	obj.str = r.ReadString()

	r.CheckByteCount(pos, bcnt, start, "TObjString")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tobjString{}
		return reflect.ValueOf(o)
	}
	Factory.add("TObjString", f)
	Factory.add("*rootio.tobjString", f)
}

var _ Object = (*tobjString)(nil)
var _ Named = (*tobjString)(nil)
var _ ROOTUnmarshaler = (*tobjString)(nil)
