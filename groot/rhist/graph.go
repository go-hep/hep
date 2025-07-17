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
	"go-hep.org/x/hep/hbook/yodacnv"
)

type tgraph struct {
	rbase.Named
	attline   rbase.AttLine
	attfill   rbase.AttFill
	attmarker rbase.AttMarker

	maxsize int32
	npoints int32
	x       []float64
	y       []float64
	funcs   root.List
	histo   *H1F
	min     float64
	max     float64
	opt     string
}

func newGraph(n int) *tgraph {
	return &tgraph{
		Named:     *rbase.NewNamed("", ""),
		attline:   *rbase.NewAttLine(),
		attfill:   *rbase.NewAttFill(),
		attmarker: *rbase.NewAttMarker(),
		maxsize:   int32(n),
		npoints:   int32(n),
		x:         make([]float64, n),
		y:         make([]float64, n),
		funcs:     rcont.NewList("", nil),
		opt:       "",
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

	groot.Named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.Named.SetTitle(v.(string))
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

	hdr := w.WriteHeader(g.Class(), g.RVersion())

	w.WriteObject(&g.Named)
	w.WriteObject(&g.attline)
	w.WriteObject(&g.attfill)
	w.WriteObject(&g.attmarker)

	w.WriteI32(g.npoints)
	{
		w.WriteI8(1)
		w.WriteArrayF64(g.x)
		w.WriteI8(1)
		w.WriteArrayF64(g.y)
	}

	w.WriteObjectAny(g.funcs)
	w.WriteObjectAny(g.histo)
	{
		w.WriteF64(g.min)
		w.WriteF64(g.max)
	}
	w.WriteString(g.opt)

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraph) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(g.Class(), g.RVersion())

	r.ReadObject(&g.Named)
	r.ReadObject(&g.attline)
	r.ReadObject(&g.attfill)
	r.ReadObject(&g.attmarker)

	g.npoints = r.ReadI32()
	g.maxsize = g.npoints
	if hdr.Vers < 2 {
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

	if hdr.Vers < 2 {
		g.min = float64(r.ReadF32())
		g.max = float64(r.ReadF32())
	} else {
		g.min = r.ReadF64()
		g.max = r.ReadF64()
	}
	if hdr.Vers > 4 {
		g.opt = r.ReadString()
	}

	r.CheckHeader(hdr)
	return r.Err()
}

func (g *tgraph) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, g.Named.RMembers()...)
	mbrs = append(mbrs, g.attline.RMembers()...)
	mbrs = append(mbrs, g.attfill.RMembers()...)
	mbrs = append(mbrs, g.attmarker.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNpoints", Value: &g.npoints},
		{Name: "fX", Value: &g.x},
		{Name: "fY", Value: &g.y},
		{Name: "fFunctions", Value: g.funcs},
		{Name: "fHistogram", Value: &g.histo},
		{Name: "fMinimum", Value: &g.min},
		{Name: "fMaximum", Value: &g.max},
		{Name: "fOption", Value: &g.opt},
	}...)

	return mbrs
}

// MarshalYODA implements the YODAMarshaler interface.
func (g *tgraph) MarshalYODA() ([]byte, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}

	s2d := hbook.NewS2D(pts...)
	s2d.Annotation()["name"] = g.Name()
	s2d.Annotation()["title"] = g.Title()
	return s2d.MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (g *tgraph) UnmarshalYODA(raw []byte) error {
	var gg hbook.S2D
	err := gg.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*g = *NewGraphFrom(&gg).(*tgraph)
	return nil
}

// Keys implements the ObjectFinder interface.
func (g *tgraph) Keys() []string {
	var keys []string
	for i := range g.funcs.Len() {
		o, ok := g.funcs.At(i).(root.Named)
		if !ok {
			continue
		}
		keys = append(keys, o.Name())
	}
	return keys
}

