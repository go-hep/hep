// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
)

// wobject wrapps a type created from a Streamer and implements the
// following interfaces:
//  - rootio.Object
//  - rootio.ROOTMarshaler
//  - rootio.ROOTUnmarshaler
type wobject struct {
	v         interface{}
	class     func(recv interface{}) string
	unmarshal func(recv interface{}, r *RBuffer) error
	marshal   func(recv interface{}, w *WBuffer) (int, error)
}

func (obj *wobject) Class() string {
	return obj.class(obj.v)
}

func (obj *wobject) UnmarshalROOT(r *RBuffer) error {
	return obj.unmarshal(obj.v, r)
}

func (obj *wobject) MarshalROOT(w *WBuffer) (int, error) {
	return obj.marshal(obj.v, w)
}

var (
	_ Object          = (*wobject)(nil)
	_ ROOTMarshaler   = (*wobject)(nil)
	_ ROOTUnmarshaler = (*wobject)(nil)
)

var builtins = map[string]reflect.Type{
	"TObject":        reflect.TypeOf((*tobject)(nil)).Elem(),
	"TString":        reflect.TypeOf(""),
	"TNamed":         reflect.TypeOf((*tnamed)(nil)).Elem(),
	"TList":          reflect.TypeOf((*tlist)(nil)).Elem(),
	"TObjArray":      reflect.TypeOf((*tobjarray)(nil)).Elem(),
	"TObjString":     reflect.TypeOf((*tobjstring)(nil)).Elem(),
	"TTree":          reflect.TypeOf((*ttree)(nil)).Elem(),
	"TBranch":        reflect.TypeOf((*tbranch)(nil)).Elem(),
	"TBranchElement": reflect.TypeOf((*tbranchElement)(nil)).Elem(),
}

var cxxbuiltins = map[string]reflect.Type{
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
	"Float16_t":  reflect.TypeOf(Float16(0)),
	"Double_t":   reflect.TypeOf(float64(0)),
	"Double32_t": reflect.TypeOf(Double32(0)),

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

var gotype2ROOTEnum = map[reflect.Type]int32{
	reflect.TypeOf(int8(0)):  kChar,
	reflect.TypeOf(int16(0)): kShort,
	reflect.TypeOf(int32(0)): kInt,
	reflect.TypeOf(int64(0)): kLong,
	//	reflect.TypeOf(int64(0)): kLong64,
	reflect.TypeOf(float32(0)): kFloat,
	reflect.TypeOf(float64(0)): kDouble,
	reflect.TypeOf(uint8(0)):   kUChar,
	//	reflect.TypeOf(uint8(0)): kCharStar,
	reflect.TypeOf(uint16(0)): kUShort,
	reflect.TypeOf(uint32(0)): kUInt,
	reflect.TypeOf(uint64(0)): kULong,
	//	reflect.TypeOf(uint64(0)): kULong64,
	reflect.TypeOf(false):       kBool,
	reflect.TypeOf(Double32(0)): kDouble32,
}
