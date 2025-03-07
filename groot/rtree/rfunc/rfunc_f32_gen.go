// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToF32 implements rfunc.Formula
type FuncToF32 struct {
	fct func() float32
}

// NewFuncToF32 return a new formula, from the provided function.
func NewFuncToF32(rvars []string, fct func() float32) *FuncToF32 {
	return &FuncToF32{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToF32) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToF32) Bind(args []any) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToF32) Func() any {
	return func() float32 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToF32)(nil)
)

// FuncF32ToF32 implements rfunc.Formula
type FuncF32ToF32 struct {
	rvars []string
	arg0  *float32
	fct   func(arg00 float32) float32
}

// NewFuncF32ToF32 return a new formula, from the provided function.
func NewFuncF32ToF32(rvars []string, fct func(arg00 float32) float32) *FuncF32ToF32 {
	return &FuncF32ToF32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32ToF32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32ToF32) Bind(args []any) error {
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
func (f *FuncF32ToF32) Func() any {
	return func() float32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF32ToF32)(nil)
)

// FuncF32F32ToF32 implements rfunc.Formula
type FuncF32F32ToF32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	fct   func(arg00 float32, arg01 float32) float32
}

// NewFuncF32F32ToF32 return a new formula, from the provided function.
func NewFuncF32F32ToF32(rvars []string, fct func(arg00 float32, arg01 float32) float32) *FuncF32F32ToF32 {
	return &FuncF32F32ToF32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32F32ToF32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32F32ToF32) Bind(args []any) error {
	if got, want := len(args), 2; got != want {
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
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF32F32ToF32) Func() any {
	return func() float32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncF32F32ToF32)(nil)
)

// FuncF32F32F32ToF32 implements rfunc.Formula
type FuncF32F32F32ToF32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32) float32
}

// NewFuncF32F32F32ToF32 return a new formula, from the provided function.
func NewFuncF32F32F32ToF32(rvars []string, fct func(arg00 float32, arg01 float32, arg02 float32) float32) *FuncF32F32F32ToF32 {
	return &FuncF32F32F32ToF32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32F32F32ToF32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32F32F32ToF32) Bind(args []any) error {
	if got, want := len(args), 3; got != want {
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
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncF32F32F32ToF32) Func() any {
	return func() float32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncF32F32F32ToF32)(nil)
)
