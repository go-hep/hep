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
	root = flag.String("f", "test-tgme.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentgme", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include "TGraphMultiErrors.h"

void gentgme(const char* fname) {
	const Int_t np = 5;
	Double_t x[np]       = {0, 1, 2, 3, 4};
	Double_t y[np]       = {0, 2, 4, 1, 3};
	Double_t exl[np]     = {0.3, 0.3, 0.3, 0.3, 0.3};
	Double_t exh[np]     = {0.3, 0.3, 0.3, 0.3, 0.3};
	Double_t eylstat[np] = {1, 0.5, 1, 0.5, 1};
	Double_t eyhstat[np] = {0.5, 1, 0.5, 1, 2};
	Double_t eylsys[np]  = {0.5, 0.4, 0.8, 0.3, 1.2};
	Double_t eyhsys[np]  = {0.6, 0.7, 0.6, 0.4, 0.8};

	auto gme = new TGraphMultiErrors(
		"gme", "TGraphMultiErrors Example",
		np, x, y, exl, exh, eylstat, eyhstat
	);
	gme->AddYError(np, eylsys, eyhsys);
	gme->SetMarkerStyle(20);
	gme->SetLineColor(kRed);
	gme->GetAttLine(0)->SetLineColor(kRed);
	gme->GetAttLine(1)->SetLineColor(kBlue);
	gme->GetAttFill(1)->SetFillStyle(0);

	auto f = TFile::Open(fname, "RECREATE");
	f->WriteTObject(gme, "gme");

	f->Write();
	f->Close();

	exit(0);
}
`
