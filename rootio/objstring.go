// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

type tobjstring struct {
	rvers int16
	obj   tobject
	str   string
}

func (*tobjstring) Class() string {
	return "TObjString"
}

func (obj *tobjstring) Name() string {
	return obj.str
}

func (obj *tobjstring) Title() string {
	return "Collectable string class"
}

func (obj *tobjstring) String() string {
	return obj.str
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (obj *tobjstring) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	obj.rvers = vers
	if err := obj.obj.UnmarshalROOT(r); err != nil {
		return err
	}
	obj.str = r.ReadString()

	r.CheckByteCount(pos, bcnt, start, "TObjString")
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &tobjstring{}
		return reflect.ValueOf(o)
	}
	Factory.add("TObjString", f)
	Factory.add("*rootio.tobjString", f)
}

var (
	_ Object          = (*tobjstring)(nil)
	_ Named           = (*tobjstring)(nil)
	_ ROOTUnmarshaler = (*tobjstring)(nil)
)
