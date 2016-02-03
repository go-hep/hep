// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/go-hep/csvutil/csvdriver"
	"github.com/go-hep/hbook"
)

var (
	nt *hbook.NTuple
)

func TestNTupleScanH1D(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := nt.ScanH1D("x", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 10,
		len:     10,
		mean:    4.5,
		rms:     2.8722813232690143,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != 1 {
			t.Errorf("error bin(%d)=%v. want=1\n", i, v)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestNTupleScanH1DWhere(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := nt.ScanH1D("x where (id > 4 && id < 10)", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}

	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 5,
		len:     10,
		mean:    7,
		rms:     1.4142135623730951,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		want := float64(0)
		if i > 4 {
			want = 1
		}
		if v != want {
			t.Errorf("error bin(%d)=%v. want=%v\n", i, v, want)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestNTupleScanH1DInt(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := nt.ScanH1D("id", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 10,
		len:     10,
		mean:    4.5,
		rms:     2.8722813232690143,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != 1 {
			t.Errorf("error bin(%d)=%v. want=1\n", i, v)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestNTupleScan(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	err := nt.Scan("id, x", func(id int64, x float64) error {
		h.Fill(x, 1)
		return nil
	})
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 10,
		len:     10,
		mean:    4.5,
		rms:     2.8722813232690143,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != 1 {
			t.Errorf("error bin(%d)=%v. want=1\n", i, v)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestNTupleScanH1DFromCSVWithCommas(t *testing.T) {
	db, err := sql.Open("csv", "testdata/simple-comma.csv")
	if err != nil {
		t.Fatalf("error opening CSV db: %v\n", err)
	}
	defer db.Close()

	nt, err := hbook.OpenNTuple(db, "csv")
	if err != nil {
		t.Fatalf("error opening ntuple: %v\n", err)
	}

	h := hbook.NewH1D(10, 0, 10)
	h, err = nt.ScanH1D("var2", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 10,
		len:     10,
		mean:    4.5,
		rms:     2.8722813232690143,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != 1 {
			t.Errorf("error bin(%d)=%v. want=1\n", i, v)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestNTupleScanH1DFromCSV(t *testing.T) {
	db, err := csvdriver.Conn{
		File:    "testdata/simple.csv",
		Comma:   ';',
		Comment: '#',
	}.Open()
	if err != nil {
		t.Fatalf("error opening CSV db: %v\n", err)
	}
	defer db.Close()

	nt, err := hbook.OpenNTuple(db, "csv")
	if err != nil {
		t.Fatalf("error opening ntuple: %v\n", err)
	}

	h := hbook.NewH1D(10, 0, 10)
	h, err = nt.ScanH1D("var2", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	want := struct {
		entries int64
		len     int
		mean    float64
		rms     float64
	}{
		entries: 10,
		len:     10,
		mean:    4.5,
		rms:     2.8722813232690143,
	}

	if h.Entries() != want.entries {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.entries)
	}
	if h.Len() != want.len {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.len)
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != 1 {
			t.Errorf("error bin(%d)=%v. want=1\n", i, v)
		}
	}

	if mean := h.Mean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.RMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestCreateNTuple(t *testing.T) {
	db, err := sql.Open("ql", "memory://ntuple.db")
	if err != nil {
		t.Fatalf("error creating db: %v\n", err)
	}
	defer db.Close()

	const ntname = "ntup"
	nt, err := hbook.CreateNTuple(db, ntname, int64(0), float64(0))
	if err != nil {
		t.Fatalf("error creating ntuple: %v\n", err)
	}

	if nt.Name() != ntname {
		t.Errorf("invalid ntuple name. got=%q want=%q\n", nt.Name(), ntname)
	}

	descr := []struct {
		n string
		t reflect.Type
	}{
		{
			n: "var1",
			t: reflect.TypeOf(int64(0)),
		},
		{
			n: "var2",
			t: reflect.TypeOf(float64(0)),
		},
	}
	if len(nt.Cols()) != len(descr) {
		t.Fatalf("invalid cols. got=%d. want=%d\n", len(nt.Cols()), len(descr))
	}

	for i := 0; i < len(descr); i++ {
		col := nt.Cols()[i]
		exp := descr[i]
		if col.Name() != exp.n {
			t.Errorf("col[%d]: invalid name. got=%q. want=%q\n",
				i, col.Name(), exp.n,
			)
		}
		if col.Type() != exp.t {
			t.Errorf("col[%d]: invalid type. got=%v. want=%v\n",
				i, col.Type(), exp.t,
			)
		}
	}
}

func TestCreateNTupleFromStruct(t *testing.T) {
	db, err := sql.Open("ql", "memory://ntuple-struct.db")
	if err != nil {
		t.Fatalf("error creating db: %v\n", err)
	}
	defer db.Close()

	type dataType struct {
		I  int64
		F  float64
		FF float64 `rio:"ff"`
		S  string  `rio:"STR" hbook:"str"`
	}

	const ntname = "ntup"
	nt, err := hbook.CreateNTuple(db, ntname, dataType{})
	if err != nil {
		t.Fatalf("error creating ntuple: %v\n", err)
	}

	if nt.Name() != ntname {
		t.Errorf("invalid ntuple name. got=%q want=%q\n", nt.Name(), ntname)
	}

	descr := []struct {
		n string
		t reflect.Type
	}{
		{
			n: "I",
			t: reflect.TypeOf(int64(0)),
		},
		{
			n: "F",
			t: reflect.TypeOf(float64(0)),
		},
		{
			n: "ff",
			t: reflect.TypeOf(float64(0)),
		},
		{
			n: "str",
			t: reflect.TypeOf(""),
		},
	}
	if len(nt.Cols()) != len(descr) {
		t.Fatalf("invalid cols. got=%d. want=%d\n", len(nt.Cols()), len(descr))
	}

	for i := 0; i < len(descr); i++ {
		col := nt.Cols()[i]
		exp := descr[i]
		if col.Name() != exp.n {
			t.Errorf("col[%d]: invalid name. got=%q. want=%q\n",
				i, col.Name(), exp.n,
			)
		}
		if col.Type() != exp.t {
			t.Errorf("col[%d]: invalid type. got=%v. want=%v\n",
				i, col.Type(), exp.t,
			)
		}
	}
}

func init() {
	var err error
	db, err := sql.Open("ql", "memory://mem.db")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	_, err = tx.Exec("create table data (id int, x float64);")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		x := float64(i)
		_, err = tx.Exec("insert into data values($1, $2);", i, x)
		if err != nil {
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	nt, err = hbook.OpenNTuple(db, "data")
	if err != nil {
		panic(err)
	}
}
