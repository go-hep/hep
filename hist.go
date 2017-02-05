// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type th1 struct {
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

func (h *th1) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 7 {
		return fmt.Errorf("rootio: TH1 version too old (%d<7)", vers)
	}

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
	h.max = r.ReadF64()
	h.min = r.ReadF64()
	h.norm = r.ReadF64()

	for _, v := range []ROOTUnmarshaler{
		&h.contour,
		&h.sumw2,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	h.opt = r.ReadString()
	if err := h.funcs.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	n := int(r.ReadI32())
	h.buffer = r.ReadFastArrayF64(n)
	h.erropt = r.ReadI32()

	r.CheckByteCount(pos, bcnt, beg, "TH1")
	return r.err
}

type th1f struct {
	th1
	arr ArrayF
}

func (th1f) Class() string {
	return "TH1F"
}

func (h *th1f) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return fmt.Errorf("rootio: TH1F version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH1")
	return r.err
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
			o := &th1f{}
			return reflect.ValueOf(o)
		}
		Factory.add("TH1F", f)
		Factory.add("*rootio.th1f", f)
	}
}

var _ Object = (*th1)(nil)
var _ Named = (*th1)(nil)
var _ ROOTUnmarshaler = (*th1)(nil)

var _ Object = (*th1f)(nil)
var _ Named = (*th1f)(nil)
var _ ROOTUnmarshaler = (*th1f)(nil)
