// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"math"

	"go-hep.org/x/hep/rio"
)

// H1D is a 1-dim histogram with weighted entries.
type H1D struct {
	bng binning1D
	ann Annotation
}

// NewH1D returns a 1-dim histogram with n bins between xmin and xmax.
func NewH1D(n int, xmin, xmax float64) *H1D {
	return &H1D{
		bng: newBinning1D(n, xmin, xmax),
		ann: make(Annotation),
	}
}

// NewH1DFromEdges returns a 1-dim histogram given a slice of edges.
// The number of bins is thus len(edges)-1.
// It panics if the length of edges is <= 1.
// It panics if the edges are not sorted.
// It panics if there are duplicate edge values.
func NewH1DFromEdges(edges []float64) *H1D {
	return &H1D{
		bng: newBinning1DFromEdges(edges),
		ann: make(Annotation),
	}
}

// NewH1DFromBins returns a 1-dim histogram given a slice of 1-dim bins.
// It panics if the length of bins is < 1.
// It panics if the bins overlap.
func NewH1DFromBins(bins ...Range) *H1D {
	return &H1D{
		bng: newBinning1DFromBins(bins),
		ann: make(Annotation),
	}
}

// Name returns the name of this histogram, if any
func (h *H1D) Name() string {
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
func (h *H1D) Annotation() Annotation {
	return h.ann
}

// Rank returns the number of dimensions for this histogram
func (h *H1D) Rank() int {
	return 1
}

// Entries returns the number of entries in this histogram
func (h *H1D) Entries() int64 {
	return h.bng.entries()
}

// EffEntries returns the number of effective entries in this histogram
func (h *H1D) EffEntries() float64 {
	return h.bng.effEntries()
}

// Binning returns the binning of this histogram
func (h *H1D) Binning() *binning1D {
	return &h.bng
}

// SumW returns the sum of weights in this histogram.
// Overflows are included in the computation.
func (h *H1D) SumW() float64 {
	return h.bng.dist.SumW()
}

// SumW2 returns the sum of squared weights in this histogram.
// Overflows are included in the computation.
func (h *H1D) SumW2() float64 {
	return h.bng.dist.SumW2()
}

// XMean returns the mean X.
// Overflows are included in the computation.
func (h *H1D) XMean() float64 {
	return h.bng.dist.mean()
}

// XVariance returns the variance in X.
// Overflows are included in the computation.
func (h *H1D) XVariance() float64 {
	return h.bng.dist.variance()
}

// XStdDev returns the standard deviation in X.
// Overflows are included in the computation.
func (h *H1D) XStdDev() float64 {
	return h.bng.dist.stdDev()
}

// XStdErr returns the standard error in X.
// Overflows are included in the computation.
func (h *H1D) XStdErr() float64 {
	return h.bng.dist.stdErr()
}

// XRMS returns the XRMS in X.
// Overflows are included in the computation.
func (h *H1D) XRMS() float64 {
	return h.bng.dist.rms()
}

// Fill fills this histogram with x and weight w.
func (h *H1D) Fill(x, w float64) {
	h.bng.fill(x, w)
}

// XMin returns the low edge of the X-axis of this histogram.
func (h *H1D) XMin() float64 {
	return h.bng.xMin()
}

// XMax returns the high edge of the X-axis of this histogram.
func (h *H1D) XMax() float64 {
	return h.bng.xMax()
}

// Scale scales the content of each bin by the given factor.
func (h *H1D) Scale(factor float64) {
	h.bng.scaleW(factor)
}

// Integral computes the integral of the histogram.
//
// The number of parameters can be 0 or 2.
// If 0, overflows are included.
// If 2, the first parameter must be the lower bound of the range in which
// the integral is computed and the second one the upper range.
//
// If the lower bound is math.Inf(-1) then the underflow bin is included.
// If the upper bound is math.Inf(+1) then the overflow bin is included.
//
// Examples:
//
//    // integral of all in-range bins + overflows
//    v := h.Integral()
//
//    // integral of all in-range bins, underflow and overflow bins included.
//    v := h.Integral(math.Inf(-1), math.Inf(+1))
//
//    // integrall of all in-range bins, overflow bin included
//    v := h.Integral(h.Binning().LowerEdge(), math.Inf(+1))
//
//    // integrall of all bins for which the lower edge is in [0.5, 5.5)
//    v := h.Integral(0.5, 5.5)
func (h *H1D) Integral(args ...float64) float64 {
	min, max := 0., 0.
	switch len(args) {
	case 0:
		return h.SumW()
	case 2:
		min = args[0]
		max = args[1]
		if min > max {
			panic("hbook: min > max")
		}
	default:
		panic("hbook: invalid number of arguments. expected 0 or 2.")
	}

	integral := 0.0
	for _, bin := range h.bng.bins {
		v := bin.xrange.Min
		if min <= v && v < max {
			integral += bin.SumW()
		}
	}
	if math.IsInf(min, -1) {
		integral += h.bng.outflows[0].SumW()
	}
	if math.IsInf(max, +1) {
		integral += h.bng.outflows[1].SumW()
	}
	return integral
}

// Value returns the content of the idx-th bin.
//
// Value implements gonum/plot/plotter.Valuer
func (h *H1D) Value(i int) float64 {
	return h.bng.bins[i].SumW()
}

// Len returns the number of bins for this histogram
//
// Len implements gonum/plot/plotter.Valuer
func (h *H1D) Len() int {
	return len(h.bng.bins)
}

// XY returns the x,y values for the i-th bin
//
// XY implements gonum/plot/plotter.XYer
func (h *H1D) XY(i int) (float64, float64) {
	bin := h.bng.bins[i]
	x := bin.xrange.Min
	y := bin.SumW()
	return x, y
}

// DataRange implements the gonum/plot.DataRanger interface
func (h *H1D) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = h.XMin()
	xmax = h.XMax()
	ymin = +math.MaxFloat64
	ymax = -math.MaxFloat64
	for _, b := range h.bng.bins {
		v := b.SumW()
		ymax = math.Max(ymax, v)
		ymin = math.Min(ymin, v)
	}
	return xmin, xmax, ymin, ymax
}

