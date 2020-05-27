// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genRFuncs()
}

type Func struct {
	Name  string
	Funcs []genroot.RFunc
}

func genRFuncs() {
	for _, typ := range []Func{
		{
			Name: "bool_f64",
			Funcs: []genroot.RFunc{
				{Name: "PredF64Ar0", Def: "func() bool"},
				{Name: "PredF64Ar1", Def: "func(x1 float64) bool"},
				{Name: "PredF64Ar2", Def: "func(x1,x2 float64) bool"},
				{Name: "PredF64Ar3", Def: "func(x1,x2,x3 float64) bool"},
			},
		},
		{
			Name: "bool_f32",
			Funcs: []genroot.RFunc{
				{Name: "PredF32Ar0", Def: "func() bool"},
				{Name: "PredF32Ar1", Def: "func(x1 float32) bool"},
				{Name: "PredF32Ar2", Def: "func(x1,x2 float32) bool"},
				{Name: "PredF32Ar3", Def: "func(x1,x2,x3 float32) bool"},
			},
		},
		{
			Name: "f64",
			Funcs: []genroot.RFunc{
				{Name: "F64Ar0", Def: "func() float64"},
				{Name: "F64Ar1", Def: "func(x1 float64) float64"},
				{Name: "F64Ar2", Def: "func(x1,x2 float64) float64"},
				{Name: "F64Ar3", Def: "func(x1,x2,x3 float64) float64"},
			},
		},
		{
			Name: "f32",
			Funcs: []genroot.RFunc{
				{Name: "F32Ar0", Def: "func() float32"},
				{Name: "F32Ar1", Def: "func(x1 float32) float32"},
				{Name: "F32Ar2", Def: "func(x1,x2 float32) float32"},
				{Name: "F32Ar3", Def: "func(x1,x2,x3 float32) float32"},
			},
		},
		{
			Name: "i32",
			Funcs: []genroot.RFunc{
				{Name: "I32Ar0", Def: "func() int32"},
				{Name: "I32Ar1", Def: "func(x1 int32) int32"},
				{Name: "I32Ar2", Def: "func(x1,x2 int32) int32"},
				{Name: "I32Ar3", Def: "func(x1,x2,x3 int32) int32"},
			},
		},
		{
			Name: "i64",
			Funcs: []genroot.RFunc{
				{Name: "I64Ar0", Def: "func() int64"},
				{Name: "I64Ar1", Def: "func(x1 int64) int64"},
				{Name: "I64Ar2", Def: "func(x1,x2 int64) int64"},
				{Name: "I64Ar3", Def: "func(x1,x2,x3 int64) int64"},
			},
		},
		{
			Name: "u32",
			Funcs: []genroot.RFunc{
				{Name: "U32Ar0", Def: "func() uint32"},
				{Name: "U32Ar1", Def: "func(x1 uint32) uint32"},
				{Name: "U32Ar2", Def: "func(x1,x2 uint32) uint32"},
				{Name: "U32Ar3", Def: "func(x1,x2,x3 uint32) uint32"},
			},
		},
		{
			Name: "u64",
			Funcs: []genroot.RFunc{
				{Name: "U64Ar0", Def: "func() uint64"},
				{Name: "U64Ar1", Def: "func(x1 uint64) uint64"},
				{Name: "U64Ar2", Def: "func(x1,x2 uint64) uint64"},
				{Name: "U64Ar3", Def: "func(x1,x2,x3 uint64) uint64"},
			},
		},
	} {
		genRFunc(typ)
	}
}

func genRFunc(typ Func) {
	f, err := os.Create("./rfunc_" + typ.Name + "_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ft, err := os.Create("./rfunc_" + typ.Name + "_gen_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer ft.Close()

	genroot.GenImports("rfunc", f,
		"fmt",
	)

	genroot.GenImports("rfunc", ft,
		"reflect",
		"testing",
	)

	for i := range typ.Funcs {
		if i > 0 {
			fmt.Fprintf(f, "\n")
			fmt.Fprintf(ft, "\n")
		}
		fct := typ.Funcs[i]
		fct.Pkg = "go-hep.org/x/hep/groot/rtree/rfunc"
		gen, err := genroot.NewRFuncGenerator(f, fct)
		if err != nil {
			log.Fatalf("could not create generator: %+v", err)
		}
		err = gen.Generate()
		if err != nil {
			log.Fatalf("could not generate code for %q: %v\n", fct.Def, err)
		}
		err = gen.GenerateTest(ft)
		if err != nil {
			log.Fatalf("could not generate test for %q: %v\n", fct.Def, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)

	err = ft.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(ft)
}
