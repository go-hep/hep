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

type AttLine struct {
	Color int16
	Style int16
	Width int16
}

func NewAttLine() *AttLine {
	return &AttLine{
		Color: 1, // FIXME(sbinet)
		Style: 1,
		Width: 1,
	}
}

func (*AttLine) Class() string {
	return "TAttLine"
}

func (*AttLine) RVersion() int16 {
	return rvers.AttLine
}

func (a *AttLine) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteVersion(a.RVersion())
	w.WriteI16(a.Color)
	w.WriteI16(a.Style)
	w.WriteI16(a.Width)
	return w.SetByteCount(pos, "TAttLine")

}

func (a *AttLine) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion()

	a.Color = r.ReadI16()
	a.Style = r.ReadI16()
	a.Width = r.ReadI16()
	r.CheckByteCount(pos, bcnt, start, "TAttLine")

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := NewAttLine()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TAttLine", f)
}

var (
	_ root.Object        = (*AttLine)(nil)
	_ rbytes.Marshaler   = (*AttLine)(nil)
	_ rbytes.Unmarshaler = (*AttLine)(nil)
)
