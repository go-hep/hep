// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rhist_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
)

func TestRWHist(t *testing.T) {

	dir, err := os.MkdirTemp("", "groot-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for i, tc := range rhist.HistoTestCases {
		fname := filepath.Join(dir, fmt.Sprintf("histos-%d.root", i))
		t.Run(tc.Name, func(t *testing.T) {
			const kname = "my-key"

			w, err := groot.Create(fname)
			if err != nil {
				t.Fatal(err)
			}

			err = w.Put(kname, tc.Want)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(w.Keys()), 1; got != want {
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

			si := r.StreamerInfos()
			if len(si) == 0 {
				t.Fatalf("empty list of streamers")
			}

			if got, want := len(r.Keys()), 1; got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			rgot, err := r.Get(kname)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := rgot.(rtests.ROOTer), tc.Want; !reflect.DeepEqual(got, want) {
				t.Fatalf("error reading back objstring.\ngot = %#v\nwant= %#v", got, want)
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
	std::cout << "retrieved object: [" << o->GetName() << "]" << std::endl;
}
`
			out, err := rtests.RunCxxROOT("rootls", []byte(rootls), fname, kname)
			if err != nil {
				t.Fatalf("ROOT/C++ could not open file %q:\n%s", fname, string(out))
			}
		})
	}
}

func TestROOT4Hist(t *testing.T) {
	f, err := groot.Open("https://github.com/scikit-hep/uproot/raw/master/tests/samples/from-geant4.root")
	if err != nil {
		t.Fatalf("could not open uproot geant4 test file: %+v", err)
	}
	defer f.Close()

	obj, err := f.Get("edep_inner")
	if err != nil {
		t.Fatalf("%+v", err)
	}

	h := obj.(*rhist.H1D)
	if got, want := h.Name(), "edep_inner"; got != want {
		t.Fatalf("invalid H1D name: got=%q, want=%q", got, want)
	}
}
