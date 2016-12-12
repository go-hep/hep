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
