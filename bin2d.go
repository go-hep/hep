// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// Bin2D models a bin in a 2-dim space.
type Bin2D struct {
	xrange Range
	yrange Range
	dist   dist2D
}

// Rank returns the number of dimensions for this bin.
func (Bin2D) Rank() int { return 2 }

func (b *Bin2D) scaleW(f float64) {
	b.dist.scaleW(f)
}

func (b *Bin2D) fill(x, y, w float64) {
	b.dist.fill(x, y, w)
}

func (b *Bin2D) Entries() int64 {
	return b.dist.Entries()
}

// EffEntries returns the effective number of entries \f$ = (\sum w)^2 / \sum w^2 \f$
func (b *Bin2D) EffEntries() float64 {
	return b.dist.EffEntries()
}

func (b *Bin2D) SumW() float64 {
	return b.dist.SumW()
}

func (b *Bin2D) SumW2() float64 {
	return b.dist.SumW2()
}

// check Bin2D implements interfaces
var _ Bin = (*Bin2D)(nil)
