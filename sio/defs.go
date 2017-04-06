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
	g_mark_record      uint32 = 0xabadcafe
	g_mark_block              = 0xdeadbeef
	g_opt_compress            = 0x00000001
	g_opt_not_compress        = 0xfffffffe
	g_align                   = 0x00000003
)

var (
	g_mark_block_b = []byte{222, 173, 190, 239}
)
