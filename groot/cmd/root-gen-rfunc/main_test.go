// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerate(t *testing.T) {
	for _, tc := range []struct {
		test string
		pkg  string
		fct  string
		name string
		want string
		err  error
	}{
		{
			test: "math_abs",
			pkg:  "math",
			fct:  "Abs",
			name: "MyAbs",
			want: "./testdata/math_abs_golden.txt",
		},
		{
			test: "math_hypot",
			pkg:  "math",
			fct:  "Hypot",
			name: "",
			want: "./testdata/math_hypot_golden.txt",
		},
		{
			test: "func1",
			pkg:  "",
			fct:  "func(x, y float64) bool",
			name: "",
			want: "./testdata/func1_golden.txt",
		},
		{
			test: "func2",
			pkg:  "",
			fct:  "func(float64, float64) bool",
			name: "MyFunc",
			want: "./testdata/func2_golden.txt",
		},
		{
			test: "invalid-arg-pkg",
			pkg:  "",
			fct:  "",
			err:  fmt.Errorf("missing package import path and/or function name"),
		},
		{
			test: "invalid-arg-func",
			pkg:  "",
			fct:  "math.Abs",
			err:  fmt.Errorf("missing function signature"),
		},
		{
			test: "invalid-expr",
			pkg:  "",
			fct:  "func F()",
			err:  fmt.Errorf(`could not generate rfunc formula: genroot: could not create rfunc generator: genroot: could not parse function signature: genroot: could not parse "func F()": 1:6: expected '(', found F`),
		},
		{
			test: "invalid-func-name",
			pkg:  "math",
			fct:  "AbsXXX",
			err:  fmt.Errorf(`could not generate rfunc formula: genroot: could not create rfunc generator: genroot: could not find AbsXXX in package "math"`),
		},
		{
			test: "invalid-func-obj",
			pkg:  "math",
			fct:  "Pi",
			err:  fmt.Errorf(`could not generate rfunc formula: genroot: could not create rfunc generator: genroot: object Pi in package "math" is not a func (*types.Const)`),
		},
	} {
		t.Run(tc.test, func(t *testing.T) {
			tmp, err := ioutil.TempFile("", "groot-gen-rfunc-")
			if err != nil {
				t.Fatal(err)
			}
			_ = tmp.Close()
			defer os.Remove(tmp.Name())

			usage := func() {}
			err = generate(tmp.Name(), tc.pkg, tc.fct, tc.name, usage)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("could not generate: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %v", err)
			case err == nil && tc.err == nil:
				// ok.
			}

			got, err := ioutil.ReadFile(tmp.Name())
			if err != nil {
				t.Fatalf("could not read generated rfunc: %+v", err)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference rfunc: %+v", err)
			}

			if got, want := string(got), string(want); got != want {
				diff := cmp.Diff(want, got)
				t.Fatalf("invalid generated code:\n%s\n", diff)
			}
		})
	}
}
