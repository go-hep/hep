// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// F64Ar0 implements rfunc.Formula
type F64Ar0 struct {
	fct func() float64
}

// NewF64Ar0 return a new formula, from the provided function.
func NewF64Ar0(rvars []string, fct func() float64) *F64Ar0 {
	return &F64Ar0{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *F64Ar0) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *F64Ar0) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *F64Ar0) Func() interface{} {
	return func() float64 {
		return f.fct()
	}
}

var (
	_ Formula = (*F64Ar0)(nil)
)

// F64Ar1 implements rfunc.Formula
type F64Ar1 struct {
	rvars []string
	arg0  *float64
	fct   func(arg00 float64) float64
}

// NewF64Ar1 return a new formula, from the provided function.
func NewF64Ar1(rvars []string, fct func(arg00 float64) float64) *F64Ar1 {
	return &F64Ar1{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *F64Ar1) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *F64Ar1) Bind(args []interface{}) error {
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
func (f *F64Ar1) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*F64Ar1)(nil)
)

// F64Ar2 implements rfunc.Formula
type F64Ar2 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	fct   func(arg00 float64, arg01 float64) float64
}

// NewF64Ar2 return a new formula, from the provided function.
func NewF64Ar2(rvars []string, fct func(arg00 float64, arg01 float64) float64) *F64Ar2 {
	return &F64Ar2{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *F64Ar2) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *F64Ar2) Bind(args []interface{}) error {
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
func (f *F64Ar2) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*F64Ar2)(nil)
)

// F64Ar3 implements rfunc.Formula
type F64Ar3 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64) float64
}

// NewF64Ar3 return a new formula, from the provided function.
func NewF64Ar3(rvars []string, fct func(arg00 float64, arg01 float64, arg02 float64) float64) *F64Ar3 {
	return &F64Ar3{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *F64Ar3) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *F64Ar3) Bind(args []interface{}) error {
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
func (f *F64Ar3) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*F64Ar3)(nil)
)
