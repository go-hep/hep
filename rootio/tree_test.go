// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io"
	"testing"
)

func TestFlatTree(t *testing.T) {
	f, err := Open("testdata/small.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	myprintf(">>> f.Get(tree)...\n")
	obj, ok := f.Get("tree")
	if !ok {
		t.Fatalf("could not retrieve tree [tree]")
	}

	key := obj.(*Key)
	if got, want := key.Name(), "tree"; got != want {
		t.Fatalf("key.Name: got=%q. want=%q", got, want)
	}

	tree := key.Value().(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}
	myprintf(">>> f.Get(tree)... [done]\n")

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

	if got, want := tree.TotBytes(), int64(40506); got != want {
		t.Fatalf("tree.totbytes: got=%v. want=%v", got, want)
	}

	if got, want := tree.ZipBytes(), int64(4184); got != want {
		t.Fatalf("tree.zipbytes: got=%v. want=%v", got, want)
	}

	type dataType struct {
		I32    int32       `rootio:"Int32"`
		I64    int64       `rootio:"Int64"`
		U32    uint32      `rootio:"UInt32"`
		U64    uint64      `rootio:"UInt64"`
		F32    float32     `rootio:"Float32"`
		F64    float64     `rootio:"Float64"`
		ArrI32 [10]int32   `rootio:"ArrayInt32"`
		ArrI64 [10]int64   `rootio:"ArrayInt64"`
		ArrU32 [10]uint32  `rootio:"ArrayUInt32"`
		ArrU64 [10]uint64  `rootio:"ArrayUInt64"`
		ArrF32 [10]float32 `rootio:"ArrayFloat32"`
		ArrF64 [10]float64 `rootio:"ArrayFloat64"`
	}

	want := func(i int64) (data dataType) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
		}
		for ii := range data.ArrI64 {
			data.ArrI64[ii] = int64(i)
		}
		for ii := range data.ArrU32 {
			data.ArrU32[ii] = uint32(i)
		}
		for ii := range data.ArrU64 {
			data.ArrU64[ii] = uint64(i)
		}
		for ii := range data.ArrF32 {
			data.ArrF32[ii] = float32(i)
		}
		for ii := range data.ArrF64 {
			data.ArrF64[ii] = float64(i)
		}
		return data
	}

	sc, err := NewScanner(tree, &dataType{})
	if err != nil {
		t.Fatal(err)
	}
	for sc.Next() {
		var data dataType
		err := sc.Scan(&data)
		if err != nil {
			t.Fatal(err)
		}
		i := sc.Entry()
		if data != want(i) {
			t.Fatalf("entry[%d]: got=%v. want=%v\n", i, data, want(i))
		}
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		t.Fatal(err)
	}
}

func TestSimpleTree(t *testing.T) {
	f, err := Open("testdata/simple.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	myprintf(">>> f.Get(tree)...\n")
	obj, ok := f.Get("tree")
	if !ok {
		t.Fatalf("could not retrieve tree [tree]")
	}

	key := obj.(*Key)
	if got, want := key.Name(), "tree"; got != want {
		t.Fatalf("key.Name: got=%q. want=%q", got, want)
	}

	tree := key.Value().(Tree)
	if got, want := tree.Name(), "tree"; got != want {
		t.Fatalf("tree.Name: got=%q. want=%q", got, want)
	}
	myprintf(">>> f.Get(tree)... [done]\n")

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
