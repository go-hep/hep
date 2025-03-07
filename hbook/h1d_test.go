// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/plot/plotter"
)

func panics(fn func()) (panicked bool, message string) {
	defer func() {
		r := recover()
		panicked = r != nil
		message = fmt.Sprint(r)
	}()
	fn()
	return
}

func TestH1D(t *testing.T) {
	h1 := NewH1D(100, 0., 100.)
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
	h := NewH1DFromEdges([]float64{
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

	bins := Bin1Ds(h.Binning.Bins)
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -4.1, want: UnderflowBin1D},
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
		{v: +4.0, want: OverflowBin1D},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}
}

func TestH1DBins(t *testing.T) {
	h := NewH1DFromBins([]Range{
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
	bins := Bin1Ds(h.Binning.Bins)
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -4.1, want: UnderflowBin1D},
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
		{v: +4.0, want: OverflowBin1D},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}

}

func TestH1DBinsWithGapsv1(t *testing.T) {
	h1 := NewH1DFromBins([]Range{
		{Min: -10, Max: -5}, {Min: -5, Max: 0}, {Min: 0, Max: 4} /*GAP*/, {Min: 5, Max: 10},
	}...)
	if got, want := h1.XMin(), -10.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h1.XMax(), 10.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}

	bins := Bin1Ds(h1.Binning.Bins)
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -20, want: UnderflowBin1D},
		{v: -10, want: 0},
		{v: -9, want: 0},
		{v: -5, want: 1},
		{v: 0, want: 2},
		{v: 4.5, want: len(bins)},
		{v: 5, want: 3},
		{v: 6, want: 3},
		{v: 10, want: OverflowBin1D},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}

	h := NewH1DFromBins([]Range{
		{Min: 0, Max: 1}, {Min: 1, Max: 2}, {Min: 3, Max: 4},
	}...)
	h.Fill(0, 1)
	h.Fill(1, 1)
	h.Fill(2, 1) // gap
	h.Fill(3, 1)
	h.Annotation()["title"] = "my-title"

	raw, err := h.marshalYODAv1()
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/h1d_gaps_v1_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(raw),
			),
		)
	}

	var href H1D
	err = href.UnmarshalYODA(want)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	raw, err = href.marshalYODAv1()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(raw),
			),
		)
	}
}

func TestH1DBinsWithGapsv2(t *testing.T) {
	h1 := NewH1DFromBins([]Range{
		{Min: -10, Max: -5}, {Min: -5, Max: 0}, {Min: 0, Max: 4} /*GAP*/, {Min: 5, Max: 10},
	}...)
	if got, want := h1.XMin(), -10.0; got != want {
		t.Errorf("got xmin=%v. want=%v", got, want)
	}
	if got, want := h1.XMax(), 10.0; got != want {
		t.Errorf("got xmax=%v. want=%v", got, want)
	}

	bins := Bin1Ds(h1.Binning.Bins)
	for _, test := range []struct {
		v    float64
		want int
	}{
		{v: -20, want: UnderflowBin1D},
		{v: -10, want: 0},
		{v: -9, want: 0},
		{v: -5, want: 1},
		{v: 0, want: 2},
		{v: 4.5, want: len(bins)},
		{v: 5, want: 3},
		{v: 6, want: 3},
		{v: 10, want: OverflowBin1D},
	} {
		idx := bins.IndexOf(test.v)
		if idx != test.want {
			t.Errorf("invalid index for %v. got=%d. want=%d\n", test.v, idx, test.want)
		}
	}

	h := NewH1DFromBins([]Range{
		{Min: 0, Max: 1}, {Min: 1, Max: 2}, {Min: 3, Max: 4},
	}...)
	h.Fill(0, 1)
	h.Fill(1, 1)
	h.Fill(2, 1) // gap
	h.Fill(3, 1)
	h.Annotation()["title"] = "my-title"

	raw, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	want, err := os.ReadFile("testdata/h1d_gaps_v2_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(raw),
			),
		)
	}

	var href H1D
	err = href.UnmarshalYODA(want)
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	raw, err = href.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(raw, want) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(raw),
			),
		)
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
			_ = NewH1DFromEdges(test)
		})
		if !panicked {
			t.Fatalf("edges %v should have panicked", test)
		}
	}
}

func TestH1DBinsWithPanics(t *testing.T) {
	for _, test := range [][]Range{
		{{Min: 0, Max: 1}, {Min: 0.5, Max: 1.5}},
		{{Min: 0, Max: 1}, {Min: -1, Max: 2}},
		{{Min: 0, Max: 1.5}, {Min: -1, Max: 1}},
		{{Min: 0, Max: 1}, {Min: 0.5, Max: 0.6}},
	} {
		panicked, _ := panics(func() {
			_ = NewH1DFromBins(test...)
		})
		if !panicked {
			t.Fatalf("bins %v should have panicked", test)
		}
	}
}

func TestH1DIntegral(t *testing.T) {
	h1 := NewH1D(6, 0, 6)
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

	h2 := NewH1D(2, 0, 1)
	h2.Fill(0.0, 1)
	h2.Fill(0.5, 1)
	for _, ibin := range []int{0, 1} {
		if got, want := h2.Value(ibin), 1.0; got != want {
			t.Errorf("got H1D.Value(%d) = %v. want %v\n", ibin, got, want)
		}
		if got, want := h2.Error(ibin), 1.0; got != want {
			t.Errorf("got H1D.Error(%d) = %v. want %v\n", ibin, got, want)
		}
	}

	if got, want := h2.Integral(), 2.0; got != want {
		t.Errorf("got H1D.Integral() == %v. want %v\n", got, want)
	}
}

