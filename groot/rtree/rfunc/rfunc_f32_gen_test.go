// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncToF32(t *testing.T) {

	var rvars []string

	fct := func() float32 {
		return 42
	}

	form := NewFuncToF32(rvars, fct)

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

	got := form.Func().(func() float32)()
	if got, want := got, float32(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF32ToF32(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 float32) float32 {
		return 42
	}

	form := NewFuncF32ToF32(rvars, fct)

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

	got := form.Func().(func() float32)()
	if got, want := got, float32(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF32F32ToF32(t *testing.T) {

	rvars := make([]string, 2)
	rvars[0] = "name-0"
	rvars[1] = "name-1"

	fct := func(arg00 float32, arg01 float32) float32 {
		return 42
	}

	form := NewFuncF32F32ToF32(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 2)
	ptrs[0] = new(float32)
	ptrs[1] = new(float32)

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

	got := form.Func().(func() float32)()
	if got, want := got, float32(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF32F32F32ToF32(t *testing.T) {

	rvars := make([]string, 3)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"

	fct := func(arg00 float32, arg01 float32, arg02 float32) float32 {
		return 42
	}

	form := NewFuncF32F32F32ToF32(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]any, 3)
	ptrs[0] = new(float32)
	ptrs[1] = new(float32)
	ptrs[2] = new(float32)

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

	got := form.Func().(func() float32)()
	if got, want := got, float32(42); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
