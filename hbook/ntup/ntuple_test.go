// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntup

import (
	"database/sql"
	"reflect"
	"testing"

	"go-hep.org/x/hep/csvutil/csvdriver"
	"go-hep.org/x/hep/hbook"
)

var (
	nt *Ntuple
)

func TestScanH1D(t *testing.T) {
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
		rms:     5.338539126015656,
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

	if mean := h.XMean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.XRMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestScanH1DWithoutH1(t *testing.T) {
	want := hbook.NewH1D(100, 0, nextULP(9))
	for i := 0; i < 10; i++ {
		want.Fill(float64(i), 1)
	}

	h, err := nt.ScanH1D("x", nil)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}
	if h.Entries() != want.Entries() {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.Entries())
	}
	if h.Len() != want.Len() {
		t.Errorf("error. got %v bins. want=%d\n", h.Len(), want.Len())
	}

	for i := 0; i < h.Len(); i++ {
		v := h.Value(i)
		if v != want.Value(i) {
			t.Errorf("error bin(%d)=%v. want=%v\n", i, v, want.Value(i))
		}
	}

	if mean := h.XMean(); mean != want.XMean() {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.XMean())
	}
	if rms := h.XRMS(); rms != want.XRMS() {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.XRMS())
	}
}

func TestScanH1DWhere(t *testing.T) {
	for _, where := range []string{
		"x where (id > 4 && id < 10)",
		"x WHERE (id > 4 && id < 10)",
		"x where (id > 4 && id < 10) order by id();",
		"x WHERE (id > 4 && id < 10) ORDER by id();",
		"x WHERE (id > 4 && id < 10) order by id();",
		"x where (id > 4 && id < 10) ORDER by id();",
	} {
		t.Run("", func(t *testing.T) {
			h := hbook.NewH1D(10, 0, 10)
			h, err := nt.ScanH1D(where, h)
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
				rms:     7.14142842854285,
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

			if mean := h.XMean(); mean != want.mean {
				t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
			}
			if rms := h.XRMS(); rms != want.rms {
				t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
			}
		})
	}
}

func TestScanH1DInt(t *testing.T) {
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
		rms:     5.338539126015656,
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

	if mean := h.XMean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.XRMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestScanH2D(t *testing.T) {
	want := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	for i := 0; i < 10; i++ {
		v := float64(i)
		want.Fill(v, v, 1)
	}

	h := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	h, err := nt.ScanH2D("id, x", h)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}

	if h.Entries() != want.Entries() {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.Entries())
	}

	type gridXYZer interface {
		Dims() (c, r int)
		Z(c, r int) float64
		X(c int) float64
		Y(r int) float64
	}

	cmpGrid := func(a, b gridXYZer) {
		ac, ar := a.Dims()
		bc, br := b.Dims()
		if ac != bc {
			t.Fatalf("got=%d want=%d", ac, bc)
		}
		if ar != br {
			t.Fatalf("got=%d want=%d", ar, br)
		}
		for i := 0; i < ar; i++ {
			ay := a.Y(i)
			by := b.Y(i)
			if ay != by {
				t.Fatalf("got=%v. want=%v", ay, by)
			}
			for j := 0; j < ac; j++ {
				if i == 0 {
					ax := a.X(j)
					bx := b.X(j)
					if ax != bx {
						t.Fatalf("got=%v. want=%v", ax, bx)
					}
				}
				az := a.Z(j, i)
				bz := b.Z(j, i)
				if az != bz {
					t.Fatalf("got=%v. want=%v", az, bz)
				}
			}
		}
	}

	cmpGrid(h.GridXYZ(), want.GridXYZ())

	if mean := h.XMean(); mean != want.XMean() {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.XMean())
	}
	if rms := h.XRMS(); rms != want.XRMS() {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.XRMS())
	}
}

func TestScanH2DWithoutH2D(t *testing.T) {
	want := hbook.NewH2D(100, 0, nextULP(9), 100, 0, nextULP(9))
	for i := 0; i < 10; i++ {
		v := float64(i)
		want.Fill(v, v, 1)
	}

	h, err := nt.ScanH2D("id, x", nil)
	if err != nil {
		t.Errorf("error running query: %v\n", err)
	}

	if h.Entries() != want.Entries() {
		t.Errorf("error. got %v entries. want=%v\n", h.Entries(), want.Entries())
	}

	type gridXYZer interface {
		Dims() (c, r int)
		Z(c, r int) float64
		X(c int) float64
		Y(r int) float64
	}

	cmpGrid := func(a, b gridXYZer) {
		ac, ar := a.Dims()
		bc, br := b.Dims()
		if ac != bc {
			t.Fatalf("got=%d want=%d", ac, bc)
		}
		if ar != br {
			t.Fatalf("got=%d want=%d", ar, br)
		}
		for i := 0; i < ar; i++ {
			ay := a.Y(i)
			by := b.Y(i)
			if ay != by {
				t.Fatalf("got=%v. want=%v", ay, by)
			}
			for j := 0; j < ac; j++ {
				if i == 0 {
					ax := a.X(j)
					bx := b.X(j)
					if ax != bx {
						t.Fatalf("got=%v. want=%v", ax, bx)
					}
				}
				az := a.Z(j, i)
				bz := b.Z(j, i)
				if az != bz {
					t.Fatalf("got=%v. want=%v", az, bz)
				}
			}
		}
	}

	cmpGrid(h.GridXYZ(), want.GridXYZ())

	if mean := h.XMean(); mean != want.XMean() {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.XMean())
	}
	if rms := h.XRMS(); rms != want.XRMS() {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.XRMS())
	}
}

