// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"testing"
	"time"

	"go-hep.org/x/hep/groot/rbytes"
)

func TestDatime(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected a panic!")
		}
		var (
			got  = err.(error).Error()
			want = "rbase: TDatime year must be >= 1995"
		)
		if got != want {
			t.Fatalf("invalid panic.\ngot= %q\nwant=%q", got, want)
		}
	}()

	d := Datime(time.Date(1980, 1, 2, 15, 4, 5, 0, time.UTC))
	w := rbytes.NewWBuffer(nil, nil, 0, nil)
	_, err := d.MarshalROOT(w)
	if err != nil {
		t.Fatalf("could not marshal TDatime: %+v", err)
	}
}
