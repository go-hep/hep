// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepevt_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hepevt"
	"golang.org/x/xerrors"
)

var (
	small = hepevt.Event{
		Nevhep: 1,
		Nhep:   8,
		Isthep: []int{3, 3, 3, 3, 3, 1, 1, 1},
		Idhep:  []int{2212, 2212, -2, 1, -24, -2, 1, 22},
		Jmohep: [][2]int{{0, 0}, {0, 0}, {2, 2}, {1, 1}, {3, 4}, {5, 5}, {5, 5}, {3, 4}},
		Jdahep: [][2]int{{4, 4}, {3, 3}, {5, 8}, {5, 8}, {6, 7}, {0, 0}, {0, 0}, {0, 0}},
		Phep: [][5]float64{
			{+0.00000000e+00, +0.00000000e+00, +7.00000000e+03, +7.0000000000e+03, +0.00000000e+00},
			{+0.00000000e+00, +0.00000000e+00, -7.00000000e+03, +7.0000000000e+03, +0.00000000e+00},
			{-3.04700000e+00, -1.90000000e+01, -5.46290000e+01, +5.7920000000e+01, +0.00000000e+00},
			{+7.50000000e-01, -1.56900000e+00, +3.21910000e+01, +3.2238000000e+01, +0.00000000e+00},
			{+1.51700000e+00, -2.06800000e+01, -2.06050000e+01, +8.5925000000e+01, +0.00000000e+00},
			{+3.96200000e+00, -4.94980000e+01, -2.66870000e+01, +5.6373000000e+01, +0.00000000e+00},
			{-2.44500000e+00, +2.88160000e+01, +6.08200000e+00, +2.9552000000e+01, +0.00000000e+00},
			{-3.81300000e+00, +1.13000000e-01, -1.83300000e+00, +4.2330000000e+00, +0.00000000e+00},
		},
		Vhep: [][4]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	}
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("testdata/small.hepevt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var evt hepevt.Event

	dec := hepevt.NewDecoder(f)
	err = dec.Decode(&evt)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(evt, small) {
		t.Fatalf("decoded event differ from reference:\ngot = %#v\nwant= %#v",
			evt, small,
		)
	}
}

func TestEncoder(t *testing.T) {
	f, err := ioutil.TempFile("", "hepevt-")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	enc := hepevt.NewEncoder(f)
	if err := enc.Encode(&small); err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/small.hepevt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("files differ.\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}

func TestDecoderFail(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		want  error
	}{
		{
			name:  "empty",
			input: "",
			want:  io.EOF,
		},
		{
			name:  "unexpected newline",
			input: "\n",
		},
		{
			name:  "expected integer",
			input: "1 t\n",
		},
		{
			name:  "newline in format does not match format",
			input: "1.1 1\n",
		},
		{
			name:  "newline in format does not match input",
			input: "1 1.1\n",
		},
		{
			name:  "newline in format does not match format",
			input: "1 1\n1\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader([]byte(tc.input))
			dec := hepevt.NewDecoder(r)
			var evt hepevt.Event
			err := dec.Decode(&evt)
			if err == nil {
				t.Fatalf("expected a failure")
			}
			if tc.want != nil && !xerrors.Is(err, tc.want) {
				t.Fatalf("unexpected error.\ngot = %v\nwant= %v", err, tc.want)
			}
			if tc.want == nil {
				t.Logf("%s: %q", tc.name, err.Error())
			}
		})
	}
}

func TestEncoderFail(t *testing.T) {
	for i := 0; i < 256; i++ {
		nbytes := i
		t.Run("", func(t *testing.T) {
			w := newWriter(nbytes)
			enc := hepevt.NewEncoder(w)
			err := enc.Encode(&small)
			if err == nil {
				t.Fatalf("expected a failure")
			}
		})
	}
}

type failWriter struct {
	n int
	c int
}

func newWriter(n int) *failWriter {
	return &failWriter{n: n, c: 0}
}

func (w *failWriter) Write(data []byte) (int, error) {
	for range data {
		w.c++
		if w.c >= w.n {
			return 0, io.ErrUnexpectedEOF
		}
	}
	return len(data), nil
}
