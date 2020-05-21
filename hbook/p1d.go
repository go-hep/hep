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

// P1D is a 1-dim profile histogram.
type P1D struct {
	bng binningP1D
	ann Annotation
}

// NewP1D returns a 1-dim profile histogram with n bins between xmin and xmax.
func NewP1D(n int, xmin, xmax float64) *P1D {
	return &P1D{
		bng: newBinningP1D(n, xmin, xmax),
		ann: make(Annotation),
	}
}

/*
// FIXME(sbinet): need support of variable-size bins
//
// NewP1DFromS2D creates a 1-dim profile histogram from a 2-dim scatter's binning.
func NewP1DFromH1D(s*S2D) *P1D {
	return &P1D{
		bng: newBinningP1D(len(h.Binning().Bins()), h.XMin(), h.XMax()),
		ann: make(Annotation),
	}
}
*/

// NewP1DFromH1D creates a 1-dim profile histogram from a 1-dim histogram's binning.
func NewP1DFromH1D(h *H1D) *P1D {
	return &P1D{
		bng: newBinningP1D(len(h.Binning.Bins), h.XMin(), h.XMax()),
		ann: make(Annotation),
	}
}

// Name returns the name of this profile histogram, if any
func (p *P1D) Name() string {
	v, ok := p.ann["name"]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

// Annotation returns the annotations attached to this profile histogram
func (p *P1D) Annotation() Annotation {
	return p.ann
}

// Rank returns the number of dimensions for this profile histogram
func (p *P1D) Rank() int {
	return 1
}

// Entries returns the number of entries in this profile histogram
func (p *P1D) Entries() int64 {
	return p.bng.entries()
}

// EffEntries returns the number of effective entries in this profile histogram
func (p *P1D) EffEntries() float64 {
	return p.bng.effEntries()
}

// Binning returns the binning of this profile histogram
func (p *P1D) Binning() *binningP1D {
	return &p.bng
}

// SumW returns the sum of weights in this profile histogram.
// Overflows are included in the computation.
func (p *P1D) SumW() float64 {
	return p.bng.dist.SumW()
}

// SumW2 returns the sum of squared weights in this profile histogram.
// Overflows are included in the computation.
func (p *P1D) SumW2() float64 {
	return p.bng.dist.SumW2()
}

// XMean returns the mean X.
// Overflows are included in the computation.
func (p *P1D) XMean() float64 {
	return p.bng.dist.xMean()
}

// XVariance returns the variance in X.
// Overflows are included in the computation.
func (p *P1D) XVariance() float64 {
	return p.bng.dist.xVariance()
}

// XStdDev returns the standard deviation in X.
// Overflows are included in the computation.
func (p *P1D) XStdDev() float64 {
	return p.bng.dist.xStdDev()
}

// XStdErr returns the standard error in X.
// Overflows are included in the computation.
func (p *P1D) XStdErr() float64 {
	return p.bng.dist.xStdErr()
}

// XRMS returns the XRMS in X.
// Overflows are included in the computation.
func (p *P1D) XRMS() float64 {
	return p.bng.dist.xRMS()
}

// Fill fills this histogram with x,y and weight w.
func (p *P1D) Fill(x, y, w float64) {
	p.bng.fill(x, y, w)
}

// XMin returns the low edge of the X-axis of this profile histogram.
func (p *P1D) XMin() float64 {
	return p.bng.xMin()
}

// XMax returns the high edge of the X-axis of this profile histogram.
func (p *P1D) XMax() float64 {
	return p.bng.xMax()
}

// Scale scales the content of each bin by the given factor.
func (p *P1D) Scale(factor float64) {
	p.bng.scaleW(factor)
}

// check various interfaces
var _ Object = (*P1D)(nil)
var _ Histogram = (*P1D)(nil)

// annToYODA creates a new Annotation with fields compatible with YODA
func (p *P1D) annToYODA() Annotation {
	ann := make(Annotation, len(p.ann))
	ann["Type"] = "Profile1D"
	ann["Path"] = "/" + p.Name()
	ann["Title"] = ""
	for k, v := range p.ann {
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
func (p *P1D) annFromYODA(ann Annotation) {
	if len(p.ann) == 0 {
		p.ann = make(Annotation, len(ann))
	}
	for k, v := range ann {
		switch k {
		case "Type":
			// noop
		case "Path":
			name := v.(string)
			name = strings.TrimPrefix(name, "/")
			p.ann["name"] = name
		case "Title":
			p.ann["title"] = v
		default:
			p.ann[k] = v
		}
	}
}

// MarshalYODA implements the YODAMarshaler interface.
func (p *P1D) MarshalYODA() ([]byte, error) {
	return p.marshalYODAv2()
}

func (p *P1D) marshalYODAv1() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := p.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_PROFILE1D %s\n", ann["Path"])
	data, err := ann.marshalYODAv1()
	if err != nil {
		return nil, err
	}
	buf.Write(data)

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t numEntries\n")
	d := p.bng.dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.Entries(),
	)

	d = p.bng.outflows[0]
	fmt.Fprintf(
		buf,
		"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.Entries(),
	)

	d = p.bng.outflows[1]
	fmt.Fprintf(
		buf,
		"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.Entries(),
	)

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t numEntries\n")
	for _, bin := range p.bng.bins {
		d := bin.dist
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
			bin.xrange.Min, bin.xrange.Max,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), d.Entries(),
		)
	}
	fmt.Fprintf(buf, "END YODA_PROFILE1D\n\n")
	return buf.Bytes(), err
}

