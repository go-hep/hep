// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "math"

// dist0D is a 0-dim distribution.
type dist0D struct {
	n     int64   // number of entries
	sumW  float64 // sum of weights
	sumW2 float64 // sum of squared weights
}

// Rank returns the number of dimensions of the distribution.
func (*dist0D) Rank() int {
	return 1
}

// Entries returns the number of Entries in the distribution.
func (d *dist0D) Entries() int64 {
	return d.n
}

// EffEntries returns the number of weighted entries, such as:
//  (\sum w)^2 / \sum w^2
func (d *dist0D) EffEntries() float64 {
	if d.sumW2 == 0 {
		return 0
	}
	return d.sumW * d.sumW / d.sumW2
}

// SumW returns the sum of weights of the distribution.
func (d *dist0D) SumW() float64 {
	return d.sumW
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist0D) SumW2() float64 {
	return d.sumW2
}

// errW returns the absolute error on sumW()
func (d *dist0D) errW() float64 {
	return math.Sqrt(d.SumW2())
}

// relErrW returns the relative error on sumW()
func (d *dist0D) relErrW() float64 {
	// FIXME(sbinet) check for low stats ?
	return d.errW() / d.SumW()
}

func (d *dist0D) fill(w float64) {
	d.n++
	d.sumW += w
	d.sumW2 += w * w
}

func (d *dist0D) scaleW(f float64) {
	d.sumW *= f
	d.sumW2 *= f * f
}

// dist1D is a 1-dim distribution.
type dist1D struct {
	dist   dist0D  // weight moments
	sumWX  float64 // 1st order weighted x moment
	sumWX2 float64 // 2nd order weighted x moment
}

// Rank returns the number of dimensions of the distribution.
func (*dist1D) Rank() int {
	return 1
}

// Entries returns the number of Entries in the distribution.
func (d *dist1D) Entries() int64 {
	return d.dist.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *dist1D) EffEntries() float64 {
	return d.dist.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *dist1D) SumW() float64 {
	return d.dist.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist1D) SumW2() float64 {
	return d.dist.SumW2()
}

// SumWX returns the 1st order weighted x moment
func (d *dist1D) SumWX() float64 {
	return d.sumWX
}

// SumWX2 returns the 2nd order weighted x moment
func (d *dist1D) SumWX2() float64 {
	return d.sumWX2
}

// errW returns the absolute error on sumW()
func (d *dist1D) errW() float64 {
	return d.dist.errW()
}

// relErrW returns the relative error on sumW()
func (d *dist1D) relErrW() float64 {
	return d.dist.relErrW()
}

// mean returns the weighted mean of the distribution
func (d *dist1D) mean() float64 {
	// FIXME(sbinet): check for low stats?
	return d.sumWX / d.SumW()
}

// variance returns the weighted variance of the distribution, defined as:
//  sig2 = ( \sum(wx^2) * \sum(w) - \sum(wx)^2 ) / ( \sum(w)^2 - \sum(w^2) )
// see: https://en.wikipedia.org/wiki/Weighted_arithmetic_mean
func (d *dist1D) variance() float64 {
	// FIXME(sbinet): check for low stats?
	num := d.sumWX2*d.SumW() - math.Pow(d.sumWX, 2)
	den := math.Pow(d.SumW(), 2) - d.SumW2()
	v := num / den
	return math.Abs(v)
}

// stdDev returns the weighted standard deviation of the distribution
func (d *dist1D) stdDev() float64 {
	return math.Sqrt(d.variance())
}

// stdErr returns the weighted standard error of the distribution
func (d *dist1D) stdErr() float64 {
	// FIXME(sbinet): check for low stats?
	// TODO(sbinet): unbiased should check that Neff>1 and divide by N-1?
	return math.Sqrt(d.variance() / d.EffEntries())
}

