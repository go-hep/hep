// Copyright 2017 The go-hep Authors.  All rights reserved.
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
package main // import "go-hep.org/x/hep/rootio/cmd/root-ls"

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"text/tabwriter"

	"go-hep.org/x/hep/rootio"
)

var (
	doProf = flag.String("profile", "", "filename of cpuprofile")
	doSI   = flag.Bool("sinfos", false, "print StreamerInfos")
	doTree = flag.Bool("t", false, "print Tree(s) (recursively)")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-ls [options] file1.root [file2.root [...]]

ex:
 $> root-ls ./testdata/graphs.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *doProf != "" {
		f, err := os.Create(*doProf)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if flag.NArg() <= 0 {
		fmt.Fprintf(os.Stderr, "error: you need to give a ROOT file\n\n")
		flag.Usage()
		os.Exit(1)
	}

	for ii, fname := range flag.Args() {

		if ii > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("=== [%s] ===\n", fname)
		f, err := rootio.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "rootio: failed to open [%s]: %v\n", fname, err)
			os.Exit(1)
		}
		defer f.Close()
		fmt.Printf("version: %v\n", f.Version())
		if *doSI {
			fmt.Printf("streamer-infos:\n")
			sinfos := f.StreamerInfos()
			for _, v := range sinfos {
				name := v.Name()
				fmt.Printf(" StreamerInfo for %q version=%d title=%q\n", name, v.ClassVersion(), v.Title())
				w := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)
				for _, elm := range v.Elements() {
					fmt.Fprintf(w, "  %s\t%s\toffset=%3d\ttype=%3d\tsize=%3d\t %s\n", elm.TypeName(), elm.Name(), elm.Offset(), elm.Type(), elm.Size(), elm.Title())
				}
				w.Flush()
			}
		}

		w := tabwriter.NewWriter(os.Stdout, 8, 4, 1, ' ', 0)
		for _, k := range f.Keys() {
			walk(w, k)
		}
		w.Flush()
	}
}

func walk(w io.Writer, k rootio.Key) {
	obj := k.Value()
	if *doTree {
		tree, ok := obj.(rootio.Tree)
		if ok {
			w := tabwriter.NewWriter(w, 8, 4, 1, ' ', 0)
			fmt.Fprintf(w, "%s\t%s\t%s\t(entries=%d)\n", k.Class(), k.Name(), k.Title(), tree.Entries())
			displayBranches(w, tree, 2)
			w.Flush()
			return
		}
	}
	fmt.Fprintf(w, "%s\t%s\t%s\t(cycle=%d)\n", k.Class(), k.Name(), k.Title(), k.Cycle())
	if dir, ok := obj.(rootio.Directory); ok {
		w := newWindent(2, w)
		for _, k := range dir.Keys() {
			walk(w, k)
		}
		w.Flush()
	}
}

type windent struct {
	hdr []byte
	w   io.Writer
}

func newWindent(n int, w io.Writer) *windent {
	return &windent{
		hdr: bytes.Repeat([]byte(" "), n),
		w:   w,
	}
}

func (w *windent) Write(data []byte) (int, error) {
	return w.w.Write(append(w.hdr, data...))
}

func (w *windent) Flush() error {
	ww, ok := w.w.(flusher)
	if !ok {
		return nil
	}
	return ww.Flush()
}

type flusher interface {
	Flush() error
}

type brancher interface {
	Branches() []rootio.Branch
}

func displayBranches(w io.Writer, bres brancher, indent int) {
	branches := bres.Branches()
	if len(branches) <= 0 {
		return
	}
	ww := newWindent(indent, w)
	for _, b := range branches {
		fmt.Fprintf(ww, "%s\t%q\t%v\n", b.Name(), b.Title(), b.Class())
		displayBranches(ww, b, 2)
	}
	ww.Flush()
}
