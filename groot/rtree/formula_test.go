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
			fct:      func(v bool) bool { return v },
			branches: []string{"B"},
			want:     []interface{}{true, false},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v int8) int8 { return v },
			branches: []string{"I8"},
			want:     []interface{}{int8(0), int8(-1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v int16) int16 { return v },
			branches: []string{"I16"},
			want:     []interface{}{int16(0), int16(-1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v int32) int32 { return v },
			branches: []string{"I32"},
			want:     []interface{}{int32(0), int32(-1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v int64) int64 { return v },
			branches: []string{"I64"},
			want:     []interface{}{int64(0), int64(-1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v uint8) uint8 { return v },
			branches: []string{"U8"},
			want:     []interface{}{uint8(0), uint8(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v uint16) uint16 { return v },
			branches: []string{"U16"},
			want:     []interface{}{uint16(0), uint16(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v uint32) uint32 { return v },
			branches: []string{"U32"},
			want:     []interface{}{uint32(0), uint32(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v uint64) uint64 { return v },
			branches: []string{"U64"},
			want:     []interface{}{uint64(0), uint64(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v float32) float32 { return v },
			branches: []string{"F32"},
			want:     []interface{}{float32(0), float32(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v float64) float64 { return v },
			branches: []string{"F64"},
			want:     []interface{}{float64(0), float64(1)},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v string) string { return v },
			branches: []string{"Str"},
			want:     []interface{}{"str-0", "str-1"},
		},
		{
			fname:    "../testdata/leaves.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(v string) [1]string { return [1]string{v} },
			branches: []string{"Str"},
			want:     []interface{}{[1]string{"str-0"}, [1]string{"str-1"}},
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x int32) int32 { return x },
			branches: []string{"ones"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: could not find all needed ReadVars (missing: [ones])`),
		},
		{
			fname:    "../testdata/simple.root",
			tname:    "tree",
			rvars:    -1,
			fct:      func(x int32, y float32, z int32) int32 { return x },
			branches: []string{"one", "twos", "ones"},
			err:      fmt.Errorf(`rtree: could not create FormulaFunc: rtree: could not find all needed ReadVars (missing: [twos ones])`),
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
				if got, want := form.Eval(), tc.want[ctx.Entry]; !reflect.DeepEqual(got, want) {
					return fmt.Errorf("entry[%d]: invalid form-eval:\ngot=%v (%T)\nwant=%v (%T)", ctx.Entry, got, got, want, want)
				}

				if got, want := reflect.ValueOf(form.Func()).Call(nil)[0].Interface(), tc.want[ctx.Entry]; !reflect.DeepEqual(got, want) {
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

				r.r.reset()

				sumBenchFormulaFunc = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormulaFunc += form.Eval().(float64)
				}
			})

			b.Run("Func", func(b *testing.B) {
				eval := form.Func().(func() float64)

				err = r.Read(func(ctx RCtx) error {
					sumBenchFormulaFunc += eval()
					return nil
				})
				if err != nil {
					b.Fatalf("error: %+v", err)
				}

				r.r.reset()

				sumBenchFormulaFunc = 0
				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; i < b.N; i++ {
					sumBenchFormulaFunc += eval()
				}
			})
		})
	}
}
