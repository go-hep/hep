// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootcnv

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/internal/diff"
)

func TestRW(t *testing.T) {
	dir, err := os.MkdirTemp("", "hepmc-rootcnv-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range []string{
		"../testdata/small.hepmc",
		"../testdata/test.hepmc",
	} {
		t.Run(tc, func(t *testing.T) {
			raw, err := os.ReadFile(tc)
			if err != nil {
				t.Fatal(err)
			}

			fname := filepath.Join(dir, filepath.Base(tc)+".root")
			o, err := groot.Create(fname)
			if err != nil {
				t.Fatalf("could not create output ROOT file: %+v", err)
			}
			defer o.Close()

			w, err := NewFlatTreeWriter(o, "tree", rtree.WithTitle("HepMC ROOT tree"))
			if err != nil {
				t.Fatalf("could not create ROOT tree writer: %+v", err)
			}

			_, err = hepmc.Copy(w, hepmc.NewASCIIReader(bytes.NewReader(raw)))
			if err != nil {
				t.Fatalf("could not copy hepmc event to ROOT: %+v", err)
			}

			err = w.Close()
			if err != nil {
				t.Fatalf("could not close ROOT tree writer: %+v", err)
			}

			err = o.Close()
			if err != nil {
				t.Fatalf("could not close ROOT file: %+v", err)
			}

			f, err := groot.Open(fname)
			if err != nil {
				t.Fatalf("could not open ROOT file: %+v", err)
			}
			defer f.Close()

			r, err := riofs.Get[rtree.Tree](f, "tree")
			if err != nil {
				t.Fatalf("could not retrieve ROOT tree: %+v", err)
			}

			rr, err := NewFlatTreeReader(r)
			if err != nil {
				t.Fatalf("could not create ROOT tree reader: %+v", err)
			}
			defer rr.Close()

			buf := new(bytes.Buffer)

			ww := hepmc.NewASCIIWriter(buf)
			_, err = hepmc.Copy(ww, rr)
			if err != nil {
				t.Fatalf("could not copy hepmc event from ROOT: %+v", err)
			}

			err = ww.Close()
			if err != nil {
				t.Fatalf("could not close hepmc writer: %+v", err)
			}

			if got, want := buf.String(), string(raw); got != want {
				d := diff.Format(string(got), string(want))
				t.Fatalf("invalid r/w round trip:\n%s", d)
			}
		})
	}
}
