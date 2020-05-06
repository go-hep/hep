// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-split splits an input file+tree into multiple output ROOT files,
// each containing N entries.
//
// Usage: root-split [options] file.root
//
// ex:
//  $> root-split -o out.root -n 10 ./testdata/chain.flat.1.root
//
// options:
//   -n int
//     	number of events to split into (default 100)
//   -o string
//     	path to output ROOT files (default "out.root")
//   -t string
//     	input tree name to split (default "tree")
//   -v	enable verbose mode
package main // import "go-hep.org/x/hep/groot/cmd/root-split"

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot/rcmd"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
)

func main() {
	log.SetPrefix("root-split: ")
	log.SetFlags(0)

	var (
		oname   = flag.String("o", "out.root", "path to output ROOT files")
		tname   = flag.String("t", "tree", "input tree name to split")
		verbose = flag.Bool("v", false, "enable verbose mode")
		nevts   = flag.Int64("n", 100, "number of events to split into")
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-split [options] file.root

ex:
 $> root-split -o out.root -n 10 ./testdata/chain.flat.1.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input file")
	}

	fname := flag.Arg(0)

	_, err := rcmd.Split(*oname, fname, *tname, *nevts, *verbose)
	if err != nil {
		log.Fatalf("could not split ROOT file: %+v", err)
	}
}
