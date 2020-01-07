// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtypes_test

import (
	"testing"

	_ "go-hep.org/x/hep/groot/rbase" // import factories for rbase types
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestFactory(t *testing.T) {
	n := rtypes.Factory.Len()
	if got, want := n, 9; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	if got, want := len(rtypes.Factory.Keys()), n; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	for _, name := range rtypes.Factory.Keys() {
		fct := rtypes.Factory.Get(name)
		obj := fct().Interface().(root.Object)
		if got, want := obj.Class(), name; got != want {
			t.Fatalf("got=%q, want=%q", got, want)
		}
	}
}
