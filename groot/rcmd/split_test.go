// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func TestSplit(t *testing.T) {
	tmp, err := ioutil.TempDir("", "groot-root-split-")
	if err != nil {
		t.Fatalf("%+v", err)
	}
	defer os.RemoveAll(tmp)

	type funcT func(t *testing.T, fname string) error
	for _, tc := range []struct {
		name    string
		n       int64
		input   funcT
		outputs []funcT
	}{
		{
			name:  "flat-tree-1",
			n:     10,
			input: makeSplitFlatTree(0, 10),
			outputs: []funcT{
				makeSplitFlatTree(0, 10),
			},
		},
		{
			name:  "flat-tree-2",
			n:     5,
			input: makeSplitFlatTree(0, 10),
			outputs: []funcT{
				makeSplitFlatTree(0, 5),
				makeSplitFlatTree(5, 10),
			},
		},
		{
			name:  "flat-tree-3",
			n:     11,
			input: makeSplitFlatTree(0, 10),
			outputs: []funcT{
				makeSplitFlatTree(0, 10),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				fname   = filepath.Join(tmp, tc.name+".in.root")
				oname   = filepath.Join(tmp, tc.name+".out.root")
				deep    = true
				verbose = true
			)

			err := tc.input(t, fname)
			if err != nil {
				t.Fatalf("%+v", err)
			}

			fnames, err := rcmd.Split(oname, fname, "dir-1/dir-11/mytree", tc.n, verbose)
			if err != nil {
				t.Fatalf("could not run root-merge: %+v", err)
			}
			if got, want := len(fnames), len(tc.outputs); got != want {
				t.Fatalf("invalid number of split files: got=%d, want=%d", got, want)
			}

			for i, wantFunc := range tc.outputs {
				oname := filepath.Join(tmp, fmt.Sprintf(tc.name+".out-%d.root", i))
				got := new(bytes.Buffer)
				err = rcmd.Dump(got, oname, deep, nil)
				if err != nil {
					t.Fatalf("could not run root-dump: %+v", err)
				}

				refname := filepath.Join(tmp, fmt.Sprintf(tc.name+"-%d.want.root", i))
				err = wantFunc(t, refname)
				if err != nil {
					t.Fatalf("%+v", err)
				}
				want := new(bytes.Buffer)
				err = rcmd.Dump(want, refname, deep, nil)
				if err != nil {
					t.Fatalf("could not run root-dump: %+v", err)
				}

				if got, want := got.String(), want.String(); got != want {
					t.Fatalf("invalid root-merge output:\ngot:\n%swant:\n%s", got, want)
				}
			}
		})
	}
}

func makeSplitFlatTree(beg, end int) func(t *testing.T, fname string) error {
	return func(t *testing.T, fname string) error {
		type Data struct {
			I32    int32
			F64    float64
			Str    string
			ArrF64 [5]float64
			N      int32
			SliF64 []float64 `groot:"SliF64[N]"`
		}

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

		for i := beg; i < end; i++ {
			evt.I32 = int32(i)
			evt.F64 = float64(i)
			evt.Str = fmt.Sprintf("evt-%0d", i)
			evt.ArrF64 = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			j := i % 5
			evt.N = int32(j)
			evt.SliF64 = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:j]
			_, err = tree.Write()
			if err != nil {
				t.Fatalf("could not write event %d: %+v", i, err)
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
