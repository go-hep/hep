// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"reflect"
	"sync"
	"testing"

	"go-hep.org/x/hep/hbook"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/plotter"
)

func ExampleH1D() {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:     0,
		Sigma:  1,
		Source: rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}

	fmt.Printf("mean:    %v\n", h.XMean())
	fmt.Printf("rms:     %v\n", h.XRMS())
	fmt.Printf("std-dev: %v\n", h.XStdDev())
	fmt.Printf("std-err: %v\n", h.XStdErr())

	// Output:
	// mean:    -0.017341752512581167
	// rms:     0.9913281479386786
	// std-dev: 0.9912260153046148
	// std-err: 0.00991226015304615
}

func TestH1D(t *testing.T) {
	h1 := hbook.NewH1D(100, 0., 100.)
	if h1 == nil {
		t.Errorf("nil pointer to H1D")
	}

	h1.Annotation()["name"] = "h1"

	n := h1.Name()
	if n != "h1" {
		t.Errorf("expected H1D.Name() == %q (got %q)\n",
			n, "h1")
	}
	xmin := h1.XMin()
	if xmin != 0. {
		t.Errorf("expected H1D.Min() == %v (got %v)\n",
			0., xmin,
		)
	}
	xmax := h1.XMax()
	if xmax != 100. {
		t.Errorf("expected H1D.Max() == %v (got %v)\n",
			100., xmax,
		)
	}
	/*
		for idx := 0; idx < nbins; idx++ {
			size := h1.Binning().BinWidth(idx)
			if size != 1. {
				t.Errorf("expected H1D.Binning.BinWidth(%v) == %v (got %v)\n",
					idx, 1., size,
				)
			}
		}
	*/
	var _ plotter.XYer = h1
	var _ plotter.Valuer = h1
}

func TestH1DEdges(t *testing.T) {
	h := hbook.NewH1DFromEdges([]float64{
		-4.0, -3.6, -3.2, -2.8, -2.4, -2.0, -1.6, -1.2, -0.8, -0.4,
		+0.0, +0.4, +0.8, +1.2, +1.6, +2.0, +2.4, +2.8, +3.2, +3.6,
		+4.0,
	})
	if got, want := h.XMin(), -4.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h.XMax(), +4.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}

	bins := hbook.Bin1Ds(h.Binning().Bins())
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -4.1, want: hbook.UnderflowBin},
		{v: -4.0, want: 0},
		{v: -3.6, want: 1},
		{v: -3.2, want: 2},
		{v: -2.8, want: 3},
		{v: -2.4, want: 4},
		{v: -2.0, want: 5},
		{v: -1.6, want: 6},
		{v: -1.2, want: 7},
		{v: -0.8, want: 8},
		{v: -0.4, want: 9},
		{v: +0.0, want: 10},
		{v: +0.4, want: 11},
		{v: +0.8, want: 12},
		{v: +1.2, want: 13},
		{v: +1.6, want: 14},
		{v: +2.0, want: 15},
		{v: +2.4, want: 16},
		{v: +2.8, want: 17},
		{v: +3.2, want: 18},
		{v: +3.6, want: 19},
		{v: +4.0, want: hbook.OverflowBin},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}
}

