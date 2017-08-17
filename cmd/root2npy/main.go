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
// The NumPy data file format is described here:
//
//  https://docs.scipy.org/doc/numpy/neps/npy-format.html
//
// Example:
//
//  $> root2npy -f $GOPATH/src/go-hep.org/x/hep/rootio/testdata/simple.root -t tree -o output.npz
//  $> python2 -c 'import sys, numpy as np; print(dict(np.load(sys.argv[1])))' ./output.npz
//  {'one':   array([1, 2, 3, 4], dtype=int32),
//   'two':   array([ 1.10000002,  2.20000005,  3.29999995,  4.4000001 ], dtype=float32),
//   'three': array([u'uno', u'dos', u'tres', u'quatro'], dtype='<U6')}
//
//  $> python3 -c 'import sys, numpy as np; print(dict(np.load(sys.argv[1])))' ./output.npz
//  {'one':   array([1, 2, 3, 4], dtype=int32),
//   'two':   array([ 1.10000002,  2.20000005,  3.29999995,  4.4000001 ], dtype=float32),
//   'three': array(['uno', 'dos', 'tres', 'quatro'], dtype='<U6')}
//
//  $> go get github.com/sbinet/npyio/cmd/npyio-ls
//  $> npyio-ls ./output.npz
//  ================================================================================
//  file: ./output.npz
//  entry: one
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<i4, Fortran:false, Shape:[4]}}
//  data = [1 2 3 4]
//
//  entry: two
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<f4, Fortran:false, Shape:[4]}}
//  data = [1.1 2.2 3.3 4.4]
//
//  entry: three
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<U6, Fortran:false, Shape:[4]}}
//  data = [uno dos tres quatro]
//
//  $> root-ls -t $GOPATH/src/go-hep.org/x/hep/rootio/testdata/simple.root
//  === [$GOPATH/src/go-hep.org/x/hep/rootio/testdata/simple.root] ===
//  version: 60600
//  TTree   tree      fake data (entries=4)
//    one   "one/I"   TBranch
//    two   "two/F"   TBranch
//    three "three/C" TBranch
//
// If you have a 10-events tree with a branch "doubles" containing an array of 3 float64,
// root2npy will convert it to a NumPy data file containing a NumPy array with a shape (10,3).
//
// Example:
//
//  $> root-ls -t $GOPATH/src/go-hep.org/x/hep/rootio/testdata/small-flat-tree.root
//  === [$GOPATH/src/go-hep.org/x/hep/rootio/testdata/small-flat-tree.root] ===
//  version: 60806
//  TTree          tree                 my tree title (entries=100)
//    Int32        "Int32/I"            TBranch
//    Int64        "Int64/L"            TBranch
//    UInt32       "UInt32/i"           TBranch
//    UInt64       "UInt64/l"           TBranch
//    Float32      "Float32/F"          TBranch
//    Float64      "Float64/D"          TBranch
//    Str          "Str/C"              TBranch
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
//  $> root2npy $GOPATH/src/go-hep.org/x/hep/rootio/testdata/small-flat-tree.root
//  root2npy: scanning leaves...
//  root2npy: >>> "SliceInt32" []int32 not supported
//  root2npy: >>> "SliceInt64" []int64 not supported
//  root2npy: >>> "SliceInt32" []int32 not supported
//  root2npy: >>> "SliceInt64" []int64 not supported
//  root2npy: >>> "SliceFloat32" []float32 not supported
//  root2npy: >>> "SliceFloat64" []float64 not supported
//  root2npy: scanning leaves... [done]
//
//  $> npyio-ls ./output.npz
//  ================================================================================
//  file: ./output.npz
//  entry: Int32
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<i4, Fortran:false, Shape:[100]}}
//  data = [0 1 2 3 4 5 6 7 8 9 10 11 ... 90 91 92 93 94 95 96 97 98 99]
//
//  entry: Int64
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<i8, Fortran:false, Shape:[100]}}
//  data = [0 1 2 3 4 5 6 7 8 9 10 11 ... 90 91 92 93 94 95 96 97 98 99]
//
//  [...]
//
//  entry: Float64
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<f8, Fortran:false, Shape:[100]}}
//  data = [0 1 2 3 4 5 6 7 8 9 10 11 ... 90 91 92 93 94 95 96 97 98 99]
//
//  entry: Str
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<U7, Fortran:false, Shape:[100]}}
//  data = [evt-000 evt-001 evt-002 evt-003 evt-004 evt-005 evt-006 evt-007 ...
//  evt-092 evt-093 evt-094 evt-095 evt-096 evt-097 evt-098 evt-099]
//
//  entry: ArrayInt32
//  npy-header: Header{Major:2, Minor:0, Descr:{Type:<i4, Fortran:false, Shape:[100 10]}}
//  data = [0 0 0 0 0 0 0 0 0 0 1 1 1 1 1 1 1 1 1 1 2 2 2 2 2 2 2 2 2 2 ...
//  ... 97 98 98 98 98 98 98 98 98 98 98 99 99 99 99 99 99 99 99 99 99]
//
//  [...]
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

	sc, err := rootio.NewTreeScannerVars(tree, nt.args...)
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
	nt.args = append(nt.args, rootio.ScanVar{Name: name})
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
