// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist_test

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleCreate_histo1D() {
	const fname = "h1d_example.root"
	defer os.Remove(fname)

	f, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	h.Fill(-10, 1) // fill underflow
	h.Fill(-20, 2)
	h.Fill(+10, 3) // fill overflow

	fmt.Printf("original histo:\n")
	fmt.Printf("w-mean:    %.7f\n", h.XMean())
	fmt.Printf("w-rms:     %.7f\n", h.XRMS())

	hroot := rhist.NewH1DFrom(h)

	err = f.Put("h1", hroot)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("h1")
	if err != nil {
		log.Fatal(err)
	}

	hr := rootcnv.H1D(robj.(rhist.H1))

	fmt.Printf("\nhisto read back:\n")
	fmt.Printf("r-mean:    %.7f\n", hr.XMean())
	fmt.Printf("r-rms:     %.7f\n", hr.XRMS())

	// Output:
	// original histo:
	// w-mean:    0.0023919
	// w-rms:     1.0628679
	//
	// histo read back:
	// r-mean:    0.0023919
	// r-rms:     1.0628679
}

func ExampleCreate_histo2D() {
	const fname = "h2d_example.root"
	defer os.Remove(fname)

	f, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	const npoints = 1000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH2D(5, -4, +4, 6, -4, +4)
	for i := 0; i < npoints; i++ {
		x := dist.Rand()
		y := dist.Rand()
		h.Fill(x, y, 1)
	}
	h.Fill(-10, -10, 1) // fill underflow
	h.Fill(-10, +10, 1)
	h.Fill(+10, -10, 1)
	h.Fill(+10, +10, 3) // fill overflow

	fmt.Printf("original histo:\n")
	fmt.Printf("w-mean-x:    %+.6f\n", h.XMean())
	fmt.Printf("w-rms-x:     %+.6f\n", h.XRMS())
	fmt.Printf("w-mean-y:    %+.6f\n", h.YMean())
	fmt.Printf("w-rms-y:     %+.6f\n", h.YRMS())

	hroot := rhist.NewH2DFrom(h)

	err = f.Put("h2", hroot)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("h2")
	if err != nil {
		log.Fatal(err)
	}

	hr := rootcnv.H2D(robj.(rhist.H2))

	fmt.Printf("\nhisto read back:\n")
	fmt.Printf("w-mean-x:    %+.6f\n", hr.XMean())
	fmt.Printf("w-rms-x:     %+.6f\n", hr.XRMS())
	fmt.Printf("w-mean-y:    %+.6f\n", hr.YMean())
	fmt.Printf("w-rms-y:     %+.6f\n", hr.YRMS())

	// Output:
	// original histo:
	// w-mean-x:    +0.046442
	// w-rms-x:     +1.231044
	// w-mean-y:    -0.018977
	// w-rms-y:     +1.253143
	//
	// histo read back:
	// w-mean-x:    +0.046442
	// w-rms-x:     +1.231044
	// w-mean-y:    -0.018977
	// w-rms-y:     +1.253143
}

func TestH1(t *testing.T) {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	h.Fill(-10, 1) // fill underflow
	h.Fill(-20, 2)
	h.Fill(+10, 1) // fill overflow
	h.Fill(+10, 2)
	h.Annotation()["name"] = "my-name"
	h.Annotation()["title"] = "my-title"

	for _, tc := range []struct {
		name   string
		h1     rhist.H1
		sumw   float64
		sumw2  float64
		sumwx  float64
		sumwx2 float64
	}{
		{
			name:   "TH1D",
			h1:     rhist.NewH1DFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
		{
			name:   "TH1F",
			h1:     rhist.NewH1FFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
		{
			name:   "TH1I",
			h1:     rhist.NewH1IFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.h1.SumW(), h.SumW(); got != want {
				t.Fatalf("sumw: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumW2(), h.SumW2(); got != want {
				t.Fatalf("sumw2: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumWX(), h.SumWX(); got != want {
				t.Fatalf("sumwx: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumWX2(), h.SumWX2(); got != want {
				t.Fatalf("sumwx2: got=%v, want=%v", got, want)
			}

			rraw, err := tc.h1.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			hh := rootcnv.H1D(tc.h1)

			hraw, err := hh.MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			var hr = rtypes.Factory.Get(tc.name)().Interface().(rhist.H1)
			if err := hr.(yodacnv.Unmarshaler).UnmarshalYODA(hraw); err != nil {
				t.Fatal(err)
			}

			rgot, err := hr.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(rgot, rraw) {
				t.Fatalf("round trip error:\nraw:\n%s\ngot:\n%s\n", rraw, rgot)
			}
		})
	}
}

func TestH2(t *testing.T) {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH2D(5, -4, +4, 6, -4, +4)
	for i := 0; i < npoints; i++ {
		x := dist.Rand()
		y := dist.Rand()
		h.Fill(x, y, 1)
	}
	h.Fill(+0, +5, 1) // N
	h.Fill(-5, +5, 2) // N-W
	h.Fill(-5, +0, 3) // W
	h.Fill(-5, -5, 4) // S-W
	h.Fill(+0, -5, 5) // S
	h.Fill(+5, -5, 6) // S-E
	h.Fill(+5, +0, 7) // E
	h.Fill(+5, +5, 8) // N-E

	h.Annotation()["name"] = "my-name"
	h.Annotation()["title"] = "my-title"

	for _, tc := range []struct {
		name   string
		h2     rhist.H2
		sumw   float64
		sumw2  float64
		sumwx  float64
		sumwx2 float64
	}{
		{
			name:   "TH2D",
			h2:     rhist.NewH2DFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
		{
			name:   "TH2F",
			h2:     rhist.NewH2FFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
		{
			name:   "TH2I",
			h2:     rhist.NewH2IFrom(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.h2.SumW(), h.SumW(); got != want {
				t.Fatalf("sumw: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumW2(), h.SumW2(); got != want {
				t.Fatalf("sumw2: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumWX(), h.SumWX(); got != want {
				t.Fatalf("sumwx: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumWX2(), h.SumWX2(); got != want {
				t.Fatalf("sumwx2: got=%v, want=%v", got, want)
			}

			rraw, err := tc.h2.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			hh := rootcnv.H2D(tc.h2)

			hraw, err := hh.MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			var hr = rtypes.Factory.Get(tc.name)().Interface().(rhist.H2)
			if err := hr.(yodacnv.Unmarshaler).UnmarshalYODA(hraw); err != nil {
				t.Fatal(err)
			}

			rgot, err := hr.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			// rounding errors... // FIXME(sbinet)
			rraw = bytes.Replace(rraw,
				[]byte("# Mean: (1.990041e-02, 2.039840e-04)"),
				[]byte("# Mean: (1.990041e-02, 2.039841e-04)"),
				-1,
			)
			if !bytes.Equal(rgot, rraw) {
				t.Fatalf("round trip error:\n%s",
					cmp.Diff(
						string(rraw),
						string(rgot),
					),
				)
			}
		})
	}
}
