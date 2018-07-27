// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
)

type taxis struct {
	rvers int16
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

func (a *taxis) BinLowEdge(i int) float64 {
	if len(a.xbins.Data) == 0 || i < 1 || i > a.nbins {
		width := (a.xmax - a.xmin) / float64(a.nbins)
		return a.xmin + float64(i-1)*width
	}
	return a.xbins.Data[i-1]
}

func (a *taxis) BinWidth(i int) float64 {
	if a.nbins <= 0 {
		return 0
	}
	if len(a.xbins.Data) <= 0 {
		return (a.xmax - a.xmin) / float64(a.nbins)
	}
	if i > a.nbins {
		i = a.nbins
	}
	if i < 1 {
		i = 1
	}
	return a.xbins.Data[i] - a.xbins.Data[i-1]
}

func (a *taxis) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(a.rvers)

	for _, v := range []ROOTMarshaler{
		&a.tnamed,
		&a.attaxis,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	w.WriteI32(int32(a.nbins))
	w.WriteF64(a.xmin)
	w.WriteF64(a.xmax)

	if _, err := a.xbins.MarshalROOT(w); err != nil {
		w.err = err
		return 0, w.err
	}

	w.WriteI32(int32(a.first))
	w.WriteI32(int32(a.last))
	w.WriteU16(a.bits2)
	w.WriteBool(a.time)
	w.WriteString(a.tfmt)

	// FIXME
	//	a.labels = nil
	//	labels := r.ReadObjectAny()
	//	if labels != nil {
	//		a.labels = labels.(*thashList)
	//	}
	//	a.modlabs = nil
	//	modlabs := r.ReadObjectAny()
	//	if modlabs != nil {
	//		a.modlabs = modlabs.(*tlist)
	//	}

	return w.SetByteCount(pos, "TAxis")
}

func (a *taxis) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	a.rvers = vers
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

var (
	_ Object          = (*taxis)(nil)
	_ Named           = (*taxis)(nil)
	_ Axis            = (*taxis)(nil)
	_ ROOTMarshaler   = (*taxis)(nil)
	_ ROOTUnmarshaler = (*taxis)(nil)
)
