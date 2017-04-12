// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"errors"
)

type Operation int

const (
	OpRead Operation = iota
	OpWrite
)

const (
	lenSB = 1
	lenDB = 2
	lenQB = 4
	lenOB = 8
)

const (
	recMarker      uint32 = 0xabadcafe
	blkMarker             = 0xdeadbeef
	ptagMarker            = 0xffffffff
	pntrMarker            = 0x00000000
	optCompress           = 0x00000001
	optNotCompress        = 0xfffffffe
	alignLen              = 0x00000003
)

var (
	blkMarkerBeg = []byte{222, 173, 190, 239}
)

var (
	errPointerIDOverflow = errors.New("sio: pointer id overflow")
)

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
