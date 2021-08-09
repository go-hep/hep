// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist_test

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
)

func TestCreate(t *testing.T) {

	dir, err := os.MkdirTemp("", "groot-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for i, tc := range []struct {
		Name string
		Skip bool
		Want []rtests.ROOTer
		ROOT string
	}{
		{
			Name: "TAxis",
			Want: []rtests.ROOTer{rhist.NewAxis("xaxis")},
			ROOT: "retrieved: [xaxis]\n",
		},
		{
			Name: "TH1I",
			ROOT: "retrieved: [h1i]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H1I {
					h := hbook.NewH1D(100, 0, 100)
					h.Annotation()["name"] = "h1i"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, 1)
					h.Fill(+200, 1)
					h.Fill(1, 1)
					h.Fill(2, 1)
					h.Fill(3, 10)
					return rhist.NewH1IFrom(h)
				}(),
			},
		},
		{
			Name: "TH1F",
			ROOT: "retrieved: [h1f]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H1F {
					h := hbook.NewH1D(100, 0, 100)
					h.Annotation()["name"] = "h1f"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, 1)
					h.Fill(+200, 1)
					h.Fill(1, 1)
					h.Fill(2, 1)
					h.Fill(3, 10)
					return rhist.NewH1FFrom(h)
				}(),
			},
		},
		{
			Name: "TH1D",
			ROOT: "retrieved: [h1d]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H1D {
					h := hbook.NewH1D(100, 0, 100)
					h.Annotation()["name"] = "h1d"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, 1)
					h.Fill(+200, 1)
					h.Fill(1, 1)
					h.Fill(2, 1)
					h.Fill(3, 10)
					return rhist.NewH1DFrom(h)
				}(),
			},
		},
		{
			Name: "TH2I",
			ROOT: "retrieved: [h2i]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H2I {
					h := hbook.NewH2D(100, 0, 100, 50, 0, 50)
					h.Annotation()["name"] = "h2i"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, -1, 1)
					h.Fill(+200, 200, 1)
					h.Fill(1, 1, 1)
					h.Fill(2, 2, 1)
					h.Fill(3, 3, 10)
					return rhist.NewH2IFrom(h)
				}(),
			},
		},
		{
			Name: "TH2F",
			ROOT: "retrieved: [h2f]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H2F {
					h := hbook.NewH2D(100, 0, 100, 50, 0, 50)
					h.Annotation()["name"] = "h2f"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, -1, 1)
					h.Fill(+200, 200, 1)
					h.Fill(1, 1, 1)
					h.Fill(2, 2, 1)
					h.Fill(3, 3, 10)
					return rhist.NewH2FFrom(h)
				}(),
			},
		},
		{
			Name: "TH2D",
			ROOT: "retrieved: [h2d]\n",
			Want: []rtests.ROOTer{
				func() *rhist.H2D {
					h := hbook.NewH2D(100, 0, 100, 50, 0, 50)
					h.Annotation()["name"] = "h2d"
					h.Annotation()["title"] = "my title"
					h.Fill(-1, -1, 1)
					h.Fill(+200, 200, 1)
					h.Fill(1, 1, 1)
					h.Fill(2, 2, 1)
					h.Fill(3, 3, 10)
					return rhist.NewH2DFrom(h)
				}(),
			},
		},
		{
			Name: "TGraph",
			ROOT: "retrieved: [tg]\n",
			Want: []rtests.ROOTer{
				func() rtests.ROOTer {
					hg := hbook.NewS2D(
						hbook.Point2D{X: 1, Y: 1},
						hbook.Point2D{X: 2, Y: 1.5},
						hbook.Point2D{X: -1, Y: +2},
					)
					hg.Annotation()["name"] = "tg"
					hg.Annotation()["title"] = "my title"
					return rhist.NewGraphFrom(hg).(rtests.ROOTer)
				}(),
			},
		},
		{
			Name: "TGraphErrors",
			ROOT: "retrieved: [tge]\n",
			Want: []rtests.ROOTer{
				func() rtests.ROOTer {
					hg := hbook.NewS2D(
						hbook.Point2D{X: 1, Y: 1, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
						hbook.Point2D{X: 2, Y: 1.5, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
						hbook.Point2D{X: -1, Y: +2, ErrX: hbook.Range{Min: 2, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 3}},
					)
					hg.Annotation()["name"] = "tge"
					hg.Annotation()["title"] = "my title"
					return rhist.NewGraphErrorsFrom(hg).(rtests.ROOTer)
				}(),
			},
		},
		{
			Name: "TGraphAsymmErrors",
			ROOT: "retrieved: [tgae]\n",
			Want: []rtests.ROOTer{
				func() rtests.ROOTer {
					hg := hbook.NewS2D(
						hbook.Point2D{X: 1, Y: 1, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
						hbook.Point2D{X: 2, Y: 1.5, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
						hbook.Point2D{X: -1, Y: +2, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
					)
					hg.Annotation()["name"] = "tgae"
					hg.Annotation()["title"] = "my title"
					return rhist.NewGraphAsymmErrorsFrom(hg).(rtests.ROOTer)
				}(),
			},
		},
	} {
		fname := filepath.Join(dir, fmt.Sprintf("out-%d.root", i))
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Skip {
				t.Skip()
			}

			w, err := riofs.Create(fname)
			if err != nil {
				t.Fatal(err)
			}

			for i := range tc.Want {
				var (
					kname = fmt.Sprintf("key-%s-%02d", tc.Name, i)
					want  = tc.Want[i]
				)

				err = w.Put(kname, want)
				if err != nil {
					t.Fatal(err)
				}
			}

			if got, want := len(w.Keys()), len(tc.Want); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			err = w.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			r, err := riofs.Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			if got, want := len(r.Keys()), len(tc.Want); got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			for i := range tc.Want {
				var (
					kname = fmt.Sprintf("key-%s-%02d", tc.Name, i)
					want  = tc.Want[i]
				)

				rgot, err := r.Get(kname)
				if err != nil {
					t.Fatal(err)
				}

				switch rgot := rgot.(type) {
				case yodacnv.Marshaler:
					got, err := rgot.MarshalYODA()
					if err != nil {
						t.Fatalf("could not marshal 'rgot' to YODA: %+v", err)
					}
					want, err := want.(yodacnv.Marshaler).MarshalYODA()
					if err != nil {
						t.Fatalf("could not marshal 'want' to YODA: %+v", err)
					}
					if !bytes.Equal(got, want) {
						t.Fatalf("error reading back value[%d].\ngot:\n%s\nwant:\n%s", i, got, want)
					}

				default:
					if got := rgot.(rtests.ROOTer); !reflect.DeepEqual(got, want) {
						t.Fatalf("error reading back value[%d].\ngot = %#v\nwant= %#v", i, got, want)
					}
				}
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			if !rtests.HasROOT {
				t.Logf("skip test with ROOT/C++")
				return
			}

			const rootls = `#include <iostream>
#include "TFile.h"
#include "TNamed.h"

void rootls(const char *fname, const char *kname) {
	auto f = TFile::Open(fname);
	auto o = f->Get<TNamed>(kname);
	if (o == NULL) {
		std:cerr << "could not retrieve [" << kname << "]" << std::endl;
		o->ClassName();
	}
	std::cout << "retrieved: [" << o->GetName() << "]" << std::endl;
}
`
			for i := range tc.Want {
				kname := fmt.Sprintf("key-%s-%02d", tc.Name, i)

				out, err := rtests.RunCxxROOT("rootls", []byte(rootls), fname, kname)
				if err != nil {
					t.Fatalf("ROOT/C++ could not open file %q:\n%s", fname, string(out))
				}
				if got := stripLine(t, out); got != tc.ROOT {
					t.Fatalf("invalid ROOT/C++ output:\ngot:\n%s\nwant:\n%s", got, tc.ROOT)
				}
			}
		})
	}
}

func stripLine(t *testing.T, raw []byte) string {
	r := bytes.NewReader(bytes.TrimSpace(raw))
	scan := bufio.NewScanner(r)
	scan.Scan()

	var o strings.Builder
	for scan.Scan() {
		o.WriteString(scan.Text() + "\n")
	}
	if err := scan.Err(); err != nil {
		t.Fatalf("could not scan text:\n%q\nerr: %+v", raw, err)
	}

	return o.String()
}
