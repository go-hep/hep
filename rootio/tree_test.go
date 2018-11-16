// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// disable test on windows because of symlinks
// +build !windows

package rootio

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFlatTree(t *testing.T) {
	t.Parallel()

	f, err := Open("testdata/small-flat-tree.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := obj.(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "my tree title"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.want {
			t.Fatalf("%v: got=[%v]. want=[%v]", table.test, table.value, table.want)
		}
	}

	entries := tree.Entries()
	if got, want := entries, int64(100); got != want {
		t.Fatalf("tree.Entries: got=%v. want=%v", got, want)
	}

	if got, want := tree.TotBytes(), int64(61368); got != want {
		t.Fatalf("tree.totbytes: got=%v. want=%v", got, want)
	}

	if got, want := tree.ZipBytes(), int64(8544); got != want {
		t.Fatalf("tree.zipbytes: got=%v. want=%v", got, want)
	}
}

type EventType struct {
	Evt EventData `rootio:"evt"`
}

type Vec3 struct {
	X int32   `rootio:"Px"`
	Y float64 `rootio:"Py"`
	Z int32   `rootio:"Pz"`
}

type EventData struct {
	Beg    string      `rootio:"Beg"`
	I16    int16       `rootio:"Int16"`
	I32    int32       `rootio:"Int32"`
	I64    int64       `rootio:"Int64"`
	U16    uint16      `rootio:"UInt16"`
	U32    uint32      `rootio:"UInt32"`
	U64    uint64      `rootio:"UInt64"`
	F32    float32     `rootio:"Float32"`
	F64    float64     `rootio:"Float64"`
	Str    string      `rootio:"Str"`
	Vec    Vec3        `rootio:"P3"`
	ArrI16 [10]int16   `rootio:"ArrayI16"`
	ArrI32 [10]int32   `rootio:"ArrayI32"`
	ArrI64 [10]int64   `rootio:"ArrayI64"`
	ArrU16 [10]uint16  `rootio:"ArrayU16"`
	ArrU32 [10]uint32  `rootio:"ArrayU32"`
	ArrU64 [10]uint64  `rootio:"ArrayU64"`
	ArrF32 [10]float32 `rootio:"ArrayF32"`
	ArrF64 [10]float64 `rootio:"ArrayF64"`
	N      int32       `rootio:"N"`
	SliI16 []int16     `rootio:"SliceI16"`
	SliI32 []int32     `rootio:"SliceI32"`
	SliI64 []int64     `rootio:"SliceI64"`
	SliU16 []uint16    `rootio:"SliceU16"`
	SliU32 []uint32    `rootio:"SliceU32"`
	SliU64 []uint64    `rootio:"SliceU64"`
	SliF32 []float32   `rootio:"SliceF32"`
	SliF64 []float64   `rootio:"SliceF64"`
	StdStr string      `rootio:"StdStr"`
	VecI16 []int16     `rootio:"StlVecI16"`
	VecI32 []int32     `rootio:"StlVecI32"`
	VecI64 []int64     `rootio:"StlVecI64"`
	VecU16 []uint16    `rootio:"StlVecU16"`
	VecU32 []uint32    `rootio:"StlVecU32"`
	VecU64 []uint64    `rootio:"StlVecU64"`
	VecF32 []float32   `rootio:"StlVecF32"`
	VecF64 []float64   `rootio:"StlVecF64"`
	VecStr []string    `rootio:"StlVecStr"`
	End    string      `rootio:"End"`
}

func (EventType) want(i int64) EventType {
	var data EventType
	data.Evt.I16 = int16(i)
	data.Evt.I32 = int32(i)
	data.Evt.I64 = int64(i)
	data.Evt.U16 = uint16(i)
	data.Evt.U32 = uint32(i)
	data.Evt.U64 = uint64(i)
	data.Evt.F32 = float32(i)
	data.Evt.F64 = float64(i)
	data.Evt.Str = fmt.Sprintf("evt-%03d", i)
	data.Evt.Vec = Vec3{
		X: int32(i - 1),
		Y: float64(i),
		Z: int32(i - 1),
	}
	data.Evt.StdStr = fmt.Sprintf("std-%03d", i)
	for ii := range data.Evt.ArrI32 {
		data.Evt.ArrI16[ii] = int16(i)
		data.Evt.ArrI32[ii] = int32(i)
		data.Evt.ArrI64[ii] = int64(i)
		data.Evt.ArrU16[ii] = uint16(i)
		data.Evt.ArrU32[ii] = uint32(i)
		data.Evt.ArrU64[ii] = uint64(i)
		data.Evt.ArrF32[ii] = float32(i)
		data.Evt.ArrF64[ii] = float64(i)
	}
	data.Evt.N = int32(i) % 10
	data.Evt.SliI16 = make([]int16, int(data.Evt.N))
	data.Evt.SliI32 = make([]int32, int(data.Evt.N))
	data.Evt.SliI64 = make([]int64, int(data.Evt.N))
	data.Evt.SliU16 = make([]uint16, int(data.Evt.N))
	data.Evt.SliU32 = make([]uint32, int(data.Evt.N))
	data.Evt.SliU64 = make([]uint64, int(data.Evt.N))
	data.Evt.SliF32 = make([]float32, int(data.Evt.N))
	data.Evt.SliF64 = make([]float64, int(data.Evt.N))
	for ii := 0; ii < int(data.Evt.N); ii++ {
		data.Evt.SliI16[ii] = int16(i)
		data.Evt.SliI32[ii] = int32(i)
		data.Evt.SliI64[ii] = int64(i)
		data.Evt.SliU16[ii] = uint16(i)
		data.Evt.SliU32[ii] = uint32(i)
		data.Evt.SliU64[ii] = uint64(i)
		data.Evt.SliF32[ii] = float32(i)
		data.Evt.SliF64[ii] = float64(i)
	}

	data.Evt.Beg = fmt.Sprintf("beg-%03d", i)
	data.Evt.VecI16 = make([]int16, int(data.Evt.N))
	data.Evt.VecI32 = make([]int32, int(data.Evt.N))
	data.Evt.VecI64 = make([]int64, int(data.Evt.N))
	data.Evt.VecU16 = make([]uint16, int(data.Evt.N))
	data.Evt.VecU32 = make([]uint32, int(data.Evt.N))
	data.Evt.VecU64 = make([]uint64, int(data.Evt.N))
	data.Evt.VecF32 = make([]float32, int(data.Evt.N))
	data.Evt.VecF64 = make([]float64, int(data.Evt.N))
	data.Evt.VecStr = make([]string, int(data.Evt.N))
	for ii := 0; ii < int(data.Evt.N); ii++ {
		data.Evt.VecI16[ii] = int16(i)
		data.Evt.VecI32[ii] = int32(i)
		data.Evt.VecI64[ii] = int64(i)
		data.Evt.VecU16[ii] = uint16(i)
		data.Evt.VecU32[ii] = uint32(i)
		data.Evt.VecU64[ii] = uint64(i)
		data.Evt.VecF32[ii] = float32(i)
		data.Evt.VecF64[ii] = float64(i)
		data.Evt.VecStr[ii] = fmt.Sprintf("vec-%03d", i)
	}
	data.Evt.End = fmt.Sprintf("end-%03d", i)
	return data
}

