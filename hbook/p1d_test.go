// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestP1D(t *testing.T) {
	p := NewP1D(10, -4, +4)
	if p == nil {
		t.Fatalf("nil pointer to P1D")
	}

	p.Annotation()["name"] = "p1d"

	for i := 0; i < 10; i++ {
		v := float64(i)
		p.Fill(v, v*2, 1)
	}
	p.Fill(-10, 10, 1)

	if got, want := p.Name(), "p1d"; got != want {
		t.Errorf("got=%q. want=%q\n", got, want)
	}

	for _, test := range []struct {
		name string
		f    func() float64
		want float64
	}{
		{
			name: "xmean",
			f:    p.XMean,
			want: 3.1818181818181817,
		},
		{
			name: "xmin",
			f:    p.XMin,
			want: -4.0,
		},
		{
			name: "xmax",
			f:    p.XMax,
			want: +4.0,
		},
		{
			name: "xrms",
			f:    p.XRMS,
			want: 5.916079783099616,
		},
		{
			name: "xstddev",
			f:    p.XStdDev,
			want: 5.231026320296655,
		},
		{
			name: "xstderr",
			f:    p.XStdErr,
			want: 1.5772137793543157,
		},
		{
			name: "xvariance",
			f:    p.XVariance,
			want: 27.363636363636363,
		},
		{
			name: "sumw",
			f:    p.SumW,
			want: 11.0,
		},
		{
			name: "sumw2",
			f:    p.SumW2,
			want: 11.0,
		},
	} {
		got := test.f()
		if got != test.want {
			t.Errorf("test: %v. got=%v. want=%v\n", test.name, got, test.want)
		}
	}
}

func TestP1DWriteYODA(t *testing.T) {
	p := NewP1D(10, -4, +4)
	if p == nil {
		t.Fatalf("nil pointer to P1D")
	}

	for i := 0; i < 10; i++ {
		v := float64(i)
		p.Fill(v, v*2, 1)
	}
	p.Fill(-10, 10, 1)

	chk, err := p.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile("testdata/p1d_v2_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("p1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestP1DReadYODAv1(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/p1d_v1_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h P1D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.marshalYODAv1()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("p1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestP1DReadYODAv2(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/p1d_v2_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h P1D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("p1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestP1DSerialization(t *testing.T) {
	pref := NewP1D(10, -4, +4)
	if pref == nil {
		t.Fatalf("nil pointer to P1D")
	}

	for i := 0; i < 10; i++ {
		v := float64(i)
		pref.Fill(v, v*2, 1)
	}
	pref.Fill(-10, 10, 1)

	pref.Annotation()["title"] = "p1d title"
	pref.Annotation()["name"] = "p1d-name"

	{
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(pref)
		if err != nil {
			t.Fatalf("could not serialize p1d: %v\n", err)
		}

		var pnew P1D
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&pnew)
		if err != nil {
			t.Fatalf("could not deserialize p1d: %v\n", err)
		}

		if !reflect.DeepEqual(pref, &pnew) {
			t.Fatalf("ref=%v\nnew=%v\n", pref, &pnew)
		}
	}
}
