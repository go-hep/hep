// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"golang.org/x/xerrors"
)

type th1 struct {
	rbase.Named
	attline   rbase.AttLine
	attfill   rbase.AttFill
	attmarker rbase.AttMarker
	ncells    int          // number of bins + under/over-flows
	xaxis     taxis        // x axis descriptor
	yaxis     taxis        // y axis descriptor
	zaxis     taxis        // z axis descriptor
	boffset   int16        // (1000*offset) for bar charts or legos
	bwidth    int16        // (1000*width) for bar charts or legos
	entries   float64      // number of entries
	tsumw     float64      // total sum of weights
	tsumw2    float64      // total sum of squares of weights
	tsumwx    float64      // total sum of weight*x
	tsumwx2   float64      // total sum of weight*x*x
	max       float64      // maximum value for plotting
	min       float64      // minimum value for plotting
	norm      float64      // normalization factor
	contour   rcont.ArrayD // array to display contour levels
	sumw2     rcont.ArrayD // array of sum of squares of weights
	opt       string       // histogram options
	funcs     rcont.List   // list of functions (fits and user)
	buffer    []float64    // entry buffer
	erropt    int32        // option for bin statistical errors
	oflow     int32        // per object flag to use under/overflows in statistics
}

func newH1() *th1 {
	return &th1{
		Named:     *rbase.NewNamed("", ""),
		attline:   *rbase.NewAttLine(),
		attfill:   *rbase.NewAttFill(),
		attmarker: *rbase.NewAttMarker(),
		xaxis:     *NewAxis("xaxis"),
		yaxis:     *NewAxis("yaxis"),
		zaxis:     *NewAxis("zaxis"),
		bwidth:    1000,
		max:       -1111,
		min:       -1111,
		funcs:     *rcont.NewList("", nil),
		oflow:     2, // kNeutral
	}
}

func (*th1) RVersion() int16 {
	return rvers.H1
}

func (h *th1) Class() string {
	return "TH1"
}

// Entries returns the number of entries for this histogram.
func (h *th1) Entries() float64 {
	return h.entries
}

// SumW returns the total sum of weights
func (h *th1) SumW() float64 {
	return h.tsumw
}

// SumW2 returns the total sum of squares of weights
func (h *th1) SumW2() float64 {
	return h.tsumw2
}

// SumWX returns the total sum of weights*x
func (h *th1) SumWX() float64 {
	return h.tsumwx
}

// SumWX2 returns the total sum of weights*x*x
func (h *th1) SumWX2() float64 {
	return h.tsumwx2
}

// SumW2s returns the array of sum of squares of weights
func (h *th1) SumW2s() []float64 {
	return h.sumw2.Data
}

