// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"bytes"
	"encoding/gob"
	"math"
	"reflect"
	"sync"
	"testing"

	"github.com/go-hep/hbook"
	"github.com/gonum/plot/plotter"
)

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
	nbins := h1.Binning().Bins()
	if nbins != 100 {
		t.Errorf("expected H1D.Binning.Bins() == %v (got %v)\n",
			100, nbins,
		)
	}
	low := h1.Binning().LowerEdge()
	if low != 0. {
		t.Errorf("expected H1D.Binning.LowerEdge() == %v (got %v)\n",
			0., low,
		)
	}
	up := h1.Binning().UpperEdge()
	if up != 100. {
		t.Errorf("expected H1D.Binning.UpperEdge() == %v (got %v)\n",
			100., up,
		)
	}

	for idx := 0; idx < nbins; idx++ {
		size := h1.Binning().BinWidth(idx)
		if size != 1. {
			t.Errorf("expected H1D.Binning.BinWidth(%v) == %v (got %v)\n",
				idx, 1., size,
			)
		}
	}

	var _ plotter.XYer = h1
	var _ plotter.Valuer = h1
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
	if integral != 5.3 {
		t.Errorf("expected H1D.Integral() == 5.3 (got %v)\n", integral)
	}
	integralall := h1.Integral(math.Inf(-1), math.Inf(+1))
	if integralall != 8.7 {
		t.Errorf("expected H1D.Integral(math.Inf(-1), math.Inf(+1)) == 8.7 (got %v)\n", integralall)
	}
	integralu := h1.Integral(math.Inf(-1), h1.Binning().UpperEdge())
	if integralu != 6.6 {
		t.Errorf("expected H1D.Integral(math.Inf(-1), h1.Binning().UpperEdge()) == 6.6 (got %v)\n", integralu)
	}
	integralo := h1.Integral(h1.Binning().LowerEdge(), math.Inf(+1))
	if integralo != 7.4 {
		t.Errorf("expected H1D.Integral(h1.Binning().LowerEdge(), math.Inf(+1)) == 7.4 (got %v)\n", integralo)
	}
	integralrange := h1.Integral(0.5, 5.5)
	if integralrange != 2.7 {
		t.Errorf("expected H1D.Integral(0.5, 5.5) == 2.7 (got %v)\n", integralrange)
	}

	mean1, rms1 := h1.Mean(), h1.RMS()

	h1.Scale(1 / integral)
	integral = h1.Integral()
	if integral != 1 {
		t.Errorf("expected H1D.Integral() == 1 (got %v)\n", integral)
	}

	mean2, rms2 := h1.Mean(), h1.RMS()

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
	if got, want := h2.Binning().BinWidth(0), 0.5; got != want {
		t.Errorf("got H1D.Binning().BinWidth == %v. want %v\n", got, want)
	}
	if got, want := h2.Integral(), 1.0; got != want {
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

	if h1.Mean() != h2.Mean() {
		t.Errorf("mean differ:\nh1=%v\nh2=%v\n", h1.Mean(), h2.Mean())
	}
	if h1.RMS() != h2.RMS() {
		t.Errorf("rms differ:\nh1=%v\nh2=%v\n", h1.RMS(), h2.RMS())
	}
	// FIXME(sbinet)
	// if h1.StdDev() != h2.StdDev() {
	//	t.Errorf("std-dev differ:\nh1=%v\nh2=%v\n", h1.StdDev(), h2.StdDev())
	//}
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

// EOF
