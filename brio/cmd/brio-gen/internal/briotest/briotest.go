// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package briotest // import "go-hep.org/x/hep/brio/cmd/brio-gen/internal/briotest"

type Hist struct {
	Name string
	Data struct {
		X float64
	}
	i int
	u uint

	i8  int8
	i16 int16
	i32 int32
	i64 int64
	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64

	f32 float32
	f64 float64

	c64  complex64
	c128 complex128

	b bool

	arrI8  [2]int8
	sliF64 []float64
	bins   []Bin
	sliPs  []*Bin

	ptr   *float64
	myu8  U8
	myu16 U16
}

func NewHist(
	name string, dataX float64,
	i int,
	u uint,

	i8 int8,
	i16 int16,
	i32 int32,
	i64 int64,
	u8 uint8,
	u16 uint16,
	u32 uint32,
	u64 uint64,

	f32 float32,
	f64 float64,

	c64 complex64,
	c128 complex128,

	b bool,

	arrI8 [2]int8,
	sliF64 []float64,
	bins []Bin,
	sliPs []*Bin,

	ptr *float64,
	myu8 U8,
	myu16 U16,
) *Hist {
	return &Hist{
		Name: name,
		Data: struct {
			X float64
		}{X: dataX},
		i: i,
		u: u,

		i8:  i8,
		i16: i16,
		i32: i32,
		i64: i64,
		u8:  u8,
		u16: u16,
		u32: u32,
		u64: u64,

		f32: f32,
		f64: f64,

		c64:  c64,
		c128: c128,

		b: b,

		arrI8:  arrI8,
		sliF64: sliF64,
		bins:   bins,
		sliPs:  sliPs,

		ptr:   ptr,
		myu8:  myu8,
		myu16: myu16,
	}
}

type Bin struct {
	x, y float64
}

func NewBin(x, y float64) *Bin { return &Bin{x: x, y: y} }

type U8 uint8
type U16 uint16
