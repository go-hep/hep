// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

// align4U32 returns sz adjusted to align at 4-byte boundaries
func align4U32(sz uint32) uint32 {
	return sz + (4-(sz&alignLen))&alignLen
}

// align4I32 returns sz adjusted to align at 4-byte boundaries
func align4I32(sz int32) int32 {
	return sz + (4-(sz&int32(alignLen)))&int32(alignLen)
}

// align4I64 returns sz adjusted to align at 4-byte boundaries
func align4I64(sz int64) int64 {
	return sz + (4-(sz&int64(alignLen)))&int64(alignLen)
}
