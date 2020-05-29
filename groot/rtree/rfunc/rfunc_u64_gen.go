// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToU64 implements rfunc.Formula
type FuncToU64 struct {
	fct func() uint64
}

// NewFuncToU64 return a new formula, from the provided function.
func NewFuncToU64(rvars []string, fct func() uint64) *FuncToU64 {
	return &FuncToU64{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToU64) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToU64) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToU64) Func() interface{} {
	return func() uint64 {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToU64)(nil)
)

// FuncU64ToU64 implements rfunc.Formula
type FuncU64ToU64 struct {
	rvars []string
	arg0  *uint64
	fct   func(arg00 uint64) uint64
}

// NewFuncU64ToU64 return a new formula, from the provided function.
func NewFuncU64ToU64(rvars []string, fct func(arg00 uint64) uint64) *FuncU64ToU64 {
	return &FuncU64ToU64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU64ToU64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU64ToU64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU64ToU64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncU64ToU64)(nil)
)

// FuncU64U64ToU64 implements rfunc.Formula
type FuncU64U64ToU64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	fct   func(arg00 uint64, arg01 uint64) uint64
}

// NewFuncU64U64ToU64 return a new formula, from the provided function.
func NewFuncU64U64ToU64(rvars []string, fct func(arg00 uint64, arg01 uint64) uint64) *FuncU64U64ToU64 {
	return &FuncU64U64ToU64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU64U64ToU64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU64U64ToU64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU64U64ToU64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncU64U64ToU64)(nil)
)

// FuncU64U64U64ToU64 implements rfunc.Formula
type FuncU64U64U64ToU64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64) uint64
}

// NewFuncU64U64U64ToU64 return a new formula, from the provided function.
func NewFuncU64U64U64ToU64(rvars []string, fct func(arg00 uint64, arg01 uint64, arg02 uint64) uint64) *FuncU64U64U64ToU64 {
	return &FuncU64U64U64ToU64{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncU64U64U64ToU64) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncU64U64U64ToU64) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncU64U64U64ToU64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncU64U64U64ToU64)(nil)
)
