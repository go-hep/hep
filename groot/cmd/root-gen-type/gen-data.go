// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root = flag.String("f", "testdata/streamers.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("streamers", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run gen-streamers:\n%s\nerror: %+v", out, err)
	}
}

const script = `
#include <vector>
#include <string>

#include "TObjString.h"
#include "TString.h"
#include "Rtypes.h"

const int ARRAYSZ = 10;

struct P3 {
	int32_t Px;
	double  Py;
	int32_t Pz;
};

struct Event {
	TString  Beg;

	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	float    F32;
	double   F64;
	Float16_t  D16;
	Double32_t D32;
	TString  Str;

	::P3       P3;
	::P3      *P3Ptr;
	TObjString ObjStr;
	TObjString *ObjStrPtr;

	int16_t  ArrayI16[ARRAYSZ];
	int32_t  ArrayI32[ARRAYSZ];
	int64_t  ArrayI64[ARRAYSZ];
	uint16_t ArrayU16[ARRAYSZ];
	uint32_t ArrayU32[ARRAYSZ];
	uint64_t ArrayU64[ARRAYSZ];
	float    ArrayF32[ARRAYSZ];
	double   ArrayF64[ARRAYSZ];
	::P3     ArrayP3s[ARRAYSZ];
	TObjString ArrayObjStr[ARRAYSZ];

	int32_t  N;
	int16_t  *SliceI16;  //[N]
	int32_t  *SliceI32;  //[N]
	int64_t  *SliceI64;  //[N]
	uint16_t *SliceU16;  //[N]
	uint32_t *SliceU32;  //[N]
	uint64_t *SliceU64;  //[N]
	float    *SliceF32;  //[N]
	double   *SliceF64;  //[N]

	std::string StdStr;

	std::vector<int16_t> StlVecI16;
	std::vector<int32_t> StlVecI32;
	std::vector<int64_t> StlVecI64;
	std::vector<uint16_t> StlVecU16;
	std::vector<uint32_t> StlVecU32;
	std::vector<uint64_t> StlVecU64;
	std::vector<float> StlVecF32;
	std::vector<double> StlVecF64;
	std::vector<std::string> StlVecStr;

	TString End;
};

void streamers(const char* fname) {
	int evtmax = 1;

	auto f = TFile::Open(fname, "RECREATE");

	Event *evt = new Event;
	Event &e = *evt;

	for (int i = 0; i != evtmax; i++) {
		e.Beg = TString::Format("beg-%03d", i);
		e.I16 = i;
		e.I32 = i;
		e.I64 = i;
		e.U16 = i;
		e.U32 = i;
		e.U64 = i;
		e.F32 = float(i);
		e.F64 = double(i);
		e.D16 = Float16_t(i);
		e.D32 = Double32_t(i);
		e.Str = TString::Format("evt-%03d", i);

		e.P3  = {i-1, double(i), i-1};
		e.P3Ptr = new P3{i, double(i), i+1};
		e.ObjStr = TObjString(TString::Format("obj-%03d", i));
		e.ObjStrPtr = new TObjString(TString::Format("obj-ptr-%03d", i));

		for (int ii = 0; ii != ARRAYSZ; ii++) {
			e.ArrayI16[ii] = i;
			e.ArrayI32[ii] = i;
			e.ArrayI64[ii] = i;
			e.ArrayU16[ii] = i;
			e.ArrayU32[ii] = i;
			e.ArrayU64[ii] = i;
			e.ArrayF32[ii] = float(i);
			e.ArrayF64[ii] = double(i);
			e.ArrayP3s[ii] = {ii,double(i),ii+1};
			e.ArrayObjStr[ii] = TObjString(TString::Format("obj-%03d", ii));
		}

		e.N = int32_t(i) % 10;
		e.SliceI16 = (int16_t*)malloc(sizeof(int16_t)*e.N);
		e.SliceI32 = (int32_t*)malloc(sizeof(int32_t)*e.N);
		e.SliceI64 = (int64_t*)malloc(sizeof(int64_t)*e.N);
		e.SliceU16 = (uint16_t*)malloc(sizeof(uint16_t)*e.N);
		e.SliceU32 = (uint32_t*)malloc(sizeof(uint32_t)*e.N);
		e.SliceU64 = (uint64_t*)malloc(sizeof(uint64_t)*e.N);
		e.SliceF32 = (float*)malloc(sizeof(float)*e.N);
		e.SliceF64 = (double*)malloc(sizeof(double)*e.N);

		for (int ii = 0; ii != e.N; ii++) {
			e.SliceI16[ii] = i;
			e.SliceI32[ii] = i;
			e.SliceI64[ii] = i;
			e.SliceU16[ii] = i;
			e.SliceU32[ii] = i;
			e.SliceU64[ii] = i;
			e.SliceF32[ii] = float(i);
			e.SliceF64[ii] = double(i);
		}

		e.StdStr = std::string(TString::Format("std-%03d", i));
		e.StlVecI16.resize(e.N);
		e.StlVecI32.resize(e.N);
		e.StlVecI64.resize(e.N);
		e.StlVecU16.resize(e.N);
		e.StlVecU32.resize(e.N);
		e.StlVecU64.resize(e.N);
		e.StlVecF32.resize(e.N);
		e.StlVecF64.resize(e.N);
		e.StlVecStr.resize(e.N);
		for (int ii =0; ii != e.N; ii++) {
			e.StlVecI16[ii] = i;
			e.StlVecI32[ii] = i;
			e.StlVecI64[ii] = i;
			e.StlVecU16[ii] = i;
			e.StlVecU32[ii] = i;
			e.StlVecU64[ii] = i;
			e.StlVecF32[ii] = float(i);
			e.StlVecF64[ii] = double(i);
			e.StlVecStr[ii] = std::string(TString::Format("vec-%03d", i));
		}
		e.End = TString::Format("end-%03d", i);
	}

	f->WriteObjectAny(evt, "Event", "evt");

	f->Write();
	f->Close();
}
`
