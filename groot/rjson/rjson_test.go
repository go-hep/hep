// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rjson_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/rjson"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/internal/diff"
)

func TestMarshal(t *testing.T) {
	for _, tc := range []struct {
		name string
		gen  func() root.Object
	}{
		{
			name: "h1d",
			gen: func() root.Object {
				h := hbook.NewH1D(100, 0, 100)
				h.Fill(1, 1)
				h.Fill(-1, 1)
				h.Fill(200, 1)
				h.Ann["name"] = "h1"
				h.Ann["title"] = "my title"

				return rhist.NewH1DFrom(h)
			},
		},
		{
			name: "h2d",
			gen: func() root.Object {
				h := hbook.NewH2D(5, 0, 5, 2, 0, 2)
				h.Fill(1, 1, 1)
				h.Fill(-1, -1, 1)
				h.Fill(200, 300, 1)
				h.Ann["name"] = "h2"
				h.Ann["title"] = "my title"

				return rhist.NewH2DFrom(h)
			},
		},
		{
			name: "graph",
			gen: func() root.Object {
				s := hbook.NewS2DFrom([]float64{1, 2, 3}, []float64{2, 4, 6})
				s.Annotation()["name"] = "s2"
				s.Annotation()["title"] = "my title"
				return rhist.NewGraphFrom(s)
			},
		},
		{
			name: "tge",
			gen: func() root.Object {
				s := hbook.NewS2D([]hbook.Point2D{
					{X: 1, Y: 2, ErrX: hbook.Range{Min: 10, Max: 20}, ErrY: hbook.Range{Min: 11, Max: 22}},
					{X: 2, Y: 4, ErrX: hbook.Range{Min: 20, Max: 30}, ErrY: hbook.Range{Min: 12, Max: 23}},
					{X: 3, Y: 6, ErrX: hbook.Range{Min: 30, Max: 40}, ErrY: hbook.Range{Min: 13, Max: 24}},
				}...)
				s.Annotation()["name"] = "s2"
				s.Annotation()["title"] = "my title"
				return rhist.NewGraphErrorsFrom(s)
			},
		},
		{
			name: "tgae",
			gen: func() root.Object {
				s := hbook.NewS2D([]hbook.Point2D{
					{X: 1, Y: 2, ErrX: hbook.Range{Min: 10, Max: 20}, ErrY: hbook.Range{Min: 11, Max: 22}},
					{X: 2, Y: 4, ErrX: hbook.Range{Min: 20, Max: 30}, ErrY: hbook.Range{Min: 12, Max: 23}},
					{X: 3, Y: 6, ErrX: hbook.Range{Min: 30, Max: 40}, ErrY: hbook.Range{Min: 13, Max: 24}},
				}...)
				s.Annotation()["name"] = "s2"
				s.Annotation()["title"] = "my title"
				return rhist.NewGraphAsymmErrorsFrom(s)
			},
		},
		{
			name: "tgme",
			gen: func() root.Object {
				s := hbook.NewS2D([]hbook.Point2D{
					{X: 1, Y: 2, ErrX: hbook.Range{Min: 10, Max: 20}, ErrY: hbook.Range{Min: 11, Max: 22}},
					{X: 2, Y: 4, ErrX: hbook.Range{Min: 20, Max: 30}, ErrY: hbook.Range{Min: 12, Max: 23}},
					{X: 3, Y: 6, ErrX: hbook.Range{Min: 30, Max: 40}, ErrY: hbook.Range{Min: 13, Max: 24}},
				}...)
				s.Annotation()["name"] = "s2"
				s.Annotation()["title"] = "my title"
				return rhist.NewGraphMultiErrorsFrom(s)
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				fname = filepath.Join("testdata", tc.name+".json")
				want  = filepath.Join("testdata", tc.name+"_golden.json")
			)

			obj := tc.gen()
			got, err := rjson.Marshal(obj)
			if err != nil {
				t.Fatalf("could not generate JSON: %+v", err)
			}

			err = os.WriteFile(fname, got, 0644)
			if err != nil {
				t.Fatalf("could not write JSON: %+v", err)
			}

			err = diff.Files(fname, want)
			if err != nil {
				t.Fatalf("invalid JSON:\n%v", err)
			}

			if !rtests.HasROOT {
				return
			}

			// make sure ROOT can read back that file as well.
			code := fmt.Sprintf(`#include <iostream>
#include <fstream>
#include <sstream>
#include <string>

#include "TBufferJSON.h"
#include "%[1]s.h"

void unmarshal(const char *fname) {
	std:ifstream input(fname);
	std::stringstream s;
	while (input >> s.rdbuf()) {}

	auto str = s.str();

	%[1]s *o = nullptr;
	TBufferJSON::FromJSON(o, str.c_str());

	if (o->ClassName() != std::string("%[1]s")) {
		std::cerr << "invalid class name: got=\"" << o->ClassName() << "\", want=\"%[1]s\"\n";
		exit(1);
	}
	o->Print();
}
`, obj.Class(),
			)
			out, err := rtests.RunCxxROOT("unmarshal", []byte(code), fname)
			if err != nil {
				t.Fatalf("could not run C++ ROOT: %+v\noutput:\n%s\ncode:\n%s", err, out, code)
			}

			defer os.Remove(fname)
		})
	}
}
