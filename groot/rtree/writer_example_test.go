// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"bytes"
	"compress/flate"
	"fmt"
	"log"
	"reflect"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/groot/rdict"
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

		var data Data
		r, err := rtree.NewReader(tree, rtree.ReadVarsFromStruct(&data))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		err = r.Read(func(ctx rtree.RCtx) error {
			fmt.Printf("entry[%d]: %+v\n", ctx.Entry, data)
			return nil
		})

		if err != nil {
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

		var data Data
		r, err := rtree.NewReader(tree, rtree.ReadVarsFromStruct(&data))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		err = r.Read(func(ctx rtree.RCtx) error {
			fmt.Printf("entry[%d]: %+v\n", ctx.Entry, data)
			return nil
		})

		if err != nil {
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

		var data Data
		r, err := rtree.NewReader(tree, rtree.ReadVarsFromStruct(&data))
		if err != nil {
			log.Fatal(err)
		}
		defer r.Close()

		err = r.Read(func(ctx rtree.RCtx) error {
			fmt.Printf("entry[%d]: %+v\n", ctx.Entry, data)
			return nil
		})

		if err != nil {
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

func Example_createEventNtupleNoSplit() {
	type P4 struct {
		Px float64 `groot:"px"`
		Py float64 `groot:"py"`
		Pz float64 `groot:"pz"`
		E  float64 `groot:"ene"`
	}

	type Particle struct {
		ID int32 `groot:"id"`
		P4 P4    `groot:"p4"`
	}

	type Event struct {
		I32 int32      `groot:"i32"`
		F64 float64    `groot:"f64"`
		Str string     `groot:"str"`
		Arr [5]float64 `groot:"arr"`
		Sli []float64  `groot:"sli"`
		P4  P4         `groot:"p4"`
		Ps  []Particle `groot:"mc"`
	}

	// register streamers
	for _, typ := range []reflect.Type{
		reflect.TypeOf(P4{}),
		reflect.TypeOf(Particle{}),
		reflect.TypeOf(Event{}),
	} {

		rdict.StreamerInfos.Add(rdict.StreamerOf(
			rdict.StreamerInfos,
			typ,
		))
	}

	const (
		fname = "../testdata/groot-event-ntuple-nosplit.root"
		nevts = 5
	)

	func() {
		f, err := groot.Create(fname)
		if err != nil {
			log.Fatalf("could not create ROOT file: %+v", err)
		}
		defer f.Close()

		var (
			evt   Event
			wvars = []rtree.WriteVar{
				{Name: "evt", Value: &evt},
			}
		)

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
			evt.Arr = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			evt.Sli = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
			evt.P4 = P4{Px: float64(i), Py: float64(i + 1), Pz: float64(i + 2), E: float64(i + 3)}
			evt.Ps = []Particle{
				{ID: int32(i), P4: evt.P4},
				{ID: int32(i + 1), P4: evt.P4},
				{ID: int32(i + 2), P4: evt.P4},
				{ID: int32(i + 3), P4: evt.P4},
				{ID: int32(i + 4), P4: evt.P4},
			}[:i]

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

	{
		fmt.Printf("-- read back ROOT file\n")
		out := new(bytes.Buffer)
		err := rcmd.List(out, fname, rcmd.ListTrees(true))
		if err != nil {
			log.Fatalf("could not list ROOT file content: %+v", err)
		}
		fmt.Printf("%s\n", out.String())
	}

	{
		var (
			deep   = true
			filter = func(name string) bool { return true }
		)
		out := new(bytes.Buffer)
		err := rcmd.Dump(out, fname, deep, filter)
		if err != nil {
			log.Fatalf("could not dump ROOT file: %+v", err)
		}
		fmt.Printf("%s\n", out.String())
	}

	// Output:
	// -- created tree "mytree":
	// branch[0]: name="evt", title="evt"
	// -- filled tree with 5 entries
	// -- read back ROOT file
	// === [../testdata/groot-event-ntuple-nosplit.root] ===
	// version: 62406
	//   TTree mytree          (entries=5)
	//     evt "evt"   TBranchElement
	//
	// key[000]: mytree;1 "" (TTree)
	// [000][evt]: {0 0 evt-0 [0 1 2 3 4] [] {0 1 2 3} []}
	// [001][evt]: {1 1 evt-1 [1 2 3 4 5] [1] {1 2 3 4} [{1 {1 2 3 4}}]}
	// [002][evt]: {2 2 evt-2 [2 3 4 5 6] [2 3] {2 3 4 5} [{2 {2 3 4 5}} {3 {2 3 4 5}}]}
	// [003][evt]: {3 3 evt-3 [3 4 5 6 7] [3 4 5] {3 4 5 6} [{3 {3 4 5 6}} {4 {3 4 5 6}} {5 {3 4 5 6}}]}
	// [004][evt]: {4 4 evt-4 [4 5 6 7 8] [4 5 6 7] {4 5 6 7} [{4 {4 5 6 7}} {5 {4 5 6 7}} {6 {4 5 6 7}} {7 {4 5 6 7}}]}

}

func Example_createEventNtupleFullSplit() {
	type P4 struct {
		Px float64 `groot:"px"`
		Py float64 `groot:"py"`
		Pz float64 `groot:"pz"`
		E  float64 `groot:"ene"`
	}

	type Particle struct {
		ID int32 `groot:"id"`
		P4 P4    `groot:"p4"`
	}

	type Event struct {
		I32 int32      `groot:"i32"`
		F64 float64    `groot:"f64"`
		Str string     `groot:"str"`
		Arr [5]float64 `groot:"arr"`
		Sli []float64  `groot:"sli"`
		P4  P4         `groot:"p4"`
		Ps  []Particle `groot:"mc"`
	}

	// register streamers
	for _, typ := range []reflect.Type{
		reflect.TypeOf(P4{}),
		reflect.TypeOf(Particle{}),
		reflect.TypeOf(Event{}),
	} {

		rdict.StreamerInfos.Add(rdict.StreamerOf(
			rdict.StreamerInfos,
			typ,
		))
	}

	const (
		fname = "../testdata/groot-event-ntuple-fullsplit.root"
		nevts = 5
	)

	func() {
		f, err := groot.Create(fname)
		if err != nil {
			log.Fatalf("could not create ROOT file: %+v", err)
		}
		defer f.Close()

		var (
			evt   Event
			wvars = rtree.WriteVarsFromStruct(&evt)
		)

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
			evt.Arr = [5]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}
			evt.Sli = []float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3), float64(i + 4)}[:i]
			evt.P4 = P4{Px: float64(i), Py: float64(i + 1), Pz: float64(i + 2), E: float64(i + 3)}
			evt.Ps = []Particle{
				{ID: int32(i), P4: evt.P4},
				{ID: int32(i + 1), P4: evt.P4},
				{ID: int32(i + 2), P4: evt.P4},
				{ID: int32(i + 3), P4: evt.P4},
				{ID: int32(i + 4), P4: evt.P4},
			}[:i]

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

	{
		fmt.Printf("-- read back ROOT file\n")
		out := new(bytes.Buffer)
		err := rcmd.List(out, fname, rcmd.ListTrees(true))
		if err != nil {
			log.Fatalf("could not list ROOT file content: %+v", err)
		}
		fmt.Printf("%s\n", out.String())
	}

	{
		var (
			deep   = true
			filter = func(name string) bool { return true }
		)
		out := new(bytes.Buffer)
		err := rcmd.Dump(out, fname, deep, filter)
		if err != nil {
			log.Fatalf("could not dump ROOT file: %+v", err)
		}
		fmt.Printf("%s\n", out.String())
	}

	// Output:
	// -- created tree "mytree":
	// branch[0]: name="i32", title="i32/I"
	// branch[1]: name="f64", title="f64/D"
	// branch[2]: name="str", title="str/C"
	// branch[3]: name="arr", title="arr[5]/D"
	// branch[4]: name="sli", title="sli"
	// branch[5]: name="p4", title="p4"
	// branch[6]: name="mc", title="mc"
	// -- filled tree with 5 entries
	// -- read back ROOT file
	// === [../testdata/groot-event-ntuple-fullsplit.root] ===
	// version: 62406
	//   TTree mytree             (entries=5)
	//     i32 "i32/I"    TBranch
	//     f64 "f64/D"    TBranch
	//     str "str/C"    TBranch
	//     arr "arr[5]/D" TBranch
	//     sli "sli"      TBranchElement
	//     p4  "p4"       TBranchElement
	//     mc  "mc"       TBranchElement
	//
	// key[000]: mytree;1 "" (TTree)
	// [000][i32]: 0
	// [000][f64]: 0
	// [000][str]: evt-0
	// [000][arr]: [0 1 2 3 4]
	// [000][sli]: []
	// [000][p4]: {0 1 2 3}
	// [000][mc]: []
	// [001][i32]: 1
	// [001][f64]: 1
	// [001][str]: evt-1
	// [001][arr]: [1 2 3 4 5]
	// [001][sli]: [1]
	// [001][p4]: {1 2 3 4}
	// [001][mc]: [{1 {1 2 3 4}}]
	// [002][i32]: 2
	// [002][f64]: 2
	// [002][str]: evt-2
	// [002][arr]: [2 3 4 5 6]
	// [002][sli]: [2 3]
	// [002][p4]: {2 3 4 5}
	// [002][mc]: [{2 {2 3 4 5}} {3 {2 3 4 5}}]
	// [003][i32]: 3
	// [003][f64]: 3
	// [003][str]: evt-3
	// [003][arr]: [3 4 5 6 7]
	// [003][sli]: [3 4 5]
	// [003][p4]: {3 4 5 6}
	// [003][mc]: [{3 {3 4 5 6}} {4 {3 4 5 6}} {5 {3 4 5 6}}]
	// [004][i32]: 4
	// [004][f64]: 4
	// [004][str]: evt-4
	// [004][arr]: [4 5 6 7 8]
	// [004][sli]: [4 5 6 7]
	// [004][p4]: {4 5 6 7}
	// [004][mc]: [{4 {4 5 6 7}} {5 {4 5 6 7}} {6 {4 5 6 7}} {7 {4 5 6 7}}]
}
