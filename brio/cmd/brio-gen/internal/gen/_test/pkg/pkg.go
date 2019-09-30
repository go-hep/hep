// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pkg

import "time"

type T1 struct {
	b    bool
	u    uint
	u8   uint8
	u16  uint16
	u32  uint32
	u64  uint64
	i    int
	i8   int8
	i16  int16
	i32  int32
	i64  int64
	f32  float32
	f64  float64
	c64  complex64
	c128 complex128
	str  string
	rune rune
	bs   []byte

	arri64  [10]int64
	arrTime [10]time.Time // implements encoding.Binary(Un)Marshaler
	slii64  []int64
	sliTime []time.Time // implements encoding.Binary(Un)Marshaler
	ptri64  *int64

	t2     T2
	t2s    []T2
	t2ptrs []*T2

	data struct {
		f64 float64
	}

	myU16 U16
}

type T2 struct {
	ID int64
}

type T3 struct {
	Time time.Time // implements encoding.Binary(Un)Marshaler
}

type U16 uint16
