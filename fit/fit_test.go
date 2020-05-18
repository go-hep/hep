// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"os"
	"path"
	"testing"
)

type chkplotFunc func(ExampleFunc func(), t *testing.T, filenames ...string)

func checkPlot(f chkplotFunc) chkplotFunc {
	return func(ex func(), t *testing.T, filenames ...string) {
		t.Helper()
		f(ex, t, filenames...)
		if t.Failed() {
			return
		}
		for _, fname := range filenames {
			_ = os.Remove(path.Join("testdata", fname))
		}
	}
}
