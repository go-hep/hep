// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
)

func TestBranchSetAddress(t *testing.T) {
	for _, tc := range []struct {
		name   string
		b      Branch
		ptr    interface{}
		panics string
		err    error
	}{
		{
			name: "0-leave",
			b:    &tbranch{named: *rbase.NewNamed("branch", "branch")},
			ptr:  nil,
			err:  fmt.Errorf("rtree: can not set address for a leaf-less branch (name=%q)", "branch"),
		},
		{
			name: "not-enough-fields",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafI{},
					&LeafF{},
				},
			},
			ptr: &struct{ i int32 }{},
			err: fmt.Errorf("rtree: fields/leaves number mismatch (name=%q, fields=%d, leaves=%d)", "branch", 1, 2),
		},
		{
			name: "too-many-fields",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafI{},
					&LeafI{},
				},
			},
			ptr: &struct{ f1, f2, f3 int32 }{},
			err: fmt.Errorf("rtree: fields/leaves number mismatch (name=%q, fields=%d, leaves=%d)", "branch", 3, 2),
		},
		{
			name: "invalid-field-type",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafF{},
				},
			},
			ptr:    &struct{ f1 int32 }{},
			panics: "invalid ptr type *struct { f1 int32 } (leaf=|*rtree.LeafF)",
		},
		{
			name: "invalid-field-type",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafI{},
					&LeafF{},
				},
			},
			ptr:    &struct{ F1, F2 int32 }{},
			panics: "invalid ptr type *int32 (leaf=|*rtree.LeafF)",
		},
		{
			name: "not-a-struct",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafI{},
					&LeafF{},
				},
			},
			ptr: []interface{}{new(int32), new(float32)},
			err: fmt.Errorf("rtree: multi-leaf branches need a pointer-to-struct (got=%s)", "[]interface {}"),
		},
		{
			name: "not-a-struct",
			b: &tbranch{
				named: *rbase.NewNamed("branch", "branch"),
				leaves: []Leaf{
					&LeafI{},
					&LeafF{},
				},
			},
			ptr: &[]interface{}{new(int32), new(float32)},
			err: fmt.Errorf("rtree: multi-leaf branches need a pointer-to-struct (got=%s)", "*[]interface {}"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic %q", tc.panics)
					}
					var got string
					switch e := err.(type) {
					case error:
						got = e.Error()
					case string:
						got = e
					default:
						got = fmt.Sprintf("%v", e)
					}

					if got, want := got, tc.panics; got != want {
						t.Fatalf("invalid panic message: got=%q, want=%q", got, want)
					}
				}()
			}
			err := tc.b.setAddress(tc.ptr)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error: got=%q, want=%q", got, want)
				}
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %+v", tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}
		})
	}
}
