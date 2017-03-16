// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type taxis struct {
	tnamed
	attaxis attaxis
	nbins   int        // number of bins
	xmin    float64    // low edge of first bin
	xmax    float64    // upper edge of last bin
	xbins   ArrayD     // bin edges array in X
	first   int        // first bin to display
	last    int        // last bin to display
	bits2   uint16     // second bit status word
	time    bool       // on/off displaying time values instead of numerics
	tfmt    string     // date&time format
	labels  *thashList // list of labels
	modlabs *tlist     // list of modified labels
}

func (a *taxis) Class() string {
	return "TAxis"
}

func (a *taxis) XMin() float64 {
	return a.xmin
}

func (a *taxis) XMax() float64 {
	return a.xmax
}

func (a *taxis) NBins() int {
	return a.nbins
}

func (a *taxis) XBins() []float64 {
	return a.xbins.Data
}

func (a *taxis) BinCenter(i int) float64 {
	if len(a.xbins.Data) == 0 || i < 1 || i > a.nbins {
		width := (a.xmax - a.xmin) / float64(a.nbins)
		return a.xmin + (float64(i)-0.5)*width
	}
	width := a.xbins.Data[i] - a.xbins.Data[i-1]
	return a.xbins.Data[i-1] + 0.5*width
}

func (a *taxis) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 9 {
		return fmt.Errorf("rootio: TAxis version too old (%d<9)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&a.tnamed,
		&a.attaxis,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	a.nbins = int(r.ReadI32())
	a.xmin = r.ReadF64()
	a.xmax = r.ReadF64()

	if err := a.xbins.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	a.first = int(r.ReadI32())
	a.last = int(r.ReadI32())
	a.bits2 = r.ReadU16()
	a.time = r.ReadBool()
	a.tfmt = r.ReadString()

	a.labels = nil
	labels := r.ReadObjectAny()
	if labels != nil {
		a.labels = labels.(*thashList)
	}

	a.modlabs = nil
	if vers >= 10 {
		modlabs := r.ReadObjectAny()
		if modlabs != nil {
			a.modlabs = modlabs.(*tlist)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TAxis")
	return r.err
}

func init() {
	{
		f := func() reflect.Value {
			o := &taxis{}
			return reflect.ValueOf(o)
		}
		Factory.add("TAxis", f)
		Factory.add("*rootio.taxis", f)
	}
}

var _ Object = (*taxis)(nil)
var _ Named = (*taxis)(nil)
var _ Axis = (*taxis)(nil)
var _ ROOTUnmarshaler = (*taxis)(nil)
