// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist

import (
	"fmt"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
	"go-hep.org/x/hep/hbook"
)

type tgraph struct {
	named rbase.Named

	maxsize int32
	npoints int32
	x       []float64
	y       []float64
	funcs   root.List
	histo   *H1F
	min     float64
	max     float64
}

func newGraph(n int) *tgraph {
	return &tgraph{
		named:   *rbase.NewNamed("", ""),
		maxsize: int32(n),
		npoints: int32(n),
		x:       make([]float64, n),
		y:       make([]float64, n),
		funcs:   rcont.NewList("", nil),
	}
}

// NewGraphFrom creates a new Graph from 2-dim hbook data points.
func NewGraphFrom(s2 *hbook.S2D) Graph {
	var (
		n     = s2.Len()
		groot = newGraph(n)
		ymin  = +math.MaxFloat64
		ymax  = -math.MaxFloat64
	)

	for i, pt := range s2.Points() {
		groot.x[i] = pt.X
		groot.y[i] = pt.Y

		ymax = math.Max(ymax, pt.Y)
		ymin = math.Min(ymin, pt.Y)
	}

	groot.named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.named.SetTitle(v.(string))
	}

	groot.min = ymin
	groot.max = ymax

	return groot
}

func (*tgraph) RVersion() int16 {
	return rvers.Graph
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

func (g *tgraph) ROOTMerge(src root.Object) error {
	switch src := src.(type) {
	case *tgraph:
		if src.maxsize > g.maxsize {
			g.maxsize = src.maxsize
		}
		g.npoints += src.npoints
		g.x = append(g.x, src.x...)
		g.y = append(g.y, src.y...)
		g.min = math.Min(g.min, src.min)
		g.max = math.Max(g.max, src.max)
		// FIXME(sbinet): handle g.funcs
		// FIXME(sbinet): handle g.histo
		// FIXME(sbinet): re-sort x,y,... slices according to x.
		return nil
	default:
		return fmt.Errorf("rhist: can not merge %T into %T", src, g)
	}
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgraph) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(g.RVersion())

	for _, a := range []rbytes.Marshaler{
		&g.named,
		&rbase.AttLine{},
		&rbase.AttFill{},
		&rbase.AttMarker{},
	} {
		if _, err := a.MarshalROOT(w); err != nil {
			return 0, err
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

	return w.SetByteCount(pos, g.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraph) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(g.Class())

	for _, a := range []rbytes.Unmarshaler{
		&g.named,
		&rbase.AttLine{},
		&rbase.AttFill{},
		&rbase.AttMarker{},
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
		xs := make([]float32, g.npoints)
		r.ReadArrayF32(xs)
		_ = r.ReadI8()
		ys := make([]float32, g.npoints)
		r.ReadArrayF32(ys)
		g.x = make([]float64, len(xs))
		g.y = make([]float64, len(ys))
		for i := range xs {
			g.x[i] = float64(xs[i])
			g.y[i] = float64(ys[i])
		}

	} else {
		_ = r.ReadI8()
		g.x = make([]float64, g.npoints)
		r.ReadArrayF64(g.x)
		_ = r.ReadI8()
		g.y = make([]float64, g.npoints)
		r.ReadArrayF64(g.y)
	}

	funcs := r.ReadObjectAny()
	if funcs != nil {
		g.funcs = funcs.(root.List)
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

	r.CheckByteCount(pos, bcnt, beg, g.Class())
	return r.Err()
}

type tgrapherrs struct {
	tgraph

	xerr []float64
	yerr []float64
}

func newGraphErrs(n int) *tgrapherrs {
	return &tgrapherrs{
		tgraph: *newGraph(n),
		xerr:   make([]float64, n),
		yerr:   make([]float64, n),
	}
}

// NewGraphErrorsFrom creates a new GraphErrors from 2-dim hbook data points.
func NewGraphErrorsFrom(s2 *hbook.S2D) GraphErrors {
	var (
		n     = s2.Len()
		groot = newGraphErrs(n)
		ymin  = +math.MaxFloat64
		ymax  = -math.MaxFloat64
	)

	for i, pt := range s2.Points() {
		groot.x[i] = pt.X
		groot.xerr[i] = pt.ErrX.Min
		groot.y[i] = pt.Y
		groot.yerr[i] = pt.ErrY.Min

		ymax = math.Max(ymax, pt.Y)
		ymin = math.Min(ymin, pt.Y)
	}

	groot.tgraph.named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.tgraph.named.SetTitle(v.(string))
	}

	groot.min = ymin
	groot.max = ymax

	return groot
}

func (*tgrapherrs) RVersion() int16 {
	return rvers.GraphErrors
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

func (g *tgrapherrs) ROOTMerge(src root.Object) error {
	switch src := src.(type) {
	case *tgrapherrs:
		err := g.tgraph.ROOTMerge(&src.tgraph)
		if err != nil {
			return fmt.Errorf("rhist: could not merge %q: %w", src.Name(), err)
		}
		g.xerr = append(g.xerr, src.xerr...)
		g.yerr = append(g.yerr, src.yerr...)
		// FIXME(sbinet): re-sort x,y,... slices according to x.
		return nil
	default:
		return fmt.Errorf("rhist: can not merge %T into %T", src, g)
	}
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgrapherrs) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, nil
	}

	pos := w.WriteVersion(g.RVersion())

	if n, err := g.tgraph.MarshalROOT(w); err != nil {
		return n, err
	}

	{
		w.WriteI8(0)
		w.WriteFastArrayF64(g.xerr)
		w.WriteI8(0)
		w.WriteFastArrayF64(g.yerr)
	}

	return w.SetByteCount(pos, g.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgrapherrs) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(g.Class())

	err := g.tgraph.UnmarshalROOT(r)
	if err != nil {
		return err
	}

	if vers < 2 {
		_ = r.ReadI8()
		xerrs := make([]float32, g.tgraph.npoints)
		r.ReadArrayF32(xerrs)
		_ = r.ReadI8()
		yerrs := make([]float32, g.tgraph.npoints)
		r.ReadArrayF32(yerrs)
		g.xerr = make([]float64, len(xerrs))
		g.yerr = make([]float64, len(yerrs))
		for i := range xerrs {
			g.xerr[i] = float64(xerrs[i])
			g.yerr[i] = float64(yerrs[i])
		}

	} else {
		_ = r.ReadI8()
		g.xerr = make([]float64, g.tgraph.npoints)
		r.ReadArrayF64(g.xerr)
		_ = r.ReadI8()
		g.yerr = make([]float64, g.tgraph.npoints)
		r.ReadArrayF64(g.yerr)
	}
	r.CheckByteCount(pos, bcnt, beg, g.Class())
	return r.Err()
}

type tgraphasymmerrs struct {
	tgraph

	xerrlo []float64
	xerrhi []float64
	yerrlo []float64
	yerrhi []float64
}

func newGraphAsymmErrs(n int) *tgraphasymmerrs {
	return &tgraphasymmerrs{
		tgraph: *newGraph(n),
		xerrlo: make([]float64, n),
		xerrhi: make([]float64, n),
		yerrlo: make([]float64, n),
		yerrhi: make([]float64, n),
	}
}

// NewGraphAsymmErrorsFrom creates a new GraphAsymErrors from 2-dim hbook data points.
func NewGraphAsymmErrorsFrom(s2 *hbook.S2D) GraphErrors {
	var (
		n     = s2.Len()
		groot = newGraphAsymmErrs(n)
		ymin  = +math.MaxFloat64
		ymax  = -math.MaxFloat64
	)

	for i, pt := range s2.Points() {
		groot.x[i] = pt.X
		groot.xerrlo[i] = pt.ErrX.Min
		groot.xerrhi[i] = pt.ErrX.Max
		groot.y[i] = pt.Y
		groot.yerrlo[i] = pt.ErrY.Min
		groot.yerrhi[i] = pt.ErrY.Max

		ymax = math.Max(ymax, pt.Y)
		ymin = math.Min(ymin, pt.Y)
	}

	groot.tgraph.named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.tgraph.named.SetTitle(v.(string))
	}

	groot.min = ymin
	groot.max = ymax

	return groot
}

