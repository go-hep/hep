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
	root = flag.String("f", "test-tconfidence-level.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentconflvl", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <vector>
#include "TConfidenceLevel.h"

void gentconflvl(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");
	auto lvl = new TConfidenceLevel(3);

	auto xs = std::vector<Double_t>{1, 2, 3};

	lvl->SetTSD(3);
	lvl->SetTSB(xs.data());
	lvl->SetTSS(xs.data());
	lvl->SetLRS(xs.data());
	lvl->SetLRB(xs.data());
	lvl->SetBtot(3);
	lvl->SetStot(2);
	lvl->SetDtot(5);

	f->WriteTObject(lvl, "clvl");

	auto limit = new TLimit;
	f->WriteObjectAny(limit, "TLimit", "limit");

	auto dsrc = new TLimitDataSource;
	f->WriteTObject(dsrc, "dsrc");

	f->Write();
	f->Close();

	exit(0);
}
`
