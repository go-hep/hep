// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdatatest

import (
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
)

// Event is a simple type to exercize streamers generation.
type Event struct {
	name string `groot:"Name"`
	u8   uint8
	u16  uint16
	u32  uint32
	u64  uint64
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

func (*Event) RVersion() int16 { return 1 }
func (*Event) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.Event" }

// Particle is a simple type to exercize streamers generation.
type Particle struct {
	name string
	pid  int
	mom  HLV
}

func (*Particle) RVersion() int16 { return 1 }
func (*Particle) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.Particle" }

// HLV is a simple type to exercize streamers generation.
type HLV struct {
	px, py, pz, e float64
}

func (*HLV) RVersion() int16 { return 1 }
func (*HLV) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.HLV" }

// Builtins exercizes all simple builtins.
type Builtins struct {
	b bool

	u8  uint8
	u16 uint16
	u32 uint32
	u64 uint64

	i8  int8
	i16 int16
	i32 int32
	i64 int64

	f32 float32
	f64 float64

	// c64  complex64  // FIXME(sbinet)
	// c128 complex128 // FIXME(sbinet)

	name string `groot:"Name"`
}

func (*Builtins) RVersion() int16 { return 1 }
func (*Builtins) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.Builtins" }

func NewBuiltins() *Builtins {
	return &Builtins{
		b: true,

		u8:  8,
		u16: 16,
		u32: 32,
		u64: 64,

		i8:  -8,
		i16: -16,
		i32: -32,
		i64: -64,

		f32: 32.32,
		f64: 64.64,

		// c64:  complex(float32(-1), float32(+2)), // FIXME(sbinet)
		// c128: complex(-2, +3),                   // FIXME(sbinet)

		name: "builtins",
	}
}

// ArrBuiltins exercizes all arrays of simple builtins.
type ArrBuiltins struct {
	b [2]bool

	u8  [2]uint8
	u16 [2]uint16
	u32 [2]uint32
	u64 [2]uint64

	i8  [2]int8
	i16 [2]int16
	i32 [2]int32
	i64 [2]int64

	f32 [2]float32
	f64 [2]float64

	// c64  [2]complex64  // FIXME(sbinet)
	// c128 [2]complex128 // FIXME(sbinet)

	// name [2]string `groot:"Name"` // FIXME(sbinet)
}

func (*ArrBuiltins) RVersion() int16 { return 1 }
func (*ArrBuiltins) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.ArrBuiltins" }

func NewArrBuiltins() *ArrBuiltins {
	return &ArrBuiltins{
		b: [2]bool{true, false},

		u8:  [2]uint8{8, 88},
		u16: [2]uint16{16, 1616},
		u32: [2]uint32{32, 3232},
		u64: [2]uint64{64, 6464},

		i8:  [2]int8{-8, -88},
		i16: [2]int16{-16, -1616},
		i32: [2]int32{-32, -3232},
		i64: [2]int64{-64, -6464},

		f32: [2]float32{32.32, -32.32},
		f64: [2]float64{64.64, +64.64},

		// c64:  complex(float32(-1), float32(+2)), // FIXME(sbinet)
		// c128: complex(-2, +3),                   // FIXME(sbinet)

		// name: [2]string{"builtins", "arrays"}, // FIXME(sbinet)
	}
}

// T1 exercizes a user type containing another user-type.
type T1 struct {
	name string `groot:"Name"`
	hlv  HLV    `groot:"MyHLV"`
}

func (*T1) RVersion() int16 { return 1 }
func (*T1) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.T1" }

func NewT1() *T1 {
	return &T1{
		name: "hello",
		hlv:  HLV{1, 2, 3, 4},
	}
}

// T2 exercizes a user type containing an array of another user-type.
type T2 struct {
	name string `groot:"Name"`
	hlvs [2]HLV `groot:"MyHLVs"`
}

func (*T2) RVersion() int16 { return 1 }
func (*T2) Class() string   { return "go-hep.org/x/hep/groot/internal/rdatatest.T2" }

func NewT2() *T2 {
	return &T2{
		name: "hello",
		hlvs: [2]HLV{{1, 2, 3, 4}, {-1, -2, -3, -4}},
	}
}

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

var (
	_ rbytes.RVersioner = (*Event)(nil)
	_ rbytes.RVersioner = (*HLV)(nil)
	_ rbytes.RVersioner = (*Particle)(nil)

	_ rbytes.RVersioner = (*Builtins)(nil)
	_ rbytes.RVersioner = (*ArrBuiltins)(nil)
	_ rbytes.RVersioner = (*T1)(nil)
	_ rbytes.RVersioner = (*T2)(nil)
)
