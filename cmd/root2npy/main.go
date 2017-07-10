// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2npy converts the content of a ROOT TTree to a NumPy data file.
//
//  Usage of root2npy:
//   -f string
//     	path to input ROOT file name
//   -o string
//     	path to output npz file name (default "output.npz")
//   -t string
//     	name of the tree to convert (default "tree")
//
package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/sbinet/npyio"

	"go-hep.org/x/hep/rootio"
)

func main() {
	log.SetPrefix("root2npy: ")
	log.SetFlags(0)

	fname := flag.String("f", "", "path to input ROOT file name")
	oname := flag.String("o", "output.npz", "path to output npz file name")
	tname := flag.String("t", "tree", "name of the tree to convert")

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing input ROOT filename argument")
	}

	f, err := rootio.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, ok := f.Get(*tname)
	if !ok {
		log.Fatalf("no object named %q in file %q", *tname, *fname)
	}

	tree, ok := obj.(rootio.Tree)
	if !ok {
		log.Fatalf("object %q in file %q is not a rootio.Tree", *tname, *fname)
	}

	var nt = ntuple{n: tree.Entries()}
	log.Printf("scanning leaves...")
	for _, leaf := range tree.Leaves() {
		if leaf.Kind() == reflect.String {
			nt.add(leaf.Name(), leaf)
			continue
		}
		if leaf.Class() == "TLeafElement" { // FIXME(sbinet): find a better, type-safe way
			log.Printf(">>> %q %v not supported", leaf.Name(), leaf.Class())
			continue
		}
		if leaf.LeafCount() != nil {
			log.Printf(">>> %q []%v not supported", leaf.Name(), leaf.TypeName())
			continue
		}
		nt.add(leaf.Name(), leaf)
	}
	log.Printf("scanning leaves... [done]")

	sc, err := rootio.NewScannerVars(tree, nt.args...)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		err = sc.Scan(nt.vars...)
		if err != nil {
			log.Fatal(err)
		}
		nt.fill()
	}

	out, err := os.Create(*oname)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	npz := zip.NewWriter(out)
	defer npz.Close()

	for _, col := range nt.cols {
		buf := new(bytes.Buffer)
		err = npyio.Write(buf, col.slice.Interface())
		if err != nil {
			log.Fatalf("error writing %q: %v\n", col.name, err)
		}

		wz, err := npz.Create(col.name)
		if err != nil {
			log.Fatalf("error creating %q: %v\n", col.name, err)
		}

		_, err = io.Copy(wz, buf)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = npz.Flush()
	if err != nil {
		log.Fatal(err)
	}

	err = npz.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

type ntuple struct {
	n    int64
	cols []column
	args []rootio.ScanVar
	vars []interface{}
}

func (nt *ntuple) add(name string, leaf rootio.Leaf) {
	n := len(nt.cols)
	nt.cols = append(nt.cols, newColumn(name, leaf, nt.n))
	col := &nt.cols[n]
	nt.args = append(nt.args, rootio.ScanVar{Name: name, Type: col.etype})
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
	leaf  rootio.Leaf
	etype reflect.Type
	shape []int
	data  reflect.Value
	slice reflect.Value
}

func newColumn(name string, leaf rootio.Leaf, n int64) column {
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
