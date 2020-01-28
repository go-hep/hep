// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGoName2Cxx(t *testing.T) {
	for _, tc := range []struct {
		name, want string
	}{
		{
			name: "go-hep.org/x/hep/hbook.H1D",
			want: "go_hep_org::x::hep::hbook::H1D",
		},
		{
			name: "go-hep.org/x.H1D",
			want: "go_hep_org::x::H1D",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := GoName2Cxx(tc.name)
			if tc.want != got {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestTypename(t *testing.T) {
	for _, tc := range []struct {
		name  string
		title string
		want  string
		ok    bool
	}{
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;go-hep.org/x/hep/hbook.H1D",
			want:  "go-hep.org/x/hep/hbook.H1D",
			ok:    true,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "",
			want:  "go_hep_org::x::hep::hbook::H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go;hbook.H1D",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Go; hbook.H1D ",
			want:  "hbook.H1D",
			ok:    false,
		},
		{
			name:  "go-hep.org/x/hep/hbook.H1D",
			title: "Rust; stl::hbook::H1D",
			want:  "stl::hbook::H1D",
			ok:    false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			name := GoName2Cxx(tc.name)
			got, ok := Typename(name, tc.title)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
			if ok != tc.ok {
				t.Fatalf("got=%q, want=%q, ok=%v (want=%v)", got, tc.want, ok, tc.ok)
			}
		})
	}

	if _, ok := Typename("go_hep_org::x::hep::groot::redm::HLV", "Go;go-hep.org/x/hep/groot/redm.Event"); ok {
		t.Fatalf("typename did not fail!")
	}
}

func TestROOTComment(t *testing.T) {
	var g genGoType
	for _, tc := range []struct {
		title string
		meta  string
		doc   string
	}{
		{
			title: "A comment",
			meta:  "",
			doc:   "A comment",
		},
		{
			title: " A comment ",
			meta:  "",
			doc:   "A comment",
		},
		{
			title: "[N]",
			meta:  "[N]",
			doc:   "",
		},
		{
			title: "[N] this is an array. ",
			meta:  "[N]",
			doc:   "this is an array.",
		},
		{
			title: "[-1,1,2]",
			meta:  "[-1,1,2]",
			doc:   "",
		},
		{
			title: "[-1,1,2] a Double32 with min,max,factor",
			meta:  "[-1,1,2]",
			doc:   "a Double32 with min,max,factor",
		},
		{
			title: "[fN][-1,1,2] an array of Double32-s with min,max,factor",
			meta:  "[fN][-1,1,2]",
			doc:   "an array of Double32-s with min,max,factor",
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			meta, doc := g.rcomment(tc.title)
			if meta != tc.meta {
				t.Fatalf("meta: got=%q, want=%q", meta, tc.meta)
			}
			if doc != tc.doc {
				t.Fatalf("doc: got=%q, want=%q", doc, tc.doc)
			}
		})
	}
}

