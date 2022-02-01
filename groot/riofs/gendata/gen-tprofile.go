// Copyright ©2022 The go-hep Authors. All rights reserved.
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
	root = flag.String("f", "test-tprofile.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentprofile", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include "TProfile.h"

void gentprofile(const char* fname) {
	auto p1d = new TProfile("p1d","Profile of pz versus px",100,-4,4,0,20);
	auto p2d = new TProfile2D("p2d","Profile of pz versus px and py",40,-4,4,40,-4,4,0,20);

	Float_t px, py, pz;
	for (Int_t i=0; i<25000; i++) {
		gRandom->Rannor(px,py);
		pz = px*px + py*py;
		p1d->Fill(px,pz,1);
		p2d->Fill(px,py,pz,1);
	}

	auto f = TFile::Open(fname, "RECREATE");
	f->WriteTObject(p1d, "p1d");
	f->WriteTObject(p2d, "p2d");

	f->Write();
	f->Close();

	exit(0);
}
`
