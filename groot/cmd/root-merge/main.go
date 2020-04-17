// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-merge merges ROOT files' content into a merged ROOT file.
//
// Usage: root-merge [options] file1.root [file2.root [file3.root [...]]]
//
// ex:
//  $> root-merge -o out.root ./testdata/chain.flat.1.root ./testdata/chain.flat.2.root
//
// options:
//   -o string
//     	path to merged output ROOT file (default "out.root")
//   -v	enable verbose mode
//
package main // import "go-hep.org/x/hep/groot/cmd/root-merge"

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
	log.SetPrefix("root-merge: ")
	log.SetFlags(0)

	var (
		oname   = flag.String("o", "out.root", "path to merged output ROOT file")
		verbose = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-merge [options] file1.root [file2.root [file3.root [...]]]

ex:
 $> root-merge -o out.root ./testdata/chain.flat.1.root ./testdata/chain.flat.2.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		log.Fatalf("missing input files")
	}

	fnames := flag.Args()

	err := rcmd.Merge(*oname, fnames, *verbose)
	if err != nil {
		log.Fatalf("could not merge ROOT files: %+v", err)
	}
}
