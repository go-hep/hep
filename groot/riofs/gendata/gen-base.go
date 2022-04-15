// Copyright ©2022 The go-hep Authors. All rights reserved.
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
	root  = flag.String("f", "test-base.root", "output ROOT file")
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
#include <stdint.h>

class Base {
public:
	int32_t I32;
};

class D1 : public Base {
public:
	int32_t D32;
};

class D2 : public Base {
public:
	int32_t I32;
};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 2;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	D1 d1;
	D2 d2;

	t->Branch("d1", &d1, bufsize, splitlvl);
	t->Branch("d2", &d2, bufsize, splitlvl);

	for (int i = 0; i != evtmax; i++) {
		d1.I32 = i+1;
		d1.D32 = i+2;
		((Base*)&d2)->I32 = i+3;
		d2.I32 = i+4;

		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
