// Copyright ©2017 The go-hep Authors. All rights reserved.
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
//  https://numpy.org/neps/nep-0001-npy-format.html
//
// Example:
//
//  $> root2npy -f $GOPATH/src/go-hep.org/x/hep/groot/testdata/simple.root -t tree -o output.npz
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
//  $> root-ls -t $GOPATH/src/go-hep.org/x/hep/groot/testdata/simple.root
//  === [$GOPATH/src/go-hep.org/x/hep/groot/testdata/simple.root] ===
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
//  $> root-ls -t $GOPATH/src/go-hep.org/x/hep/groot/testdata/small-flat-tree.root
//  === [$GOPATH/src/go-hep.org/x/hep/groot/testdata/small-flat-tree.root] ===
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
//  $> root2npy $GOPATH/src/go-hep.org/x/hep/groot/testdata/small-flat-tree.root
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
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/sbinet/npyio"

	"go-hep.org/x/hep/groot"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
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

	err := process(*oname, *fname, *tname)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func process(oname, fname, tname string) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ROOT file: %w", err)
	}
	defer f.Close()

	obj, err := f.Get(tname)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	tree, ok := obj.(rtree.Tree)
	if !ok {
		return fmt.Errorf("object %q in file %q is not a rtree.Tree", tname, fname)
	}

	var (
		nt    = ntuple{n: tree.Entries()}
		rvars = rtree.NewReadVars(tree)
	)
	log.Printf("scanning leaves...")
	for _, rvar := range rvars {
		rv := reflect.ValueOf(rvar.Value).Elem()
		switch rv.Kind() {
		case reflect.Struct, reflect.Slice:
			log.Printf(">>> %q %T not supported", rvar.Name, rv.Interface())
			continue
		}
		nt.add(rvar)
	}
	log.Printf("scanning leaves... [done]")

	r, err := rtree.NewReader(tree, nt.args)
	if err != nil {
		return fmt.Errorf("could not create ROOT reader: %w", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		nt.fill()
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not read ROOT data: %w", err)
	}

	out, err := os.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create NumPy file: %w", err)
	}
	defer out.Close()

	npz := zip.NewWriter(out)
	defer npz.Close()

	work := make([]byte, 1*1024*1024)
	for _, col := range nt.cols {
		buf := new(bytes.Buffer)
		err = npyio.Write(buf, col.slice.Interface())
		if err != nil {
			return fmt.Errorf("could not write %q: %w", col.name, err)
		}

		wz, err := npz.Create(col.name)
		if err != nil {
			return fmt.Errorf("could not create column %q: %w", col.name, err)
		}

		_, err = io.CopyBuffer(wz, buf, work)
		if err != nil {
			return fmt.Errorf("could not save column %q: %w", col.name, err)
		}
	}

	err = npz.Flush()
	if err != nil {
		return fmt.Errorf("could not flush NumPy zip-file: %w", err)
	}

	err = npz.Close()
	if err != nil {
		return fmt.Errorf("could not close NumPy zip-file: %w", err)
	}

	err = out.Close()
	if err != nil {
		return fmt.Errorf("could not close NumPy file: %w", err)
	}

	return nil
}

type ntuple struct {
	n    int64
	cols []column
	args []rtree.ReadVar
}

func (nt *ntuple) add(rvar rtree.ReadVar) {
	nt.cols = append(nt.cols, newColumn(rvar, nt.n))
	nt.args = append(nt.args, rvar)
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
	etype reflect.Type
	shape []int
	data  reflect.Value
	slice reflect.Value
}

func newColumn(rvar rtree.ReadVar, n int64) column {
	etype := reflect.TypeOf(rvar.Value).Elem()
	shape := []int{int(n)}
	rtype := reflect.SliceOf(etype)
	return column{
		name:  rvar.Name,
		i:     0,
		etype: etype,
		shape: shape,
		data:  reflect.ValueOf(rvar.Value).Elem(),
		slice: reflect.MakeSlice(rtype, int(n), int(n)),
	}
}

func (col *column) fill() {
	col.slice.Index(int(col.i)).Set(col.data)
	col.i++
}
