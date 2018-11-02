// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import "time"

func init() {
	nowUTC = func() time.Time {
		return datime2time(1576331001)
	}
}
