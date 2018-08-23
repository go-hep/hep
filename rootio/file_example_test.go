// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio_test

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/rootio"
)

func ExampleCreate_empty() {
	const fname = "empty.root"
	defer os.Remove(fname)

	w, err := rootio.Create(fname)
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
	r, err := rootio.Open(fname)
	if err != nil {
		log.Fatalf("could not open empty file: %v", err)
	}
	defer r.Close()

	fmt.Printf("file: %q\n", r.Name())

	sinfos := r.StreamerInfos()
	fmt.Printf("streamer infos: %d\n", len(sinfos))

	// Output:
	// file: "empty.root"
	// streamer infos: 0
}

func ExampleCreate() {
	const fname = "objstring.root"
	defer os.Remove(fname)

	w, err := rootio.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer w.Close()

	var (
		k = "my-objstring"
		v = rootio.NewObjString("Hello World from Go-HEP!")
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

	r, err := rootio.Open(fname)
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
	rv := obj.(rootio.ObjString)
	fmt.Printf("objstring=%q\n", rv)

	// Output:
	// wkeys: 1
	// rkeys: 1
	// key: name="my-objstring", type="TObjString"
	// objstring="Hello World from Go-HEP!"
}
