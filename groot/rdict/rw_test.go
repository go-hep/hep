// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rmeta"
	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestWRBuffer(t *testing.T) {
	for _, tc := range []struct {
		name string
		want rtests.ROOTer
	}{
		{
			name: "TStreamerBase",
			want: &StreamerBase{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("TAttLine", "Line attributes"),
					etype:  0,
					esize:  0,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "BASE",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				vbase: 1,
			},
		},
		{
			name: "TStreamerBasicType",
			want: &StreamerBasicType{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fEntries", "Number of entries"),
					etype:  16,
					esize:  8,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "Long64_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerBasicType",
			want: &StreamerBasicType{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fEntries", "Array of entries"),
					etype:  rmeta.OffsetL + rmeta.ULong,
					esize:  40,
					arrlen: 5,
					arrdim: 1,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "ULong_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerBasicPointer",
			want: &StreamerBasicPointer{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fClusterRangeEnd", "[fNClusterRange] Last entry of a cluster range."),
					etype:  56,
					esize:  8,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "Long64_t*",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
				cvers: 19,
				cname: "fNClusterRange",
				ccls:  "TTree",
			},
		},
		{
			name: "TStreamerBasicType",
			want: &StreamerBasicType{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fEntries", "DynArray of entries"),
					etype:  rmeta.OffsetP + rmeta.ULong,
					esize:  8,
					arrlen: 0,
					arrdim: 1,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "ULong_t",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerLoop",
			want: &StreamerLoop{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fLoop", "A streamer loop"),
				},
				cvers:  1,
				cname:  "fArrayCount",
				cclass: "MyArrayCount",
			},
		},
		{
			name: "TStreamerObject",
			want: &StreamerObject{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fBranches", "List of branches"),
					etype:  61,
					esize:  64,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "TObjArray",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerObjectAnyPointer",
			want: &StreamerObjectAnyPointer{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fObjAnyPtr", "A pointer to any object"),
				},
			},
		},
		{
			name: "TStreamerObjectAny",
			want: &StreamerObjectAny{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fIndexValues", "Sorted index values"),

					etype:  62,
					esize:  24,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "TArrayD",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerString",
			want: &StreamerString{
				StreamerElement: StreamerElement{
					named:  *rbase.NewNamed("fName", "object identifier"),
					etype:  65,
					esize:  24,
					arrlen: 0,
					arrdim: 0,
					maxidx: [5]int32{0, 0, 0, 0, 0},
					offset: 0,
					ename:  "TString",
					xmin:   0,
					xmax:   0,
					factor: 0,
				},
			},
		},
		{
			name: "TStreamerSTL",
			want: &StreamerSTL{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fStdSet", "A std::set<int>"),
					etype: rmeta.STL,
					ename: "std::set<int>",
				},
				vtype: rmeta.STLset,
				ctype: rmeta.STLset,
			},
		},
		{
			name: "TStreamerSTL",
			want: &StreamerSTL{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fStdMultimap", "A std::multimap<int,int>"),
					etype: rmeta.STL,
					ename: "std::multimap<int,int>",
				},
				vtype: rmeta.STLmultimap,
				ctype: rmeta.STLmultimap,
			},
		},
		{
			name: "TStreamerSTLstring",
			want: &StreamerSTLstring{
				StreamerSTL: StreamerSTL{
					StreamerElement: StreamerElement{
						named: *rbase.NewNamed("fStdString", "A std::string"),
						etype: rmeta.STL,
						ename: "std::string",
					},
					vtype: rmeta.STLstring,
					ctype: rmeta.STLstring,
				},
			},
		},
		{
			name: "TStreamerArtificial",
			want: &StreamerArtificial{
				StreamerElement: StreamerElement{
					named: *rbase.NewNamed("fArtificial", "An artificial streamer"),
					ename: "std::artificial",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			class := tc.want.Class()
			obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
			{
				rbuf.SetErr(io.EOF)
				err = obj.UnmarshalROOT(rbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
				rbuf.SetErr(nil)
			}
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}
