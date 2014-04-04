package dao

import (
	"math"
)

// H1D is a 1-dim histogram with weighted entries.
type H1D struct {
	bins    []Bin1D // in-range bins
	allbins []Bin1D // in-range bins and under/over-flow bins
	axis    Axis
	entries int64      // number of entries for this histogram
	ann     Annotation // Annotations for this histogram (title, labels,...)
}

// NewH1D returns a 1-dim histogram with nbins bins between low and high.
func NewH1D(nbins int, low, high float64) *H1D {
	h := &H1D{
		bins:    nil,
		allbins: make([]Bin1D, nbins+2),
		axis:    NewEvenBinAxis(nbins, low, high),
		entries: 0,
		ann:     make(Annotation),
	}
	h.bins = h.allbins[2:]
	return h
}

// Name returns the name of this histogram, if any
func (h *H1D) Name() string {
	n := h.ann["name"].(string)
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

// Axis returns the axis of this histgram.
func (h *H1D) Axis() Axis {
	return h.axis
}

// Entries returns the number of entries in this histogram
func (h *H1D) Entries() int64 {
	return h.entries
}

// Fill fills this histogram with x and weight w.
func (h *H1D) Fill(x, w float64) {
	//fmt.Printf("H1D.fill(x=%v, w=%v)...\n", x, w)
	idx := h.axis.CoordToIndex(x)
	switch idx {
	case UnderflowBin:
		h.allbins[0].fill(x, w)
	case OverflowBin:
		h.allbins[1].fill(x, w)
	default:
		h.bins[idx].fill(x, w)
	}
	h.entries += 1
	//fmt.Printf("H1D.fill(x=%v, w=%v)...[done]\n", x, w)
}

// Content returns the content of the idx-th bin.
func (h *H1D) Content(idx int) float64 {
	return h.bins[idx].sw
}

// Mean returns the mean of this histogram.
func (h *H1D) Mean() float64 {
	summeans := 0.0
	sumweights := 0.0
	idx := 0
	for idx = 0; idx < len(h.bins); idx++ {
		summeans = summeans + h.bins[idx].swc
		sumweights = sumweights + h.bins[idx].sw
	}
	return summeans / sumweights
}

// RMS returns the root mean squared of this histogram.
func (h *H1D) RMS() float64 {
	summeans := 0.0
	summean2 := 0.0
	sumweights := 0.0
	idx := 0
	for idx = 0; idx < len(h.bins); idx++ {
		summeans = summeans + h.bins[idx].swc
		sumweights = sumweights + h.bins[idx].sw
		if h.bins[idx].sw != 0. {
			summean2 = summean2 + h.bins[idx].swc*h.bins[idx].swc/h.bins[idx].sw
		}
	}
	invw := 1. / sumweights
	return math.Sqrt(invw * (summean2 - (summeans*summeans)*invw))
}

// Max returns the maximum y value of this histogram.
func (h *H1D) Max() float64 {
	ymax := math.Inf(-1)
	for idx := range h.bins {
		c := h.bins[idx].sw
		if c > ymax {
			ymax = c
		}
	}
	return ymax
}

// Min returns the minimum y value of this histogram.
func (h *H1D) Min() float64 {
	ymin := math.Inf(1)
	for idx := range h.bins {
		c := h.bins[idx].sw
		if c < ymin {
			ymin = c
		}
	}
	return ymin
}

// check various interfaces
var _ Object = (*H1D)(nil)
var _ Histogram = (*H1D)(nil)

// EOF