func (*tgraphasymmerrs) RVersion() int16 {
	return rvers.GraphAsymmErrors
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

func (g *tgraphasymmerrs) ROOTMerge(src root.Object) error {
	switch src := src.(type) {
	case *tgraphasymmerrs:
		err := g.tgraph.ROOTMerge(&src.tgraph)
		if err != nil {
			return fmt.Errorf("rhist: could not merge %q: %w", src.Name(), err)
		}
		g.xerrlo = append(g.xerrlo, src.xerrlo...)
		g.xerrhi = append(g.xerrhi, src.xerrhi...)
		g.yerrlo = append(g.yerrlo, src.yerrlo...)
		g.yerrhi = append(g.yerrhi, src.yerrhi...)
		// FIXME(sbinet): re-sort x,y,... slices according to x.
		return nil
	default:
		return fmt.Errorf("rhist: can not merge %T into %T", src, g)
	}
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself to a ROOT buffer
func (g *tgraphasymmerrs) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(g.RVersion())

	if n, err := g.tgraph.MarshalROOT(w); err != nil {
		return n, err
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

	return w.SetByteCount(pos, g.Class())
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraphasymmerrs) UnmarshalROOT(r *rbytes.RBuffer) error {
	beg := r.Pos()

	vers, pos, bcnt := r.ReadVersion(g.Class())

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
		// up to version 2, order is: xlo,ylo,xhi,yhi
		xerrlo := make([]float32, n)
		yerrlo := make([]float32, n)
		xerrhi := make([]float32, n)
		yerrhi := make([]float32, n)
		_ = r.ReadI8()
		r.ReadArrayF32(xerrlo)
		_ = r.ReadI8()
		r.ReadArrayF32(yerrlo)
		_ = r.ReadI8()
		r.ReadArrayF32(xerrhi)
		_ = r.ReadI8()
		r.ReadArrayF32(yerrhi)
		for i := range xerrlo {
			g.xerrlo[i] = float64(xerrlo[i])
			g.xerrhi[i] = float64(xerrhi[i])
			g.yerrlo[i] = float64(yerrlo[i])
			g.yerrhi[i] = float64(yerrhi[i])
		}
	case vers == 2:
		// version 2, order is: xlo,ylo,xhi,yhi (but in float64)
		_ = r.ReadI8()
		r.ReadArrayF64(g.xerrlo)
		_ = r.ReadI8()
		r.ReadArrayF64(g.yerrlo)
		_ = r.ReadI8()
		r.ReadArrayF64(g.xerrhi)
		_ = r.ReadI8()
		r.ReadArrayF64(g.yerrhi)
	default:
		// version 3 and higher: xlo,xhi,ylo,yhi
		// ie: the order of the fields in the TGraphAsymmErrors class.
		_ = r.ReadI8()
		r.ReadArrayF64(g.xerrlo)
		_ = r.ReadI8()
		r.ReadArrayF64(g.xerrhi)
		_ = r.ReadI8()
		r.ReadArrayF64(g.yerrlo)
		_ = r.ReadI8()
		r.ReadArrayF64(g.yerrhi)
	}
	r.CheckByteCount(pos, bcnt, beg, g.Class())
	return r.Err()
}

func init() {
	{
		f := func() reflect.Value {
			o := newGraph(0)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TGraph", f)
	}
	{
		f := func() reflect.Value {
			o := newGraphErrs(0)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TGraphErrors", f)
	}
	{
		f := func() reflect.Value {
			o := newGraphAsymmErrs(0)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TGraphAsymmErrors", f)
	}
}

var (
	_ root.Object        = (*tgraph)(nil)
	_ root.Named         = (*tgraph)(nil)
	_ root.Merger        = (*tgraph)(nil)
	_ Graph              = (*tgraph)(nil)
	_ rbytes.Marshaler   = (*tgraph)(nil)
	_ rbytes.Unmarshaler = (*tgraph)(nil)

	_ root.Object        = (*tgrapherrs)(nil)
	_ root.Named         = (*tgrapherrs)(nil)
	_ root.Merger        = (*tgrapherrs)(nil)
	_ Graph              = (*tgrapherrs)(nil)
	_ GraphErrors        = (*tgrapherrs)(nil)
	_ rbytes.Marshaler   = (*tgrapherrs)(nil)
	_ rbytes.Unmarshaler = (*tgrapherrs)(nil)

	_ root.Object        = (*tgraphasymmerrs)(nil)
	_ root.Named         = (*tgraphasymmerrs)(nil)
	_ root.Merger        = (*tgraphasymmerrs)(nil)
	_ Graph              = (*tgraphasymmerrs)(nil)
	_ GraphErrors        = (*tgraphasymmerrs)(nil)
	_ rbytes.Marshaler   = (*tgraphasymmerrs)(nil)
	_ rbytes.Unmarshaler = (*tgraphasymmerrs)(nil)
)