// Get implements the ObjectFinder interface.
func (g *tgraph) Get(name string) (root.Object, error) {
	for i := range g.funcs.Len() {
		o, ok := g.funcs.At(i).(root.Named)
		if !ok {
			continue
		}
		if o.Name() == name {
			return g.funcs.At(i), nil
		}
	}

	return nil, fmt.Errorf("no object named %q", name)
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

	groot.tgraph.Named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.tgraph.Named.SetTitle(v.(string))
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

	hdr := w.WriteHeader(g.Class(), g.RVersion())

	w.WriteObject(&g.tgraph)
	{
		w.WriteI8(1)
		w.WriteArrayF64(g.xerr)
		w.WriteI8(1)
		w.WriteArrayF64(g.yerr)
	}

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgrapherrs) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(g.Class(), g.RVersion())

	r.ReadObject(&g.tgraph)

	if hdr.Vers < 2 {
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

	r.CheckHeader(hdr)
	return r.Err()
}

func (g *tgrapherrs) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, g.tgraph.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fEX", Value: &g.xerr},
		{Name: "fEY", Value: &g.yerr},
	}...)

	return mbrs
}

// MarshalYODA implements the YODAMarshaler interface.
func (g *tgrapherrs) MarshalYODA() ([]byte, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}
	for i := range pts {
		xlo, xhi := g.XError(i)
		ylo, yhi := g.YError(i)
		pt := &pts[i]
		pt.ErrX = hbook.Range{Min: xlo, Max: xhi}
		pt.ErrY = hbook.Range{Min: ylo, Max: yhi}
	}

	s2d := hbook.NewS2D(pts...)
	s2d.Annotation()["name"] = g.Name()
	s2d.Annotation()["title"] = g.Title()
	return s2d.MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (g *tgrapherrs) UnmarshalYODA(raw []byte) error {
	var gg hbook.S2D
	err := gg.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*g = *NewGraphErrorsFrom(&gg).(*tgrapherrs)
	return nil
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

	groot.tgraph.Named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.tgraph.Named.SetTitle(v.(string))
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

	hdr := w.WriteHeader(g.Class(), g.RVersion())

	w.WriteObject(&g.tgraph)
	{
		w.WriteI8(1)
		w.WriteArrayF64(g.xerrlo)
		w.WriteI8(1)
		w.WriteArrayF64(g.xerrhi)
		w.WriteI8(1)
		w.WriteArrayF64(g.yerrlo)
		w.WriteI8(1)
		w.WriteArrayF64(g.yerrhi)
	}

	return w.SetHeader(hdr)
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (g *tgraphasymmerrs) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(g.Class(), g.RVersion())

	r.ReadObject(&g.tgraph)

	n := int(g.tgraph.npoints)
	g.xerrlo = make([]float64, n)
	g.xerrhi = make([]float64, n)
	g.yerrlo = make([]float64, n)
	g.yerrhi = make([]float64, n)
	switch {
	case hdr.Vers < 2:
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
	case hdr.Vers == 2:
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

	r.CheckHeader(hdr)
	return r.Err()
}

func (g *tgraphasymmerrs) RMembers() (mbrs []rbytes.Member) {
	mbrs = append(mbrs, g.tgraph.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fEXlow", Value: &g.xerrlo},
		{Name: "fEXhigh", Value: &g.xerrhi},
		{Name: "fEYlow", Value: &g.yerrlo},
		{Name: "fEYhigh", Value: &g.yerrhi},
	}...)

	return mbrs
}

