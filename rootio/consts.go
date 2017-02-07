// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import "reflect"

// start of payload in a TFile (in bytes)
const kBEGIN = 100

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

// constants from core/foundation/inc/ESTLType.h
const (
	kNotSTL      = 0
	kSTLvector   = 1
	kSTLlist     = 2
	kSTLdeque    = 3
	kSTLmap      = 4
	kSTLmultimap = 5
	kSTLset      = 6
	kSTLmultiset = 7
	kSTLbitset   = 8
	// Here the c++11 containers start. Order counts. For example,
	// tstreamerelements in written rootfiles carry a value and we cannot
	// introduce shifts.
	kSTLforwardlist       = 9
	kSTLunorderedset      = 10
	kSTLunorderedmultiset = 11
	kSTLunorderedmap      = 12
	kSTLunorderedmultimap = 13
	kSTLend               = 14
	kSTLany               = 300 /* TVirtualStreamerInfo::kSTL */
)

const (
	kNullTag = 0
	// on tag :
	kNewClassTag    = 0xFFFFFFFF
	kClassMask      = 0x80000000
	kMapOffset      = 2
	kByteCountVMask = 0x4000

	kIsOnHeap     = 0x01000000
	kNotDeleted   = 0x02000000
	kZombie       = 0x04000000
	kBitMask      = 0x00ffffff
	kIsReferenced = 1 << 4

	//baskets
	kDisplacementMask = 0xFF000000
)

var ptrSize = reflect.TypeOf((*int)(nil)).Size()
