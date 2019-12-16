// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"gonum.org/v1/plot/cmpimg"
)

func TestPrint(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-root-print-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(dir)

	refname := filepath.Join(dir, "ref.root")
	ref, err := groot.Create(refname)
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer ref.Close()

	dir111, err := riofs.Dir(ref).Mkdir("dir-1/dir-11/dir-111")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dir121, err := riofs.Dir(ref).Mkdir("dir-1/dir-12/dir-121")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	dir2, err := riofs.Dir(ref).Mkdir("dir-2")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h00 := hbook.NewH1D(10, 0, 10)
	h00.Annotation()["name"] = "h00"
	h00.Fill(5, 5)
	err = ref.Put("h00", rootcnv.FromH1D(h00))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h111 := hbook.NewH1D(10, 0, 10)
	h111.Annotation()["name"] = "h111"
	h111.Fill(5, 5)
	err = dir111.Put("h111", rootcnv.FromH1D(h111))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h121 := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	h121.Annotation()["name"] = "h121"
	h121.Fill(5, 5, 5)
	err = dir121.Put("h121", rootcnv.FromH2D(h121))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h21 := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	h21.Annotation()["name"] = "h21"
	h21.Fill(2, 1, 5)
	err = dir2.Put("h21", rootcnv.FromH2D(h21))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h22 := hbook.NewH2D(10, 0, 10, 10, 0, 10)
	h22.Annotation()["name"] = "h22"
	h22.Fill(2, 2, 5)
	err = dir2.Put("h22", rootcnv.FromH2D(h22))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	g22 := hbook.NewS2DFrom([]float64{1, 2, 3}, []float64{11, 12, 13})
	g22.Annotation()["name"] = "g22"
	err = dir2.Put("g22", rootcnv.FromS2D(g22))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	g23 := hbook.NewS2D([]hbook.Point2D{
		{X: 10, ErrX: hbook.Range{Min: 2, Max: 3}, Y: 10, ErrY: hbook.Range{Min: 2, Max: 3}},
		{X: 11, ErrX: hbook.Range{Min: 2, Max: 3}, Y: 11, ErrY: hbook.Range{Min: 2, Max: 3}},
	}...)
	g23.Annotation()["name"] = "g23"
	err = dir2.Put("g23", rootcnv.FromS2D(g23))
	if err != nil {
		t.Fatalf("%+v", err)
	}

	err = ref.Close()
	if err != nil {
		t.Fatalf("%+v", err)
	}

	for _, tc := range []struct {
		fname string
		otype string
		want  []string
	}{
		{
			fname: refname,
			otype: "png",
			want: []string{
				"h00.png",
				"h111.png",
				"h121.png",
				"h21.png",
				"h22.png",
				"g22.png",
				"g23.png",
			},
		},
		{
			fname: refname + ":g.*",
			otype: "png",
			want: []string{
				"g22.png",
				"g23.png",
			},
		},
		{
			fname: refname + ":dir",
			otype: "png",
			want: []string{
				"h111.png",
				"h121.png",
				"h21.png",
				"h22.png",
				"g22.png",
				"g23.png",
			},
		},
		{
			fname: refname + ":dir-2",
			otype: "png",
			want: []string{
				"h21.png",
				"h22.png",
				"g22.png",
				"g23.png",
			},
		},
		{
			fname: refname + ":dir-111",
			otype: "png",
			want: []string{
				"h111.png",
			},
		},
		{
			fname: refname + ":/dir-111",
			otype: "png",
			want: []string{
				"h111.png",
			},
		},
		{
			fname: refname + ":^/dir-111",
			otype: "png",
			want:  []string{},
		},
	} {
		tname := tc.fname
		tname = tname[len(dir)+1:]
		t.Run(tname, func(t *testing.T) {
			odir, err := ioutil.TempDir("", "groot-root-print-out-")
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer os.RemoveAll(odir)

			const verbose = false
			err = rootprint(odir, []string{tc.fname}, tc.otype, verbose)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			files, err := filepath.Glob(filepath.Join(odir, "*."+tc.otype))
			if err != nil {
				t.Fatalf("%+v", err)
			}

			if got, want := len(files), len(tc.want); got != want {
				t.Fatalf("invalid number of files: got=%d, want=%d", got, want)
			}

			for _, name := range files {
				got, err := ioutil.ReadFile(name)
				if err != nil {
					t.Fatalf("could not read file %q: %+v", name, err)
				}
				want, err := ioutil.ReadFile(filepath.Join("testdata", filepath.Base(name)))
				if err != nil {
					t.Fatalf("could not read ref file %q: %+v", name, err)
				}
				ok, err := cmpimg.Equal(tc.otype, got, want)
				if err != nil {
					t.Fatalf("could not compare %q: %+v", name, err)
				}
				if !ok {
					t.Fatalf("file %q does not compare equal", name)
				}
			}
		})
	}
}
