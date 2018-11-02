// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs_test

import (
	"compress/flate"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

func ExampleCreate_empty() {
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
