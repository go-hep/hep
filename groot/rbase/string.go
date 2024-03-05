// Copyright ©2024 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

// String is a bare-bone string value implementing root.Object
type String struct {
	str string
}

func NewString(v string) *String {
	return &String{v}
}

func (*String) Class() string {
	return "*rbase.String"
}

func (*String) RVersion() int16 {
	return 0
}

func (v *String) String() string {
	return v.str
}

func (obj *String) UnmarshalROOT(r *rbytes.RBuffer) error {
	obj.str = r.ReadString()
	return r.Err()
}

func (obj *String) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteString(obj.str)
	end := w.Pos()
	return int(end - pos), w.Err()
}

func init() {
	f := func() reflect.Value {
		o := &String{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("*rbase.String", f)
}

var (
	_ root.Object        = (*String)(nil)
	_ rbytes.Marshaler   = (*String)(nil)
	_ rbytes.Unmarshaler = (*String)(nil)
)
