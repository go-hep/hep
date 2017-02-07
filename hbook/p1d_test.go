// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"reflect"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/gonum/matrix/mat64"
	"github.com/gonum/stat/distmv"
)

func TestP1D(t *testing.T) {
	p := hbook.NewP1D(10, -4, +4)
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

func ExampleP1D() {
	const npoints = 1000

	p := hbook.NewP1D(100, -10, 10)
	dist, ok := distmv.NewNormal(
		[]float64{0, 1},
		mat64.NewSymDense(2, []float64{4, 0, 0, 2}),
		rand.New(rand.NewSource(1234)),
	)
	if !ok {
		log.Fatalf("error creating distmv.Normal")
	}

	v := make([]float64, 2)
	// Draw some random values from the standard
	// normal distribution.
	for i := 0; i < npoints; i++ {
		v = dist.Rand(v)
		p.Fill(v[0], v[1], 1)
	}

	fmt.Printf("mean:    %v\n", p.XMean())
	fmt.Printf("rms:     %v\n", p.XRMS())
	fmt.Printf("std-dev: %v\n", p.XStdDev())
	fmt.Printf("std-err: %v\n", p.XStdErr())

	// Output:
	// mean:    -0.04449868272082065
	// rms:     2.1327992781495637
	// std-dev: 2.1334019855956714
	// std-err: 0.06746409439208055
}

func TestP1DWriteYODA(t *testing.T) {
	p := hbook.NewP1D(10, -4, +4)
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

	ref, err := ioutil.ReadFile("testdata/p1d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h2d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestP1DReadYODA(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/p1d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h hbook.P1D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h1d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestP1DSerialization(t *testing.T) {
	pref := hbook.NewP1D(10, -4, +4)
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

		var pnew hbook.P1D
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
