// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestGenerate(t *testing.T) {
	newCtx := func(file, tree string) Context {
		return Context{
			Package: "event",
			Defs: map[string]*StructDef{
				"DataReader": {
					Name:   "DataReader",
					Fields: nil,
				},
			},
			GenDataReader: true,
			File:          file,
			Tree:          tree,
			Verbose:       false,
		}
	}

	for _, tc := range []struct {
		ctx  Context
		want string
	}{
		{
			ctx:  newCtx("../../testdata/simple.root", "tree"),
			want: "testdata/simple.root.txt",
		},
	} {
		t.Run(tc.ctx.File, func(t *testing.T) {
			o := new(strings.Builder)
			err := process(o, tc.ctx)
			if err != nil {
				t.Fatalf("error: %+v", err)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read golden file %q: %+v", tc.want, err)
			}

			if got, want := o.String(), string(want); got != want {
				t.Fatalf("invalid code generation:\n=== got ===\n%v\n=== want ===\n%v\n===\n", got, want)
			}
		})
	}
}