func TestH1DBins(t *testing.T) {
	h := hbook.NewH1DFromBins([]hbook.Range{
		{Min: -4.0, Max: -3.6}, {Min: -3.6, Max: -3.2}, {Min: -3.2, Max: -2.8}, {Min: -2.8, Max: -2.4}, {Min: -2.4, Max: -2.0},
		{Min: -2.0, Max: -1.6}, {Min: -1.6, Max: -1.2}, {Min: -1.2, Max: -0.8}, {Min: -0.8, Max: -0.4}, {Min: -0.4, Max: +0.0},
		{Min: +0.0, Max: +0.4}, {Min: +0.4, Max: +0.8}, {Min: +0.8, Max: +1.2}, {Min: +1.2, Max: +1.6}, {Min: +1.6, Max: +2.0},
		{Min: +2.0, Max: +2.4}, {Min: +2.4, Max: +2.8}, {Min: +2.8, Max: +3.2}, {Min: +3.2, Max: +3.6}, {Min: +3.6, Max: +4.0},
	}...)
	if got, want := h.XMin(), -4.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h.XMax(), +4.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}
	bins := hbook.Bin1Ds(h.Binning().Bins())
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -4.1, want: hbook.UnderflowBin},
		{v: -4.0, want: 0},
		{v: -3.6, want: 1},
		{v: -3.2, want: 2},
		{v: -2.8, want: 3},
		{v: -2.4, want: 4},
		{v: -2.0, want: 5},
		{v: -1.6, want: 6},
		{v: -1.2, want: 7},
		{v: -0.8, want: 8},
		{v: -0.4, want: 9},
		{v: +0.0, want: 10},
		{v: +0.4, want: 11},
		{v: +0.8, want: 12},
		{v: +1.2, want: 13},
		{v: +1.6, want: 14},
		{v: +2.0, want: 15},
		{v: +2.4, want: 16},
		{v: +2.8, want: 17},
		{v: +3.2, want: 18},
		{v: +3.6, want: 19},
		{v: +4.0, want: hbook.OverflowBin},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}

}

func TestH1DBinsWithGaps(t *testing.T) {
	h1 := hbook.NewH1DFromBins([]hbook.Range{
		{Min: -10, Max: -5}, {Min: -5, Max: 0}, {Min: 0, Max: 4} /*GAP*/, {Min: 5, Max: 10},
	}...)
	if got, want := h1.XMin(), -10.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h1.XMax(), 10.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}

	bins := hbook.Bin1Ds(h1.Binning().Bins())
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -20, want: hbook.UnderflowBin},
		{v: -10, want: 0},
		{v: -9, want: 0},
		{v: -5, want: 1},
		{v: 0, want: 2},
		{v: 4.5, want: len(bins)},
		{v: 5, want: 3},
		{v: 6, want: 3},
		{v: 10, want: hbook.OverflowBin},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}

	h := hbook.NewH1DFromBins([]hbook.Range{
		{Min: 0, Max: 1}, {Min: 1, Max: 2}, {Min: 3, Max: 4},
	}...)
	h.Fill(0, 1)
	h.Fill(1, 1)
	h.Fill(2, 1) // gap
	h.Fill(3, 1)

	raw, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/h1d_gaps_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("got:\n%s\nwant:\n%s\n", string(raw), string(want))
	}

	var href hbook.H1D
	err = href.UnmarshalYODA(want)
	if err != nil {
		t.Fatal(err)
	}

	raw, err = href.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("got:\n%s\nwant:\n%s\n", string(raw), string(want))
	}
}

func TestH1DEdgesWithPanics(t *testing.T) {
	for _, test := range [][]float64{
		{0},
		{0, 1, 0.5, 2},
		{0, 1, 1},
		{0, 1, 0, 1},
		{0, 1, 2, 2},
		{0, 1, 2, 2, 2},
	} {
		panicked, _ := panics(func() {
			_ = hbook.NewH1DFromEdges(test)
		})
		if !panicked {
			t.Fatalf("edges %v should have panicked", test)
		}
	}
}

func TestH1DBinsWithPanics(t *testing.T) {
	for _, test := range [][]hbook.Range{
		{{Min: 0, Max: 1}, {Min: 0.5, Max: 1.5}},
		{{Min: 0, Max: 1}, {Min: -1, Max: 2}},
		{{Min: 0, Max: 1.5}, {Min: -1, Max: 1}},
		{{Min: 0, Max: 1}, {Min: 0.5, Max: 0.6}},
	} {
		panicked, _ := panics(func() {
			_ = hbook.NewH1DFromBins(test...)
		})
		if !panicked {
			t.Fatalf("bins %v should have panicked", test)
		}
	}
}

