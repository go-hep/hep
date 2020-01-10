// +build ignore

// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "test-small.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	const fname = "gentree.C"

	flag.Parse()

	err := ioutil.WriteFile(fname, []byte(script), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fname)

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <string.h>
#include <stdio.h>

const int ARRAYSZ  = 10;
const int MAXSLICE = 20;
const int MAXSTR   = 32;

#define OFFSET 0

struct Event {
	bool     Bool;
	char     Str[MAXSTR];
	int8_t   I8;
	int16_t  I16;
	int32_t  I32;
	int64_t  I64;
	uint8_t  U8;
	uint16_t U16;
	uint32_t U32;
	uint64_t U64;
	float    F32;
	double   F64;

	Double32_t D32;

	bool     ArrayBs[ARRAYSZ];
//	char     ArrayStr[ARRAYSZ][MAXSTR];
	int8_t   ArrayI8[ARRAYSZ];
	int16_t  ArrayI16[ARRAYSZ];
	int32_t  ArrayI32[ARRAYSZ];
	int64_t  ArrayI64[ARRAYSZ];
	uint8_t  ArrayU8[ARRAYSZ];
	uint16_t ArrayU16[ARRAYSZ];
	uint32_t ArrayU32[ARRAYSZ];
	uint64_t ArrayU64[ARRAYSZ];
	float    ArrayF32[ARRAYSZ];
	double   ArrayF64[ARRAYSZ];

	int32_t  N;
	bool     SliceBs[MAXSLICE];   //[N]
//	char     SliceStr[MAXSLICE][MAXSTR]; //[N]
	int8_t   SliceI8[MAXSLICE];   //[N]
	int16_t  SliceI16[MAXSLICE];  //[N]
	int32_t  SliceI32[MAXSLICE];  //[N]
	int64_t  SliceI64[MAXSLICE];  //[N]
	uint8_t  SliceU8[MAXSLICE];   //[N]
	uint16_t SliceU16[MAXSLICE];  //[N]
	uint32_t SliceU32[MAXSLICE];  //[N]
	uint64_t SliceU64[MAXSLICE];  //[N]
	float    SliceF32[MAXSLICE];  //[N]
	double   SliceF64[MAXSLICE];  //[N]
};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 10;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("B",   &e.Bool, "B/O");
	t->Branch("Str",  e.Str,  "Str/C");
	t->Branch("I8",  &e.I8,   "I8/B");
	t->Branch("I16", &e.I16,  "I16/S");
	t->Branch("I32", &e.I32,  "I32/I");
	t->Branch("I64", &e.I64,  "I64/L");
	t->Branch("U8",  &e.U8,   "U8/b");
	t->Branch("U16", &e.U16,  "U16/s");
	t->Branch("U32", &e.U32,  "U32/i");
	t->Branch("U64", &e.U64,  "U64/l");
	t->Branch("F32", &e.F32,  "F32/F");
	t->Branch("F64", &e.F64,  "F64/D");
	t->Branch("D32", &e.D32,  "D32/d[0,0,32]");

	// static arrays

	t->Branch("ArrBs",  e.ArrayBs,  "ArrBs[10]/O");
//	t->Branch("ArrStr", e.ArrayStr, "ArrStr[10][32]/C");
	t->Branch("ArrI8",  e.ArrayI8,  "ArrI8[10]/B");
	t->Branch("ArrI16", e.ArrayI16, "ArrI16[10]/S");
	t->Branch("ArrI32", e.ArrayI32, "ArrI32[10]/I");
	t->Branch("ArrI64", e.ArrayI64, "ArrI64[10]/L");
	t->Branch("ArrU8",  e.ArrayU8,  "ArrU8[10]/b");
	t->Branch("ArrU16", e.ArrayU16, "ArrU16[10]/s");
	t->Branch("ArrU32", e.ArrayU32, "ArrU32[10]/i");
	t->Branch("ArrU64", e.ArrayU64, "ArrU64[10]/l");
	t->Branch("ArrF32", e.ArrayF32, "ArrF32[10]/F");
	t->Branch("ArrF64", e.ArrayF64, "ArrF64[10]/D");

	// dynamic arrays
	
	t->Branch("N", &e.N, "N/I");
	t->Branch("SliBs",  e.SliceBs,  "SliBs[N]/O");
//	t->Branch("SliStr", e.SliceStr, "SliStr[N][32]/C");
	t->Branch("SliI8",  e.SliceI8,  "SliI8[N]/B");
	t->Branch("SliI16", e.SliceI16, "SliI16[N]/S");
	t->Branch("SliI32", e.SliceI32, "SliI32[N]/I");
	t->Branch("SliI64", e.SliceI64, "SliI64[N]/L");
	t->Branch("SliU8",  e.SliceU8,  "SliU8[N]/b");
	t->Branch("SliU16", e.SliceU16, "SliU16[N]/s");
	t->Branch("SliU32", e.SliceU32, "SliU32[N]/i");
	t->Branch("SliU64", e.SliceU64, "SliU64[N]/l");
	t->Branch("SliF32", e.SliceF32, "SliF32[N]/F");
	t->Branch("SliF64", e.SliceF64, "SliF64[N]/D");


	for (int j = 0; j != evtmax; j++) {
		int i = j + OFFSET;
		e.Bool = (i % 2) == 0;
		strncpy(e.Str, TString::Format("str-%d", i).Data(), 32);
		e.I8  = -i;
		e.I16 = -i;
		e.I32 = -i;
		e.I64 = -i;
		e.U8  = i;
		e.U16 = i;
		e.U32 = i;
		e.U64 = i;
		e.F32 = float(i);
		e.F64 = double(i);
		e.D32 = double(i);

		for (int ii = 0; ii != ARRAYSZ; ii++) {
			e.ArrayBs[ii]  = ii == i;
//			sprintf(e.ArrayStr[ii], "ars-%d-%d", i, ii);
			e.ArrayI8[ii]  = -i;
			e.ArrayI16[ii] = -i;
			e.ArrayI32[ii] = -i;
			e.ArrayI64[ii] = -i;
			e.ArrayU8[ii]  = i;
			e.ArrayU16[ii] = i;
			e.ArrayU32[ii] = i;
			e.ArrayU64[ii] = i;
			e.ArrayF32[ii] = float(i);
			e.ArrayF64[ii] = double(i);
		}

		e.N = int32_t(i) % 10;
		for (int ii = 0; ii != e.N; ii++) {
			e.SliceBs[ii]  = (ii+1) == i;
//			strncpy(e.SliceStr[ii], TString::Format("sli-%d-%d", i, ii).Data(), 32);
			e.SliceI8[ii]  = -i;
			e.SliceI16[ii] = -i;
			e.SliceI32[ii] = -i;
			e.SliceI64[ii] = -i;
			e.SliceU8[ii]  = i;
			e.SliceU16[ii] = i;
			e.SliceU32[ii] = i;
			e.SliceU64[ii] = i;
			e.SliceF32[ii] = float(i);
			e.SliceF64[ii] = double(i);
		}

		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
