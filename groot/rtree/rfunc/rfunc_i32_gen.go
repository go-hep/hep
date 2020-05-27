// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// I32Ar0 implements rfunc.Formula
type I32Ar0 struct {
	fct func() int32
}

// NewI32Ar0 return a new formula, from the provided function.
func NewI32Ar0(rvars []string, fct func() int32) *I32Ar0 {
	return &I32Ar0{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *I32Ar0) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *I32Ar0) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *I32Ar0) Func() interface{} {
	return func() int32 {
		return f.fct()
	}
}

var (
	_ Formula = (*I32Ar0)(nil)
)

// I32Ar1 implements rfunc.Formula
type I32Ar1 struct {
	rvars []string
	arg0  *int32
	fct   func(arg00 int32) int32
}

// NewI32Ar1 return a new formula, from the provided function.
func NewI32Ar1(rvars []string, fct func(arg00 int32) int32) *I32Ar1 {
	return &I32Ar1{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *I32Ar1) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *I32Ar1) Bind(args []interface{}) error {
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
func (f *I32Ar1) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*I32Ar1)(nil)
)

// I32Ar2 implements rfunc.Formula
type I32Ar2 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	fct   func(arg00 int32, arg01 int32) int32
}

// NewI32Ar2 return a new formula, from the provided function.
func NewI32Ar2(rvars []string, fct func(arg00 int32, arg01 int32) int32) *I32Ar2 {
	return &I32Ar2{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *I32Ar2) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *I32Ar2) Bind(args []interface{}) error {
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
func (f *I32Ar2) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*I32Ar2)(nil)
)

// I32Ar3 implements rfunc.Formula
type I32Ar3 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32) int32
}

// NewI32Ar3 return a new formula, from the provided function.
func NewI32Ar3(rvars []string, fct func(arg00 int32, arg01 int32, arg02 int32) int32) *I32Ar3 {
	return &I32Ar3{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *I32Ar3) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *I32Ar3) Bind(args []interface{}) error {
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
func (f *I32Ar3) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*I32Ar3)(nil)
)
