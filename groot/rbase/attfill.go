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

type AttFill struct {
	Color int16
	Style int16
}

func NewAttFill() *AttFill {
	return &AttFill{
		Color: 0,
		Style: 1001, // FIXME(sbinet)
	}
}

func (*AttFill) Class() string {
	return "TAttFill"
}

func (*AttFill) RVersion() int16 {
	return rvers.AttFill
}

func (a *AttFill) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteVersion(a.RVersion())
	w.WriteI16(a.Color)
	w.WriteI16(a.Style)
	return w.SetByteCount(pos, "TAttFill")

}

func (a *AttFill) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	a.Color = r.ReadI16()
	a.Style = r.ReadI16()
	r.CheckByteCount(pos, bcnt, start, "TAttFill")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewAttFill()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TAttFill", f)
}

var (
	_ root.Object        = (*AttFill)(nil)
	_ rbytes.Marshaler   = (*AttFill)(nil)
	_ rbytes.Unmarshaler = (*AttFill)(nil)
)
