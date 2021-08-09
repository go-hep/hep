// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/astrogo/fitsio"
	"go-hep.org/x/hep/groot/rcmd"
)

func TestConvert(t *testing.T) {
	tmp, err := os.MkdirTemp("", "fits2root-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name string
		cols []fitsio.Column
		data interface{}
		want string
	}{
		{
			name: "bools",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "L",
				},
			},
			data: []bool{true, false, true, false, true},
			want: `key[000]: test;1 "" (TTree)
[000][col]: true
[001][col]: false
[002][col]: true
[003][col]: false
[004][col]: true
`,
		},
		{
			name: "i8",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "B",
				},
			},
			data: []int8{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "i16",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "I",
				},
			},
			data: []int16{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "i32",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "J",
				},
			},
			data: []int32{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "i64",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "K",
				},
			},
			data: []int64{-10, -11, -12, -13, -14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: -10
[001][col]: -11
[002][col]: -12
[003][col]: -13
[004][col]: -14
`,
		},
		{
			name: "u8",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "B",
				},
			},
			data: []uint8{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "u16",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "I",
				},
			},
			data: []uint16{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "u32",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "J",
				},
			},
			data: []uint32{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "u64",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "K",
				},
			},
			data: []uint64{10, 11, 12, 13, 14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: 10
[001][col]: 11
[002][col]: 12
[003][col]: 13
[004][col]: 14
`,
		},
		{
			name: "f32",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "E",
				},
			},
			data: []float32{-10, -11, -12, -13, -14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: -10
[001][col]: -11
[002][col]: -12
[003][col]: -13
[004][col]: -14
`,
		},
		{
			name: "f64",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "D",
				},
			},
			data: []float64{-10, -11, -12, -13, -14},
			want: `key[000]: test;1 "" (TTree)
[000][col]: -10
[001][col]: -11
[002][col]: -12
[003][col]: -13
[004][col]: -14
`,
		},
		{
			name: "strings",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "10A",
				},
			},
			data: []string{"a", "", "c ", " d", "eée", " "},
			want: `key[000]: test;1 "" (TTree)
[000][col]: a
[001][col]: 
[002][col]: c 
[003][col]:  d
[004][col]: eée
[005][col]:  
`,
		},
		{
			name: "2df64",
			cols: []fitsio.Column{
				{
					Name:   "col",
					Format: "2D",
				},
			},
			data: [][2]float64{{10, 11}, {12, 13}, {14, 15}, {16, 17}, {18, 19}},
			want: `key[000]: test;1 "" (TTree)
[000][col]: [10 11]
[001][col]: [12 13]
[002][col]: [14 15]
[003][col]: [16 17]
[004][col]: [18 19]
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				fname = filepath.Join(tmp, tc.name+".fits")
				oname = filepath.Join(tmp, tc.name+".root")
			)

			// create
			func() {
				w, err := os.Create(fname)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				defer w.Close()

				f, err := fitsio.Create(w)
				if err != nil {
					t.Fatalf("could not create input FITS file: %+v", err)
				}
				defer f.Close()

				phdu, err := fitsio.NewPrimaryHDU(nil)
				if err != nil {
					t.Fatalf("could not create primary hdu: %+v", err)
				}
				err = f.Write(phdu)
				if err != nil {
					t.Fatalf("could not write primary hdu: %+v", err)
				}

				tbl, err := fitsio.NewTable("test", tc.cols, fitsio.BINARY_TBL)
				if err != nil {
					t.Fatalf("could not create FITS table: %+v", err)
				}
				defer tbl.Close()

				rslice := reflect.ValueOf(tc.data)
				for i := 0; i < rslice.Len(); i++ {
					data := rslice.Index(i).Addr()
					err = tbl.Write(data.Interface())
					if err != nil {
						t.Fatalf("could not write row [%v]: %+v", i, err)
					}
				}

				err = f.Write(tbl)
				if err != nil {
					t.Fatalf("could not write FITS table: %+v", err)
				}

				err = f.Close()
				if err != nil {
					t.Fatalf("could not close FITS file: %+v", err)
				}
			}()

			err := process(oname, "test", fname)
			if err != nil {
				t.Fatalf("could not convert FITS file: %+v", err)
			}

			// read-back
			func() {
				const deep = true
				got := new(strings.Builder)
				err := rcmd.Dump(got, oname, deep, nil)
				if err != nil {
					t.Fatalf("could not run root-dump: %+v", err)
				}

				if got, want := got.String(), tc.want; got != want {
					t.Fatalf("fits2root conversion failed:\ngot:\n%s\nwant:\n%s\n", got, want)
				}
			}()
		})
	}
}