// rms returns the weighted RMS of the distribution, defined as:
//  rms = \sqrt{\sum{w . x^2} / \sum{w}}
func (d *dist1D) rms() float64 {
	// FIXME(sbinet): check for low stats?
	meansq := d.sumWX2 / d.SumW()
	return math.Sqrt(meansq)
}

func (d *dist1D) fill(x, w float64) {
	d.dist.fill(w)
	d.sumWX += w * x
	d.sumWX2 += w * x * x
}

func (d *dist1D) scaleW(f float64) {
	d.dist.scaleW(f)
	d.sumWX *= f
	d.sumWX2 *= f * f
}

func (d *dist1D) scaleX(f float64) {
	d.sumWX *= f
	d.sumWX2 *= f * f
}

// dist2D is a 2-dim distribution.
type dist2D struct {
	x      dist1D  // x moments
	y      dist1D  // y moments
	sumWXY float64 // 2nd-order cross-term
}

// Rank returns the number of dimensions of the distribution.
func (*dist2D) Rank() int {
	return 2
}

// Entries returns the number of entries in the distribution.
func (d *dist2D) Entries() int64 {
	return d.x.Entries()
}

// EffEntries returns the effective number of entries in the distribution.
func (d *dist2D) EffEntries() float64 {
	return d.x.EffEntries()
}

// SumW returns the sum of weights of the distribution.
func (d *dist2D) SumW() float64 {
	return d.x.SumW()
}

// SumW2 returns the sum of squared weights of the distribution.
func (d *dist2D) SumW2() float64 {
	return d.x.SumW2()
}

func (d *dist2D) sumWX() float64 {
	return d.x.SumW()
}

func (d *dist2D) sumWX2() float64 {
	return d.x.SumW2()
}

func (d *dist2D) sumWY() float64 {
	return d.y.SumW()
}

func (d *dist2D) sumWY2() float64 {
	return d.y.SumW2()
}

// errW returns the absolute error on sumW()
func (d *dist2D) errW() float64 {
	return d.x.errW()
}

// relErrW returns the relative error on sumW()
func (d *dist2D) relErrW() float64 {
	return d.x.relErrW()
}

// meanX returns the weighted mean of the distribution
func (d *dist2D) meanX() float64 {
	return d.x.mean()
}

// meanY returns the weighted mean of the distribution
func (d *dist2D) meanY() float64 {
	return d.y.mean()
}

// varianceX returns the weighted variance of the distribution
func (d *dist2D) varianceX() float64 {
	return d.x.variance()
}

// varianceY returns the weighted variance of the distribution
func (d *dist2D) varianceY() float64 {
	return d.y.variance()
}

// stdDevX returns the weighted standard deviation of the distribution
func (d *dist2D) stdDevX() float64 {
	return d.x.stdDev()
}

// stdDevY returns the weighted standard deviation of the distribution
func (d *dist2D) stdDevY() float64 {
	return d.y.stdDev()
}

// stdErrX returns the weighted standard error of the distribution
func (d *dist2D) stdErrX() float64 {
	return d.x.stdErr()
}

// stdErrY returns the weighted standard error of the distribution
func (d *dist2D) stdErrY() float64 {
	return d.y.stdErr()
}

// rmsX returns the weighted RMS of the distribution
func (d *dist2D) rmsX() float64 {
	return d.x.rms()
}

// rmsY returns the weighted RMS of the distribution
func (d *dist2D) rmsY() float64 {
	return d.y.rms()
}

func (d *dist2D) fill(x, y, w float64) {
	d.x.fill(x, w)
	d.y.fill(y, w)
	d.sumWXY += w * x * y
}

func (d *dist2D) scaleW(f float64) {
	d.x.scaleW(f)
	d.y.scaleW(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleX(f float64) {
	d.x.scaleX(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleY(f float64) {
	d.y.scaleX(f)
	d.sumWXY *= f
}

func (d *dist2D) scaleXY(fx, fy float64) {
	d.scaleX(fx)
	d.scaleY(fy)
}
