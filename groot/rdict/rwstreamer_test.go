// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rvers"
)

func TestRWStream(t *testing.T) {
	const kind = rbytes.ObjectWise // FIXME(sbinet): also test MemberWise.

	for _, tc := range []struct {
		name string
		skip bool
		si   *StreamerInfo
		ptr  interface{}
		deps []rbytes.StreamerInfo
		err  error
	}{
		{
			name: "bits",
			ptr: &struct {
				F uint32
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Bits,
						esize: 4,
						ename: "Bits_t",
					}},
				},
			},
		},
		{
			name: "bool",
			ptr: &struct {
				F bool
			}{true},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Bool,
						esize: 1,
						ename: "bool",
					}},
				},
			},
		},
		{
			name: "uint8",
			ptr: &struct {
				F uint8
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Uint8,
						esize: 1,
						ename: "uint8",
					}},
				},
			},
		},
		{
			name: "uint16",
			ptr: &struct {
				F uint16
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Uint16,
						esize: 2,
						ename: "uint16",
					}},
				},
			},
		},
		{
			name: "uint32",
			ptr: &struct {
				F uint32
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Uint32,
						esize: 4,
						ename: "uint32",
					}},
				},
			},
		},
		{
			name: "uint64",
			ptr: &struct {
				F uint64
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Uint64,
						esize: 8,
						ename: "uint64",
					}},
				},
			},
		},
		{
			name: "int8",
			ptr: &struct {
				F int8
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Int8,
						esize: 1,
						ename: "int8",
					}},
				},
			},
		},
		{
			name: "int16",
			ptr: &struct {
				F int16
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Int16,
						esize: 2,
						ename: "int16",
					}},
				},
			},
		},
		{
			name: "int32",
			ptr: &struct {
				F int32
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Int32,
						esize: 4,
						ename: "int32",
					}},
				},
			},
		},
		{
			name: "int64",
			ptr: &struct {
				F int64
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Int64,
						esize: 8,
						ename: "int64",
					}},
				},
			},
		},
		{
			name: "float32",
			ptr: &struct {
				F float32
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Float32,
						esize: 4,
						ename: "float32",
					}},
				},
			},
		},
		{
			name: "float64",
			ptr: &struct {
				F float64
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{StreamerElement{
						named: *rbase.NewNamed("F", ""),
						etype: rmeta.Float64,
						esize: 8,
						ename: "float64",
					}},
				},
			},
		},
		{
			name: "float16",
			ptr: &struct {
				F root.Float16
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("F", "[0,42]"),
						Type:  rmeta.Float16,
						Size:  2,
						EName: "Float16_t",
					}.New()},
				},
			},
		},
		{
			name: "double32",
			ptr: &struct {
				F root.Double32
			}{42},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("F", "[0,42]"),
						Type:  rmeta.Double32,
						Size:  4,
						EName: "Double32_t",
					}.New()},
				},
			},
		},
		{
			name: "pchar",
			ptr: &struct {
				F string
			}{"hello"},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("F", ""),
						Type:  rmeta.CharStar,
						Size:  8,
						EName: "char*",
					}.New()},
				},
			},
		},
		{
			name: "TString",
			ptr: &struct {
				F string
			}{"hello"},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{Element{
						Name:  *rbase.NewNamed("F", ""),
						Type:  rmeta.TString,
						Size:  24,
						EName: "TString",
					}.New()},
				},
			},
		},
		{
			name: "STL-string",
			ptr: &struct {
				F string
			}{"hello"},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerSTLstring{
						StreamerSTL: StreamerSTL{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("F", ""),
								Type:   rmeta.STLstring,
								Size:   32,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "string",
							}.New(),
							vtype: rmeta.ESTLType(rmeta.STLstring),
							ctype: rmeta.STLstring,
						},
					},
				},
			},
		},
		{
			name: "TObject",
			ptr: &struct {
				F rbase.Object
			}{*rbase.NewObject()},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObject{Element{
						Name:  *rbase.NewNamed("F", ""),
						Type:  rmeta.TObject,
						Size:  16,
						EName: "TObject",
					}.New()},
				},
			},
		},
		{
			name: "TNamed",
			ptr: &struct {
				F rbase.Named
			}{*rbase.NewNamed("hello", "world")},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObject{Element{
						Name:  *rbase.NewNamed("F", ""),
						Type:  rmeta.TNamed,
						Size:  64,
						EName: "TNamed",
					}.New()},
				},
			},
		},
		{
			name: "arr-bool",
			ptr: &struct {
				F [3]bool
			}{[3]bool{true, false, true}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Bool,
						Size:   3,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "bool",
					}.New()},
				},
			},
		},
		{
			name: "arr-int8",
			ptr: &struct {
				F [3]int8
			}{[3]int8{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Int8,
						Size:   3,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int8",
					}.New()},
				},
			},
		},
		{
			name: "arr-int16",
			ptr: &struct {
				F [3]int16
			}{[3]int16{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Int16,
						Size:   3 * 2,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int16",
					}.New()},
				},
			},
		},
		{
			name: "arr-int32",
			ptr: &struct {
				F [3]int32
			}{[3]int32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Int32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int32",
					}.New()},
				},
			},
		},
		{
			name: "arr-int64",
			ptr: &struct {
				F [3]int64
			}{[3]int64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Int64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int64",
					}.New()},
				},
			},
		},
		{
			name: "arr-uint8",
			ptr: &struct {
				F [3]uint8
			}{[3]uint8{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Uint8,
						Size:   3,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint8",
					}.New()},
				},
			},
		},
		{
			name: "arr-uint16",
			ptr: &struct {
				F [3]uint16
			}{[3]uint16{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Uint16,
						Size:   3 * 2,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint16",
					}.New()},
				},
			},
		},
		{
			name: "arr-uint32",
			ptr: &struct {
				F [3]uint32
			}{[3]uint32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Uint32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint32",
					}.New()},
				},
			},
		},
		{
			name: "arr-uint64",
			ptr: &struct {
				F [3]uint64
			}{[3]uint64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Uint64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint64",
					}.New()},
				},
			},
		},
		{
			name: "arr-float32",
			ptr: &struct {
				F [3]float32
			}{[3]float32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Float32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "float32",
					}.New()},
				},
			},
		},
		{
			name: "arr-float64",
			ptr: &struct {
				F [3]float64
			}{[3]float64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.OffsetL + rmeta.Float64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "float64",
					}.New()},
				},
			},
		},
		{
			name: "arr-float16",
			ptr: &struct {
				F [3]root.Float16
			}{[3]root.Float16{42, 42, 42}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", "[3]f/[0,42]"),
						Type:   rmeta.OffsetL + rmeta.Float16,
						Size:   3 * 2,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "Float16_t",
					}.New()},
				},
			},
		},
		{
			name: "arr-double32",
			ptr: &struct {
				F [3]root.Double32
			}{[3]root.Double32{42, 42, 42}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", "[3]/d[0, 42]"),
						Type:   rmeta.OffsetL + rmeta.Double32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "Double32_t",
					}.New()},
				},
			},
		},
		{
			name: "arr-pchar",
			ptr: &struct {
				F [3]string
			}{[3]string{"hello", "world", "Go-HEP"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.CharStar,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "char*",
					}.New()},
				},
			},
		},
		{
			name: "arr-TString",
			ptr: &struct {
				F [3]string
			}{[3]string{"hello", "world", "Go-HEP"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.TString,
						Size:   3 * 24,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TString",
					}.New()},
				},
			},
		},
		{
			name: "arr-TObject",
			ptr: &struct {
				F [3]rbase.Object
			}{[3]rbase.Object{*rbase.NewObject(), *rbase.NewObject(), *rbase.NewObject()}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObject{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.TObject,
						Size:   3 * 16,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TObject",
					}.New()},
				},
			},
		},
		{
			name: "arr-TNamed",
			ptr: &struct {
				F [3]rbase.Named
			}{[3]rbase.Named{*rbase.NewNamed("n1", "t1"), *rbase.NewNamed("n2", "t2"), *rbase.NewNamed("n3", "t3")}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObject{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.TNamed,
						Size:   3 * 64,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TNamed",
					}.New()},
				},
			},
		},
		{
			name: "arr-TObjString",
			ptr: func() interface{} {
				type T struct {
					F [3]rbase.ObjString
				}
				return &T{
					F: [3]rbase.ObjString{
						*rbase.NewObjString("str-1"),
						*rbase.NewObjString("str-2"),
						*rbase.NewObjString("str-3"),
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObject{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.Object,
						Size:   3 * 40,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TObjString",
					}.New()},
				},
			},
		},
		{
			name: "arr-Pos",
			ptr: func() interface{} {
				type Pos struct {
					X float32
					Y float64
				}
				type T struct {
					F [3]Pos
				}
				return &T{
					F: [3]Pos{
						{1, 2},
						{3, 4},
						{5, 6},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{Element{
						Name:   *rbase.NewNamed("F", "[3]"),
						Type:   rmeta.OffsetL + rmeta.Any,
						Size:   3 * (4 + 8),
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "Pos",
					}.New()},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Pos", "Pos"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("X", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("Y", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "sli-counter-32",
			ptr: &struct {
				N int32
				F []bool
			}{3, []bool{true, false, true}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Counter,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Bool,
						Size:  1,
						EName: "bool*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-counter-64",
			ptr: &struct {
				N int64
				F []bool
			}{3, []bool{true, false, true}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Counter,
						Size:  8,
						EName: "int64",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Bool,
						Size:  1,
						EName: "bool*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-bool",
			ptr: &struct {
				N int32
				F []bool
			}{3, []bool{true, false, true}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Bool,
						Size:  1,
						EName: "bool*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-int8",
			ptr: &struct {
				N int32
				F []int8
			}{3, []int8{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int8,
						Size:  1,
						EName: "int8*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-int16",
			ptr: &struct {
				N int32
				F []int16
			}{3, []int16{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int16,
						Size:  2,
						EName: "int16*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-int32",
			ptr: &struct {
				N int32
				F []int32
			}{3, []int32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int32,
						Size:  4,
						EName: "int32*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-int64",
			ptr: &struct {
				N int32
				F []int64
			}{3, []int64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int64,
						Size:  8,
						EName: "int64*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-uint8",
			ptr: &struct {
				N int32
				F []uint8
			}{3, []uint8{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint8,
						Size:  1,
						EName: "uint8*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-uint16",
			ptr: &struct {
				N int32
				F []uint16
			}{3, []uint16{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint16,
						Size:  2,
						EName: "uint16*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-uint32",
			ptr: &struct {
				N int32
				F []uint32
			}{3, []uint32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint32,
						Size:  4,
						EName: "uint32*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-uint64",
			ptr: &struct {
				N int32
				F []uint64
			}{3, []uint64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint64,
						Size:  8,
						EName: "uint64*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-float32",
			ptr: &struct {
				N int32
				F []float32
			}{3, []float32{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Float32,
						Size:  4,
						EName: "float32*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-float64",
			ptr: &struct {
				N int32
				F []float64
			}{3, []float64{42, 43, 44}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Float64,
						Size:  8,
						EName: "float64*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-float16",
			ptr: &struct {
				N int32
				F []root.Float16
			}{3, []root.Float16{42, 42, 42}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]/f[0,42]"),
						Type:  rmeta.OffsetP + rmeta.Float16,
						Size:  2,
						EName: "Float16_t*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-double32",
			ptr: &struct {
				N int32
				F []root.Double32
			}{3, []root.Double32{42, 42, 42}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]/d[0,42]"),
						Type:  rmeta.OffsetP + rmeta.Double32,
						Size:  4,
						EName: "Double32_t*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-pchar",
			ptr: &struct {
				N int32
				F []string
			}{3, []string{"s11", "s222", "s333"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerBasicPointer(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Type:  rmeta.OffsetP + rmeta.CharStar,
						Size:  4,
						EName: "char**",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-TString",
			ptr: &struct {
				N int32
				F []string
			}{3, []string{"s1", "s2", "s3"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerLoop(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Size:  4,
						EName: "TString*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-TObject",
			ptr: &struct {
				N int32
				F []rbase.Object
			}{3, []rbase.Object{*rbase.NewObject(), *rbase.NewObject(), *rbase.NewObject()}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerLoop(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Size:  4,
						EName: "TObject*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-TNamed",
			ptr: &struct {
				N int32
				F []rbase.Named
			}{3, []rbase.Named{*rbase.NewNamed("s1", "t1"), *rbase.NewNamed("s2", "t2"), *rbase.NewNamed("s3", "t3")}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerLoop(Element{
						Name:  *rbase.NewNamed("F", "[N]"),
						Size:  4,
						EName: "TNamed*",
					}.New(), 1, "N", "T"),
				},
			},
		},
		{
			name: "sli-Pos",
			ptr: func() interface{} {
				type Pos struct {
					Px float32
					Py float64
				}
				type T struct {
					N  int32
					Ps []Pos `groot:"Ps[N]"`
				}
				return &T{
					N:  2,
					Ps: []Pos{{1, 2}, {3, 4}},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBasicType{Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Int32,
						Size:  4,
						EName: "int32",
					}.New()},
					NewStreamerLoop(Element{
						Name:  *rbase.NewNamed("Ps", "[N]"),
						Size:  4,
						EName: "Pos*",
					}.New(), 1, "N", "T"),
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Pos", "Pos"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("Px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("Py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "std-string",
			ptr: &struct {
				F string
			}{"hello"},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerSTLstring{
						StreamerSTL: StreamerSTL{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("F", ""),
								Type:   rmeta.Streamer,
								Size:   32,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "string",
							}.New(),
							vtype: rmeta.ESTLType(rmeta.STLstring),
							ctype: rmeta.STLstring,
						},
					},
				},
			},
		},
		{
			name: "std::vector<bool>",
			ptr: &struct {
				F []bool
			}{[]bool{true, false, true}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<bool>",
					}.New(), rmeta.STLvector, rmeta.Bool),
				},
			},
		},
		{
			name: "std::vector<int8>",
			ptr: &struct {
				F []int8
			}{[]int8{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<int8>",
					}.New(), rmeta.STLvector, rmeta.Int8),
				},
			},
		},
		{
			name: "std::vector<int16>",
			ptr: &struct {
				F []int16
			}{[]int16{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<int16>",
					}.New(), rmeta.STLvector, rmeta.Int16),
				},
			},
		},
		{
			name: "std::vector<int32>",
			ptr: &struct {
				F []int32
			}{[]int32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<int32>",
					}.New(), rmeta.STLvector, rmeta.Int32),
				},
			},
		},
		{
			name: "std::vector<int64>",
			ptr: &struct {
				F []int64
			}{[]int64{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<int64>",
					}.New(), rmeta.STLvector, rmeta.Int64),
				},
			},
		},
		{
			name: "std::vector<uint8>",
			ptr: &struct {
				F []uint8
			}{[]uint8{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<uint8>",
					}.New(), rmeta.STLvector, rmeta.Uint8),
				},
			},
		},
		{
			name: "std::vector<uint16>",
			ptr: &struct {
				F []uint16
			}{[]uint16{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<uint16>",
					}.New(), rmeta.STLvector, rmeta.Uint16),
				},
			},
		},
		{
			name: "std::vector<uint32>",
			ptr: &struct {
				F []uint32
			}{[]uint32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<uint32>",
					}.New(), rmeta.STLvector, rmeta.Uint32),
				},
			},
		},
		{
			name: "std::vector<uint64>",
			ptr: &struct {
				F []uint64
			}{[]uint64{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<uint64>",
					}.New(), rmeta.STLvector, rmeta.Uint64),
				},
			},
		},
		{
			name: "std::vector<float32>",
			ptr: &struct {
				F []float32
			}{[]float32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<float32>",
					}.New(), rmeta.STLvector, rmeta.Float32),
				},
			},
		},
		{
			name: "std::vector<float64>",
			ptr: &struct {
				F []float64
			}{[]float64{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<float64>",
					}.New(), rmeta.STLvector, rmeta.Float64),
				},
			},
		},
		{
			name: "std::vector<TString>",
			ptr: &struct {
				F []string
			}{[]string{"hello", "world", "Go-HEP"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<TString>",
					}.New(), rmeta.STLvector, rmeta.TString),
				},
			},
		},
		{
			name: "std::vector<TObject>",
			ptr: &struct {
				F []rbase.Object
			}{[]rbase.Object{*rbase.NewObject(), *rbase.NewObject(), *rbase.NewObject()}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<TObject>",
					}.New(), rmeta.STLvector, rmeta.TObject),
				},
			},
		},
		{
			name: "std::vector<TNamed>",
			ptr: &struct {
				F []rbase.Named
			}{[]rbase.Named{*rbase.NewNamed("v1", "t1"), *rbase.NewNamed("v2", "t2"), *rbase.NewNamed("v3", "t3")}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<TNamed>",
					}.New(), rmeta.STLvector, rmeta.TNamed),
				},
			},
		},
		{
			name: "std::vector<string>",
			ptr: &struct {
				F []string
			}{[]string{"hello", "world", "Go-HEP"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<TString>",
					}.New(), rmeta.STLvector, rmeta.STLstring),
				},
			},
		},
		{
			name: "particle",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
				}
				return &Particle{Pos: P2{42, 66}, Name: "HEP"}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("pos", ""),
							Type:   rmeta.Any,
							Size:   4 + 8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "P2",
						}.New(),
					},
					&StreamerString{Element{
						Name:  *rbase.NewNamed("name", ""),
						Type:  rmeta.TString,
						Size:  24,
						EName: "TString",
					}.New()},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "event-objstring",
			ptr: func() interface{} {
				type Particle struct {
					Name rbase.ObjString
				}
				type T struct {
					P Particle
				}
				return &T{
					P: Particle{
						Name: *rbase.NewObjString("part-1"),
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("P", ""),
							Type:   rmeta.Any,
							Size:   40,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "Particle",
						}.New(),
					},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObject{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("Name", ""),
								Type:   rmeta.Object,
								Size:   40,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "TObjString",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "event-particle",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
				}
				type T struct {
					Particle Particle `groot:"particle"`
				}
				return &T{
					Particle: Particle{Pos: P2{142, 166}, Name: "HEP-1"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAny{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("particle", ""),
							Type:   rmeta.Any,
							Size:   4 + 8 + 24,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "Particle",
						}.New(),
					},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObjectAny{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("pos", ""),
								Type:   rmeta.Any,
								Size:   4 + 8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "P2",
							}.New(),
						},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "std::vector<P2>",
			ptr: func() interface{} {
				type P2 struct {
					Px float32
					Py float64
				}
				type T struct {
					F []P2
				}
				return &T{
					F: []P2{{1, 2}, {3, 4}},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<P2>",
					}.New(), rmeta.STLvector, rmeta.Any),
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "event-std::vector<particle>",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  []P2   `groot:"pos"`
					Name string `groot:"name"`
				}
				type T struct {
					Particles []Particle `groot:"particles"`
				}
				return &T{
					Particles: []Particle{
						{Pos: []P2{{142, 166}, {143, 167}}, Name: "HEP-1"},
						{Pos: []P2{{242, 266}, {}, {243, 267}}, Name: "HEP-2"},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("particles", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<Particle>",
					}.New(), rmeta.STLvector, rmeta.Any),
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						NewCxxStreamerSTL(Element{
							Name:   *rbase.NewNamed("pos", ""),
							Type:   rmeta.Streamer,
							Size:   24,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "vector<P2>",
						}.New(), rmeta.STLvector, rmeta.Any),
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "base-object",
			ptr: func() interface{} {
				type T struct {
					rbase.Object `groot:",base"`
					F1           float64
				}
				return &T{Object: *rbase.NewObject(), F1: 42}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBase{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("TObject", ""),
							Type:   rmeta.Base,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "BASE",
						}.New(),
						vbase: rvers.Named,
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F1", ""),
							Type:   rmeta.Float64,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "float64",
						}.New(),
					},
				},
			},
		},
		{
			name: "base-named",
			ptr: func() interface{} {
				type T struct {
					rbase.Named `groot:",base"`
					F1          float64
				}
				return &T{Named: *rbase.NewNamed("n1", "t1"), F1: 42}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBase{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("TNamed", ""),
							Type:   rmeta.Base,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "BASE",
						}.New(),
						vbase: rvers.Named,
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F1", ""),
							Type:   rmeta.Float64,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "float64",
						}.New(),
					},
				},
			},
		},
		{
			name: "base-objstring",
			ptr: func() interface{} {
				type T struct {
					rbase.ObjString `groot:",base"`
					F1              float64
				}
				return &T{ObjString: *rbase.NewObjString("hello-obj-string"), F1: 42}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerBase{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("TObjString", ""),
							Type:   rmeta.Base,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "BASE",
						}.New(),
						vbase: rvers.ObjString,
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F1", ""),
							Type:   rmeta.Float64,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "float64",
						}.New(),
					},
				},
			},
		},
		{
			name: "ptr-to-object",
			ptr: func() interface{} {
				type T struct {
					F1 *rbase.Named
					F2 *rbase.Named
					F3 *rbase.Named
				}
				f1 := rbase.NewNamed("n1", "t1")
				f3 := f1
				return &T{F1: f1, F2: nil, F3: f3}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectPointer{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F1", ""),
							Type:   rmeta.ObjectP,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "TNamed*",
						}.New(),
					},
					&StreamerObjectPointer{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F2", ""),
							Type:   rmeta.ObjectP,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "TNamed*",
						}.New(),
					},
					&StreamerObjectPointer{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("F3", ""),
							Type:   rmeta.ObjectP,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "TNamed*",
						}.New(),
					},
				},
			},
		},
		{
			name: "ptr-to-any",
			ptr: &PtrToAny_T{
				Pos: &rcont.ArrayD{Data: []float64{1, 2, 3}},
				Nil: nil,
			},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("Particle", "Particle"),
				clsver: int32(((*PtrToAny_T)(nil)).RVersion()),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerObjectAnyPointer{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("pos", ""),
							Type:   rmeta.AnyP,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "TArrayD*",
						}.New(),
					},
					&StreamerObjectAnyPointer{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("nil", ""),
							Type:   rmeta.AnyP,
							Size:   8,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "TArrayD*",
						}.New(),
					},
				},
			},
		},
		{
			name: "event-particles",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
				}
				type T struct {
					Particles []Particle `groot:"particles"`
				}
				return &T{
					Particles: []Particle{
						{Pos: P2{142, 166}, Name: "HEP-1"},
						{Pos: P2{242, 266}, Name: "HEP-2"},
						{Pos: P2{342, 366}, Name: "HEP-3"},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("particles", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<Particle>",
					}.New(), rmeta.STLvector, rmeta.Any),
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObjectAny{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("pos", ""),
								Type:   rmeta.Any,
								Size:   4 + 8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "P2",
							}.New(),
						},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
		{
			name: "event-particles-tags",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Tag struct {
					Name string `groot:"name"`
					Alg  string `groot:"alg"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
					Tags []Tag  `groot:"tags"`
				}
				type T struct {
					Particles []Particle `groot:"particles"`
				}
				return &T{
					Particles: []Particle{
						{Pos: P2{142, 166}, Name: "HEP-1", Tags: []Tag{{"Tag11", "Alg11"}, {"Tag12", "Alg12"}}},
						{Pos: P2{242, 266}, Name: "HEP-2", Tags: []Tag{{"Tag21", "Alg21"}, {"Tag22", "Alg22"}}},
						{Pos: P2{342, 366}, Name: "HEP-3", Tags: []Tag{{"Tag31", "Alg31"}, {"Tag32", "Alg32"}}},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("particles", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<Particle>",
					}.New(), rmeta.STLvector, rmeta.Any),
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObjectAny{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("pos", ""),
								Type:   rmeta.Any,
								Size:   4 + 8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "P2",
							}.New(),
						},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
						NewCxxStreamerSTL(Element{
							Name:   *rbase.NewNamed("tags", ""),
							Type:   rmeta.Streamer,
							Size:   24,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "vector<Tag>",
						}.New(), rmeta.STLvector, rmeta.Any),
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("Tag", "Tag"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("alg", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
			},
		},
		{
			name: "set<i32>",
			ptr: func() interface{} {
				type T struct {
					F []int32
				}
				return &T{
					F: []int32{1, 2, 3},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "set<int32>",
					}.New(), rmeta.STLset, rmeta.Object),
				},
			},
		},
		{
			name: "set<TString>",
			ptr: func() interface{} {
				type T struct {
					F []string
				}
				return &T{
					F: []string{"s1", "s22", "s333"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "set<TString>",
					}.New(), rmeta.STLset, rmeta.Object),
				},
			},
		},
		{
			name: "set<vector<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "set<vector<float> >",
					}.New(), rmeta.STLset, rmeta.Object),
				},
			},
		},
		{
			name: "set<set<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "set<set<float> >",
					}.New(), rmeta.STLset, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_set<i32>",
			ptr: func() interface{} {
				type T struct {
					F []int32
				}
				return &T{
					F: []int32{1, 2, 3},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_set<int32>",
					}.New(), rmeta.STLunorderedset, rmeta.Int32),
				},
			},
		},
		{
			name: "unordered_set<TString>",
			ptr: func() interface{} {
				type T struct {
					F []string
				}
				return &T{
					F: []string{"s1", "s22", "s333"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_set<TString>",
					}.New(), rmeta.STLunorderedset, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_set<vector<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_set<vector<float> >",
					}.New(), rmeta.STLunorderedset, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_set<unordered_set<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_set<unordered_set<float> >",
					}.New(), rmeta.STLunorderedset, rmeta.Object),
				},
			},
		},
		{
			name: "map<i32,i32>",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]int32 `groot:"m"`
				}
				return &T{
					Map: map[int32]int32{
						1: 10,
						2: 20,
						3: 30,
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "map<int32,int32>",
					}.New(), rmeta.STLmap, rmeta.Object),
				},
			},
		},
		{
			name: "map<i32,string>",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]string `groot:"m"`
				}
				return &T{
					Map: map[int32]string{
						1: "one",
						2: "two",
						3: "three",
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "map<int32,TString>",
					}.New(), rmeta.STLmap, rmeta.Object),
				},
			},
		},
		{
			name: "map<i32,map<i32,string> >",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]map[int32]string `groot:"m"`
				}
				return &T{
					Map: map[int32]map[int32]string{
						1: {
							1: "one",
							2: "two",
							3: "three",
						},
						2: {
							1: "un",
							2: "deux",
							3: "trois",
						},
						3: {
							1: "eins",
							2: "zwei",
							3: "drei",
						},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "map<int32,map<int32,string> >",
					}.New(), rmeta.STLmap, rmeta.Object),
				},
			},
		},
		{
			name: "map<i32,vector<string> >",
			ptr: func() interface{} {
				type T struct {
					Map map[int32][]string `groot:"m"`
				}
				return &T{
					Map: map[int32][]string{
						1: {"one", "un", "eins"},
						2: {"two", "deux", "zwei"},
						3: {"three", "trois", "drei"},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   48,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "map<int32,vector<string> >",
					}.New(), rmeta.STLmap, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_map<i32,i32>",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]int32 `groot:"m"`
				}
				return &T{
					Map: map[int32]int32{
						1: 10,
						2: 20,
						3: 30,
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_map<int32,int32>",
					}.New(), rmeta.STLunorderedmap, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_map<i32,string>",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]string `groot:"m"`
				}
				return &T{
					Map: map[int32]string{
						1: "one",
						2: "two",
						3: "three",
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_map<int32,TString>",
					}.New(), rmeta.STLunorderedmap, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_map<i32,unordered_map<i32,string> >",
			ptr: func() interface{} {
				type T struct {
					Map map[int32]map[int32]string `groot:"m"`
				}
				return &T{
					Map: map[int32]map[int32]string{
						1: {
							1: "one",
							2: "two",
							3: "three",
						},
						2: {
							1: "un",
							2: "deux",
							3: "trois",
						},
						3: {
							1: "eins",
							2: "zwei",
							3: "drei",
						},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_map<int32,unordered_map<int32,string> >",
					}.New(), rmeta.STLunorderedmap, rmeta.Object),
				},
			},
		},
		{
			name: "unordered_map<i32,vector<string> >",
			ptr: func() interface{} {
				type T struct {
					Map map[int32][]string `groot:"m"`
				}
				return &T{
					Map: map[int32][]string{
						1: {"one", "un", "eins"},
						2: {"two", "deux", "zwei"},
						3: {"three", "trois", "drei"},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("m", ""),
						Type:   rmeta.Streamer,
						Size:   56,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "unordered_map<int32,vector<string> >",
					}.New(), rmeta.STLunorderedmap, rmeta.Object),
				},
			},
		},
		{
			name: "list<i32>",
			ptr: func() interface{} {
				type T struct {
					F []int32
				}
				return &T{
					F: []int32{1, 2, 3},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "list<int32>",
					}.New(), rmeta.STLlist, rmeta.Int32),
				},
			},
		},
		{
			name: "list<TString>",
			ptr: func() interface{} {
				type T struct {
					F []string
				}
				return &T{
					F: []string{"s1", "s22", "s333"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "list<TString>",
					}.New(), rmeta.STLlist, rmeta.Object),
				},
			},
		},
		{
			name: "list<vector<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "list<vector<float> >",
					}.New(), rmeta.STLlist, rmeta.Object),
				},
			},
		},
		{
			name: "list<list<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   24,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "list<list<float> >",
					}.New(), rmeta.STLlist, rmeta.Object),
				},
			},
		},
		{
			name: "deque<i32>",
			ptr: func() interface{} {
				type T struct {
					F []int32
				}
				return &T{
					F: []int32{1, 2, 3},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   80,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "deque<int32>",
					}.New(), rmeta.STLdeque, rmeta.Int32),
				},
			},
		},
		{
			name: "deque<TString>",
			ptr: func() interface{} {
				type T struct {
					F []string
				}
				return &T{
					F: []string{"s1", "s22", "s333"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   80,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "deque<TString>",
					}.New(), rmeta.STLdeque, rmeta.Object),
				},
			},
		},
		{
			name: "deque<vector<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   80,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "deque<vector<float> >",
					}.New(), rmeta.STLdeque, rmeta.Object),
				},
			},
		},
		{
			name: "deque<deque<float> >",
			ptr: func() interface{} {
				type T struct {
					F [][]float32
				}
				return &T{
					F: [][]float32{
						{1, 2, 3, 4},
						{5, 6},
						{7, 8, 9, 10},
					},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("F", ""),
						Type:   rmeta.Streamer,
						Size:   80,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "deque<deque<float> >",
					}.New(), rmeta.STLdeque, rmeta.Object),
				},
			},
		},
		{
			name: "rmeta-stl-string",
			ptr: &struct {
				F string
			}{"Go-HEP"},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerSTLstring{
						StreamerSTL: *NewCxxStreamerSTL(Element{
							Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
							Type:   rmeta.STL,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "vector<string>",
						}.New(), rmeta.ESTLType(rmeta.STLstring), rmeta.STLstring),
					},
				},
			},
		},
		{
			name: "rmeta-stl-vector<float>",
			ptr: &struct {
				F []float32
			}{[]float32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<float>",
					}.New(), rmeta.STLvector, rmeta.Float32),
				},
			},
		},
		{
			name: "rmeta-stl-set<float>",
			ptr: &struct {
				F []float32
			}{[]float32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "set<float>",
					}.New(), rmeta.STLset, rmeta.Float32),
				},
			},
		},
		{
			name: "rmeta-stl-list<float>",
			ptr: &struct {
				F []float32
			}{[]float32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "list<float>",
					}.New(), rmeta.STLlist, rmeta.Float32),
				},
			},
		},
		{
			name: "rmeta-stl-deque<float>",
			ptr: &struct {
				F []float32
			}{[]float32{1, 2, 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "deque<float>",
					}.New(), rmeta.STLdeque, rmeta.Float32),
				},
			},
		},
		{
			name: "rmeta-stl-map<int,float>",
			ptr: &struct {
				F map[int32]float32
			}{map[int32]float32{1: 1, 2: 2, 3: 3}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "map<int,float>",
					}.New(), rmeta.STLmap, rmeta.Object),
				},
			},
		},
		{
			name: "rmeta-stl-vector<string>",
			ptr: &struct {
				F []string
			}{[]string{"hello", "world", "Go-HEP"}},
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					NewCxxStreamerSTL(Element{
						Name:   *rbase.NewNamed("This", "Used to call the proper TStreamerInfo case"),
						Type:   rmeta.STL,
						MaxIdx: [5]int32{0, 0, 0, 0, 0},
						EName:  "vector<string>",
					}.New(), rmeta.STLvector, rmeta.STLstring),
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipping %s", tc.name)
			}

			for _, dep := range tc.deps {
				StreamerInfos.Add(dep)
			}
			defer func() {
				StreamerInfos.Lock()
				defer StreamerInfos.Unlock()
				for _, dep := range tc.deps {
					key := streamerDbKey{
						class:   dep.Name(),
						version: dep.ClassVersion(),
					}
					delete(StreamerInfos.db, key)
				}
			}()

			err := tc.si.BuildStreamers()
			if err != nil {
				t.Fatalf("could not build streamers: %+v", err)
			}

			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			enc, err := tc.si.NewEncoder(kind, wbuf)
			if err != nil {
				t.Fatalf("could not create encoder: %+v", err)
			}

			err = enc.EncodeROOT(tc.ptr)
			if err != nil {
				t.Fatalf("could not encode value: %+v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			dec, err := tc.si.NewDecoder(kind, rbuf)
			if err != nil {
				t.Fatalf("could not create decoder: %+v", err)
			}

			rv := reflect.New(reflect.TypeOf(tc.ptr).Elem()).Elem()
			got := rv.Addr().Interface()
			err = dec.DecodeROOT(got)
			if err != nil {
				t.Fatalf("could not decode value: %+v", err)
			}

			if got, want := got, tc.ptr; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid round-trip:\ngot= %#v\nwant=%#v", got, want)
			}
		})
	}
}

func TestRWStreamerInfo(t *testing.T) {
	const kind = rbytes.ObjectWise // FIXME(sbinet): also test MemberWise.

	for _, tc := range []struct {
		name string
		skip bool
		si   *StreamerInfo
		ptr  interface{}
		deps []rbytes.StreamerInfo
		err  error
	}{
		{
			name: "event",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
				}
				type T struct {
					Name     string   `groot:"name"`
					Particle Particle `groot:"particle"`
				}
				return &T{
					Name:     "Go-HEP",
					Particle: Particle{Pos: P2{142, 166}, Name: "HEP-1"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{Element{
						Name:  *rbase.NewNamed("name", ""),
						Type:  rmeta.TString,
						Size:  24,
						EName: "TString",
					}.New()},
					&StreamerObjectAny{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("particle", ""),
							Type:   rmeta.Any,
							Size:   4 + 8 + 24,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "Particle",
						}.New(),
					},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObjectAny{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("pos", ""),
								Type:   rmeta.Any,
								Size:   4 + 8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "P2",
							}.New(),
						},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipping %s", tc.name)
			}

			for _, dep := range tc.deps {
				StreamerInfos.Add(dep)
			}
			defer func() {
				StreamerInfos.Lock()
				defer StreamerInfos.Unlock()
				for _, dep := range tc.deps {
					key := streamerDbKey{
						class:   dep.Name(),
						version: dep.ClassVersion(),
					}
					delete(StreamerInfos.db, key)
				}
			}()

			err := tc.si.BuildStreamers()
			if err != nil {
				t.Fatalf("could not build streamers: %+v", err)
			}

			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)

			w, err := tc.si.NewWStreamer(kind)
			if err != nil {
				t.Fatalf("could not create write-streamer: %+v", err)
			}
			err = w.(rbytes.Binder).Bind(tc.ptr)
			if err != nil {
				t.Fatalf("could not bind write-streamer: %+v", err)
			}

			err = w.WStreamROOT(wbuf)
			if err != nil {
				t.Fatalf("could not write value: %+v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			r, err := tc.si.NewRStreamer(kind)
			if err != nil {
				t.Fatalf("could not create read-streamer: %+v", err)
			}

			rv := reflect.New(reflect.TypeOf(tc.ptr).Elem()).Elem()
			got := rv.Addr().Interface()

			err = r.(rbytes.Binder).Bind(got)
			if err != nil {
				t.Fatalf("could not bind read-streamer: %+v", err)
			}

			err = r.RStreamROOT(rbuf)
			if err != nil {
				t.Fatalf("could not read value: %+v", err)
			}

			if got, want := got, tc.ptr; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid round-trip:\ngot= %#v\nwant=%#v", got, want)
			}
		})
	}
}

func TestRWStreamerElem(t *testing.T) {
	const kind = rbytes.ObjectWise // FIXME(sbinet): also test MemberWise.

	for _, tc := range []struct {
		name string
		skip bool
		si   *StreamerInfo
		ptr  interface{}
		deps []rbytes.StreamerInfo
		err  error
	}{
		{
			name: "event",
			ptr: func() interface{} {
				type P2 struct {
					Px float32 `groot:"px"`
					Py float64 `groot:"py"`
				}
				type Particle struct {
					Pos  P2     `groot:"pos"`
					Name string `groot:"name"`
				}

				type T struct {
					Name string `groot:"name"`

					B   bool    `groot:"b"`
					I8  int8    `groot:"i8"`
					I16 int16   `groot:"i16"`
					I32 int32   `groot:"i32"`
					I64 int64   `groot:"i64"`
					U8  uint8   `groot:"u8"`
					U16 uint16  `groot:"u16"`
					U32 uint32  `groot:"u32"`
					U64 uint64  `groot:"u64"`
					F32 float32 `groot:"f32"`
					F64 float64 `groot:"f64"`

					ArrB   [3]bool    `groot:"arrB[3]"`
					ArrI8  [3]int8    `groot:"arrU8[3]"`
					ArrI16 [3]int16   `groot:"arrU16[3]"`
					ArrI32 [3]int32   `groot:"arrU32[3]"`
					ArrI64 [3]int64   `groot:"arrU64[3]"`
					ArrU8  [3]uint8   `groot:"arrU8[3]"`
					ArrU16 [3]uint16  `groot:"arrU16[3]"`
					ArrU32 [3]uint32  `groot:"arrU32[3]"`
					ArrU64 [3]uint64  `groot:"arrU64[3]"`
					ArrF32 [3]float32 `groot:"arrF32[3]"`
					ArrF64 [3]float64 `groot:"arrF64[3]"`

					N      int32     `groot:"N"`
					SliB   []bool    `groot:"sliB[N]"`
					SliI8  []int8    `groot:"sliU8[N]"`
					SliI16 []int16   `groot:"sliU16[N]"`
					SliI32 []int32   `groot:"sliU32[N]"`
					SliI64 []int64   `groot:"sliU64[N]"`
					SliU8  []uint8   `groot:"sliU8[N]"`
					SliU16 []uint16  `groot:"sliU16[N]"`
					SliU32 []uint32  `groot:"sliU32[N]"`
					SliU64 []uint64  `groot:"sliU64[N]"`
					SliF32 []float32 `groot:"sliF32[N]"`
					SliF64 []float64 `groot:"sliF64[N]"`

					Particle Particle `groot:"particle"`
				}
				return &T{
					Name: "Go-HEP",

					B:   true,
					I8:  -8,
					I16: -16,
					I32: -32,
					I64: -64,
					U8:  8,
					U16: 16,
					U32: 32,
					U64: 64,
					F32: -32,
					F64: -64,

					ArrB:   [3]bool{true, false, true},
					ArrI8:  [3]int8{-18, -28, -38},
					ArrI16: [3]int16{-18, -28, -38},
					ArrI32: [3]int32{-132, -232, -332},
					ArrI64: [3]int64{-164, -264, -364},
					ArrU8:  [3]uint8{18, 28, 38},
					ArrU16: [3]uint16{18, 28, 38},
					ArrU32: [3]uint32{132, 232, 332},
					ArrU64: [3]uint64{164, 264, 364},
					ArrF32: [3]float32{-132, -232, -332},
					ArrF64: [3]float64{-164, -264, -364},

					N:      2,
					SliB:   []bool{true, false},
					SliI8:  []int8{-18, -28},
					SliI16: []int16{-18, -28},
					SliI32: []int32{-132, -232},
					SliI64: []int64{-164, -264},
					SliU8:  []uint8{18, 28},
					SliU16: []uint16{18, 28},
					SliU32: []uint32{132, 232},
					SliU64: []uint64{164, 264},
					SliF32: []float32{-132, -232},
					SliF64: []float64{-164, -264},

					Particle: Particle{Pos: P2{142, 166}, Name: "HEP-1"},
				}
			}(),
			si: &StreamerInfo{
				named:  *rbase.NewNamed("T", "T"),
				objarr: rcont.NewObjArray(),
				elems: []rbytes.StreamerElement{
					&StreamerString{Element{
						Name:  *rbase.NewNamed("name", ""),
						Type:  rmeta.TString,
						Size:  24,
						EName: "TString",
					}.New()},

					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("b", ""),
							Type:  rmeta.Bool,
							EName: "bool",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("i8", ""),
							Type:  rmeta.Int8,
							EName: "int8_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("i16", ""),
							Type:  rmeta.Int16,
							EName: "int16_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("i32", ""),
							Type:  rmeta.Int32,
							EName: "int32_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("i64", ""),
							Type:  rmeta.Int64,
							EName: "int64_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("u8", ""),
							Type:  rmeta.Uint8,
							EName: "uint8_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("u16", ""),
							Type:  rmeta.Uint16,
							EName: "uint16_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("u32", ""),
							Type:  rmeta.Uint32,
							EName: "uint32_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("u64", ""),
							Type:  rmeta.Uint64,
							EName: "uint64_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("f32", ""),
							Type:  rmeta.Float32,
							EName: "float32_t",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("f64", ""),
							Type:  rmeta.Float64,
							EName: "float64_t",
						}.New(),
					},

					// arrays
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrB", ""),
							Type:   rmeta.OffsetL + rmeta.Bool,
							Size:   3 * 1,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "bool*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrI8", ""),
							Type:   rmeta.OffsetL + rmeta.Int8,
							Size:   3 * 1,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "int8_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrI16", ""),
							Type:   rmeta.OffsetL + rmeta.Int16,
							Size:   3 * 2,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "int16_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrI32", ""),
							Type:   rmeta.OffsetL + rmeta.Int32,
							Size:   3 * 4,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "int32_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrI64", ""),
							Type:   rmeta.OffsetL + rmeta.Int64,
							Size:   3 * 8,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "int64_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrU8", ""),
							Type:   rmeta.OffsetL + rmeta.Uint8,
							Size:   3 * 1,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "uint8_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrU16", ""),
							Type:   rmeta.OffsetL + rmeta.Uint16,
							Size:   3 * 2,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "uint16_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrU32", ""),
							Type:   rmeta.OffsetL + rmeta.Uint32,
							Size:   3 * 4,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "uint32_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrU64", ""),
							Type:   rmeta.OffsetL + rmeta.Uint64,
							Size:   3 * 8,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "uint64_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrF32", ""),
							Type:   rmeta.OffsetL + rmeta.Float32,
							Size:   3 * 4,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "float32_t*",
						}.New(),
					},
					&StreamerBasicType{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("arrF64", ""),
							Type:   rmeta.OffsetL + rmeta.Float64,
							Size:   3 * 8,
							ArrLen: 3,
							ArrDim: 1,
							MaxIdx: [5]int32{3, 0, 0, 0, 0},
							EName:  "float64_t*",
						}.New(),
					},

					// var-len arrays
					&StreamerBasicType{
						StreamerElement: Element{
							Name:  *rbase.NewNamed("N", ""),
							Type:  rmeta.Counter,
							Size:  4,
							EName: "int32_t",
						}.New(),
					},
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliB", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Bool,
							Size:  1,
							EName: "bool*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliI8", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int8,
							Size:  1,
							EName: "int8_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliI16", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int16,
							Size:  1,
							EName: "int16_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliI32", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int32,
							Size:  1,
							EName: "int32_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliI64", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Int64,
							Size:  1,
							EName: "int64_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliU8", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint8,
							Size:  1,
							EName: "uint8_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliU16", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint16,
							Size:  1,
							EName: "uint16_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliU32", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint32,
							Size:  1,
							EName: "uint32_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliU64", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Uint64,
							Size:  1,
							EName: "uint64_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliF32", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Float32,
							Size:  1,
							EName: "float32_t*",
						}.New(), 1, "N", "T",
					),
					NewStreamerBasicPointer(
						Element{
							Name:  *rbase.NewNamed("sliF64", "[N]"),
							Type:  rmeta.OffsetP + rmeta.Float64,
							Size:  1,
							EName: "float64_t*",
						}.New(), 1, "N", "T",
					),

					&StreamerObjectAny{
						StreamerElement: Element{
							Name:   *rbase.NewNamed("particle", ""),
							Type:   rmeta.Any,
							Size:   4 + 8 + 24,
							MaxIdx: [5]int32{0, 0, 0, 0, 0},
							EName:  "Particle",
						}.New(),
					},
				},
			},
			deps: []rbytes.StreamerInfo{
				&StreamerInfo{
					named:  *rbase.NewNamed("Particle", "Particle"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerObjectAny{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("pos", ""),
								Type:   rmeta.Any,
								Size:   4 + 8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "P2",
							}.New(),
						},
						&StreamerString{Element{
							Name:  *rbase.NewNamed("name", ""),
							Type:  rmeta.TString,
							Size:  24,
							EName: "TString",
						}.New()},
					},
				},
				&StreamerInfo{
					named:  *rbase.NewNamed("P2", "P2"),
					objarr: rcont.NewObjArray(),
					elems: []rbytes.StreamerElement{
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("px", ""),
								Type:   rmeta.Float32,
								Size:   4,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float32",
							}.New(),
						},
						&StreamerBasicType{
							StreamerElement: Element{
								Name:   *rbase.NewNamed("py", ""),
								Type:   rmeta.Float64,
								Size:   8,
								MaxIdx: [5]int32{0, 0, 0, 0, 0},
								EName:  "float64",
							}.New(),
						},
					},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skipf("skipping %s", tc.name)
			}

			for _, dep := range tc.deps {
				StreamerInfos.Add(dep)
			}
			defer func() {
				StreamerInfos.Lock()
				defer StreamerInfos.Unlock()
				for _, dep := range tc.deps {
					key := streamerDbKey{
						class:   dep.Name(),
						version: dep.ClassVersion(),
					}
					delete(StreamerInfos.db, key)
				}
			}()

			err := tc.si.BuildStreamers()
			if err != nil {
				t.Fatalf("could not build streamers: %+v", err)
			}

			for i, se := range tc.si.Elements() {
				t.Run(se.Name(), func(t *testing.T) {
					ptr := reflect.ValueOf(tc.ptr).Elem().Field(i).Addr().Interface()

					wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
					w, err := WStreamerOf(tc.si, i, kind)
					if err != nil {
						t.Fatalf("could not create write-streamer: %+v", err)
					}

					err = w.(rbytes.Binder).Bind(ptr)
					if err != nil {
						t.Fatalf("could not bind write-streamer: %+v", err)
					}

					if rv := reflect.ValueOf(ptr).Elem(); rv.Kind() == reflect.Slice {
						w.(*wstreamerElem).wop.cfg.count = func() int { return rv.Len() }
					}

					err = w.WStreamROOT(wbuf)
					if err != nil {
						t.Fatalf("could not write value: %+v", err)
					}

					rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
					r, err := RStreamerOf(tc.si, i, kind)
					if err != nil {
						t.Fatalf("could not create read-streamer: %+v", err)
					}

					rv := reflect.New(reflect.TypeOf(ptr).Elem()).Elem()
					got := rv.Addr().Interface()

					err = r.(rbytes.Binder).Bind(got)
					if err != nil {
						t.Fatalf("could not bind read-streamer: %+v", err)
					}

					if rv := reflect.ValueOf(ptr).Elem(); rv.Kind() == reflect.Slice {
						err = r.(rbytes.Counter).Count(func() int { return rv.Len() })
						if err != nil {
							t.Fatalf("could not set read-streamer counter: %+v", err)
						}
					}

					err = r.RStreamROOT(rbuf)
					if err != nil {
						t.Fatalf("could not read value: %+v", err)
					}

					if got, want := got, ptr; !reflect.DeepEqual(got, want) {
						t.Fatalf("invalid round-trip:\ngot= %#v\nwant=%#v", got, want)
					}
				})
			}
		})
	}
}

type PtrToAny_T struct {
	Pos *rcont.ArrayD `groot:"pos"`
	Nil *rcont.ArrayD `groot:"nil"`
}

func (*PtrToAny_T) RVersion() int16 { return 41 }
func (*PtrToAny_T) Class() string   { return "Particle" }

func (p *PtrToAny_T) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(p.RVersion())
	if err := w.WriteObjectAny(p.Pos); err != nil {
		return 0, err
	}
	if err := w.WriteObjectAny(p.Nil); err != nil {
		return 0, err
	}

	return w.SetByteCount(pos, p.Class())
}

func (p *PtrToAny_T) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(p.Class())
	if vers != p.RVersion() {
		return fmt.Errorf("invalid particle version: got=%d, want=%d", vers, p.RVersion())
	}

	p.Pos = nil
	if obj := r.ReadObjectAny(); obj != nil {
		p.Pos = obj.(*rcont.ArrayD)
	}

	p.Nil = nil
	if obj := r.ReadObjectAny(); obj != nil {
		p.Nil = obj.(*rcont.ArrayD)
	}

	r.CheckByteCount(pos, bcnt, beg, p.Class())
	return r.Err()
}

var (
	_ root.Object        = (*PtrToAny_T)(nil)
	_ rbytes.Marshaler   = (*PtrToAny_T)(nil)
	_ rbytes.Unmarshaler = (*PtrToAny_T)(nil)
)
