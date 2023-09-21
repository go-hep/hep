// Copyright ©2024 The go-hep Authors. All rights reserved.
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
	root = flag.String("f", "test-tscatter.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentscatter", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
void gentscatter(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");

	const int n = 5;
	double xs[n] = {0, 1, 2, 3, 4};
	double ys[n] = {0, 2, 4, 6, 8};
	double cs[n] = {1, 3, 5, 7, 9};
	double ss[n] = {2, 4, 6, 8, 10};

	auto s = new TScatter(n, xs, ys, cs, ss);
	s->SetMarkerStyle(20);
	s->SetMarkerColor(kRed);
	s->SetTitle("Scatter plot;X;Y");
	s->SetName("scatter");

	s->Draw("A"); // generate underlying TH2F.
	auto h = s->GetHistogram();
	if (h == NULL) {
		exit(1);
	}

	f->WriteTObject(s);
	f->Write();
	f->Close();

	exit(0);
}
`
