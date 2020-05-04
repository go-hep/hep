// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func TestFormula(t *testing.T) {
	for _, tc := range []struct {
		fname   string
		tname   string
		rvars   int
		expr    string
		imports []string
		want    []interface{}
		err     error
	}{
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "one",
			want:  []interface{}{int32(1), int32(2)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "float64(one) + float64(two*100)",
			want:  []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: 0,
			expr:  "float64(one) + float64(two*100)",
			want:  []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: 1,
			expr:  "float64(one) + float64(two*100)",
			want:  []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "one*one",
			want:  []interface{}{int32(1), int32(4)},
		},
		{
			fname:   "../testdata/simple.root",
			tname:   "tree",
			rvars:   -1,
			expr:    "math.Sqrt(float64(one*one))",
			imports: []string{"math"},
			want:    []interface{}{float64(1), float64(2)},
		},
		{
			fname:   "../testdata/simple.root",
			tname:   "tree",
			rvars:   -1,
			expr:    `fmt.Sprintf("%d", one)`,
			imports: []string{"fmt"},
			want:    []interface{}{"1", "2"},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "ArrU64",
			want:  []interface{}{[10]uint64{}, [10]uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "ArrU64[0]",
			want:  []interface{}{uint64(0), uint64(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "D16",
			want:  []interface{}{root.Float16(0.0), root.Float16(1.0)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "D32",
			want:  []interface{}{root.Double32(0.0), root.Double32(1.0)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "ArrD32[0]",
			want:  []interface{}{root.Double32(0), root.Double32(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			expr:  "float64(D32)+float64(len(SliI64))",
			want:  []interface{}{0.0, 2.0},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "float64(one",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not parse expression: rtree: could not parse formula "float64(one": 1:12: missing ',' before newline in argument list`),
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "ones",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not analyze formula type: rtree: could not analyze formula: repl.go:1:16: undefined identifier: ones`),
		},
		// {
		// 	fname:   "../testdata/simple.root",
		// 	tname:   "tree",
		// 	rvars:   -1,
		// 	expr:    "one",
		// 	imports: []string{"golang.org/x/exp/rand"},
		// 	err:     fmt.Errorf(`rtree: could not create Formula: rtree: no known stdlib import for "go-hep.org/x/hep/groot"`),
		// },
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "one+three",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not analyze formula type: rtree: could not analyze formula: repl.go:3:20: mismatched types in binary operation + between <int32> and <string>: one + three`),
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			expr:  "math.Sqrt(float64(one))",
			err:   fmt.Errorf(`rtree: could not create Formula: rtree: could not analyze formula type: rtree: could not analyze formula: repl.go:2:16: undefined "math" in math.Sqrt <*ast.SelectorExpr>`),
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

			var rvars []ReadVar
			switch tc.rvars {
			case -1:
				rvars = NewReadVars(tree)
			case 0:
				rvars = nil
			default:
				rvars = NewReadVars(tree)[:tc.rvars]
			}

			r, err := NewReader(tree, rvars, WithRange(0, 2))
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
					t.Fatalf("could not run form-eval: %+v", e)
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

func TestFormulaFunc(t *testing.T) {
	for _, tc := range []struct {
		fname    string
		tname    string
		rvars    int
		fct      interface{}
		branches []string
		want     []interface{}
		err      error
	}{
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x int32) int32 { return x },
			branches: []string{"one"},
			want:     []interface{}{int32(1), int32(2)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			fct: func(x1 int32, x2 float32) float64 {
				return float64(x1) + float64(x2*100)
			},
			branches: []string{"one", "two"},
			want:     []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: 0,
			fct: func(x1 int32, x2 float32) float64 {
				return float64(x1) + float64(x2*100)
			},
			branches: []string{"one", "two"},
			want:     []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: 1,
			fct: func(x1 int32, x2 float32) float64 {
				return float64(x1) + float64(x2*100)
			},
			branches: []string{"one", "two"},
			want:     []interface{}{float64(111), float64(222)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			fct: func(x1 int32) int32 {
				return x1 * x1
			},
			branches: []string{"one"},
			want:     []interface{}{int32(1), int32(4)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			fct: func(x1 int32) float64 {
				return math.Sqrt(float64(x1 * x1))
			},
			branches: []string{"one"},
			want:     []interface{}{float64(1), float64(2)},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			rvars: -1,
			fct: func(x1 int32) string {
				return fmt.Sprintf("%d", x1)
			},
			branches: []string{"one"},
			want:     []interface{}{"1", "2"},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x [10]uint64) [10]uint64 {
				return x
			},
			branches: []string{"ArrU64"},
			want:     []interface{}{[10]uint64{}, [10]uint64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x [10]uint64) uint64 {
				return x[0]
			},
			branches: []string{"ArrU64"},
			want:     []interface{}{uint64(0), uint64(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x root.Float16) root.Float16 {
				return x
			},
			branches: []string{"D16"},
			want:     []interface{}{root.Float16(0.0), root.Float16(1.0)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x root.Double32) root.Double32 {
				return x
			},
			branches: []string{"D32"},
			want:     []interface{}{root.Double32(0.0), root.Double32(1.0)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x [10]root.Double32) root.Double32 {
				return x[0]
			},
			branches: []string{"ArrD32"},
			want:     []interface{}{root.Double32(0), root.Double32(1)},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func(x1 root.Double32, x2 []int64) float64 {
				return float64(x1) + float64(len(x2))
			},
			branches: []string{"D32", "SliI64"},
			want:     []interface{}{0.0, 2.0},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			rvars: -1,
			fct: func() float64 {
				return 42.0
			},
			branches: nil,
			want:     []interface{}{42.0, 42.0},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() bool { return true },
			branches: nil,
			want:     []interface{}{true, true},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() int8 { return -42 },
			branches: nil,
			want:     []interface{}{int8(-42), int8(-42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() int16 { return -42 },
			branches: nil,
			want:     []interface{}{int16(-42), int16(-42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() int32 { return -42 },
			branches: nil,
			want:     []interface{}{int32(-42), int32(-42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() int64 { return -42 },
			branches: nil,
			want:     []interface{}{int64(-42), int64(-42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() uint8 { return 42 },
			branches: nil,
			want:     []interface{}{uint8(42), uint8(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() uint16 { return 42 },
			branches: nil,
			want:     []interface{}{uint16(42), uint16(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() uint32 { return 42 },
			branches: nil,
			want:     []interface{}{uint32(42), uint32(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() uint64 { return 42 },
			branches: nil,
			want:     []interface{}{uint64(42), uint64(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() float32 { return 42 },
			branches: nil,
			want:     []interface{}{float32(42), float32(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() float64 { return 42 },
			branches: nil,
			want:     []interface{}{float64(42), float64(42)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() string { return "hello" },
			branches: nil,
			want:     []interface{}{"hello", "hello"},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func() [1]string { return [1]string{"hello"} },
			branches: nil,
			want:     []interface{}{[1]string{"hello"}, [1]string{"hello"}},
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x int32) int32 { return x },
			branches: []string{"ones"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: could not find all needed ReadVars`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x1 int32, x2 float64) float64 { return 0 },
			branches: []string{"one", "two"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: argument type 1 mismatch: func=float64, read-var[two]=float32`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      "not a func",
			branches: []string{"one", "two"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: FormulaFunc expects a func`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x1 int32, x2 float64) float64 { return 0 },
			branches: []string{"one"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: num-branches/func-arity mismatch`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x1 int32) float64 { return 0 },
			branches: []string{"one", "two"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: num-branches/func-arity mismatch`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x1 int32) (a, b float64) { return },
			branches: []string{"one"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: invalid number of return values`),
		},
	} {
		t.Run("", func(t *testing.T) {
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

			var rvars []ReadVar
			switch tc.rvars {
			case -1:
				rvars = NewReadVars(tree)
			case 0:
				rvars = nil
			default:
				rvars = NewReadVars(tree)[:tc.rvars]
			}

			r, err := NewReader(tree, rvars, WithRange(0, 2))
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			form, err := r.FormulaFunc(tc.branches, tc.fct)
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
					t.Fatalf("could not run form-eval: %+v", e)
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

var sumBenchFormula float64

func BenchmarkFormula(b *testing.B) {
	for _, tc := range []struct {
		name string
		expr string
		imps []string
	}{
		{
			name: "f0",
			expr: "42.0",
		},
		{
			name: "f1",
			expr: "F64",
		},
		{
			name: "f2",
			expr: "2*F64",
		},
		{
			name: "f3",
			expr: "math.Abs(2*F64)",
			imps: []string{"math"},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			f, err := riofs.Open("../testdata/leaves.root")
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get("tree")
			if err != nil {
				b.Fatal(err)
			}
			tree := o.(Tree)

			rvars := []ReadVar{{Name: "F64", Value: new(float64)}}

			r, err := NewReader(tree, rvars)
			if err != nil {
				b.Fatal(err)
			}

			form, err := r.Formula(tc.expr, tc.imps)
			if err != nil {
				b.Fatal(err)
			}

			b.Run("Eval", func(b *testing.B) {
				err = r.Read(func(ctx RCtx) error {
					sumBenchFormula += form.Eval().(float64)
					return nil
				})
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				err = r.scan.SeekEntry(0)
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				sumBenchFormula = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormula += form.Eval().(float64)
				}
			})

			b.Run("Func", func(b *testing.B) {
				eval := form.Func().(func() float64)

				err = r.Read(func(ctx RCtx) error {
					sumBenchFormula += eval()
					return nil
				})
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				err = r.scan.SeekEntry(0)
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				sumBenchFormula = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormula += eval()
				}
			})
		})
	}
}

var sumBenchFormulaFunc float64

func BenchmarkFormulaFunc(b *testing.B) {
	for _, tc := range []struct {
		name string
		fct  interface{}
		brs  []string
	}{
		{
			name: "f0",
			fct:  func() float64 { return 42 },
		},
		{
			name: "f1",
			fct:  func(x float64) float64 { return x },
			brs:  []string{"F64"},
		},
		{
			name: "f2",
			fct:  func(x float64) float64 { return 2 * x },
			brs:  []string{"F64"},
		},
		{
			name: "f3",
			fct:  func(x float64) float64 { return math.Abs(2 * x) },
			brs:  []string{"F64"},
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			f, err := riofs.Open("../testdata/leaves.root")
			if err != nil {
				b.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get("tree")
			if err != nil {
				b.Fatal(err)
			}
			tree := o.(Tree)

			rvars := []ReadVar{{Name: "F64", Value: new(float64)}}

			r, err := NewReader(tree, rvars)
			if err != nil {
				b.Fatal(err)
			}

			form, err := r.FormulaFunc(tc.brs, tc.fct)
			if err != nil {
				b.Fatal(err)
			}
			b.Run("Eval", func(b *testing.B) {
				err = r.Read(func(ctx RCtx) error {
					sumBenchFormulaFunc += form.Eval().(float64)
					return nil
				})
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				err = r.scan.SeekEntry(0)
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				sumBenchFormula = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormulaFunc += form.Eval().(float64)
				}
			})

			b.Run("Func", func(b *testing.B) {
				eval := form.Func().(func() float64)

				err = r.Read(func(ctx RCtx) error {
					sumBenchFormula += eval()
					return nil
				})
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				err = r.scan.SeekEntry(0)
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				sumBenchFormula = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormula += eval()
				}
			})
		})
	}
}
