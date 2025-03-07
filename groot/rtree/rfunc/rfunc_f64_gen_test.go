// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncI32ToF64(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 int32) float64 {
		return 42
	}

	form := NewFuncI32ToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 1)
	ptrs[0] = new(int32)

	{
		bad := make([]any, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = any(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, any(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF32ToF64(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 float32) float64 {
		return 42
	}

	form := NewFuncF32ToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 1)
	ptrs[0] = new(float32)

	{
		bad := make([]any, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = any(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, any(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncToF64(t *testing.T) {

	var rvars []string

	fct := func() float64 {
		return 42
	}

	form := NewFuncToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	var ptrs []any

	{
		bad := make([]any, 1)
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF64ToF64(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 float64) float64 {
		return 42
	}

	form := NewFuncF64ToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 1)
	ptrs[0] = new(float64)

	{
		bad := make([]any, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = any(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, any(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF64F64ToF64(t *testing.T) {

	rvars := make([]string, 2)
	rvars[0] = "name-0"
	rvars[1] = "name-1"

	fct := func(arg00 float64, arg01 float64) float64 {
		return 42
	}

	form := NewFuncF64F64ToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 2)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)

	{
		bad := make([]any, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = any(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, any(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF64F64F64ToF64(t *testing.T) {

	rvars := make([]string, 3)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"

	fct := func(arg00 float64, arg01 float64, arg02 float64) float64 {
		return 42
	}

	form := NewFuncF64F64F64ToF64(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 3)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)

	{
		bad := make([]any, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = any(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, any(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
