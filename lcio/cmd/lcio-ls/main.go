// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"go-hep.org/x/hep/lcio"
)

func main() {
	log.SetPrefix("lcio-ls: ")
	log.SetFlags(0)

	var (
		fname = ""
		nevts = flag.Int64("n", -1, "number of events to inspect")
	)

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	log.Printf("inspecting file [%s]...", fname)
	r, err := lcio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	rhdr := r.RunHeader()
	ehdr := r.EventHeader()

	evts := 0
	for ievt := int64(0); r.Next() && (*nevts < 0 || ievt < *nevts); ievt++ {
		if hdr := r.RunHeader(); !reflect.DeepEqual(hdr, rhdr) {
			fmt.Printf("%v\n", &hdr)
			rhdr = hdr
		}
		if hdr := r.EventHeader(); !reflect.DeepEqual(hdr, ehdr) {
			fmt.Printf("%v\n", &hdr)
			ehdr = hdr
		}
		evts++
	}
	err = r.Err()
	if err == io.EOF && evts == 0 {
		if hdr := r.RunHeader(); !reflect.DeepEqual(hdr, rhdr) {
			fmt.Printf("%v\n", &hdr)
		}
		if hdr := r.EventHeader(); !reflect.DeepEqual(hdr, ehdr) {
			fmt.Printf("%v\n", &hdr)
		}
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
