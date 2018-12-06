// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rmeta // import "go-hep.org/x/hep/groot/rmeta"

// constants for the streamers
const (
	Base       = 0
	Char       = 1
	Short      = 2
	Int        = 3
	Long       = 4
	Float      = 5
	Counter    = 6
	CharStar   = 7
	Double     = 8
	Double32   = 9
	LegacyChar = 10
	UChar      = 11
	UShort     = 12
	UInt       = 13
	ULong      = 14
	Bits       = 15
	Long64     = 16
	ULong64    = 17
	Bool       = 18
	Float16    = 19
	OffsetL    = 20
	OffsetP    = 40
	Object     = 61
	Any        = 62
	Objectp    = 63
	ObjectP    = 64
	TString    = 65
	TObject    = 66
	TNamed     = 67
	Anyp       = 68
	AnyP       = 69
	AnyPnoVT   = 70
	STLp       = 71

	Skip  = 100
	SkipL = 120
	SkipP = 140

	Conv  = 200
	ConvL = 220
	ConvP = 240

	STL       = 300
	STLstring = 365

	Streamer   = 500
	StreamLoop = 501
)

// aliases for Go
const (
	Int8    = Char
	Int16   = Short
	Int32   = Int
	Int64   = Long
	Uint8   = UChar
	Uint16  = UShort
	Uint32  = UInt
	Uint64  = ULong
	Float32 = Float
	Float64 = Double
)

// constants from core/foundation/inc/ESTLType.h
const (
	NotSTL      = 0
	STLvector   = 1
	STLlist     = 2
	STLdeque    = 3
	STLmap      = 4
	STLmultimap = 5
	STLset      = 6
	STLmultiset = 7
	STLbitset   = 8
	// Here the c++11 containers start. Order counts. For example,
	// tstreamerelements in written rootfiles carry a value and we cannot
	// introduce shifts.
	STLforwardlist       = 9
	STLunorderedset      = 10
	STLunorderedmultiset = 11
	STLunorderedmap      = 12
	STLunorderedmultimap = 13
	STLend               = 14
	STLany               = 300 /* TVirtualStreamerInfo::kSTL */
)
