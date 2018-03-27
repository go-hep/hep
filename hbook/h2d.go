// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
)

// H2D is a 2-dim histogram with weighted entries.
type H2D struct {
	bng binning2D
	ann Annotation
}

// NewH2D creates a new 2-dim histogram.
func NewH2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) *H2D {
	return &H2D{
		bng: newBinning2D(nx, xlow, xhigh, ny, ylow, yhigh),
		ann: make(Annotation),
	}
}

// NewH2DFromEdges creates a new 2-dim histogram from slices
// of edges in x and y.
// The number of bins in x and y is thus len(edges)-1.
// It panics if the length of edges is <=1 (in any dimension.)
// It panics if the edges are not sorted (in any dimension.)
// It panics if there are duplicate edge values (in any dimension.)
func NewH2DFromEdges(xedges []float64, yedges []float64) *H2D {
	return &H2D{
		bng: newBinning2DFromEdges(xedges, yedges),
		ann: make(Annotation),
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
	return h.bng.entries()
}

// EffEntries returns the number of effective entries in this histogram
func (h *H2D) EffEntries() float64 {
	return h.bng.effEntries()
}

// Binning returns the binning of this histogram
func (h *H2D) Binning() *binning2D {
	return &h.bng
}

// SumW returns the sum of weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW() float64 {
	return h.bng.dist.SumW()
}

// SumW2 returns the sum of squared weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW2() float64 {
	return h.bng.dist.SumW2()
}

// XMean returns the mean X.
// Overflows are included in the computation.
func (h *H2D) XMean() float64 {
	return h.bng.dist.xMean()
}

// YMean returns the mean Y.
// Overflows are included in the computation.
func (h *H2D) YMean() float64 {
	return h.bng.dist.yMean()
}

// XVariance returns the variance in X.
// Overflows are included in the computation.
func (h *H2D) XVariance() float64 {
	return h.bng.dist.xVariance()
}

// YVariance returns the variance in Y.
// Overflows are included in the computation.
func (h *H2D) YVariance() float64 {
	return h.bng.dist.yVariance()
}

// XStdDev returns the standard deviation in X.
// Overflows are included in the computation.
func (h *H2D) XStdDev() float64 {
	return h.bng.dist.xStdDev()
}

// YStdDev returns the standard deviation in Y.
// Overflows are included in the computation.
func (h *H2D) YStdDev() float64 {
	return h.bng.dist.yStdDev()
}

// XStdErr returns the standard error in X.
// Overflows are included in the computation.
func (h *H2D) XStdErr() float64 {
	return h.bng.dist.xStdErr()
}

// YStdErr returns the standard error in Y.
// Overflows are included in the computation.
func (h *H2D) YStdErr() float64 {
	return h.bng.dist.yStdErr()
}

// XRMS returns the RMS in X.
// Overflows are included in the computation.
func (h *H2D) XRMS() float64 {
	return h.bng.dist.xRMS()
}

// YRMS returns the RMS in Y.
// Overflows are included in the computation.
func (h *H2D) YRMS() float64 {
	return h.bng.dist.yRMS()
}

// Fill fills this histogram with (x,y) and weight w.
func (h *H2D) Fill(x, y, w float64) {
	h.bng.fill(x, y, w)
}

// XMin returns the low edge of the X-axis of this histogram.
func (h *H2D) XMin() float64 {
	return h.bng.xMin()
}

// XMax returns the high edge of the X-axis of this histogram.
func (h *H2D) XMax() float64 {
	return h.bng.xMax()
}

// YMin returns the low edge of the Y-axis of this histogram.
func (h *H2D) YMin() float64 {
	return h.bng.yMin()
}

// YMax returns the high edge of the Y-axis of this histogram.
func (h *H2D) YMax() float64 {
	return h.bng.yMax()
}

// Integral computes the integral of the histogram.
//
// Overflows are included in the computation.
func (h *H2D) Integral() float64 {
	return h.SumW()
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
	return g.h.bng.nx, g.h.bng.ny
}

func (g h2dGridXYZ) Z(c, r int) float64 {
	idx := r*g.h.bng.nx + c
	return g.h.bng.bins[idx].SumW()
}

func (g h2dGridXYZ) X(c int) float64 {
	return g.h.bng.bins[c].XMid()
}

func (g h2dGridXYZ) Y(r int) float64 {
	idx := r * g.h.bng.nx
	return g.h.bng.bins[idx].YMid()
}

// check various interfaces
var _ Object = (*H2D)(nil)
var _ Histogram = (*H2D)(nil)

