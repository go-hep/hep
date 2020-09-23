// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

func TestTypeFromSI(t *testing.T) {
	rdict.StreamerInfos.Add(rdict.NewCxxStreamerInfo("TypeFromSI_Pos1", 1, 0, []rbytes.StreamerElement{
		&rdict.StreamerBasicType{
			StreamerElement: rdict.Element{
				Name:   *rbase.NewNamed("px", ""),
				Type:   rmeta.Float32,
				Size:   4,
				MaxIdx: [5]int32{0, 0, 0, 0, 0},
				EName:  "float32",
			}.New(),
		},
		&rdict.StreamerBasicType{
			StreamerElement: rdict.Element{
				Name:   *rbase.NewNamed("py", ""),
				Type:   rmeta.Float64,
				Size:   8,
				MaxIdx: [5]int32{0, 0, 0, 0, 0},
				EName:  "float64",
			}.New(),
		},
	}))

	rdict.StreamerInfos.Add(rdict.NewCxxStreamerInfo("TypeFromSI_Pos2", 1, 0, []rbytes.StreamerElement{
		&rdict.StreamerBasicType{
			StreamerElement: rdict.Element{
				Name:   *rbase.NewNamed("px", ""),
				Type:   rmeta.Float32,
				Size:   4,
				MaxIdx: [5]int32{0, 0, 0, 0, 0},
				EName:  "float32",
			}.New(),
		},
		&rdict.StreamerBasicType{
			StreamerElement: rdict.Element{
				Name:   *rbase.NewNamed("py", ""),
				Type:   rmeta.Float64,
				Size:   8,
				MaxIdx: [5]int32{0, 0, 0, 0, 0},
				EName:  "float64",
			}.New(),
		},
	}))
	type TypeFromSI_Pos2 struct {
		Px float32 `groot:"px"`
		Py float64 `groot:"py"`
	}
	rtypes.Factory.Add("TypeFromSI_Pos2", func() reflect.Value {
		return reflect.ValueOf(&TypeFromSI_Pos2{})
	})

	for _, tc := range []struct {
		name string
		si   rbytes.StreamerInfo
		want reflect.Type
	}{
		{
			name: "TObject",
			si: func() rbytes.StreamerInfo {
				const name = "TObject"
				si, ok := rdict.StreamerInfos.Get(name, -1)
				if !ok {
					t.Fatalf("could not load streamer for %q", name)
				}
				return si
			}(),
			want: reflect.TypeOf((*rbase.Object)(nil)).Elem(),
		},
		{
			name: "TNamed",
			si: func() rbytes.StreamerInfo {
				const name = "TNamed"
				si, ok := rdict.StreamerInfos.Get(name, -1)
				if !ok {
					t.Fatalf("could not load streamer for %q", name)
				}
				return si
			}(),
			want: reflect.TypeOf((*rbase.Named)(nil)).Elem(),
		},
		{
			name: "TObjString",
			si: func() rbytes.StreamerInfo {
				const name = "TObjString"
				si, ok := rdict.StreamerInfos.Get(name, -1)
				if !ok {
					t.Fatalf("could not load streamer for %q", name)
				}
				return si
			}(),
			want: reflect.TypeOf((*rbase.ObjString)(nil)).Elem(),
		},
		{
			name: "TObjArray",
			si: func() rbytes.StreamerInfo {
				const name = "TObjArray"
				si, ok := rdict.StreamerInfos.Get(name, -1)
				if !ok {
					t.Fatalf("could not load streamer for %q", name)
				}
				return si
			}(),
			want: reflect.TypeOf((*rcont.ObjArray)(nil)).Elem(),
		},
		{
			name: "TArrayD",
			si: func() rbytes.StreamerInfo {
				const name = "TArrayD"
				si, ok := rdict.StreamerInfos.Get(name, -1)
				if !ok {
					t.Fatalf("could not load streamer for %q", name)
				}
				return si
			}(),
			want: reflect.TypeOf((*rcont.ArrayD)(nil)).Elem(),
		},
		{
			name: "MyArrayD",
			si: rdict.NewCxxStreamerInfo("MyArrayD", 1, 0, []rbytes.StreamerElement{
				&rdict.StreamerSTLstring{
					StreamerSTL: *rdict.NewCxxStreamerSTL(
						rdict.Element{
							Name:  *rbase.NewNamed("Name", ""),
							Type:  rmeta.Streamer,
							Size:  32,
							EName: "string",
						}.New(),
						rmeta.ESTLType(rmeta.STLstring),
						rmeta.STLstring,
					),
				},
				&rdict.StreamerObject{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("Array", ""),
						Type:  rmeta.Object,
						Size:  40,
						EName: "TArrayD",
					}.New(),
				},
			}),
			want: func() reflect.Type {
				return reflect.TypeOf((*struct {
					ROOT_Name  string       `groot:"Name"`
					ROOT_Array rcont.ArrayD `groot:"Array"`
				})(nil)).Elem()
			}(),
		},
		{
			name: "string-old-streamer-v2",
			si:   rdict.NewCxxStreamerInfo("string", 2, 0, nil),
			want: reflect.TypeOf(""),
		},
		{
			name: "vector<double>-old-streamer-v6",
			si: rdict.NewCxxStreamerInfo("vector<double>", 6, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("vector<double>", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "vector<double>",
				}.New(), rmeta.STLvector, rmeta.Float64),
			}),
			want: reflect.TypeOf((*[]float64)(nil)).Elem(),
		},
		{
			name: "pair<int,float>",
			si: rdict.NewCxxStreamerInfo("pair<int,float>", 1, 0, []rbytes.StreamerElement{
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("first", ""),
						Type:  rmeta.Int32,
						EName: "int32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("second", ""),
						Type:  rmeta.Float32,
						EName: "float32_t",
					}.New(),
				},
			}),
			want: reflect.TypeOf((*struct {
				ROOT_first  int32   `groot:"first"`
				ROOT_second float32 `groot:"second"`
			})(nil)).Elem(),
		},
		{
			name: "map<int,float>",
			si: rdict.NewCxxStreamerInfo("map<int,float>", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,float>",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32]float32)(nil)).Elem(),
		},
		{
			name: "map<int,string>",
			si: rdict.NewCxxStreamerInfo("map<int,string>", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,string>",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32]string)(nil)).Elem(),
		},
		{
			name: "map<int,TNamed>",
			si: rdict.NewCxxStreamerInfo("map<int,TNamed>", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,TNamed>",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32]rbase.Named)(nil)).Elem(),
		},
		{
			name: "map<TNamed,int>",
			si: rdict.NewCxxStreamerInfo("map<TNamed,int>", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<TNamed,int>",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[rbase.Named]int32)(nil)).Elem(),
		},
		{
			name: "map<int,vector<TNamed> >",
			si: rdict.NewCxxStreamerInfo("map<int,vector<TNamed> >", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,vector<TNamed> >",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32][]rbase.Named)(nil)).Elem(),
		},
		{
			name: "map<int,vector<string> >",
			si: rdict.NewCxxStreamerInfo("map<int,vector<string> >", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,vector<string> >",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32][]string)(nil)).Elem(),
		},
		{
			name: "map<int,map<int,vector<string> > >",
			si: rdict.NewCxxStreamerInfo("map<int,map<int,vector<string> > >", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("This", ""),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,map<int,vector<string> > >",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*map[int32]map[int32][]string)(nil)).Elem(),
		},
		{
			name: "map<int,string>-with-pair-dict",
			si: rdict.NewCxxStreamerInfo("MyEvent", 1, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("mapI32Str", " (pair<int,string>)"),
					Type:  rmeta.Streamer,
					Size:  48,
					EName: "map<int,string>",
				}.New(), rmeta.STLmap, rmeta.Object),
			}),
			want: reflect.TypeOf((*struct {
				ROOT_mapI32Str map[int32]string `groot:"mapI32Str"`
			})(nil)).Elem(),
		},
		{
			name: "vector<TObject>",
			si: rdict.NewCxxStreamerInfo("vector<TObject>", rvers.StreamerInfo, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("vector<TObject>", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "vector<TObject>",
				}.New(), rmeta.STLvector, rmeta.TObject),
			}),
			want: reflect.TypeOf((*[]rbase.Object)(nil)).Elem(),
		},
		{
			name: "vector<TNamed>",
			si: rdict.NewCxxStreamerInfo("vector<TNamed>", rvers.StreamerInfo, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("vector<TNamed>", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "vector<TNamed>",
				}.New(), rmeta.STLvector, rmeta.TNamed),
			}),
			want: reflect.TypeOf((*[]rbase.Named)(nil)).Elem(),
		},
		{
			name: "vector<TObjString>",
			si: rdict.NewCxxStreamerInfo("vector<TObjString>", rvers.StreamerInfo, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("vector<TObjString>", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "vector<TObjString>",
				}.New(), rmeta.STLvector, rmeta.Object),
			}),
			want: reflect.TypeOf((*[]rbase.ObjString)(nil)).Elem(),
		},
		{
			name: "bitset<256>",
			si: rdict.NewCxxStreamerInfo("MyBitset", rvers.StreamerInfo, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("bs", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "bitset<256>",
				}.New(), rmeta.STLbitset, 0),
			}),
			want: reflect.TypeOf((*struct {
				ROOT_bs []uint8 `groot:"bs"`
			})(nil)).Elem(),
		},
		{
			name: "vector<bitset<256> >",
			si: rdict.NewCxxStreamerInfo("MyBitset", rvers.StreamerInfo, 0, []rbytes.StreamerElement{
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:  *rbase.NewNamed("bs", ""),
					Type:  rmeta.Streamer,
					Size:  24,
					EName: "vector<bitset<256> >",
				}.New(), rmeta.STLvector, 0),
			}),
			want: reflect.TypeOf((*struct {
				ROOT_bs [][]uint8 `groot:"bs"`
			})(nil)).Elem(),
		},
		{
			name: "event",
			si: rdict.NewCxxStreamerInfo("event", 1, 0, []rbytes.StreamerElement{
				rdict.NewStreamerBase(
					rdict.Element{
						Name:  *rbase.NewNamed("TObject", ""),
						Type:  rmeta.Base,
						EName: "BASE",
					}.New(),
					rvers.Named,
				),
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("b", ""),
						Type:  rmeta.Bool,
						EName: "bool",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("i8", ""),
						Type:  rmeta.Int8,
						EName: "int8_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("i16", ""),
						Type:  rmeta.Int16,
						EName: "int16_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("i32", ""),
						Type:  rmeta.Int32,
						EName: "int32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("i64", ""),
						Type:  rmeta.Int64,
						EName: "int64_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("u8", ""),
						Type:  rmeta.Uint8,
						EName: "uint8_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("u16", ""),
						Type:  rmeta.Uint16,
						EName: "uint16_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("u32", ""),
						Type:  rmeta.Uint32,
						EName: "uint32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("u64", ""),
						Type:  rmeta.Uint64,
						EName: "uint64_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("N", ""),
						Type:  rmeta.Counter,
						Size:  4,
						EName: "int32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("NN", ""),
						Type:  rmeta.Counter,
						Size:  8,
						EName: "int64_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("bits", ""),
						Type:  rmeta.Bits,
						Size:  4,
						EName: "Bits_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("f32", ""),
						Type:  rmeta.Float32,
						EName: "float32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("f64", ""),
						Type:  rmeta.Float64,
						EName: "float64_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("f16", ""),
						Type:  rmeta.Float16,
						EName: "Float16_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("d32", ""),
						Type:  rmeta.Double32,
						EName: "Double32_t",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("str", ""),
						Type:  rmeta.TString,
						EName: "TString",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("cstr", ""),
						Type:  rmeta.CharStar,
						EName: "char*",
					}.New(),
				},

				// arrays
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrB", ""),
						Type:   rmeta.OffsetL + rmeta.Bool,
						Size:   3 * 1,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "bool*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrI8", ""),
						Type:   rmeta.OffsetL + rmeta.Int8,
						Size:   3 * 1,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int8_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrI16", ""),
						Type:   rmeta.OffsetL + rmeta.Int16,
						Size:   3 * 2,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int16_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrI32", ""),
						Type:   rmeta.OffsetL + rmeta.Int32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int32_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrI64", ""),
						Type:   rmeta.OffsetL + rmeta.Int64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "int64_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrU8", ""),
						Type:   rmeta.OffsetL + rmeta.Uint8,
						Size:   3 * 1,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint8_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrU16", ""),
						Type:   rmeta.OffsetL + rmeta.Uint16,
						Size:   3 * 2,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint16_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrU32", ""),
						Type:   rmeta.OffsetL + rmeta.Uint32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint32_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrU64", ""),
						Type:   rmeta.OffsetL + rmeta.Uint64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "uint64_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrF32", ""),
						Type:   rmeta.OffsetL + rmeta.Float32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "float32_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrF64", ""),
						Type:   rmeta.OffsetL + rmeta.Float64,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "float64_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrF16", ""),
						Type:   rmeta.OffsetL + rmeta.Float16,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "Float16_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrD32", ""),
						Type:   rmeta.OffsetL + rmeta.Double32,
						Size:   3 * 4,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "Double32_t*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrStr", ""),
						Type:   rmeta.OffsetL + rmeta.TString,
						Size:   3 * 24,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TString*",
					}.New(),
				},
				&rdict.StreamerBasicType{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrCstr", ""),
						Type:   rmeta.OffsetL + rmeta.CharStar,
						Size:   3 * 8,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "char**",
					}.New(),
				},
				//	&rdict.StreamerBasicType{
				//		StreamerElement: rdict.Element{
				//			Name:   *rbase.NewNamed("arrN32", ""),
				//			Type:   rmeta.OffsetL + rmeta.Counter,
				//			Size:   3 * 4,
				//			ArrLen: 3,
				//			ArrDim: 1,
				//		MaxIdx: [5]int32{1, 0, 0, 0, 0},
				//			EName:  "Counter_t",
				//		}.New(),
				//	},
				//	&rdict.StreamerBasicType{
				//		StreamerElement: rdict.Element{
				//			Name:   *rbase.NewNamed("arrN64", ""),
				//			Type:   rmeta.OffsetL + rmeta.Counter,
				//			Size:   3 * 8,
				//			ArrLen: 3,
				//			ArrDim: 1,
				//		MaxIdx: [5]int32{1, 0, 0, 0, 0},
				//			EName:  "Counter_t",
				//		}.New(),
				//	},

				&rdict.StreamerObjectPointer{
					StreamerElement: rdict.Element{
						Name:   *rbase.NewNamed("arrTObj", ""),
						Type:   rmeta.OffsetL + rmeta.TObject,
						Size:   3 * 16,
						ArrLen: 3,
						ArrDim: 1,
						MaxIdx: [5]int32{3, 0, 0, 0, 0},
						EName:  "TObject*",
					}.New(),
				},

				//	var-len arrays
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliB", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Bool,
						Size:  1,
						EName: "bool*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliI8", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int8,
						Size:  1,
						EName: "int8_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliI16", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int16,
						Size:  1,
						EName: "int16_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliI32", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int32,
						Size:  1,
						EName: "int32_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliI64", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Int64,
						Size:  1,
						EName: "int64_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliU8", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint8,
						Size:  1,
						EName: "uint8_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliU16", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint16,
						Size:  1,
						EName: "uint16_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliU32", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint32,
						Size:  1,
						EName: "uint32_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliU64", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Uint64,
						Size:  1,
						EName: "uint64_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliF32", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Float32,
						Size:  1,
						EName: "float32_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliF64", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Float64,
						Size:  1,
						EName: "float64_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliF16", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Float16,
						Size:  1,
						EName: "Float16_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliD32", "[N]"),
						Type:  rmeta.OffsetP + rmeta.Double32,
						Size:  1,
						EName: "Double32_t*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerLoop(
					rdict.Element{
						Name:  *rbase.NewNamed("sliStr", "[N]"),
						Size:  4,
						EName: "TString*",
					}.New(), 1, "N", "event",
				),
				rdict.NewStreamerBasicPointer(
					rdict.Element{
						Name:  *rbase.NewNamed("sliCstr", "[N]"),
						Type:  rmeta.OffsetP + rmeta.CharStar,
						Size:  1,
						EName: "char**",
					}.New(), 1, "N", "event",
				),

				// std-vector
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecB", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<bool>",
				}.New(), rmeta.STLvector, rmeta.Bool),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecI8", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<int8_t>",
				}.New(), rmeta.STLvector, rmeta.Int8),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecI16", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<int16_t>",
				}.New(), rmeta.STLvector, rmeta.Int16),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecI32", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<int32_t>",
				}.New(), rmeta.STLvector, rmeta.Int32),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecI64", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<int64_t>",
				}.New(), rmeta.STLvector, rmeta.Int64),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecU8", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<uint8_t>",
				}.New(), rmeta.STLvector, rmeta.Uint8),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecU16", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<uint16_t>",
				}.New(), rmeta.STLvector, rmeta.Uint16),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecU32", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<uint32_t>",
				}.New(), rmeta.STLvector, rmeta.Uint32),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecU64", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<uint64_t>",
				}.New(), rmeta.STLvector, rmeta.Uint64),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecF32", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<float32_t>",
				}.New(), rmeta.STLvector, rmeta.Float32),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecF64", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<float64_t>",
				}.New(), rmeta.STLvector, rmeta.Float64),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecF16", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<Float16_t>",
				}.New(), rmeta.STLvector, rmeta.Float16),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecD32", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<Double32_t>",
				}.New(), rmeta.STLvector, rmeta.Double32),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecStr", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<string>",
				}.New(), rmeta.STLvector, rmeta.STLstring),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecCstr", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<char*>",
				}.New(), rmeta.STLvector, rmeta.CharStar),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecNamed1", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<TNamed>",
				}.New(), rmeta.STLvector, rmeta.Object),
				rdict.NewCxxStreamerSTL(rdict.Element{
					Name:   *rbase.NewNamed("stdVecNamed2", ""),
					Type:   rmeta.Streamer,
					Size:   24,
					MaxIdx: [5]int32{0, 0, 0, 0, 0},
					EName:  "vector<TNamed>",
				}.New(), rmeta.STLvector, rmeta.TNamed),

				// obj-ptr
				&rdict.StreamerObjectPointer{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("ptrObj", ""),
						Type:  rmeta.ObjectP,
						Size:  8,
						EName: "TObject*",
					}.New(),
				},
				&rdict.StreamerObjectAnyPointer{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("ptrPos1", ""),
						Type:  rmeta.AnyP,
						Size:  8,
						EName: "TypeFromSI_Pos1*",
					}.New(),
				},
				&rdict.StreamerObjectAnyPointer{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("ptrPos2", ""),
						Type:  rmeta.AnyP,
						Size:  8,
						EName: "TypeFromSI_Pos2*",
					}.New(),
				},
				&rdict.StreamerObjectAnyPointer{
					StreamerElement: rdict.Element{
						Name:  *rbase.NewNamed("ptrArrF", ""),
						Type:  rmeta.AnyP,
						Size:  8,
						EName: "TArrayF*",
					}.New(),
				},
			}),
			want: reflect.TypeOf((*struct {
				ROOT_TObject rbase.Object  `groot:"TObject"`
				ROOT_b       bool          `groot:"b"`
				ROOT_i8      int8          `groot:"i8"`
				ROOT_i16     int16         `groot:"i16"`
				ROOT_i32     int32         `groot:"i32"`
				ROOT_i64     int64         `groot:"i64"`
				ROOT_u8      uint8         `groot:"u8"`
				ROOT_u16     uint16        `groot:"u16"`
				ROOT_u32     uint32        `groot:"u32"`
				ROOT_u64     uint64        `groot:"u64"`
				ROOT_N       int32         `groot:"N"`
				ROOT_NN      int64         `groot:"NN"`
				ROOT_bits    uint32        `groot:"bits"`
				ROOT_f32     float32       `groot:"f32"`
				ROOT_f64     float64       `groot:"f64"`
				ROOT_f16     root.Float16  `groot:"f16"`
				ROOT_d32     root.Double32 `groot:"d32"`
				ROOT_str     string        `groot:"str"`
				ROOT_cstr    string        `groot:"cstr"`
				// arrays
				ROOT_arrB    [3]bool          `groot:"arrB[3]"`
				ROOT_arrI8   [3]int8          `groot:"arrI8[3]"`
				ROOT_arrI16  [3]int16         `groot:"arrI16[3]"`
				ROOT_arrI32  [3]int32         `groot:"arrI32[3]"`
				ROOT_arrI64  [3]int64         `groot:"arrI64[3]"`
				ROOT_arrU8   [3]uint8         `groot:"arrU8[3]"`
				ROOT_arrU16  [3]uint16        `groot:"arrU16[3]"`
				ROOT_arrU32  [3]uint32        `groot:"arrU32[3]"`
				ROOT_arrU64  [3]uint64        `groot:"arrU64[3]"`
				ROOT_arrF32  [3]float32       `groot:"arrF32[3]"`
				ROOT_arrF64  [3]float64       `groot:"arrF64[3]"`
				ROOT_arrF16  [3]root.Float16  `groot:"arrF16[3]"`
				ROOT_arrD32  [3]root.Double32 `groot:"arrD32[3]"`
				ROOT_arrStr  [3]string        `groot:"arrStr[3]"`
				ROOT_arrCstr [3]string        `groot:"arrCstr[3]"`
				//	ROOT_arrN32  [3]int32         `groot:"arrN32[3]"`
				//	ROOT_arrN64  [3]int64         `groot:"arrN64[3]"`
				ROOT_arrTObj [3]rbase.Object `groot:"arrTObj[3]"`

				// slices
				ROOT_sliB    []bool          `groot:"sliB[N]"`
				ROOT_sliI8   []int8          `groot:"sliI8[N]"`
				ROOT_sliI16  []int16         `groot:"sliI16[N]"`
				ROOT_sliI32  []int32         `groot:"sliI32[N]"`
				ROOT_sliI64  []int64         `groot:"sliI64[N]"`
				ROOT_sliU8   []uint8         `groot:"sliU8[N]"`
				ROOT_sliU16  []uint16        `groot:"sliU16[N]"`
				ROOT_sliU32  []uint32        `groot:"sliU32[N]"`
				ROOT_sliU64  []uint64        `groot:"sliU64[N]"`
				ROOT_sliF32  []float32       `groot:"sliF32[N]"`
				ROOT_sliF64  []float64       `groot:"sliF64[N]"`
				ROOT_sliF16  []root.Float16  `groot:"sliF16[N]"`
				ROOT_sliD32  []root.Double32 `groot:"sliD32[N]"`
				ROOT_sliStr  []string        `groot:"sliStr[N]"`
				ROOT_sliCstr []string        `groot:"sliCstr[N]"`

				// std::vectors
				ROOT_stdVecB      []bool          `groot:"stdVecB"`
				ROOT_stdVecI8     []int8          `groot:"stdVecI8"`
				ROOT_stdVecI16    []int16         `groot:"stdVecI16"`
				ROOT_stdVecI32    []int32         `groot:"stdVecI32"`
				ROOT_stdVecI64    []int64         `groot:"stdVecI64"`
				ROOT_stdVecU8     []uint8         `groot:"stdVecU8"`
				ROOT_stdVecU16    []uint16        `groot:"stdVecU16"`
				ROOT_stdVecU32    []uint32        `groot:"stdVecU32"`
				ROOT_stdVecU64    []uint64        `groot:"stdVecU64"`
				ROOT_stdVecF32    []float32       `groot:"stdVecF32"`
				ROOT_stdVecF64    []float64       `groot:"stdVecF64"`
				ROOT_stdVecF16    []root.Float16  `groot:"stdVecF16"`
				ROOT_stdVecD32    []root.Double32 `groot:"stdVecD32"`
				ROOT_stdVecStr    []string        `groot:"stdVecStr"`
				ROOT_stdVecCstr   []string        `groot:"stdVecCstr"`
				ROOT_stdVecNamed1 []rbase.Named   `groot:"stdVecNamed1"`
				ROOT_stdVecNamed2 []rbase.Named   `groot:"stdVecNamed2"`

				// obj-ptr
				ROOT_ptrObj  *rbase.Object `groot:"ptrObj"`
				ROOT_ptrPos1 *struct {
					ROOT_px float32 `groot:"px"`
					ROOT_py float64 `groot:"py"`
				} `groot:"ptrPos1"`
				ROOT_ptrPos2 *TypeFromSI_Pos2 `groot:"ptrPos2"`
				ROOT_ptrArrF *rcont.ArrayF    `groot:"ptrArrF"`
			})(nil)).Elem(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.si.BuildStreamers()
			if err != nil {
				t.Fatalf("could not build streamers: %+v", err)
			}

			got, err := rdict.TypeFromSI(rdict.StreamerInfos, tc.si)
			if err != nil {
				t.Fatalf("could not load type: %+v", err)
			}

			if got != tc.want {
				t.Fatalf("invalid Go type:\ngot= %T\nwant=%T",
					reflect.New(got).Elem().Interface(),
					reflect.New(tc.want).Elem().Interface(),
				)
			}
		})
	}
}
