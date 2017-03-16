// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
)

// H2F implements ROOT TH2F
type H2F struct {
	th2
	arr ArrayF
}

// Class returns the ROOT class name.
func (*H2F) Class() string {
	return "TH2F"
}

func (h *H2F) Array() ArrayF {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2F) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2F) NbinsX() int {
	return h.th1.xaxis.nbins
}

// NbinsY returns the number of bins in Y.
func (h *H2F) NbinsY() int {
	return h.th1.yaxis.nbins
}

func (h *H2F) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH2F version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2F")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2F{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2F", f)
	Factory.add("*rootio.H2F", f)
}

var _ Object = (*H2F)(nil)
var _ Named = (*H2F)(nil)
var _ ROOTUnmarshaler = (*H2F)(nil)

// H2D implements ROOT TH2D
type H2D struct {
	th2
	arr ArrayD
}

// Class returns the ROOT class name.
func (*H2D) Class() string {
	return "TH2D"
}

func (h *H2D) Array() ArrayD {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2D) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2D) NbinsX() int {
	return h.th1.xaxis.nbins
}

// NbinsY returns the number of bins in Y.
func (h *H2D) NbinsY() int {
	return h.th1.yaxis.nbins
}

func (h *H2D) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH2D version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2D")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2D{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2D", f)
	Factory.add("*rootio.H2D", f)
}

var _ Object = (*H2D)(nil)
var _ Named = (*H2D)(nil)
var _ ROOTUnmarshaler = (*H2D)(nil)

// H2I implements ROOT TH2I
type H2I struct {
	th2
	arr ArrayI
}

// Class returns the ROOT class name.
func (*H2I) Class() string {
	return "TH2I"
}

func (h *H2I) Array() ArrayI {
	return h.arr
}

// Rank returns the number of dimensions of this histogram.
func (h *H2I) Rank() int {
	return 2
}

// NbinsX returns the number of bins in X.
func (h *H2I) NbinsX() int {
	return h.th1.xaxis.nbins
}

// NbinsY returns the number of bins in Y.
func (h *H2I) NbinsY() int {
	return h.th1.yaxis.nbins
}

func (h *H2I) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: TH2I version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TH2I")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &H2I{}
		return reflect.ValueOf(o)
	}
	Factory.add("TH2I", f)
	Factory.add("*rootio.H2I", f)
}

var _ Object = (*H2I)(nil)
var _ Named = (*H2I)(nil)
var _ ROOTUnmarshaler = (*H2I)(nil)
