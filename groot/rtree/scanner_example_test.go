// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"io"
	"log"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleTreeScanner() {
	log.SetPrefix("groot: ")
	log.SetFlags(0)

	f, err := riofs.Open("../testdata/small-flat-tree.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)

	// like for the encoding/json package, struct fields need to
	// be exported to be properly handled by rtree.Scanner.
	// Thus, if the ROOT branch name is lower-case, use the "groot"
	// struct-tag like shown below.
	type Data struct {
		I64    int64       `groot:"Int64"`
		F64    float64     `groot:"Float64"`
		Str    string      `groot:"Str"`
		ArrF64 [10]float64 `groot:"ArrayFloat64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliceFloat64"`
	}

	sc, err := rtree.NewTreeScanner(tree, &Data{})
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

		fmt.Printf("entry[%d]: %+v\n", sc.Entry(), data)
		if sc.Entry() == 9 {
			break
		}
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		log.Fatal(err)
	}

	// Output:
	// entry[0]: {I64:0 F64:0 Str:evt-000 ArrF64:[0 0 0 0 0 0 0 0 0 0] N:0 SliF64:[]}
	// entry[1]: {I64:1 F64:1 Str:evt-001 ArrF64:[1 1 1 1 1 1 1 1 1 1] N:1 SliF64:[1]}
	// entry[2]: {I64:2 F64:2 Str:evt-002 ArrF64:[2 2 2 2 2 2 2 2 2 2] N:2 SliF64:[2 2]}
	// entry[3]: {I64:3 F64:3 Str:evt-003 ArrF64:[3 3 3 3 3 3 3 3 3 3] N:3 SliF64:[3 3 3]}
	// entry[4]: {I64:4 F64:4 Str:evt-004 ArrF64:[4 4 4 4 4 4 4 4 4 4] N:4 SliF64:[4 4 4 4]}
	// entry[5]: {I64:5 F64:5 Str:evt-005 ArrF64:[5 5 5 5 5 5 5 5 5 5] N:5 SliF64:[5 5 5 5 5]}
	// entry[6]: {I64:6 F64:6 Str:evt-006 ArrF64:[6 6 6 6 6 6 6 6 6 6] N:6 SliF64:[6 6 6 6 6 6]}
	// entry[7]: {I64:7 F64:7 Str:evt-007 ArrF64:[7 7 7 7 7 7 7 7 7 7] N:7 SliF64:[7 7 7 7 7 7 7]}
	// entry[8]: {I64:8 F64:8 Str:evt-008 ArrF64:[8 8 8 8 8 8 8 8 8 8] N:8 SliF64:[8 8 8 8 8 8 8 8]}
	// entry[9]: {I64:9 F64:9 Str:evt-009 ArrF64:[9 9 9 9 9 9 9 9 9 9] N:9 SliF64:[9 9 9 9 9 9 9 9 9]}
}

func ExampleTreeScanner_withVars() {
	log.SetPrefix("groot: ")
	log.SetFlags(0)

	f, err := riofs.Open("../testdata/small-flat-tree.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)

	rvars := []rtree.ReadVar{
		{Name: "Int64"},
		{Name: "Float64"},
		{Name: "Str"},
		{Name: "ArrayFloat64"},
		{Name: "N"},
		{Name: "SliceFloat64"},
	}
	sc, err := rtree.NewTreeScannerVars(tree, rvars...)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		var (
			i64 int64
			f64 float64
			str string
			arr [10]float64
			n   int32
			sli []float64
		)
		err := sc.Scan(&i64, &f64, &str, &arr, &n, &sli)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"entry[%d]: i64=%v f64=%v str=%q arr=%v n=%d sli=%v\n",
			sc.Entry(),
			i64, f64, str, arr, n, sli,
		)
		if sc.Entry() == 9 {
			break
		}
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		log.Fatal(err)
	}

	// Output:
	// entry[0]: i64=0 f64=0 str="evt-000" arr=[0 0 0 0 0 0 0 0 0 0] n=0 sli=[]
	// entry[1]: i64=1 f64=1 str="evt-001" arr=[1 1 1 1 1 1 1 1 1 1] n=1 sli=[1]
	// entry[2]: i64=2 f64=2 str="evt-002" arr=[2 2 2 2 2 2 2 2 2 2] n=2 sli=[2 2]
	// entry[3]: i64=3 f64=3 str="evt-003" arr=[3 3 3 3 3 3 3 3 3 3] n=3 sli=[3 3 3]
	// entry[4]: i64=4 f64=4 str="evt-004" arr=[4 4 4 4 4 4 4 4 4 4] n=4 sli=[4 4 4 4]
	// entry[5]: i64=5 f64=5 str="evt-005" arr=[5 5 5 5 5 5 5 5 5 5] n=5 sli=[5 5 5 5 5]
	// entry[6]: i64=6 f64=6 str="evt-006" arr=[6 6 6 6 6 6 6 6 6 6] n=6 sli=[6 6 6 6 6 6]
	// entry[7]: i64=7 f64=7 str="evt-007" arr=[7 7 7 7 7 7 7 7 7 7] n=7 sli=[7 7 7 7 7 7 7]
	// entry[8]: i64=8 f64=8 str="evt-008" arr=[8 8 8 8 8 8 8 8 8 8] n=8 sli=[8 8 8 8 8 8 8 8]
	// entry[9]: i64=9 f64=9 str="evt-009" arr=[9 9 9 9 9 9 9 9 9 9] n=9 sli=[9 9 9 9 9 9 9 9 9]
}

