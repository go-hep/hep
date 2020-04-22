// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"compress/gzip"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
)

func TestYODA2ROOT(t *testing.T) {

	h1 := hbook.NewH1D(10, -4, 4)
	h1.Annotation()["name"] = "h1-name"
	h1.Annotation()["title"] = "h1-title"
	h1.Fill(1, 1)
	h1.Fill(2, 1)

	h2 := hbook.NewH2D(10, -4, 4, 20, -5, 5)
	h2.Annotation()["name"] = "h2-name"
	h2.Annotation()["title"] = "h2-title"
	h2.Fill(1, 1, 1)
	h2.Fill(2, 2, 2)

	s2 := hbook.NewS2DFrom([]float64{1, 2, 3, 4, 5}, []float64{1, 4, 9, 16, 25})
	s2.Annotation()["name"] = "s2-name"
	s2.Annotation()["title"] = "s2-title"

	anon := hbook.NewS2DFrom([]float64{10, 20}, []float64{10, 40})
	anon.Annotation()["title"] = "no-title"

	for _, tc := range []struct {
		yfname string
		rfname string
	}{
		{
			yfname: "f1.yoda",
			rfname: "f1.root",
		},
		{
			yfname: "f2.yoda.gz",
			rfname: "f2.root",
		},
	} {
		t.Run(tc.yfname, func(t *testing.T) {
			yfname := tc.yfname
			rfname := tc.rfname

			defer os.Remove(yfname)
			defer os.Remove(rfname)

			f, err := os.Create(yfname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			var yf io.WriteCloser = f

			if strings.HasSuffix(yfname, ".gz") {
				yf = gzip.NewWriter(f)
			}

			err = yodacnv.Write(yf, h1, h2, s2, anon)
			if err != nil {
				t.Fatal(err)
			}

			err = yf.Close()
			if err != nil {
				t.Fatal(err)
			}

			if strings.HasSuffix(yfname, ".gz") {
				err = f.Close()
				if err != nil {
					t.Fatal(err)
				}
			}

			o, err := groot.Create(rfname)
			if err != nil {
				t.Fatal(err)
			}
			defer o.Close()

			err = convert(o, yfname)
			if err != nil {
				t.Fatal(err)
			}

			err = o.Close()
			if err != nil {
				t.Fatal(err)
			}

			rf, err := groot.Open(rfname)
			if err != nil {
				t.Fatal(err)
			}

			robj, err := rf.Get("h1-name")
			if err != nil {
				t.Fatal(err)
			}

			rh1 := rootcnv.H1D(robj.(rhist.H1))

			if got, want := rh1.XMean(), h1.XMean(); got != want {
				t.Fatalf("h1 round-trip failed: got: %v, want: %v", got, want)
			}

			robj, err = rf.Get("h2-name")
			if err != nil {
				t.Fatal(err)
			}

			rh2 := rootcnv.H2D(robj.(rhist.H2))

			if got, want := rh2.XMean(), h2.XMean(); got != want {
				t.Fatalf("h2 round-trip failed: got: %v, want: %v", got, want)
			}

			robj, err = rf.Get("s2-name")
			if err != nil {
				t.Fatal(err)
			}

			rs2 := rootcnv.S2D(robj.(rhist.GraphErrors))

			if !reflect.DeepEqual(rs2, s2) {
				t.Fatalf("s2 round-trip failed")
			}

			robj, err = rf.Get("yoda-scatter-003")
			if err != nil {
				t.Fatal(err)
			}

			ranon := rootcnv.S2D(robj.(rhist.GraphErrors))

			if !reflect.DeepEqual(ranon.Points(), anon.Points()) {
				t.Fatalf("s2-anon round-trip failed")
			}
		})
	}

}
