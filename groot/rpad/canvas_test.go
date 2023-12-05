// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rpad_test

import (
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rpad"
	"go-hep.org/x/hep/groot/rvers"
)

func TestCanvasRead(t *testing.T) {
	f, err := groot.Open("../testdata/tcanvas.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("c1")
	if err != nil {
		t.Fatal(err)
	}
	c := o.(*rpad.Canvas)

	if got, want := c.Name(), "c1"; got != want {
		t.Fatalf("invalid name: got=%q, want=%q", got, want)
	}

	if got, want := c.Title(), "c1-title"; got != want {
		t.Fatalf("invalid title: got=%q, want=%q", got, want)
	}

	if got, want := c.Class(), "TCanvas"; got != want {
		t.Fatalf("invalid class: got=%q, want=%q", got, want)
	}

	if got, want := int(c.RVersion()), rvers.Canvas; got != want {
		t.Fatalf("invalid version: got=%d, want=%d", got, want)
	}
}