// RioMarshal implements rio.RioMarshaler
func (h *H1D) RioMarshal(w io.Writer) error {
	data, err := h.MarshalBinary()
	if err != nil {
		return err
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(len(data)))
	_, err = w.Write(buf[:])
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// RioUnmarshal implements rio.RioUnmarshaler
func (h *H1D) RioUnmarshal(r io.Reader) error {
	buf := make([]byte, 8)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	n := int64(binary.LittleEndian.Uint64(buf))
	buf = make([]byte, int(n))
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	return h.UnmarshalBinary(buf)
}

// RioVersion implements rio.RioStreamer
func (h *H1D) RioVersion() rio.Version {
	return 0
}

// annToYODA creates a new Annotation with fields compatible with YODA
func (h *H1D) annToYODA() Annotation {
	ann := make(Annotation, len(h.ann))
	ann["Type"] = "Histo1D"
	ann["Path"] = "/" + h.Name()
	ann["Title"] = ""
	for k, v := range h.ann {
		ann[k] = v
	}
	return ann
}

// MarshalYODA implements the YODAMarshaler interface.
func (h *H1D) MarshalYODA() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := h.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_HISTO1D %s\n", ann["Path"])
	data, err := ann.MarshalYODA()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	fmt.Fprintf(buf, "# Mean: %e\n", h.XMean())
	fmt.Fprintf(buf, "# Area: %e\n", h.Integral())

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")
	d := h.bng.dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)
	d = h.bng.outflows[0]
	fmt.Fprintf(
		buf,
		"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)

	d = h.bng.outflows[1]
	fmt.Fprintf(
		buf,
		"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
	)

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t sumw\t sumw2\t sumwx\t sumwx2\t numEntries\n")
	for _, bin := range h.bng.bins {
		d := bin.dist
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			bin.xrange.Min, bin.xrange.Max,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_HISTO1D\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (h *H1D) UnmarshalYODA(data []byte) error {
	var err error
	var path string
	r := bytes.NewBuffer(data)
	_, err = fmt.Fscanf(r, "BEGIN YODA_HISTO1D %s\n", &path)
	if err != nil {
		return err
	}
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# Mean:"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid H1D-YODA data")
	}
	err = ann.UnmarshalYODA(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %v\nhbook: %q", err, string(r.Bytes()[:pos+1]))
	}
	h.ann = ann
	r.Next(pos)

	var ctx struct {
		total bool
		under bool
		over  bool
		bins  bool
	}

	// sets of xlow values, to infer number of bins in X.
	xset := make(map[float64]int)

	var (
		dist   dist1D
		oflows [2]dist1D
		bins   []Bin1D
		xmin   = math.Inf(+1)
		xmax   = math.Inf(-1)
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
		case bytes.HasPrefix(buf, []byte("END YODA_HISTO1D")):
			break scanLoop
		case !ctx.total && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.total = true
			d := &dist
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.sumW, &d.dist.sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
		case !ctx.under && bytes.HasPrefix(buf, []byte("Underflow\t")):
			ctx.under = true
			d := &oflows[0]
			_, err = fmt.Fscanf(
				rbuf,
				"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.sumW, &d.dist.sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
		case !ctx.over && bytes.HasPrefix(buf, []byte("Overflow\t")):
			ctx.over = true
			d := &oflows[1]
			_, err = fmt.Fscanf(
				rbuf,
				"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%d\n",
				&d.dist.sumW, &d.dist.sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			ctx.bins = true
		case ctx.bins:
			var bin Bin1D
			d := &bin.dist
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&bin.xrange.Min, &bin.xrange.Max,
				&d.dist.sumW, &d.dist.sumW2,
				&d.sumWX, &d.sumWX2,
				&d.dist.n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %v\nhbook: %q", err, string(buf))
			}
			xset[bin.xrange.Min] = 1
			xmin = math.Min(xmin, bin.xrange.Min)
			xmax = math.Max(xmax, bin.xrange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid H1D-YODA data: %q", string(buf))
		}
	}
	h.bng = binning1D{
		bins:     bins,
		dist:     dist,
		outflows: oflows,
		xrange:   Range{xmin, xmax},
	}
	return err
}

// check various interfaces
var _ Object = (*H1D)(nil)
var _ Histogram = (*H1D)(nil)

// serialization interfaces
var _ rio.Marshaler = (*H1D)(nil)
var _ rio.Unmarshaler = (*H1D)(nil)
var _ rio.Streamer = (*H1D)(nil)

func init() {
	gob.Register((*H1D)(nil))
}
