// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "fmt"

const g_rootio_debug = false

func myprintf(format string, args ...interface{}) (n int, err error) {
	if g_rootio_debug {
		return fmt.Printf(format, args...)
	}
	return
}
