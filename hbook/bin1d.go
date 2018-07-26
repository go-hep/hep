// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "sort"

// Bin1D models a bin in a 1-dim space.
type Bin1D struct {
	xrange Range
	dist   dist1D
}

// Rank returns the number of dimensions for this bin.
func (Bin1D) Rank() int { return 1 }

func (b *Bin1D) scaleW(f float64) {
	b.dist.scaleW(f)
}

func (b *Bin1D) fill(x, w float64) {
	b.dist.fill(x, w)
}

// Entries returns the number of entries in this bin.
func (b *Bin1D) Entries() int64 {
	return b.dist.Entries()
}

// EffEntries returns the effective number of entries \f$ = (\sum w)^2 / \sum w^2 \f$
func (b *Bin1D) EffEntries() float64 {
	return b.dist.EffEntries()
}

// SumW returns the sum of weights in this bin.
func (b *Bin1D) SumW() float64 {
	return b.dist.SumW()
}

// SumW2 returns the sum of squared weights in this bin.
func (b *Bin1D) SumW2() float64 {
	return b.dist.SumW2()
}

// ErrW returns the absolute error on SumW()
func (b *Bin1D) ErrW() float64 {
	return b.dist.errW()
}

// XEdges returns the [low,high] edges of this bin.
func (b *Bin1D) XEdges() Range {
	return b.xrange
}

// XMin returns the lower limit of the bin (inclusive).
func (b *Bin1D) XMin() float64 {
	return b.xrange.Min
}

// XMax returns the upper limit of the bin (exclusive).
func (b *Bin1D) XMax() float64 {
	return b.xrange.Max
}

// XMid returns the geometric center of the bin.
// i.e.: 0.5*(high+low)
func (b *Bin1D) XMid() float64 {
	return 0.5 * (b.xrange.Min + b.xrange.Max)
}

// XWidth returns the (signed) width of the bin
func (b *Bin1D) XWidth() float64 {
	return b.xrange.Max - b.xrange.Min
}

// XFocus returns the mean position in the bin, or the midpoint (if the
// sum of weights for this bin is 0).
func (b *Bin1D) XFocus() float64 {
	if b.SumW() == 0 {
		return b.XMid()
	}
	return b.XMean()
}

// XMean returns the mean X.
func (b *Bin1D) XMean() float64 {
	return b.dist.mean()
}

// XVariance returns the variance in X.
func (b *Bin1D) XVariance() float64 {
	return b.dist.variance()
}

// XStdDev returns the standard deviation in X.
func (b *Bin1D) XStdDev() float64 {
	return b.dist.stdDev()
}

// XStdErr returns the standard error in X.
func (b *Bin1D) XStdErr() float64 {
	return b.dist.stdErr()
}

// XRMS returns the RMS in X.
func (b *Bin1D) XRMS() float64 {
	return b.dist.rms()
}

// Bin1Ds is a sorted slice of Bin1D implementing sort.Interface.
type Bin1Ds []Bin1D

func (p Bin1Ds) Len() int           { return len(p) }
func (p Bin1Ds) Less(i, j int) bool { return p[i].xrange.Min < p[j].xrange.Min }
func (p Bin1Ds) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// IndexOf returns the index of the Bin1D containing the value v.
// It returns UndeflowBin if v is smaller than the smallest bin value.
// It returns OverflowBin if v is greater than the greatest bin value.
// It returns len(bins) if v falls within a bins gap.
func (p Bin1Ds) IndexOf(v float64) int {
	i := sort.Search(len(p), func(i int) bool { return v < p[i].xrange.Max })
	if i == len(p) {
		return OverflowBin
	}
	rng := p[i].xrange
	if i == 0 && v < rng.Min {
		return UnderflowBin
	}
	if rng.Min <= v && v < rng.Max {
		return i
	}
	return len(p)
}
