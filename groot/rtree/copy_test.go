// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/internal/rcmd"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func TestCopyTree(t *testing.T) {
	const deep = true
	tmp, err := ioutil.TempDir("", "groot-rtree-copy-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		file     string
		tree     string
		branches map[string]int
		nevts    int64
		want     string
	}{
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			nevts: 4,
			want: `key[000]: tree;1 "fake data" (TTree)
[000][one]: 1
[000][two]: 1.1
[000][three]: uno
[001][one]: 2
[001][two]: 2.2
[001][three]: dos
[002][one]: 3
[002][two]: 3.3
[002][three]: tres
[003][one]: 4
[003][two]: 4.4
[003][three]: quatro
`,
		},
		{
			file: "../testdata/simple.root",
			tree: "tree",
			branches: map[string]int{
				"one":   1,
				"two":   1,
				"three": 1,
			},
			nevts: 4,
			want: `key[000]: tree;1 "fake data" (TTree)
[000][one]: 1
[000][two]: 1.1
[000][three]: uno
[001][one]: 2
[001][two]: 2.2
[001][three]: dos
[002][one]: 3
[002][two]: 3.3
[002][three]: tres
[003][one]: 4
[003][two]: 4.4
[003][three]: quatro
`,
		},
		{
			file: "../testdata/simple.root",
			tree: "tree",
			branches: map[string]int{
				"one":   1,
				"three": 1,
			},
			nevts: 4,
			want: `key[000]: tree;1 "fake data" (TTree)
[000][one]: 1
[000][three]: uno
[001][one]: 2
[001][three]: dos
[002][one]: 3
[002][three]: tres
[003][one]: 4
[003][three]: quatro
`,
		},
		{
			file: "../testdata/simple.root",
			tree: "tree",
			branches: map[string]int{
				"one":   1,
				"two":   1,
				"three": 1,
			},
			nevts: 0,
			want:  "key[000]: tree;1 \"fake data\" (TTree)\n",
		},
		{
			file: "../testdata/simple.root",
			tree: "tree",
			branches: map[string]int{
				"one":   1,
				"two":   1,
				"three": 1,
			},
			nevts: -1,
			want:  "key[000]: tree;1 \"fake data\" (TTree)\n",
		},
		{
			file:  "../testdata/leaves.root",
			tree:  "tree",
			nevts: 4,
			want: `key[000]: tree;1 "my tree title" (TTree)
[000][B]: true
[000][Str]: str-0
[000][I8]: 0
[000][I16]: 0
[000][I32]: 0
[000][I64]: 0
[000][U8]: 0
[000][U16]: 0
[000][U32]: 0
[000][U64]: 0
[000][F32]: 0
[000][F64]: 0
[000][ArrBs]: [true false false false false false false false false false]
[000][ArrI8]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrI16]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrI32]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrI64]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrU8]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrU16]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrU32]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrU64]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrF32]: [0 0 0 0 0 0 0 0 0 0]
[000][ArrF64]: [0 0 0 0 0 0 0 0 0 0]
[000][N]: 0
[000][SliBs]: []
[000][SliI8]: []
[000][SliI16]: []
[000][SliI32]: []
[000][SliI64]: []
[000][SliU8]: []
[000][SliU16]: []
[000][SliU32]: []
[000][SliU64]: []
[000][SliF32]: []
[000][SliF64]: []
[001][B]: false
[001][Str]: str-1
[001][I8]: -1
[001][I16]: -1
[001][I32]: -1
[001][I64]: -1
[001][U8]: 1
[001][U16]: 1
[001][U32]: 1
[001][U64]: 1
[001][F32]: 1
[001][F64]: 1
[001][ArrBs]: [false true false false false false false false false false]
[001][ArrI8]: [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1]
[001][ArrI16]: [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1]
[001][ArrI32]: [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1]
[001][ArrI64]: [-1 -1 -1 -1 -1 -1 -1 -1 -1 -1]
[001][ArrU8]: [1 1 1 1 1 1 1 1 1 1]
[001][ArrU16]: [1 1 1 1 1 1 1 1 1 1]
[001][ArrU32]: [1 1 1 1 1 1 1 1 1 1]
[001][ArrU64]: [1 1 1 1 1 1 1 1 1 1]
[001][ArrF32]: [1 1 1 1 1 1 1 1 1 1]
[001][ArrF64]: [1 1 1 1 1 1 1 1 1 1]
[001][N]: 1
[001][SliBs]: [true]
[001][SliI8]: [-1]
[001][SliI16]: [-1]
[001][SliI32]: [-1]
[001][SliI64]: [-1]
[001][SliU8]: [1]
[001][SliU16]: [1]
[001][SliU32]: [1]
[001][SliU64]: [1]
[001][SliF32]: [1]
[001][SliF64]: [1]
[002][B]: true
[002][Str]: str-2
[002][I8]: -2
[002][I16]: -2
[002][I32]: -2
[002][I64]: -2
[002][U8]: 2
[002][U16]: 2
[002][U32]: 2
[002][U64]: 2
[002][F32]: 2
[002][F64]: 2
[002][ArrBs]: [false false true false false false false false false false]
[002][ArrI8]: [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2]
[002][ArrI16]: [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2]
[002][ArrI32]: [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2]
[002][ArrI64]: [-2 -2 -2 -2 -2 -2 -2 -2 -2 -2]
[002][ArrU8]: [2 2 2 2 2 2 2 2 2 2]
[002][ArrU16]: [2 2 2 2 2 2 2 2 2 2]
[002][ArrU32]: [2 2 2 2 2 2 2 2 2 2]
[002][ArrU64]: [2 2 2 2 2 2 2 2 2 2]
[002][ArrF32]: [2 2 2 2 2 2 2 2 2 2]
[002][ArrF64]: [2 2 2 2 2 2 2 2 2 2]
[002][N]: 2
[002][SliBs]: [false true]
[002][SliI8]: [-2 -2]
[002][SliI16]: [-2 -2]
[002][SliI32]: [-2 -2]
[002][SliI64]: [-2 -2]
[002][SliU8]: [2 2]
[002][SliU16]: [2 2]
[002][SliU32]: [2 2]
[002][SliU64]: [2 2]
[002][SliF32]: [2 2]
[002][SliF64]: [2 2]
[003][B]: false
[003][Str]: str-3
[003][I8]: -3
[003][I16]: -3
[003][I32]: -3
[003][I64]: -3
[003][U8]: 3
[003][U16]: 3
[003][U32]: 3
[003][U64]: 3
[003][F32]: 3
[003][F64]: 3
[003][ArrBs]: [false false false true false false false false false false]
[003][ArrI8]: [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3]
[003][ArrI16]: [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3]
[003][ArrI32]: [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3]
[003][ArrI64]: [-3 -3 -3 -3 -3 -3 -3 -3 -3 -3]
[003][ArrU8]: [3 3 3 3 3 3 3 3 3 3]
[003][ArrU16]: [3 3 3 3 3 3 3 3 3 3]
[003][ArrU32]: [3 3 3 3 3 3 3 3 3 3]
[003][ArrU64]: [3 3 3 3 3 3 3 3 3 3]
[003][ArrF32]: [3 3 3 3 3 3 3 3 3 3]
[003][ArrF64]: [3 3 3 3 3 3 3 3 3 3]
[003][N]: 3
[003][SliBs]: [false false true]
[003][SliI8]: [-3 -3 -3]
[003][SliI16]: [-3 -3 -3]
[003][SliI32]: [-3 -3 -3]
[003][SliI64]: [-3 -3 -3]
[003][SliU8]: [3 3 3]
[003][SliU16]: [3 3 3]
[003][SliU32]: [3 3 3]
[003][SliU64]: [3 3 3]
[003][SliF32]: [3 3 3]
[003][SliF64]: [3 3 3]
`,
		},
	} {
		t.Run(tc.file, func(t *testing.T) {
			f, err := groot.Open(tc.file)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer f.Close()

			obj, err := riofs.Dir(f).Get(tc.tree)
			if err != nil {
				t.Fatalf("could not get input tree: %+v", err)
			}

			src := obj.(rtree.Tree)

			wvars := rtree.WriteVarsFromTree(src)
			if len(tc.branches) > 0 {
				all := wvars
				wvars = make([]rtree.WriteVar, 0, len(tc.branches))
				for _, wvar := range all {
					if _, ok := tc.branches[wvar.Name]; !ok {
						continue
					}
					wvars = append(wvars, wvar)
				}
			}

			oname := filepath.Join(tmp, fmt.Sprintf("copy-%s.root", filepath.Base(tc.file)))
			o, err := groot.Create(oname)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer o.Close()

			dst, err := rtree.NewWriter(o, src.Name(), wvars, rtree.WithTitle(src.Title()))
			if err != nil {
				t.Fatalf("could not create tree writer: %+v", err)
			}

			_, err = rtree.CopyN(dst, src, tc.nevts)
			if err != nil {
				t.Fatalf("could not copy tree: %+v", err)
			}

			err = dst.Close()
			if err != nil {
				t.Fatalf("could not close tree: %+v", err)
			}

			err = o.Close()
			if err != nil {
				t.Fatalf("could not close file: %+v", err)
			}

			got := new(bytes.Buffer)
			err = rcmd.Dump(got, oname, deep, nil)
			if err != nil {
				t.Errorf("could not dump output file: %+v", err)
			}

			if got, want := got.String(), tc.want; got != want {
				t.Fatalf("invalid root-dump output:\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}
