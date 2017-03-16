// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"math"
	"reflect"
)

// H1F implements ROOT TH1F
type H1F struct {
	th1
	arr ArrayF
}

// Class returns the ROOT class name.
func (*H1F) Class() string {
	return "TH1F"
}

func (h *H1F) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH1F version too old (%d<2)", vers)
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

	r.CheckByteCount(pos, bcnt, beg, "TH1F")
	return r.err
}

func (h *H1F) Array() ArrayF {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1F) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1F) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1F) XAxis() Axis {
	return &h.th1.xaxis
}

// BinCenter returns the bin center value
func (h *H1F) BinCenter(i int) float64 {
	return h.th1.xaxis.BinCenter(i)
}

// BinContent returns the bin content
func (h *H1F) BinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// BinError returns the bin error
func (h *H1F) BinError(i int) float64 {
	return math.Sqrt(float64(h.th1.sumw2.Data[i]))
}

func init() {
	f := func() reflect.Value {
		o := &H1F{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1F", f)
	Factory.add("*rootio.H1F", f)
}

var _ Object = (*H1F)(nil)
var _ Named = (*H1F)(nil)
var _ ROOTUnmarshaler = (*H1F)(nil)

// H1D implements ROOT TH1D
type H1D struct {
	th1
	arr ArrayD
}

// Class returns the ROOT class name.
func (*H1D) Class() string {
	return "TH1D"
}

func (h *H1D) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH1D version too old (%d<2)", vers)
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

	r.CheckByteCount(pos, bcnt, beg, "TH1D")
	return r.err
}

func (h *H1D) Array() ArrayD {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1D) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1D) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1D) XAxis() Axis {
	return &h.th1.xaxis
}

// BinCenter returns the bin center value
func (h *H1D) BinCenter(i int) float64 {
	return h.th1.xaxis.BinCenter(i)
}

// BinContent returns the bin content
func (h *H1D) BinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// BinError returns the bin error
func (h *H1D) BinError(i int) float64 {
	return math.Sqrt(float64(h.th1.sumw2.Data[i]))
}

func init() {
	f := func() reflect.Value {
		o := &H1D{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1D", f)
	Factory.add("*rootio.H1D", f)
}

var _ Object = (*H1D)(nil)
var _ Named = (*H1D)(nil)
var _ ROOTUnmarshaler = (*H1D)(nil)

// H1I implements ROOT TH1I
type H1I struct {
	th1
	arr ArrayI
}

// Class returns the ROOT class name.
func (*H1I) Class() string {
	return "TH1I"
}

func (h *H1I) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH1I version too old (%d<2)", vers)
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

	r.CheckByteCount(pos, bcnt, beg, "TH1I")
	return r.err
}

func (h *H1I) Array() ArrayI {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H1I) Rank() int {
	return 1
}

// NbinsX returns the number of bins in X.
func (h *H1I) NbinsX() int {
	return h.th1.xaxis.nbins
}

// XAxis returns the axis along X.
func (h *H1I) XAxis() Axis {
	return &h.th1.xaxis
}

// BinCenter returns the bin center value
func (h *H1I) BinCenter(i int) float64 {
	return h.th1.xaxis.BinCenter(i)
}

// BinContent returns the bin content
func (h *H1I) BinContent(i int) float64 {
	return float64(h.arr.Data[i])
}

// BinError returns the bin error
func (h *H1I) BinError(i int) float64 {
	return math.Sqrt(float64(h.th1.sumw2.Data[i]))
}

func init() {
	f := func() reflect.Value {
		o := &H1I{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH1I", f)
	Factory.add("*rootio.H1I", f)
}

var _ Object = (*H1I)(nil)
var _ Named = (*H1I)(nil)
var _ ROOTUnmarshaler = (*H1I)(nil)
