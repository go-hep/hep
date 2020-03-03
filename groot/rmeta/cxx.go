// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rmeta

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/root"
)

var GoType2ROOTEnum = map[reflect.Type]Enum{
	reflect.TypeOf(int8(0)):  Char,
	reflect.TypeOf(int16(0)): Short,
	reflect.TypeOf(int32(0)): Int,
	reflect.TypeOf(int64(0)): Long,
	//	reflect.TypeOf(int64(0)): Long64,
	reflect.TypeOf(float32(0)): Float,
	reflect.TypeOf(float64(0)): Double,
	reflect.TypeOf(uint8(0)):   UChar,
	//	reflect.TypeOf(uint8(0)): CharStar,
	reflect.TypeOf(uint16(0)): UShort,
	reflect.TypeOf(uint32(0)): UInt,
	reflect.TypeOf(uint64(0)): ULong,
	//	reflect.TypeOf(uint64(0)): ULong64,
	reflect.TypeOf(false):            Bool,
	reflect.TypeOf(root.Float16(0)):  Float16,
	reflect.TypeOf(root.Double32(0)): Double32,
	reflect.TypeOf(""):               TString,
}

var GoType2Cxx = map[string]string{
	"bool":    "bool",
	"byte":    "unsigned char",
	"uint":    "unsigned int",
	"uint8":   "unsigned char",
	"uint16":  "unsigned short",
	"uint32":  "unsigned int",
	"uint64":  "unsigned long",
	"int":     "int",
	"int8":    "char",
	"int16":   "short",
	"int32":   "int",
	"int64":   "long",
	"float32": "float",
	"float64": "double",
}

var CxxBuiltins = map[string]reflect.Type{
	"bool": reflect.TypeOf(false),

	/*
		"uint":   reflect.TypeOf(uint(0)),
		"uint8":  reflect.TypeOf(uint8(0)),
		"uint16": reflect.TypeOf(uint16(0)),
		"uint32": reflect.TypeOf(uint32(0)),
		"uint64": reflect.TypeOf(uint64(0)),

		"int":   reflect.TypeOf(int(0)),
		"int8":  reflect.TypeOf(int8(0)),
		"int16": reflect.TypeOf(int16(0)),
		"int32": reflect.TypeOf(int32(0)),
		"int64": reflect.TypeOf(int64(0)),

		"float32": reflect.TypeOf(float32(0)),
		"float64": reflect.TypeOf(float64(0)),
	*/

	// C/C++ builtins

	"unsigned":       reflect.TypeOf(uint(0)),
	"unsigned char":  reflect.TypeOf(uint8(0)),
	"unsigned short": reflect.TypeOf(uint16(0)),
	"unsigned int":   reflect.TypeOf(uint32(0)),
	"unsigned long":  reflect.TypeOf(uint64(0)),

	//"int":   reflect.TypeOf(int(0)),
	"char":  reflect.TypeOf(int8(0)),
	"short": reflect.TypeOf(int16(0)),
	"int":   reflect.TypeOf(int32(0)),
	"long":  reflect.TypeOf(int64(0)),

	"float":  reflect.TypeOf(float32(0)),
	"double": reflect.TypeOf(float64(0)),

	"string": reflect.TypeOf(""),

	// ROOT builtins
	"Bool_t": reflect.TypeOf(true),

	"Byte_t": reflect.TypeOf(uint8(0)),

	"Char_t":    reflect.TypeOf(int8(0)),
	"UChar_t":   reflect.TypeOf(uint8(0)),
	"Short_t":   reflect.TypeOf(int16(0)),
	"UShort_t":  reflect.TypeOf(uint16(0)),
	"Int_t":     reflect.TypeOf(int32(0)),
	"UInt_t":    reflect.TypeOf(uint32(0)),
	"Seek_t":    reflect.TypeOf(int64(0)),  // FIXME(sbinet): not portable
	"Long_t":    reflect.TypeOf(int64(0)),  // FIXME(sbinet): not portable
	"ULong_t":   reflect.TypeOf(uint64(0)), // FIXME(sbinet): not portable
	"Long64_t":  reflect.TypeOf(int64(0)),
	"ULong64_t": reflect.TypeOf(uint64(0)),

	"Float_t":    reflect.TypeOf(float32(0)),
	"Float16_t":  reflect.TypeOf(root.Float16(0)),
	"Double_t":   reflect.TypeOf(float64(0)),
	"Double32_t": reflect.TypeOf(root.Double32(0)),

	"Version_t": reflect.TypeOf(int16(0)),
	"Option_t":  reflect.TypeOf(""),
	"Ssiz_t":    reflect.TypeOf(int(0)),
	"Real_t":    reflect.TypeOf(float32(0)),

	"Axis_t": reflect.TypeOf(float64(0)),
	"Stat_t": reflect.TypeOf(float64(0)),

	"Font_t":   reflect.TypeOf(int16(0)),
	"Style_t":  reflect.TypeOf(int16(0)),
	"Marker_t": reflect.TypeOf(int16(0)),
	"Width_t":  reflect.TypeOf(int16(0)),
	"Color_t":  reflect.TypeOf(int16(0)),
	"SCoord_t": reflect.TypeOf(int16(0)),
	"Coord_t":  reflect.TypeOf(float64(0)),
	"Angle_t":  reflect.TypeOf(float32(0)),
	"Size_t":   reflect.TypeOf(float32(0)),
}

