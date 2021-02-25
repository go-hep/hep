// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root = flag.String("f", "padding-struct.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <string.h>
#include <stdio.h>

struct Pad {
	int8_t  x1;
	int64_t x2;
	int8_t  x3;
};

struct Nop {
	int64_t x1;
	int8_t  x2;
	int8_t  x3;
};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 5;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "tree w/ & w/o padding");

	Pad pad;
	Nop nop;

	t->Branch("pad", &pad, "x1/B:x2/L:x3/B");
	t->Branch("nop", &nop, "x1/L:x2/B:x3/B");

	for (int j = 0; j != evtmax; j++) {
		pad.x1 = j;
		pad.x2 = j;
		pad.x3 = j;

		nop.x1 = j;
		nop.x2 = j;
		nop.x3 = j;

		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
