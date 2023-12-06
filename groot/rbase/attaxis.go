// Copyright ©2017 The go-hep Authors. All rights reserved.
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

type AttAxis struct {
	Ndivs       int32   // number of divisions (10000*n3 + 100*n2 + n1)
	AxisColor   int16   // color of the line axis
	LabelColor  int16   // color of labels
	LabelFont   int16   // font for labels
	LabelOffset float32 // offset of labels
	LabelSize   float32 // size of labels
	Ticks       float32 // length of tick marks
	TitleOffset float32 // offset of axis title
	TitleSize   float32 // size of axis title
	TitleColor  int16   // color of axis title
	TitleFont   int16   // font for axis title
}

func NewAttAxis() *AttAxis {
	return &AttAxis{
		Ndivs:       510, // FIXME(sbinet)
		AxisColor:   1,
		LabelColor:  1,
		LabelFont:   42,
		LabelOffset: 0.005,
		LabelSize:   0.035,
		Ticks:       0.03,
		TitleOffset: 1,
		TitleSize:   0.035,
		TitleColor:  1,
		TitleFont:   42,
	}
}

func (*AttAxis) RVersion() int16 {
	return rvers.AttAxis
}

func (*AttAxis) Class() string {
	return "TAttAxis"
}

func (a *AttAxis) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(a.Class(), a.RVersion())
	w.WriteI32(a.Ndivs)
	w.WriteI16(a.AxisColor)
	w.WriteI16(a.LabelColor)
	w.WriteI16(a.LabelFont)
	w.WriteF32(a.LabelOffset)
	w.WriteF32(a.LabelSize)
	w.WriteF32(a.Ticks)
	w.WriteF32(a.TitleOffset)
	w.WriteF32(a.TitleSize)
	w.WriteI16(a.TitleColor)
	w.WriteI16(a.TitleFont)

	return w.SetHeader(hdr)
}

func (a *AttAxis) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(a.Class(), a.RVersion())
	a.Ndivs = r.ReadI32()
	a.AxisColor = r.ReadI16()
	a.LabelColor = r.ReadI16()
	a.LabelFont = r.ReadI16()
	a.LabelOffset = r.ReadF32()
	a.LabelSize = r.ReadF32()
	a.Ticks = r.ReadF32()
	a.TitleOffset = r.ReadF32()
	a.TitleSize = r.ReadF32()
	a.TitleColor = r.ReadI16()
	a.TitleFont = r.ReadI16()

	r.CheckHeader(hdr)

	return r.Err()
}

func (a *AttAxis) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNdivisions", Value: &a.Ndivs},
		{Name: "fAxisColor", Value: &a.AxisColor},
		{Name: "fLabelColor", Value: &a.LabelColor},
		{Name: "fLabelFont", Value: &a.LabelFont},
		{Name: "fLabelOffset", Value: &a.LabelOffset},
		{Name: "fLabelSize", Value: &a.LabelSize},
		{Name: "fTickLength", Value: &a.Ticks},
		{Name: "fTitleOffset", Value: &a.TitleOffset},
		{Name: "fTitleSize", Value: &a.TitleSize},
		{Name: "fTitleColor", Value: &a.TitleColor},
		{Name: "fTitleFont", Value: &a.TitleFont},
	}...)
	return mbrs
}

func init() {
	f := func() reflect.Value {
		o := NewAttAxis()
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TAttAxis", f)
}

var (
	_ root.Object        = (*AttAxis)(nil)
	_ rbytes.Marshaler   = (*AttAxis)(nil)
	_ rbytes.Unmarshaler = (*AttAxis)(nil)
)
