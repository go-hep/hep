// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-diff compares the content of 2 ROOT files, including the content of
// their Trees (for all entries), if any.
//
// Example:
//
//	$> root-diff ./ref.root ./chk.root
//	$> root-diff -k=key1,tree,my-tree ./ref.root ./chk.root
//
//	$> root-diff -h
//	Usage: root-diff [options] a.root b.root
//
//	ex:
//	 $> root-diff ./testdata/small-flat-tree.root ./testdata/small-flat-tree.root
//
//	options:
//	  -k string
//	    	comma-separated list of keys to inspect and compare (default=all common keys)
package main // import "go-hep.org/x/hep/groot/cmd/root-diff"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rcmd"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
)

func main() {
	keysFlag := flag.String("k", "", "comma-separated list of keys to inspect and compare (default=all common keys)")

	log.SetPrefix("root-diff: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: root-diff [options] a.root b.root

ex:
 $> root-diff ./testdata/small-flat-tree.root ./testdata/small-flat-tree.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		log.Fatalf("need 2 input ROOT files to compare")
	}

	err := rootdiff(flag.Arg(0), flag.Arg(1), *keysFlag)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func rootdiff(ref, chk string, keysFlag string) error {
	fref, err := groot.Open(ref)
	if err != nil {
		return fmt.Errorf("could not open reference file: %w", err)
	}
	defer fref.Close()

	fchk, err := groot.Open(chk)
	if err != nil {
		return fmt.Errorf("could not open check file: %w", err)
	}
	defer fchk.Close()

	var keys []string
	if keysFlag != "" {
		keys = strings.Split(keysFlag, ",")
	}

	err = rcmd.Diff(nil, fchk, fref, keys)
	if err != nil {
		return fmt.Errorf("files differ: %w", err)
	}

	return nil
}
