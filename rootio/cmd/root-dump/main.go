// Copyright 2017 The go-hep Authors.  All rights reserved.
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
//
package main // import "go-hep.org/x/hep/rootio/cmd/root-dump"

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"go-hep.org/x/hep/rootio"
)

var (
	deepFlag = flag.Bool("deep", true, "enable deep dumping of values (including Trees' entries)")
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

	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatalf("need at least one input ROOT file")
	}

	for _, fname := range flag.Args() {
		err := dump(os.Stdout, fname, *deepFlag)
		if err != nil {
			log.Fatalf("error dumping file %q: %v", fname, err)
		}
	}
}

func dump(w io.Writer, fname string, deep bool) error {
	fmt.Fprintf(w, ">>> file[%s]\n", fname)
	f, err := rootio.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	return dumpDir(w, f, deep)
}

func dumpDir(w io.Writer, dir rootio.Directory, deep bool) error {
	for i, key := range dir.Keys() {
		obj, err := key.Object()
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "key[%03d]: %s;%d %q (%s)", i, key.Name(), key.Cycle(), key.Title(), obj.Class())
		if deep {
			switch obj := obj.(type) {
			case rootio.Tree:
				fmt.Fprintf(w, "\n")
				err = dumpTree(w, obj)
			case rootio.Directory:
				fmt.Fprintf(w, "\n")
				err = dumpDir(w, obj, deep)
			default:
				fmt.Fprintf(w, " => ignoring key of type %T\n", obj)
				continue
			}
		} else {
			fmt.Fprintf(w, "\n")
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func dumpTree(w io.Writer, t rootio.Tree) error {

	var vars []rootio.ScanVar
	for _, b := range t.Branches() {
		for _, leaf := range b.Leaves() {
			ptr := newValue(leaf)
			vars = append(vars, rootio.ScanVar{Name: b.Name(), Leaf: leaf.Name(), Value: ptr})
		}
	}

	sc, err := rootio.NewScannerVars(t, vars...)
	if err != nil {
		return err
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan()
		if err != nil {
			return fmt.Errorf("error scanning entry %d: %v", sc.Entry(), err)
		}
		for _, v := range vars {
			rv := reflect.Indirect(reflect.ValueOf(v.Value))
			fmt.Fprintf(w, "[%03d][%s]: %v\n", sc.Entry(), v.Name, rv.Interface())
		}
	}
	return nil
}

func newValue(leaf rootio.Leaf) interface{} {
	etype := leaf.Type()
	switch {
	case leaf.LeafCount() != nil:
		etype = reflect.SliceOf(etype)
	case leaf.Len() > 1 && leaf.Kind() != reflect.String:
		etype = reflect.ArrayOf(leaf.Len(), etype)
	}
	return reflect.New(etype).Interface()
}
