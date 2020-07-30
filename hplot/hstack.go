// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// HStack implements the plot.Plotter interface,
// drawing a stack of histograms.
type HStack struct {
	hs []*H1D

	// LogY allows rendering with a log-scaled Y axis.
	// When enabled, histogram bins with no entries will be discarded from
	// the histogram's DataRange.
	// The lowest Y value for the DataRange will be corrected to leave an
	// arbitrary amount of height for the smallest bin entry so it is visible
	// on the final plot.
	LogY bool

	// Stack specifies how histograms are displayed.
	// Default is to display histograms stacked on top of each other.
	// If not stacked, individual histogram uncertainty bands will be
	// displayed when defined.
	// If stacked, individual uncertainty bands will not be diplayed
	// but the total band can be displayed with the hplot.WithBand(true)
	// option.
	Stack HStackKind

	// Band displays a colored band between the y-min and y-max error bars.
	// Error bars are computed as the bin-by-bin quadratic sum of individual
	// histogram uncertainties.
	Band *BinnedErrBand
}

// HStackKind customizes how a HStack should behave.
type HStackKind int

const (
	// HStackOn instructs HStack to display histograms
	// stacked on top of each other.
	HStackOn HStackKind = iota
	// HStackOff instructs HStack to display histograms
	// with the stack disabled.
	HStackOff
)

func (hsk HStackKind) yoffs(i int, ys []float64, v float64) {
	switch hsk {
	case HStackOn:
		ys[i] += v
	case HStackOff:
		// no-op
	default:
		panic(fmt.Errorf("hplot: unknow HStackKind value %d", hsk))
	}
}

// NewHStack creates a new histogram stack from the provided list of histograms.
// NewHStack panicks if the list of histograms is empty.
// NewHStack panicks if the histograms have different binning.
func NewHStack(histos []*H1D, opts ...Options) *HStack {
	if len(histos) == 0 {
		panic(fmt.Errorf("hplot: not enough histograms to make a stack"))
	}

	cfg := newConfig(opts)

	hstack := &HStack{
		hs:    make([]*H1D, len(histos)),
		Stack: HStackOn,
		LogY:  cfg.log.y,
	}
	copy(hstack.hs, histos)

	ref := hstack.hs[0].Hist.Binning.Bins
	for _, h := range hstack.hs {
		h.LogY = cfg.log.y
		hstack.checkBins(ref, h.Hist.Binning.Bins)
	}

	if cfg.band {
		hstack.Band = NewH1D(hstack.summedH1D(),
			WithBand(true),
			WithLogY(cfg.log.y),
		).Band
	}

	return hstack
}

// summedH1D returns the summed histogram
func (hstack *HStack) summedH1D() *hbook.H1D {
	bookHtot := hstack.hs[0].Hist
	for _, h := range hstack.hs[1:] {
		bookHtot = hbook.AddH1D(bookHtot, h.Hist)
	}
	return bookHtot
}

// DataRange returns the minimum and maximum X and Y values
func (hstack *HStack) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(+1)
	xmax = math.Inf(-1)
	ymin = math.Inf(+1)
	ymax = math.Inf(-1)
	ylow := math.Inf(+1) // ylow will hold the smallest non-zero y value.

	yoffs := make([]float64, len(hstack.hs[0].Hist.Binning.Bins))
	for _, h := range hstack.hs {
		for i, bin := range h.Hist.Binning.Bins {
			sumw := bin.SumW()
			xmax = math.Max(bin.XMax(), xmax)
			xmin = math.Min(bin.XMin(), xmin)
			ymax = math.Max(yoffs[i]+sumw, ymax)
			//if !(yoffs[i]+sumw <= 0 && hstack.LogY){
			ymin = math.Min(yoffs[i]+sumw, ymin)
			//}
			if bin.SumW() != 0 && !(sumw <= 0 && hstack.LogY) {
				ylow = math.Min(bin.SumW(), ylow)
			}
			hstack.Stack.yoffs(i, yoffs, sumw)
		}
	}

	// Handle error bands the error band of individual histograms
	switch hstack.Stack {
	case HStackOff:
		for _, h := range hstack.hs {
			if h.Band != nil {
				_, _, ymin1, ymax1 := h.Band.DataRange()
				for i := range h.Hist.Binning.Bins {
					ymin = math.Min(yoffs[i]+ymin1, ymin)
					ymax = math.Max(yoffs[i]+ymax1, ymax)
				}
			}
		}
	case HStackOn:
		if hstack.Band != nil {
			_, _, ymin1, ymax1 := hstack.Band.DataRange()
			ymin = math.Min(ymin, ymin1)
			ymax = math.Max(ymax, ymax1)
		}
	}

	if hstack.LogY {
		if ymin <= 0 && !math.IsInf(ylow, +1) {
			// Reserve a bit of space for the smallest bin to be displayed still.
			ymin = ylow * 0.5
		}
	}

	return xmin, xmax, ymin, ymax
}

