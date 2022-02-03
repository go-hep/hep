// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist_test

import (
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/hbook"
)

func TestGraph(t *testing.T) {
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tg")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rhist.Graph)
	if !ok {
		t.Fatalf("'tg' not a rhist.Graph: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
	}
}

func TestGraphErrors(t *testing.T) {
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tge")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rhist.GraphErrors)
	if !ok {
		t.Fatalf("'tge' not a rhist.GraphErrors: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
		xlo, xhi := g.XError(i)
		if want := 0.1 * v; want != xlo {
			t.Errorf("xerr[%d].low=%v want=%v", i, xlo, want)
		}
		if want := 0.1 * v; want != xhi {
			t.Errorf("xerr[%d].high=%v want=%v", i, xhi, want)
		}
		ylo, yhi := g.YError(i)
		if want := 0.2 * v; want != ylo {
			t.Errorf("yerr[%d].low=%v want=%v", i, ylo, want)
		}
		if want := 0.2 * v; want != yhi {
			t.Errorf("yerr[%d].high=%v want=%v", i, yhi, want)
		}
	}
}

func TestGraphAsymmErrors(t *testing.T) {
	f, err := groot.Open("../testdata/graphs.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tgae")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rhist.GraphErrors)
	if !ok {
		t.Fatalf("'tgae' not a rhist.GraphErrors: %T", obj)
	}

	if n, want := g.Len(), int(4); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	for i, v := range []float64{1, 2, 3, 4} {
		x, y := g.XY(i)
		if x != v {
			t.Errorf("x[%d]=%v. want=%v", i, x, v)
		}
		if y != 2*v {
			t.Errorf("y[%d]=%v. want=%v", i, y, 2*v)
		}
		xlo, xhi := g.XError(i)
		if want := 0.1 * v; want != xlo {
			t.Errorf("xerr[%d].low=%v want=%v", i, xlo, want)
		}
		if want := 0.2 * v; want != xhi {
			t.Errorf("xerr[%d].high=%v want=%v", i, xhi, want)
		}
		ylo, yhi := g.YError(i)
		if want := 0.3 * v; want != ylo {
			t.Errorf("yerr[%d].low=%v want=%v", i, ylo, want)
		}
		if want := 0.4 * v; want != yhi {
			t.Errorf("yerr[%d].high=%v want=%v", i, yhi, want)
		}
	}
}

func TestGraphMultiErrors(t *testing.T) {
	f, err := groot.Open("../testdata/tgme.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("gme")
	if err != nil {
		t.Fatal(err)
	}
	g, ok := obj.(rhist.GraphErrors)
	if !ok {
		t.Fatalf("'gme' not a rhist.GraphErrors: %T", obj)
	}

	if n, want := g.Len(), int(5); n != want {
		t.Errorf("npts=%d. want=%d\n", n, want)
	}

	var (
		xs   = []float64{0, 1, 2, 3, 4}
		ys   = []float64{0, 2, 4, 1, 3}
		xlos = []float64{0.3, 0.3, 0.3, 0.3, 0.3}
		xhis = []float64{0.3, 0.3, 0.3, 0.3, 0.3}
		ylos = []float64{1, 0.5, 1, 0.5, 1}
		yhis = []float64{0.5, 1, 0.5, 1, 2}
	)
	for i := 0; i < g.Len(); i++ {
		x, y := g.XY(i)
		if x != xs[i] {
			t.Errorf("x[%d]=%v. want=%v", i, x, xs[i])
		}
		if y != ys[i] {
			t.Errorf("y[%d]=%v. want=%v", i, y, ys[i])
		}
		xlo, xhi := g.XError(i)
		if want := xlos[i]; want != xlo {
			t.Errorf("xerr[%d].low=%v want=%v", i, xlo, want)
		}
		if want := xhis[i]; want != xhi {
			t.Errorf("xerr[%d].high=%v want=%v", i, xhi, want)
		}
		ylo, yhi := g.YError(i)
		if want := ylos[i]; want != ylo {
			t.Errorf("yerr[%d].low=%v want=%v", i, ylo, want)
		}
		if want := yhis[i]; want != yhi {
			t.Errorf("yerr[%d].high=%v want=%v", i, yhi, want)
		}
	}
}

func TestInvalidGraphMerger(t *testing.T) {
	var (
		gr = hbook.NewS2D([]hbook.Point2D{
			{X: 0, Y: 0, ErrX: hbook.Range{Min: 1, Max: 1}, ErrY: hbook.Range{Min: 2, Max: 2}},
			{X: 1, Y: 1, ErrX: hbook.Range{Min: 1, Max: 1}, ErrY: hbook.Range{Min: 2, Max: 2}},
		}...)
		src = rbase.NewObjString("foo")
	)
	for _, tc := range []struct {
		name string
		obj  root.Merger
		want string
	}{
		{
			name: "graph",
			obj:  rhist.NewGraphFrom(gr).(root.Merger),
			want: "rhist: can not merge *rbase.ObjString into *rhist.tgraph",
		},
		{
			name: "graph-errs",
			obj:  rhist.NewGraphErrorsFrom(gr).(root.Merger),
			want: "rhist: can not merge *rbase.ObjString into *rhist.tgrapherrs",
		},
		{
			name: "graph-asymmerr",
			obj:  rhist.NewGraphAsymmErrorsFrom(gr).(root.Merger),
			want: "rhist: can not merge *rbase.ObjString into *rhist.tgraphasymmerrs",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.obj.ROOTMerge(src)
			if err == nil {
				t.Fatalf("expected an error")
			}

			if got, want := err.Error(), tc.want; got != want {
				t.Fatalf("invalid ROOTMerge error. got=%q, want=%q", got, want)
			}
		})
	}
}
