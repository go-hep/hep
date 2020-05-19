// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-dump dumps the content of a ROOT file, including the content of
// the Trees (for all entries), if any.
//
// Example:
//
//  $> root-dump ./testdata/small-flat-tree.root
//  >>> file[./testdata/small-flat-tree.root]
//  key[000]: tree;1 "my tree title" (TTree)
//  [000][Int32]: 0
//  [000][Int64]: 0
//  [000][UInt32]: 0
//  [000][UInt64]: 0
//  [000][Float32]: 0
//  [000][Float64]: 0
//  [000][Str]: evt-000
//  [000][ArrayInt32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayInt64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayFloat32]: [0 0 0 0 0 0 0 0 0 0]
//  [000][ArrayFloat64]: [0 0 0 0 0 0 0 0 0 0]
//  [000][N]: 0
//  [000][SliceInt32]: []
//  [000][SliceInt64]: []
//  [...]
//
//  $> root-dump -h
//  Usage: root-dump [options] f0.root [f1.root [...]]
//
//  ex:
//   $> root-dump ./testdata/small-flat-tree.root
//   $> root-dump -deep=0 ./testdata/small-flat-tree.root
//
//  options:
//    -deep
//      	enable deep dumping of values (including Trees' entries) (default true)
//    -name string
//      	regex of object names to dump
//
package main // import "go-hep.org/x/hep/groot/cmd/root-dump"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"go-hep.org/x/hep/groot/rcmd"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
)

var (
	deepFlag = flag.Bool("deep", true, "enable deep dumping of values (including Trees' entries)")
	nameFlag = flag.String("name", "", "regex of object names to dump")
)

func main() {
	log.SetPrefix("root-dump: ")
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: root-dump [options] f0.root [f1.root [...]]

ex:
 $> root-dump ./testdata/small-flat-tree.root
 $> root-dump -deep=0 ./testdata/small-flat-tree.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *nameFlag != "" {
		reName = regexp.MustCompile(*nameFlag)
	}

	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatalf("need at least one input ROOT file")
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for _, fname := range flag.Args() {
		err := dump(out, fname, *deepFlag)
		if err != nil {
			out.Flush()
			log.Fatalf("error dumping file %q: %+v", fname, err)
		}
	}
}

func dump(w io.Writer, fname string, deep bool) error {
	fmt.Fprintf(w, ">>> file[%s]\n", fname)
	return rcmd.Dump(w, fname, deep, match)
}

var reName *regexp.Regexp

func match(name string) bool {
	if reName == nil {
		return true
	}
	return reName.MatchString(name)
}
