// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func checkPlot(t *testing.T, ref string) {
	fname := strings.Replace(ref, "_golden", "", 1)

	if *cmpimg.GenerateTestData {
		got, _ := ioutil.ReadFile(fname)
		ioutil.WriteFile(ref, got, 0644)
	}

	want, err := ioutil.ReadFile(ref)
	if err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	ext := filepath.Ext(ref)[1:]
	if ok, err := cmpimg.Equal(ext, got, want); !ok || err != nil {
		if err != nil {
			t.Fatalf("error: comparing %q with reference file: %v\n", fname, err)
		} else {
			t.Fatalf("error: %q differ with reference file\n", fname)
		}
	}
	os.Remove(fname)
}

func TestSubPlot(t *testing.T) {
	cmpimg.CheckPlot(Example_subplot, t, "sub_plot.png")
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
