// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

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
