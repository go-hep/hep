// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

type tgraph struct {
	rvers int16
	named tnamed

	maxsize int32
	npoints int32
	x       []float64
	y       []float64
	funcs   List
	histo   *H1F
	min     float64
	max     float64
}

func (g *tgraph) Class() string {
	return "TGraph"
}

func (g *tgraph) Name() string {
	return g.named.Name()
}

func (g *tgraph) Title() string {
	return g.named.Title()
}

func (g *tgraph) Len() int {
	return int(len(g.x))
}

func (g *tgraph) XY(i int) (float64, float64) {
	return g.x[i], g.y[i]
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgraph) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()

	w.WriteVersion(g.rvers)

	for _, a := range []ROOTMarshaler{
		&g.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		if _, err := a.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}
	}

	w.WriteI32(g.npoints)
	{
		w.WriteI8(0)
		w.WriteFastArrayF64(g.x)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.y)
	}

	w.WriteObjectAny(g.funcs)
	w.WriteObjectAny(g.histo)
	{
		w.WriteF64(g.min)
		w.WriteF64(g.max)
	}

	return w.SetByteCount(pos, "TGraph")
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraph) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	g.rvers = vers

	for _, a := range []ROOTUnmarshaler{
		&g.named,
		&attline{},
		&attfill{},
		&attmarker{},
	} {
		err := a.UnmarshalROOT(r)
		if err != nil {
			return err
		}
	}

	g.npoints = r.ReadI32()
	g.maxsize = g.npoints
	if vers < 2 {
		_ = r.ReadI8()
		xs := r.ReadFastArrayF32(int(g.npoints))
		_ = r.ReadI8()
		ys := r.ReadFastArrayF32(int(g.npoints))
		g.x = make([]float64, len(xs))
		g.y = make([]float64, len(ys))
		for i := range xs {
			g.x[i] = float64(xs[i])
			g.y[i] = float64(ys[i])
		}

	} else {
		_ = r.ReadI8()
		g.x = r.ReadFastArrayF64(int(g.npoints))
		_ = r.ReadI8()
		g.y = r.ReadFastArrayF64(int(g.npoints))
	}

	funcs := r.ReadObjectAny()
	if funcs != nil {
		g.funcs = funcs.(List)
	}

	histo := r.ReadObjectAny()
	if histo != nil {
		g.histo = histo.(*H1F)
	}

	if vers < 2 {
		g.min = float64(r.ReadF32())
		g.max = float64(r.ReadF32())
	} else {
		g.min = r.ReadF64()
		g.max = r.ReadF64()
	}

	r.CheckByteCount(pos, bcnt, beg, "TGraph")
	return r.Err()
}

type tgrapherrs struct {
	rvers int16
	tgraph

	xerr []float64
	yerr []float64
}

func (g *tgrapherrs) Class() string {
	return "TGraphErrors"
}

func (g *tgrapherrs) XError(i int) (float64, float64) {
	return g.xerr[i], g.xerr[i]
}

func (g *tgrapherrs) YError(i int) (float64, float64) {
	return g.yerr[i], g.yerr[i]
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgrapherrs) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, nil
	}

	pos := w.Pos()
	w.WriteVersion(g.rvers)

	if n, err := g.tgraph.MarshalROOT(w); err != nil {
		w.err = err
		return n, w.err
	}

	{
		w.WriteI8(0)
		w.WriteFastArrayF64(g.xerr)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.yerr)
	}

	return w.SetByteCount(pos, "TGraphErrors")
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgrapherrs) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	g.rvers = vers

	err := g.tgraph.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	if vers < 2 {
		_ = r.ReadI8()
		xerrs := r.ReadFastArrayF32(int(g.tgraph.npoints))
		_ = r.ReadI8()
		yerrs := r.ReadFastArrayF32(int(g.tgraph.npoints))
		g.xerr = make([]float64, len(xerrs))
		g.yerr = make([]float64, len(yerrs))
		for i := range xerrs {
			g.xerr[i] = float64(xerrs[i])
			g.yerr[i] = float64(yerrs[i])
		}

	} else {
		_ = r.ReadI8()
		g.xerr = r.ReadFastArrayF64(int(g.tgraph.npoints))
		_ = r.ReadI8()
		g.yerr = r.ReadFastArrayF64(int(g.tgraph.npoints))
	}
	r.CheckByteCount(pos, bcnt, beg, "TGraphErrors")
	return r.Err()
}

