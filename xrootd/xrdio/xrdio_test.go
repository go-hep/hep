// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdio_test

import (
	"errors"
	"os"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdio"
)

func TestFileCloseNil(t *testing.T) {
	var f *xrdio.File
	err := f.Close()
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("invalid error: got=%v, want=%v", err, os.ErrInvalid)
	}
}