// STLNameFor creates a regular C++ STL container name given a STL enum type
// and a ROOT enum value for the contained element.
func STLNameFor(vtype ESTLType, ctype Enum) string {
	ename := rmeta2Name(ctype)
	if strings.HasSuffix(ename, ">") {
		ename += " "
	}

	typfmt := "%s"
	switch vtype {
	case STLvector:
		typfmt = "vector<%s>"
	case STLlist:
		typfmt = "list<%s>"
	case STLdeque:
		typfmt = "deque<%s>"
	case STLmap:
		typfmt = "map<%s>"
	case STLmultimap:
		typfmt = "multimap<%s>"
	case STLset:
		typfmt = "set<%s>"
	case STLmultiset:
		typfmt = "multiset<%s>"
	case STLbitset:
		typfmt = "bitset<%s>"
	case STLforwardlist:
		typfmt = "forward_list<%s>"
	case STLunorderedset:
		typfmt = "unordered_set<%s>"
	case STLunorderedmultiset:
		typfmt = "unordered_multiset<%s>"
	case STLunorderedmap:
		typfmt = "unordered_map<%s>"
	case STLunorderedmultimap:
		typfmt = "unordered_multimap<%s>"
	}

	return fmt.Sprintf(typfmt, ename)
}

func rmeta2Name(t Enum) string {
	switch t {
	case Char:
		return "char"
	case Short:
		return "short"
	case Int:
		return "int"
	case Long:
		return "long"
	case Float:
		return "float"
	case Counter:
		return "int"
	case CharStar:
		return "char*"
	case Double:
		return "double"
	case Double32:
		return "Double32_t"
	case LegacyChar:
		return "char" // FIXME(sbinet)
	case UChar:
		return "unsigned char"
	case UShort:
		return "unsigned short"
	case UInt:
		return "unsigned int"
	case ULong:
		return "unsigned long"
	case Bits:
		return "bits" // FIXME(sbinet)
	case Long64:
		return "Long64_t"
	case ULong64:
		return "ULong64_t"
	case Bool:
		return "bool"
	case Float16:
		return "Float16_t"
	case Object:
		// class derived from TObject
		return "TObject" // FIXME(sbinet): better handling?
	case TString:
		return "TString"
	case TObject:
		return "TObject"
	case TNamed:
		return "TNamed"
	case STLstring:
		return "string"
	}
	panic(fmt.Errorf("not implemented: t=%d (%v)", t, t))
}
