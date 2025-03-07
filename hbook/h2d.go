// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"strings"
)

// H2D is a 2-dim histogram with weighted entries.
type H2D struct {
	Binning Binning2D
	Ann     Annotation
}

// NewH2D creates a new 2-dim histogram.
func NewH2D(nx int, xlow, xhigh float64, ny int, ylow, yhigh float64) *H2D {
	return &H2D{
		Binning: newBinning2D(nx, xlow, xhigh, ny, ylow, yhigh),
		Ann:     make(Annotation),
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
		Binning: newBinning2DFromEdges(xedges, yedges),
		Ann:     make(Annotation),
	}
}

// Name returns the name of this histogram, if any
func (h *H2D) Name() string {
	v, ok := h.Ann["name"]
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
	return h.Ann
}

// Rank returns the number of dimensions for this histogram
func (h *H2D) Rank() int {
	return 2
}

// Entries returns the number of entries in this histogram
func (h *H2D) Entries() int64 {
	return h.Binning.entries()
}

// EffEntries returns the number of effective entries in this histogram
func (h *H2D) EffEntries() float64 {
	return h.Binning.effEntries()
}

// SumW returns the sum of weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW() float64 {
	return h.Binning.Dist.SumW()
}

// SumW2 returns the sum of squared weights in this histogram.
// Overflows are included in the computation.
func (h *H2D) SumW2() float64 {
	return h.Binning.Dist.SumW2()
}

// SumWX returns the 1st order weighted x moment
// Overflows are included in the computation.
func (h *H2D) SumWX() float64 {
	return h.Binning.Dist.SumWX()
}

// SumWX2 returns the 2nd order weighted x moment
// Overflows are included in the computation.
func (h *H2D) SumWX2() float64 {
	return h.Binning.Dist.SumWX2()
}

// SumWY returns the 1st order weighted y moment
// Overflows are included in the computation.
func (h *H2D) SumWY() float64 {
	return h.Binning.Dist.SumWY()
}

// SumWY2 returns the 2nd order weighted y moment
// Overflows are included in the computation.
func (h *H2D) SumWY2() float64 {
	return h.Binning.Dist.SumWY2()
}

// SumWXY returns the 1st order weighted x*y moment
// Overflows are included in the computation.
func (h *H2D) SumWXY() float64 {
	return h.Binning.Dist.SumWXY()
}

// XMean returns the mean X.
// Overflows are included in the computation.
func (h *H2D) XMean() float64 {
	return h.Binning.Dist.xMean()
}

// YMean returns the mean Y.
// Overflows are included in the computation.
func (h *H2D) YMean() float64 {
	return h.Binning.Dist.yMean()
}

// XVariance returns the variance in X.
// Overflows are included in the computation.
func (h *H2D) XVariance() float64 {
	return h.Binning.Dist.xVariance()
}

// YVariance returns the variance in Y.
// Overflows are included in the computation.
func (h *H2D) YVariance() float64 {
	return h.Binning.Dist.yVariance()
}

// XStdDev returns the standard deviation in X.
// Overflows are included in the computation.
func (h *H2D) XStdDev() float64 {
	return h.Binning.Dist.xStdDev()
}

// YStdDev returns the standard deviation in Y.
// Overflows are included in the computation.
func (h *H2D) YStdDev() float64 {
	return h.Binning.Dist.yStdDev()
}

// XStdErr returns the standard error in X.
// Overflows are included in the computation.
func (h *H2D) XStdErr() float64 {
	return h.Binning.Dist.xStdErr()
}

// YStdErr returns the standard error in Y.
// Overflows are included in the computation.
func (h *H2D) YStdErr() float64 {
	return h.Binning.Dist.yStdErr()
}

// XRMS returns the RMS in X.
// Overflows are included in the computation.
func (h *H2D) XRMS() float64 {
	return h.Binning.Dist.xRMS()
}

// YRMS returns the RMS in Y.
// Overflows are included in the computation.
func (h *H2D) YRMS() float64 {
	return h.Binning.Dist.yRMS()
}

// Fill fills this histogram with (x,y) and weight w.
func (h *H2D) Fill(x, y, w float64) {
	h.Binning.fill(x, y, w)
}

