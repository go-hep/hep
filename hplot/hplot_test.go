// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

type chkplotFunc func(ExampleFunc func(), t *testing.T, filenames ...string)

func checkPlot(f chkplotFunc) chkplotFunc {
	return func(ex func(), t *testing.T, filenames ...string) {
		t.Helper()
		f(ex, t, filenames...)
		if t.Failed() {
			return
		}
		for _, fname := range filenames {
			_ = os.Remove(path.Join("testdata", fname))
		}
	}
}

func TestSubPlot(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(Example_subplot, t, "sub_plot.png")
}

func TestLatexPlot(t *testing.T) {
	Example_latexplot()
	ref, err := ioutil.ReadFile("testdata/latex_plot_golden.tex")
	if err != nil {
		t.Fatal(err)
	}
	chk, err := ioutil.ReadFile("testdata/latex_plot.tex")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ref, chk) {
		t.Fatal("files testdata/latex_plot{,_golden}.tex differ\n")
	}
	os.Remove("testdata/latex_plot.tex")
}

func TestShow(t *testing.T) {
	p := hplot.New()
	p.Title.Text = "title"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	for _, tc := range []struct {
		w, h   vg.Length
		format string
	}{
		{-1, -1, ""},
		{-1, -1, "png"},
		{-1, -1, "jpg"},
		{-1, -1, "pdf"},
		{-1, -1, "eps"},
		{-1, -1, "tex"},
		{-1, 10 * vg.Centimeter, ""},
		{10 * vg.Centimeter, -1, ""},
	} {
		t.Run(tc.format, func(t *testing.T) {
			_, err := hplot.Show(p, tc.w, tc.h, tc.format)
			if err != nil {
				t.Fatalf("%+v", err)
			}
		})
	}
}
