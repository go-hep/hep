// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

const (
	kIsOnHeap     = 0x01000000
	kNotDeleted   = 0x02000000
	kZombie       = 0x04000000
	kBitMask      = 0x00ffffff
	kIsReferenced = 1 << 4
)
