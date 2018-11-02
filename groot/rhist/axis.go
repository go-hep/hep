// Copyright 2017 The go-hep Authors. All rights reserved.
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

type taxis struct {
	rbase.Named
	attaxis rbase.AttAxis
	nbins   int             // number of bins
	xmin    float64         // low edge of first bin
	xmax    float64         // upper edge of last bin
	xbins   rcont.ArrayD    // bin edges array in X
	first   int             // first bin to display
	last    int             // last bin to display
	bits2   uint16          // second bit status word
	time    bool            // on/off displaying time values instead of numerics
	tfmt    string          // date&time format
	labels  *rcont.HashList // list of labels
	modlabs *rcont.List     // list of modified labels
}

func NewAxis(name string) *taxis {
	return &taxis{
		Named:   *rbase.NewNamed(name, ""),
		attaxis: *rbase.NewAttAxis(),
	}
}

func (*taxis) RVersion() int16 {
	return rvers.Axis
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

func (a *taxis) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.Pos()
	w.WriteVersion(a.RVersion())

	for _, v := range []rbytes.Marshaler{
		&a.Named,
		&a.attaxis,
	} {
		if _, err := v.MarshalROOT(w); err != nil {
			return 0, err
		}
	}

	w.WriteI32(int32(a.nbins))
	w.WriteF64(a.xmin)
	w.WriteF64(a.xmax)

	if _, err := a.xbins.MarshalROOT(w); err != nil {
		return 0, err
	}

	w.WriteI32(int32(a.first))
	w.WriteI32(int32(a.last))
	w.WriteU16(a.bits2)
	w.WriteBool(a.time)
	w.WriteString(a.tfmt)

	if err := w.WriteObjectAny(a.labels); err != nil {
		return 0, err
	}

	if a.RVersion() >= 10 {
		if err := w.WriteObjectAny(a.modlabs); err != nil {
			return 0, err
		}
	}

	return w.SetByteCount(pos, "TAxis")
}

func (a *taxis) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 9 {
		return fmt.Errorf("rhist: TAxis version too old (%d<9)", vers)
	}

	for _, v := range []rbytes.Unmarshaler{
		&a.Named,
		&a.attaxis,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			return err
		}
	}

	a.nbins = int(r.ReadI32())
	a.xmin = r.ReadF64()
	a.xmax = r.ReadF64()

	if err := a.xbins.UnmarshalROOT(r); err != nil {
		return err
	}

	a.first = int(r.ReadI32())
	a.last = int(r.ReadI32())
	a.bits2 = r.ReadU16()
	a.time = r.ReadBool()
	a.tfmt = r.ReadString()

	a.labels = nil
	labels := r.ReadObjectAny()
	if labels != nil {
		a.labels = labels.(*rcont.HashList)
	}

	a.modlabs = nil
	if vers >= 10 {
		modlabs := r.ReadObjectAny()
		if modlabs != nil {
			a.modlabs = modlabs.(*rcont.List)
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "TAxis")
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := NewAxis("")
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TAxis", f)
	}
}

var (
	_ root.Object        = (*taxis)(nil)
	_ root.Named         = (*taxis)(nil)
	_ Axis               = (*taxis)(nil)
	_ rbytes.Marshaler   = (*taxis)(nil)
	_ rbytes.Unmarshaler = (*taxis)(nil)
)