func (p *P1D) marshalYODAv2() ([]byte, error) {
	buf := new(bytes.Buffer)
	ann := p.annToYODA()
	fmt.Fprintf(buf, "BEGIN YODA_PROFILE1D_V2 %s\n", ann["Path"])
	data, err := ann.marshalYODAv2()
	if err != nil {
		return nil, err
	}
	buf.Write(data)
	buf.Write([]byte("---\n"))

	fmt.Fprintf(buf, "# ID\t ID\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t numEntries\n")
	d := p.bng.dist
	fmt.Fprintf(
		buf,
		"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), float64(d.Entries()),
	)

	d = p.bng.outflows[0]
	fmt.Fprintf(
		buf,
		"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), float64(d.Entries()),
	)

	d = p.bng.outflows[1]
	fmt.Fprintf(
		buf,
		"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
		d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), float64(d.Entries()),
	)

	// bins
	fmt.Fprintf(buf, "# xlow\t xhigh\t sumw\t sumw2\t sumwx\t sumwx2\t sumwy\t sumwy2\t numEntries\n")
	for _, bin := range p.bng.bins {
		d := bin.dist
		fmt.Fprintf(
			buf,
			"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
			bin.xrange.Min, bin.xrange.Max,
			d.SumW(), d.SumW2(), d.SumWX(), d.SumWX2(), d.SumWY(), d.SumWY2(), float64(d.Entries()),
		)
	}
	fmt.Fprintf(buf, "END YODA_PROFILE1D_V2\n\n")
	return buf.Bytes(), err
}

// UnmarshalYODA implements the YODAUnmarshaler interface.
func (p *P1D) UnmarshalYODA(data []byte) error {
	r := newRBuffer(data)
	_, vers, err := readYODAHeader(r, "BEGIN YODA_PROFILE1D")
	if err != nil {
		return err
	}
	switch vers {
	case 1:
		return p.unmarshalYODAv1(r)
	case 2:
		return p.unmarshalYODAv2(r)
	default:
		return fmt.Errorf("hbook: invalid YODA version %v", vers)
	}
}

