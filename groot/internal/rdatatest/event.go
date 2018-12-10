// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdatatest

import (
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

// Event is a simple type to exercize streamers generation.
type Event struct {
	Name string
	u8   uint8
	u16  uint16
	u32  uint32
	u64  uint64 `groot:"U64"`
	i8   int8
	i16  int16
	i32  int32
	i64  int64
	f32  float32
	f64  float64

	b  bool
	bb byte

	u8s  []uint8
	u16s []uint16
	u32s []uint32
	u64s []uint64
	i8s  []int8
	i16s []int16
	i32s []int32
	i64s []int64
	f32s []float32
	f64s []float64
	bs   []bool
	bbs  []byte

	arru8s  [10]uint8
	arru16s [10]uint16
	arru32s [10]uint32
	arru64s [10]uint64
	arri8s  [10]int8
	arri16s [10]int16
	arri32s [10]int32
	arri64s [10]int64
	arrf32s [10]float32
	arrf64s [10]float64
	arrbs   [10]bool
	arrbbs  [10]byte

	SliF64 []float64   `groot:"SliceF64"`
	SliStr []string    `groot:"SliceStr"`
	SliHLV []HLV       `groot:"SliceHLV"`
	ArrF64 [10]float64 `groot:"ArrF64"`
}

func (Event) RVersion() int16 { return 1 }

// Particle is a simple type to exercize streamers generation.
type Particle struct {
	name string
	pid  int
	mom  HLV
}

func (Particle) RVersion() int16 { return 1 }

// HLV is a simple type to exercize streamers generation.
type HLV struct {
	px, py, pz, e float64
}

func (HLV) RVersion() int16 { return 1 }

// FIXME(sbinet)
//  - support types that "inherit" from TObject
//  - support types that contain a TList
//  - support types that contain a pointer to TObject
//  - support types that contain a pointer to any-object

type TObject struct {
	rbase.Object
	name string
}

type TList struct {
	rbase.Object
	objs []root.Object
	list rcont.List
}

type TClonesArray struct {
	rbase.Object
	//	clones rcont.ClonesArray // FIXME(sbinet)
}
