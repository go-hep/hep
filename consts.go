// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

// start of payload in a TFile (in bytes)
const g_kBEGIN = 100

// constants for the streamers
const (
	kBase       = 0
	kChar       = 1
	kShort      = 2
	kInt        = 3
	kLong       = 4
	kFloat      = 5
	kCounter    = 6
	kCharStar   = 7
	kDouble     = 8
	kDouble32   = 9
	kLegacyChar = 10
	kUChar      = 11
	kUShort     = 12
	kUInt       = 13
	kULong      = 14
	kBits       = 15
	kLong64     = 16
	kULong64    = 17
	kBool       = 18
	kFloat16    = 19
	kOffsetL    = 20
	kOffsetP    = 40
	kObject     = 61
	kAny        = 62
	kObjectp    = 63
	kObjectP    = 64
	kTString    = 65
	kTObject    = 66
	kTNamed     = 67
	kAnyp       = 68
	kAnyP       = 69
	kAnyPnoVT   = 70
	kSTLp       = 71

	kSkip  = 100
	kSkipL = 120
	kSkipP = 140

	kConv  = 200
	kConvL = 220
	kConvP = 240

	kSTL       = 300
	kSTLstring = 365

	kStreamer   = 500
	kStreamLoop = 501

	kByteCountMask = 0x40000000
)

const (
	kNullTag = 0
	// on tag :
	kNewClassTag    = 0xFFFFFFFF
	kClassMask      = 0x80000000
	kMapOffset      = 2
	kByteCountVMask = 0x4000

	kIsReferenced = 1 << 4
	kIsOnHeap     = 0x01000000

	//baskets
	kDisplacementMask = 0xFF000000
)
