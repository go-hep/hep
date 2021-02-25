// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "std-bitset.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	tmp, err := ioutil.TempDir("", "groot-")
	if err != nil {
		log.Fatalf("could not created tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	dict, err := rtests.GenROOTDictCode(event, link)
	if err != nil {
		log.Fatalf("could not run ROOT dict: %+v", err)
	}

	out, err := rtests.RunCxxROOT("gentree", []byte(event+string(dict)+script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const event = `
#ifndef EVENT_H
#define EVENT_H 1

#include <bitset>
#include <vector>

struct Event {
	std::bitset<8> Bs8;
	std::vector<std::bitset<8> > VecBs8;

	void clear() {
		this->VecBs8.clear();
	}
};

#endif // EVENT_H
`

const link = `
#ifdef __CINT__

#pragma link off all globals;
#pragma link off all classes;
#pragma link off all functions;

#pragma link C++ class Event+;

#endif
`

const script = `
void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 2;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("evt", &e, bufsize, splitlvl);

	// 0
	e.clear();
	e.Bs8 = std::bitset<8>("00010001");
	e.VecBs8.push_back(std::bitset<8>("11101110"));
	t->Fill();

	// 1
	e.clear();
	e.Bs8 = std::bitset<8>("10011001");
	e.VecBs8.push_back(std::bitset<8>("00010001"));
	e.VecBs8.push_back(std::bitset<8>("11101110"));
	t->Fill();

	// 2
	e.clear();
	e.Bs8 = std::bitset<8>("01100110");
	e.VecBs8.push_back(std::bitset<8>("10011001"));
	e.VecBs8.push_back(std::bitset<8>("01100110"));
	e.VecBs8.push_back(std::bitset<8>("11001100"));
	t->Fill();

	f->Write();
	f->Close();

	exit(0);
}
`
