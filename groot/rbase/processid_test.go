// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase_test

import (
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
)

func TestProcessID(t *testing.T) {
	var pid rbase.ProcessID

	if got, want := pid.Name(), ""; got != want {
		t.Fatalf("invalid name. got=%q, want=%q", got, want)
	}
	if got, want := pid.Title(), ""; got != want {
		t.Fatalf("invalid title. got=%q, want=%q", got, want)
	}

	pid.SetName("my-name")
	pid.SetTitle("my-title")

	if got, want := pid.Name(), "my-name"; got != want {
		t.Fatalf("invalid name. got=%q, want=%q", got, want)
	}
	if got, want := pid.Title(), "my-title"; got != want {
		t.Fatalf("invalid title. got=%q, want=%q", got, want)
	}
	if got, want := pid.String(), `TProcessID{Name: my-name, Title: my-title}`; got != want {
		t.Fatalf("invalid string representation. got=%s, want=%s", got, want)
	}

	t.Run("read-from-root", func(t *testing.T) {
		f, err := groot.Open("../testdata/pid.root")
		if err != nil {
			t.Fatalf("%+v", err)
		}
		defer f.Close()

		o, err := f.Get("type-TProcessID")
		if err != nil {
			t.Fatalf("%+v", err)
		}

		pid := o.(*rbase.ProcessID)

		if got, want := pid.Name(), "my-pid"; got != want {
			t.Fatalf("invalid name. got=%q, want=%q", got, want)
		}
		if got, want := pid.Title(), "my-title"; got != want {
			t.Fatalf("invalid title. got=%q, want=%q", got, want)
		}
	})
}
