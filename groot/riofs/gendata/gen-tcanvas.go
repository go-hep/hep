// Copyright ©2023 The go-hep Authors. All rights reserved.
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
	root = flag.String("f", "test-tcanvas.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentcanvas", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
void gentcanvas(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");
	auto c = new TCanvas("c1", "c1-title", 300, 400);

	c->AddExec("ex1", ".ls");
	c->AddExec("ex2", ".ls");

	const Int_t np = 5;
	Double_t x[np]       = {0, 1, 2, 3, 4};
	Double_t y[np]       = {0, 2, 4, 1, 3};

	auto gr = new TGraph(np, x, y);
	gr->Draw();
	gr->Fit("pol1");

	c->SetFixedAspectRatio();

	f->WriteTObject(c);
	f->Write();
	f->Close();

	exit(0);
}
`
