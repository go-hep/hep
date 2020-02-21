// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

type FeldmanCousins struct {
	obj rbase.Object

	CL      float64 // confidence level as a fraction [e.g. 90% = 0.9]
	UpLim   float64 // calculated upper limit
	LoLim   float64 // calculated lower limit
	Nobs    float64 // input number of observed events
	Nbkg    float64 // input number of background events
	MuMin   float64 // minimum value of signal to use in calculating the tables
	MuMax   float64 // maximum value of signal to use in calculating the tables
	MuStep  float64 // step in signal to use when generating tables
	NMuStep int32
	NMax    int32
	Quick   int32
}

func NewFeldmanCousins() *FeldmanCousins {
	return &FeldmanCousins{
		obj: *rbase.NewObject(),
	}
}

func (*FeldmanCousins) RVersion() int16 {
	return rvers.FeldmanCousins
}

func (*FeldmanCousins) Class() string {
	return "TFeldmanCousins"
}

func (fc *FeldmanCousins) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(fc.RVersion())
	if _, err := fc.obj.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteF64(fc.CL)
	w.WriteF64(fc.UpLim)
	w.WriteF64(fc.LoLim)
	w.WriteF64(fc.Nobs)
	w.WriteF64(fc.Nbkg)
	w.WriteF64(fc.MuMin)
	w.WriteF64(fc.MuMax)
	w.WriteF64(fc.MuStep)
	w.WriteI32(fc.NMuStep)
	w.WriteI32(fc.NMax)
	w.WriteI32(fc.Quick)

	return w.SetByteCount(pos, fc.Class())
}

func (fc *FeldmanCousins) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()

	_, pos, bcnt := r.ReadVersion(fc.Class())

	if err := fc.obj.UnmarshalROOT(r); err != nil {
		return err
	}

	fc.CL = r.ReadF64()
	fc.UpLim = r.ReadF64()
	fc.LoLim = r.ReadF64()
	fc.Nobs = r.ReadF64()
	fc.Nbkg = r.ReadF64()
	fc.MuMin = r.ReadF64()
	fc.MuMax = r.ReadF64()
	fc.MuStep = r.ReadF64()
	fc.NMuStep = r.ReadI32()
	fc.NMax = r.ReadI32()
	fc.Quick = r.ReadI32()

	r.CheckByteCount(pos, bcnt, beg, fc.Class())
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := NewFeldmanCousins()
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TFeldmanCousins", f)
	}
}

var (
	_ root.Object        = (*FeldmanCousins)(nil)
	_ rbytes.Marshaler   = (*FeldmanCousins)(nil)
	_ rbytes.Unmarshaler = (*FeldmanCousins)(nil)
)
