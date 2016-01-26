// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplt_test

import (
	"database/sql"
	"testing"

	_ "github.com/cznic/ql/driver"
	"github.com/go-hep/hbook"
	"github.com/go-hep/hbook/hplt"
)

var (
	db *sql.DB
)

func TestScanH1D(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := hplt.ScanH1D(db, "select x from data", h)
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

func TestScanH1DWhere(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := hplt.ScanH1D(db, "select x from data where id > 4", h)
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

func TestScanH1DInt(t *testing.T) {
	h := hbook.NewH1D(10, 0, 10)
	h, err := hplt.ScanH1D(db, "select id from data", h)
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
func init() {
	var err error
	db, err = sql.Open("ql", "memory://mem.db")
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

}
