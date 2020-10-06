// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "test-ndim.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <string.h>
#include <stdio.h>

#define OFFSET 0
#define N0 2
#define N1 3
#define N2 4
#define N3 5

struct Event {
	bool     ArrayBs[2][3][4][5];
	int8_t   ArrayI8[2][3][4][5];
	int16_t  ArrayI16[2][3][4][5];
	int32_t  ArrayI32[2][3][4][5];
	int64_t  ArrayI64[2][3][4][5];
	uint8_t  ArrayU8[2][3][4][5];
	uint16_t ArrayU16[2][3][4][5];
	uint32_t ArrayU32[2][3][4][5];
	uint64_t ArrayU64[2][3][4][5];
	float    ArrayF32[2][3][4][5];
	double   ArrayF64[2][3][4][5];

	Float16_t    ArrayD16[2][3][4][5];
	Double32_t   ArrayD32[2][3][4][5];

};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 2;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("ArrBs",  e.ArrayBs,  "ArrBs[2][3][4][5]/O");
	t->Branch("ArrI8",  e.ArrayI8,  "ArrI8[2][3][4][5]/B");
	t->Branch("ArrI16", e.ArrayI16, "ArrI16[2][3][4][5]/S");
	t->Branch("ArrI32", e.ArrayI32, "ArrI32[2][3][4][5]/I");
	t->Branch("ArrI64", e.ArrayI64, "ArrI64[2][3][4][5]/L");
	t->Branch("ArrU8",  e.ArrayU8,  "ArrU8[2][3][4][5]/b");
	t->Branch("ArrU16", e.ArrayU16, "ArrU16[2][3][4][5]/s");
	t->Branch("ArrU32", e.ArrayU32, "ArrU32[2][3][4][5]/i");
	t->Branch("ArrU64", e.ArrayU64, "ArrU64[2][3][4][5]/l");
	t->Branch("ArrF32", e.ArrayF32, "ArrF32[2][3][4][5]/F");
	t->Branch("ArrF64", e.ArrayF64, "ArrF64[2][3][4][5]/D");

	t->Branch("ArrD16", e.ArrayD16, "ArrD16[2][3][4][5]/f[0,0,16]");
	t->Branch("ArrD32", e.ArrayD32, "ArrD32[2][3][4][5]/d[0,0,32]");

	for (int j = 0; j != evtmax; j++) {
		int i = j + OFFSET;
		for (int i0 = 0; i0 != N0; i0++) {
			for (int i1 = 0; i1 != N1; i1++) {
				for (int i2 = 0; i2 != N2; i2++) {
					for (int i3 = 0; i3 != N3; i3++) {
						e.ArrayBs[i0][i1][i2][i3]  = i3%2 == 0;
						e.ArrayI8[i0][i1][i2][i3]  = -i;
						e.ArrayI16[i0][i1][i2][i3] = -i;
						e.ArrayI32[i0][i1][i2][i3] = -i;
						e.ArrayI64[i0][i1][i2][i3] = -i;
						e.ArrayU8[i0][i1][i2][i3]  = i;
						e.ArrayU16[i0][i1][i2][i3] = i;
						e.ArrayU32[i0][i1][i2][i3] = i;
						e.ArrayU64[i0][i1][i2][i3] = i;
						e.ArrayF32[i0][i1][i2][i3] = float(i);
						e.ArrayF64[i0][i1][i2][i3] = double(i);
						e.ArrayD16[i0][i1][i2][i3] = float(i);
						e.ArrayD32[i0][i1][i2][i3] = double(i);
						i++;
					}
				}
			}
		}
		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