// FillN fills this histogram with the provided slices (xs,ys) and weights ws.
// if ws is nil, the histogram will be filled with entries of weight 1.
// Otherwise, FillN panics if the slices lengths differ.
func (h *H2D) FillN(xs, ys, ws []float64) {
	switch ws {
	case nil:
		if len(xs) != len(ys) {
			panic(fmt.Errorf("hbook: lengths mismatch"))
		}
		for i := range xs {
			x := xs[i]
			y := ys[i]
			h.Binning.fill(x, y, 1)
		}
	default:
		if len(xs) != len(ys) {
			panic(fmt.Errorf("hbook: lengths mismatch"))
		}
		if len(xs) != len(ws) {
			panic(fmt.Errorf("hbook: lengths mismatch"))
		}
		for i := range xs {
			x := xs[i]
			y := ys[i]
			w := ws[i]
			h.Binning.fill(x, y, w)
		}
	}
}

// Bin returns the bin at coordinates (x,y) for this 2-dim histogram.
// Bin returns nil for under/over flow bins.
func (h *H2D) Bin(x, y float64) *Bin2D {
	idx := h.Binning.coordToIndex(x, y)
	if idx < 0 {
		return nil
	}
	return &h.Binning.Bins[idx]
}

// XMin returns the low edge of the X-axis of this histogram.
func (h *H2D) XMin() float64 {
	return h.Binning.xMin()
}

// XMax returns the high edge of the X-axis of this histogram.
func (h *H2D) XMax() float64 {
	return h.Binning.xMax()
}

// YMin returns the low edge of the Y-axis of this histogram.
func (h *H2D) YMin() float64 {
	return h.Binning.yMin()
}

// YMax returns the high edge of the Y-axis of this histogram.
func (h *H2D) YMax() float64 {
	return h.Binning.yMax()
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
	return g.h.Binning.Nx, g.h.Binning.Ny
}

func (g h2dGridXYZ) Z(c, r int) float64 {
	idx := r*g.h.Binning.Nx + c
	return g.h.Binning.Bins[idx].SumW()
}

func (g h2dGridXYZ) X(c int) float64 {
	return g.h.Binning.Bins[c].XMid()
}

func (g h2dGridXYZ) Y(r int) float64 {
	idx := r * g.h.Binning.Nx
	return g.h.Binning.Bins[idx].YMid()
}

// check various interfaces
var _ Object = (*H2D)(nil)
var _ Histogram = (*H2D)(nil)

// annToYODA creates a new Annotation with fields compatible with YODA
func (h *H2D) annToYODA() Annotation {
	ann := make(Annotation, len(h.Ann))
	ann["Type"] = "Histo2D"
	ann["Path"] = "/" + h.Name()
	ann["Title"] = ""
	for k, v := range h.Ann {
		if k == "name" {
			continue
		}
		if k == "title" {
			ann["Title"] = v
			continue
		}
		ann[k] = v
	}
	return ann
}

// annFromYODA creates a new Annotation from YODA compatible fields
func (h *H2D) annFromYODA(ann Annotation) {
	if len(h.Ann) == 0 {
		h.Ann = make(Annotation, len(ann))
	}
	for k, v := range ann {
		switch k {
		case "Type":
			// noop
		case "Path":
			name := v.(string)
			name = strings.TrimPrefix(name, "/")
			h.Ann["name"] = name
		case "Title":
			h.Ann["title"] = v
		default:
			h.Ann[k] = v
		}
	}
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H2D) MarshalYODA() ([]byte, error) {
	return h.marshalYODAv2()
}

func (h *H2D) marshalYODAv1() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := h.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D %s\n", ann["Path"])
	data, err := ann.marshalYODAv1()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	fmt.Fprintf(buf, "# Mean: (%e, %e)\n", h.XMean(), h.YMean())
	fmt.Fprintf(buf, "# Volume: %e\n", h.Integral())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	d := h.Binning.Dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), d.Entries(),
	)

	// outflows
	fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := range h.Binning.Nx {
		for iy := range h.Binning.Ny {
			bin := h.Binning.Bins[iy*h.Binning.Nx+ix]
			d := bin.Dist
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				bin.XRange.Min, bin.XRange.Max, bin.YRange.Min, bin.YRange.Max,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), d.Entries(),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D\n\n")
	return buf.Bytes(), err
}

