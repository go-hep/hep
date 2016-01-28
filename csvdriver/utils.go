// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvdriver

import "database/sql/driver"

func params(args []driver.Value) []interface{} {
	qargs := make([]interface{}, len(args))
	for i, arg := range args {
		qargs[i] = arg
	}
	return qargs
}
