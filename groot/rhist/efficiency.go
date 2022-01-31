// Copyright ©2022 The go-hep Authors. All rights reserved.
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

// Efficiency handles efficiency histograms.
type Efficiency struct {
	named   rbase.Named
	attline rbase.AttLine
	attfill rbase.AttFill
	attmark rbase.AttMarker

	betaAlpha float64 // global parameter for prior beta distribution (default = 1)
	betaBeta  float64 // global parameter for prior beta distribution (default = 1)

	betaBinParams [][2]float64 // parameter for prior beta distribution different bin by bin

	confLvl float64     // confidence level (default = 0.683, 1 sigma)
	funcs   *rcont.List // ->pointer to list of functions

	passedHist H1      // histogram for events which passed certain criteria
	statOpt    int32   // defines how the confidence intervals are determined
	totHist    H1      // histogram for total number of events
	weight     float64 // weight for all events (default = 1)
}

func (*Efficiency) Class() string {
	return "TEfficiency"
}

func (*Efficiency) RVersion() int16 {
	return rvers.Efficiency
}

// MarshalROOT implements rbytes.Marshaler
func (o *Efficiency) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

	for _, v := range []rbytes.Marshaler{
		&o.named, &o.attline, &o.attfill, &o.attmark,
	} {
		n, err := v.MarshalROOT(w)
		if err != nil {
			return n, err
		}
	}
	w.WriteF64(o.betaAlpha)
	w.WriteF64(o.betaBeta)
	if n, err := writeVecPairF64(w, o.betaBinParams); err != nil {
		return n, err
	}
	w.WriteF64(o.confLvl)
	if err := w.WriteObject(o.funcs); err != nil { // obj-ptr
		return 0, err
	}
	if err := w.WriteObjectAny(o.passedHist); err != nil { // obj-ptr
		return 0, err
	}
	w.WriteI32(o.statOpt)
	if err := w.WriteObjectAny(o.totHist); err != nil { // obj-ptr
		return 0, err
	}
	w.WriteF64(o.weight)

	return w.SetByteCount(pos, o.Class())
}

func writeVecPairF64(w *rbytes.WBuffer, vs [][2]float64) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(rvers.StreamerInfo | rbytes.StreamedMemberWise)
	w.WriteI16(0)        // class version
	w.WriteU32(0xd7bed2) // checksum
	w.WriteI32(int32(len(vs)))
	for i := range vs {
		w.WriteF64(vs[i][0])
	}
	for i := range vs {
		w.WriteF64(vs[i][1])
	}
	return w.SetByteCount(pos, "vector<pair<double,double> >")
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (o *Efficiency) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion(o.Class())
	if vers > rvers.Efficiency {
		panic(fmt.Errorf("rhist: invalid TEfficiency version=%d > %d", vers, rvers.Efficiency))
	}

	for _, v := range []rbytes.Unmarshaler{
		&o.named, &o.attline, &o.attfill, &o.attmark,
	} {
		err := v.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}
	o.betaAlpha = r.ReadF64()
	o.betaBeta = r.ReadF64()
	if err := readVecPairF64(r, &o.betaBinParams); err != nil {
		return err
	}
	o.confLvl = r.ReadF64()
	{
		o.funcs = nil
		if oo := r.ReadObject("TList"); oo != nil { // obj-ptr
			o.funcs = oo.(*rcont.List)
		}
	}
	{
		o.passedHist = nil
		if oo := r.ReadObjectAny(); oo != nil { // obj-ptr
			o.passedHist = oo.(H1)
		}
	}
	o.statOpt = r.ReadI32()
	{
		o.totHist = nil
		if oo := r.ReadObjectAny(); oo != nil { // obj-ptr
			o.totHist = oo.(H1)
		}
	}
	o.weight = r.ReadF64()

	r.CheckByteCount(pos, bcnt, start, o.Class())
	return r.Err()
}

func readVecPairF64(r *rbytes.RBuffer, vs *[][2]float64) error {
	if r.Err() != nil {
		return r.Err()
	}

	start := r.Pos()
	const typename = "vector<pair<double,double> >"
	vers, pos, bcnt := r.ReadVersion(typename)
	mbrwise := vers&rbytes.StreamedMemberWise != 0
	if mbrwise {
		vers &^= rbytes.StreamedMemberWise
	}
	if vers != rvers.StreamerInfo {
		r.SetErr(fmt.Errorf(
			"rbytes: invalid version for %q. got=%v, want=%v",
			typename, vers, rvers.StreamerInfo,
		))
		return r.Err()
	}
	if mbrwise {
		clvers := r.ReadI16()
		switch {
		case clvers == 1:
			// TODO
		case clvers <= 0:
			/*chksum*/ _ = r.ReadU32()
		}
	}
	n := int(r.ReadI32())
	if n == 0 {
		*vs = nil
		r.CheckByteCount(pos, bcnt, start, typename)
		return r.Err()
	}

	*vs = make([][2]float64, n)
	switch {
	case mbrwise:
		p := make([]float64, n)
		r.ReadArrayF64(p)
		for i := range *vs {
			(*vs)[i][0] = p[i]
		}
		r.ReadArrayF64(p)
		for i := range *vs {
			(*vs)[i][1] = p[i]
		}
	default:
		for i := range *vs {
			(*vs)[i][0] = r.ReadF64()
			(*vs)[i][1] = r.ReadF64()
		}
	}
	r.CheckByteCount(pos, bcnt, start, typename)
	return r.Err()
}

func init() {
	f := func() reflect.Value {
		var o Efficiency
		return reflect.ValueOf(&o)
	}
	rtypes.Factory.Add("TEfficiency", f)
}

var (
	_ root.Object        = (*Efficiency)(nil)
	_ rbytes.RVersioner  = (*Efficiency)(nil)
	_ rbytes.Marshaler   = (*Efficiency)(nil)
	_ rbytes.Unmarshaler = (*Efficiency)(nil)
)
