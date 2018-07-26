// +build ignore

// Copyright 2018 The go-hep Authors. All rights reserved.
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
	root = flag.String("f", "dirs.root", "output ROOT file")
)

func main() {
	const fname = "gendirs.C"

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
