// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func TestFormula(t *testing.T) {
	for _, tc := range []struct {
		fname   string
		tname   string
		expr    string
		imports []string
		want    []interface{}
		err     error
	}{
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			expr:  "one",
			want:  []interface{}{int32(1), int32(2)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			expr:  "one*one",
			want:  []interface{}{int32(1), int32(4)},
		},
		{
			fname:   "../testdata/simple.root",
			tname:   "tree",
			expr:    "math.Sqrt(float64(one*one))",
			imports: []string{"math"},
			want:    []interface{}{float64(1), float64(2)},
		},
		{
			fname:   "../testdata/simple.root",
			tname:   "tree",
			expr:    `fmt.Sprintf("%d", one)`,
			imports: []string{"fmt"},
			want:    []interface{}{"1", "2"},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			expr:  "ArrU64",
			want:  []interface{}{[10]uint64{}, [10]uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			expr:  "ArrU64[0]",
			want:  []interface{}{uint64(0), uint64(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			expr:  "D32",
			want:  []interface{}{root.Double32(0), root.Double32(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			expr:  "float64(D32)+float64(len(SliI64))",
			want:  []interface{}{0.0, 2.0},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			expr:  "ones",
			err:   fmt.Errorf("rtree: could not create Formula: rtree: could not define formula eval-func: 6:19: undefined: ones"),
		},
		{
			fname:   "../testdata/simple.root",
			tname:   "tree",
			expr:    "one",
			imports: []string{"go-hep.org/x/hep/groot"},
			err:     fmt.Errorf(`rtree: could not create Formula: rtree: no known stdlib import for "go-hep.org/x/hep/groot"`),
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			expr:  "one+three",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not define formula eval-func: 6:19: mismatched types .int32 and .string`),
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			expr:  "math.Sqrt(float64(one))",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not define formula eval-func: 6:19: undefined: math`),
		},
	} {
		t.Run(tc.expr, func(t *testing.T) {
			f, err := riofs.Open(tc.fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := riofs.Dir(f).Get(tc.tname)
			if err != nil {
				t.Fatal(err)
			}

			tree := o.(Tree)

			r, err := NewReader(tree, NewReadVars(tree), WithRange(0, 2))
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			form, err := r.Formula(tc.expr, tc.imports)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error.\ngot= %v\nwant=%v", got, want)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %v (got=nil)", tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}

			defer func() {
				e := recover()
				if e != nil {
					t.Fatalf("could not run form-eval:\n%s\n%+v", form.prog, e)
				}
			}()

			err = r.Read(func(ctx RCtx) error {
				got := form.Eval()
				if got, want := got, tc.want[ctx.Entry]; !reflect.DeepEqual(got, want) {
					return fmt.Errorf("entry[%d]: invalid form-eval:\ngot=%v (%T)\nwant=%v (%T)", ctx.Entry, got, got, want, want)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("error: %+v", err)
			}
		})
	}
}
