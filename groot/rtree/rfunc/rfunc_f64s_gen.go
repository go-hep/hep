// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncF32sToF64s implements rfunc.Formula
type FuncF32sToF64s struct {
	rvars []string
	arg0  *[]float32
	fct   func(arg00 []float32) []float64
}

// NewFuncF32sToF64s return a new formula, from the provided function.
func NewFuncF32sToF64s(rvars []string, fct func(arg00 []float32) []float64) *FuncF32sToF64s {
	return &FuncF32sToF64s{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32sToF64s) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32sToF64s) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*[]float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*[]float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF32sToF64s) Func() any {
	return func() []float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF32sToF64s)(nil)
)

// FuncF64sToF64s implements rfunc.Formula
type FuncF64sToF64s struct {
	rvars []string
	arg0  *[]float64
	fct   func(arg00 []float64) []float64
}

// NewFuncF64sToF64s return a new formula, from the provided function.
func NewFuncF64sToF64s(rvars []string, fct func(arg00 []float64) []float64) *FuncF64sToF64s {
	return &FuncF64sToF64s{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64sToF64s) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64sToF64s) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*[]float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*[]float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF64sToF64s) Func() any {
	return func() []float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF64sToF64s)(nil)
)
