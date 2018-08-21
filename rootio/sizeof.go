// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

// tstringSizeof returns the size in bytes of the TString structure.
func tstringSizeof(v string) int32 {
	n := int32(len(v))
	if n > 254 {
		return n + 1 + 4
	}
	return n + 1
}

// datimeSizeof returns the size in bytes of the TDatime structure.
func datimeSizeof() int32 {
	return 4
}
