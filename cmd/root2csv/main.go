// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2csv converts the content of a ROOT TTree to a CSV file.
//
//  Usage of root2csv:
//    -f string
//      	path to input ROOT file name
//    -o string
//      	path to output CSV file name (default "output.csv")
//    -t string
//      	name of the tree to convert (default "tree")
//
// By default, root2csv will write out a CSV file with ';' as a column delimiter.
// root2csv ignores the branches of the TTree that are not supported by CSV:
//  - slices/arrays
//  - C++ objects
//
// Example:
//  $> root2csv -o out.csv -t tree -f testdata/small-flat-tree.root
//  $> head out.csv
//  ## Automatically generated from "testdata/small-flat-tree.root"
//  Int32;Int64;UInt32;UInt64;Float32;Float64;Str;N
//  0;0;0;0;0;0;evt-000;0
//  1;1;1;1;1;1;evt-001;1
//  2;2;2;2;2;2;evt-002;2
//  3;3;3;3;3;3;evt-003;3
//  4;4;4;4;4;4;evt-004;4
//  5;5;5;5;5;5;evt-005;5
//  6;6;6;6;6;6;evt-006;6
//  7;7;7;7;7;7;evt-007;7
package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"

	"go-hep.org/x/hep/csvutil"
	"go-hep.org/x/hep/groot"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("root2csv: ")
	log.SetFlags(0)

	fname := flag.String("f", "", "path to input ROOT file name")
	oname := flag.String("o", "output.csv", "path to output CSV file name")
	tname := flag.String("t", "tree", "name of the tree to convert")

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing input ROOT filename argument")
	}

	f, err := groot.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get(*tname)
	if err != nil {
		log.Fatal(err)
	}

	tree, ok := obj.(rtree.Tree)
	if !ok {
		log.Fatalf("object %q in file %q is not a rtree.Tree", *tname, *fname)
	}

	var nt = ntuple{n: tree.Entries()}
	log.Printf("scanning leaves...")
	for _, leaf := range tree.Leaves() {
		kind := leaf.Type().Kind()
		switch kind {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			log.Printf(">>> %q %v not supported (%v)", leaf.Name(), leaf.Class(), kind)
			continue
		case reflect.Bool,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String:
		default:
			log.Printf(">>> %q %v not supported (%v) (unknown!)", leaf.Name(), leaf.Class(), kind)
			continue
		}

		nt.add(leaf.Name(), leaf)
	}
	log.Printf("scanning leaves... [done]")

	sc, err := rtree.NewTreeScannerVars(tree, nt.args...)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	nrows := 0
	for sc.Next() {
		err = sc.Scan(nt.vars...)
		if err != nil {
			log.Fatal(err)
		}
		nt.fill()
		nrows++
	}

	tbl, err := csvutil.Create(*oname)
	if err != nil {
		log.Fatal(err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	names := make([]string, len(nt.cols))
	for i, col := range nt.cols {
		names[i] = col.name
	}
	err = tbl.WriteHeader(fmt.Sprintf(
		"## Automatically generated from %q\n%s\n",
		*fname,
		strings.Join(names, string(tbl.Writer.Comma)),
	))
	if err != nil {
		log.Fatalf("could not write header: %v", err)
	}

	row := make([]interface{}, len(nt.cols))
	for irow := 0; irow < nrows; irow++ {
		for i, col := range nt.cols {
			row[i] = col.slice.Index(irow).Interface()
		}
		err = tbl.WriteRow(row...)
		if err != nil {
			log.Fatalf("error writing row %d: %v", irow, err)
		}
	}

	err = tbl.Close()
	if err != nil {
		log.Fatalf("could not close CSV file: %v", err)
	}
}

type ntuple struct {
	n    int64
	cols []column
	args []rtree.ScanVar
	vars []interface{}
}

func (nt *ntuple) add(name string, leaf rtree.Leaf) {
	n := len(nt.cols)
	nt.cols = append(nt.cols, newColumn(name, leaf, nt.n))
	col := &nt.cols[n]
	nt.args = append(nt.args, rtree.ScanVar{Name: name, Leaf: leaf.Name()})
	nt.vars = append(nt.vars, col.data.Addr().Interface())
}

func (nt *ntuple) fill() {
	for i := range nt.cols {
		col := &nt.cols[i]
		col.fill()
	}
}

type column struct {
	name  string
	i     int64
	leaf  rtree.Leaf
	etype reflect.Type
	shape []int
	data  reflect.Value
	slice reflect.Value
}

func newColumn(name string, leaf rtree.Leaf, n int64) column {
	etype := leaf.Type()
	shape := []int{int(n)}
	if leaf.Len() > 1 && leaf.Kind() != reflect.String {
		etype = reflect.ArrayOf(leaf.Len(), etype)
		shape = append(shape, leaf.Len())
	}
	rtype := reflect.SliceOf(etype)
	return column{
		name:  name,
		i:     0,
		leaf:  leaf,
		etype: etype,
		shape: shape,
		data:  reflect.New(etype).Elem(),
		slice: reflect.MakeSlice(rtype, int(n), int(n)),
	}
}

func (col *column) fill() {
	col.slice.Index(int(col.i)).Set(col.data)
	col.i++
}
