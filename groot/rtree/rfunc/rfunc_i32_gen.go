// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToI32 implements rfunc.Formula
type FuncToI32 struct {
	fct func() int32
}

// NewFuncToI32 return a new formula, from the provided function.
func NewFuncToI32(rvars []string, fct func() int32) *FuncToI32 {
	return &FuncToI32{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToI32) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToI32) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToI32) Func() interface{} {
	return func() int32 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToI32)(nil)
)

// FuncI32ToI32 implements rfunc.Formula
type FuncI32ToI32 struct {
	rvars []string
	arg0  *int32
	fct   func(arg00 int32) int32
}

// NewFuncI32ToI32 return a new formula, from the provided function.
func NewFuncI32ToI32(rvars []string, fct func(arg00 int32) int32) *FuncI32ToI32 {
	return &FuncI32ToI32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI32ToI32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI32ToI32) Bind(args []interface{}) error {
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
func (f *FuncI32ToI32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncI32ToI32)(nil)
)

// FuncI32I32ToI32 implements rfunc.Formula
type FuncI32I32ToI32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	fct   func(arg00 int32, arg01 int32) int32
}

// NewFuncI32I32ToI32 return a new formula, from the provided function.
func NewFuncI32I32ToI32(rvars []string, fct func(arg00 int32, arg01 int32) int32) *FuncI32I32ToI32 {
	return &FuncI32I32ToI32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI32I32ToI32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI32I32ToI32) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
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
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI32I32ToI32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncI32I32ToI32)(nil)
)

// FuncI32I32I32ToI32 implements rfunc.Formula
type FuncI32I32I32ToI32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32) int32
}

// NewFuncI32I32I32ToI32 return a new formula, from the provided function.
func NewFuncI32I32I32ToI32(rvars []string, fct func(arg00 int32, arg01 int32, arg02 int32) int32) *FuncI32I32I32ToI32 {
	return &FuncI32I32I32ToI32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncI32I32I32ToI32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncI32I32I32ToI32) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
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
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncI32I32I32ToI32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncI32I32I32ToI32)(nil)
)
