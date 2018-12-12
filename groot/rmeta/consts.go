// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rmeta // import "go-hep.org/x/hep/groot/rmeta"

//go:generate stringer -type Enum consts.go

// Enum is the set of ROOT streamer enums
type Enum int32

// constants for the streamers
const (
	Base       Enum = 0
	Char       Enum = 1
	Short      Enum = 2
	Int        Enum = 3
	Long       Enum = 4
	Float      Enum = 5
	Counter    Enum = 6
	CharStar   Enum = 7
	Double     Enum = 8
	Double32   Enum = 9
	LegacyChar Enum = 10
	UChar      Enum = 11
	UShort     Enum = 12
	UInt       Enum = 13
	ULong      Enum = 14
	Bits       Enum = 15
	Long64     Enum = 16
	ULong64    Enum = 17
	Bool       Enum = 18
	Float16    Enum = 19
	OffsetL    Enum = 20
	OffsetP    Enum = 40
	Object     Enum = 61
	Any        Enum = 62
	Objectp    Enum = 63
	ObjectP    Enum = 64
	TString    Enum = 65
	TObject    Enum = 66
	TNamed     Enum = 67
	Anyp       Enum = 68
	AnyP       Enum = 69
	AnyPnoVT   Enum = 70
	STLp       Enum = 71

	Skip  Enum = 100
	SkipL Enum = 120
	SkipP Enum = 140

	Conv  Enum = 200
	ConvL Enum = 220
	ConvP Enum = 240

	STL       Enum = 300
	STLstring Enum = 365

	Streamer   Enum = 500
	StreamLoop Enum = 501
)

// aliases for Go
const (
	Int8    Enum = Char
	Int16   Enum = Short
	Int32   Enum = Int
	Int64   Enum = Long
	Uint8   Enum = UChar
	Uint16  Enum = UShort
	Uint32  Enum = UInt
	Uint64  Enum = ULong
	Float32 Enum = Float
	Float64 Enum = Double
)

type ESTLType int32

// constants from core/foundation/inc/ESTLType.h
const (
	NotSTL      ESTLType = 0
	STLvector   ESTLType = 1
	STLlist     ESTLType = 2
	STLdeque    ESTLType = 3
	STLmap      ESTLType = 4
	STLmultimap ESTLType = 5
	STLset      ESTLType = 6
	STLmultiset ESTLType = 7
	STLbitset   ESTLType = 8
	// Here the c++11 containers start. Order counts. For example,
	// tstreamerelements in written rootfiles carry a value and we cannot
	// introduce shifts.
	STLforwardlist       ESTLType = 9
	STLunorderedset      ESTLType = 10
	STLunorderedmultiset ESTLType = 11
	STLunorderedmap      ESTLType = 12
	STLunorderedmultimap ESTLType = 13
	STLend               ESTLType = 14
	STLany               ESTLType = 300 /* TVirtualStreamerInfo::kSTL */
)