// MarshalYODA implements the YODAMarshaler interface.
func (g *tgraphasymmerrs) MarshalYODA() ([]byte, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}
	for i := range pts {
		xlo, xhi := g.XError(i)
		ylo, yhi := g.YError(i)
		pt := &pts[i]
		pt.ErrX = hbook.Range{Min: xlo, Max: xhi}
		pt.ErrY = hbook.Range{Min: ylo, Max: yhi}
	}

	s2d := hbook.NewS2D(pts...)
	s2d.Annotation()["name"] = g.Name()
	s2d.Annotation()["title"] = g.Title()
	return s2d.MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (g *tgraphasymmerrs) UnmarshalYODA(raw []byte) error {
	var gg hbook.S2D
	err := gg.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*g = *NewGraphAsymmErrorsFrom(&gg).(*tgraphasymmerrs)
	return nil
}

// tgraphmultierrs is a graph with asymmetric error bars and
// multiple y error dimensions.
type tgraphmultierrs struct {
	tgraph

	nyerr      int32           // The amount of different y-errors
	sumErrMode int32           // How y errors are summed: kOnlyFirst = Only First; kSquareSum = Squared Sum; kSum =
	xerrlo     []float64       // array of X low errors
	xerrhi     []float64       // array of X high errors
	yerrlo     []rcont.ArrayD  // two dimensional array of Y low errors
	yerrhi     []rcont.ArrayD  // two dimensional array of Y high errors
	attfills   []rbase.AttFill // the AttFill attributes of the different errors
	attlines   []rbase.AttLine // the AttLine attributes of the different errors
}

// NewGraphMultiErrorsFrom creates a new GraphMultiErrors from 2-dim hbook data points.
func NewGraphMultiErrorsFrom(s2 *hbook.S2D) GraphErrors {
	return newGraphMultiErrorsFrom(s2)
}

func newGraphMultiErrs(n, ny int) *tgraphmultierrs {
	g := &tgraphmultierrs{
		tgraph:   *newGraph(n),
		nyerr:    int32(ny),
		xerrlo:   make([]float64, n),
		xerrhi:   make([]float64, n),
		yerrlo:   make([]rcont.ArrayD, ny),
		yerrhi:   make([]rcont.ArrayD, ny),
		attfills: make([]rbase.AttFill, n),
		attlines: make([]rbase.AttLine, n),
	}
	for i := range ny {
		g.yerrlo[i].Data = make([]float64, n)
		g.yerrhi[i].Data = make([]float64, n)
	}
	return g
}

func newGraphMultiErrorsFrom(s2 *hbook.S2D) GraphErrors {
	var (
		n     = s2.Len()
		groot = newGraphMultiErrs(n, 1)
		ymin  = +math.MaxFloat64
		ymax  = -math.MaxFloat64
	)
	for i, pt := range s2.Points() {
		groot.x[i] = pt.X
		groot.xerrlo[i] = pt.ErrX.Min
		groot.xerrhi[i] = pt.ErrX.Max
		groot.y[i] = pt.Y
		groot.yerrlo[0].Data[i] = pt.ErrY.Min
		groot.yerrhi[0].Data[i] = pt.ErrY.Max

		ymax = math.Max(ymax, pt.Y)
		ymin = math.Min(ymin, pt.Y)
	}

	groot.tgraph.Named.SetName(s2.Name())
	if v, ok := s2.Annotation()["title"]; ok {
		groot.tgraph.Named.SetTitle(v.(string))
	}

	groot.min = ymin
	groot.max = ymax

	return groot
}

func (*tgraphmultierrs) Class() string {
	return "TGraphMultiErrors"
}

func (*tgraphmultierrs) RVersion() int16 {
	return rvers.GraphMultiErrors
}

func (g *tgraphmultierrs) XError(i int) (float64, float64) {
	return g.xerrlo[i], g.xerrhi[i]
}

func (g *tgraphmultierrs) YError(i int) (float64, float64) {
	return g.yerrlo[0].At(i), g.yerrhi[0].At(i)
}