func TestScan(t *testing.T) {
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
		rms:     5.338539126015656,
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

	if mean := h.XMean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.XRMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestScanH1DFromCSVWithCommas(t *testing.T) {
	db, err := sql.Open("csv", "testdata/simple-comma.csv")
	if err != nil {
		t.Fatalf("error opening CSV db: %v\n", err)
	}
	defer db.Close()

	nt, err := Open(db, "csv")
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
		rms:     5.338539126015656,
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

	if mean := h.XMean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.XRMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestScanH1DFromCSV(t *testing.T) {
	db, err := csvdriver.Conn{
		File:    "testdata/simple.csv",
		Comma:   ';',
		Comment: '#',
	}.Open()
	if err != nil {
		t.Fatalf("error opening CSV db: %v\n", err)
	}
	defer db.Close()

	nt, err := Open(db, "csv")
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
		rms:     5.338539126015656,
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

	if mean := h.XMean(); mean != want.mean {
		t.Errorf("error: mean=%v. want=%v\n", mean, want.mean)
	}
	if rms := h.XRMS(); rms != want.rms {
		t.Errorf("error: rms=%v. want=%v\n", rms, want.rms)
	}
}

func TestScanInvalid(t *testing.T) {
	for _, tc := range []struct {
		name string
		fct  interface{}
	}{
		{
			name: "nil func",
			fct:  nil,
		},
		{
			name: "not a func",
			fct:  0,
		},
		{
			name: "0-arity",
			fct:  func() {},
		},
		{
			name: "invalid func",
			fct:  func() int { return 0 },
		},
		{
			name: "2-arity",
			fct:  func() (error, int) { return nil, 1 },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := nt.Scan("id, x", tc.fct)
			if err == nil {
				t.Fatalf("expected an error")
			}
		})
	}
}

func TestCreate(t *testing.T) {
	db, err := sql.Open("ql", "memory://ntuple.db")
	if err != nil {
		t.Fatalf("error creating db: %v\n", err)
	}
	defer db.Close()

	const ntname = "ntup"
	nt, err := Create(db, ntname, int64(0), float64(0))
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

func TestCreateFromStruct(t *testing.T) {
	db, err := sql.Open("ql", "memory://ntuple-struct.db")
	if err != nil {
		t.Fatalf("error creating db: %v\n", err)
	}
	defer db.Close()

	type dataType struct {
		I   int64
		F   float64
		FF  float64 `rio:"ff" hbook:"-"`
		S   string  `rio:"STR" hbook:"str"`
		not string
	}

	const ntname = "ntup"
	nt, err := Create(db, ntname, dataType{})
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

func TestCreateInvalid(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		cols []interface{}
	}{
		{
			name: "missing-col-def.db",
			err:  ErrMissingColDef,
		},
		{
			name: "one-value.db",
			cols: []interface{}{int64(0)},
		},
		{
			name: "err-chan.db",
			cols: []interface{}{make(chan int)},
			err:  errChanType,
		},
		{
			name: "err-struct-chan.db",
			cols: []interface{}{func() interface{} {
				type Person struct {
					Field chan int
				}
				return Person{Field: make(chan int)}
			}(),
			},
			err: errChanType,
		},
		//		{
		//			name: "err-iface.db",
		//			cols: []interface{}{(io.Writer)(os.Stdout)},
		//			err:  errIfaceType,
		//		},
		//		{
		//			name: "err-eface.db",
		//			cols: []interface{}{interface{}(nil)},
		//			err:  errIfaceType,
		//		},
		{
			name: "err-map.db",
			cols: []interface{}{make(map[string]int)},
			err:  errMapType,
		},
		{
			name: "err-struct-map.db",
			cols: []interface{}{func() interface{} {
				type Person struct {
					Field map[string]int
				}
				return Person{Field: make(map[string]int)}
			}(),
			},
			err: errMapType,
		},
		{
			name: "err-slice.db",
			cols: []interface{}{make([]int, 2)},
			err:  errSliceType,
		},
		{
			name: "err-struct-slice.db",
			cols: []interface{}{func() interface{} {
				type Person struct {
					Field []int
				}
				return Person{Field: make([]int, 2)}
			}(),
			},
			err: errSliceType,
		},
		{
			name: "err-struct.db",
			cols: []interface{}{func() interface{} {
				type Name struct {
					Name string
				}
				type Person struct {
					Name Name
				}
				var p Person
				p.Name.Name = "bob"
				return p
			}(),
			},
			err: errStructType,
		},
		{
			name: "err-estruct.db",
			cols: []interface{}{func() interface{} {
				type Name struct{}
				type Anon struct {
					Name Name
				}
				var v Anon
				return v
			}(),
			},
			err: errStructType,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			db, err := sql.Open("ql", "memory://"+tc.name)
			if err != nil {
				t.Fatalf("error creating db: %v\n", err)
			}
			defer db.Close()

			const ntname = "ntup"
			nt, err := Create(db, ntname, tc.cols...)
			if tc.err != nil && err == nil {
				t.Fatalf("expected an error")
			}
			if err != tc.err {
				t.Fatalf("got=%v. want=%v", err, tc.err)
			}
			if nt != nil {
				defer nt.DB().Close()
			}
		})
	}
}

func init() {
	var err error
	db, err := sql.Open("ql", "memory://mem.db")
	if err != nil {
		panic(err)
	}

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
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

	nt, err = Open(db, "data")
	if err != nil {
		panic(err)
	}
}
