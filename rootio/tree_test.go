// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestFlatTree(t *testing.T) {
	f, err := Open("testdata/small-flat-tree.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, ok := f.Get("tree")
	if !ok {
		t.Fatalf("could not retrieve tree [tree]")
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

	obj, ok := f.Get("tree")
	if !ok {
		t.Errorf("%s: could not retrieve tree [tree]", name)
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
	f, err := Open("testdata/simple.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	obj, ok := f.Get("tree")
	if !ok {
		t.Fatalf("could not retrieve tree [tree]")
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
