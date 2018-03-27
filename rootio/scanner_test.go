// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

type ScannerData struct {
	I32    int32       `rootio:"Int32"`
	I64    int64       `rootio:"Int64"`
	U32    uint32      `rootio:"UInt32"`
	U64    uint64      `rootio:"UInt64"`
	F32    float32     `rootio:"Float32"`
	F64    float64     `rootio:"Float64"`
	Str    string      `rootio:"Str"`
	ArrI32 [10]int32   `rootio:"ArrayInt32"`
	ArrI64 [10]int64   `rootio:"ArrayInt64"`
	ArrU32 [10]uint32  `rootio:"ArrayUInt32"`
	ArrU64 [10]uint64  `rootio:"ArrayUInt64"`
	ArrF32 [10]float32 `rootio:"ArrayFloat32"`
	ArrF64 [10]float64 `rootio:"ArrayFloat64"`
	N      int32       `rootio:"N"`
	SliI32 []int32     `rootio:"SliceInt32"`
	SliI64 []int64     `rootio:"SliceInt64"`
	SliU32 []uint32    `rootio:"SliceUInt32"`
	SliU64 []uint64    `rootio:"SliceUInt64"`
	SliF32 []float32   `rootio:"SliceFloat32"`
	SliF64 []float64   `rootio:"SliceFloat64"`
}

func TestTreeScannerStruct(t *testing.T) {
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

	want := func(i int64) (data ScannerData) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		data.Str = fmt.Sprintf("evt-%03d", i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
			data.ArrI64[ii] = int64(i)
			data.ArrU32[ii] = uint32(i)
			data.ArrU64[ii] = uint64(i)
			data.ArrF32[ii] = float32(i)
			data.ArrF64[ii] = float64(i)
		}
		data.N = int32(i) % 10
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		for ii := 0; ii < int(data.N); ii++ {
			data.SliI32[ii] = int32(i)
			data.SliI64[ii] = int64(i)
			data.SliU32[ii] = uint32(i)
			data.SliU64[ii] = uint64(i)
			data.SliF32[ii] = float32(i)
			data.SliF64[ii] = float64(i)
		}
		return data
	}

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
}

func TestScannerStruct(t *testing.T) {
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

	want := func(i int64) (data ScannerData) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		data.Str = fmt.Sprintf("evt-%03d", i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
			data.ArrI64[ii] = int64(i)
			data.ArrU32[ii] = uint32(i)
			data.ArrU64[ii] = uint64(i)
			data.ArrF32[ii] = float32(i)
			data.ArrF64[ii] = float64(i)
		}
		data.N = int32(i) % 10
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		for ii := 0; ii < int(data.N); ii++ {
			data.SliI32[ii] = int32(i)
			data.SliI64[ii] = int64(i)
			data.SliU32[ii] = uint32(i)
			data.SliU64[ii] = uint64(i)
			data.SliF32[ii] = float32(i)
			data.SliF64[ii] = float64(i)
		}
		return data
	}

	var data ScannerData
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
}

func TestScannerVars(t *testing.T) {
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

	want := func(i int64) (data ScannerData) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		data.Str = fmt.Sprintf("evt-%03d", i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
			data.ArrI64[ii] = int64(i)
			data.ArrU32[ii] = uint32(i)
			data.ArrU64[ii] = uint64(i)
			data.ArrF32[ii] = float32(i)
			data.ArrF64[ii] = float64(i)
		}
		data.N = int32(i) % 10
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		for ii := 0; ii < int(data.N); ii++ {
			data.SliI32[ii] = int32(i)
			data.SliI64[ii] = int64(i)
			data.SliU32[ii] = uint32(i)
			data.SliU64[ii] = uint64(i)
			data.SliF32[ii] = float32(i)
			data.SliF64[ii] = float64(i)
		}
		return data
	}

	var data ScannerData
	scanVars := []ScanVar{
		{Name: "Int32", Value: &data.I32},
		{Name: "Int64", Value: &data.I64},
		{Name: "UInt32", Value: &data.U32},
		{Name: "UInt64", Value: &data.U64},
		{Name: "Float32", Value: &data.F32},
		{Name: "Float64", Value: &data.F64},
		{Name: "Str", Value: &data.Str},
		{Name: "ArrayInt32", Value: &data.ArrI32},
		{Name: "ArrayInt64", Value: &data.ArrI64},
		{Name: "ArrayUInt32", Value: &data.ArrU32},
		{Name: "ArrayUInt64", Value: &data.ArrU64},
		{Name: "ArrayFloat32", Value: &data.ArrF32},
		{Name: "ArrayFloat64", Value: &data.ArrF64},
		{Name: "N", Value: &data.N},
		{Name: "SliceInt32", Value: &data.SliI32},
		{Name: "SliceInt64", Value: &data.SliI64},
		{Name: "SliceUInt32", Value: &data.SliU32},
		{Name: "SliceUInt64", Value: &data.SliU64},
		{Name: "SliceFloat32", Value: &data.SliF32},
		{Name: "SliceFloat64", Value: &data.SliF64},
	}
	sc, err := NewScannerVars(tree, scanVars...)
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

