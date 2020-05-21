// Copyright ©2019 The go-hep Authors. All rights reserved.
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

func ExampleRecordReader_withChain() {
	chain, closer, err := rtree.ChainOf("tree", "../testdata/chain.1.root", "../testdata/chain.2.root")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = closer()
		if err != nil {
			log.Fatalf("could not close chain: %+v", err)
		}
	}()

	rr := rarrow.NewRecordReader(chain, rarrow.WithStart(10), rarrow.WithEnd(20))
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
	// rec[0][evt]: {["beg-010"] [10] [[10 10 10 10 10 10 10 10 10 10]] [0] [[]] ["std-010"] [[]] [[]] ["end-010"]}
	// rec[1][evt]: {["beg-011"] [11] [[11 11 11 11 11 11 11 11 11 11]] [1] [[11]] ["std-011"] [[11]] [["vec-011"]] ["end-011"]}
	// rec[2][evt]: {["beg-012"] [12] [[12 12 12 12 12 12 12 12 12 12]] [2] [[12 12]] ["std-012"] [[12 12]] [["vec-012" "vec-012"]] ["end-012"]}
	// rec[3][evt]: {["beg-013"] [13] [[13 13 13 13 13 13 13 13 13 13]] [3] [[13 13 13]] ["std-013"] [[13 13 13]] [["vec-013" "vec-013" "vec-013"]] ["end-013"]}
	// rec[4][evt]: {["beg-014"] [14] [[14 14 14 14 14 14 14 14 14 14]] [4] [[14 14 14 14]] ["std-014"] [[14 14 14 14]] [["vec-014" "vec-014" "vec-014" "vec-014"]] ["end-014"]}
	// rec[5][evt]: {["beg-015"] [15] [[15 15 15 15 15 15 15 15 15 15]] [5] [[15 15 15 15 15]] ["std-015"] [[15 15 15 15 15]] [["vec-015" "vec-015" "vec-015" "vec-015" "vec-015"]] ["end-015"]}
	// rec[6][evt]: {["beg-016"] [16] [[16 16 16 16 16 16 16 16 16 16]] [6] [[16 16 16 16 16 16]] ["std-016"] [[16 16 16 16 16 16]] [["vec-016" "vec-016" "vec-016" "vec-016" "vec-016" "vec-016"]] ["end-016"]}
	// rec[7][evt]: {["beg-017"] [17] [[17 17 17 17 17 17 17 17 17 17]] [7] [[17 17 17 17 17 17 17]] ["std-017"] [[17 17 17 17 17 17 17]] [["vec-017" "vec-017" "vec-017" "vec-017" "vec-017" "vec-017" "vec-017"]] ["end-017"]}
	// rec[8][evt]: {["beg-018"] [18] [[18 18 18 18 18 18 18 18 18 18]] [8] [[18 18 18 18 18 18 18 18]] ["std-018"] [[18 18 18 18 18 18 18 18]] [["vec-018" "vec-018" "vec-018" "vec-018" "vec-018" "vec-018" "vec-018" "vec-018"]] ["end-018"]}
	// rec[9][evt]: {["beg-019"] [19] [[19 19 19 19 19 19 19 19 19 19]] [9] [[19 19 19 19 19 19 19 19 19]] ["std-019"] [[19 19 19 19 19 19 19 19 19]] [["vec-019" "vec-019" "vec-019" "vec-019" "vec-019" "vec-019" "vec-019" "vec-019" "vec-019"]] ["end-019"]}
}
