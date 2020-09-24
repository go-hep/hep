// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
)

type ScannerData struct {
	B      bool              `groot:"B"`
	Str    string            `groot:"Str"`
	I8     int8              `groot:"I8"`
	I16    int16             `groot:"I16"`
	I32    int32             `groot:"I32"`
	I64    int64             `groot:"I64"`
	U8     uint8             `groot:"U8"`
	U16    uint16            `groot:"U16"`
	U32    uint32            `groot:"U32"`
	U64    uint64            `groot:"U64"`
	F32    float32           `groot:"F32"`
	F64    float64           `groot:"F64"`
	D16    root.Float16      `groot:"D16"`
	D32    root.Double32     `groot:"D32"`
	ArrBs  [10]bool          `groot:"ArrBs[10]"`
	ArrI8  [10]int8          `groot:"ArrI8[10]"`
	ArrI16 [10]int16         `groot:"ArrI16[10]"`
	ArrI32 [10]int32         `groot:"ArrI32[10]"`
	ArrI64 [10]int64         `groot:"ArrI64[10]"`
	ArrU8  [10]uint8         `groot:"ArrU8[10]"`
	ArrU16 [10]uint16        `groot:"ArrU16[10]"`
	ArrU32 [10]uint32        `groot:"ArrU32[10]"`
	ArrU64 [10]uint64        `groot:"ArrU64[10]"`
	ArrF32 [10]float32       `groot:"ArrF32[10]"`
	ArrF64 [10]float64       `groot:"ArrF64[10]"`
	ArrD16 [10]root.Float16  `groot:"ArrD16[10]"`
	ArrD32 [10]root.Double32 `groot:"ArrD32[10]"`
	N      int32             `groot:"N"`
	SliBs  []bool            `groot:"SliBs[N]"`
	SliI8  []int8            `groot:"SliI8[N]"`
	SliI16 []int16           `groot:"SliI16[N]"`
	SliI32 []int32           `groot:"SliI32[N]"`
	SliI64 []int64           `groot:"SliI64[N]"`
	SliU8  []uint8           `groot:"SliU8[N]"`
	SliU16 []uint16          `groot:"SliU16[N]"`
	SliU32 []uint32          `groot:"SliU32[N]"`
	SliU64 []uint64          `groot:"SliU64[N]"`
	SliF32 []float32         `groot:"SliF32[N]"`
	SliF64 []float64         `groot:"SliF64[N]"`
	SliD16 []root.Float16    `groot:"SliD16[N]"`
	SliD32 []root.Double32   `groot:"SliD32[N]"`
}

