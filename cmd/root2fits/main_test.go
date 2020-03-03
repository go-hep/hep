// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/astrogo/fitsio"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func TestConvert(t *testing.T) {
	tmp, err := ioutil.TempDir("", "root2fits-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name  string
		nevts int
		data  func(i int) interface{}
		want  string
	}{
		{
			name:  "builtins",
			nevts: 5,
			data: func(i int) interface{} {
				type D struct {
					B   bool
					I8  int8
					I16 int16
					I32 int32
					I64 int64
					U8  uint8
					U16 uint16
					U32 uint32
					U64 uint64
					F32 float32
					F64 float64
					Str string
				}
				return &D{
					B:   i%2 == 0,
					I8:  int8(-i),
					I16: int16(-i),
					I32: int32(-i),
					I64: int64(-i),
					U8:  uint8(i),
					U16: uint16(i),
					U32: uint32(i),
					U64: uint64(i),
					F32: float32(i),
					F64: float64(i),
					Str: fmt.Sprintf("%05d", i),
				}
			},
			want: `== 00001/00005 =================================================================
B          | true
I8         | 0
I16        | 0
I32        | 0
I64        | 0
U8         | 0
U16        | 0
U32        | 0
U64        | 0
F32        | 0
F64        | 0
Str        | 00000
== 00002/00005 =================================================================
B          | false
I8         | 255
I16        | -1
I32        | -1
I64        | -1
U8         | 1
U16        | 1
U32        | 1
U64        | 1
F32        | 1
F64        | 1
Str        | 00001
== 00003/00005 =================================================================
B          | true
I8         | 254
I16        | -2
I32        | -2
I64        | -2
U8         | 2
U16        | 2
U32        | 2
U64        | 2
F32        | 2
F64        | 2
Str        | 00002
== 00004/00005 =================================================================
B          | false
I8         | 253
I16        | -3
I32        | -3
I64        | -3
U8         | 3
U16        | 3
U32        | 3
U64        | 3
F32        | 3
F64        | 3
Str        | 00003
== 00005/00005 =================================================================
B          | true
I8         | 252
I16        | -4
I32        | -4
I64        | -4
U8         | 4
U16        | 4
U32        | 4
U64        | 4
F32        | 4
F64        | 4
Str        | 00004
`,
		},
		{
			name:  "arrays",
			nevts: 5,
			data: func(i int) interface{} {
				type D struct {
					B   [3]bool
					I8  [3]int8
					I16 [3]int16
					I32 [3]int32
					I64 [3]int64
					U8  [3]uint8
					U16 [3]uint16
					U32 [3]uint32
					U64 [3]uint64
					F32 [3]float32
					F64 [3]float64
				}
				return &D{
					B:   [3]bool{i%2 == 0, (i+1)%2 == 0, (i+2)%2 == 0},
					I8:  [3]int8{int8(-i), int8(-i - 1), int8(-i - 2)},
					I16: [3]int16{int16(-i), int16(-i - 1), int16(-i - 2)},
					I32: [3]int32{int32(-i), int32(-i - 1), int32(-i - 2)},
					I64: [3]int64{int64(-i), int64(-i - 1), int64(-i - 2)},
					U8:  [3]uint8{uint8(i), uint8(i + 1), uint8(i + 2)},
					U16: [3]uint16{uint16(i), uint16(i + 1), uint16(i + 2)},
					U32: [3]uint32{uint32(i), uint32(i + 1), uint32(i + 2)},
					U64: [3]uint64{uint64(i), uint64(i + 1), uint64(i + 2)},
					F32: [3]float32{float32(10 + i), float32(20 + i), float32(30 + i)},
					F64: [3]float64{float64(10 + i), float64(20 + i), float64(30 + i)},
				}
			},
			want: `== 00001/00005 =================================================================
B          | [true false true]
I8         | [0 255 254]
I16        | [0 -1 -2]
I32        | [0 -1 -2]
I64        | [0 -1 -2]
U8         | [0 1 2]
U16        | [0 1 2]
U32        | [0 1 2]
U64        | [0 1 2]
F32        | [10 20 30]
F64        | [10 20 30]
== 00002/00005 =================================================================
B          | [false true false]
I8         | [255 254 253]
I16        | [-1 -2 -3]
I32        | [-1 -2 -3]
I64        | [-1 -2 -3]
U8         | [1 2 3]
U16        | [1 2 3]
U32        | [1 2 3]
U64        | [1 2 3]
F32        | [11 21 31]
F64        | [11 21 31]
== 00003/00005 =================================================================
B          | [true false true]
I8         | [254 253 252]
I16        | [-2 -3 -4]
I32        | [-2 -3 -4]
I64        | [-2 -3 -4]
U8         | [2 3 4]
U16        | [2 3 4]
U32        | [2 3 4]
U64        | [2 3 4]
F32        | [12 22 32]
F64        | [12 22 32]
== 00004/00005 =================================================================
B          | [false true false]
I8         | [253 252 251]
I16        | [-3 -4 -5]
I32        | [-3 -4 -5]
I64        | [-3 -4 -5]
U8         | [3 4 5]
U16        | [3 4 5]
U32        | [3 4 5]
U64        | [3 4 5]
F32        | [13 23 33]
F64        | [13 23 33]
== 00005/00005 =================================================================
B          | [true false true]
I8         | [252 251 250]
I16        | [-4 -5 -6]
I32        | [-4 -5 -6]
I64        | [-4 -5 -6]
U8         | [4 5 6]
U16        | [4 5 6]
U32        | [4 5 6]
U64        | [4 5 6]
F32        | [14 24 34]
F64        | [14 24 34]
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				fname = filepath.Join(tmp, tc.name+".root")
				tname = "tree"
				oname = filepath.Join(tmp, tc.name+".fits")
			)
			// create
			func() {
				f, err := groot.Create(fname)
				if err != nil {
					t.Fatalf("could not create write ROOT file %q: %v", fname, err)
				}
				defer f.Close()

				ptr := tc.data(0)
				wvars := rtree.WriteVarsFromStruct(ptr)
				tw, err := rtree.NewWriter(f, tname, wvars)
				if err != nil {
					t.Fatalf("could not create tree writer: %v", err)
				}

				for i := 0; i < int(tc.nevts); i++ {
					want := reflect.ValueOf(tc.data(i)).Elem().Interface()
					for j, wvar := range wvars {
						v := reflect.ValueOf(wvar.Value).Elem()
						want := reflect.ValueOf(want).Field(j)
						v.Set(want)
					}
					_, err = tw.Write()
					if err != nil {
						t.Fatalf("could not write event %d: %v", i, err)
					}
				}

				err = tw.Close()
				if err != nil {
					t.Fatalf("could not close tree writer: %v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("could not close write ROOT file %q: %v", fname, err)
				}
			}()

			err := process(oname, tname, fname)
			if err != nil {
				t.Fatalf("could not convert ROOT tree to FITS table: %+v", err)
			}

			got := new(strings.Builder)
			err = display(got, tname, oname)
			if err != nil {
				t.Fatalf("could not display FITS table content: %+v", err)
			}

			if got, want := got.String(), tc.want; got != want {
				t.Fatalf("invalid FITS table content.\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}

func display(o io.Writer, hname, fname string) error {
	r, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open file %q: %w", fname, err)
	}
	defer r.Close()

	f, err := fitsio.Open(r)
	if err != nil {
		return fmt.Errorf("could not open FITS file %q: %w", fname, err)
	}
	defer f.Close()

	hdu := f.Get(hname)
	if hdu.Type() == fitsio.IMAGE_HDU {
		return fmt.Errorf("HDU %q not a table", hname)
	}

	table := hdu.(*fitsio.Table)
	ncols := len(table.Cols())
	nrows := table.NumRows()
	rows, err := table.Read(0, nrows)
	if err != nil {
		return fmt.Errorf("could not read FITS table range: %w", err)
	}
	hdrline := strings.Repeat("=", 80-15)
	maxname := 10
	for _, col := range table.Cols() {
		if len(col.Name) > maxname {
			maxname = len(col.Name)
		}
	}

	data := make([]interface{}, ncols)
	names := make([]string, ncols)
	for i, col := range table.Cols() {
		names[i] = col.Name
		data[i] = reflect.New(col.Type()).Interface()
	}

	rowfmt := fmt.Sprintf("%%-%ds | %%v\n", maxname)
	for irow := 0; rows.Next(); irow++ {
		err = rows.Scan(data...)
		if err != nil {
			return fmt.Errorf("could not read row %d: %w", irow, err)
		}
		fmt.Fprintf(o, "== %05d/%05d %s\n", irow+1, nrows, hdrline)
		for i := 0; i < ncols; i++ {
			rv := reflect.Indirect(reflect.ValueOf(data[i]))
			fmt.Fprintf(o, rowfmt, names[i], rv.Interface())
		}
	}

	err = rows.Err()
	if err != nil {
		return fmt.Errorf("could not scan table: %w", err)
	}

	return nil
}
