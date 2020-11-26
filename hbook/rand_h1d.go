// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"sort"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

// Rand1D represents a 1D distribution created from a hbook.H1D histogram.
type Rand1D struct {
	y    *distuv.Uniform
	bins []Bin1D
	cdf  []float64
}

// NewRand1D creates a Rand1D from the provided histogram.
// If src is nil, the global x/exp/rand source will be used.
func NewRand1D(h *H1D, src rand.Source) Rand1D {
	var (
		sum  float64
		n    = len(h.Binning.Bins)
		cdf  = make([]float64, 1, n+1)
		bins = make([]Bin1D, n, n+1)
	)

	copy(bins, h.Binning.Bins)
	for _, xbin := range bins {
		sum += xbin.SumW()
		cdf = append(cdf, sum)
	}
	cdf[len(cdf)-1] = sum
	norm := 1 / sum
	for i := range cdf {
		cdf[i] *= norm
	}
	// append a bin with similar shape to ease boundary case.
	bins = append(bins, bins[n-1])
	bins[n].Range.Min = bins[n-1].Range.Max
	bins[n].Range.Max = bins[n-1].Range.Max + bins[n-1].Range.Width()

	return Rand1D{
		y:    &distuv.Uniform{Min: 0, Max: 1, Src: src},
		bins: bins,
		cdf:  cdf,
	}
}

func (d *Rand1D) search(v float64) int {
	return sort.Search(len(d.cdf)-1, func(i int) bool { return v < d.cdf[i+1] })
}

// Rand returns a random sample drawn from the distribution.
func (d *Rand1D) Rand() float64 {
	var (
		y = d.y.Rand()
		i = d.search(y)
		x = d.bins[i].XEdges().Min
	)
	if y > d.cdf[i] {
		xbin := d.bins[i+1]
		cdf1 := d.cdf[i+1]
		cdf0 := d.cdf[i]
		x += xbin.XWidth() * (y - cdf0) / (cdf1 - cdf0)
	}
	return x
}

// CDF computes the value of the cumulative density function at x.
func (d *Rand1D) CDF(x float64) float64 {
	i := Bin1Ds(d.bins).IndexOf(x)
	switch i {
	case UnderflowBin1D:
		return 0
	case OverflowBin1D:
		return 1
	default:
		return d.cdf[i]
	}
}
