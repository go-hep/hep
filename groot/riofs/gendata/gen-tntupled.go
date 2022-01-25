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
	root = flag.String("f", "test-tntupled.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentntuple", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
void gentntuple(const char* fname) {
	int bufsize = 32000;
	int evtmax = 10;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TNtupleD("ntup", "my ntuple title", "x:y");

	for (int i = 0; i != evtmax; i++) {
		t->Fill(i, i+0.5);
	}
	f->Write();
	f->Close();

	exit(0);
}
`
