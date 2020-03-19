// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"compress/flate"
	"fmt"
	"io"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func Example_createFlatNtuple() {
	type Data struct {
		I32    int32
		F64    float64
		Str    string
		ArrF64 [5]float64
		N      int32
		SliF64 []float64
	}
	const (
		fname = "../testdata/groot-flat-ntuple.root"
		nevts = 5
	)
	func() {
		f, err := groot.Create(fname)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		defer f.Close()

		var evt Data

		wvars := []rtree.WriteVar{
			{Name: "I32", Value: &evt.I32},
			{Name: "F64", Value: &evt.F64},
			{Name: "Str", Value: &evt.Str},
			{Name: "ArrF64", Value: &evt.ArrF64},
			{Name: "N", Value: &evt.N},
			{Name: "SliF64", Value: &evt.SliF64, Count: "N"},
		}
		tree, err := rtree.NewWriter(f, "mytree", wvars)
		if err != nil {
			log.Fatalf("could not create tree writer: %+v", err)
		}
		defer tree.Close()

		fmt.Printf("-- created tree %q:\n", tree.Name())
		for i, b := range tree.Branches() {
			fmt.Printf("branch[%d]: name=%q, title=%q\n", i, b.Name(), b.Title())
		}

		for i := 0; i < nevts; i++ {
			evt.I32 = int32(i)
			evt.F64 = float64(i)
			evt.Str = fmt.Sprintf("evt-%0d", i)
			evt.ArrF64 = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			evt.N = int32(i)
			evt.SliF64 = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
			_, err = tree.Write()
			if err != nil {
				log.Fatalf("could not write event %d: %+v", i, err)
			}
		}
		fmt.Printf("-- filled tree with %d entries\n", tree.Entries())

		err = tree.Close()
		if err != nil {
			log.Fatalf("could not write tree: %+v", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("could not close tree: %+v", err)
		}
	}()

	func() {
		fmt.Printf("-- read back ROOT file\n")
		f, err := groot.Open(fname)
		if err != nil {
			log.Fatalf("could not open ROOT file: %+v", err)
		}
		defer f.Close()

		obj, err := f.Get("mytree")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		tree := obj.(rtree.Tree)

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
	}()

	// Output:
	// -- created tree "mytree":
	// branch[0]: name="I32", title="I32/I"
	// branch[1]: name="F64", title="F64/D"
	// branch[2]: name="Str", title="Str/C"
	// branch[3]: name="ArrF64", title="ArrF64[5]/D"
	// branch[4]: name="N", title="N/I"
	// branch[5]: name="SliF64", title="SliF64[N]/D"
	// -- filled tree with 5 entries
	// -- read back ROOT file
	// entry[0]: {I32:0 F64:0 Str:evt-0 ArrF64:[0 1 2 3 4] N:0 SliF64:[]}
	// entry[1]: {I32:1 F64:1 Str:evt-1 ArrF64:[1 2 3 4 5] N:1 SliF64:[1]}
	// entry[2]: {I32:2 F64:2 Str:evt-2 ArrF64:[2 3 4 5 6] N:2 SliF64:[2 3]}
	// entry[3]: {I32:3 F64:3 Str:evt-3 ArrF64:[3 4 5 6 7] N:3 SliF64:[3 4 5]}
	// entry[4]: {I32:4 F64:4 Str:evt-4 ArrF64:[4 5 6 7 8] N:4 SliF64:[4 5 6 7]}
}

func Example_createFlatNtupleWithLZMA() {
	type Data struct {
		I32    int32
		F64    float64
		Str    string
		ArrF64 [5]float64
		N      int32
		SliF64 []float64
	}
	const (
		fname = "../testdata/groot-flat-ntuple-with-lzma.root"
		nevts = 5
	)
	func() {
		f, err := groot.Create(fname)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		defer f.Close()

		var evt Data

		wvars := []rtree.WriteVar{
			{Name: "I32", Value: &evt.I32},
			{Name: "F64", Value: &evt.F64},
			{Name: "Str", Value: &evt.Str},
			{Name: "ArrF64", Value: &evt.ArrF64},
			{Name: "N", Value: &evt.N},
			{Name: "SliF64", Value: &evt.SliF64, Count: "N"},
		}
		tree, err := rtree.NewWriter(f, "mytree", wvars, rtree.WithLZMA(flate.BestCompression), rtree.WithBasketSize(32*1024))
		if err != nil {
			log.Fatalf("could not create tree writer: %+v", err)
		}
		defer tree.Close()

		fmt.Printf("-- created tree %q:\n", tree.Name())
		for i, b := range tree.Branches() {
			fmt.Printf("branch[%d]: name=%q, title=%q\n", i, b.Name(), b.Title())
		}

		for i := 0; i < nevts; i++ {
			evt.I32 = int32(i)
			evt.F64 = float64(i)
			evt.Str = fmt.Sprintf("evt-%0d", i)
			evt.ArrF64 = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			evt.N = int32(i)
			evt.SliF64 = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
			_, err = tree.Write()
			if err != nil {
				log.Fatalf("could not write event %d: %+v", i, err)
			}
		}
		fmt.Printf("-- filled tree with %d entries\n", tree.Entries())

		err = tree.Close()
		if err != nil {
			log.Fatalf("could not write tree: %+v", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("could not close tree: %+v", err)
		}
	}()

	func() {
		fmt.Printf("-- read back ROOT file\n")
		f, err := groot.Open(fname)
		if err != nil {
			log.Fatalf("could not open ROOT file: %+v", err)
		}
		defer f.Close()

		obj, err := f.Get("mytree")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		tree := obj.(rtree.Tree)

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
	}()

	// Output:
	// -- created tree "mytree":
	// branch[0]: name="I32", title="I32/I"
	// branch[1]: name="F64", title="F64/D"
	// branch[2]: name="Str", title="Str/C"
	// branch[3]: name="ArrF64", title="ArrF64[5]/D"
	// branch[4]: name="N", title="N/I"
	// branch[5]: name="SliF64", title="SliF64[N]/D"
	// -- filled tree with 5 entries
	// -- read back ROOT file
	// entry[0]: {I32:0 F64:0 Str:evt-0 ArrF64:[0 1 2 3 4] N:0 SliF64:[]}
	// entry[1]: {I32:1 F64:1 Str:evt-1 ArrF64:[1 2 3 4 5] N:1 SliF64:[1]}
	// entry[2]: {I32:2 F64:2 Str:evt-2 ArrF64:[2 3 4 5 6] N:2 SliF64:[2 3]}
	// entry[3]: {I32:3 F64:3 Str:evt-3 ArrF64:[3 4 5 6 7] N:3 SliF64:[3 4 5]}
	// entry[4]: {I32:4 F64:4 Str:evt-4 ArrF64:[4 5 6 7 8] N:4 SliF64:[4 5 6 7]}
}

func Example_createFlatNtupleFromStruct() {
	type Data struct {
		I32    int32
		F64    float64
		Str    string
		ArrF64 [5]float64
		N      int32
		SliF64 []float64 `groot:"SliF64[N]"`
	}
	const (
		fname = "../testdata/groot-flat-ntuple-with-struct.root"
		nevts = 5
	)
	func() {
		f, err := groot.Create(fname)
		if err != nil {
			log.Fatalf("%+v", err)
		}
		defer f.Close()

		var evt Data

		tree, err := rtree.NewWriter(f, "mytree", rtree.WriteVarsFromStruct(&evt))
		if err != nil {
			log.Fatalf("could not create tree writer: %+v", err)
		}
		defer tree.Close()

		fmt.Printf("-- created tree %q:\n", tree.Name())
		for i, b := range tree.Branches() {
			fmt.Printf("branch[%d]: name=%q, title=%q\n", i, b.Name(), b.Title())
		}

		for i := 0; i < nevts; i++ {
			evt.I32 = int32(i)
			evt.F64 = float64(i)
			evt.Str = fmt.Sprintf("evt-%0d", i)
			evt.ArrF64 = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			evt.N = int32(i)
			evt.SliF64 = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
			_, err = tree.Write()
			if err != nil {
				log.Fatalf("could not write event %d: %+v", i, err)
			}
		}
		fmt.Printf("-- filled tree with %d entries\n", tree.Entries())

		err = tree.Close()
		if err != nil {
			log.Fatalf("could not write tree: %+v", err)
		}

		err = f.Close()
		if err != nil {
			log.Fatalf("could not close tree: %+v", err)
		}
	}()

	func() {
		fmt.Printf("-- read back ROOT file\n")
		f, err := groot.Open(fname)
		if err != nil {
			log.Fatalf("could not open ROOT file: %+v", err)
		}
		defer f.Close()

		obj, err := f.Get("mytree")
		if err != nil {
			log.Fatalf("%+v", err)
		}

		tree := obj.(rtree.Tree)

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
	}()

	// Output:
	// -- created tree "mytree":
	// branch[0]: name="I32", title="I32/I"
	// branch[1]: name="F64", title="F64/D"
	// branch[2]: name="Str", title="Str/C"
	// branch[3]: name="ArrF64", title="ArrF64[5]/D"
	// branch[4]: name="N", title="N/I"
	// branch[5]: name="SliF64", title="SliF64[N]/D"
	// -- filled tree with 5 entries
	// -- read back ROOT file
	// entry[0]: {I32:0 F64:0 Str:evt-0 ArrF64:[0 1 2 3 4] N:0 SliF64:[]}
	// entry[1]: {I32:1 F64:1 Str:evt-1 ArrF64:[1 2 3 4 5] N:1 SliF64:[1]}
	// entry[2]: {I32:2 F64:2 Str:evt-2 ArrF64:[2 3 4 5 6] N:2 SliF64:[2 3]}
	// entry[3]: {I32:3 F64:3 Str:evt-3 ArrF64:[3 4 5 6 7] N:3 SliF64:[3 4 5]}
	// entry[4]: {I32:4 F64:4 Str:evt-4 ArrF64:[4 5 6 7 8] N:4 SliF64:[4 5 6 7]}
}
