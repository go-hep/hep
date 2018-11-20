// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestXrdCp(t *testing.T) {
	dir, err := ioutil.TempDir("", "xrootd-xrdcp-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dst := filepath.Join(dir, "chain.1.root")
	src := "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio/testdata/chain.1.root"

	const (
		recursive = false
		verbose   = true
	)

	err = xrdcopy(dst, src, recursive, verbose)
	if err != nil {
		t.Fatalf("could not copy remote file: %v", err)
	}
}

func BenchmarkXrdCp_Small(b *testing.B) {
	benchmarkXrdCp(b, "root://ccxrootdgotest.in2p3.fr:9001/tmp/rootio/testdata/chain.1.root")
}

func BenchmarkXrdCp_Medium(b *testing.B) {
	benchmarkXrdCp(b, "root://eospublic.cern.ch//eos/root-eos/cms_opendata_2012_nanoaod/SMHiggsToZZTo4L.root")
}

func BenchmarkXrdCp_Large(b *testing.B) {
	benchmarkXrdCp(b, "root://eospublic.cern.ch//eos/root-eos/cms_opendata_2012_nanoaod/Run2012B_DoubleElectron.root")
}

func benchmarkXrdCp(b *testing.B, src string) {
	dir, err := ioutil.TempDir("", "xrootd-xrdcp-")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dst := filepath.Join(dir, filepath.Base(src))

	const (
		recursive = false
		verbose   = false
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		os.RemoveAll(dst)
		err = xrdcopy(dst, src, recursive, verbose)
		if err != nil {
			b.Fatalf("could not copy remote file: %v", err)
		}
	}
}
