// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-ls lists the content of a ROOT file.
//
// Usage: root-ls [options] file1.root [file2.root [...]]
//
// ex:
//
//  $> root-ls -t ./testdata/graphs.root ./testdata/small-flat-tree.root
//  === [./testdata/graphs.root] ===
//  version: 60806
//  TGraph            tg      graph without errors         (cycle=1)
//  TGraphErrors      tge     graph with errors            (cycle=1)
//  TGraphAsymmErrors tgae    graph with asymmetric errors (cycle=1)
//
//  === [./testdata/small-flat-tree.root] ===
//  version: 60804
//  TTree          tree                 my tree title (entries=100)
//    Int32        "Int32/I"            TBranch
//    Int64        "Int64/L"            TBranch
//    UInt32       "UInt32/i"           TBranch
//    UInt64       "UInt64/l"           TBranch
//    Float32      "Float32/F"          TBranch
//    Float64      "Float64/D"          TBranch
//    ArrayInt32   "ArrayInt32[10]/I"   TBranch
//    ArrayInt64   "ArrayInt64[10]/L"   TBranch
//    ArrayUInt32  "ArrayInt32[10]/i"   TBranch
//    ArrayUInt64  "ArrayInt64[10]/l"   TBranch
//    ArrayFloat32 "ArrayFloat32[10]/F" TBranch
//    ArrayFloat64 "ArrayFloat64[10]/D" TBranch
//    N            "N/I"                TBranch
//    SliceInt32   "SliceInt32[N]/I"    TBranch
//    SliceInt64   "SliceInt64[N]/L"    TBranch
//    SliceUInt32  "SliceInt32[N]/i"    TBranch
//    SliceUInt64  "SliceInt64[N]/l"    TBranch
//    SliceFloat32 "SliceFloat32[N]/F"  TBranch
//    SliceFloat64 "SliceFloat64[N]/D"  TBranch
//
package main // import "go-hep.org/x/hep/groot/cmd/root-ls"

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"

	"go-hep.org/x/hep/groot/rcmd"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	_ "go-hep.org/x/hep/groot/ztypes"
)

var (
	fset = flag.NewFlagSet("ls", flag.ContinueOnError)

	siFlag   = fset.Bool("sinfos", false, "print StreamerInfos")
	treeFlag = fset.Bool("t", false, "print Tree(s) (recursively)")
	cpuFlag  = fset.String("cpu-profile", "", "path to CPU profile output file")

	usage = `Usage: root-ls [options] file1.root [file2.root [...]]

ex:
 $> root-ls ./testdata/graphs.root
 $> root-ls -t -sinfos ./testdata/graphs.root

options:
`
)

func main() {
	os.Exit(run(os.Stdout, os.Stderr, os.Args[1:]))
}

func run(stdout, stderr io.Writer, args []string) int {
	fset.Usage = func() {
		fmt.Fprint(stderr, usage)
		fset.PrintDefaults()
	}

	err := fset.Parse(args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		log.Printf("could not parse args %q: %+v", args, err)
		return 1
	}

	if *cpuFlag != "" {
		f, err := os.Create(*cpuFlag)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("could not start CPU profiling: %+v", err)
		}
		defer pprof.StopCPUProfile()
	}

	if fset.NArg() <= 0 {
		fmt.Fprintf(stderr, "error: you need to give a ROOT file\n\n")
		fset.Usage()
		return 1
	}

	out := bufio.NewWriter(stdout)
	defer out.Flush()

	opts := []rcmd.ListOption{
		rcmd.ListStreamers(*siFlag),
		rcmd.ListTrees(*treeFlag),
	}

	for ii, fname := range fset.Args() {
		if ii > 0 {
			fmt.Fprintf(out, "\n")
		}
		err := rcmd.List(out, fname, opts...)
		if err != nil {
			out.Flush()
			log.Printf("%+v", err)
			return 1
		}
	}

	return 0
}
