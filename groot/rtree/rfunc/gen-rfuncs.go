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
			Name: "bool",
			Funcs: []genroot.RFunc{
				{Def: "func() bool"},
				// float32
				{Def: "func(x1 float32) bool"},
				{Def: "func(x1,x2 float32) bool"},
				{Def: "func(x1,x2,x3 float32) bool"},
				// float64
				{Def: "func(x1 float64) bool"},
				{Def: "func(x1,x2 float64) bool"},
				{Def: "func(x1,x2,x3 float64) bool"},
			},
		},
		{
			Name: "i32",
			Funcs: []genroot.RFunc{
				{Def: "func() int32"},
				{Def: "func(x1 int32) int32"},
				{Def: "func(x1,x2 int32) int32"},
				{Def: "func(x1,x2,x3 int32) int32"},
			},
		},
		{
			Name: "i64",
			Funcs: []genroot.RFunc{
				{Def: "func() int64"},
				{Def: "func(x1 int64) int64"},
				{Def: "func(x1,x2 int64) int64"},
				{Def: "func(x1,x2,x3 int64) int64"},
			},
		},
		{
			Name: "u32",
			Funcs: []genroot.RFunc{
				{Def: "func() uint32"},
				{Def: "func(x1 uint32) uint32"},
				{Def: "func(x1,x2 uint32) uint32"},
				{Def: "func(x1,x2,x3 uint32) uint32"},
			},
		},
		{
			Name: "u64",
			Funcs: []genroot.RFunc{
				{Def: "func() uint64"},
				{Def: "func(x1 uint64) uint64"},
				{Def: "func(x1,x2 uint64) uint64"},
				{Def: "func(x1,x2,x3 uint64) uint64"},
			},
		},
		{
			Name: "f32",
			Funcs: []genroot.RFunc{
				// float32
				{Def: "func() float32"},
				{Def: "func(x1 float32) float32"},
				{Def: "func(x1,x2 float32) float32"},
				{Def: "func(x1,x2,x3 float32) float32"},
			},
		},
		{
			Name: "f64",
			Funcs: []genroot.RFunc{
				// int32
				{Def: "func(x1 int32) float64"},
				// float32
				{Def: "func(x1 float32) float64"},
				// float64
				{Def: "func() float64"},
				{Def: "func(x1 float64) float64"},
				{Def: "func(x1,x2 float64) float64"},
				{Def: "func(x1,x2,x3 float64) float64"},
			},
		},
		// slices
		{
			Name: "f64s",
			Funcs: []genroot.RFunc{
				// float32s
				{Def: "func(xs []float32) []float64"},
				// float64s
				{Def: "func(xs []float64) []float64"},
			},
		},
	} {
		genRFunc(typ)
	}
}

func genRFunc(typ Func) {
	fcodeName := "./rfunc_" + typ.Name + "_gen.go"
	yearCode := genroot.ExtractYear(fcodeName)

	ftestName := "./rfunc_" + typ.Name + "_gen_test.go"
	yearTest := genroot.ExtractYear(ftestName)

	f, err := os.Create(fcodeName)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	ft, err := os.Create(ftestName)
	if err != nil {
		log.Fatal(err)
	}
	defer ft.Close()

	genroot.GenImports(yearCode, "rfunc", f,
		"fmt",
	)

	genroot.GenImports(yearTest, "rfunc", ft,
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