func TestEventTree(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		name  string
		fname string
	}{
		{
			name:  "nosplit",
			fname: "testdata/small-evnt-tree-nosplit.root",
		},
		{
			name:  "fullsplit",
			fname: "testdata/small-evnt-tree-fullsplit.root",
		},
		//		{
		//			name:  "nosplit-xrootd",
		//			fname: XrdRemote("testdata/small-evnt-tree-nosplit.root"),
		//		},
		//		{
		//			name:  "fullsplit-xrootd",
		//			fname: XrdRemote("testdata/small-evnt-tree-fullsplit.root"),
		//		},
	} {
		testEventTree(t, test.name, test.fname)
	}
}

func testEventTree(t *testing.T, name, fname string) {
	f, err := Open(fname)
	if err != nil {
		t.Errorf("%s: %v", name, err.Error())
		return
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Errorf("%s: %v", name, err)
		return
	}

	tree := obj.(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Errorf("%s: tree.Name: got=%q. want=%q", name, got, want)
		return
	}

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "my tree title"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.want {
			t.Errorf("%s: %v: got=[%v]. want=[%v]", name, table.test, table.value, table.want)
			return
		}
	}

	entries := tree.Entries()
	if got, want := entries, int64(100); got != want {
		t.Errorf("%s: tree.Entries: got=%v. want=%v", name, got, want)
		return
	}

	want := EventType{}.want

	sc, err := NewTreeScanner(tree, &EventType{})
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()
	var d1 EventType
	ievt := 0
	for sc.Next() {
		err := sc.Scan(&d1)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			return
		}
		i := sc.Entry()
		if !reflect.DeepEqual(d1, want(i)) {
			t.Errorf("%s: entry[%d]:\ngot= %#v.\nwant=%#v\n", name, i, d1, want(i))
			return
		}

		var d2 EventType
		err = sc.Scan(&d2)
		if err != nil {
			t.Errorf("%s: %v", name, err)
			return
		}
		if !reflect.DeepEqual(d2, want(i)) {
			t.Errorf("%s: entry[%d]:\ngot= %#v.\nwant=%#v\n", name, i, d2, want(i))
			return
		}
		ievt++
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		t.Errorf("%s: %v", name, err)
		return
	}
	if ievt != int(tree.Entries()) {
		t.Errorf("%s: read %d entries. want=%d", name, ievt, tree.Entries())
		return
	}
}

func TestSimpleTree(t *testing.T) {
	t.Parallel()

	f, err := Open("testdata/simple.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := obj.(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "fake data"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.want {
			t.Fatalf("%v: got=[%v]. want=[%v]", table.test, table.value, table.want)
		}
	}

	entries := tree.Entries()
	if got, want := entries, int64(4); got != want {
		t.Fatalf("tree.Entries: got=%v. want=%v", got, want)
	}

	if got, want := tree.TotBytes(), int64(288); got != want {
		t.Fatalf("tree.totbytes: got=%v. want=%v", got, want)
	}

	if got, want := tree.ZipBytes(), int64(288); got != want {
		t.Fatalf("tree.zipbytes: got=%v. want=%v", got, want)
	}
}

func TestSimpleTreeOverHTTP(t *testing.T) {
	t.Parallel()

	f, err := Open("https://github.com/go-hep/hep/raw/master/groot/testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := obj.(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}

	for _, table := range []struct {
		test  string
		value string
		want  string
	}{
		{"Name", tree.Name(), "tree"}, // name when created
		{"Title", tree.Title(), "fake data"},
		{"Class", tree.Class(), "TTree"},
	} {
		if table.value != table.want {
			t.Fatalf("%v: got=[%v]. want=[%v]", table.test, table.value, table.want)
		}
	}

	entries := tree.Entries()
	if got, want := entries, int64(4); got != want {
		t.Fatalf("tree.Entries: got=%v. want=%v", got, want)
	}

	if got, want := tree.TotBytes(), int64(288); got != want {
		t.Fatalf("tree.totbytes: got=%v. want=%v", got, want)
	}

	if got, want := tree.ZipBytes(), int64(288); got != want {
		t.Fatalf("tree.zipbytes: got=%v. want=%v", got, want)
	}
}

func TestTreeWithBasketWithTKeyData(t *testing.T) {
	for _, fname := range []string{
		"testdata/PhaseSpaceSimulation.root",
		//		XrdRemote("testdata/PhaseSpaceSimulation.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := Open(fname)
			if err != nil {
				t.Skipf("error: %v", err)
			}
			defer f.Close()

			obj, err := f.Get("PhaseSpaceTree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)
			if got, want := tree.Name(), "PhaseSpaceTree"; got != want {
				t.Fatalf("tree.Name: got=%q. want=%q", got, want)
			}

			entries := tree.Entries()
			if got, want := entries, int64(50000); got != want {
				t.Fatalf("tree.Entries: got=%v. want=%v", got, want)
			}
		})
	}
}