func TestH1DNegativeWeights(t *testing.T) {
	h1 := NewH1D(5, 0, 100)
	h1.Fill(10, -200)
	h1.Fill(20, 1)
	h1.Fill(30, 0.2)
	h1.Fill(10, +200)

	h2 := NewH1D(5, 0, 100)
	h2.Fill(20, 1)
	h2.Fill(30, 0.2)

	const tol = 1e-12
	if x1, x2 := h1.XMean(), h2.XMean(); !scalar.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("mean differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	if x1, x2 := h1.XRMS(), h2.XRMS(); !scalar.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("rms differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	/* FIXME
	if x1, x2 := h1.StdDevX(), h2.StdDevX(); !scalar.EqualWithinAbs(x1, x2, tol) {
		t.Errorf("std-dev differ:\nh1=%v\nh2=%v\n", x1, x2)
	}
	*/
}

func TestH1DSerialization(t *testing.T) {
	const nentries = 50
	href := NewH1D(100, 0., 100.)
	for range nentries {
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

		var hnew H1D
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

		var hnew H1D
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
	h := NewH1D(10, -4, 4)
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

	ref, err := os.ReadFile("testdata/h1d_v2_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		fatalf := t.Fatalf
		if runtime.GOOS == "darwin" {
			// ignore errors for darwin and mac-silicon
			fatalf = t.Logf
		}
		fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestH1DReadYODAv1(t *testing.T) {
	ref, err := os.ReadFile("testdata/h1d_v1_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h H1D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.marshalYODAv1()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestH1DReadYODAv2(t *testing.T) {
	ref, err := os.ReadFile("testdata/h1d_v2_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var h H1D
	err = h.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := h.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(ref),
				string(chk),
			),
		)
	}
}

func TestH1DBin(t *testing.T) {
	h := NewH1DFromEdges([]float64{
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

	h.Fill(-4.0, 1)
	h.Fill(-3.6, 1)
	h.Fill(-3.6, 1)
	h.Fill(-3.1, 1)
	h.Fill(-3.1, 1)
	h.Fill(-3.1, 1)

	for _, tc := range []struct {
		x   float64
		bin int
	}{
		{-4.0, 1},
		{-3.9, 1},
		{-3.6, 2},
		{-3.1, 3},
		{-10, -1},
		{+4, -1},
	} {
		t.Run(fmt.Sprintf("x=%v", tc.x), func(t *testing.T) {
			bin := h.Bin(tc.x)
			if tc.bin < 0 && bin == nil {
				// ok
				return
			}
			if bin == nil {
				t.Fatalf("unexpected nil bin")
			}

			if bin.EffEntries() != float64(tc.bin) {
				t.Fatalf("x=%v got=%v %v, want=%d", tc.x, bin.EffEntries(), bin.Entries(), tc.bin)
			}
		})
	}
}

func TestH1DFillN(t *testing.T) {
	h1 := NewH1D(10, 0, 10)
	h2 := NewH1D(10, 0, 10)

	xs := []float64{1, 2, 3, 4}
	ws := []float64{1, 2, 1, 1}

	for i := range xs {
		h1.Fill(xs[i], ws[i])
	}
	h2.FillN(xs, ws)

	if s1, s2 := h1.SumW(), h2.SumW(); s1 != s2 {
		t.Fatalf("invalid sumw: h1=%v, h2=%v", s1, s2)
	}

	for i := range xs {
		h1.Fill(xs[i], 1)
	}
	h2.FillN(xs, nil)

	if s1, s2 := h1.SumW(), h2.SumW(); s1 != s2 {
		t.Fatalf("invalid sumw: h1=%v, h2=%v", s1, s2)
	}

	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatalf("expected a panic!")
			}
			const want = "hbook: lengths mismatch"
			if got, want := err.(error).Error(), want; got != want {
				t.Fatalf("invalid panic message:\ngot= %q\nwant=%q", got, want)
			}
		}()

		h2.FillN(xs, []float64{1})
	}()
}

func TestH1DClone(t *testing.T) {
	h1 := NewH1D(10, 0, 10)
	h1.FillN(
		[]float64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
	)
	h1.Ann["hello"] = "world"

	msg1, err := h1.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h1: %+v", err)
	}

	h2 := h1.Clone()
	msg2, err := h2.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h2: %+v", err)
	}

	if !bytes.Equal(msg1, msg2) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(msg1),
				string(msg2),
			),
		)
	}

	h1.Ann["world"] = "bye"
	h1.FillN(
		[]float64{-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 11},
		[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 1, 1},
	)

	msg3, err := h1.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h1: %+v", err)
	}

	msg4, err := h2.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h2: %+v", err)
	}

	if bytes.Equal(msg1, msg3) {
		t.Fatalf("msg1/msg3 should differ")
	}

	if !bytes.Equal(msg4, msg2) {
		t.Fatalf("h1d file differ:\n%s\n",
			cmp.Diff(
				string(msg4),
				string(msg2),
			),
		)
	}
}
