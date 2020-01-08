// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-cp selects and copies keys from a ROOT file to another ROOT file.
//
// Usage: root-cp [options] file1.root[:REGEXP] [file2.root[:REGEXP] [...]] out.root
//
// ex:
//
//  $> root-cp f.root out.root
//  $> root-cp f1.root f2.root f3.root out.root
//  $> root-cp f1.root:hist.* f2.root:h2 out.root
//
package main // import "go-hep.org/x/hep/groot/cmd/root-cp"

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
	log.SetPrefix("root-cp: ")
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-cp [options] file1.root[:REGEXP] [file2.root[:REGEXP] [...]] out.root

ex:
 $> root-cp f.root out.root
 $> root-cp f1.root f2.root f3.root out.root
 $> root-cp f1.root:hist.* f2.root:h2 out.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() < 2 {
		log.Printf("error: you need to give input and output ROOT files\n\n")
		flag.Usage()
		os.Exit(1)
	}

	dst := flag.Arg(flag.NArg() - 1)
	srcs := flag.Args()[:flag.NArg()-1]

	err := rcmd.Copy(dst, srcs)
	if err != nil {
		log.Fatal(err)
	}
}
