// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
)

// FuncToBool implements rfunc.Formula
type FuncToBool struct {
	fct func() bool
}

// NewFuncToBool return a new formula, from the provided function.
func NewFuncToBool(rvars []string, fct func() bool) *FuncToBool {
	return &FuncToBool{
		fct: fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncToBool) RVars() []string { return nil }

// Bind implements rfunc.Formula
func (f *FuncToBool) Bind(args []any) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

// Func implements rfunc.Formula
func (f *FuncToBool) Func() any {
	return func() bool {
		return f.fct()
	}
}

var (
	_ Formula = (*FuncToBool)(nil)
)

// FuncF32ToBool implements rfunc.Formula
type FuncF32ToBool struct {
	rvars []string
	arg0  *float32
	fct   func(arg00 float32) bool
}

// NewFuncF32ToBool return a new formula, from the provided function.
func NewFuncF32ToBool(rvars []string, fct func(arg00 float32) bool) *FuncF32ToBool {
	return &FuncF32ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32ToBool) Bind(args []any) error {
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
func (f *FuncF32ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF32ToBool)(nil)
)

// FuncF32F32ToBool implements rfunc.Formula
type FuncF32F32ToBool struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	fct   func(arg00 float32, arg01 float32) bool
}

// NewFuncF32F32ToBool return a new formula, from the provided function.
func NewFuncF32F32ToBool(rvars []string, fct func(arg00 float32, arg01 float32) bool) *FuncF32F32ToBool {
	return &FuncF32F32ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32F32ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32F32ToBool) Bind(args []any) error {
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
func (f *FuncF32F32ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncF32F32ToBool)(nil)
)

// FuncF32F32F32ToBool implements rfunc.Formula
type FuncF32F32F32ToBool struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32) bool
}

// NewFuncF32F32F32ToBool return a new formula, from the provided function.
func NewFuncF32F32F32ToBool(rvars []string, fct func(arg00 float32, arg01 float32, arg02 float32) bool) *FuncF32F32F32ToBool {
	return &FuncF32F32F32ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF32F32F32ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF32F32F32ToBool) Bind(args []any) error {
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
func (f *FuncF32F32F32ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncF32F32F32ToBool)(nil)
)

// FuncF64ToBool implements rfunc.Formula
type FuncF64ToBool struct {
	rvars []string
	arg0  *float64
	fct   func(arg00 float64) bool
}

// NewFuncF64ToBool return a new formula, from the provided function.
func NewFuncF64ToBool(rvars []string, fct func(arg00 float64) bool) *FuncF64ToBool {
	return &FuncF64ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64ToBool) Bind(args []any) error {
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
func (f *FuncF64ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*FuncF64ToBool)(nil)
)

// FuncF64F64ToBool implements rfunc.Formula
type FuncF64F64ToBool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	fct   func(arg00 float64, arg01 float64) bool
}

// NewFuncF64F64ToBool return a new formula, from the provided function.
func NewFuncF64F64ToBool(rvars []string, fct func(arg00 float64, arg01 float64) bool) *FuncF64F64ToBool {
	return &FuncF64F64ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64F64ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64F64ToBool) Bind(args []any) error {
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
func (f *FuncF64F64ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*FuncF64F64ToBool)(nil)
)

// FuncF64F64F64ToBool implements rfunc.Formula
type FuncF64F64F64ToBool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64) bool
}

// NewFuncF64F64F64ToBool return a new formula, from the provided function.
func NewFuncF64F64F64ToBool(rvars []string, fct func(arg00 float64, arg01 float64, arg02 float64) bool) *FuncF64F64F64ToBool {
	return &FuncF64F64F64ToBool{
		rvars: rvars,
		fct:   fct,
	}
}

// RVars implements rfunc.Formula
func (f *FuncF64F64F64ToBool) RVars() []string { return f.rvars }

// Bind implements rfunc.Formula
func (f *FuncF64F64F64ToBool) Bind(args []any) error {
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
func (f *FuncF64F64F64ToBool) Func() any {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*FuncF64F64F64ToBool)(nil)
)
