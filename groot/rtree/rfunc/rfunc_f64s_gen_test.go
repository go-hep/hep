// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncF32sToF64s(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 []float32) []float64 {
		return []float64{42}
	}

	form := NewFuncF32sToF64s(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 1)
	ptrs[0] = new([]float32)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() []float64)()
	if got, want := got, []float64([]float64{42}); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncF64sToF64s(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 []float64) []float64 {
		return []float64{42}
	}

	form := NewFuncF64sToF64s(rvars, fct)

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 1)
	ptrs[0] = new([]float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() []float64)()
	if got, want := got, []float64([]float64{42}); !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