func TestGenCxxStreamerInfo(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "TObject",
			want: `NewCxxStreamerInfo("TObject", 1, 0x901bc02d, []rbytes.StreamerElement{
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fUniqueID", "object unique identifier"),
		Type:   rmeta.UInt,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "unsigned int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fBits", "bit field status word"),
		Type:   rmeta.Bits,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "unsigned int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
})`,
		},
		{
			name: "TObjString",
			want: `NewCxxStreamerInfo("TObjString", 1, 0x9c8e4800, []rbytes.StreamerElement{
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TObject", "Basic ROOT object"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, -1877229523, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1),
	&StreamerString{StreamerElement: Element{
		Name:   *rbase.NewNamed("fString", "wrapped TString"),
		Type:   rmeta.TString,
		Size:   24,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TString",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
})`,
		},
		{
			name: "TArray",
			want: `NewCxxStreamerInfo("TArray", 1, 0x7021b2, []rbytes.StreamerElement{
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fN", "Number of array elements"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
})`,
		},
		{
			name: "TArrayC",
			want: `NewCxxStreamerInfo("TArrayC", 1, 0xae879936, []rbytes.StreamerElement{
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TArray", "Abstract array base class"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 7348658, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1),
	NewStreamerBasicPointer(Element{
		Name:   *rbase.NewNamed("fArray", "[fN] Array of fN chars"),
		Type:   41,
		Size:   1,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "char*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1, "fN", "TArray"),
})`,
		},
		{
			name: "TTree",
			want: `NewCxxStreamerInfo("TTree", 20, 0x7264e07f, []rbytes.StreamerElement{
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TNamed", "The basis for a named object (name, title)"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, -541636036, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1),
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TAttLine", "Line attributes"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, -1811462839, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 2),
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TAttFill", "Fill area attributes"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, -2545006, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 2),
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TAttMarker", "Marker attributes"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 689802220, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 2),
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fEntries", "Number of entries"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fTotBytes", "Total number of bytes in all branches before compression"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fZipBytes", "Total number of bytes in all branches after compression"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fSavedBytes", "Number of autosaved bytes"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fFlushedBytes", "Number of auto-flushed bytes"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fWeight", "Tree weight (see TTree::SetWeight)"),
		Type:   rmeta.Double,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "double",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fTimerInterval", "Timer interval in milliseconds"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fScanField", "Number of runs before prompting in Scan"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fUpdate", "Update frequency for EntryLoop"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fDefaultEntryOffsetLen", "Initial Length of fEntryOffset table in the basket buffers"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fNClusterRange", "Number of Cluster range in addition to the one defined by 'AutoFlush'"),
		Type:   rmeta.Counter,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fMaxEntries", "Maximum number of entries in case of circular buffers"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fMaxEntryLoop", "Maximum number of entries to process"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fMaxVirtualSize", "Maximum total size of buffers kept in memory"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fAutoSave", "Autosave tree when fAutoSave entries written or -fAutoSave (compressed) bytes produced"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fAutoFlush", "Auto-flush tree when fAutoFlush entries written or -fAutoFlush (compressed) bytes produced"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fEstimate", "Number of entries to estimate histogram limits"),
		Type:   rmeta.Long64,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	NewStreamerBasicPointer(Element{
		Name:   *rbase.NewNamed("fClusterRangeEnd", "[fNClusterRange] Last entry of a cluster range."),
		Type:   56,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 20, "fNClusterRange", "TTree"),
	NewStreamerBasicPointer(Element{
		Name:   *rbase.NewNamed("fClusterSize", "[fNClusterRange] Number of entries in each cluster for a given range."),
		Type:   56,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "Long64_t*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 20, "fNClusterRange", "TTree"),
	&StreamerObjectAny{StreamerElement: Element{
		Name:   *rbase.NewNamed("fIOFeatures", "IO features to define for newly-written baskets and branches."),
		Type:   rmeta.Any,
		Size:   1,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "ROOT::TIOFeatures",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObject{StreamerElement: Element{
		Name:   *rbase.NewNamed("fBranches", "List of Branches"),
		Type:   rmeta.Object,
		Size:   64,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TObjArray",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObject{StreamerElement: Element{
		Name:   *rbase.NewNamed("fLeaves", "Direct pointers to individual branch leaves"),
		Type:   rmeta.Object,
		Size:   64,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TObjArray",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fAliases", "List of aliases for expressions based on the tree branches."),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TList*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectAny{StreamerElement: Element{
		Name:   *rbase.NewNamed("fIndexValues", "Sorted index values"),
		Type:   rmeta.Any,
		Size:   24,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TArrayD",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectAny{StreamerElement: Element{
		Name:   *rbase.NewNamed("fIndex", "Index of sorted values"),
		Type:   rmeta.Any,
		Size:   24,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TArrayI",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fTreeIndex", "Pointer to the tree Index (if any)"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TVirtualIndex*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fFriends", "pointer to list of friend elements"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TList*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fUserInfo", "pointer to a list of user objects associated to this Tree"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TList*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fBranchRef", "Branch supporting the TRefTable (if any)"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TBranchRef*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
})`,
		},
		{
			name: "TRefTable",
			want: `NewCxxStreamerInfo("TRefTable", 3, 0x8c895b85, []rbytes.StreamerElement{
	NewStreamerBase(Element{
		Name:   *rbase.NewNamed("TObject", "Basic ROOT object"),
		Type:   rmeta.Base,
		Size:   0,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, -1877229523, 0, 0, 0},
		Offset: 0,
		EName:  "BASE",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1),
	&StreamerBasicType{StreamerElement: Element{
		Name:   *rbase.NewNamed("fSize", "dummy for backward compatibility"),
		Type:   rmeta.Int,
		Size:   4,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "int",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fParents", "array of Parent objects  (eg TTree branch) holding the referenced objects"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TObjArray*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	&StreamerObjectPointer{StreamerElement: Element{
		Name:   *rbase.NewNamed("fOwner", "Object owning this TRefTable"),
		Type:   rmeta.ObjectP,
		Size:   8,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "TObject*",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New()},
	NewCxxStreamerSTL(Element{
		Name:   *rbase.NewNamed("fProcessGUIDs", "UUIDs of TProcessIDs used in fParentIDs"),
		Type:   rmeta.Streamer,
		Size:   24,
		ArrLen: 0,
		ArrDim: 0,
		MaxIdx: [5]int32{0, 0, 0, 0, 0},
		Offset: 0,
		EName:  "vector<string>",
		XMin:   0.000000,
		XMax:   0.000000,
		Factor: 0.000000,
	}.New(), 1, 61),
})`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new(strings.Builder)
			si, ok := StreamerInfos.Get(tc.name, -1)
			if !ok {
				t.Fatalf("could not get streamer for %q", tc.name)
			}
			err := GenCxxStreamerInfo(got, si, true)
			if err != nil {
				t.Fatalf("could not generate textual representation of %q: %+v", tc.name, err)
			}

			if got, want := got.String(), tc.want; got != want {
				diff := cmp.Diff(got, want)
				t.Fatalf("invalid streamer representation for %q:\n%s", tc.name, diff)
			}
		})
	}

}
