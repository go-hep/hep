// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

const (
	kNullTag = 0
	// on tag :
	kNewClassTag    = 0xFFFFFFFF
	kClassMask      = 0x80000000
	kMapOffset      = 2
	kByteCountVMask = 0x4000
	kByteCountMask  = 0x40000000

	kIsOnHeap     = 0x01000000
	kNotDeleted   = 0x02000000
	kZombie       = 0x04000000
	kBitMask      = 0x00ffffff
	kIsReferenced = 1 << 4

	//baskets
	kDisplacementMask = 0xFF000000
)

const (
	StreamedMemberWise = 0x4000
)
