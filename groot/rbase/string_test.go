// Copyright ©2024 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase_test

import (
	"testing"

	"go-hep.org/x/hep/groot/rbase"
)

func TestStringVersion(t *testing.T) {
	var (
		str  = rbase.NewString("string")
		got  = str.RVersion()
		want int16
	)
	if got != want {
		t.Fatalf("invalid streamer version: got=%d, want=%d", got, want)
	}
}
