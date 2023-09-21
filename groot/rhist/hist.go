// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
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

	hdr := w.WriteHeader(h.Class(), h.RVersion())
	w.WriteObject(&h.Named)
	w.WriteObject(&h.attline)
	w.WriteObject(&h.attfill)
	w.WriteObject(&h.attmarker)

	w.WriteI32(int32(h.ncells))

	w.WriteObject(&h.xaxis)
	w.WriteObject(&h.yaxis)
	w.WriteObject(&h.zaxis)
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
	w.WriteObject(&h.contour)
	w.WriteObject(&h.sumw2)
	w.WriteString(h.opt)
	w.WriteObject(&h.funcs)

	w.WriteI32(int32(len(h.buffer)))
	w.WriteI8(0) // FIXME(sbinet)
	w.WriteArrayF64(h.buffer)
	w.WriteI32(h.erropt)
	if h.RVersion() > 7 {
		w.WriteI32(h.oflow)
	}

	return w.SetHeader(hdr)
}

func (h *th1) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(h.Class())
	if hdr.Vers > rvers.H1 {
		panic(fmt.Errorf("rhist: invalid TH1 version=%d > %d", hdr.Vers, rvers.H1))
	}

	r.ReadObject(&h.Named)
	r.ReadObject(&h.attline)
	r.ReadObject(&h.attfill)
	r.ReadObject(&h.attmarker)

	h.ncells = int(r.ReadI32())

	r.ReadObject(&h.xaxis)
	r.ReadObject(&h.yaxis)
	r.ReadObject(&h.zaxis)

	h.boffset = r.ReadI16()
	h.bwidth = r.ReadI16()
	h.entries = r.ReadF64()
	h.tsumw = r.ReadF64()
	h.tsumw2 = r.ReadF64()
	h.tsumwx = r.ReadF64()
	h.tsumwx2 = r.ReadF64()
	if hdr.Vers < 2 {
		h.max = float64(r.ReadF32())
		h.min = float64(r.ReadF32())
		h.norm = float64(r.ReadF32())
		n := int(r.ReadI32())
		h.contour.Data = rbytes.ResizeF64(h.contour.Data, n)
		r.ReadArrayF64(h.contour.Data)
	} else {
		h.max = r.ReadF64()
		h.min = r.ReadF64()
		h.norm = r.ReadF64()
		r.ReadObject(&h.contour)
	}

	r.ReadObject(&h.sumw2)
	h.opt = r.ReadString()
	r.ReadObject(&h.funcs)

	if hdr.Vers > 3 {
		n := int(r.ReadI32())
		_ = r.ReadI8()
		h.buffer = rbytes.ResizeF64(h.buffer, n)
		r.ReadArrayF64(h.buffer)
		if hdr.Vers > 6 {
			h.erropt = r.ReadI32()
		}
		if hdr.Vers > 7 {
			h.oflow = r.ReadI32()
		}
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func (h *th1) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, h.Named.RMembers()...)
	mbrs = append(mbrs, h.attline.RMembers()...)
	mbrs = append(mbrs, h.attfill.RMembers()...)
	mbrs = append(mbrs, h.attmarker.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNcells", Value: &h.ncells},
		{Name: "fXaxis", Value: &h.xaxis},
		{Name: "fYaxis", Value: &h.yaxis},
		{Name: "fZaxis", Value: &h.zaxis},
		{Name: "fBarOffset", Value: &h.boffset},
		{Name: "fBarWidth", Value: &h.bwidth},
		{Name: "fEntries", Value: &h.entries},
		{Name: "fTsumw", Value: &h.tsumw},
		{Name: "fTsumw2", Value: &h.tsumw2},
		{Name: "fTsumwx", Value: &h.tsumwx},
		{Name: "fTsumwx2", Value: &h.tsumwx2},
		{Name: "fMaximum", Value: &h.max},
		{Name: "fMinimum", Value: &h.min},
		{Name: "fNormFactor", Value: &h.norm},
		{Name: "fContour", Value: &h.contour.Data},
		{Name: "fSumw2", Value: &h.sumw2.Data},
		{Name: "fOption", Value: &h.opt},
		{Name: "fFunctions", Value: &h.funcs},
		{Name: "fBufferSize", Value: len(h.buffer)}, // FIXME(sbinet)
		{Name: "fBuffer", Value: &h.buffer},
		{Name: "fBinStatErrOpt", Value: &h.erropt},
		{Name: "fStatOverflows", Value: &h.oflow},
	}...)

	return mbrs
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

	hdr := w.WriteHeader(h.Class(), h.RVersion())

	w.WriteObject(&h.th1)
	w.WriteF64(h.scale)
	w.WriteF64(h.tsumwy)
	w.WriteF64(h.tsumwy2)
	w.WriteF64(h.tsumwxy)

	return w.SetHeader(hdr)
}

func (h *th2) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(h.Class())
	if hdr.Vers > rvers.H2 {
		panic(fmt.Errorf("rhist: invalid TH2 version=%d > %d", hdr.Vers, rvers.H2))
	}
	if hdr.Vers < 3 {
		return fmt.Errorf("rhist: TH2 version too old (%d<3)", hdr.Vers)
	}

	r.ReadObject(&h.th1)
	h.scale = r.ReadF64()
	h.tsumwy = r.ReadF64()
	h.tsumwy2 = r.ReadF64()
	h.tsumwxy = r.ReadF64()

	r.CheckHeader(hdr)
	return r.Err()
}

func (h *th2) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, h.th1.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fScalefactor", Value: &h.scale},
		{Name: "fTsumwy", Value: &h.tsumwy},
		{Name: "fTsumwy2", Value: &h.tsumwy2},
		{Name: "fTsumwxy", Value: &h.tsumwxy},
	}...)

	return mbrs
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
	_ rbytes.RSlicer     = (*th1)(nil)

	_ root.Object        = (*th2)(nil)
	_ root.Named         = (*th2)(nil)
	_ rbytes.Marshaler   = (*th2)(nil)
	_ rbytes.Unmarshaler = (*th2)(nil)
)
