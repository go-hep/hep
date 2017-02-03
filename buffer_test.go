// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

func TestRBuffer(t *testing.T) {
	data := make([]byte, 32)
	r := NewRBuffer(data, nil, 0)

	if got, want := r.Len(), int64(32); got != want {
		t.Fatalf("got len=%v. want=%v", got, want)
	}
	start := r.Pos()
	if start != 0 {
		t.Fatalf("got start=%v. want=%v", start, 0)
	}

	_ = r.ReadI16()
	if r.Err() != nil {
		t.Fatalf("error reading int16: %v", r.Err())
	}

	pos := r.Pos()
	if pos != 2 {
		t.Fatalf("got pos=%v. want=%v", pos, 16)
	}

	pos = 0
	data = make([]byte, 2*(2+4+8))
	r = NewRBuffer(data, nil, 0)
	for _, n := range []int{2, 4, 8} {
		beg := r.Pos()
		if beg != pos {
			t.Errorf("pos[%d] error: got=%d, want=%d\n", n, beg, pos)
		}
		switch n {
		case 2:
			_ = r.ReadI16()
			_ = r.ReadU16()
		case 4:
			_ = r.ReadI32()
			_ = r.ReadU32()
		case 8:
			_ = r.ReadI64()
			_ = r.ReadU64()
		}
		end := r.Pos()
		pos += int64(2 * n)

		if got, want := end-beg, int64(2*n); got != want {
			t.Errorf("%d-bytes: got=%d. want=%d\n", n, got, want)
		}
	}
}

