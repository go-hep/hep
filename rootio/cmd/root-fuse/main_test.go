// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFUSESimple(t *testing.T) {
	tmp, err := ioutil.TempDir("", "root-fuse-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remote %q: %v", tmp, err)
		}
	}()

	ready := make(chan struct{})
	quit := make(chan os.Signal)
	defer func() {
		select {
		case quit <- os.Interrupt:
		default:
		}
	}()
	const verbose = false
	go run(tmp, "../../testdata/simple.root", verbose, ready, quit)

	<-ready
	f, err := os.Open(tmp)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dirs, err := f.Readdir(-1)
	if err != nil {
		t.Fatal(err)
	}

	if len(dirs) != 1 {
		t.Fatalf("got %d dir-entries. want=1", len(dirs))
	}
	if got, want := dirs[0].Name(), "tree"; got != want {
		t.Fatalf("invalid entry name. got %q. want %q", got, want)
	}

	got, err := ioutil.ReadFile(filepath.Join(tmp, "tree"))
	if err != nil {
		t.Fatalf("could not read /tmp/tree: %v", err)
	}
	want := []byte(`name:  tree
title: fake data
type:  TTree
`)

	if !bytes.Equal(got, want) {
		t.Fatalf("/tmp/tree contents differ.\ngot = %q\nwant= %q\n", got, want)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("could not close root-dir: %v", err)
	}

	quit <- os.Interrupt
	<-quit
}

func TestFUSEDirs(t *testing.T) {
	tmp, err := ioutil.TempDir("", "root-fuse-")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := os.RemoveAll(tmp)
		if err != nil {
			t.Logf("could not remote %q: %v", tmp, err)
		}
	}()

	ready := make(chan struct{})
	quit := make(chan os.Signal)
	defer func() {
		select {
		case quit <- os.Interrupt:
		default:
		}
	}()
	const verbose = false
	go run(tmp, "../../testdata/dirs-6.14.00.root", verbose, ready, quit)

	<-ready
	f, err := os.Open(tmp)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dirs, err := f.Readdir(-1)
	if err != nil {
		t.Fatal(err)
	}

	if len(dirs) != 3 {
		t.Fatalf("got %d dir-entries. want=3", len(dirs))
	}
	for i, dir := range []string{"dir1", "dir2", "dir3"} {
		if got, want := dirs[i].Name(), dir; got != want {
			t.Fatalf("invalid entry name. got %q. want %q", got, want)
		}
	}

	got, err := ioutil.ReadFile(filepath.Join(tmp, "dir1", "dir11", "h1"))
	if err != nil {
		t.Fatalf("could not read /tmp/dir1/dir11/h1: %v", err)
	}
	want := []byte(`name:  h1
title: h1
type:  TH1F
`)

	if !bytes.Equal(got, want) {
		t.Fatalf("/tmp/dir1/dir11/h1 contents differ.\ngot = %q\nwant= %q\n", got, want)
	}

	err = f.Close()
	if err != nil {
		t.Fatalf("could not close root-dir: %v", err)
	}

	quit <- os.Interrupt
	<-quit
}