func (ScannerData) want(i int64) (data ScannerData) {
	data.B = i%2 == 0
	data.Str = fmt.Sprintf("str-%d", i)
	data.I8 = int8(-i)
	data.I16 = int16(-i)
	data.I32 = int32(-i)
	data.I64 = int64(-i)
	data.U8 = uint8(i)
	data.U16 = uint16(i)
	data.U32 = uint32(i)
	data.U64 = uint64(i)
	data.F32 = float32(i)
	data.F64 = float64(i)
	data.D16 = root.Float16(i)
	data.D32 = root.Double32(i)
	for ii := range data.ArrI32 {
		data.ArrBs[ii] = ii == int(i)
		data.ArrI8[ii] = int8(-i)
		data.ArrI16[ii] = int16(-i)
		data.ArrI32[ii] = int32(-i)
		data.ArrI64[ii] = int64(-i)
		data.ArrU8[ii] = uint8(i)
		data.ArrU16[ii] = uint16(i)
		data.ArrU32[ii] = uint32(i)
		data.ArrU64[ii] = uint64(i)
		data.ArrF32[ii] = float32(i)
		data.ArrF64[ii] = float64(i)
		data.ArrD16[ii] = root.Float16(i)
		data.ArrD32[ii] = root.Double32(i)
	}
	data.N = int32(i) % 10
	data.SliBs = make([]bool, int(data.N))
	data.SliI8 = make([]int8, int(data.N))
	data.SliI16 = make([]int16, int(data.N))
	data.SliI32 = make([]int32, int(data.N))
	data.SliI64 = make([]int64, int(data.N))
	data.SliU8 = make([]uint8, int(data.N))
	data.SliU16 = make([]uint16, int(data.N))
	data.SliU32 = make([]uint32, int(data.N))
	data.SliU64 = make([]uint64, int(data.N))
	data.SliF32 = make([]float32, int(data.N))
	data.SliF64 = make([]float64, int(data.N))
	data.SliD16 = make([]root.Float16, int(data.N))
	data.SliD32 = make([]root.Double32, int(data.N))
	for ii := 0; ii < int(data.N); ii++ {
		data.SliBs[ii] = (ii + 1) == int(i)
		data.SliI8[ii] = int8(-i)
		data.SliI16[ii] = int16(-i)
		data.SliI32[ii] = int32(-i)
		data.SliI64[ii] = int64(-i)
		data.SliU8[ii] = uint8(i)
		data.SliU16[ii] = uint16(i)
		data.SliU32[ii] = uint32(i)
		data.SliU64[ii] = uint64(i)
		data.SliF32[ii] = float32(i)
		data.SliF64[ii] = float64(i)
		data.SliD16[ii] = root.Float16(i)
		data.SliD32[ii] = root.Double32(i)
	}
	return data
}
func TestTreeScannerStruct(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			want := ScannerData{}.want

			sc, err := NewTreeScanner(tree, &ScannerData{})
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			var d1 ScannerData
			for sc.Next() {
				err := sc.Scan(&d1)
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(d1, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, d1, want(i))
				}

				var d2 ScannerData
				err = sc.Scan(&d2)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(d2, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, d2, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerStruct(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			var (
				want = ScannerData{}.want
				data ScannerData
			)
			sc, err := NewScanner(tree, &data)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}

				// test a second time
				err = sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i = sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerVars(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)

			want := ScannerData{}.want

			var (
				data  ScannerData
				rvars = ReadVarsFromStruct(&data)
			)
			sc, err := NewScannerVars(tree, rvars...)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestTreeScannerVarsMultipleTimes(t *testing.T) {
	for _, fname := range []string{
		"../testdata/mc_105986.ZZ.root",
		rtests.XrdRemote("testdata/mc_105986.ZZ.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Skip(err)
			}

			obj, err := f.Get("mini")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			for i := 0; i < 10; i++ {
				sc, err := NewTreeScannerVars(tree, ReadVar{Name: "lep_pt"})
				if err != nil {
					t.Fatal(err)
				}
				defer sc.Close()

				for sc.Next() {
					var data []float32
					err := sc.Scan(&data)
					if err != nil {
						t.Fatalf("could not scan data i=%d evt=%v err=%v", i, sc.Entry(), err)
					}
				}
				err = sc.Err()
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestTreeScannerVars(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)
			want := ScannerData{}.want
			rvars := ReadVarsFromStruct(new(ScannerData))
			sc, err := NewTreeScannerVars(tree, rvars...)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			var d1 ScannerData
			for sc.Next() {
				err := sc.Scan(
					&d1.B,
					&d1.Str,
					&d1.I8,
					&d1.I16,
					&d1.I32,
					&d1.I64,
					&d1.U8,
					&d1.U16,
					&d1.U32,
					&d1.U64,
					&d1.F32,
					&d1.F64,
					&d1.D16,
					&d1.D32,
					&d1.ArrBs,
					&d1.ArrI8,
					&d1.ArrI16,
					&d1.ArrI32,
					&d1.ArrI64,
					&d1.ArrU8,
					&d1.ArrU16,
					&d1.ArrU32,
					&d1.ArrU64,
					&d1.ArrF32,
					&d1.ArrF64,
					&d1.ArrD16,
					&d1.ArrD32,
					&d1.N,
					&d1.SliBs,
					&d1.SliI8,
					&d1.SliI16,
					&d1.SliI32,
					&d1.SliI64,
					&d1.SliU8,
					&d1.SliU16,
					&d1.SliU32,
					&d1.SliU64,
					&d1.SliF32,
					&d1.SliF64,
					&d1.SliD16,
					&d1.SliD32,
				)
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(d1, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, d1, want(i))
				}

				var d2 ScannerData
				err = sc.Scan(
					&d2.B,
					&d2.Str,
					&d2.I8,
					&d2.I16,
					&d2.I32,
					&d2.I64,
					&d2.U8,
					&d2.U16,
					&d2.U32,
					&d2.U64,
					&d2.F32,
					&d2.F64,
					&d2.D16,
					&d2.D32,
					&d2.ArrBs,
					&d2.ArrI8,
					&d2.ArrI16,
					&d2.ArrI32,
					&d2.ArrI64,
					&d2.ArrU8,
					&d2.ArrU16,
					&d2.ArrU32,
					&d2.ArrU64,
					&d2.ArrF32,
					&d2.ArrF64,
					&d2.ArrD16,
					&d2.ArrD32,
					&d2.N,
					&d2.SliBs,
					&d2.SliI8,
					&d2.SliI16,
					&d2.SliI32,
					&d2.SliI64,
					&d2.SliU8,
					&d2.SliU16,
					&d2.SliU32,
					&d2.SliU64,
					&d2.SliF32,
					&d2.SliF64,
					&d2.SliD16,
					&d2.SliD32,
				)
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(d2, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, d2, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerVarsMultipleTimes(t *testing.T) {
	for _, fname := range []string{
		"../testdata/mc_105986.ZZ.root",
		rtests.XrdRemote("testdata/mc_105986.ZZ.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Skip(err)
			}

			obj, err := f.Get("mini")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			var pt []float32
			for i := 0; i < 10; i++ {
				sc, err := NewScannerVars(tree, ReadVar{Name: "lep_pt", Value: &pt})
				if err != nil {
					t.Fatal(err)
				}
				defer sc.Close()

				for sc.Next() {
					err := sc.Scan()
					if err != nil {
						t.Error(err)
					}
				}
				err = sc.Err()
				if err != nil {
					t.Error(err)
				}
			}
		})
	}
}

func TestTreeScannerStructWithCounterLeaf(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)

			type Data struct {
				Sli []int32 `groot:"SliI32"`
			}
			var data Data

			want := func(i int64) Data {
				var data Data
				n := int32(i) % 10
				data.Sli = make([]int32, int(n))
				for ii := 0; ii < int(n); ii++ {
					data.Sli[ii] = int32(-i)
				}
				return data
			}

			sc, err := NewTreeScanner(tree, &data)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan(&data)
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerStructWithCounterLeaf(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)

			type Data struct {
				Sli []int32 `groot:"SliI32"`
			}
			var data Data

			want := func(i int64) Data {
				var data Data
				n := int32(i) % 10
				data.Sli = make([]int32, int(n))
				for ii := 0; ii < int(n); ii++ {
					data.Sli[ii] = int32(-i)
				}
				return data
			}

			sc, err := NewScanner(tree, &data)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestTreeScannerVarsWithCounterLeaf(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)

			want := func(i int64) []int32 {
				n := int32(i) % 10
				data := make([]int32, int(n))
				for ii := 0; ii < int(n); ii++ {
					data[ii] = int32(-i)
				}
				return data
			}

			rvar := ReadVar{Name: "SliI32"}
			sc, err := NewTreeScannerVars(tree, rvar)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				var data []int32
				err := sc.Scan(&data)
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerVarsWithCounterLeaf(t *testing.T) {
	for _, fname := range []string{
		"../testdata/x-flat-tree.root",
		rtests.XrdRemote("testdata/x-flat-tree.root"),
	} {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}

			tree := obj.(Tree)

			want := func(i int64) []int32 {
				n := int32(i) % 10
				data := make([]int32, int(n))
				for ii := 0; ii < int(n); ii++ {
					data[ii] = int32(-i)
				}
				return data
			}

			var data []int32
			rvar := ReadVar{Name: "SliI32", Value: &data}
			sc, err := NewScannerVars(tree, rvar)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func TestScannerStructWithStdVectorBool(t *testing.T) {
	files, err := filepath.Glob("../testdata/stdvec-bool-*.root")
	if err != nil {
		t.Fatal(err)
	}

	for _, fname := range files {
		t.Run(fname, func(t *testing.T) {
			t.Parallel()

			f, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(Tree)

			type Data struct {
				Bool    bool     `groot:"Bool"`
				ArrBool [10]bool `groot:"ArrayBool"`
				N       int32    `groot:"N"`
				SliBool []bool   `groot:"SliceBool[N]"`
				StlBool []bool   `groot:"StlVecBool"`
			}
			type Event struct {
				Data Data `groot:"evt"`
			}

			want := func(i int64) Event {
				var data Data
				data.Bool = i%2 == 0
				for ii := range data.ArrBool {
					data.ArrBool[ii] = i%2 == 0
				}
				data.N = int32(i) % 10
				switch i {
				case 0:
					data.SliBool = nil
					data.StlBool = nil
				default:
					data.SliBool = make([]bool, int(data.N))
					data.StlBool = make([]bool, int(data.N))
				}
				for ii := 0; ii < int(data.N); ii++ {
					data.SliBool[ii] = i%2 == 0
					data.StlBool[ii] = i%2 == 0
				}
				return Event{data}
			}

			var data Event
			sc, err := NewScanner(tree, &data)
			if err != nil {
				t.Fatal(err)
			}
			defer sc.Close()
			for sc.Next() {
				err := sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i := sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}

				// test a second time
				err = sc.Scan()
				if err != nil {
					t.Fatal(err)
				}
				i = sc.Entry()
				if !reflect.DeepEqual(data, want(i)) {
					t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
				}
			}
			if err := sc.Err(); err != nil && err != io.EOF {
				t.Fatal(err)
			}
		})
	}
}

func BenchmarkTreeScannerStruct(b *testing.B) {
	f, err := riofs.Open("../testdata/x-flat-tree.root")
	if err != nil {
		b.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	type Data struct {
		F64 float64 `groot:"F64"`
	}

	var data Data
	s, err := NewTreeScanner(tree, &data)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var sum float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.SeekEntry(0)
		for s.Next() {
			err = s.Scan(&data)
			if err != nil {
				b.Fatal(err)
			}
			sum += data.F64
		}
	}
}

func BenchmarkScannerStruct(b *testing.B) {
	f, err := riofs.Open("../testdata/x-flat-tree.root")
	if err != nil {
		b.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	type Data struct {
		F64 float64 `groot:"F64"`
	}

	var data Data
	s, err := NewScanner(tree, &data)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var sum float64
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.SeekEntry(0)
		for s.Next() {
			err = s.Scan()
			if err != nil {
				b.Fatal(err)
			}
			sum += data.F64
		}
	}
}

func BenchmarkTreeScannerVars(b *testing.B) {
	f, err := riofs.Open("../testdata/x-flat-tree.root")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}

	tree := obj.(Tree)

	rvars := []ReadVar{
		{Name: "F64"},
	}
	s, err := NewTreeScannerVars(tree, rvars...)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var data ScannerData
	var sum float64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.SeekEntry(0)
		for s.Next() {
			err := s.Scan(&data.F64)
			if err != nil {
				b.Fatal(err)
			}
			sum += data.F64
		}
	}
}

func BenchmarkScannerVars(b *testing.B) {
	f, err := riofs.Open("../testdata/x-flat-tree.root")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}

	tree := obj.(Tree)

	var f64 float64
	rvars := []ReadVar{
		{Name: "F64", Value: &f64},
	}
	s, err := NewScannerVars(tree, rvars...)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var sum float64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.SeekEntry(0)
		for s.Next() {
			err := s.Scan()
			if err != nil {
				b.Fatal(err)
			}
			sum += f64
		}
	}
}