// annToYODA creates a new Annotation with fields compatible with YODA
func (h *H2D) annToYODA() Annotation {
	ann := make(Annotation, len(h.ann))
	ann["Type"] = "Histo2D"
	ann["Path"] = "/" + h.Name()
	ann["Title"] = ""
	for k, v := range h.ann {
		if k == "name" {
			continue
		}
		ann[k] = v
	}
	return ann
}

// annFromYODA creates a new Annotation from YODA compatible fields
func (h *H2D) annFromYODA(ann Annotation) {
	if len(h.ann) == 0 {
		h.ann = make(Annotation, len(ann))
	}
	for k, v := range ann {
		switch k {
		case "Type":
			// noop
		case "Path":
			h.ann["name"] = string(v.(string)[1:]) // skip leading '/'
		default:
			h.ann[k] = v
		}
	}
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H2D) MarshalYODA() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := h.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D %s\n", ann["Path"])
	data, err := ann.MarshalYODA()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	fmt.Fprintf(buf, "# Mean: (%e, %e)\n", h.XMean(), h.YMean())
	fmt.Fprintf(buf, "# Volume: %e\n", h.Integral())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	d := h.bng.dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.sumWXY, d.Entries(),
	)

	// outflows
	fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := 0; ix < h.bng.nx; ix++ {
		for iy := 0; iy < h.bng.ny; iy++ {
			bin := h.bng.bins[iy*h.bng.nx+ix]
			d := bin.dist
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				bin.xrange.Min, bin.xrange.Max, bin.yrange.Min, bin.yrange.Max,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.sumWXY, d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H2D) UnmarshalYODA(data []byte) error {
	r := bytes.NewBuffer(data)
	_, err := readYODAHeader(r, "BEGIN YODA_HISTO2D")
	if err != nil {
		return err
	}
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# Mean:"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid H2D-YODA data")
	}
	err = ann.UnmarshalYODA(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %v\nhbook: %q", err, string(r.Bytes()[:pos+1]))
	}
	h.annFromYODA(ann)
	r.Next(pos)

	var ctx struct {
		dist bool
		bins bool
	}

	// sets of xlow and ylow values, to infer number of bins in X and Y.
	xset := make(map[float64]int)
	yset := make(map[float64]int)

	var (
		dist dist2D
		bins []Bin2D
		xmin = math.Inf(+1)
		xmax = math.Inf(-1)
		ymin = math.Inf(+1)
		ymax = math.Inf(-1)
	)
	s := bufio.NewScanner(r)
scanLoop:
	for s.Scan() {
		buf := s.Bytes()
		if len(buf) == 0 || buf[0] == '#' {
			continue
		}
		rbuf := bytes.NewReader(buf)
		switch {
		case bytes.HasPrefix(buf, []byte("END YODA_HISTO2D")):
			break scanLoop
		case !ctx.dist && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.dist = true
			d := &dist
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&d.x.dist.sumW, &d.x.dist.sumW2,
				&d.x.sumWX, &d.x.sumWX2,
				&d.y.sumWX, &d.y.sumWX2,
				&d.sumWXY, &d.x.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			d.y.dist = d.x.dist
			ctx.bins = true
		case ctx.bins:
			var bin Bin2D
			d := &bin.dist
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&bin.xrange.Min, &bin.xrange.Max, &bin.yrange.Min, &bin.yrange.Max,
				&d.x.dist.sumW, &d.x.dist.sumW2,
				&d.x.sumWX, &d.x.sumWX2,
				&d.y.sumWX, &d.y.sumWX2,
				&d.sumWXY, &d.x.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			d.y.dist = d.x.dist
			xset[bin.xrange.Min] = 1
			yset[bin.yrange.Min] = 1
			xmin = math.Min(xmin, bin.xrange.Min)
			xmax = math.Max(xmax, bin.xrange.Max)
			ymin = math.Min(ymin, bin.yrange.Min)
			ymax = math.Max(ymax, bin.yrange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid H2D-YODA data: %q", string(buf))
		}
	}
	h.bng = newBinning2D(len(xset), xmin, xmax, len(yset), ymin, ymax)
	h.bng.dist = dist
	// YODA bins are transposed wrt ours
	for ix := 0; ix < h.bng.nx; ix++ {
		for iy := 0; iy < h.bng.ny; iy++ {
			h.bng.bins[iy*h.bng.nx+ix] = bins[ix*h.bng.ny+iy]
		}
	}
	return err
}