func (h *H2D) marshalYODAv2() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := h.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_HISTO2D_V2 %s\n", ann["Path"])
	data, err := ann.marshalYODAv2()
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	buf.Write([]byte("---\n"))

	fmt.Fprintf(buf, "# Mean: (%e, %e)\n", h.XMean(), h.YMean())
	fmt.Fprintf(buf, "# Volume: %e\n", h.Integral())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	d := h.Binning.Dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), float64(d.Entries()),
	)

	// outflows
	fmt.Fprintf(buf, "# 2D outflow persistency not currently supported until API is stable\n")

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t ylow\t yhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t sumwxy\t numEntries\n")
	for ix := range h.Binning.Nx {
		for iy := range h.Binning.Ny {
			bin := h.Binning.Bins[iy*h.Binning.Nx+ix]
			d := bin.Dist
			fmt.Fprintf(
				buf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				bin.XRange.Min, bin.XRange.Max, bin.YRange.Min, bin.YRange.Max,
				d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.SumWXY(), float64(d.Entries()),
			)
		}
	}
	fmt.Fprintf(buf, "END YODA_HISTO2D_V2\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H2D) UnmarshalYODA(data []byte) error {
	r := newRBuffer(data)
	_, vers, err := readYODAHeader(r, "BEGIN YODA_HISTO2D")
	if err != nil {
		return err
	}
	switch vers {
	case 1:
		return h.unmarshalYODAv1(r)
	case 2:
		return h.unmarshalYODAv2(r)
	default:
		return fmt.Errorf("hbook: invalid YODA version %v", vers)
	}
}

func (h *H2D) unmarshalYODAv1(r *rbuffer) error {
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# Mean:"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid H2D-YODA data")
	}
	err := ann.unmarshalYODAv1(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %q\nhbook: %w", string(r.Bytes()[:pos+1]), err)
	}
	h.annFromYODA(ann)
	r.next(pos)

	var ctx struct {
		dist bool
		bins bool
	}

	// sets of xlow and ylow values, to infer number of bins in X and Y.
	xset := make(map[float64]int)
	yset := make(map[float64]int)

	var (
		dist Dist2D
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
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.Stats.SumWXY, &d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist = d.X.Dist
			ctx.bins = true
		case ctx.bins:
			var bin Bin2D
			d := &bin.Dist
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&bin.XRange.Min, &bin.XRange.Max, &bin.YRange.Min, &bin.YRange.Max,
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.Stats.SumWXY, &d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist = d.X.Dist
			xset[bin.XRange.Min] = 1
			yset[bin.YRange.Min] = 1
			xmin = math.Min(xmin, bin.XRange.Min)
			xmax = math.Max(xmax, bin.XRange.Max)
			ymin = math.Min(ymin, bin.YRange.Min)
			ymax = math.Max(ymax, bin.YRange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid H2D-YODA data: %q", string(buf))
		}
	}
	h.Binning = newBinning2D(len(xset), xmin, xmax, len(yset), ymin, ymax)
	h.Binning.Dist = dist
	// YODA bins are transposed wrt ours
	for ix := range h.Binning.Nx {
		for iy := range h.Binning.Ny {
			h.Binning.Bins[iy*h.Binning.Nx+ix] = bins[ix*h.Binning.Ny+iy]
		}
	}
	return err
}

func (h *H2D) unmarshalYODAv2(r *rbuffer) error {
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# Mean:"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid H2D-YODA data")
	}
	err := ann.unmarshalYODAv2(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %q\nhbook: %w", string(r.Bytes()[:pos+1]), err)
	}
	h.annFromYODA(ann)
	r.next(pos)

	var ctx struct {
		dist bool
		bins bool
	}

	// sets of xlow and ylow values, to infer number of bins in X and Y.
	xset := make(map[float64]int)
	yset := make(map[float64]int)

	var (
		dist Dist2D
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
		case bytes.HasPrefix(buf, []byte("END YODA_HISTO2D_V2")):
			break scanLoop
		case !ctx.dist && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.dist = true
			d := &dist
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.Stats.SumWXY, &n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist = d.X.Dist
			ctx.bins = true
		case ctx.bins:
			var bin Bin2D
			d := &bin.Dist
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&bin.XRange.Min, &bin.XRange.Max, &bin.YRange.Min, &bin.YRange.Max,
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.Stats.SumWXY, &n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist = d.X.Dist
			xset[bin.XRange.Min] = 1
			yset[bin.YRange.Min] = 1
			xmin = math.Min(xmin, bin.XRange.Min)
			xmax = math.Max(xmax, bin.XRange.Max)
			ymin = math.Min(ymin, bin.YRange.Min)
			ymax = math.Max(ymax, bin.YRange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid H2D-YODA data: %q", string(buf))
		}
	}
	h.Binning = newBinning2D(len(xset), xmin, xmax, len(yset), ymin, ymax)
	h.Binning.Dist = dist
	// YODA bins are transposed wrt ours
	for ix := range h.Binning.Nx {
		for iy := range h.Binning.Ny {
			h.Binning.Bins[iy*h.Binning.Nx+ix] = bins[ix*h.Binning.Ny+iy]
		}
	}
	return err
}