func ExampleScanner_withStruct() {
	log.SetPrefix("groot: ")
	log.SetFlags(0)

	f, err := riofs.Open("../testdata/small-flat-tree.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)

	// like for the encoding/json package, struct fields need to
	// be exported to be properly handled by rtree.Scanner.
	// Thus, if the ROOT branch name is lower-case, use the "groot"
	// struct-tag like shown below.
	type Data struct {
		I64    int64       `groot:"Int64"`
		F64    float64     `groot:"Float64"`
		Str    string      `groot:"Str"`
		ArrF64 [10]float64 `groot:"ArrayFloat64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliceFloat64"`
	}

	var data Data
	sc, err := rtree.NewScanner(tree, &data)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		err := sc.Scan()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("entry[%d]: %+v\n", sc.Entry(), data)
		if sc.Entry() == 9 {
			break
		}
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		log.Fatal(err)
	}

	// Output:
	// entry[0]: {I64:0 F64:0 Str:evt-000 ArrF64:[0 0 0 0 0 0 0 0 0 0] N:0 SliF64:[]}
	// entry[1]: {I64:1 F64:1 Str:evt-001 ArrF64:[1 1 1 1 1 1 1 1 1 1] N:1 SliF64:[1]}
	// entry[2]: {I64:2 F64:2 Str:evt-002 ArrF64:[2 2 2 2 2 2 2 2 2 2] N:2 SliF64:[2 2]}
	// entry[3]: {I64:3 F64:3 Str:evt-003 ArrF64:[3 3 3 3 3 3 3 3 3 3] N:3 SliF64:[3 3 3]}
	// entry[4]: {I64:4 F64:4 Str:evt-004 ArrF64:[4 4 4 4 4 4 4 4 4 4] N:4 SliF64:[4 4 4 4]}
	// entry[5]: {I64:5 F64:5 Str:evt-005 ArrF64:[5 5 5 5 5 5 5 5 5 5] N:5 SliF64:[5 5 5 5 5]}
	// entry[6]: {I64:6 F64:6 Str:evt-006 ArrF64:[6 6 6 6 6 6 6 6 6 6] N:6 SliF64:[6 6 6 6 6 6]}
	// entry[7]: {I64:7 F64:7 Str:evt-007 ArrF64:[7 7 7 7 7 7 7 7 7 7] N:7 SliF64:[7 7 7 7 7 7 7]}
	// entry[8]: {I64:8 F64:8 Str:evt-008 ArrF64:[8 8 8 8 8 8 8 8 8 8] N:8 SliF64:[8 8 8 8 8 8 8 8]}
	// entry[9]: {I64:9 F64:9 Str:evt-009 ArrF64:[9 9 9 9 9 9 9 9 9 9] N:9 SliF64:[9 9 9 9 9 9 9 9 9]}
}

func ExampleScanner_withVars() {
	log.SetPrefix("groot: ")
	log.SetFlags(0)

	f, err := riofs.Open("../testdata/small-flat-tree.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := obj.(rtree.Tree)

	var (
		i64 int64
		f64 float64
		str string
		arr [10]float64
		n   int32
		sli []float64
	)
	rvars := []rtree.ReadVar{
		{Name: "Int64", Value: &i64},
		{Name: "Float64", Value: &f64},
		{Name: "Str", Value: &str},
		{Name: "ArrayFloat64", Value: &arr},
		{Name: "N", Value: &n},
		{Name: "SliceFloat64", Value: &sli},
	}
	sc, err := rtree.NewScannerVars(tree, rvars...)
	if err != nil {
		log.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		err := sc.Scan()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf(
			"entry[%d]: i64=%v f64=%v str=%q arr=%v n=%d sli=%v\n",
			sc.Entry(),
			i64, f64, str, arr, n, sli,
		)
		if sc.Entry() == 9 {
			break
		}
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		log.Fatal(err)
	}

	// Output:
	// entry[0]: i64=0 f64=0 str="evt-000" arr=[0 0 0 0 0 0 0 0 0 0] n=0 sli=[]
	// entry[1]: i64=1 f64=1 str="evt-001" arr=[1 1 1 1 1 1 1 1 1 1] n=1 sli=[1]
	// entry[2]: i64=2 f64=2 str="evt-002" arr=[2 2 2 2 2 2 2 2 2 2] n=2 sli=[2 2]
	// entry[3]: i64=3 f64=3 str="evt-003" arr=[3 3 3 3 3 3 3 3 3 3] n=3 sli=[3 3 3]
	// entry[4]: i64=4 f64=4 str="evt-004" arr=[4 4 4 4 4 4 4 4 4 4] n=4 sli=[4 4 4 4]
	// entry[5]: i64=5 f64=5 str="evt-005" arr=[5 5 5 5 5 5 5 5 5 5] n=5 sli=[5 5 5 5 5]
	// entry[6]: i64=6 f64=6 str="evt-006" arr=[6 6 6 6 6 6 6 6 6 6] n=6 sli=[6 6 6 6 6 6]
	// entry[7]: i64=7 f64=7 str="evt-007" arr=[7 7 7 7 7 7 7 7 7 7] n=7 sli=[7 7 7 7 7 7 7]
	// entry[8]: i64=8 f64=8 str="evt-008" arr=[8 8 8 8 8 8 8 8 8 8] n=8 sli=[8 8 8 8 8 8 8 8]
	// entry[9]: i64=9 f64=9 str="evt-009" arr=[9 9 9 9 9 9 9 9 9 9] n=9 sli=[9 9 9 9 9 9 9 9 9]
}
