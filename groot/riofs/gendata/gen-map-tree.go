// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"

	"go-hep.org/x/hep/groot/internal/rtests"
)

var (
	root  = flag.String("f", "stdmap.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const script = `
#include <map>
#include <vector>
#include <string>

const int ARRAYSZ = 10;

struct P3 {
	int32_t Px;
	double  Py;
	int32_t Pz;
};

struct Event {
	std::map<int32_t, int32_t> mi32;

	std::map<std::string, int32_t>     msi32;
	std::map<std::string, std::string> mss;
//	std::map<std::string, P3>          msp3;

//	std::map<std::string, std::vector<std::string> > msvs;
//	std::map<std::string, std::vector<int32_t> >     msvi32;
//	std::map<std::string, std::vector<P3> >          msvp3;

	void clear() {
		this->mi32.clear();

		this->msi32.clear();
		this->mss.clear();
//		this->msp3.clear();

//		this->msvs.clear();
//		this->msvi32.clear();
//		this->msvp3.clear();
	}
};

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 10;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	Event e;

	t->Branch("evt", &e, bufsize, splitlvl);

	for (int i = 0; i != evtmax; i++) {
		e.clear();
		for (int ii = 0; ii < i; ii++) {
			e.mi32[int32_t(ii)] = int32_t(ii);

			std::string key = std::string(TString::Format("key-%03d", ii).Data());

			e.msi32[key] = int32_t(ii);
			e.mss[key] = std::string(TString::Format("val-%03d", ii).Data());
//			e.msp3[key] = P3{ii, double(ii+1), ii+2};

//			e.msvs[key] = std::vector<std::string>({
//					{TString::Format("val-%03d", ii).Data()}
//					,{TString::Format("val-%03d", ii+1).Data()}
//					,{TString::Format("val-%03d", ii+2).Data()}
//			});
//			e.msvi32[key] = std::vector<int32_t>({1, ii, 3, ii});
//			e.msvp3[key] = std::vector<P3>({{ii, double(ii+1), ii+2}, {ii+1, double(ii+2), ii+3}});
		}

		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
