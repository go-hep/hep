// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// U32Ar0 implements rfunc.Formula
type U32Ar0 struct {
	fct func() uint32
}

// NewU32Ar0 return a new formula, from the provided function.
func NewU32Ar0(rvars []string, fct func() uint32) *U32Ar0 {
	return &U32Ar0{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *U32Ar0) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *U32Ar0) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *U32Ar0) Func() interface{} {
	return func() uint32 {
		return f.fct()
	}
}

var (
	_ Formula = (*U32Ar0)(nil)
)

// U32Ar1 implements rfunc.Formula
type U32Ar1 struct {
	rvars []string
	arg0  *uint32
	fct   func(arg00 uint32) uint32
}

// NewU32Ar1 return a new formula, from the provided function.
func NewU32Ar1(rvars []string, fct func(arg00 uint32) uint32) *U32Ar1 {
	return &U32Ar1{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *U32Ar1) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *U32Ar1) Bind(args []interface{}) error {
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
func (f *U32Ar1) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*U32Ar1)(nil)
)

// U32Ar2 implements rfunc.Formula
type U32Ar2 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	fct   func(arg00 uint32, arg01 uint32) uint32
}

// NewU32Ar2 return a new formula, from the provided function.
func NewU32Ar2(rvars []string, fct func(arg00 uint32, arg01 uint32) uint32) *U32Ar2 {
	return &U32Ar2{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *U32Ar2) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *U32Ar2) Bind(args []interface{}) error {
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
func (f *U32Ar2) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*U32Ar2)(nil)
)

// U32Ar3 implements rfunc.Formula
type U32Ar3 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32) uint32
}

// NewU32Ar3 return a new formula, from the provided function.
func NewU32Ar3(rvars []string, fct func(arg00 uint32, arg01 uint32, arg02 uint32) uint32) *U32Ar3 {
	return &U32Ar3{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *U32Ar3) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *U32Ar3) Bind(args []interface{}) error {
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
func (f *U32Ar3) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*U32Ar3)(nil)
)
