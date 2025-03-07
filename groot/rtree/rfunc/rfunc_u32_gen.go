// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToU32 implements rfunc.Formula
type FuncToU32 struct {
	fct func() uint32
}

// NewFuncToU32 return a new formula, from the provided function.
func NewFuncToU32(rvars []string, fct func() uint32) *FuncToU32 {
	return &FuncToU32{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToU32) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToU32) Bind(args []any) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToU32) Func() any {
	return func() uint32 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToU32)(nil)
)

// FuncU32ToU32 implements rfunc.Formula
type FuncU32ToU32 struct {
	rvars []string
	arg0  *uint32
	fct   func(arg00 uint32) uint32
}

// NewFuncU32ToU32 return a new formula, from the provided function.
func NewFuncU32ToU32(rvars []string, fct func(arg00 uint32) uint32) *FuncU32ToU32 {
	return &FuncU32ToU32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU32ToU32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU32ToU32) Bind(args []any) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU32ToU32) Func() any {
	return func() uint32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncU32ToU32)(nil)
)

// FuncU32U32ToU32 implements rfunc.Formula
type FuncU32U32ToU32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	fct   func(arg00 uint32, arg01 uint32) uint32
}

// NewFuncU32U32ToU32 return a new formula, from the provided function.
func NewFuncU32U32ToU32(rvars []string, fct func(arg00 uint32, arg01 uint32) uint32) *FuncU32U32ToU32 {
	return &FuncU32U32ToU32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU32U32ToU32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU32U32ToU32) Bind(args []any) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU32U32ToU32) Func() any {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncU32U32ToU32)(nil)
)

// FuncU32U32U32ToU32 implements rfunc.Formula
type FuncU32U32U32ToU32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32) uint32
}

// NewFuncU32U32U32ToU32 return a new formula, from the provided function.
func NewFuncU32U32U32ToU32(rvars []string, fct func(arg00 uint32, arg01 uint32, arg02 uint32) uint32) *FuncU32U32U32ToU32 {
	return &FuncU32U32U32ToU32{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU32U32U32ToU32) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU32U32U32ToU32) Bind(args []any) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU32U32U32ToU32) Func() any {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncU32U32U32ToU32)(nil)
)
