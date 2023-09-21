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
		nbins:   1,
		xmin:    0,
		xmax:    1,
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

	hdr := w.WriteHeader(a.Class(), a.RVersion())

	w.WriteObject(&a.Named)
	w.WriteObject(&a.attaxis)

	w.WriteI32(int32(a.nbins))
	w.WriteF64(a.xmin)
	w.WriteF64(a.xmax)
	w.WriteObject(&a.xbins)
	w.WriteI32(int32(a.first))
	w.WriteI32(int32(a.last))
	w.WriteU16(a.bits2)
	w.WriteBool(a.time)
	w.WriteString(a.tfmt)

	w.WriteObjectAny(a.labels)
	if a.RVersion() >= 10 {
		w.WriteObjectAny(a.modlabs)
	}

	return w.SetHeader(hdr)
}

func (a *taxis) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(a.Class())
	if hdr.Vers > rvers.Axis {
		panic(fmt.Errorf("rhist: invalid TAxis version=%d > %d", hdr.Vers, rvers.Axis))
	}
	const minVers = 6
	if hdr.Vers < minVers {
		return fmt.Errorf("rhist: TAxis version too old (%d<%d)", hdr.Vers, minVers)
	}

	r.ReadObject(&a.Named)
	r.ReadObject(&a.attaxis)

	a.nbins = int(r.ReadI32())
	a.xmin = r.ReadF64()
	a.xmax = r.ReadF64()
	r.ReadObject(&a.xbins)
	a.first = int(r.ReadI32())
	a.last = int(r.ReadI32())
	if hdr.Vers >= 9 {
		a.bits2 = r.ReadU16()
	}
	a.time = r.ReadBool()
	a.tfmt = r.ReadString()

	a.labels = nil
	if hdr.Vers >= 9 {
		labels := r.ReadObjectAny()
		if labels != nil {
			a.labels = labels.(*rcont.HashList)
		}
	}

	a.modlabs = nil
	if hdr.Vers >= 10 {
		modlabs := r.ReadObjectAny()
		if modlabs != nil {
			a.modlabs = modlabs.(*rcont.List)
		}
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func (a *taxis) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, a.Named.RMembers()...)
	mbrs = append(mbrs, a.attaxis.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNbins", Value: &a.nbins},
		{Name: "fXmin", Value: &a.xmin},
		{Name: "fXmax", Value: &a.xmax},
		{Name: "fXbins", Value: &a.xbins.Data},
		{Name: "fFirst", Value: &a.first},
		{Name: "fLast", Value: &a.last},
		{Name: "fBits2", Value: &a.bits2},
		{Name: "fTimeDisplay", Value: &a.time},
		{Name: "fTimeFormat", Value: &a.tfmt},
		{Name: "fLabels", Value: &a.labels},
		{Name: "fModLabs", Value: &a.modlabs},
	}...)
	return mbrs
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
