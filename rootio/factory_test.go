// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"strings"
	"testing"
)

func TestFactory(t *testing.T) {
	n := Factory.Len()
	if got, want := len(Factory.Keys()), n; got != want {
		t.Fatalf("got=%d, want=%d", got, want)
	}

	for _, name := range Factory.Keys() {
		if strings.HasPrefix(name, "*rootio") || strings.HasPrefix(name, "rootio") {
			continue
		}
		fct := Factory.Get(name)
		obj := fct().Interface().(Object)
		if got, want := obj.Class(), name; got != want {
			t.Fatalf("got=%q, want=%q", got, want)
		}
	}
}
