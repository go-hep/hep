// +build ignore

// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var (
	root = flag.String("f", "tlv.root", "output ROOT file")
)

func main() {
	const fname = "gentlv.C"

	flag.Parse()
	err := ioutil.WriteFile(fname, []byte(script), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(fname)

	cmd := exec.Command("root.exe", "-b", fmt.Sprintf("./%s(%q)", fname, *root))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

const script = `
#include <vector>
#include "TLorentzVector.h"
#include "TTree.h"

const int ARRAYSZ = 10;

struct Event {
	TLorentzVector P4;
	std::vector<TLorentzVector> P4s;
};

void gentlv(const char* fname) {
	int bufsize = 32000;
	int evtmax = 10;
	int split = 99;
	int nosplit = 0;

	auto f = TFile::Open(fname, "RECREATE");
	auto t1 = new TTree("t1", "my split tree");
	auto t2 = new TTree("t2", "my no split tree");
	auto t3 = new TTree("t3", "my flat tree");
	auto t4 = new TTree("t4", "my split flat tree");

	Event evt;

	t1->Branch("evt", &evt, bufsize, split);
	t2->Branch("evt", &evt, bufsize, nosplit);
	t3->Branch("tlv", &evt.P4, bufsize, split);
	t3->Branch("p4s", &evt.P4s, bufsize, nosplit);
	t4->Branch("tlv", &evt.P4, bufsize, split);
	t4->Branch("p4s", &evt.P4s, bufsize, split);

	for (int i = 0; i != evtmax; i++) {
		evt.P4.SetPxPyPzE(1, 2, 3, 4);
		evt.P4s.resize(2);
		evt.P4s[0].SetPxPyPzE(i, i+1, i+2, i+3);
		evt.P4s[1].SetPxPyPzE(i*2, i*2 + 1, i*2 + 2, i*2 + 3);

		t1->Fill();
		t2->Fill();
		t3->Fill();
		t4->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
