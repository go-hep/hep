// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type AttMarker struct {
	Color int16
	Style int16
	Width float32
}

func NewAttMarker() *AttMarker {
	return &AttMarker{
		Color: 1,
		Style: 1,
		Width: 1,
	}
}

func (*AttMarker) Class() string {
	return "TAttMarker"
}

func (*AttMarker) RVersion() int16 {
	return rvers.AttMarker
}

func (a *AttMarker) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(a.Class(), a.RVersion())
	w.WriteI16(a.Color)
	w.WriteI16(a.Style)
	w.WriteF32(a.Width)
	return w.SetHeader(hdr)
}

func (a *AttMarker) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(a.Class())
	if hdr.Vers > rvers.AttMarker {
		panic(fmt.Errorf("rbase: invalid attmarker version=%d > %d", hdr.Vers, rvers.AttMarker))
	}

	a.Color = r.ReadI16()
	a.Style = r.ReadI16()
	a.Width = r.ReadF32()

	r.CheckHeader(hdr)
	return r.Err()
}

func (a *AttMarker) RMembers() []rbytes.Member {
	return []rbytes.Member{
		{Name: "fMarkerColor", Value: &a.Color},
		{Name: "fMarkerStyle", Value: &a.Style},
		{Name: "fMarkerSize", Value: &a.Width},
	}
}

func init() {
	f := func() reflect.Value {
		o := NewAttMarker()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TAttMarker", f)
}

var (
	_ root.Object        = (*AttMarker)(nil)
	_ rbytes.Marshaler   = (*AttMarker)(nil)
	_ rbytes.Unmarshaler = (*AttMarker)(nil)
)