// Plot implements the Plotter interface, drawing a line
// that connects each point in the Line.
func (hstack *HStack) Plot(c draw.Canvas, p *plot.Plot) {
	for i := range hstack.hs {
		hstack.hs[i].LogY = hstack.LogY
	}

	yoffs := make([]float64, len(hstack.hs[0].Hist.Binning.Bins))
	for i, h := range hstack.hs {
		hstack.hplot(c, p, h, yoffs, hstack.Stack, i)
	}
}

func (hstack *HStack) checkBins(refs, bins []hbook.Bin1D) {
	if len(refs) != len(bins) {
		panic("hplot: bins length mismatch")
	}
	for i := range refs {
		ref := refs[i]
		bin := bins[i]
		if ref.Range != bin.Range {
			panic("hplot: bin range mismatch")
		}
	}
}

func (hs *HStack) hplot(c draw.Canvas, p *plot.Plot, h *H1D, yoffs []float64, hsk HStackKind, ih int) {
	trX, trY := p.Transforms(&c)
	var pts []vg.Point
	hist := h.Hist
	bins := h.Hist.Binning.Bins
	nbins := len(bins)

	yfct := func(i int, sumw float64) (ymin, ymax vg.Length) {
		return trY(yoffs[i]), trY(yoffs[i] + sumw)
	}
	if h.LogY {
		yfct = func(i int, sumw float64) (ymin, ymax vg.Length) {
			ymin = c.Min.Y
			if yoffs[i] > 0 {
				ymin = trY(yoffs[i])
			}
			ymax = c.Min.Y
			if sumw+yoffs[i] > 0 {
				ymax = trY(yoffs[i] + sumw)
			}
			return ymin, ymax
		}
	}

	for i, bin := range bins {
		xmin := trX(bin.XMin())
		xmax := trX(bin.XMax())
		sumw := bin.SumW()
		ymin, ymax := yfct(i, sumw)
		switch i {
		case 0:
			pts = append(pts, vg.Point{X: xmin, Y: ymin})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})

		case nbins - 1:
			lft := bins[i-1]
			xlft := trX(lft.XMax())
			_, ylft := yfct(i-1, lft.SumW())
			pts = append(pts, vg.Point{X: xlft, Y: ylft})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymin})

		default:
			lft := bins[i-1]
			xlft := trX(lft.XMax())
			_, ylft := yfct(i-1, lft.SumW())
			pts = append(pts, vg.Point{X: xlft, Y: ylft})
			pts = append(pts, vg.Point{X: xmin, Y: ymax})
			pts = append(pts, vg.Point{X: xmax, Y: ymax})
		}

		if h.GlyphStyle.Radius != 0 {
			x := trX(bin.XMid())
			_, y := yfct(i, bin.SumW())
			c.DrawGlyph(h.GlyphStyle, vg.Point{X: x, Y: y})
		}
	}

	if h.FillColor != nil {
		poly := pts
		for i := range yoffs {
			j := len(yoffs) - 1 - i
			bin := bins[j]
			xmin := trX(bin.XMin())
			xmax := trX(bin.XMax())
			ymin, _ := yfct(j, bin.SumW())
			poly = append(poly, vg.Point{X: xmax, Y: ymin})
			poly = append(poly, vg.Point{X: xmin, Y: ymin})
		}
		c.FillPolygon(h.FillColor, c.ClipPolygonXY(poly))
	}

	// Plot individual histo band when not stacked or total band
	switch hsk {
	case HStackOff:
		if h.Band != nil {
			h.Band.Plot(c, p)
		}
	case HStackOn:
		if ih == len(hs.hs)-1 && hs.Band != nil {
			hs.Band.Plot(c, p)
		}
	}

	c.StrokeLines(h.LineStyle, c.ClipLinesXY(pts)...)

	if h.YErrs != nil {
		yerrs := h.withYErrBars(yoffs)

		yerrs.LineStyle = h.YErrs.LineStyle
		yerrs.CapWidth = h.YErrs.CapWidth

		yerrs.Plot(c, p)
	}

	if h.Infos.Style != HInfoNone {
		fnt, err := vg.MakeFont(DefaultStyle.Fonts.Name, DefaultStyle.Fonts.Tick.Size)
		if err == nil {
			sty := draw.TextStyle{Font: fnt}
			legend := histLegend{
				ColWidth:  DefaultStyle.Fonts.Tick.Size,
				TextStyle: sty,
			}

			for i := uint32(0); i < 32; i++ {
				switch h.Infos.Style & (1 << i) {
				case HInfoEntries:
					legend.Add("Entries", hist.Entries())
				case HInfoMean:
					legend.Add("Mean", hist.XMean())
				case HInfoRMS:
					legend.Add("RMS", hist.XRMS())
				case HInfoStdDev:
					legend.Add("Std Dev", hist.XStdDev())
				default:
				}
			}
			legend.Top = true

			legend.draw(c)
		}
	}

	// handle stack, if any.
	for i, bin := range bins {
		sumw := bin.SumW()
		hsk.yoffs(i, yoffs, sumw)
	}
}

var (
	_ plot.Plotter = (*HStack)(nil)
)
