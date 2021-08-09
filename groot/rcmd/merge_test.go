// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
)

func TestMerge(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-root-merge-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	type funcT func(t *testing.T, fname string) error
	for _, tc := range []struct {
		name   string
		inputs []funcT
		output funcT
		panics string
	}{
		{
			name:   "flat-tree-1",
			inputs: []funcT{makeFlatTree(1)},
			output: makeFlatTree(1),
		},
		{
			name:   "flat-tree-2",
			inputs: []funcT{makeFlatTree(1), makeFlatTree(1)},
			output: makeFlatTree(2),
		},
		{
			name:   "h1f-1",
			inputs: []funcT{makeH1F(1)},
			output: makeH1F(1),
		},
		{
			name:   "h1f-2",
			inputs: []funcT{makeH1F(1), makeH1F(1)},
			output: makeH1F(2),
		},
		{
			name:   "h1d-1",
			inputs: []funcT{makeH1D(1)},
			output: makeH1D(1),
		},
		{
			name:   "h1d-2",
			inputs: []funcT{makeH1D(1), makeH1D(1)},
			output: makeH1D(2),
		},
		{
			name:   "h1i-1",
			inputs: []funcT{makeH1I(1)},
			output: makeH1I(1),
		},
		{
			name:   "h1i-2",
			inputs: []funcT{makeH1I(1), makeH1I(1)},
			output: makeH1I(2),
		},
		{
			name:   "h2d-1",
			inputs: []funcT{makeH2D(1)},
			output: makeH2D(1),
		},
		{
			name:   "h2d-2",
			inputs: []funcT{makeH2D(1), makeH2D(1)},
			output: makeH2D(2),
			panics: "not implemented", // FIXME(sbinet)
		},
		{
			name:   "graph-1",
			inputs: []funcT{makeGraph(0, 1)},
			output: makeGraph(0, 1),
		},
		{
			name:   "graph-2",
			inputs: []funcT{makeGraph(0, 1), makeGraph(1, 2)},
			output: makeGraph(0, 2),
		},
		{
			name:   "graph-err-1",
			inputs: []funcT{makeGraphErr(0, 1)},
			output: makeGraphErr(0, 1),
		},
		{
			name:   "graph-err-2",
			inputs: []funcT{makeGraphErr(0, 1), makeGraphErr(1, 2)},
			output: makeGraphErr(0, 2),
		},
		{
			name:   "graph-asymmerr-1",
			inputs: []funcT{makeGraphAsymmErr(0, 1)},
			output: makeGraphAsymmErr(0, 1),
		},
		{
			name:   "graph-asymmerr-2",
			inputs: []funcT{makeGraphAsymmErr(0, 1), makeGraphAsymmErr(1, 2)},
			output: makeGraphAsymmErr(0, 2),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				fnames  []string
				oname   = filepath.Join(tmp, tc.name+".out.root")
				verbose = true
				deep    = true
			)
			for i, fct := range tc.inputs {
				fname := filepath.Join(tmp, fmt.Sprintf("%s-%02d.root", tc.name, i))
				err := fct(t, fname)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				fnames = append(fnames, fname)
			}
			refname := filepath.Join(tmp, tc.name+".want.root")
			err := tc.output(t, refname)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(string), tc.panics; got != want {
						t.Fatalf("invalid panic message. got=%q, want=%q", got, want)
					}
				}()
			}

			err = rcmd.Merge(oname, fnames, verbose)
			if err != nil {
				t.Fatalf("could not run root-merge: %+v", err)
			}

			got := new(bytes.Buffer)
			err = rcmd.Dump(got, oname, deep, nil)
			if err != nil {
				t.Fatalf("could not run root-dump: %+v", err)
			}

			want := new(bytes.Buffer)
			err = rcmd.Dump(want, refname, deep, nil)
			if err != nil {
				t.Fatalf("could not run root-dump: %+v", err)
			}

			if got, want := got.String(), want.String(); got != want {
				t.Fatalf("invalid root-merge output:\ngot:\n%swant:\n%s", got, want)
			}
		})
	}
}

