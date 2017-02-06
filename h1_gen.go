// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
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