func (p *P1D) unmarshalYODAv1(r *rbuffer) error {
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# ID\t ID\t"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid P1D-YODA data")
	}
	err := ann.unmarshalYODAv1(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %q\nhbook: %w", string(r.Bytes()[:pos+1]), err)
	}
	p.annFromYODA(ann)
	r.next(pos)

	var ctx struct {
		total bool
		under bool
		over  bool
		bins  bool
	}

	// sets of xlow values, to infer number of bins in X.
	xset := make(map[float64]int)

	var (
		dist   Dist2D
		oflows [2]Dist2D
		bins   []BinP1D
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
		case bytes.HasPrefix(buf, []byte("END YODA_PROFILE1D")):
			break scanLoop
		case !ctx.total && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.total = true
			d := &dist
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist.N = d.X.Dist.N
		case !ctx.under && bytes.HasPrefix(buf, []byte("Underflow\t")):
			ctx.under = true
			d := &oflows[0]
			_, err = fmt.Fscanf(
				rbuf,
				"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist.N = d.X.Dist.N
		case !ctx.over && bytes.HasPrefix(buf, []byte("Overflow\t")):
			ctx.over = true
			d := &oflows[1]
			_, err = fmt.Fscanf(
				rbuf,
				"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist.N = d.X.Dist.N
			ctx.bins = true
		case ctx.bins:
			var bin BinP1D
			d := &bin.dist
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%d\n",
				&bin.xrange.Min, &bin.xrange.Max,
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&d.X.Dist.N,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.Y.Dist.N = d.X.Dist.N
			xset[bin.xrange.Min] = 1
			xmin = math.Min(xmin, bin.xrange.Min)
			xmax = math.Max(xmax, bin.xrange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid P1D-YODA data: %q", string(buf))
		}
	}
	p.bng = newBinningP1D(len(xset), xmin, xmax)
	p.bng.dist = dist
	p.bng.bins = bins
	p.bng.outflows = oflows
	return err
}

func (p *P1D) unmarshalYODAv2(r *rbuffer) error {
	ann := make(Annotation)

	// pos of end of annotations
	pos := bytes.Index(r.Bytes(), []byte("\n# ID\t ID\t"))
	if pos < 0 {
		return fmt.Errorf("hbook: invalid P1D-YODA data")
	}
	err := ann.unmarshalYODAv2(r.Bytes()[:pos+1])
	if err != nil {
		return fmt.Errorf("hbook: %q\nhbook: %w", string(r.Bytes()[:pos+1]), err)
	}
	p.annFromYODA(ann)
	r.next(pos)

	var ctx struct {
		total bool
		under bool
		over  bool
		bins  bool
	}

	// sets of xlow values, to infer number of bins in X.
	xset := make(map[float64]int)

	var (
		dist   Dist2D
		oflows [2]Dist2D
		bins   []BinP1D
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
		case bytes.HasPrefix(buf, []byte("END YODA_PROFILE1D_V2")):
			break scanLoop
		case !ctx.total && bytes.HasPrefix(buf, []byte("Total   \t")):
			ctx.total = true
			d := &dist
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"Total   \tTotal   \t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist.N = d.X.Dist.N
		case !ctx.under && bytes.HasPrefix(buf, []byte("Underflow\t")):
			ctx.under = true
			d := &oflows[0]
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"Underflow\tUnderflow\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist.N = d.X.Dist.N
		case !ctx.over && bytes.HasPrefix(buf, []byte("Overflow\t")):
			ctx.over = true
			d := &oflows[1]
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"Overflow\tOverflow\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist.N = d.X.Dist.N
			ctx.bins = true
		case ctx.bins:
			var bin BinP1D
			d := &bin.dist
			var n float64
			_, err = fmt.Fscanf(
				rbuf,
				"%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\t%e\n",
				&bin.xrange.Min, &bin.xrange.Max,
				&d.X.Dist.SumW, &d.X.Dist.SumW2,
				&d.X.Stats.SumWX, &d.X.Stats.SumWX2,
				&d.Y.Stats.SumWX, &d.Y.Stats.SumWX2,
				&n,
			)
			if err != nil {
				return fmt.Errorf("hbook: %q\nhbook: %w", string(buf), err)
			}
			d.X.Dist.N = int64(n)
			d.Y.Dist.N = d.X.Dist.N
			xset[bin.xrange.Min] = 1
			xmin = math.Min(xmin, bin.xrange.Min)
			xmax = math.Max(xmax, bin.xrange.Max)
			bins = append(bins, bin)

		default:
			return fmt.Errorf("hbook: invalid P1D-YODA data: %q", string(buf))
		}
	}
	p.bng = newBinningP1D(len(xset), xmin, xmax)
	p.bng.dist = dist
	p.bng.bins = bins
	p.bng.outflows = oflows
	return err
}

// binningP1D is a 1-dim binning for 1-dim profile histograms.
type binningP1D struct {
	bins     []BinP1D
	dist     Dist2D
	outflows [2]Dist2D
	xrange   Range
	xstep    float64
}

func newBinningP1D(n int, xmin, xmax float64) binningP1D {
	if xmin >= xmax {
		panic("hbook: invalid X-axis limits")
	}
	if n <= 0 {
		panic("hbook: X-axis with zero bins")
	}
	bng := binningP1D{
		bins:   make([]BinP1D, n),
		xrange: Range{Min: xmin, Max: xmax},
	}
	bng.xstep = float64(n) / bng.xrange.Width()
	width := bng.xrange.Width() / float64(n)
	for i := range bng.bins {
		bin := &bng.bins[i]
		bin.xrange.Min = xmin + float64(i)*width
		bin.xrange.Max = xmin + float64(i+1)*width
	}

	return bng
}

