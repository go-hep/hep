// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcolors"
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
		Color: rcolors.Blue + 2,
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

	hdr := w.WriteHeader(a.Class(), a.RVersion())
	w.WriteI16(a.Color)
	w.WriteI16(a.Style)
	w.WriteI16(a.Width)
	return w.SetHeader(hdr)
}

func (a *AttLine) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(a.Class(), a.RVersion())

	a.Color = r.ReadI16()
	a.Style = r.ReadI16()
	a.Width = r.ReadI16()

	r.CheckHeader(hdr)
	return r.Err()
}

func (a *AttLine) RMembers() []rbytes.Member {
	return []rbytes.Member{
		{Name: "fLineColor", Value: &a.Color},
		{Name: "fLineStyle", Value: &a.Style},
		{Name: "fLineWidth", Value: &a.Width},
	}
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