type tgraphasymmerrs struct {
	rvers int16
	tgraph

	xerrlo []float64
	xerrhi []float64
	yerrlo []float64
	yerrhi []float64
}

func (g *tgraphasymmerrs) Class() string {
	return "TGraphAsymmErrors"
}

func (g *tgraphasymmerrs) XError(i int) (float64, float64) {
	return g.xerrlo[i], g.xerrhi[i]
}

func (g *tgraphasymmerrs) YError(i int) (float64, float64) {
	return g.yerrlo[i], g.yerrhi[i]
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgraphasymmerrs) MarshalROOT(w *WBuffer) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	pos := w.Pos()
	w.WriteVersion(g.rvers)

	if n, err := g.tgraph.MarshalROOT(w); err != nil {
		w.err = err
		return n, w.err
	}

	{
		w.WriteI8(0)
		w.WriteFastArrayF64(g.xerrlo)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.xerrhi)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.yerrlo)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.yerrhi)
	}

	return w.SetByteCount(pos, "TGraphAsymmErrors")
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraphasymmerrs) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()
	g.rvers = vers

	err := g.tgraph.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	n := int(g.tgraph.npoints)
	switch {
	case vers < 2:
		// up to version 2, order is: xlo,ylo,xhi,yhi
		_ = r.ReadI8()
		xerrlo := r.ReadFastArrayF32(n)
		_ = r.ReadI8()
		yerrlo := r.ReadFastArrayF32(n)
		_ = r.ReadI8()
		xerrhi := r.ReadFastArrayF32(n)
		_ = r.ReadI8()
		yerrhi := r.ReadFastArrayF32(n)
		g.xerrlo = make([]float64, n)
		g.xerrhi = make([]float64, n)
		g.yerrlo = make([]float64, n)
		g.yerrhi = make([]float64, n)
		for i := range xerrlo {
			g.xerrlo[i] = float64(xerrlo[i])
			g.xerrhi[i] = float64(xerrhi[i])
			g.yerrlo[i] = float64(yerrlo[i])
			g.yerrhi[i] = float64(yerrhi[i])
		}
	case vers == 2:
		// version 2, order is: xlo,ylo,xhi,yhi (but in float64)
		_ = r.ReadI8()
		g.xerrlo = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.yerrlo = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.xerrhi = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.yerrhi = r.ReadFastArrayF64(n)
	default:
		// version 3 and higher: xlo,xhi,ylo,yhi
		// ie: the order of the fields in the TGraphAsymmErrors class.
		_ = r.ReadI8()
		g.xerrlo = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.xerrhi = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.yerrlo = r.ReadFastArrayF64(n)
		_ = r.ReadI8()
		g.yerrhi = r.ReadFastArrayF64(n)
	}
	r.CheckByteCount(pos, bcnt, beg, "TGraphAsymmErrors")
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := &tgraph{}
			return reflect.ValueOf(o)
		}
		Factory.add("TGraph", f)
		Factory.add("*rootio.tgraph", f)
	}
	{
		f := func() reflect.Value {
			o := &tgrapherrs{}
			return reflect.ValueOf(o)
		}
		Factory.add("TGraphErrors", f)
		Factory.add("*rootio.tgrapherrs", f)
	}
	{
		f := func() reflect.Value {
			o := &tgraphasymmerrs{}
			return reflect.ValueOf(o)
		}
		Factory.add("TGraphAsymmErrors", f)
		Factory.add("*rootio.tgraphasymmerrs", f)
	}
}

var (
	_ Object          = (*tgraph)(nil)
	_ Named           = (*tgraph)(nil)
	_ Graph           = (*tgraph)(nil)
	_ ROOTMarshaler   = (*tgraph)(nil)
	_ ROOTUnmarshaler = (*tgraph)(nil)

	_ Object          = (*tgrapherrs)(nil)
	_ Named           = (*tgrapherrs)(nil)
	_ Graph           = (*tgrapherrs)(nil)
	_ GraphErrors     = (*tgrapherrs)(nil)
	_ ROOTMarshaler   = (*tgrapherrs)(nil)
	_ ROOTUnmarshaler = (*tgrapherrs)(nil)

	_ Object          = (*tgraphasymmerrs)(nil)
	_ Named           = (*tgraphasymmerrs)(nil)
	_ Graph           = (*tgraphasymmerrs)(nil)
	_ GraphErrors     = (*tgraphasymmerrs)(nil)
	_ ROOTMarshaler   = (*tgraphasymmerrs)(nil)
	_ ROOTUnmarshaler = (*tgraphasymmerrs)(nil)
)