func TestTreeScannerVarsMultipleTimes(t *testing.T) {
	f, err := Open("testdata/mc_105986.ZZ.root")
	if err != nil {
		t.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		t.Fatal(err)
	}
	tree := obj.(Tree)

	for i := 0; i < 10; i++ {
		sc, err := NewTreeScannerVars(tree, ScanVar{Name: "lep_pt"})
		if err != nil {
			t.Fatal(err)
		}
		defer sc.Close()

		for sc.Next() {
			var data []float32
			err := sc.Scan(&data)
			if err != nil {
				t.Error(err)
			}
		}
		err = sc.Err()
		if err != nil {
			t.Error(err)
		}
	}
}

func TestTreeScannerVars(t *testing.T) {
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

	want := func(i int64) (data ScannerData) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		data.Str = fmt.Sprintf("evt-%03d", i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
			data.ArrI64[ii] = int64(i)
			data.ArrU32[ii] = uint32(i)
			data.ArrU64[ii] = uint64(i)
			data.ArrF32[ii] = float32(i)
			data.ArrF64[ii] = float64(i)
		}
		data.N = int32(i) % 10
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		for ii := 0; ii < int(data.N); ii++ {
			data.SliI32[ii] = int32(i)
			data.SliI64[ii] = int64(i)
			data.SliU32[ii] = uint32(i)
			data.SliU64[ii] = uint64(i)
			data.SliF32[ii] = float32(i)
			data.SliF64[ii] = float64(i)
		}
		return data
	}

	scanVars := []ScanVar{
		{Name: "Int32"},
		{Name: "Int64"},
		{Name: "UInt32"},
		{Name: "UInt64"},
		{Name: "Float32"},
		{Name: "Float64"},
		{Name: "Str"},
		{Name: "ArrayInt32"},
		{Name: "ArrayInt64"},
		{Name: "ArrayUInt32"},
		{Name: "ArrayUInt64"},
		{Name: "ArrayFloat32"},
		{Name: "ArrayFloat64"},
		{Name: "N"},
		{Name: "SliceInt32"},
		{Name: "SliceInt64"},
		{Name: "SliceUInt32"},
		{Name: "SliceUInt64"},
		{Name: "SliceFloat32"},
		{Name: "SliceFloat64"},
	}
	sc, err := NewTreeScannerVars(tree, scanVars...)
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()
	var d1 ScannerData
	for sc.Next() {
		err := sc.Scan(
			&d1.I32, &d1.I64, &d1.U32, &d1.U64, &d1.F32, &d1.F64,
			&d1.Str,
			&d1.ArrI32, &d1.ArrI64, &d1.ArrU32, &d1.ArrU64, &d1.ArrF32, &d1.ArrF64,
			&d1.N,
			&d1.SliI32, &d1.SliI64, &d1.SliU32, &d1.SliU64, &d1.SliF32, &d1.SliF64,
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
			&d2.I32, &d2.I64, &d2.U32, &d2.U64, &d2.F32, &d2.F64,
			&d2.Str,
			&d2.ArrI32, &d2.ArrI64, &d2.ArrU32, &d2.ArrU64, &d2.ArrF32, &d2.ArrF64,
			&d2.N,
			&d2.SliI32, &d2.SliI64, &d2.SliU32, &d2.SliU64, &d2.SliF32, &d2.SliF64,
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
}

func TestScannerVarsMultipleTimes(t *testing.T) {
	f, err := Open("testdata/mc_105986.ZZ.root")
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
		sc, err := NewScannerVars(tree, ScanVar{Name: "lep_pt", Value: &pt})
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
}

func TestTreeScannerStructWithCounterLeaf(t *testing.T) {
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

	type Data struct {
		Sli []int32 `rootio:"SliceInt32"`
	}
	var data Data

	want := func(i int64) Data {
		var data Data
		n := int32(i) % 10
		data.Sli = make([]int32, int(n))
		for ii := 0; ii < int(n); ii++ {
			data.Sli[ii] = int32(i)
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
}

func TestScannerStructWithCounterLeaf(t *testing.T) {
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

	type Data struct {
		Sli []int32 `rootio:"SliceInt32"`
	}
	var data Data

	want := func(i int64) Data {
		var data Data
		n := int32(i) % 10
		data.Sli = make([]int32, int(n))
		for ii := 0; ii < int(n); ii++ {
			data.Sli[ii] = int32(i)
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
}

func TestTreeScannerVarsWithCounterLeaf(t *testing.T) {
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

	want := func(i int64) []int32 {
		n := int32(i) % 10
		data := make([]int32, int(n))
		for ii := 0; ii < int(n); ii++ {
			data[ii] = int32(i)
		}
		return data
	}

	scanVar := ScanVar{Name: "SliceInt32"}
	sc, err := NewTreeScannerVars(tree, scanVar)
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
}

func TestScannerVarsWithCounterLeaf(t *testing.T) {
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

	want := func(i int64) []int32 {
		n := int32(i) % 10
		data := make([]int32, int(n))
		for ii := 0; ii < int(n); ii++ {
			data[ii] = int32(i)
		}
		return data
	}

	var data []int32
	scanVar := ScanVar{Name: "SliceInt32", Value: &data}
	sc, err := NewScannerVars(tree, scanVar)
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

func TestScannerStructWithStdVectorBool(t *testing.T) {
	files, err := filepath.Glob("testdata/stdvec-bool-*.root")
	if err != nil {
		t.Fatal(err)
	}

	for _, fname := range files {
		t.Run(fname, func(t *testing.T) {
			f, err := Open(fname)
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
				Bool    bool     `rootio:"Bool"`
				ArrBool [10]bool `rootio:"ArrayBool"`
				N       int32    `rootio:"N"`
				SliBool []bool   `rootio:"SliceBool"`
				StlBool []bool   `rootio:"StlVecBool"`
			}
			type Event struct {
				Data Data `rootio:"evt"`
			}

			want := func(i int64) Event {
				var data Data
				data.Bool = i%2 == 0
				for ii := range data.ArrBool {
					data.ArrBool[ii] = i%2 == 0
				}
				data.N = int32(i) % 10
				data.SliBool = make([]bool, int(data.N))
				data.StlBool = make([]bool, int(data.N))
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
	f, err := Open("testdata/small-flat-tree.root")
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
		F64 float64 `rootio:"Float64"`
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
		s.SeekEntry(0)
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
	f, err := Open("testdata/small-flat-tree.root")
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
		F64 float64 `rootio:"Float64"`
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
		s.SeekEntry(0)
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
	f, err := Open("testdata/small-flat-tree.root")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		b.Fatal(err)
	}

	tree := obj.(Tree)

	scanVars := []ScanVar{
		{Name: "Float64"},
	}
	s, err := NewTreeScannerVars(tree, scanVars...)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var data ScannerData
	var sum float64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.SeekEntry(0)
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
	f, err := Open("testdata/small-flat-tree.root")
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
	scanVars := []ScanVar{
		{Name: "Float64", Value: &f64},
	}
	s, err := NewScannerVars(tree, scanVars...)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	var sum float64

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.SeekEntry(0)
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
	f, err := Open("testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	sc, err := NewTreeScannerVars(tree, ScanVar{Name: "mcWeight"})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sc.SeekEntry(0)
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
	f, err := Open("testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	var mc float32
	sc, err := NewScannerVars(tree, ScanVar{Name: "mcWeight", Value: &mc})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sc.SeekEntry(0)
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
	f, err := Open("testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	sc, err := NewTreeScannerVars(tree, ScanVar{Name: "lep_pt"})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sc.SeekEntry(0)
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
	f, err := Open("testdata/mc_105986.ZZ.root")
	if err != nil {
		b.Skip(err)
	}

	obj, err := f.Get("mini")
	if err != nil {
		b.Fatal(err)
	}
	tree := obj.(Tree)

	var pt []float32
	sc, err := NewScannerVars(tree, ScanVar{Name: "lep_pt", Value: &pt})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	var sum float32

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sc.SeekEntry(0)
		for sc.Next() {
			err := sc.Scan()
			if err != nil {
				b.Error(err)
			}
			sum += pt[0]
		}
	}
}
