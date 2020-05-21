// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

type Operation int

const (
	OpRead Operation = iota
	OpWrite
)

const (
	recMarker      uint32 = 0xabadcafe
	blkMarker      uint32 = 0xdeadbeef
	ptagMarker     uint32 = 0xffffffff
	pntrMarker     uint32 = 0x00000000
	optCompress    uint32 = 0x00000001
	optNotCompress uint32 = 0xfffffffe
	alignLen       uint32 = 0x00000003
)

var (
	blkMarkerBeg = []byte{222, 173, 190, 239}
)

// align4U32 returns sz adjusted to align at 4-byte boundaries
func align4U32(sz uint32) uint32 {
	return sz + (4-(sz&alignLen))&alignLen
}