// MarshalROOT implements rbytes.Marshaler
func (g *tgraphmultierrs) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	hdr := w.WriteHeader(g.Class(), g.RVersion())

	w.WriteObject(&g.tgraph)
	w.WriteI32(g.nyerr)
	w.WriteI32(g.sumErrMode)
	w.WriteI8(1) // is-array
	w.WriteArrayF64(g.xerrlo[:g.tgraph.npoints])
	w.WriteI8(1) // is-array
	w.WriteArrayF64(g.xerrhi[:g.tgraph.npoints])
	writeStdVectorTArrayD(w, g.yerrlo)
	writeStdVectorTArrayD(w, g.yerrhi)
	writeStdVectorTAttFill(w, g.attfills)
	writeStdVectorTAttLine(w, g.attlines)
	return w.SetHeader(hdr)
}

// UnmarshalROOT implements rbytes.Unmarshaler
func (g *tgraphmultierrs) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	hdr := r.ReadHeader(g.Class(), g.RVersion())

	r.ReadObject(&g.tgraph)

	g.nyerr = r.ReadI32()
	g.sumErrMode = r.ReadI32()
	_ = r.ReadI8() // is-array
	g.xerrlo = rbytes.ResizeF64(nil, int(g.tgraph.npoints))
	r.ReadArrayF64(g.xerrlo)
	_ = r.ReadI8() // is-array
	g.xerrhi = rbytes.ResizeF64(nil, int(g.tgraph.npoints))
	r.ReadArrayF64(g.xerrhi)
	readStdVectorTArrayD(r, &g.yerrlo)
	readStdVectorTArrayD(r, &g.yerrhi)
	readStdVectorTAttFill(r, &g.attfills)
	readStdVectorTAttLine(r, &g.attlines)

	r.CheckHeader(hdr)
	return r.Err()
}

func (g *tgraphmultierrs) RMembers() (mbrs []rbytes.Member) {
	var (
		yerrlo = make([][]float64, len(g.yerrlo))
		yerrhi = make([][]float64, len(g.yerrhi))
	)
	for i, v := range g.yerrlo {
		yerrlo[i] = v.Data
	}

	for i, v := range g.yerrhi {
		yerrhi[i] = v.Data
	}

	var (
		attfills = make([]*rbase.AttFill, len(g.attfills))
		attlines = make([]*rbase.AttLine, len(g.attlines))
	)
	for i := range g.attfills {
		attfills[i] = &g.attfills[i]
	}
	for i := range g.attlines {
		attlines[i] = &g.attlines[i]
	}
	mbrs = append(mbrs, g.tgraph.RMembers()...)
	mbrs = append(mbrs, []rbytes.Member{
		{Name: "fNYErrors", Value: &g.nyerr},
		{Name: "fSumErrorsMode", Value: &g.sumErrMode},
		{Name: "fExL", Value: &g.xerrlo},
		{Name: "fExH", Value: &g.xerrhi},
		{Name: "fEyL", Value: &yerrlo},
		{Name: "fEyH", Value: &yerrhi},
		{Name: "fAttFill", Value: attfills},
		{Name: "fAttLine", Value: attlines},
	}...)

	return mbrs
}

