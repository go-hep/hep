// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rarrow"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleRecordReader() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree)
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [1]
	// rec[0][two]: [1.1]
	// rec[0][three]: ["uno"]
	// rec[1][one]: [2]
	// rec[1][two]: [2.2]
	// rec[1][three]: ["dos"]
	// rec[2][one]: [3]
	// rec[2][two]: [3.3]
	// rec[2][three]: ["tres"]
	// rec[3][one]: [4]
	// rec[3][two]: [4.4]
	// rec[3][three]: ["quatro"]
}

func ExampleRecordReader_withChunk() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree, rarrow.WithChunk(3))
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [1 2 3]
	// rec[0][two]: [1.1 2.2 3.3]
	// rec[0][three]: ["uno" "dos" "tres"]
	// rec[1][one]: [4]
	// rec[1][two]: [4.4]
	// rec[1][three]: ["quatro"]
}

func ExampleRecordReader_allTree() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree, rarrow.WithChunk(-1))
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [1 2 3 4]
	// rec[0][two]: [1.1 2.2 3.3 4.4]
	// rec[0][three]: ["uno" "dos" "tres" "quatro"]
}

func ExampleRecordReader_withStart() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree, rarrow.WithStart(1))
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [2]
	// rec[0][two]: [2.2]
	// rec[0][three]: ["dos"]
	// rec[1][one]: [3]
	// rec[1][two]: [3.3]
	// rec[1][three]: ["tres"]
	// rec[2][one]: [4]
	// rec[2][two]: [4.4]
	// rec[2][three]: ["quatro"]
}

func ExampleRecordReader_withEnd() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree, rarrow.WithEnd(2))
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [1]
	// rec[0][two]: [1.1]
	// rec[0][three]: ["uno"]
	// rec[1][one]: [2]
	// rec[1][two]: [2.2]
	// rec[1][three]: ["dos"]
}

func ExampleRecordReader_withStartEnd() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := riofs.Dir(f).Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	rr := rarrow.NewRecordReader(tree, rarrow.WithStart(1), rarrow.WithEnd(2))
	defer rr.Release()

	recs := 0
	for rr.Next() {
		rec := rr.Record()
		for i, col := range rec.Columns() {
			fmt.Printf("rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
		}
		recs++
	}

	// Output:
	// rec[0][one]: [2]
	// rec[0][two]: [2.2]
	// rec[0][three]: ["dos"]
}
