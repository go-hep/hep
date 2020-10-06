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
	root  = flag.String("f", "test-ndim-slice.root", "output ROOT file")
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

#define MAXSLICE 5

#define OFFSET 0
#define N1 2
#define N2 3
#define N3 4

struct Event {
	int32_t  N;

	bool     SliceBs[MAXSLICE][2][3][4];  //[N]
	int8_t   SliceI8[MAXSLICE][2][3][4];  //[N]
	int16_t  SliceI16[MAXSLICE][2][3][4]; //[N]
	int32_t  SliceI32[MAXSLICE][2][3][4]; //[N]
	int64_t  SliceI64[MAXSLICE][2][3][4]; //[N]
	uint8_t  SliceU8[MAXSLICE][2][3][4];  //[N]
	uint16_t SliceU16[MAXSLICE][2][3][4]; //[N]
	uint32_t SliceU32[MAXSLICE][2][3][4]; //[N]
	uint64_t SliceU64[MAXSLICE][2][3][4]; //[N]
	float    SliceF32[MAXSLICE][2][3][4]; //[N]
	double   SliceF64[MAXSLICE][2][3][4]; //[N]

	Float16_t    SliceD16[MAXSLICE][2][3][4]; //[N]
	Double32_t   SliceD32[MAXSLICE][2][3][4]; //[N]

};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 2;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("N", &e.N, "N/I");

	t->Branch("SliBs",  e.SliceBs,  "SliBs[N][2][3][4]/O");
	t->Branch("SliI8",  e.SliceI8,  "SliI8[N][2][3][4]/B");
	t->Branch("SliI16", e.SliceI16, "SliI16[N][2][3][4]/S");
	t->Branch("SliI32", e.SliceI32, "SliI32[N][2][3][4]/I");
	t->Branch("SliI64", e.SliceI64, "SliI64[N][2][3][4]/L");
	t->Branch("SliU8",  e.SliceU8,  "SliU8[N][2][3][4]/b");
	t->Branch("SliU16", e.SliceU16, "SliU16[N][2][3][4]/s");
	t->Branch("SliU32", e.SliceU32, "SliU32[N][2][3][4]/i");
	t->Branch("SliU64", e.SliceU64, "SliU64[N][2][3][4]/l");
	t->Branch("SliF32", e.SliceF32, "SliF32[N][2][3][4]/F");
	t->Branch("SliF64", e.SliceF64, "SliF64[N][2][3][4]/D");

	t->Branch("SliD16", e.SliceD16, "SliD16[N][2][3][4]/f[0,0,16]");
	t->Branch("SliD32", e.SliceD32, "SliD32[N][2][3][4]/d[0,0,32]");

	for (int j = 0; j != evtmax; j++) {
		int i = j + OFFSET;
		e.N = i+1;
		for (int i0 = 0; i0 != e.N; i0++) {
			for (int i1 = 0; i1 != N1; i1++) {
				for (int i2 = 0; i2 != N2; i2++) {
					for (int i3 = 0; i3 != N3; i3++) {
						e.SliceBs[i0][i1][i2][i3]  = i3%2 == 0;
						e.SliceI8[i0][i1][i2][i3]  = -i;
						e.SliceI16[i0][i1][i2][i3] = -i;
						e.SliceI32[i0][i1][i2][i3] = -i;
						e.SliceI64[i0][i1][i2][i3] = -i;
						e.SliceU8[i0][i1][i2][i3]  = i;
						e.SliceU16[i0][i1][i2][i3] = i;
						e.SliceU32[i0][i1][i2][i3] = i;
						e.SliceU64[i0][i1][i2][i3] = i;
						e.SliceF32[i0][i1][i2][i3] = float(i);
						e.SliceF64[i0][i1][i2][i3] = double(i);
						e.SliceD16[i0][i1][i2][i3] = float(i);
						e.SliceD32[i0][i1][i2][i3] = double(i);
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