func TestH1DIntegral(t *testing.T) {
	h1 := hbook.NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1.3)
	h1.Fill(0, 1)
	h1.Fill(0.5, 1.6)
	h1.Fill(1.2, 0.7)
	h1.Fill(2.1, 0.3)
	h1.Fill(4.2, 1.2)
	h1.Fill(5.9, 0.5)
	h1.Fill(6, 2.1)

	integral := h1.Integral()
	if integral != 8.7 {
		t.Errorf("expected H1D.Integral() == 8.7 (got %v)\n", integral)
	}
	if got, want := h1.Integral(h1.XMin(), h1.XMax()), 5.3; got != want {
		t.Errorf("H1D.Integral(xmin,xmax) = %v. want=%v\n", got, want)
	}

	integralall := h1.Integral(math.Inf(-1), math.Inf(+1))
	if integralall != 8.7 {
		t.Errorf("expected H1D.Integral(math.Inf(-1), math.Inf(+1)) == 8.7 (got %v)\n", integralall)
	}
	integralu := h1.Integral(math.Inf(-1), h1.XMax())
	if integralu != 6.6 {
		t.Errorf("expected H1D.Integral(math.Inf(-1), h1.Binning().UpperEdge()) == 6.6 (got %v)\n", integralu)
	}
	integralo := h1.Integral(h1.XMin(), math.Inf(+1))
	if integralo != 7.4 {
		t.Errorf("expected H1D.Integral(h1.Binning().LowerEdge(), math.Inf(+1)) == 7.4 (got %v)\n", integralo)
	}
	integralrange := h1.Integral(0.5, 5.5)
	if integralrange != 2.7 {
		t.Errorf("expected H1D.Integral(0.5, 5.5) == 2.7 (got %v)\n", integralrange)
	}

	mean1, rms1 := h1.XMean(), h1.XRMS()

	h1.Scale(1 / integral)
	integral = h1.Integral()
	if integral != 1 {
		t.Errorf("expected H1D.Integral() == 1 (got %v)\n", integral)
	}

	mean2, rms2 := h1.XMean(), h1.XRMS()

	if math.Abs(mean1-mean2)/mean1 > 1e-12 {
		t.Errorf("mean has changed while rescaling (mean1, mean2) = (%v, %v)", mean1, mean2)
	}
	if math.Abs(rms1-rms2)/rms1 > 1e-12 {
		t.Errorf("rms has changed while rescaling (rms1, rms2) = (%v, %v)", rms1, rms2)
	}

	h2 := hbook.NewH1D(2, 0, 1)
	h2.Fill(0.0, 1)
	h2.Fill(0.5, 1)
	for _, ibin := range []int{0, 1} {
		if got, want := h2.Value(ibin), 1.0; got != want {
			t.Errorf("got H1D.Value(%d) = %v. want %v\n", ibin, got, want)
		}
	}

	if got, want := h2.Integral(), 2.0; got != want {
		t.Errorf("got H1D.Integral() == %v. want %v\n", got, want)
	}
}

func TestH1DNegativeWeights(t *testing.T) {
	h1 := hbook.NewH1D(5, 0, 100)
	h1.Fill(10, -200)
	h1.Fill(20, 1)
	h1.Fill(30, 0.2)
	h1.Fill(10, +200)

	h2 := hbook.NewH1D(5, 0, 100)
	h2.Fill(20, 1)
	h2.Fill(30, 0.2)

	const tol = 1e-12
	if x1, x2 := h1.XMean(), h2.XMean(); !floats.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("mean differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	if x1, x2 := h1.XRMS(), h2.XRMS(); !floats.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("rms differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	/* FIXME
	if x1, x2 := h1.StdDevX(), h2.StdDevX(); !floats.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("std-dev differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	*/
}

func BenchmarkH1DSTFillConst(b *testing.B) {
	b.StopTimer()
	h1 := hbook.NewH1D(100, 0., 100.)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		h1.Fill(10., 1.)
	}
}

