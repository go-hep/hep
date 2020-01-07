// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/internal/rcmd"
)

func TestDump(t *testing.T) {
	const deep = true
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "../../testdata/simple.root",
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
			name: "../../testdata/root_numpy_struct.root",
			want: `key[000]: test;1 "identical leaf names in different branches" (TTree)
[000][branch1.intleaf]: 10
[000][branch1.floatleaf]: 15.5
[000][branch2.intleaf]: 20
[000][branch2.floatleaf]: 781.2
`,
		},
		{
			name: "../../testdata/padding.root",
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new(strings.Builder)
			err := rcmd.Dump(got, tc.name, deep, nil)
			if err != nil {
				t.Fatalf("could not run root-dump: %+v", err)
			}

			if got, want := got.String(), tc.want; got != want {
				t.Fatalf("invalid root-dump output:\ngot:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}