// MarshalYODA implements the YODAMarshaler interface.
func (g *tgraphmultierrs) MarshalYODA() ([]byte, error) {
	pts := make([]hbook.Point2D, g.Len())
	for i := range pts {
		x, y := g.XY(i)
		pts[i].X = x
		pts[i].Y = y
	}
	for i := range pts {
		xlo, xhi := g.XError(i)
		ylo, yhi := g.YError(i)
		pt := &pts[i]
		pt.ErrX = hbook.Range{Min: xlo, Max: xhi}
		pt.ErrY = hbook.Range{Min: ylo, Max: yhi}
	}

	// FIXME(sbinet): add a yoda-compatible representation
	// for multi-errors?
	s2d := hbook.NewS2D(pts...)
	s2d.Annotation()["name"] = g.Name()
	s2d.Annotation()["title"] = g.Title()
	return s2d.MarshalYODA()
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (g *tgraphmultierrs) UnmarshalYODA(raw []byte) error {
	var gg hbook.S2D
	err := gg.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*g = *newGraphMultiErrorsFrom(&gg).(*tgraphmultierrs)
	return nil
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
	{
		f := func() reflect.Value {
			o := newGraphMultiErrs(0, 0)
			return reflect.ValueOf(o)
		}
		rtypes.Factory.Add("TGraphMultiErrors", f)
	}
}

var (
	_ root.Object         = (*tgraph)(nil)
	_ root.Named          = (*tgraph)(nil)
	_ root.Merger         = (*tgraph)(nil)
	_ root.ObjectFinder   = (*tgraph)(nil)
	_ Graph               = (*tgraph)(nil)
	_ rbytes.Marshaler    = (*tgraph)(nil)
	_ rbytes.Unmarshaler  = (*tgraph)(nil)
	_ rbytes.RSlicer      = (*tgraph)(nil)
	_ yodacnv.Marshaler   = (*tgraph)(nil)
	_ yodacnv.Unmarshaler = (*tgraph)(nil)

	_ root.Object         = (*tgrapherrs)(nil)
	_ root.Named          = (*tgrapherrs)(nil)
	_ root.Merger         = (*tgrapherrs)(nil)
	_ root.ObjectFinder   = (*tgrapherrs)(nil)
	_ Graph               = (*tgrapherrs)(nil)
	_ GraphErrors         = (*tgrapherrs)(nil)
	_ rbytes.Marshaler    = (*tgrapherrs)(nil)
	_ rbytes.Unmarshaler  = (*tgrapherrs)(nil)
	_ rbytes.RSlicer      = (*tgrapherrs)(nil)
	_ yodacnv.Marshaler   = (*tgrapherrs)(nil)
	_ yodacnv.Unmarshaler = (*tgrapherrs)(nil)

	_ root.Object         = (*tgraphasymmerrs)(nil)
	_ root.Named          = (*tgraphasymmerrs)(nil)
	_ root.Merger         = (*tgraphasymmerrs)(nil)
	_ root.ObjectFinder   = (*tgraphasymmerrs)(nil)
	_ Graph               = (*tgraphasymmerrs)(nil)
	_ GraphErrors         = (*tgraphasymmerrs)(nil)
	_ rbytes.Marshaler    = (*tgraphasymmerrs)(nil)
	_ rbytes.Unmarshaler  = (*tgraphasymmerrs)(nil)
	_ rbytes.RSlicer      = (*tgraphasymmerrs)(nil)
	_ yodacnv.Marshaler   = (*tgraphasymmerrs)(nil)
	_ yodacnv.Unmarshaler = (*tgraphasymmerrs)(nil)

	_ root.Object         = (*tgraphmultierrs)(nil)
	_ root.Named          = (*tgraphmultierrs)(nil)
	_ root.Merger         = (*tgraphmultierrs)(nil)
	_ root.ObjectFinder   = (*tgraphmultierrs)(nil)
	_ Graph               = (*tgraphmultierrs)(nil)
	_ GraphErrors         = (*tgraphmultierrs)(nil)
	_ rbytes.Marshaler    = (*tgraphmultierrs)(nil)
	_ rbytes.Unmarshaler  = (*tgraphmultierrs)(nil)
	_ rbytes.RSlicer      = (*tgraphmultierrs)(nil)
	_ yodacnv.Marshaler   = (*tgraphmultierrs)(nil)
	_ yodacnv.Unmarshaler = (*tgraphmultierrs)(nil)
)

func writeStdVectorTArrayD(w *rbytes.WBuffer, vs []rcont.ArrayD) {
	if w.Err() != nil {
		return
	}
	const typename = "vector<TArrayD>"
	hdr := w.WriteHeader(typename, rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(vs)))
	for i := range vs {
		w.WriteObject(&vs[i])
	}
	_, _ = w.SetHeader(hdr)
}