func (bng *binningP1D) entries() int64 {
	return bng.dist.Entries()
}

func (bng *binningP1D) effEntries() float64 {
	return bng.dist.EffEntries()
}

// xMin returns the low edge of the X-axis
func (bng *binningP1D) xMin() float64 {
	return bng.xrange.Min
}

// xMax returns the high edge of the X-axis
func (bng *binningP1D) xMax() float64 {
	return bng.xrange.Max
}

func (bng *binningP1D) fill(x, y, w float64) {
	idx := bng.coordToIndex(x)
	bng.dist.fill(x, y, w)
	if idx < 0 {
		bng.outflows[-idx-1].fill(x, y, w)
		return
	}
	bng.bins[idx].fill(x, y, w)
}

// coordToIndex returns the bin index corresponding to the coordinate x.
func (bng *binningP1D) coordToIndex(x float64) int {
	switch {
	default:
		i := int((x - bng.xrange.Min) * bng.xstep)
		return i
	case x < bng.xrange.Min:
		return UnderflowBin1D
	case x >= bng.xrange.Max:
		return OverflowBin1D
	}
}

func (bng *binningP1D) scaleW(f float64) {
	bng.dist.scaleW(f)
	bng.outflows[0].scaleW(f)
	bng.outflows[1].scaleW(f)
	for i := range bng.bins {
		bin := &bng.bins[i]
		bin.scaleW(f)
	}
}

// Bins returns the slice of bins for this binning.
func (bng *binningP1D) Bins() []BinP1D {
	return bng.bins
}

// BinP1D models a bin in a 1-dim space.
type BinP1D struct {
	xrange Range
	dist   Dist2D
}

// Rank returns the number of dimensions for this bin.
func (BinP1D) Rank() int { return 1 }

func (b *BinP1D) scaleW(f float64) {
	b.dist.scaleW(f)
}

func (b *BinP1D) fill(x, y, w float64) {
	b.dist.fill(x, y, w)
}

// Entries returns the number of entries in this bin.
func (b *BinP1D) Entries() int64 {
	return b.dist.Entries()
}

// EffEntries returns the effective number of entries \f$ = (\sum w)^2 / \sum w^2 \f$
func (b *BinP1D) EffEntries() float64 {
	return b.dist.EffEntries()
}

// SumW returns the sum of weights in this bin.
func (b *BinP1D) SumW() float64 {
	return b.dist.SumW()
}

// SumW2 returns the sum of squared weights in this bin.
func (b *BinP1D) SumW2() float64 {
	return b.dist.SumW2()
}

// XEdges returns the [low,high] edges of this bin.
func (b *BinP1D) XEdges() Range {
	return b.xrange
}

// XMin returns the lower limit of the bin (inclusive).
func (b *BinP1D) XMin() float64 {
	return b.xrange.Min
}

// XMax returns the upper limit of the bin (exclusive).
func (b *BinP1D) XMax() float64 {
	return b.xrange.Max
}

// XMid returns the geometric center of the bin.
// i.e.: 0.5*(high+low)
func (b *BinP1D) XMid() float64 {
	return 0.5 * (b.xrange.Min + b.xrange.Max)
}

// XWidth returns the (signed) width of the bin
func (b *BinP1D) XWidth() float64 {
	return b.xrange.Max - b.xrange.Min
}

// XFocus returns the mean position in the bin, or the midpoint (if the
// sum of weights for this bin is 0).
func (b *BinP1D) XFocus() float64 {
	if b.SumW() == 0 {
		return b.XMid()
	}
	return b.XMean()
}

// XMean returns the mean X.
func (b *BinP1D) XMean() float64 {
	return b.dist.xMean()
}

// XVariance returns the variance in X.
func (b *BinP1D) XVariance() float64 {
	return b.dist.xVariance()
}

// XStdDev returns the standard deviation in X.
func (b *BinP1D) XStdDev() float64 {
	return b.dist.xStdDev()
}

// XStdErr returns the standard error in X.
func (b *BinP1D) XStdErr() float64 {
	return b.dist.xStdErr()
}

// XRMS returns the RMS in X.
func (b *BinP1D) XRMS() float64 {
	return b.dist.xRMS()
}
