// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sbinet/npyio"
)

func TestRecord(t *testing.T) {
	for _, tc := range []struct {
		name   string
		forder bool // Fortran order
		want   interface{}
	}{
		{
			name: "float64_4x3x2-c-order",
			want: [4][3][2]float64{
				{{10, 11}, {12, 13}, {14, 15}},
				{{16, 17}, {18, 19}, {20, 21}},
				{{22, 23}, {24, 25}, {26, 27}},
				{{28, 29}, {30, 31}, {32, 33}},
			},
		},
		{
			name:   "float64_4x3x2-f-order",
			forder: true,
			want: [2][3][4]float64{
				{{10, 11, 12, 13}, {14, 15, 16, 17}, {18, 19, 20, 21}},
				{{22, 23, 24, 25}, {26, 27, 28, 29}, {30, 31, 32, 33}},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "npy2root-")
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer os.RemoveAll(tmp)

			fname := filepath.Join(tmp, "data.npy")
			src, err := os.Create(fname)
			if err != nil {
				t.Fatalf("could not create NumPy data file: %+v", err)
			}
			defer src.Close()

			err = npyio.Write(src, tc.want)
			if err != nil {
				t.Fatalf("could not save NumPy data file: %+v", err)
			}

			err = src.Close()
			if err != nil {
				t.Fatalf("could not close NumPy data file: %+v", err)
			}

			src, err = os.Open(fname)
			if err != nil {
				t.Fatalf("could not reopen NumPy data file: %+v", err)
			}
			defer src.Close()

			npy, err := npyio.NewReader(src)
			if err != nil {
				t.Fatalf("could not create numpy file reader: %+v", err)
			}

			npy.Header.Descr.Fortran = tc.forder

			rec := NewRecord(npy)
			defer rec.Release()

			rec.Retain()
			rec.Release()

			if got, want := rec.NumRows(), int64(4); got != want {
				t.Fatalf("invalid number of rows: got=%d, want=%d", got, want)
			}

			if got, want := rec.NumCols(), int64(1); got != want {
				t.Fatalf("invalid number of cols: got=%d, want=%d", got, want)
			}

			if got, want := rec.ColumnName(0), "numpy"; got != want {
				t.Fatalf("invalid column name: got=%q, want=%q", got, want)
			}
		})
	}
}
