// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type ObjString struct {
	obj Object
	str string
}

// NewObjString creates a new ObjString.
func NewObjString(s string) *ObjString {
	return &ObjString{
		obj: *NewObject(),
		str: s,
	}
}

func (*ObjString) RVersion() int16 {
	return rvers.ObjString
}

func (*ObjString) Class() string {
	return "TObjString"
}

func (obj *ObjString) Name() string {
	return obj.str
}

func (*ObjString) Title() string {
	return "Collectable string class"
}

func (obj *ObjString) String() string {
	return obj.str
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (obj *ObjString) UnmarshalROOT(r *rbytes.RBuffer) error {
	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()
	if err := obj.obj.UnmarshalROOT(r); err != nil {
		return err
	}
	obj.str = r.ReadString()

	r.CheckByteCount(pos, bcnt, start, "TObjString")
	return r.Err()
}

func (obj *ObjString) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}
	pos := w.WriteVersion(obj.RVersion())
	if _, err := obj.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteString(obj.str)

	return w.SetByteCount(pos, "TObjString")
}

func init() {
	f := func() reflect.Value {
		o := &ObjString{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TObjString", f)
}

var (
	_ root.Object        = (*ObjString)(nil)
	_ root.Named         = (*ObjString)(nil)
	_ root.ObjString     = (*ObjString)(nil)
	_ rbytes.Marshaler   = (*ObjString)(nil)
	_ rbytes.Unmarshaler = (*ObjString)(nil)
)
