// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/rootio"
)

// ExampleChain shows how to create a chain made of 2 trees.
func ExampleChain() {
	const name = "tree"

	f1, err := rootio.Open("testdata/chain.1.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f1.Close()

	o1, err := f1.Get(name)
	if err != nil {
		log.Fatal(err)
	}
	t1 := o1.(rootio.Tree)

	f2, err := rootio.Open("testdata/chain.2.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	o2, err := f2.Get(name)
	if err != nil {
		log.Fatal(err)
	}
	t2 := o2.(rootio.Tree)

	chain := rootio.Chain(t1, t2)

	type Data struct {
		Event struct {
			Beg       string      `rootio:"Beg"`
			F64       float64     `rootio:"F64"`
			ArrF64    [10]float64 `rootio:"ArrayF64"`
			N         int32       `rootio:"N"`
			SliF64    []float64   `rootio:"SliceF64"`
			StdStr    string      `rootio:"StdStr"`
			StlVecF64 []float64   `rootio:"StlVecF64"`
			StlVecStr []string    `rootio:"StlVecStr"`
			End       string      `rootio:"End"`
		} `rootio:"evt"`
	}

	sc, err := rootio.NewTreeScanner(chain, &Data{})
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		var data Data
		err := sc.Scan(&data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("entry[%02d]: beg=%q f64=%v\n", sc.Entry(), data.Event.Beg, data.Event.F64)
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("error during scan: %v", err)
	}

	// Output:
	// entry[00]: beg="beg-000" f64=0
	// entry[01]: beg="beg-001" f64=1
	// entry[02]: beg="beg-002" f64=2
	// entry[03]: beg="beg-003" f64=3
	// entry[04]: beg="beg-004" f64=4
	// entry[05]: beg="beg-005" f64=5
	// entry[06]: beg="beg-006" f64=6
	// entry[07]: beg="beg-007" f64=7
	// entry[08]: beg="beg-008" f64=8
	// entry[09]: beg="beg-009" f64=9
	// entry[10]: beg="beg-010" f64=10
	// entry[11]: beg="beg-011" f64=11
	// entry[12]: beg="beg-012" f64=12
	// entry[13]: beg="beg-013" f64=13
	// entry[14]: beg="beg-014" f64=14
	// entry[15]: beg="beg-015" f64=15
	// entry[16]: beg="beg-016" f64=16
	// entry[17]: beg="beg-017" f64=17
	// entry[18]: beg="beg-018" f64=18
	// entry[19]: beg="beg-019" f64=19
}

// ExampleChainOf shows how to create a chain made of trees from 2 files.
func ExampleChainOf() {
	const name = "tree"

	chain, closer, err := rootio.ChainOf(name, "testdata/chain.1.root", "testdata/chain.2.root")
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	type Data struct {
		Event struct {
			Beg       string      `rootio:"Beg"`
			F64       float64     `rootio:"F64"`
			ArrF64    [10]float64 `rootio:"ArrayF64"`
			N         int32       `rootio:"N"`
			SliF64    []float64   `rootio:"SliceF64"`
			StdStr    string      `rootio:"StdStr"`
			StlVecF64 []float64   `rootio:"StlVecF64"`
			StlVecStr []string    `rootio:"StlVecStr"`
			End       string      `rootio:"End"`
		} `rootio:"evt"`
	}

	sc, err := rootio.NewTreeScanner(chain, &Data{})
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		var data Data
		err := sc.Scan(&data)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("entry[%02d]: beg=%q f64=%v\n", sc.Entry(), data.Event.Beg, data.Event.F64)
	}

	if err := sc.Err(); err != nil {
		log.Fatalf("error during scan: %v", err)
	}

	// Output:
	// entry[00]: beg="beg-000" f64=0
	// entry[01]: beg="beg-001" f64=1
	// entry[02]: beg="beg-002" f64=2
	// entry[03]: beg="beg-003" f64=3
	// entry[04]: beg="beg-004" f64=4
	// entry[05]: beg="beg-005" f64=5
	// entry[06]: beg="beg-006" f64=6
	// entry[07]: beg="beg-007" f64=7
	// entry[08]: beg="beg-008" f64=8
	// entry[09]: beg="beg-009" f64=9
	// entry[10]: beg="beg-010" f64=10
	// entry[11]: beg="beg-011" f64=11
	// entry[12]: beg="beg-012" f64=12
	// entry[13]: beg="beg-013" f64=13
	// entry[14]: beg="beg-014" f64=14
	// entry[15]: beg="beg-015" f64=15
	// entry[16]: beg="beg-016" f64=16
	// entry[17]: beg="beg-017" f64=17
	// entry[18]: beg="beg-018" f64=18
	// entry[19]: beg="beg-019" f64=19
}
