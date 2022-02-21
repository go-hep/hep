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
	root  = flag.String("f", "tdatime.root", "output ROOT file")
	split = flag.Int("split", 0, "default split-level for TTree")
)

func main() {
	flag.Parse()

	out, err := rtests.RunCxxROOT("gentree", []byte(class+script), *root, *split)
	if err != nil {
		log.Fatalf("could not run ROOT macro:\noutput:\n%v\nerror: %+v", string(out), err)
	}
}

const class = `
#ifndef EVENT_H
#define EVENT_H 1

#include "TObject.h"
#include "TDatime.h"

class TFoo : public TObject {
public:
    TDatime d;
ClassDef(TFoo, 1)
};

class TBar : public TObject {
public:
    TDatime d;

	// pad object to increase its size to fit version header that
	// TDatime doesn't stream.
	char pad[6] = "12345";
ClassDef(TBar, 1);
};

class Date {
public:
    TDatime d;

	// pad object to increase its size to fit version header that
	// TDatime doesn't stream.
	char pad[6] = "12345";
ClassDef(Date, 1);
};

#endif // EVENT_H
`

const script = `
#include "TFile.h"
#include "TTree.h"

void gentree(const char* fname, int splitlvl = 99) {
	int bufsize = 32000;
	int evtmax = 2;

	auto f = TFile::Open(fname, "RECREATE");
	auto t = new TTree("tree", "my tree title");

	auto tda = TDatime(2006, 1, 2, 15, 4, 5);
	TFoo foo; foo.d = tda;
	TBar bar; bar.d = tda;
	Date dat; dat.d = tda;

	f->WriteObjectAny(&tda, "TDatime", "tda");
	f->WriteTObject(&foo,              "foo");
	f->WriteTObject(&bar,              "bar");
	f->WriteObjectAny(&dat, "Date",    "dat");

	t->Branch("b0", &tda, bufsize, splitlvl);
	t->Branch("b1", &foo, bufsize, splitlvl);
	t->Branch("b2", &bar, bufsize, splitlvl);
	t->Branch("b3", &dat, bufsize, splitlvl);

	for (int i = 0; i != evtmax; i++) {
		  tda.Set(2006, 1, 2+i, 15, 4, 5);
		foo.d.Set(2006, 1, 2+i, 15, 4, 5);
		bar.d.Set(2006, 1, 2+i, 15, 4, 5);
		dat.d.Set(2006, 1, 2+i, 15, 4, 5);
		t->Fill();
	}

	f->Write();
	f->Close();

	exit(0);
}
`
