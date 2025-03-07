// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncI32ToF64 implements rfunc.Formula
type FuncI32ToF64 struct {
	rvars []string
	arg0  *int32
	fct   func(arg00 int32) float64
}

// NewFuncI32ToF64 return a new formula, from the provided function.
func NewFuncI32ToF64(rvars []string, fct func(arg00 int32) float64) *FuncI32ToF64 {
	return &FuncI32ToF64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI32ToF64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI32ToF64) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI32ToF64) Func() any {
	return func() float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncI32ToF64)(nil)
)

// FuncF32ToF64 implements rfunc.Formula
type FuncF32ToF64 struct {
	rvars []string
	arg0  *float32
	fct   func(arg00 float32) float64
}

// NewFuncF32ToF64 return a new formula, from the provided function.
func NewFuncF32ToF64(rvars []string, fct func(arg00 float32) float64) *FuncF32ToF64 {
	return &FuncF32ToF64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32ToF64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32ToF64) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF32ToF64) Func() any {
	return func() float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF32ToF64)(nil)
)

// FuncToF64 implements rfunc.Formula
type FuncToF64 struct {
	fct func() float64
}

// NewFuncToF64 return a new formula, from the provided function.
func NewFuncToF64(rvars []string, fct func() float64) *FuncToF64 {
	return &FuncToF64{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToF64) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToF64) Bind(args []any) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToF64) Func() any {
	return func() float64 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToF64)(nil)
)

// FuncF64ToF64 implements rfunc.Formula
type FuncF64ToF64 struct {
	rvars []string
	arg0  *float64
	fct   func(arg00 float64) float64
}

// NewFuncF64ToF64 return a new formula, from the provided function.
func NewFuncF64ToF64(rvars []string, fct func(arg00 float64) float64) *FuncF64ToF64 {
	return &FuncF64ToF64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64ToF64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64ToF64) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF64ToF64) Func() any {
	return func() float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF64ToF64)(nil)
)

// FuncF64F64ToF64 implements rfunc.Formula
type FuncF64F64ToF64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	fct   func(arg00 float64, arg01 float64) float64
}

// NewFuncF64F64ToF64 return a new formula, from the provided function.
func NewFuncF64F64ToF64(rvars []string, fct func(arg00 float64, arg01 float64) float64) *FuncF64F64ToF64 {
	return &FuncF64F64ToF64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64F64ToF64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64F64ToF64) Bind(args []any) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF64F64ToF64) Func() any {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncF64F64ToF64)(nil)
)

// FuncF64F64F64ToF64 implements rfunc.Formula
type FuncF64F64F64ToF64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64) float64
}

// NewFuncF64F64F64ToF64 return a new formula, from the provided function.
func NewFuncF64F64F64ToF64(rvars []string, fct func(arg00 float64, arg01 float64, arg02 float64) float64) *FuncF64F64F64ToF64 {
	return &FuncF64F64F64ToF64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64F64F64ToF64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64F64F64ToF64) Bind(args []any) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF64F64F64ToF64) Func() any {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncF64F64F64ToF64)(nil)
)
