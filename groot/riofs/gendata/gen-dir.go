// Copyright ©2018 The go-hep Authors. All rights reserved.
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
	root = flag.String("f", "dirs.root", "output ROOT file")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gendirs", []byte(script), *root)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
void gendirs(const char* fname) {
	auto f = TFile::Open(fname, "RECREATE");

	auto dir1 = f->mkdir("dir1");
	f->mkdir("dir2");
	f->mkdir("dir3");

	dir1->cd();
	auto dir11 = dir1->mkdir("dir11");
	dir11->cd();

	auto h1 = new TH1F("h1", "h1", 100, 0, 100);
	h1->FillRandom("gaus", 5);

	f->Write();
	f->Close();

	exit(0);
}
`
