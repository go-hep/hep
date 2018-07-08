// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// lcio-ls displays the content of a LCIO file.
//
// The default behaviour is to only display RunHeaders and EventHeaders.
// Events' contents can be printed out with the --print-event flag.
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
		fname      = ""
		nevts      = flag.Int64("n", -1, "number of events to inspect")
		printEvent = flag.Bool("print-event", false, "enable event(s) printout")
	)

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}
	inspect(os.Stdout, fname, *nevts, *printEvent)
}

func inspect(w io.Writer, fname string, nevts int64, printEvent bool) {
	log.Printf("inspecting file [%s]...", fname)
	r, err := lcio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	rhdr := r.RunHeader()
	ehdr := r.EventHeader()

	evts := 0
	for ievt := int64(0); r.Next() && (nevts < 0 || ievt < nevts); ievt++ {
		if hdr := r.RunHeader(); !reflect.DeepEqual(hdr, rhdr) {
			fmt.Fprintf(w, "%v\n", &hdr)
			rhdr = hdr
		}
		if hdr := r.EventHeader(); !reflect.DeepEqual(hdr, ehdr) {
			if !printEvent {
				fmt.Fprintf(w, "%v\n", &hdr)
			}
			ehdr = hdr
		}
		if printEvent {
			evt := r.Event()
			fmt.Fprintf(w, "%v\n", &evt)
		}
		evts++
	}
	err = r.Err()
	if err == io.EOF && evts == 0 {
		if hdr := r.RunHeader(); !reflect.DeepEqual(hdr, rhdr) {
			fmt.Fprintf(w, "%v\n", &hdr)
		}
		if hdr := r.EventHeader(); !reflect.DeepEqual(hdr, ehdr) {
			fmt.Fprintf(w, "%v\n", &hdr)
		}
	}
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
}