func makeFlatTree(n int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		type Data struct {
			I32    int32
			F64    float64
			Str    string
			ArrF64 [5]float64
			N      int32
			SliF64 []float64 `groot:"SliF64[N]"`
		}
		const (
			nevts = 5
		)

		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		dir, err := riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		var evt Data
		tree, err := rtree.NewWriter(dir, "mytree", rtree.WriteVarsFromStruct(&evt))
		if err != nil {
			t.Fatalf("could not create tree writer: %+v", err)
		}

		for j := 0; j < n; j++ {
			for i := 0; i < nevts; i++ {
				evt.I32 = int32(i)
				evt.F64 = float64(i)
				evt.Str = fmt.Sprintf("evt-%0d", i)
				evt.ArrF64 = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
				evt.N = int32(i)
				evt.SliF64 = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
				_, err = tree.Write()
				if err != nil {
					t.Fatalf("could not write event %d: %+v", i, err)
				}
			}
		}

		err = tree.Close()
		if err != nil {
			t.Fatalf("could not write tree: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeH1F(n int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		h := hbook.NewH1D(10, 0, 10)
		h.Annotation()["title"] = "h1f"
		for i := 0; i < n; i++ {
			h.Fill(5, 1)
			h.Fill(6, 2)
		}

		err = dir21.Put("h1f", rhist.NewH1FFrom(h))
		if err != nil {
			t.Fatalf("could not save H1F: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeH1D(n int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		h := hbook.NewH1D(10, 0, 10)
		h.Annotation()["title"] = "h1d"
		for i := 0; i < n; i++ {
			h.Fill(5, 1)
			h.Fill(6, 2)
		}

		err = dir21.Put("h1d", rootcnv.FromH1D(h))
		if err != nil {
			t.Fatalf("could not save H1D: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeH1I(n int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		h := hbook.NewH1D(10, 0, 10)
		h.Annotation()["title"] = "h1i"
		for i := 0; i < n; i++ {
			h.Fill(5, 1)
			h.Fill(6, 2)
		}

		err = dir21.Put("h1i", rhist.NewH1IFrom(h))
		if err != nil {
			t.Fatalf("could not save H1I: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeH2D(n int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		h := hbook.NewH2D(10, 0, 10, 10, 0, 10)
		h.Annotation()["title"] = "h2d"
		for i := 0; i < n; i++ {
			h.Fill(5, 5, 1)
			h.Fill(6, 6, 2)
		}

		err = dir21.Put("h2d", rootcnv.FromH2D(h))
		if err != nil {
			t.Fatalf("could not save H2D: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeGraph(beg, end int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		var (
			xs []float64
			ys []float64
		)
		for i := beg; i < end; i++ {
			for j := 0; j < 10; j++ {
				xs = append(xs, float64(10*(1+i)+j))
				ys = append(ys, float64(10*(1+i)+j))
			}
		}

		gr := hbook.NewS2DFrom(xs, ys)
		gr.Annotation()["title"] = "graph"
		err = dir21.Put("graph", rhist.NewGraphFrom(gr))
		if err != nil {
			t.Fatalf("could not save S2D: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeGraphErr(beg, end int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		var (
			pts []hbook.Point2D
		)
		for i := beg; i < end; i++ {
			for j := 0; j < 10; j++ {
				var (
					x    = float64(10*(1+i) + j)
					y    = float64(10*(1+i) + j)
					xerr = 2.5
					yerr = 3.5
				)
				pts = append(pts, hbook.Point2D{
					X:    x,
					Y:    y,
					ErrX: hbook.Range{Min: xerr, Max: xerr},
					ErrY: hbook.Range{Min: yerr, Max: yerr},
				})
			}
		}

		gr := hbook.NewS2D(pts...)
		gr.Annotation()["title"] = "graph"
		err = dir21.Put("graph", rhist.NewGraphErrorsFrom(gr))
		if err != nil {
			t.Fatalf("could not save S2D: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}

func makeGraphAsymmErr(beg, end int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		f, err := groot.Create(fname)
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		_, err = riofs.Dir(f).Mkdir("dir-1/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		dir21, err := riofs.Dir(f).Mkdir("dir-2/dir-11")
		if err != nil {
			t.Fatalf("could not create directory: %+v", err)
		}

		var (
			pts []hbook.Point2D
		)
		for i := beg; i < end; i++ {
			for j := 0; j < 10; j++ {
				var (
					x = float64(10*(1+i) + j)
					y = float64(10*(1+i) + j)
				)
				pts = append(pts, hbook.Point2D{
					X:    x,
					Y:    y,
					ErrX: hbook.Range{Min: 1.5, Max: 2.5},
					ErrY: hbook.Range{Min: 1.5, Max: 2.5},
				})
			}
		}

		gr := hbook.NewS2D(pts...)
		gr.Annotation()["title"] = "graph"
		err = dir21.Put("graph", rhist.NewGraphAsymmErrorsFrom(gr))
		if err != nil {
			t.Fatalf("could not save S2D: %+v", err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close file: %+v", err)
		}

		return nil
	}
}