func BenchmarkTreeScannerVarsBigFileScalar(b *testing.B) {
	f, err := riofs.Open("../testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	sc, err := NewTreeScannerVars(tree, ReadVar{Name: "mcWeight"})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sc.SeekEntry(0)
		for sc.Next() {
			var data float32
			err := sc.Scan(&data)
			if err != nil {
				b.Error(err)
			}
			sum += data
		}
	}
}

func BenchmarkScannerVarsBigFileScalar(b *testing.B) {
	f, err := riofs.Open("../testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	var mc float32
	sc, err := NewScannerVars(tree, ReadVar{Name: "mcWeight", Value: &mc})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sc.SeekEntry(0)
		for sc.Next() {
			err := sc.Scan()
			if err != nil {
				b.Error(err)
			}
			sum += mc
		}
	}
}

func BenchmarkTreeScannerVarsBigFileSlice(b *testing.B) {
	f, err := riofs.Open("../testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	sc, err := NewTreeScannerVars(tree, ReadVar{Name: "lep_pt"})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sc.SeekEntry(0)
		for sc.Next() {
			var data []float32
			err := sc.Scan(&data)
			if err != nil {
				b.Error(err)
			}
			sum += data[0]
		}
	}
}

func BenchmarkScannerVarsBigFileSlice(b *testing.B) {
	f, err := riofs.Open("../testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	var pt []float32
	sc, err := NewScannerVars(tree, ReadVar{Name: "lep_pt", Value: &pt})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sc.SeekEntry(0)
		for sc.Next() {
			err := sc.Scan()
			if err != nil {
				b.Error(err)
			}
			sum += pt[0]
		}
	}
}