func TestReadRBuffer(t *testing.T) {
	for _, test := range []struct {
		name string
		file string
		want ROOTUnmarshaler
	}{
		{
			name: "TNamed",
			want: &tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "my-name", title: "my-title"},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-cmssw.dat",
			want: &tnamed{
				obj:  tobject{id: 0x0, bits: 0x3000000},
				name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-cmssw-2.dat",
			want: &tnamed{
				obj:  tobject{id: 0x0, bits: 0x3500000},
				name: "edmTriggerResults_TriggerResults__HLT.present", title: "edmTriggerResults_TriggerResults__HLT.present",
			},
		},
		{
			name: "TNamed",
			file: "testdata/tnamed-long-string.dat",
			want: &tnamed{
				obj:   tobject{id: 0x0, bits: 0x3000000},
				name:  strings.Repeat("*", 256),
				title: "my-title",
			},
		},
		{
			name: "TList",
			want: &tlist{
				name: "list-name",
				objs: []Object{
					&tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "n0", title: "t0"},
					&tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "n1", title: "t1"},
				},
			},
		},
		{
			name: "TObjArray",
			want: &objarray{
				obj:  tobject{id: 0x0, bits: 0x3000000},
				name: "my-objs",
				arr: []Object{
					&tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "n0", title: "t0"},
					&tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "n1", title: "t1"},
					&tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "n2", title: "t2"},
				},
				last: 2,
			},
		},
		{
			name: "TList",
			file: "testdata/tlist-tsi.dat",
			want: &tlist{
				name: "",
				objs: []Object{
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TTree", title: ""},
						chksum: 0xa2a28f2,
						clsver: 19,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TNamed", title: "The basis for a named object (name, title)"},
									etype:  67,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttLine", title: "Line attributes"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttFill", title: "Fill area attributes"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttMarker", title: "Marker attributes"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 2,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fEntries", title: "Number of entries"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fTotBytes", title: "Total number of bytes in all branches before compression"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fZipBytes", title: "Total number of bytes in all branches after compression"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fSavedBytes", title: "Number of autosaved bytes"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFlushedBytes", title: "Number of autoflushed bytes"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fWeight", title: "Tree weight (see TTree::SetWeight)"},
									etype:  8,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Double_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fTimerInterval", title: "Timer interval in milliseconds"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fScanField", title: "Number of runs before prompting in Scan"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fUpdate", title: "Update frequency for EntryLoop"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fDefaultEntryOffsetLen", title: "Initial Length of fEntryOffset table in the basket buffers"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fNClusterRange", title: "Number of Cluster range in addition to the one defined by 'AutoFlush'"},
									etype:  6,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaxEntries", title: "Maximum number of entries in case of circular buffers"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaxEntryLoop", title: "Maximum number of entries to process"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaxVirtualSize", title: "Maximum total size of buffers kept in memory"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fAutoSave", title: "Autosave tree when fAutoSave bytes produced"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fAutoFlush", title: "Autoflush tree when fAutoFlush entries written"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fEstimate", title: "Number of entries to estimate histogram limits"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fClusterRangeEnd", title: "[fNClusterRange] Last entry of a cluster range."},
									etype:  56,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t*",
								},
								cvers: 19,
								cname: "fNClusterRange",
								ccls:  "TTree",
							},
							&tstreamerBasicPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fClusterSize", title: "[fNClusterRange] Number of entries in each cluster for a given range."},
									etype:  56,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t*",
								},
								cvers: 19,
								cname: "fNClusterRange",
								ccls:  "TTree",
							},
							&tstreamerObject{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBranches", title: "List of Branches"},
									etype:  61,
									esize:  64,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray",
								},
							},
							&tstreamerObject{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLeaves", title: "Direct pointers to individual branch leaves"},
									etype:  61,
									esize:  64,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fAliases", title: "List of aliases for expressions based on the tree branches."},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TList*",
								},
							},
							&tstreamerObjectAny{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fIndexValues", title: "Sorted index values"},
									etype:  62,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TArrayD",
								},
							},
							&tstreamerObjectAny{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fIndex", title: "Index of sorted values"},
									etype:  62,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TArrayI",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fTreeIndex", title: "Pointer to the tree Index (if any)"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TVirtualIndex*",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFriends", title: "pointer to list of friend elements"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TList*",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fUserInfo", title: "pointer to a list of user objects associated to this Tree"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TList*",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBranchRef", title: "Branch supporting the TRefTable (if any)"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TBranchRef*",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TNamed", title: ""},
						chksum: 0xfbe93f79,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TObject", title: "Basic ROOT object"},
									etype:  66,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerString{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fName", title: "object identifier"},
									etype:  65,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TString",
								},
							},
							&tstreamerString{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fTitle", title: "object title"},
									etype:  65,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TString",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TObject", title: ""},
						chksum: 0x52d96731,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fUniqueID", title: "object unique identifier"},
									etype:  13,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "UInt_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBits", title: "bit field status word"},
									etype:  15,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "UInt_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttLine", title: ""},
						chksum: 0x51a23e92,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLineColor", title: "line color"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLineStyle", title: "line style"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLineWidth", title: "line width"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttFill", title: ""},
						chksum: 0x47c56358,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFillColor", title: "fill area color"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFillStyle", title: "fill area style"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttMarker", title: ""},
						chksum: 0xfacd2184,
						clsver: 2,
						elems: []StreamerElement{
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMarkerColor", title: "Marker color index"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMarkerStyle", title: "Marker style"},
									etype:  2,
									esize:  2,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "short",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMarkerSize", title: "Marker size"},
									etype:  5,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "float",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TBranch", title: ""},
						chksum: 0x911cc38e,
						clsver: 12,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TNamed", title: "The basis for a named object (name, title)"},
									etype:  67,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TAttFill", title: "Fill area attributes"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fCompress", title: "Compression level and algorithm"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBasketSize", title: "Initial Size of  Basket Buffer"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fEntryOffsetLen", title: "Initial Length of fEntryOffset table in the basket buffers"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fWriteBasket", title: "Last basket number written"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fEntryNumber", title: "Current entry number (last one filled in this branch)"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fOffset", title: "Offset of this branch"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaxBaskets", title: "Maximum number of Baskets so far"},
									etype:  6,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fSplitLevel", title: "Branch split level"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fEntries", title: "Number of entries"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFirstEntry", title: "Number of the first entry in this branch"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fTotBytes", title: "Total number of bytes in all leaves before compression"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fZipBytes", title: "Total number of bytes in all leaves after compression"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerObject{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBranches", title: "-> List of Branches of this branch"},
									etype:  61,
									esize:  64,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray",
								},
							},
							&tstreamerObject{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLeaves", title: "-> List of leaves of this branch"},
									etype:  61,
									esize:  64,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray",
								},
							},
							&tstreamerObject{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBaskets", title: "-> List of baskets of this branch"},
									etype:  61,
									esize:  64,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray",
								},
							},
							&tstreamerBasicPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBasketBytes", title: "[fMaxBaskets] Length of baskets on file"},
									etype:  43,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t*",
								},
								cvers: 12,
								cname: "fMaxBaskets",
								ccls:  "TBranch",
							},
							&tstreamerBasicPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBasketEntry", title: "[fMaxBaskets] Table of first entry in eack basket"},
									etype:  56,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t*",
								},
								cvers: 12,
								cname: "fMaxBaskets",
								ccls:  "TBranch",
							},
							&tstreamerBasicPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fBasketSeek", title: "[fMaxBaskets] Addresses of baskets on file"},
									etype:  56,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t*",
								},
								cvers: 12,
								cname: "fMaxBaskets",
								ccls:  "TBranch",
							},
							&tstreamerString{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fFileName", title: "Name of file where buffers are stored (\"\" if in same file as Tree header)"},
									etype:  65,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TString",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeafI", title: ""},
						chksum: 0xd0548a75,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeaf", title: "Leaf: description of a Branch data type"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 2,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMinimum", title: "Minimum value if leaf range is specified"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaximum", title: "Maximum value if leaf range is specified"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeaf", title: ""},
						chksum: 0x2b643927,
						clsver: 2,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TNamed", title: "The basis for a named object (name, title)"},
									etype:  67,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLen", title: "Number of fixed length elements"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLenType", title: "Number of bytes for this data type"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fOffset", title: "Offset in ClonesArray object (if one)"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fIsRange", title: "(=kTRUE if leaf has a range, kFALSE otherwise)"},
									etype:  18,
									esize:  1,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Bool_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fIsUnsigned", title: "(=kTRUE if unsigned, kFALSE otherwise)"},
									etype:  18,
									esize:  1,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Bool_t",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLeafCount", title: "Pointer to Leaf count if variable length (we do not own the counter)"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TLeaf*",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeafL", title: ""},
						chksum: 0x74651570,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeaf", title: "Leaf: description of a Branch data type"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 2,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMinimum", title: "Minimum value if leaf range is specified"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaximum", title: "Maximum value if leaf range is specified"},
									etype:  16,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Long64_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeafF", title: ""},
						chksum: 0x51705bd0,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeaf", title: "Leaf: description of a Branch data type"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 2,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMinimum", title: "Minimum value if leaf range is specified"},
									etype:  5,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Float_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaximum", title: "Maximum value if leaf range is specified"},
									etype:  5,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Float_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeafD", title: ""},
						chksum: 0x9716fde,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TLeaf", title: "Leaf: description of a Branch data type"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 2,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMinimum", title: "Minimum value if leaf range is specified"},
									etype:  8,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Double_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fMaximum", title: "Maximum value if leaf range is specified"},
									etype:  8,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Double_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TList", title: ""},
						chksum: 0x79c882a7,
						clsver: 5,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TSeqCollection", title: "Sequenceable collection ABC"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 0,
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TSeqCollection", title: ""},
						chksum: 0xd79a0d4d,
						clsver: 0,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TCollection", title: "Collection abstract base class"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 3,
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TCollection", title: ""},
						chksum: 0x8fd14d5e,
						clsver: 3,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TObject", title: "Basic ROOT object"},
									etype:  66,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerString{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fName", title: "name of the collection"},
									etype:  65,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TString",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fSize", title: "number of elements in collection"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TString", title: ""},
						chksum: 0x17419,
						clsver: 2,
						elems:  nil,
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TBranchRef", title: ""},
						chksum: 0xae295353,
						clsver: 1,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TBranch", title: "Branch descriptor"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 12,
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fRefTable", title: "pointer to the TRefTable"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TRefTable*",
								},
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TRefTable", title: ""},
						chksum: 0xac58de3a,
						clsver: 3,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TObject", title: "Basic ROOT object"},
									etype:  66,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 1,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fSize", title: "dummy for backward compatibility"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fParents", title: "array of Parent objects  (eg TTree branch) holding the referenced objects"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObjArray*",
								},
							},
							&tstreamerObjectPointer{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fOwner", title: "Object owning this TRefTable"},
									etype:  64,
									esize:  8,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "TObject*",
								},
							},
							&tstreamerSTL{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fProcessGUIDs", title: "UUIDs of TProcessIDs used in fParentIDs"},
									etype:  500,
									esize:  24,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "vector<string>",
								},
								vtype: 1,
								ctype: 61,
							},
						},
					},
					&tstreamerInfo{
						named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TObjArray", title: ""},
						chksum: 0xf6eac680,
						clsver: 3,
						elems: []StreamerElement{
							&tstreamerBase{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "TSeqCollection", title: "Sequenceable collection ABC"},
									etype:  0,
									esize:  0,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "BASE",
								},
								vbase: 0,
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLowerBound", title: "Lower bound of the array"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
							&tstreamerBasicType{
								tstreamerElement: tstreamerElement{
									named:  tnamed{obj: tobject{id: 0x0, bits: 0x3000000}, name: "fLast", title: "Last element in array containing an object"},
									etype:  3,
									esize:  4,
									arrlen: 0,
									arrdim: 0,
									maxidx: [5]int32{0, 0, 0, 0, 0},
									ename:  "Int_t",
								},
							},
						},
					},
					&tlist{
						name: "listOfRules",
						objs: []Object{
							&tobjString{
								obj: tobject{id: 0x0, bits: 0x3000000},
								str: "type=read sourceClass=\"TTree\" targetClass=\"TTree\" version=\"[-16]\" source=\"\" target=\"fDefaultEntryOffsetLen\" code=\"{ fDefaultEntryOffsetLen = 1000; }\" ",
							},
							&tobjString{
								obj: tobject{id: 0x0, bits: 0x3000000},
								str: "type=read sourceClass=\"TTree\" targetClass=\"TTree\" version=\"[-18]\" source=\"\" target=\"fNClusterRange\" code=\"{ fNClusterRange = 0; }\" ",
							},
						},
					},
				},
			},
		},
		{
			name: "TArrayI",
			file: "testdata/tarrayi.dat",
			want: &ArrayI{Data: []int32{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayL64",
			file: "testdata/tarrayl64.dat",
			want: &ArrayL64{Data: []int64{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayF",
			file: "testdata/tarrayf.dat",
			want: &ArrayF{Data: []float32{0, 1, 2, 3, 4}},
		},
		{
			name: "TArrayD",
			file: "testdata/tarrayd.dat",
			want: &ArrayD{Data: []float64{0, 1, 2, 3, 4}},
		},
	} {
		test := test
		file := test.file
		if file == "" {
			file = "testdata/" + strings.ToLower(test.name) + ".dat"
		}
		t.Run("read-buffer="+test.name, func(t *testing.T) {
			testReadRBuffer(t, test.name, file, test.want)
		})
	}
}

func testReadRBuffer(t *testing.T, name, file string, want interface{}) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	r := NewRBuffer(data, nil, 0)
	obj := Factory.get(name)().Interface().(ROOTUnmarshaler)
	err = obj.UnmarshalROOT(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(obj, want) {
		t.Fatalf("error: %q\ngot= %+v\nwant=%+v\n", file, obj, want)
	}
}