func TestUprootTrees(t *testing.T) {
	type Data struct {
		N     int32      `rootio:"n"`
		B     bool       `rootio:"b"`
		Arrb  [3]bool    `rootio:"ab"`
		Ab    []bool     `rootio:"Ab"`
		I1    int8       `rootio:"i1"`
		Arri1 [3]int8    `rootio:"ai1"`
		Ai1   []int8     `rootio:"Ai1"`
		U1    int8       `rootio:"u1"`
		Arru1 [3]int8    `rootio:"au1"`
		Au1   []int8     `rootio:"Au1"`
		I2    int16      `rootio:"i2"`
		Arri2 [3]int16   `rootio:"ai2"`
		Ai2   []int16    `rootio:"Ai2"`
		U2    int16      `rootio:"u2"`
		Arru2 [3]int16   `rootio:"au2"`
		Au2   []int16    `rootio:"Au2"`
		I4    int32      `rootio:"i4"`
		Arri4 [3]int32   `rootio:"ai4"`
		Ai4   []int32    `rootio:"Ai4"`
		U4    int32      `rootio:"u4"`
		Arru4 [3]int32   `rootio:"au4"`
		Au4   []int32    `rootio:"Au4"`
		I8    int64      `rootio:"i8"`
		Arri8 [3]int64   `rootio:"ai8"`
		Ai8   []int64    `rootio:"Ai8"`
		U8    int64      `rootio:"u8"`
		Arru8 [3]int64   `rootio:"au8"`
		Au8   []int64    `rootio:"Au8"`
		F4    float32    `rootio:"f4"`
		Arrf4 [3]float32 `rootio:"af4"`
		Af4   []float32  `rootio:"Af4"`
		F8    float64    `rootio:"f8"`
		Arrf8 [3]float64 `rootio:"af8"`
		Af8   []float64  `rootio:"Af8"`
		Str   string     `rootio:"str"`
	}

	var want = [...]Data{
		{
			N:     0,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{},
			I1:    -15,
			Arri1: [3]int8{-14, -13, -12},
			Ai1:   []int8{},
			U1:    0,
			Arru1: [3]int8{1, 2, 3},
			Au1:   []int8{},
			I2:    -15,
			Arri2: [3]int16{-14, -13, -12},
			Ai2:   []int16{},
			U2:    0,
			Arru2: [3]int16{1, 2, 3},
			Au2:   []int16{},
			I4:    -15,
			Arri4: [3]int32{-14, -13, -12},
			Ai4:   []int32{},
			U4:    0,
			Arru4: [3]int32{1, 2, 3},
			Au4:   []int32{},
			I8:    -15,
			Arri8: [3]int64{-14, -13, -12},
			Ai8:   []int64{},
			U8:    0,
			Arru8: [3]int64{1, 2, 3},
			Au8:   []int64{},
			F4:    -14.9,
			Arrf4: [3]float32{-13.9, -12.9, -11.9},
			Af4:   []float32{},
			F8:    -14.9,
			Arrf8: [3]float64{-13.9, -12.9, -11.9},
			Af8:   []float64{},
			Str:   "hey-0",
		},
		{
			N:     1,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true},
			I1:    -14,
			Arri1: [3]int8{-13, -12, -11},
			Ai1:   []int8{-15},
			U1:    1,
			Arru1: [3]int8{2, 3, 4},
			Au1:   []int8{0},
			I2:    -14,
			Arri2: [3]int16{-13, -12, -11},
			Ai2:   []int16{-15},
			U2:    1,
			Arru2: [3]int16{2, 3, 4},
			Au2:   []int16{0},
			I4:    -14,
			Arri4: [3]int32{-13, -12, -11},
			Ai4:   []int32{-15},
			U4:    1,
			Arru4: [3]int32{2, 3, 4},
			Au4:   []int32{0},
			I8:    -14,
			Arri8: [3]int64{-13, -12, -11},
			Ai8:   []int64{-15},
			U8:    1,
			Arru8: [3]int64{2, 3, 4},
			Au8:   []int64{0},
			F4:    -13.9,
			Arrf4: [3]float32{-12.9, -11.9, -10.9},
			Af4:   []float32{-15},
			F8:    -13.9,
			Arrf8: [3]float64{-12.9, -11.9, -10.9},
			Af8:   []float64{-15},
			Str:   "hey-1",
		},
		{
			N:     2,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true},
			I1:    -13,
			Arri1: [3]int8{-12, -11, -10},
			Ai1:   []int8{-15, -13},
			U1:    2,
			Arru1: [3]int8{3, 4, 5},
			Au1:   []int8{0, 2},
			I2:    -13,
			Arri2: [3]int16{-12, -11, -10},
			Ai2:   []int16{-15, -13},
			U2:    2,
			Arru2: [3]int16{3, 4, 5},
			Au2:   []int16{0, 2},
			I4:    -13,
			Arri4: [3]int32{-12, -11, -10},
			Ai4:   []int32{-15, -13},
			U4:    2,
			Arru4: [3]int32{3, 4, 5},
			Au4:   []int32{0, 2},
			I8:    -13,
			Arri8: [3]int64{-12, -11, -10},
			Ai8:   []int64{-15, -13},
			U8:    2,
			Arru8: [3]int64{3, 4, 5},
			Au8:   []int64{0, 2},
			F4:    -12.9,
			Arrf4: [3]float32{-11.9, -10.9, -9.9},
			Af4:   []float32{-15, -13.9},
			F8:    -12.9,
			Arrf8: [3]float64{-11.9, -10.9, -9.9},
			Af8:   []float64{-15, -13.9},
			Str:   "hey-2",
		},
		{
			N:     3,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true, true, true},
			I1:    -12,
			Arri1: [3]int8{-11, -10, -9},
			Ai1:   []int8{-15, -13, -11},
			U1:    3,
			Arru1: [3]int8{4, 5, 6},
			Au1:   []int8{0, 2, 4},
			I2:    -12,
			Arri2: [3]int16{-11, -10, -9},
			Ai2:   []int16{-15, -13, -11},
			U2:    3,
			Arru2: [3]int16{4, 5, 6},
			Au2:   []int16{0, 2, 4},
			I4:    -12,
			Arri4: [3]int32{-11, -10, -9},
			Ai4:   []int32{-15, -13, -11},
			U4:    3,
			Arru4: [3]int32{4, 5, 6},
			Au4:   []int32{0, 2, 4},
			I8:    -12,
			Arri8: [3]int64{-11, -10, -9},
			Ai8:   []int64{-15, -13, -11},
			U8:    3,
			Arru8: [3]int64{4, 5, 6},
			Au8:   []int64{0, 2, 4},
			F4:    -11.9,
			Arrf4: [3]float32{-10.9, -9.9, -8.9},
			Af4:   []float32{-15, -13.9, -12.8},
			F8:    -11.9,
			Arrf8: [3]float64{-10.9, -9.9, -8.9},
			Af8:   []float64{-15, -13.9, -12.8},
			Str:   "hey-3",
		},
		{
			N:     4,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true, true, true},
			I1:    -11,
			Arri1: [3]int8{-10, -9, -8},
			Ai1:   []int8{-15, -13, -11, -9},
			U1:    4,
			Arru1: [3]int8{5, 6, 7},
			Au1:   []int8{0, 2, 4, 6},
			I2:    -11,
			Arri2: [3]int16{-10, -9, -8},
			Ai2:   []int16{-15, -13, -11, -9},
			U2:    4,
			Arru2: [3]int16{5, 6, 7},
			Au2:   []int16{0, 2, 4, 6},
			I4:    -11,
			Arri4: [3]int32{-10, -9, -8},
			Ai4:   []int32{-15, -13, -11, -9},
			U4:    4,
			Arru4: [3]int32{5, 6, 7},
			Au4:   []int32{0, 2, 4, 6},
			I8:    -11,
			Arri8: [3]int64{-10, -9, -8},
			Ai8:   []int64{-15, -13, -11, -9},
			U8:    4,
			Arru8: [3]int64{5, 6, 7},
			Au8:   []int64{0, 2, 4, 6},
			F4:    -10.9,
			Arrf4: [3]float32{-9.9, -8.9, -7.9},
			Af4:   []float32{-15, -13.9, -12.8, -11.7},
			F8:    -10.9,
			Arrf8: [3]float64{-9.9, -8.9, -7.9},
			Af8:   []float64{-15, -13.9, -12.8, -11.7},
			Str:   "hey-4",
		},
		{
			N:     0,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{},
			I1:    -10,
			Arri1: [3]int8{-9, -8, -7},
			Ai1:   []int8{},
			U1:    5,
			Arru1: [3]int8{6, 7, 8},
			Au1:   []int8{},
			I2:    -10,
			Arri2: [3]int16{-9, -8, -7},
			Ai2:   []int16{},
			U2:    5,
			Arru2: [3]int16{6, 7, 8},
			Au2:   []int16{},
			I4:    -10,
			Arri4: [3]int32{-9, -8, -7},
			Ai4:   []int32{},
			U4:    5,
			Arru4: [3]int32{6, 7, 8},
			Au4:   []int32{},
			I8:    -10,
			Arri8: [3]int64{-9, -8, -7},
			Ai8:   []int64{},
			U8:    5,
			Arru8: [3]int64{6, 7, 8},
			Au8:   []int64{},
			F4:    -9.9,
			Arrf4: [3]float32{-8.9, -7.9, -6.9},
			Af4:   []float32{},
			F8:    -9.9,
			Arrf8: [3]float64{-8.9, -7.9, -6.9},
			Af8:   []float64{},
			Str:   "hey-5",
		},
		{
			N:     1,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false},
			I1:    -9,
			Arri1: [3]int8{-8, -7, -6},
			Ai1:   []int8{-10},
			U1:    6,
			Arru1: [3]int8{7, 8, 9},
			Au1:   []int8{5},
			I2:    -9,
			Arri2: [3]int16{-8, -7, -6},
			Ai2:   []int16{-10},
			U2:    6,
			Arru2: [3]int16{7, 8, 9},
			Au2:   []int16{5},
			I4:    -9,
			Arri4: [3]int32{-8, -7, -6},
			Ai4:   []int32{-10},
			U4:    6,
			Arru4: [3]int32{7, 8, 9},
			Au4:   []int32{5},
			I8:    -9,
			Arri8: [3]int64{-8, -7, -6},
			Ai8:   []int64{-10},
			U8:    6,
			Arru8: [3]int64{7, 8, 9},
			Au8:   []int64{5},
			F4:    -8.9,
			Arrf4: [3]float32{-7.9, -6.9, -5.9},
			Af4:   []float32{-10},
			F8:    -8.9,
			Arrf8: [3]float64{-7.9, -6.9, -5.9},
			Af8:   []float64{-10},
			Str:   "hey-6",
		},
		{
			N:     2,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false},
			I1:    -8,
			Arri1: [3]int8{-7, -6, -5},
			Ai1:   []int8{-10, -8},
			U1:    7,
			Arru1: [3]int8{8, 9, 10},
			Au1:   []int8{5, 7},
			I2:    -8,
			Arri2: [3]int16{-7, -6, -5},
			Ai2:   []int16{-10, -8},
			U2:    7,
			Arru2: [3]int16{8, 9, 10},
			Au2:   []int16{5, 7},
			I4:    -8,
			Arri4: [3]int32{-7, -6, -5},
			Ai4:   []int32{-10, -8},
			U4:    7,
			Arru4: [3]int32{8, 9, 10},
			Au4:   []int32{5, 7},
			I8:    -8,
			Arri8: [3]int64{-7, -6, -5},
			Ai8:   []int64{-10, -8},
			U8:    7,
			Arru8: [3]int64{8, 9, 10},
			Au8:   []int64{5, 7},
			F4:    -7.9,
			Arrf4: [3]float32{-6.9, -5.9, -4.9},
			Af4:   []float32{-10, -8.9},
			F8:    -7.9,
			Arrf8: [3]float64{-6.9, -5.9, -4.9},
			Af8:   []float64{-10, -8.9},
			Str:   "hey-7",
		},
		{
			N:     3,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false, false, false},
			I1:    -7,
			Arri1: [3]int8{-6, -5, -4},
			Ai1:   []int8{-10, -8, -6},
			U1:    8,
			Arru1: [3]int8{9, 10, 11},
			Au1:   []int8{5, 7, 9},
			I2:    -7,
			Arri2: [3]int16{-6, -5, -4},
			Ai2:   []int16{-10, -8, -6},
			U2:    8,
			Arru2: [3]int16{9, 10, 11},
			Au2:   []int16{5, 7, 9},
			I4:    -7,
			Arri4: [3]int32{-6, -5, -4},
			Ai4:   []int32{-10, -8, -6},
			U4:    8,
			Arru4: [3]int32{9, 10, 11},
			Au4:   []int32{5, 7, 9},
			I8:    -7,
			Arri8: [3]int64{-6, -5, -4},
			Ai8:   []int64{-10, -8, -6},
			U8:    8,
			Arru8: [3]int64{9, 10, 11},
			Au8:   []int64{5, 7, 9},
			F4:    -6.9,
			Arrf4: [3]float32{-5.9, -4.9, -3.9},
			Af4:   []float32{-10, -8.9, -7.8},
			F8:    -6.9,
			Arrf8: [3]float64{-5.9, -4.9, -3.9000000000000004},
			Af8:   []float64{-10, -8.9, -7.8},
			Str:   "hey-8",
		},
		{
			N:     4,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false, false, false},
			I1:    -6,
			Arri1: [3]int8{-5, -4, -3},
			Ai1:   []int8{-10, -8, -6, -4},
			U1:    9,
			Arru1: [3]int8{10, 11, 12},
			Au1:   []int8{5, 7, 9, 11},
			I2:    -6,
			Arri2: [3]int16{-5, -4, -3},
			Ai2:   []int16{-10, -8, -6, -4},
			U2:    9,
			Arru2: [3]int16{10, 11, 12},
			Au2:   []int16{5, 7, 9, 11},
			I4:    -6,
			Arri4: [3]int32{-5, -4, -3},
			Ai4:   []int32{-10, -8, -6, -4},
			U4:    9,
			Arru4: [3]int32{10, 11, 12},
			Au4:   []int32{5, 7, 9, 11},
			I8:    -6,
			Arri8: [3]int64{-5, -4, -3},
			Ai8:   []int64{-10, -8, -6, -4},
			U8:    9,
			Arru8: [3]int64{10, 11, 12},
			Au8:   []int64{5, 7, 9, 11},
			F4:    -5.9,
			Arrf4: [3]float32{-4.9, -3.9, -2.9},
			Af4:   []float32{-10, -8.9, -7.8, -6.7},
			F8:    -5.9,
			Arrf8: [3]float64{-4.9, -3.9000000000000004, -2.9000000000000004},
			Af8:   []float64{-10, -8.9, -7.8, -6.7},
			Str:   "hey-9",
		},
		{
			N:     0,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{},
			I1:    -5,
			Arri1: [3]int8{-4, -3, -2},
			Ai1:   []int8{},
			U1:    10,
			Arru1: [3]int8{11, 12, 13},
			Au1:   []int8{},
			I2:    -5,
			Arri2: [3]int16{-4, -3, -2},
			Ai2:   []int16{},
			U2:    10,
			Arru2: [3]int16{11, 12, 13},
			Au2:   []int16{},
			I4:    -5,
			Arri4: [3]int32{-4, -3, -2},
			Ai4:   []int32{},
			U4:    10,
			Arru4: [3]int32{11, 12, 13},
			Au4:   []int32{},
			I8:    -5,
			Arri8: [3]int64{-4, -3, -2},
			Ai8:   []int64{},
			U8:    10,
			Arru8: [3]int64{11, 12, 13},
			Au8:   []int64{},
			F4:    -4.9,
			Arrf4: [3]float32{-3.9, -2.9, -1.9},
			Af4:   []float32{},
			F8:    -4.9,
			Arrf8: [3]float64{-3.9000000000000004, -2.9000000000000004, -1.9000000000000004},
			Af8:   []float64{},
			Str:   "hey-10",
		},
		{
			N:     1,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true},
			I1:    -4,
			Arri1: [3]int8{-3, -2, -1},
			Ai1:   []int8{-5},
			U1:    11,
			Arru1: [3]int8{12, 13, 14},
			Au1:   []int8{10},
			I2:    -4,
			Arri2: [3]int16{-3, -2, -1},
			Ai2:   []int16{-5},
			U2:    11,
			Arru2: [3]int16{12, 13, 14},
			Au2:   []int16{10},
			I4:    -4,
			Arri4: [3]int32{-3, -2, -1},
			Ai4:   []int32{-5},
			U4:    11,
			Arru4: [3]int32{12, 13, 14},
			Au4:   []int32{10},
			I8:    -4,
			Arri8: [3]int64{-3, -2, -1},
			Ai8:   []int64{-5},
			U8:    11,
			Arru8: [3]int64{12, 13, 14},
			Au8:   []int64{10},
			F4:    -3.9,
			Arrf4: [3]float32{-2.9, -1.9, -0.9},
			Af4:   []float32{-5},
			F8:    -3.9000000000000004,
			Arrf8: [3]float64{-2.9000000000000004, -1.9000000000000004, -0.9000000000000004},
			Af8:   []float64{-5},
			Str:   "hey-11",
		},
		{
			N:     2,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true},
			I1:    -3,
			Arri1: [3]int8{-2, -1, 0},
			Ai1:   []int8{-5, -3},
			U1:    12,
			Arru1: [3]int8{13, 14, 15},
			Au1:   []int8{10, 12},
			I2:    -3,
			Arri2: [3]int16{-2, -1, 0},
			Ai2:   []int16{-5, -3},
			U2:    12,
			Arru2: [3]int16{13, 14, 15},
			Au2:   []int16{10, 12},
			I4:    -3,
			Arri4: [3]int32{-2, -1, 0},
			Ai4:   []int32{-5, -3},
			U4:    12,
			Arru4: [3]int32{13, 14, 15},
			Au4:   []int32{10, 12},
			I8:    -3,
			Arri8: [3]int64{-2, -1, 0},
			Ai8:   []int64{-5, -3},
			U8:    12,
			Arru8: [3]int64{13, 14, 15},
			Au8:   []int64{10, 12},
			F4:    -2.9,
			Arrf4: [3]float32{-1.9, -0.9, 0.1},
			Af4:   []float32{-5, -3.9},
			F8:    -2.9000000000000004,
			Arrf8: [3]float64{-1.9000000000000004, -0.9000000000000004, 0.09999999999999964},
			Af8:   []float64{-5, -3.9},
			Str:   "hey-12",
		},
		{
			N:     3,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true, true, true},
			I1:    -2,
			Arri1: [3]int8{-1, 0, 1},
			Ai1:   []int8{-5, -3, -1},
			U1:    13,
			Arru1: [3]int8{14, 15, 16},
			Au1:   []int8{10, 12, 14},
			I2:    -2,
			Arri2: [3]int16{-1, 0, 1},
			Ai2:   []int16{-5, -3, -1},
			U2:    13,
			Arru2: [3]int16{14, 15, 16},
			Au2:   []int16{10, 12, 14},
			I4:    -2,
			Arri4: [3]int32{-1, 0, 1},
			Ai4:   []int32{-5, -3, -1},
			U4:    13,
			Arru4: [3]int32{14, 15, 16},
			Au4:   []int32{10, 12, 14},
			I8:    -2,
			Arri8: [3]int64{-1, 0, 1},
			Ai8:   []int64{-5, -3, -1},
			U8:    13,
			Arru8: [3]int64{14, 15, 16},
			Au8:   []int64{10, 12, 14},
			F4:    -1.9,
			Arrf4: [3]float32{-0.9, 0.1, 1.1},
			Af4:   []float32{-5, -3.9, -2.8},
			F8:    -1.9000000000000004,
			Arrf8: [3]float64{-0.9000000000000004, 0.09999999999999964, 1.0999999999999996},
			Af8:   []float64{-5, -3.9, -2.8},
			Str:   "hey-13",
		},
		{
			N:     4,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true, true, true},
			I1:    -1,
			Arri1: [3]int8{0, 1, 2},
			Ai1:   []int8{-5, -3, -1, 1},
			U1:    14,
			Arru1: [3]int8{15, 16, 17},
			Au1:   []int8{10, 12, 14, 16},
			I2:    -1,
			Arri2: [3]int16{0, 1, 2},
			Ai2:   []int16{-5, -3, -1, 1},
			U2:    14,
			Arru2: [3]int16{15, 16, 17},
			Au2:   []int16{10, 12, 14, 16},
			I4:    -1,
			Arri4: [3]int32{0, 1, 2},
			Ai4:   []int32{-5, -3, -1, 1},
			U4:    14,
			Arru4: [3]int32{15, 16, 17},
			Au4:   []int32{10, 12, 14, 16},
			I8:    -1,
			Arri8: [3]int64{0, 1, 2},
			Ai8:   []int64{-5, -3, -1, 1},
			U8:    14,
			Arru8: [3]int64{15, 16, 17},
			Au8:   []int64{10, 12, 14, 16},
			F4:    -0.9,
			Arrf4: [3]float32{0.1, 1.1, 2.1},
			Af4:   []float32{-5, -3.9, -2.8, -1.7},
			F8:    -0.9000000000000004,
			Arrf8: [3]float64{0.09999999999999964, 1.0999999999999996, 2.0999999999999996},
			Af8:   []float64{-5, -3.9, -2.8, -1.7},
			Str:   "hey-14",
		},
		{
			N:     0,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{},
			I1:    0,
			Arri1: [3]int8{1, 2, 3},
			Ai1:   []int8{},
			U1:    15,
			Arru1: [3]int8{16, 17, 18},
			Au1:   []int8{},
			I2:    0,
			Arri2: [3]int16{1, 2, 3},
			Ai2:   []int16{},
			U2:    15,
			Arru2: [3]int16{16, 17, 18},
			Au2:   []int16{},
			I4:    0,
			Arri4: [3]int32{1, 2, 3},
			Ai4:   []int32{},
			U4:    15,
			Arru4: [3]int32{16, 17, 18},
			Au4:   []int32{},
			I8:    0,
			Arri8: [3]int64{1, 2, 3},
			Ai8:   []int64{},
			U8:    15,
			Arru8: [3]int64{16, 17, 18},
			Au8:   []int64{},
			F4:    0.1,
			Arrf4: [3]float32{1.1, 2.1, 3.1},
			Af4:   []float32{},
			F8:    0.09999999999999964,
			Arrf8: [3]float64{1.0999999999999996, 2.0999999999999996, 3.0999999999999996},
			Af8:   []float64{},
			Str:   "hey-15",
		},
		{
			N:     1,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false},
			I1:    1,
			Arri1: [3]int8{2, 3, 4},
			Ai1:   []int8{0},
			U1:    16,
			Arru1: [3]int8{17, 18, 19},
			Au1:   []int8{15},
			I2:    1,
			Arri2: [3]int16{2, 3, 4},
			Ai2:   []int16{0},
			U2:    16,
			Arru2: [3]int16{17, 18, 19},
			Au2:   []int16{15},
			I4:    1,
			Arri4: [3]int32{2, 3, 4},
			Ai4:   []int32{0},
			U4:    16,
			Arru4: [3]int32{17, 18, 19},
			Au4:   []int32{15},
			I8:    1,
			Arri8: [3]int64{2, 3, 4},
			Ai8:   []int64{0},
			U8:    16,
			Arru8: [3]int64{17, 18, 19},
			Au8:   []int64{15},
			F4:    1.1,
			Arrf4: [3]float32{2.1, 3.1, 4.1},
			Af4:   []float32{0},
			F8:    1.0999999999999996,
			Arrf8: [3]float64{2.0999999999999996, 3.0999999999999996, 4.1},
			Af8:   []float64{0},
			Str:   "hey-16",
		},
		{
			N:     2,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false},
			I1:    2,
			Arri1: [3]int8{3, 4, 5},
			Ai1:   []int8{0, 2},
			U1:    17,
			Arru1: [3]int8{18, 19, 20},
			Au1:   []int8{15, 17},
			I2:    2,
			Arri2: [3]int16{3, 4, 5},
			Ai2:   []int16{0, 2},
			U2:    17,
			Arru2: [3]int16{18, 19, 20},
			Au2:   []int16{15, 17},
			I4:    2,
			Arri4: [3]int32{3, 4, 5},
			Ai4:   []int32{0, 2},
			U4:    17,
			Arru4: [3]int32{18, 19, 20},
			Au4:   []int32{15, 17},
			I8:    2,
			Arri8: [3]int64{3, 4, 5},
			Ai8:   []int64{0, 2},
			U8:    17,
			Arru8: [3]int64{18, 19, 20},
			Au8:   []int64{15, 17},
			F4:    2.1,
			Arrf4: [3]float32{3.1, 4.1, 5.1},
			Af4:   []float32{0, 1.1},
			F8:    2.0999999999999996,
			Arrf8: [3]float64{3.0999999999999996, 4.1, 5.1},
			Af8:   []float64{0, 1.1},
			Str:   "hey-17",
		},
		{
			N:     3,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false, false, false},
			I1:    3,
			Arri1: [3]int8{4, 5, 6},
			Ai1:   []int8{0, 2, 4},
			U1:    18,
			Arru1: [3]int8{19, 20, 21},
			Au1:   []int8{15, 17, 19},
			I2:    3,
			Arri2: [3]int16{4, 5, 6},
			Ai2:   []int16{0, 2, 4},
			U2:    18,
			Arru2: [3]int16{19, 20, 21},
			Au2:   []int16{15, 17, 19},
			I4:    3,
			Arri4: [3]int32{4, 5, 6},
			Ai4:   []int32{0, 2, 4},
			U4:    18,
			Arru4: [3]int32{19, 20, 21},
			Au4:   []int32{15, 17, 19},
			I8:    3,
			Arri8: [3]int64{4, 5, 6},
			Ai8:   []int64{0, 2, 4},
			U8:    18,
			Arru8: [3]int64{19, 20, 21},
			Au8:   []int64{15, 17, 19},
			F4:    3.1,
			Arrf4: [3]float32{4.1, 5.1, 6.1},
			Af4:   []float32{0, 1.1, 2.2},
			F8:    3.0999999999999996,
			Arrf8: [3]float64{4.1, 5.1, 6.1},
			Af8:   []float64{0, 1.1, 2.2},
			Str:   "hey-18",
		},
		{
			N:     4,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false, false, false},
			I1:    4,
			Arri1: [3]int8{5, 6, 7},
			Ai1:   []int8{0, 2, 4, 6},
			U1:    19,
			Arru1: [3]int8{20, 21, 22},
			Au1:   []int8{15, 17, 19, 21},
			I2:    4,
			Arri2: [3]int16{5, 6, 7},
			Ai2:   []int16{0, 2, 4, 6},
			U2:    19,
			Arru2: [3]int16{20, 21, 22},
			Au2:   []int16{15, 17, 19, 21},
			I4:    4,
			Arri4: [3]int32{5, 6, 7},
			Ai4:   []int32{0, 2, 4, 6},
			U4:    19,
			Arru4: [3]int32{20, 21, 22},
			Au4:   []int32{15, 17, 19, 21},
			I8:    4,
			Arri8: [3]int64{5, 6, 7},
			Ai8:   []int64{0, 2, 4, 6},
			U8:    19,
			Arru8: [3]int64{20, 21, 22},
			Au8:   []int64{15, 17, 19, 21},
			F4:    4.1,
			Arrf4: [3]float32{5.1, 6.1, 7.1},
			Af4:   []float32{0, 1.1, 2.2, 3.3},
			F8:    4.1,
			Arrf8: [3]float64{5.1, 6.1, 7.1},
			Af8:   []float64{0, 1.1, 2.2, 3.3},
			Str:   "hey-19",
		},
		{
			N:     0,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{},
			I1:    5,
			Arri1: [3]int8{6, 7, 8},
			Ai1:   []int8{},
			U1:    20,
			Arru1: [3]int8{21, 22, 23},
			Au1:   []int8{},
			I2:    5,
			Arri2: [3]int16{6, 7, 8},
			Ai2:   []int16{},
			U2:    20,
			Arru2: [3]int16{21, 22, 23},
			Au2:   []int16{},
			I4:    5,
			Arri4: [3]int32{6, 7, 8},
			Ai4:   []int32{},
			U4:    20,
			Arru4: [3]int32{21, 22, 23},
			Au4:   []int32{},
			I8:    5,
			Arri8: [3]int64{6, 7, 8},
			Ai8:   []int64{},
			U8:    20,
			Arru8: [3]int64{21, 22, 23},
			Au8:   []int64{},
			F4:    5.1,
			Arrf4: [3]float32{6.1, 7.1, 8.1},
			Af4:   []float32{},
			F8:    5.1,
			Arrf8: [3]float64{6.1, 7.1, 8.1},
			Af8:   []float64{},
			Str:   "hey-20",
		},
		{
			N:     1,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true},
			I1:    6,
			Arri1: [3]int8{7, 8, 9},
			Ai1:   []int8{5},
			U1:    21,
			Arru1: [3]int8{22, 23, 24},
			Au1:   []int8{20},
			I2:    6,
			Arri2: [3]int16{7, 8, 9},
			Ai2:   []int16{5},
			U2:    21,
			Arru2: [3]int16{22, 23, 24},
			Au2:   []int16{20},
			I4:    6,
			Arri4: [3]int32{7, 8, 9},
			Ai4:   []int32{5},
			U4:    21,
			Arru4: [3]int32{22, 23, 24},
			Au4:   []int32{20},
			I8:    6,
			Arri8: [3]int64{7, 8, 9},
			Ai8:   []int64{5},
			U8:    21,
			Arru8: [3]int64{22, 23, 24},
			Au8:   []int64{20},
			F4:    6.1,
			Arrf4: [3]float32{7.1, 8.1, 9.1},
			Af4:   []float32{5},
			F8:    6.1,
			Arrf8: [3]float64{7.1, 8.1, 9.1},
			Af8:   []float64{5},
			Str:   "hey-21",
		},
		{
			N:     2,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true},
			I1:    7,
			Arri1: [3]int8{8, 9, 10},
			Ai1:   []int8{5, 7},
			U1:    22,
			Arru1: [3]int8{23, 24, 25},
			Au1:   []int8{20, 22},
			I2:    7,
			Arri2: [3]int16{8, 9, 10},
			Ai2:   []int16{5, 7},
			U2:    22,
			Arru2: [3]int16{23, 24, 25},
			Au2:   []int16{20, 22},
			I4:    7,
			Arri4: [3]int32{8, 9, 10},
			Ai4:   []int32{5, 7},
			U4:    22,
			Arru4: [3]int32{23, 24, 25},
			Au4:   []int32{20, 22},
			I8:    7,
			Arri8: [3]int64{8, 9, 10},
			Ai8:   []int64{5, 7},
			U8:    22,
			Arru8: [3]int64{23, 24, 25},
			Au8:   []int64{20, 22},
			F4:    7.1,
			Arrf4: [3]float32{8.1, 9.1, 10.1},
			Af4:   []float32{5, 6.1},
			F8:    7.1,
			Arrf8: [3]float64{8.1, 9.1, 10.1},
			Af8:   []float64{5, 6.1},
			Str:   "hey-22",
		},
		{
			N:     3,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{true, true, true},
			I1:    8,
			Arri1: [3]int8{9, 10, 11},
			Ai1:   []int8{5, 7, 9},
			U1:    23,
			Arru1: [3]int8{24, 25, 26},
			Au1:   []int8{20, 22, 24},
			I2:    8,
			Arri2: [3]int16{9, 10, 11},
			Ai2:   []int16{5, 7, 9},
			U2:    23,
			Arru2: [3]int16{24, 25, 26},
			Au2:   []int16{20, 22, 24},
			I4:    8,
			Arri4: [3]int32{9, 10, 11},
			Ai4:   []int32{5, 7, 9},
			U4:    23,
			Arru4: [3]int32{24, 25, 26},
			Au4:   []int32{20, 22, 24},
			I8:    8,
			Arri8: [3]int64{9, 10, 11},
			Ai8:   []int64{5, 7, 9},
			U8:    23,
			Arru8: [3]int64{24, 25, 26},
			Au8:   []int64{20, 22, 24},
			F4:    8.1,
			Arrf4: [3]float32{9.1, 10.1, 11.1},
			Af4:   []float32{5, 6.1, 7.2},
			F8:    8.1,
			Arrf8: [3]float64{9.1, 10.1, 11.1},
			Af8:   []float64{5, 6.1, 7.2},
			Str:   "hey-23",
		},
		{
			N:     4,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{true, true, true, true},
			I1:    9,
			Arri1: [3]int8{10, 11, 12},
			Ai1:   []int8{5, 7, 9, 11},
			U1:    24,
			Arru1: [3]int8{25, 26, 27},
			Au1:   []int8{20, 22, 24, 26},
			I2:    9,
			Arri2: [3]int16{10, 11, 12},
			Ai2:   []int16{5, 7, 9, 11},
			U2:    24,
			Arru2: [3]int16{25, 26, 27},
			Au2:   []int16{20, 22, 24, 26},
			I4:    9,
			Arri4: [3]int32{10, 11, 12},
			Ai4:   []int32{5, 7, 9, 11},
			U4:    24,
			Arru4: [3]int32{25, 26, 27},
			Au4:   []int32{20, 22, 24, 26},
			I8:    9,
			Arri8: [3]int64{10, 11, 12},
			Ai8:   []int64{5, 7, 9, 11},
			U8:    24,
			Arru8: [3]int64{25, 26, 27},
			Au8:   []int64{20, 22, 24, 26},
			F4:    9.1,
			Arrf4: [3]float32{10.1, 11.1, 12.1},
			Af4:   []float32{5, 6.1, 7.2, 8.3},
			F8:    9.1,
			Arrf8: [3]float64{10.1, 11.1, 12.1},
			Af8:   []float64{5, 6.1, 7.2, 8.3},
			Str:   "hey-24",
		},
		{
			N:     0,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{},
			I1:    10,
			Arri1: [3]int8{11, 12, 13},
			Ai1:   []int8{},
			U1:    25,
			Arru1: [3]int8{26, 27, 28},
			Au1:   []int8{},
			I2:    10,
			Arri2: [3]int16{11, 12, 13},
			Ai2:   []int16{},
			U2:    25,
			Arru2: [3]int16{26, 27, 28},
			Au2:   []int16{},
			I4:    10,
			Arri4: [3]int32{11, 12, 13},
			Ai4:   []int32{},
			U4:    25,
			Arru4: [3]int32{26, 27, 28},
			Au4:   []int32{},
			I8:    10,
			Arri8: [3]int64{11, 12, 13},
			Ai8:   []int64{},
			U8:    25,
			Arru8: [3]int64{26, 27, 28},
			Au8:   []int64{},
			F4:    10.1,
			Arrf4: [3]float32{11.1, 12.1, 13.1},
			Af4:   []float32{},
			F8:    10.1,
			Arrf8: [3]float64{11.1, 12.1, 13.1},
			Af8:   []float64{},
			Str:   "hey-25",
		},
		{
			N:     1,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false},
			I1:    11,
			Arri1: [3]int8{12, 13, 14},
			Ai1:   []int8{10},
			U1:    26,
			Arru1: [3]int8{27, 28, 29},
			Au1:   []int8{25},
			I2:    11,
			Arri2: [3]int16{12, 13, 14},
			Ai2:   []int16{10},
			U2:    26,
			Arru2: [3]int16{27, 28, 29},
			Au2:   []int16{25},
			I4:    11,
			Arri4: [3]int32{12, 13, 14},
			Ai4:   []int32{10},
			U4:    26,
			Arru4: [3]int32{27, 28, 29},
			Au4:   []int32{25},
			I8:    11,
			Arri8: [3]int64{12, 13, 14},
			Ai8:   []int64{10},
			U8:    26,
			Arru8: [3]int64{27, 28, 29},
			Au8:   []int64{25},
			F4:    11.1,
			Arrf4: [3]float32{12.1, 13.1, 14.1},
			Af4:   []float32{10},
			F8:    11.1,
			Arrf8: [3]float64{12.1, 13.1, 14.1},
			Af8:   []float64{10},
			Str:   "hey-26",
		},
		{
			N:     2,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false},
			I1:    12,
			Arri1: [3]int8{13, 14, 15},
			Ai1:   []int8{10, 12},
			U1:    27,
			Arru1: [3]int8{28, 29, 30},
			Au1:   []int8{25, 27},
			I2:    12,
			Arri2: [3]int16{13, 14, 15},
			Ai2:   []int16{10, 12},
			U2:    27,
			Arru2: [3]int16{28, 29, 30},
			Au2:   []int16{25, 27},
			I4:    12,
			Arri4: [3]int32{13, 14, 15},
			Ai4:   []int32{10, 12},
			U4:    27,
			Arru4: [3]int32{28, 29, 30},
			Au4:   []int32{25, 27},
			I8:    12,
			Arri8: [3]int64{13, 14, 15},
			Ai8:   []int64{10, 12},
			U8:    27,
			Arru8: [3]int64{28, 29, 30},
			Au8:   []int64{25, 27},
			F4:    12.1,
			Arrf4: [3]float32{13.1, 14.1, 15.1},
			Af4:   []float32{10, 11.1},
			F8:    12.1,
			Arrf8: [3]float64{13.1, 14.1, 15.1},
			Af8:   []float64{10, 11.1},
			Str:   "hey-27",
		},
		{
			N:     3,
			B:     true,
			Arrb:  [3]bool{false, true, false},
			Ab:    []bool{false, false, false},
			I1:    13,
			Arri1: [3]int8{14, 15, 16},
			Ai1:   []int8{10, 12, 14},
			U1:    28,
			Arru1: [3]int8{29, 30, 31},
			Au1:   []int8{25, 27, 29},
			I2:    13,
			Arri2: [3]int16{14, 15, 16},
			Ai2:   []int16{10, 12, 14},
			U2:    28,
			Arru2: [3]int16{29, 30, 31},
			Au2:   []int16{25, 27, 29},
			I4:    13,
			Arri4: [3]int32{14, 15, 16},
			Ai4:   []int32{10, 12, 14},
			U4:    28,
			Arru4: [3]int32{29, 30, 31},
			Au4:   []int32{25, 27, 29},
			I8:    13,
			Arri8: [3]int64{14, 15, 16},
			Ai8:   []int64{10, 12, 14},
			U8:    28,
			Arru8: [3]int64{29, 30, 31},
			Au8:   []int64{25, 27, 29},
			F4:    13.1,
			Arrf4: [3]float32{14.1, 15.1, 16.1},
			Af4:   []float32{10, 11.1, 12.2},
			F8:    13.1,
			Arrf8: [3]float64{14.1, 15.1, 16.1},
			Af8:   []float64{10, 11.1, 12.2},
			Str:   "hey-28",
		},
		{
			N:     4,
			B:     false,
			Arrb:  [3]bool{true, false, true},
			Ab:    []bool{false, false, false, false},
			I1:    14,
			Arri1: [3]int8{15, 16, 17},
			Ai1:   []int8{10, 12, 14, 16},
			U1:    29,
			Arru1: [3]int8{30, 31, 32},
			Au1:   []int8{25, 27, 29, 31},
			I2:    14,
			Arri2: [3]int16{15, 16, 17},
			Ai2:   []int16{10, 12, 14, 16},
			U2:    29,
			Arru2: [3]int16{30, 31, 32},
			Au2:   []int16{25, 27, 29, 31},
			I4:    14,
			Arri4: [3]int32{15, 16, 17},
			Ai4:   []int32{10, 12, 14, 16},
			U4:    29,
			Arru4: [3]int32{30, 31, 32},
			Au4:   []int32{25, 27, 29, 31},
			I8:    14,
			Arri8: [3]int64{15, 16, 17},
			Ai8:   []int64{10, 12, 14, 16},
			U8:    29,
			Arru8: [3]int64{30, 31, 32},
			Au8:   []int64{25, 27, 29, 31},
			F4:    14.1,
			Arrf4: [3]float32{15.1, 16.1, 17.1},
			Af4:   []float32{10, 11.1, 12.2, 13.3},
			F8:    14.1,
			Arrf8: [3]float64{15.1, 16.1, 17.1},
			Af8:   []float64{10, 11.1, 12.2, 13.3},
			Str:   "hey-29",
		},
	}

	files, err := filepath.Glob("./testdata/uproot/sample-*.root")
	if err != nil {
		t.Fatal(err)
	}

	for _, fname := range files {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			var d Data
			f, err := Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			obj, err := f.Get("sample")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			s, err := NewScanner(tree, &d)
			if err != nil {
				t.Fatal(err)
			}
			defer s.Close()

			for s.Next() {
				err = s.Scan()
				if err != nil {
					t.Fatalf("error scanning entry %d: %v", s.Entry(), err)
				}
				i := int(s.Entry())
				if !reflect.DeepEqual(d, want[i]) {
					t.Fatalf("entry %d differ.\ngot= %v\nwant=%v\n", s.Entry(), d, want[i])
				}
			}
			err = s.Err()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
