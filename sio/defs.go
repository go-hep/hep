// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

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
	optCompress           = 0x00000001
	optNotCompress        = 0xfffffffe
	alignLen              = 0x00000003
)

var (
	blkMarkerBeg = []byte{222, 173, 190, 239}
)
