package dao

import (
	"math"
)

type Bin1D struct {
	entries int64
	sw      float64 // sum of weights
	swc     float64 // sum of weights times coord
	sw2     float64 // sum of squared weights
}

func (b *Bin1D) Entries() int64 {
	return b.entries
}

func (b *Bin1D) fill(coord, weight float64) {
	b.entries += 1
	b.sw += weight
	b.swc += weight * coord
	b.sw2 += weight * weight
}

func (b *Bin1D) increment(n int64, height, error, centre float64) {
	b.entries += n
	b.sw += height
	b.swc += centre * height
	b.sw2 += error * error
}

func (b *Bin1D) Set(n int64, height, error, centre float64) {
	b.entries = n
	b.sw = height
	b.swc = centre * height
	b.sw2 = error * error
}

func (b *Bin1D) Err() float64 {
	return math.Sqrt(b.sw2)
}

func (b *Bin1D) Err2() float64 {
	return b.sw2
}

func (b *Bin1D) Centre() float64 {
	if b.sw == 0. {
		return 0.
	}
	return b.swc / b.sw
}

func (b *Bin1D) Scale(factor float64) {
	b.sw *= factor
	b.swc *= factor
	b.sw2 *= factor * factor
}

// EOF
