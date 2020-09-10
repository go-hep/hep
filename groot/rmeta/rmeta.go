// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rmeta provides tools to interoperate with ROOT Meta.
package rmeta // import "go-hep.org/x/hep/groot/rmeta"

import (
	"strings"
)

// CxxTemplateArgsOf extracts the typenames of a C++ templated typename.
// Ex:
//  std::map<K,V> -> []string{"K", "V"}
//  std::vector<T> -> []string{"T"}
//  Foo<T1,T2,std::map<K,V>> -> []string{"T1", "T2", "std::map<K,V>"}
func CxxTemplateArgsOf(typename string) []string {
	name := strings.TrimSpace(typename)
	name = name[strings.Index(name, "<")+1:] // drop heading 'xxx<'
	name = name[:len(name)-1]                // drop trailing '>'
	name = strings.TrimSpace(name)

	switch strings.Count(name, ",") {
	case 0:
		return []string{name}
	case 1:
		// easy case of std::map<K,V> where none of K or V are templated.
		i := strings.Index(name, ",")
		k := strings.TrimSpace(name[:i])
		v := strings.TrimSpace(name[i+1:])
		return []string{k, v}
	default:
		var (
			types []string
			bldr  strings.Builder
			tmpl  int
		)
		for _, s := range name {
			switch s {
			case '<':
				tmpl++
				bldr.WriteRune(s)
			case '>':
				tmpl--
				bldr.WriteRune(s)
			case ',':
				switch {
				case tmpl > 0:
					bldr.WriteRune(s)
				default:
					types = append(types, strings.TrimSpace(bldr.String()))
					bldr.Reset()
				}
			default:
				bldr.WriteRune(s)

			}
		}
		types = append(types, strings.TrimSpace(bldr.String()))
		return types
	}
}

// TypeName2Enum returns the Enum corresponding to the provided C++ (or Go) typename.
func TypeName2Enum(typename string) (Enum, bool) {
	switch typename {
	case "bool", "_Bool", "Bool_t":
		return Bool, true
	case "byte", "uint8", "uint8_t", "unsigned char", "UChar_t", "Byte_t":
		return Uint8, true
	case "uint16", "uint16_t", "unsigned short", "UShort_t":
		return Uint16, true
	case "uint32", "uint32_t", "unsigned int", "UInt_t":
		return Uint32, true
	case "uint64", "uint64_t", "unsigned long", "ULong_t", "ULong64_t":
		return Uint64, true

	case "char*":
		return CharStar, true
	case "Bits_t":
		return Bits, true

	case "int8", "int8_t", "char", "Char_t":
		return Int8, true
	case "int16", "int16_t", "short", "Short_t", "Version_t",
		"Font_t", "Style_t", "Marker_t", "Width_t",
		"Color_t",
		"SCoord_t":
		return Int16, true
	case "int32", "int32_t", "int", "Int_t":
		return Int32, true
	case "int64", "int64_t", "long", "Long_t", "Long64_t",
		"Seek_t":
		return Int64, true

	case "float32", "float", "Float_t", "float32_t",
		"Angle_t", "Size_t":
		return Float32, true
	case "float64", "double", "Double_t", "float64_t",
		"Coord_t":
		return Float64, true
	case "Float16_t", "Float16":
		return Float16, true
	case "Double32_t", "Double32":
		return Double32, true

	case "string", "TString", "Option_t":
		return TString, true
	case "std::string":
		return STLstring, true
	case "TObject":
		return TObject, true
	case "TNamed":
		return TNamed, true
	}

	return -1, false
}
