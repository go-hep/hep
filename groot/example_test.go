// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package groot_test

import (
	"compress/flate"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleCreate_emptyFile() {
	const fname = "empty.root"
	defer os.Remove(fname)

	w, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	// empty file. close it.
	err = w.Close()
	if err != nil {
		log.Fatalf("could not close empty file: %v", err)
	}

	// read back.
	r, err := groot.Open(fname)
	if err != nil {
		log.Fatalf("could not open empty file: %v", err)
	}
	defer r.Close()

	fmt.Printf("file: %q\n", r.Name())

	// Output:
	// file: "empty.root"
}

func ExampleCreate() {
	const fname = "objstring.root"
	defer os.Remove(fname)

	w, err := groot.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var (
		k = "my-objstring"
		v = rbase.NewObjString("Hello World from Go-HEP!")
	)

	err = w.Put(k, v)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wkeys: %d\n", len(w.Keys()))

	err = w.Close()
	if err != nil {
		log.Fatalf("could not close file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer r.Close()

	fmt.Printf("rkeys: %d\n", len(r.Keys()))

	for _, k := range r.Keys() {
		fmt.Printf("key: name=%q, type=%q\n", k.Name(), k.ClassName())
	}

	obj, err := r.Get(k)
	if err != nil {
		log.Fatal(err)
	}
	rv := obj.(root.ObjString)
	fmt.Printf("objstring=%q\n", rv)

	// Output:
	// wkeys: 1
	// rkeys: 1
	// key: name="my-objstring", type="TObjString"
	// objstring="Hello World from Go-HEP!"
}

func ExampleCreate_withZlib() {
	const fname = "objstring-zlib.root"
	defer os.Remove(fname)

	w, err := groot.Create(fname, riofs.WithZlib(flate.BestCompression))
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var (
		k = "my-objstring"
		v = rbase.NewObjString("Hello World from Go-HEP!")
	)

	err = w.Put(k, v)
	if err != nil {
		log.Fatal(err)
	}

	err = w.Close()
	if err != nil {
		log.Fatalf("could not close writable file: %v", err)
	}

	r, err := groot.Open(fname)
	if err != nil {
		log.Fatalf("could not open file: %v", err)
	}
	defer r.Close()

	for _, k := range r.Keys() {
		fmt.Printf("key: name=%q, type=%q\n", k.Name(), k.ClassName())
	}

	obj, err := r.Get(k)
	if err != nil {
		log.Fatalf("could not get key %q: %v", k, err)
	}
	rv := obj.(root.ObjString)
	fmt.Printf("objstring=%q\n", rv)

	// Output:
	// key: name="my-objstring", type="TObjString"
	// objstring="Hello World from Go-HEP!"
}

func ExampleOpen() {
	f, err := groot.Open("testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, key := range f.Keys() {
		fmt.Printf("key:  %q cycle=%d title=%q\n", key.Name(), key.Cycle(), key.Title())
	}

	obj, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)
	fmt.Printf("tree: %q, entries=%d\n", tree.Name(), tree.Entries())

	// Output:
	// key:  "tree" cycle=1 title="fake data"
	// tree: "tree", entries=4
}

func ExampleOpen_overXRootD() {
	f, err := groot.Open("root://eospublic.cern.ch//eos/root-eos/cms_opendata_2012_nanoaod/Run2012B_DoubleMuParked.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, key := range f.Keys() {
		fmt.Printf("key:  %q cycle=%d title=%q\n", key.Name(), key.Cycle(), key.Title())
	}

	obj, err := f.Get("Events")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)
	fmt.Printf("tree: %q, entries=%d\n", tree.Name(), tree.Entries())

	// Output:
	// key:  "Events" cycle=1 title="Events"
	// tree: "Events", entries=29308627
}

func ExampleOpen_graph() {
	f, err := groot.Open("testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tg")
	if err != nil {
		log.Fatal(err)
	}

	g := obj.(rhist.Graph)
	fmt.Printf("name:  %q\n", g.Name())
	fmt.Printf("title: %q\n", g.Title())
	fmt.Printf("#pts:  %d\n", g.Len())
	for i := 0; i < g.Len(); i++ {
		x, y := g.XY(i)
		fmt.Printf("(x,y)[%d] = (%+e, %+e)\n", i, x, y)
	}

	// Output:
	// name:  "tg"
	// title: "graph without errors"
	// #pts:  4
	// (x,y)[0] = (+1.000000e+00, +2.000000e+00)
	// (x,y)[1] = (+2.000000e+00, +4.000000e+00)
	// (x,y)[2] = (+3.000000e+00, +6.000000e+00)
	// (x,y)[3] = (+4.000000e+00, +8.000000e+00)
}
