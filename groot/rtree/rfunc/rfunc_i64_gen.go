// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToI64 implements rfunc.Formula
type FuncToI64 struct {
	fct func() int64
}

// NewFuncToI64 return a new formula, from the provided function.
func NewFuncToI64(rvars []string, fct func() int64) *FuncToI64 {
	return &FuncToI64{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToI64) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToI64) Bind(args []any) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToI64) Func() any {
	return func() int64 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToI64)(nil)
)

// FuncI64ToI64 implements rfunc.Formula
type FuncI64ToI64 struct {
	rvars []string
	arg0  *int64
	fct   func(arg00 int64) int64
}

// NewFuncI64ToI64 return a new formula, from the provided function.
func NewFuncI64ToI64(rvars []string, fct func(arg00 int64) int64) *FuncI64ToI64 {
	return &FuncI64ToI64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI64ToI64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI64ToI64) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI64ToI64) Func() any {
	return func() int64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncI64ToI64)(nil)
)

// FuncI64I64ToI64 implements rfunc.Formula
type FuncI64I64ToI64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	fct   func(arg00 int64, arg01 int64) int64
}

// NewFuncI64I64ToI64 return a new formula, from the provided function.
func NewFuncI64I64ToI64(rvars []string, fct func(arg00 int64, arg01 int64) int64) *FuncI64I64ToI64 {
	return &FuncI64I64ToI64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI64I64ToI64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI64I64ToI64) Bind(args []any) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI64I64ToI64) Func() any {
	return func() int64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncI64I64ToI64)(nil)
)

// FuncI64I64I64ToI64 implements rfunc.Formula
type FuncI64I64I64ToI64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64) int64
}

// NewFuncI64I64I64ToI64 return a new formula, from the provided function.
func NewFuncI64I64I64ToI64(rvars []string, fct func(arg00 int64, arg01 int64, arg02 int64) int64) *FuncI64I64I64ToI64 {
	return &FuncI64I64I64ToI64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI64I64I64ToI64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI64I64I64ToI64) Bind(args []any) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI64I64I64ToI64) Func() any {
	return func() int64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncI64I64I64ToI64)(nil)
)
