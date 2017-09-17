// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

type tgraph struct {
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

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraph) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()

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

	r.ReadI32(&g.npoints)
	g.maxsize = g.npoints
	g.x = make([]float64, g.npoints)
	g.y = make([]float64, g.npoints)
	if vers < 2 {
		var i8 int8
		r.ReadI8(&i8)
		xs := make([]float32, g.npoints)
		r.ReadFastArrayF32(xs)
		r.ReadI8(&i8)
		ys := make([]float32, g.npoints)
		r.ReadFastArrayF32(ys)
		for i := range xs {
			g.x[i] = float64(xs[i])
			g.y[i] = float64(ys[i])
		}

	} else {
		var i8 int8
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.x)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.y)
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
		var f32 float32
		r.ReadF32(&f32)
		g.min = float64(f32)
		r.ReadF32(&f32)
		g.max = float64(f32)
	} else {
		r.ReadF64(&g.min)
		r.ReadF64(&g.max)
	}

	r.CheckByteCount(pos, bcnt, beg, "TGraph")
	return r.Err()
}

type tgrapherrs struct {
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

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgrapherrs) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()

	err := g.tgraph.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	g.xerr = make([]float64, g.tgraph.npoints)
	g.yerr = make([]float64, g.tgraph.npoints)
	if vers < 2 {
		var i8 int8
		xerrs := make([]float32, g.tgraph.npoints)
		yerrs := make([]float32, g.tgraph.npoints)
		r.ReadI8(&i8)
		r.ReadFastArrayF32(xerrs)
		r.ReadI8(&i8)
		r.ReadFastArrayF32(yerrs)
		for i := range xerrs {
			g.xerr[i] = float64(xerrs[i])
			g.yerr[i] = float64(yerrs[i])
		}

	} else {
		var i8 int8
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.xerr)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.yerr)
	}
	r.CheckByteCount(pos, bcnt, beg, "TGraphErrors")
	return r.Err()
}

type tgraphasymmerrs struct {
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

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraphasymmerrs) UnmarshalROOT(r *RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion()

	err := g.tgraph.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	n := int(g.tgraph.npoints)
	g.xerrlo = make([]float64, n)
	g.xerrhi = make([]float64, n)
	g.yerrlo = make([]float64, n)
	g.yerrhi = make([]float64, n)
	switch {
	case vers < 2:
		var i8 int8
		xerrlo := make([]float32, n)
		xerrhi := make([]float32, n)
		yerrlo := make([]float32, n)
		yerrhi := make([]float32, n)
		// up to version 2, order is: xlo,ylo,xhi,yhi
		r.ReadI8(&i8)
		r.ReadFastArrayF32(xerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF32(yerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF32(xerrhi)
		r.ReadI8(&i8)
		r.ReadFastArrayF32(yerrhi)
		for i := range xerrlo {
			g.xerrlo[i] = float64(xerrlo[i])
			g.xerrhi[i] = float64(xerrhi[i])
			g.yerrlo[i] = float64(yerrlo[i])
			g.yerrhi[i] = float64(yerrhi[i])
		}
	case vers == 2:
		var i8 int8
		// version 2, order is: xlo,ylo,xhi,yhi (but in float64)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.xerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.yerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.xerrhi)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.yerrhi)
	default:
		var i8 int8
		// version 3 and higher: xlo,xhi,ylo,yhi
		// ie: the order of the fields in the TGraphAsymmErrors class.
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.xerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.xerrhi)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.yerrlo)
		r.ReadI8(&i8)
		r.ReadFastArrayF64(g.yerrhi)
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

var _ Object = (*tgraph)(nil)
var _ Named = (*tgraph)(nil)
var _ Graph = (*tgraph)(nil)
var _ ROOTUnmarshaler = (*tgraph)(nil)

var _ Object = (*tgrapherrs)(nil)
var _ Named = (*tgrapherrs)(nil)
var _ Graph = (*tgrapherrs)(nil)
var _ GraphErrors = (*tgrapherrs)(nil)
var _ ROOTUnmarshaler = (*tgrapherrs)(nil)

var _ Object = (*tgraphasymmerrs)(nil)
var _ Named = (*tgraphasymmerrs)(nil)
var _ Graph = (*tgraphasymmerrs)(nil)
var _ GraphErrors = (*tgraphasymmerrs)(nil)
var _ ROOTUnmarshaler = (*tgraphasymmerrs)(nil)
