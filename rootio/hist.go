// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

type th1 struct {
	rvers int16
	tnamed
	attline   attline
	attfill   attfill
	attmarker attmarker
	ncells    int       // number of bins + under/over-flows
	xaxis     taxis     // x axis descriptor
	yaxis     taxis     // y axis descriptor
	zaxis     taxis     // z axis descriptor
	boffset   int16     // (1000*offset) for bar charts or legos
	bwidth    int16     // (1000*width) for bar charts or legos
	entries   float64   // number of entries
	tsumw     float64   // total sum of weights
	tsumw2    float64   // total sum of squares of weights
	tsumwx    float64   // total sum of weight*x
	tsumwx2   float64   // total sum of weight*x*x
	max       float64   // maximum value for plotting
	min       float64   // minimum value for plotting
	norm      float64   // normalization factor
	contour   ArrayD    // array to display contour levels
	sumw2     ArrayD    // array of sum of squares of weights
	opt       string    // histogram options
	funcs     tlist     // list of functions (fits and user)
	buffer    []float64 // entry buffer
	erropt    int32     // option for bin statistical errors

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

func (h *th1) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	for _, v := range []ROOTUnmarshaler{
		&h.tnamed,
		&h.attline,
		&h.attfill,
		&h.attmarker,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	h.ncells = int(r.ReadI32())

	for _, v := range []ROOTUnmarshaler{
		&h.xaxis,
		&h.yaxis,
		&h.zaxis,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
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
			r.err = err
			return r.err
		}
	}

	if err := h.sumw2.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	h.opt = r.ReadString()
	if err := h.funcs.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	n := int(r.ReadI32())
	_ = r.ReadI8()
	h.buffer = r.ReadFastArrayF64(n)
	if vers > 6 {
		h.erropt = r.ReadI32()
	}

	r.CheckByteCount(pos, bcnt, beg, "TH1")
	return r.err
}

type th2 struct {
	rvers int16
	th1
	scale   float64 // scale factor
	tsumwy  float64 // total sum of weight*y
	tsumwy2 float64 // total sum of weight*y*y
	tsumwxy float64 // total sum of weight*x*y
}

func (*th2) Class() string {
	return "TH2"
}

func (h *th2) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	h.rvers = vers
	if vers < 3 {
		return errorf("rootio: TH2 version too old (%d<3)", vers)
	}

	if err := h.th1.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	h.scale = r.ReadF64()
	h.tsumwy = r.ReadF64()
	h.tsumwy2 = r.ReadF64()
	h.tsumwxy = r.ReadF64()

	r.CheckByteCount(pos, bcnt, beg, "TH2")
	return r.err
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

type dist0D struct {
	n      int64
	sumw   float64
	sumw2  float64
	sumwx  float64
	sumwx2 float64
}

func (d dist0D) Entries() int64 {
	return d.n
}

func (d dist0D) SumW() float64 {
	return d.sumw
}

func (d dist0D) SumW2() float64 {
	return d.sumw2
}

func (d dist0D) SumWX() float64 {
	return d.sumwx
}

func (d dist0D) SumWX2() float64 {
	return d.sumwx2
}

// dist2D is a 2-dim distribution.
type dist2D struct {
	x      dist0D  // x moments
	y      dist0D  // y moments
	sumWXY float64 // 2nd-order cross-term
}

// Entries returns the number of entries in the distribution.
func (d *dist2D) Entries() int64 {
	return d.x.Entries()
}

// SumW returns the sum of weights of the distribution.
func (d *dist2D) SumW() float64 {
	return d.x.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist2D) SumW2() float64 {
	return d.x.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *dist2D) SumWX() float64 {
	return d.x.SumWX()
}

// SumWX2 returns the 2nd order weighted x moment
func (d *dist2D) SumWX2() float64 {
	return d.x.SumWX2()
}

// SumWY returns the 1st order weighted y moment
func (d *dist2D) SumWY() float64 {
	return d.y.SumWX()
}

// SumWY2 returns the 2nd order weighted y moment
func (d *dist2D) SumWY2() float64 {
	return d.y.SumWX2()
}

func init() {
	{
		f := func() reflect.Value {
			o := &th1{}
			return reflect.ValueOf(o)
		}
		Factory.add("TH1", f)
		Factory.add("*rootio.th1", f)
	}
	{
		f := func() reflect.Value {
			o := &th2{}
			return reflect.ValueOf(o)
		}
		Factory.add("TH2", f)
		Factory.add("*rootio.th2", f)
	}
}

var _ Object = (*th1)(nil)
var _ Named = (*th1)(nil)
var _ ROOTUnmarshaler = (*th1)(nil)

var _ Object = (*th2)(nil)
var _ Named = (*th2)(nil)
var _ ROOTUnmarshaler = (*th2)(nil)