func writeStdVectorTAttFill(w *rbytes.WBuffer, vs []rbase.AttFill) {
	if w.Err() != nil {
		return
	}
	const typename = "vector<TAttFill>"
	hdr := w.WriteHeader(typename, rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(vs)))
	for i := range vs {
		w.WriteObject(&vs[i])
	}
	_, _ = w.SetHeader(hdr)
}

func writeStdVectorTAttLine(w *rbytes.WBuffer, vs []rbase.AttLine) {
	if w.Err() != nil {
		return
	}
	const typename = "vector<TAttLine>"
	hdr := w.WriteHeader(typename, rvers.StreamerBaseSTL)
	w.WriteI32(int32(len(vs)))
	for i := range vs {
		w.WriteObject(&vs[i])
	}
	_, _ = w.SetHeader(hdr)
}

func readStdVectorTArrayD(r *rbytes.RBuffer, vs *[]rcont.ArrayD) {
	if r.Err() != nil {
		return
	}

	hdr := r.ReadHeader("vector<TArrayD>", rvers.StreamerBaseSTL)

	// FIXME(sbinet): use rbytes.Resize[T]
	n := int(r.ReadI32())
	if n == 0 {
		*vs = nil
		r.CheckHeader(hdr)
		return
	}
	*vs = make([]rcont.ArrayD, n)
	for i := range *vs {
		r.ReadObject(&(*vs)[i])
	}

	r.CheckHeader(hdr)
}

func readStdVectorTAttFill(r *rbytes.RBuffer, vs *[]rbase.AttFill) {
	if r.Err() != nil {
		return
	}

	hdr := r.ReadHeader("vector<TAttFill>", rvers.StreamerBaseSTL)
	if hdr.MemberWise {
		clvers := r.ReadI16()
		switch {
		case clvers == 1:
			// TODO
		case clvers <= 0:
			/*chksum*/ _ = r.ReadU32()
		}
	}

	// FIXME(sbinet): use rbytes.Resize[T]
	n := int(r.ReadI32())
	if n == 0 {
		*vs = nil
		r.CheckHeader(hdr)
		return
	}

	*vs = make([]rbase.AttFill, n)
	switch {
	case hdr.MemberWise:
		p := make([]int16, n)
		r.ReadArrayI16(p)
		for i := range *vs {
			(*vs)[i].Color = p[i]
		}
		r.ReadArrayI16(p)
		for i := range *vs {
			(*vs)[i].Style = p[i]
		}
	default:
		for i := range *vs {
			r.ReadObject(&(*vs)[i])
		}
	}

	r.CheckHeader(hdr)
}

func readStdVectorTAttLine(r *rbytes.RBuffer, vs *[]rbase.AttLine) {
	if r.Err() != nil {
		return
	}

	hdr := r.ReadHeader("vector<TAttLine>", rvers.StreamerBaseSTL)
	if hdr.MemberWise {
		clvers := r.ReadI16()
		switch {
		case clvers == 1:
			// TODO
		case clvers <= 0:
			/*chksum*/ _ = r.ReadU32()
		}
	}

	// FIXME(sbinet): use rbytes.Resize[T]
	n := int(r.ReadI32())
	if n == 0 {
		*vs = nil
		r.CheckHeader(hdr)
		return
	}
	*vs = make([]rbase.AttLine, n)
	switch {
	case hdr.MemberWise:
		p := make([]int16, n)
		r.ReadArrayI16(p)
		for i := range *vs {
			(*vs)[i].Color = p[i]
		}
		r.ReadArrayI16(p)
		for i := range *vs {
			(*vs)[i].Style = p[i]
		}
		r.ReadArrayI16(p)
		for i := range *vs {
			(*vs)[i].Width = p[i]
		}
	default:
		for i := range *vs {
			r.ReadObject(&(*vs)[i])
		}
	}

	r.CheckHeader(hdr)
}
