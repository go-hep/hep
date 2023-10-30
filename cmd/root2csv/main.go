// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2csv converts the content of a ROOT TTree to a CSV file.
//
//	Usage of root2csv:
//	  -f string
//	    	path to input ROOT file name
//	  -o string
//	    	path to output CSV file name (default "output.csv")
//	  -t string
//	    	name of the tree or graph to convert (default "tree")
//
// By default, root2csv will write out a CSV file with ';' as a column delimiter.
// root2csv ignores the branches of the TTree that are not supported by CSV:
//   - slices/arrays
//   - C++ objects
//
// Example:
//
//	$> root2csv -o out.csv -t tree -f testdata/small-flat-tree.root
//	$> head out.csv
//	## Automatically generated from "testdata/small-flat-tree.root"
//	Int32;Int64;UInt32;UInt64;Float32;Float64;Str;N
//	0;0;0;0;0;0;evt-000;0
//	1;1;1;1;1;1;evt-001;1
//	2;2;2;2;2;2;evt-002;2
//	3;3;3;3;3;3;evt-003;3
//	4;4;4;4;4;4;evt-004;4
//	5;5;5;5;5;5;evt-005;5
//	6;6;6;6;6;6;evt-006;6
//	7;7;7;7;7;7;evt-007;7
package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"strings"

	"go-hep.org/x/hep/csvutil"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("root2csv: ")
	log.SetFlags(0)

	fname := flag.String("f", "", "path to input ROOT file name")
	oname := flag.String("o", "output.csv", "path to output CSV file name")
	tname := flag.String("t", "tree", "name of the tree or graph to convert")

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing input ROOT filename argument")
	}

	err := process(*oname, *fname, *tname)
	if err != nil {
		log.Fatal(err)
	}
}

func process(oname, fname, tname string) error {

	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ROOT file: %w", err)
	}
	defer f.Close()

	obj, err := riofs.Dir(f).Get(tname)
	if err != nil {
		return fmt.Errorf("could not get ROOT object: %w", err)
	}

	switch obj := obj.(type) {
	case rtree.Tree:
		return processTree(oname, fname, obj)
	case rhist.GraphErrors: // Note: test rhist.GraphErrors before rhist.Graph
		return processGraphErrors(oname, fname, obj)
	case rhist.Graph:
		return processGraph(oname, fname, obj)
	default:
		return fmt.Errorf("object %q in file %q is not a rtree.Tree nor a rhist.Graph", tname, fname)
	}
}

func processTree(oname, fname string, tree rtree.Tree) error {
	var nt = ntuple{n: tree.Entries()}
	log.Printf("scanning leaves...")
	for _, leaf := range tree.Leaves() {
		kind := leaf.Type().Kind()
		switch kind {
		case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
			log.Printf(">>> %q %v not supported (%v)", leaf.Name(), leaf.Class(), kind)
			continue
		case reflect.String:
			// ok
		case reflect.Bool,
			reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			if leaf.LeafCount() != nil {
				log.Printf(">>> %q %v not supported (slice)", leaf.Name(), leaf.Class())
				continue
			}
			if leaf.Len() > 1 {
				log.Printf(">>> %q %v not supported (array)", leaf.Name(), leaf.Class())
				continue
			}
		default:
			log.Printf(">>> %q %v not supported (%v) (unknown!)", leaf.Name(), leaf.Class(), kind)
			continue
		}

		nt.add(leaf.Name(), leaf)
	}
	log.Printf("scanning leaves... [done]")

	r, err := rtree.NewReader(tree, nt.args)
	if err != nil {
		return fmt.Errorf("could not create tree reader: %w", err)
	}
	defer r.Close()

	nrows := 0
	err = r.Read(func(ctx rtree.RCtx) error {
		nt.fill()
		nrows++
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not read tree: %w", err)
	}

	tbl, err := csvutil.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output CSV file: %w", err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	names := make([]string, len(nt.cols))
	for i, col := range nt.cols {
		names[i] = col.name
	}
	err = tbl.WriteHeader(fmt.Sprintf(
		"## Automatically generated from %q\n%s\n",
		fname,
		strings.Join(names, string(tbl.Writer.Comma)),
	))
	if err != nil {
		return fmt.Errorf("could not write CSV header: %w", err)
	}

	row := make([]interface{}, len(nt.cols))
	for irow := 0; irow < nrows; irow++ {
		for i, col := range nt.cols {
			row[i] = col.slice.Index(irow).Interface()
		}
		err = tbl.WriteRow(row...)
		if err != nil {
			return fmt.Errorf("could not write row %d to CSV file: %w", irow, err)
		}
	}

	err = tbl.Close()
	if err != nil {
		return fmt.Errorf("could not close CSV output file: %w", err)
	}

	return nil
}

type ntuple struct {
	n    int64
	cols []column
	args []rtree.ReadVar
	vars []interface{}
}

func (nt *ntuple) add(name string, leaf rtree.Leaf) {
	n := len(nt.cols)
	nt.cols = append(nt.cols, newColumn(name, leaf, nt.n))
	col := &nt.cols[n]
	nt.args = append(nt.args, rtree.ReadVar{
		Name:  name,
		Leaf:  leaf.Name(),
		Value: col.data.Addr().Interface(),
	})
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

func processGraph(oname, fname string, g rhist.Graph) error {
	names := []string{"x", "y"}

	tbl, err := csvutil.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output CSV file: %w", err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	err = tbl.WriteHeader(fmt.Sprintf(
		"## Automatically generated from %q\n%s\n",
		fname,
		strings.Join(names, string(tbl.Writer.Comma)),
	))
	if err != nil {
		return fmt.Errorf("could not write CSV header: %w", err)
	}

	n := g.Len()
	for i := 0; i < n; i++ {
		var (
			x, y = g.XY(i)
		)
		err = tbl.WriteRow(x, y)
		if err != nil {
			return fmt.Errorf("could not write row %d to CSV file: %w", i, err)
		}
	}

	err = tbl.Close()
	if err != nil {
		return fmt.Errorf("could not close CSV output file: %w", err)
	}

	return nil
}

func processGraphErrors(oname, fname string, g rhist.GraphErrors) error {
	names := []string{"x", "y", "ex-lo", "ex-hi", "ey-lo", "ey-hi"}

	tbl, err := csvutil.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output CSV file: %w", err)
	}
	defer tbl.Close()
	tbl.Writer.Comma = ';'

	err = tbl.WriteHeader(fmt.Sprintf(
		"## Automatically generated from %q\n%s\n",
		fname,
		strings.Join(names, string(tbl.Writer.Comma)),
	))
	if err != nil {
		return fmt.Errorf("could not write CSV header: %w", err)
	}

	n := g.Len()
	for i := 0; i < n; i++ {
		var (
			x, y     = g.XY(i)
			xlo, xhi = g.XError(i)
			ylo, yhi = g.YError(i)
		)
		err = tbl.WriteRow(x, y, xlo, xhi, ylo, yhi)
		if err != nil {
			return fmt.Errorf("could not write row %d to CSV file: %w", i, err)
		}
	}

	err = tbl.Close()
	if err != nil {
		return fmt.Errorf("could not close CSV output file: %w", err)
	}

	return nil
}
