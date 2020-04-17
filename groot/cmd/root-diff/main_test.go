// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import "testing"

func TestROOTDiff(t *testing.T) {
	const allkeys = ""
	err := rootdiff("../../testdata/small-flat-tree.root", "../../testdata/small-flat-tree.root", allkeys)
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