func (h *th1) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())
	for _, v := range []rbytes.Marshaler{
		&h.Named,
		&h.attline,
		&h.attfill,
		&h.attmarker,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	w.WriteI32(int32(h.ncells))

	for _, v := range []rbytes.Marshaler{
		&h.xaxis,
		&h.yaxis,
		&h.zaxis,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	w.WriteI16(h.boffset)
	w.WriteI16(h.bwidth)
	w.WriteF64(h.entries)
	w.WriteF64(h.tsumw)
	w.WriteF64(h.tsumw2)
	w.WriteF64(h.tsumwx)
	w.WriteF64(h.tsumwx2)
	w.WriteF64(h.max)
	w.WriteF64(h.min)
	w.WriteF64(h.norm)
	if _, err := h.contour.MarshalROOT(w); err != nil {
		return 0, err
	}

	if _, err := h.sumw2.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteString(h.opt)
	if _, err := h.funcs.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteI32(int32(len(h.buffer)))
	w.WriteI8(0) // FIXME(sbinet)
	w.WriteFastArrayF64(h.buffer)
	w.WriteI32(h.erropt)
	if h.RVersion() > 7 {
		w.WriteI32(h.oflow)
	}

	return w.SetByteCount(pos, h.Class())
}

func (h *th1) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	for _, v := range []rbytes.Unmarshaler{
		&h.Named,
		&h.attline,
		&h.attfill,
		&h.attmarker,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	h.ncells = int(r.ReadI32())

	for _, v := range []rbytes.Unmarshaler{
		&h.xaxis,
		&h.yaxis,
		&h.zaxis,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	h.boffset = r.ReadI16()
	h.bwidth = r.ReadI16()
	h.entries = r.ReadF64()
	h.tsumw = r.ReadF64()
	h.tsumw2 = r.ReadF64()
	h.tsumwx = r.ReadF64()
	h.tsumwx2 = r.ReadF64()
	if vers < 2 {
		h.max = float64(r.ReadF32())
		h.min = float64(r.ReadF32())
		h.norm = float64(r.ReadF32())
		n := int(r.ReadI32())
		h.contour.Data = r.ReadFastArrayF64(n)
	} else {
		h.max = r.ReadF64()
		h.min = r.ReadF64()
		h.norm = r.ReadF64()
		if err := h.contour.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	if err := h.sumw2.UnmarshalROOT(r); err != nil {
		return err
	}

	h.opt = r.ReadString()
	if err := h.funcs.UnmarshalROOT(r); err != nil {
		return err
	}

	n := int(r.ReadI32())
	_ = r.ReadI8()
	h.buffer = r.ReadFastArrayF64(n)
	if vers > 6 {
		h.erropt = r.ReadI32()
	}
	if vers > 7 {
		h.oflow = r.ReadI32()
	}

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

type th2 struct {
	th1
	scale   float64 // scale factor
	tsumwy  float64 // total sum of weight*y
	tsumwy2 float64 // total sum of weight*y*y
	tsumwxy float64 // total sum of weight*x*y
}

func newH2() *th2 {
	return &th2{
		th1: *newH1(),
	}
}

func (*th2) RVersion() int16 {
	return rvers.H2
}

func (*th2) Class() string {
	return "TH2"
}

func (h *th2) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(h.RVersion())

	if _, err := h.th1.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteF64(h.scale)
	w.WriteF64(h.tsumwy)
	w.WriteF64(h.tsumwy2)
	w.WriteF64(h.tsumwxy)

	return w.SetByteCount(pos, h.Class())
}

func (h *th2) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(h.Class())
	if vers < 3 {
		return xerrors.Errorf("rhist: TH2 version too old (%d<3)", vers)
	}

	if err := h.th1.UnmarshalROOT(r); err != nil {
		return err
	}

	h.scale = r.ReadF64()
	h.tsumwy = r.ReadF64()
	h.tsumwy2 = r.ReadF64()
	h.tsumwxy = r.ReadF64()

	r.CheckByteCount(pos, bcnt, beg, h.Class())
	return r.Err()
}

// SumWY returns the total sum of weights*y
func (h *th2) SumWY() float64 {
	return h.tsumwy
}

// SumWY2 returns the total sum of weights*y*y
func (h *th2) SumWY2() float64 {
	return h.tsumwy2
}

// SumWXY returns the total sum of weights*x*y
func (h *th2) SumWXY() float64 {
	return h.tsumwxy
}

func init() {
	{
		f := func() reflect.Value {
			o := newH1()
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TH1", f)
	}
	{
		f := func() reflect.Value {
			o := newH2()
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TH2", f)
	}
}

var (
	_ root.Object        = (*th1)(nil)
	_ root.Named         = (*th1)(nil)
	_ rbytes.Marshaler   = (*th1)(nil)
	_ rbytes.Unmarshaler = (*th1)(nil)

	_ root.Object        = (*th2)(nil)
	_ root.Named         = (*th2)(nil)
	_ rbytes.Marshaler   = (*th2)(nil)
	_ rbytes.Unmarshaler = (*th2)(nil)
)
