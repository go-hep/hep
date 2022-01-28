// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go-hep.org/x/hep/groot/rcmd"
)

func TestDump(t *testing.T) {
	const deep = true
	loadRef := func(fname string) string {
		t.Helper()
		raw, err := os.ReadFile(fname)
		if err != nil {
			t.Fatalf("could not load reference file %q: %+v", fname, err)
		}
		return string(raw)
	}

	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "../testdata/simple.root",
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
			name: "../testdata/root_numpy_struct.root",
			want: `key[000]: test;1 "identical leaf names in different branches" (TTree)
[000][branch1.intleaf]: 10
[000][branch1.floatleaf]: 15.5
[000][branch2.intleaf]: 20
[000][branch2.floatleaf]: 781.2
`,
		},
		{
			name: "../testdata/tntuple.root",
			want: `key[000]: ntup;1 "my ntuple title" (TNtuple)
[000][x]: 0
[000][y]: 0.5
[001][x]: 1
[001][y]: 1.5
[002][x]: 2
[002][y]: 2.5
[003][x]: 3
[003][y]: 3.5
[004][x]: 4
[004][y]: 4.5
[005][x]: 5
[005][y]: 5.5
[006][x]: 6
[006][y]: 6.5
[007][x]: 7
[007][y]: 7.5
[008][x]: 8
[008][y]: 8.5
[009][x]: 9
[009][y]: 9.5
`,
		},
		{
			name: "../testdata/tntupled.root",
			want: `key[000]: ntup;1 "my ntuple title" (TNtupleD)
[000][x]: 0
[000][y]: 0.5
[001][x]: 1
[001][y]: 1.5
[002][x]: 2
[002][y]: 2.5
[003][x]: 3
[003][y]: 3.5
[004][x]: 4
[004][y]: 4.5
[005][x]: 5
[005][y]: 5.5
[006][x]: 6
[006][y]: 6.5
[007][x]: 7
[007][y]: 7.5
[008][x]: 8
[008][y]: 8.5
[009][x]: 9
[009][y]: 9.5
`,
		},
		{
			name: "../testdata/padding.root",
			want: `key[000]: tree;1 "tree w/ & w/o padding" (TTree)
[000][pad.x1]: 0
[000][pad.x2]: 548655054794
[000][pad.x3]: 0
[000][nop.x1]: 0
[000][nop.x2]: 0
[000][nop.x3]: 0
[001][pad.x1]: 1
[001][pad.x2]: 72058142692982730
[001][pad.x3]: 0
[001][nop.x1]: 1
[001][nop.x2]: 1
[001][nop.x3]: 1
[002][pad.x1]: 2
[002][pad.x2]: 144115736730910666
[002][pad.x3]: 0
[002][nop.x1]: 2
[002][nop.x2]: 2
[002][nop.x3]: 2
[003][pad.x1]: 3
[003][pad.x2]: 216173330768838602
[003][pad.x3]: 0
[003][nop.x1]: 3
[003][nop.x2]: 3
[003][nop.x3]: 3
[004][pad.x1]: 4
[004][pad.x2]: 288230924806766538
[004][pad.x3]: 0
[004][nop.x1]: 4
[004][nop.x2]: 4
[004][nop.x3]: 4
`,
		},
		{
			name: "../testdata/small-flat-tree.root",
			want: loadRef("testdata/small-flat-tree.root.txt"),
		},
		{
			name: "../testdata/small-evnt-tree-fullsplit.root",
			want: loadRef("testdata/small-evnt-tree-fullsplit.root.txt"),
		},
		{
			name: "../testdata/small-evnt-tree-nosplit.root",
			want: loadRef("testdata/small-evnt-tree-nosplit.root.txt"),
		},
		{
			name: "../testdata/leaves.root",
			want: loadRef("testdata/leaves.root.txt"),
		},
		{
			name: "../testdata/embedded-std-vector.root",
			want: `key[000]: modules;1 "Module Tree Analysis" (TTree)
[000][hits_n]: 10
[000][hits_time_mc]: [12.206399 11.711122 11.73492 12.45704 11.558057 11.56502 11.687759 11.528914 12.893241 11.429288]
[001][hits_n]: 11
[001][hits_time_mc]: [11.718019 12.985347 12.23121 11.825082 12.405976 15.339471 11.939051 12.935032 13.661691 11.969542 11.893113]
[002][hits_n]: 15
[002][hits_time_mc]: [12.231329 12.214683 12.194867 12.246092 11.859249 19.35934 12.155213 12.226966 -4.712372 11.851829 11.8806925 11.8204975 11.866335 13.285733 -4.6470475]
[003][hits_n]: 9
[003][hits_time_mc]: [11.33844 11.725604 12.774131 12.108594 12.192085 12.120591 12.129445 12.18349 11.591005]
[004][hits_n]: 13
[004][hits_time_mc]: [12.156414 12.641215 11.678816 12.329707 11.578169 12.512748 11.840462 14.120602 11.875188 14.133265 14.105912 14.905052 11.813884]
`,
		},
		{
			// recovered baskets
			name: "../testdata/uproot/issue21.root",
			want: loadRef("../testdata/uproot/issue21.root.txt"),
		},
		{
			name: "../testdata/treeCharExample.root",
			want: `key[000]: nominal;1 "tree" (TTree)
[000][d_fakeEvent]: false
[000][d_lep_ECIDS]: [1 110]
[001][d_fakeEvent]: false
[001][d_lep_ECIDS]: [1 1]
[002][d_fakeEvent]: false
[002][d_lep_ECIDS]: [110 110]
[003][d_fakeEvent]: false
[003][d_lep_ECIDS]: [1 110]
[004][d_fakeEvent]: false
[004][d_lep_ECIDS]: [1 110]
[005][d_fakeEvent]: false
[005][d_lep_ECIDS]: [110 110]
[006][d_fakeEvent]: false
[006][d_lep_ECIDS]: [1 0]
[007][d_fakeEvent]: false
[007][d_lep_ECIDS]: [1 110]
[008][d_fakeEvent]: false
[008][d_lep_ECIDS]: [1 1]
[009][d_fakeEvent]: false
[009][d_lep_ECIDS]: [1 1]
`,
		},
		{
			// no embedded streamer for std::string
			name: "../testdata/no-streamer-string.root",
			want: loadRef("testdata/no-streamer-string.root.txt"),
		},
		{
			// 'This' streamer of vector<vector<double>>
			name: "../testdata/vec-vec-double.root",
			want: `key[000]: t;1 "" (TTree)
[000][x]: []
[001][x]: [[] []]
[002][x]: [[10] [] [10 20]]
[003][x]: [[20 -21 -22]]
[004][x]: [[200] [-201] [202]]
`,
		},
		{
			// Geant4 w/ recover baskets
			name: "../testdata/g4-like.root",
			want: `key[000]: mytree;1 "my title" (TTree)
[000][i32]: 1
[000][f64]: 1
[000][slif64]: []
[001][i32]: 2
[001][f64]: 2
[001][slif64]: [1]
[002][i32]: 3
[002][f64]: 3
[002][slif64]: [2 3]
[003][i32]: 4
[003][f64]: 4
[003][slif64]: [3 4 5]
[004][i32]: 5
[004][f64]: 5
[004][slif64]: [4 5 6 7]
`,
		},
		{
			// std::bitset<N>, std::vector<std::bitset<N>>
			name: "../testdata/std-bitset.root",
			want: `key[000]: tree;1 "my tree title" (TTree)
[000][evt]: {[0 0 0 1 0 0 0 1] [[1 1 1 0 1 1 1 0]]}
[001][evt]: {[1 0 0 1 1 0 0 1] [[0 0 0 1 0 0 0 1] [1 1 1 0 1 1 1 0]]}
[002][evt]: {[0 1 1 0 0 1 1 0] [[1 0 0 1 1 0 0 1] [0 1 1 0 0 1 1 0] [1 1 0 0 1 1 0 0]]}
`,
		},
		{
			// std-map w/ split-level=0, mbr-wise
			name: "../testdata/std-map-split0.root",
			want: loadRef("testdata/std-map-split0.root.txt"),
		},
		{
			// std-map w/ split-level=1, mbr-wise
			name: "../testdata/std-map-split1.root",
			want: loadRef("testdata/std-map-split1.root.txt"),
		},
		{
			// n-dim arrays
			// FIXME(sbinet): arrays of Float16_t and Double32_t are flatten.
			// This is because of:
			// https://sft.its.cern.ch/jira/browse/ROOT-10149
			name: "../testdata/ndim.root",
			want: loadRef("testdata/ndim.root.txt"),
		},
		{
			// slices of n-dim arrays
			// FIXME(sbinet): arrays of Float16_t and Double32_t are flatten.
			// This is because of:
			// https://sft.its.cern.ch/jira/browse/ROOT-10149
			name: "../testdata/ndim-slice.root",
			want: loadRef("testdata/ndim-slice.root.txt"),
		},
		{
			name: "../testdata/tformula.root",
			want: `key[000]: func1;1 "[0] + [1]*x" (TF1) => "TF1{Formula: TFormula{[p0]+[p1]*x}}"
key[001]: func2;1 "func2" (TF1) => "TF1{Params: TF1Parameters{Values: [10 20], Names: [p0 p1]}}"
key[002]: func3;1 "func3" (TF1) => "TF1{Params: TF1Parameters{Values: [1 -0.3 0 1], Names: [p0 p1 p2 p3]}}"
key[003]: func4;1 "func4" (TF1) => "TF1{Params: TF1Parameters{Values: [0 0 0 0 0 0], Names: [p0 p1 p2 p3 p4 p5]}}"
key[004]: fconv;1 "" (TF1Convolution) => "TF1Convolution{Func1: TF1{Formula: TFormula{exp([Constant]+[Slope]*x)}}, Func2: TF1{Formula: TFormula{[Constant]*exp(-0.5*((x-[Mean])/[Sigma])*((x-[Mean])/[Sigma]))}}}"
key[005]: fnorm;1 "" (TF1NormSum) => "TF1Convolution{Funcs: []{TF1{Formula: TFormula{[p0]+[p1]*x}}, TF1{Params: TF1Parameters{Values: [10 20], Names: [p0 p1]}}}, Coeffs: [10 20]}"
`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new(strings.Builder)
			err := rcmd.Dump(got, tc.name, deep, nil)
			if err != nil {
				t.Fatalf("could not run root-dump: %+v", err)
			}

			if got, want := got.String(), tc.want; got != want {
				diff := cmp.Diff(want, got)
				t.Fatalf("invalid root-dump output: -- (-ref +got)\n%s", diff)
			}
		})
	}
}