func BenchmarkH1DFillFlat(b *testing.B) {
	b.StopTimer()
	h1 := hbook.NewH1D(100, 0., 100.)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		h1.Fill(rnd()*100., 1.)
	}
}

func BenchmarkH1DFillFlatGo(b *testing.B) {
	b.StopTimer()
	h1 := hbook.NewH1D(100, 0., 100.)
	wg := new(sync.WaitGroup)
	//wg.Add(b.N)
	b.StartTimer()

	// throttle...
	q := make(chan struct{}, 1000)
	for i := 0; i < b.N; i++ {
		q <- struct{}{}
		go func() {
			wg.Add(1)
			h1.Fill(rnd()*100., 1.)
			<-q
			wg.Done()
		}()
	}
	wg.Wait()
}

func st_process_evts(n int, hists []*hbook.H1D, process func(hists []*hbook.H1D)) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			process(hists)
			wg.Done()
		}()
	}
	wg.Wait()
}

func st_process_evts_const(hists []*hbook.H1D) {
	for _, h := range hists {
		h.Fill(10., 1.)
	}
}
func BenchmarkNH1DFillConst(b *testing.B) {
	b.StopTimer()
	hists := make([]*hbook.H1D, 100)
	for i := 0; i < 100; i++ {
		hists[i] = hbook.NewH1D(100, 0., 100.)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		st_process_evts(100, hists, st_process_evts_const)
	}
}

func st_process_evts_flat(hists []*hbook.H1D) {
	for _, h := range hists {
		h.Fill(rnd()*100., 1.)
	}
}

func BenchmarkNH1DFillFlat(b *testing.B) {
	b.StopTimer()
	hists := make([]*hbook.H1D, 100)
	for i := 0; i < 100; i++ {
		hists[i] = hbook.NewH1D(100, 0., 100.)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		st_process_evts(100, hists, st_process_evts_flat)
	}
}

func TestH1DSerialization(t *testing.T) {
	const nentries = 50
	href := hbook.NewH1D(100, 0., 100.)
	for i := 0; i < nentries; i++ {
		href.Fill(rnd()*100., 1.)
	}
	href.Annotation()["title"] = "histo title"
	href.Annotation()["name"] = "histo name"

	// test gob.GobDecode/gob.GobEncode interface
	func() {
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(href)
		if err != nil {
			t.Fatalf("could not serialize histogram: %v", err)
		}

		var hnew hbook.H1D
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&hnew)
		if err != nil {
			t.Fatalf("could not deserialize histogram: %v", err)
		}

		if !reflect.DeepEqual(href, &hnew) {
			t.Fatalf("ref=%v\nnew=%v\n", href, &hnew)
		}
	}()

	// test rio.Marshaler/Unmarshaler
	func() {
		buf := new(bytes.Buffer)
		err := href.RioMarshal(buf)
		if err != nil {
			t.Fatalf("could not serialize histogram: %v", err)
		}

		var hnew hbook.H1D
		err = hnew.RioUnmarshal(buf)
		if err != nil {
			t.Fatalf("could not deserialize histogram: %v", err)
		}

		if !reflect.DeepEqual(href, &hnew) {
			t.Fatalf("ref=%v\nnew=%v\n", href, &hnew)
		}
	}()

}

func TestH1DWriteYODA(t *testing.T) {
	h := hbook.NewH1D(10, -4, 4)
	h.Fill(1, 1)
	h.Fill(2, 1)
	h.Fill(-3, 1)
	h.Fill(-4, 1)
	h.Fill(0, 1)
	h.Fill(0, 1)
	h.Fill(10, 1)
	h.Fill(-10, 1)

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile("testdata/h1d_golden.yoda")
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

func TestH1DReadYODA(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/h1d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h hbook.H1D
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
