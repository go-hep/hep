// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerate(t *testing.T) {
	newCtx := func(file, tree string) *Context {
		const (
			gen     = true
			verbose = false
		)
		return newContext("event", file, tree, gen, verbose)
	}

	for _, tc := range []struct {
		ctx  *Context
		want string
	}{
		{
			ctx:  newCtx("../../testdata/simple.root", "tree"),
			want: "testdata/simple.root.txt",
		},
		{
			ctx:  newCtx("../../testdata/small-flat-tree.root", "tree"),
			want: "testdata/small-flat-tree.root.txt",
		},
		{
			ctx:  newCtx("../../testdata/small-evnt-tree-fullsplit.root", "tree"),
			want: "testdata/small-evnt-tree.root.txt",
		},
		{
			ctx:  newCtx("../../testdata/small-evnt-tree-nosplit.root", "tree"),
			want: "testdata/small-evnt-tree.root.txt",
		},
		{
			ctx:  newCtx("../../testdata/x-flat-tree.root", "tree"),
			want: "testdata/x-flat-tree.root.txt",
		},
		{
			ctx:  newCtx("../../testdata/leaves.root", "tree"),
			want: "testdata/leaves.root.txt",
		},
	} {
		t.Run(tc.ctx.File, func(t *testing.T) {
			o := new(strings.Builder)
			err := process(o, tc.ctx)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read golden file %q: %+v", tc.want, err)
			}

			if got, want := o.String(), string(want); got != want {
				diff := cmp.Diff(got, want)
				t.Fatalf("invalid code generation:\n%s", diff)
			}
		})
	}
}
