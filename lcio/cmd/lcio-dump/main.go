// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"os"

	"github.com/kr/pretty"
	"go-hep.org/x/hep/lcio"
)

func main() {
	log.SetPrefix("lcio-dump: ")
	log.SetFlags(0)

	var (
		fname = ""
		nevts = flag.Int64("n", -1, "number of events to dump")
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

	for ievt := int64(0); r.Next() && (*nevts < 0 || ievt < *nevts); ievt++ {
		evt := r.Event()
		log.Printf("%s", pretty.Sprintf("ievt[%d]: % #v\n", ievt, evt))
	}
	err = r.Err()
	if err != nil {
		log.Fatal(err)
	}
}
