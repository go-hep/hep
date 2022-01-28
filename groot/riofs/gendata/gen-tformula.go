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
	root = flag.String("f", "test-tformula.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentformula", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include "TFormula.h"
#include "TF1.h"
#include "TF1Convolution.h"

double func2(double *x, double *par) {
	return par[0] + par[1]*x[0];
};

void gentformula(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");

	auto f1 = new TF1("func1", "[0] + [1]*x", 0, 10);
	f1->SetParNames("p0", "p1");
	f1->SetParameters(1, 2);
	f1->SetParameter(0, 10);
	f1->SetParameter(1, 20);
	f1->SetChisquare(0.2);
	f1->SetNDF(2);
	f1->SetNumberFitPoints(101);
	f1->SetNormalized(true);

	auto f2 = new TF1("func2", func2, 0, 10, 2);
	f2->SetParNames("p0", "p1");
	f2->SetParameters(1, 2);
	f2->SetParameter(0, 10);
	f2->SetParameter(1, 20);
	f2->SetChisquare(0.2);
	f2->SetNDF(2);
	f2->SetNumberFitPoints(101);

	auto conv = new TF1Convolution("expo", "gaus", -1, 6, true);
	conv->SetRange(-1., 6.);
	conv->SetNofPointsFFT(1000);
	auto f3 = new TF1("func3", *conv, 0, 5, conv->GetNpar());
	f3->SetParNames("p0", "p1", "p2", "p3");
	f3->SetParameters(1., -0.3, 0., 1.);
	f3->SetChisquare(0.2);

	auto norm = new TF1NormSum(f1, f2, 10, 20); 
	auto f4 = new TF1("func4", *norm, 0, 5, norm->GetNpar());
	f4->SetChisquare(0.2);

	f->WriteTObject(f1);
	f->WriteTObject(f2);
	f->WriteTObject(f3);
	f->WriteTObject(f4);
	f->WriteTObject(conv, "fconv");
	f->WriteTObject(norm, "fnorm");

	f->Write();
	f->Close();

	exit(0);
}
`
