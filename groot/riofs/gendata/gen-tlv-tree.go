// Copyright ©2021 The go-hep Authors. All rights reserved.
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
	root  = flag.String("f", "tlv.root", "output ROOT file")
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
#include "TLorentzVector.h"

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 10;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	{
		TLorentzVector *p4 = new TLorentzVector;
		p4->SetPxPyPzE(10, 20, 30, 40);
		f->WriteTObject(p4, "tlv");
	}

	TLorentzVector *tlv = new TLorentzVector;

	t->Branch("p4", &tlv, bufsize, splitlvl);

	for (int i = 0; i != evtmax; i++) {
		tlv->SetPxPyPzE(0+i, 1+i, 2+i, 3+i);
		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