func TestTreeScannerSeekEntry(t *testing.T) {
	t.Parallel()

	fname := "../testdata/chain.1.root"
	f, err := riofs.Open(fname)
	if err != nil {
		t.Fatalf("could not open ROOT file %q: %v", fname, err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := obj.(Tree)

	type Data struct {
		Event struct {
			Beg       string      `groot:"Beg"`
			F64       float64     `groot:"F64"`
			ArrF64    [10]float64 `groot:"ArrayF64"`
			N         int32       `groot:"N"`
			SliF64    []float64   `groot:"SliceF64"`
			StdStr    string      `groot:"StdStr"`
			StlVecF64 []float64   `groot:"StlVecF64"`
			StlVecStr []string    `groot:"StlVecStr"`
			End       string      `groot:"End"`
		} `groot:"evt"`
	}

	sc, err := NewTreeScanner(tree, &Data{})
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	for _, entry := range []int64{0, 1, 2, 0, 1, 2, 9, 0, 9, 1} {
		err := sc.SeekEntry(entry)
		if err != nil {
			t.Fatalf("could not seek to entry %d: %v", entry, err)
		}
		if !sc.Next() {
			t.Fatalf("could not read entry %d", entry)
		}
		var d Data
		err = sc.Scan(&d)
		if err != nil {
			t.Fatal(err)
		}
		i := sc.Entry()
		if i != entry {
			t.Fatalf("did not seek to entry %d. got=%d, want=%d", entry, i, entry)
		}
		if d.Event.F64 != float64(i) {
			t.Fatalf("entry[%d]:\ngot= %#v\nwant=%#v\n", i, d.Event.F64, float64(i))
		}
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		t.Fatal(err)
	}
}

func TestNewScanVars(t *testing.T) {
	f, err := riofs.Open("../testdata/leaves.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(Tree)

	vars := NewScanVars(tree)
	want := []ReadVar{
		{Name: "B", Leaf: "B", Value: new(bool)},
		{Name: "Str", Leaf: "Str", Value: new(string)},
		{Name: "I8", Leaf: "I8", Value: new(int8)},
		{Name: "I16", Leaf: "I16", Value: new(int16)},
		{Name: "I32", Leaf: "I32", Value: new(int32)},
		{Name: "I64", Leaf: "I64", Value: new(int64)},
		{Name: "U8", Leaf: "U8", Value: new(uint8)},
		{Name: "U16", Leaf: "U16", Value: new(uint16)},
		{Name: "U32", Leaf: "U32", Value: new(uint32)},
		{Name: "U64", Leaf: "U64", Value: new(uint64)},
		{Name: "F32", Leaf: "F32", Value: new(float32)},
		{Name: "F64", Leaf: "F64", Value: new(float64)},
		{Name: "D16", Leaf: "D16", Value: new(root.Float16)},
		{Name: "D32", Leaf: "D32", Value: new(root.Double32)},
		// arrays
		{Name: "ArrBs", Leaf: "ArrBs", Value: new([10]bool)},
		{Name: "ArrI8", Leaf: "ArrI8", Value: new([10]int8)},
		{Name: "ArrI16", Leaf: "ArrI16", Value: new([10]int16)},
		{Name: "ArrI32", Leaf: "ArrI32", Value: new([10]int32)},
		{Name: "ArrI64", Leaf: "ArrI64", Value: new([10]int64)},
		{Name: "ArrU8", Leaf: "ArrU8", Value: new([10]uint8)},
		{Name: "ArrU16", Leaf: "ArrU16", Value: new([10]uint16)},
		{Name: "ArrU32", Leaf: "ArrU32", Value: new([10]uint32)},
		{Name: "ArrU64", Leaf: "ArrU64", Value: new([10]uint64)},
		{Name: "ArrF32", Leaf: "ArrF32", Value: new([10]float32)},
		{Name: "ArrF64", Leaf: "ArrF64", Value: new([10]float64)},
		{Name: "ArrD16", Leaf: "ArrD16", Value: new([10]root.Float16)},
		{Name: "ArrD32", Leaf: "ArrD32", Value: new([10]root.Double32)},
		// slices
		{Name: "N", Leaf: "N", Value: new(int32)},
		{Name: "SliBs", Leaf: "SliBs", Value: new([]bool)},
		{Name: "SliI8", Leaf: "SliI8", Value: new([]int8)},
		{Name: "SliI16", Leaf: "SliI16", Value: new([]int16)},
		{Name: "SliI32", Leaf: "SliI32", Value: new([]int32)},
		{Name: "SliI64", Leaf: "SliI64", Value: new([]int64)},
		{Name: "SliU8", Leaf: "SliU8", Value: new([]uint8)},
		{Name: "SliU16", Leaf: "SliU16", Value: new([]uint16)},
		{Name: "SliU32", Leaf: "SliU32", Value: new([]uint32)},
		{Name: "SliU64", Leaf: "SliU64", Value: new([]uint64)},
		{Name: "SliF32", Leaf: "SliF32", Value: new([]float32)},
		{Name: "SliF64", Leaf: "SliF64", Value: new([]float64)},
		{Name: "SliD16", Leaf: "SliD16", Value: new([]root.Float16)},
		{Name: "SliD32", Leaf: "SliD32", Value: new([]root.Double32)},
	}

	n := len(want)
	if len(vars) < n {
		n = len(vars)
	}

	for i := 0; i < n; i++ {
		got := vars[i]
		if got.Name != want[i].Name {
			t.Fatalf("invalid read-var name[%d]: got=%q, want=%q", i, got.Name, want[i].Name)
		}
		if got.Leaf != want[i].Leaf {
			t.Fatalf("invalid read-var (name=%q) leaf-name[%d]: got=%q, want=%q", got.Name, i, got.Leaf, want[i].Leaf)
		}
		if got, want := reflect.TypeOf(got.Value), reflect.TypeOf(want[i].Value); got != want {
			t.Fatalf("invalid read-var (name=%q) type[%d]: got=%v, want=%v", vars[i].Name, i, got, want)
		}
	}

	if len(want) != len(vars) {
		t.Fatalf("invalid lengths. got=%d, want=%d", len(vars), len(want))
	}
}

func TestG4LikeTree(t *testing.T) {
	t.Parallel()
	fname := rtests.XrdRemote("testdata/g4-like.root")

	f, err := riofs.Open(fname)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, err := f.Get("mytree")
	if err != nil {
		t.Fatal(err)
	}

	tree := obj.(Tree)

	type EventData struct {
		I32 int32     `groot:"i32"`
		F64 float64   `groot:"f64"`
		Sli []float64 `groot:"slif64"`
	}

	want := func(i int64) (data EventData) {
		data.I32 = int32(i + 1)
		data.F64 = float64(i + 1)
		data.Sli = make([]float64, i)
		for ii := range data.Sli {
			data.Sli[ii] = float64(ii) + float64(i)
		}
		return data
	}

	data := EventData{
		Sli: make([]float64, 0),
	}
	rvars := []ReadVar{
		{Name: "i32", Value: &data.I32},
		{Name: "f64", Value: &data.F64},
		{Name: "slif64", Value: &data.Sli},
	}
	sc, err := NewScannerVars(tree, rvars...)
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()
	for sc.Next() {
		err := sc.Scan()
		if err != nil {
			t.Fatal(err)
		}
		i := sc.Entry()
		if !reflect.DeepEqual(data, want(i)) {
			t.Fatalf("entry[%d]:\ngot= %#v.\nwant=%#v\n", i, data, want(i))
		}
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		t.Fatal(err)
	}
}

func TestMultiLeafBranchWithReadVars(t *testing.T) {
	t.Parallel()

	f, err := riofs.Open("../testdata/root_numpy_struct.root")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer f.Close()

	obj, err := f.Get("test")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	tree := obj.(Tree)

	type Data struct {
		b1l1 int32
		b1l2 float32
		b2l1 int32
		b2l2 float32
	}

	var (
		data Data
		want = []Data{
			{10, 15.5, 20, 781.2},
		}
	)

	rvars := []ReadVar{
		{
			Name:  "branch1",
			Leaf:  "intleaf",
			Value: &data.b1l1,
		},
		{
			Name:  "branch1",
			Leaf:  "floatleaf",
			Value: &data.b1l2,
		},
		{
			Name:  "branch2",
			Leaf:  "intleaf",
			Value: &data.b2l1,
		},
		{
			Name:  "branch2",
			Leaf:  "floatleaf",
			Value: &data.b2l2,
		},
	}

	sc, err := NewScannerVars(tree, rvars...)
	if err != nil {
		t.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan()
		if err != nil {
			t.Fatalf("could not scan entry %d: %+v", sc.Entry(), err)
		}

		if got, want := data, want[sc.Entry()]; !reflect.DeepEqual(got, want) {
			t.Fatalf("invalid entry %d:\ngot= %#v\nwant=%#v", sc.Entry(), got, want)
		}
	}
}

func TestMultiLeafBranchWithTreeReadVars(t *testing.T) {
	t.Parallel()

	f, err := riofs.Open("../testdata/root_numpy_struct.root")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer f.Close()

	obj, err := f.Get("test")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	tree := obj.(Tree)

	type B struct {
		L1 int32   `groot:"intleaf"`
		L2 float32 `groot:"floatleaf"`
	}

	type Data struct {
		B1 B `groot:"branch1"`
		B2 B `groot:"branch2"`
	}

	var (
		data Data
		want = []Data{
			{B{10, 15.5}, B{20, 781.2}},
		}
	)

	rvars := []ReadVar{
		{
			Name:  "branch1",
			Value: &data.B1,
		},
		{
			Name:  "branch2",
			Value: &data.B2,
		},
	}

	sc, err := NewTreeScannerVars(tree, rvars...)
	if err != nil {
		t.Fatalf("could not create scanner: %+v", err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan(&data.B1, &data.B2)
		if err != nil {
			t.Fatalf("could not scan entry %d: %+v", sc.Entry(), err)
		}

		if got, want := data, want[sc.Entry()]; !reflect.DeepEqual(got, want) {
			t.Fatalf("invalid entry %d:\ngot= %#v\nwant=%#v", sc.Entry(), got, want)
		}
	}
}
