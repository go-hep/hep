// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// H2D is a 2-dim histogram with weighted entries.
type H2D struct {
	axis axis2D
	ann  Annotation
}

// NewH2D creates a new 2-dim histogram.
func NewH2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) *H2D {
	return &H2D{
		axis: newAxis2D(nx, xlow, xhigh, ny, ylow, yhigh),
		ann:  make(Annotation),
	}
}

// Name returns the name of this histogram, if any
func (h *H2D) Name() string {
	v, ok := h.ann["name"]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

// Annotation returns the annotations attached to this histogram
func (h *H2D) Annotation() Annotation {
	return h.ann
}

// Rank returns the number of dimensions for this histogram
func (h *H2D) Rank() int {
	return 2
}

// Entries returns the number of entries in this histogram
func (h *H2D) Entries() int64 {
	return h.axis.entries()
}

// EffEntries returns the number of effective entries in this histogram
func (h *H2D) EffEntries() float64 {
	return h.axis.effEntries()
}

// SumW returns the sum of weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW() float64 {
	return h.axis.dist.SumW()
}

// SumW2 returns the sum of squared weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW2() float64 {
	return h.axis.dist.SumW2()
}

// MeanX returns the mean X.
// Overflows are included in the computation.
func (h *H2D) MeanX() float64 {
	return h.axis.dist.meanX()
}

// MeanY returns the mean Y.
// Overflows are included in the computation.
func (h *H2D) MeanY() float64 {
	return h.axis.dist.meanY()
}

// VarianceX returns the variance in X.
// Overflows are included in the computation.
func (h *H2D) VarianceX() float64 {
	return h.axis.dist.varianceX()
}

// VarianceY returns the variance in Y.
// Overflows are included in the computation.
func (h *H2D) VarianceY() float64 {
	return h.axis.dist.varianceY()
}

// StdDevX returns the standard deviation in X.
// Overflows are included in the computation.
func (h *H2D) StdDevX() float64 {
	return h.axis.dist.stdDevX()
}

// StdDevY returns the standard deviation in Y.
// Overflows are included in the computation.
func (h *H2D) StdDevY() float64 {
	return h.axis.dist.stdDevY()
}

// StdErrX returns the standard error in X.
// Overflows are included in the computation.
func (h *H2D) StdErrX() float64 {
	return h.axis.dist.stdErrX()
}

// StdErrY returns the standard error in Y.
// Overflows are included in the computation.
func (h *H2D) StdErrY() float64 {
	return h.axis.dist.stdErrY()
}

// RMSX returns the RMS in X.
// Overflows are included in the computation.
func (h *H2D) RMSX() float64 {
	return h.axis.dist.rmsX()
}

// RMSY returns the RMS in Y.
// Overflows are included in the computation.
func (h *H2D) RMSY() float64 {
	return h.axis.dist.rmsY()
}

// Fill fills this histogram with (x,y) and weight w.
func (h *H2D) Fill(x, y, w float64) {
	h.axis.fill(x, y, w)
}

// MinX returns the low edge of the X-axis of this histogram.
func (h *H2D) MinX() float64 {
	return h.axis.minX()
}

// MaxX returns the high edge of the X-axis of this histogram.
func (h *H2D) MaxX() float64 {
	return h.axis.maxX()
}

// MinY returns the low edge of the Y-axis of this histogram.
func (h *H2D) MinY() float64 {
	return h.axis.minY()
}

// MaxY returns the high edge of the Y-axis of this histogram.
func (h *H2D) MaxY() float64 {
	return h.axis.maxY()
}

// GridXYZ returns an anonymous struct value that implements
// gonum/plot/plotter.GridXYZ and is ready to plot.
func (h *H2D) GridXYZ() h2dGridXYZ {
	return h2dGridXYZ{h}
}

type h2dGridXYZ struct {
	h *H2D
}

func (g h2dGridXYZ) Dims() (c, r int) {
	return g.h.axis.nx, g.h.axis.ny
}

func (g h2dGridXYZ) Z(c, r int) float64 {
	idx := r*g.h.axis.nx + c
	return g.h.axis.bins[idx].SumW()
}

func (g h2dGridXYZ) X(c int) float64 {
	return g.h.axis.bins[c].xrange.Min
}

func (g h2dGridXYZ) Y(r int) float64 {
	idx := r * g.h.axis.nx
	return g.h.axis.bins[idx].yrange.Min
}

// check various interfaces
var _ Object = (*H2D)(nil)
