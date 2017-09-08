// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-diff compares the content of 2 ROOT files, including the content of
// their Trees (for all entries), if any.
package main // import "go-hep.org/x/hep/rootio/cmd/root-diff"
import (
	"flag"
	"log"

	"go-hep.org/x/hep/rootio"
)

func main() {
	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		log.Fatalf("need 2 input ROOT files to compare")
	}

	fref, err := rootio.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer fref.Close()

	fchk, err := rootio.Open(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer fchk.Close()
}
