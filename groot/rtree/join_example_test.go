// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"log"
	"strings"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleJoin() {

	get := func(fname, tname string) (rtree.Tree, func() error) {
		f, err := groot.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		t, err := f.Get(tname)
		if err != nil {
			_ = f.Close()
			log.Fatal(err)
		}
		return t.(rtree.Tree), f.Close
	}
	chk := func(f func() error) {
		err := f()
		if err != nil {
			log.Fatal(err)
		}
	}

	t1, f1 := get("../testdata/join1.root", "j1")
	defer chk(f1)

	t2, f2 := get("../testdata/join2.root", "j2")
	defer chk(f2)

	t3, f3 := get("../testdata/join3.root", "j3")
	defer chk(f3)

	join, err := rtree.Join(t1, t2, t3)
	if err != nil {
		log.Fatalf("could not join trees: %+v", err)
	}

	fmt.Printf("t1:   %s (nevts=%d)\n", t1.Name(), t1.Entries())
	fmt.Printf("t2:   %s (nevts=%d)\n", t2.Name(), t2.Entries())
	fmt.Printf("t3:   %s (nevts=%d)\n", t3.Name(), t3.Entries())
	fmt.Printf("join: %s\n", join.Name())
	fmt.Printf("entries: %d\n", join.Entries())

	rvars := rtree.NewReadVars(join)
	r, err := rtree.NewReader(join, rvars)
	if err != nil {
		log.Fatalf("could not create reader for joined trees: %+v", err)
	}
	defer r.Close()

	rf1, err := r.FormulaFunc(
		[]string{"b10", "b30", "b20"},
		func(b1, b3, b2 float64) float64 {
			return b1 + b2 + b3
		},
	)
	if err != nil {
		log.Fatalf("could not bind formula: %+v", err)
	}
	fct1 := rf1.Func().(func() float64)

	err = r.Read(func(rctx rtree.RCtx) error {
		for _, rv := range rvars {
			fmt.Printf("join[%03d][%s]: %v\n", rctx.Entry, rv.Name, rv.Deref())
		}
		fmt.Printf("join[%03d][fun]: %v\n", rctx.Entry, fct1())
		return nil
	})

	if err != nil {
		log.Fatalf("could not process events: %+v", err)
	}

	// Output:
	// t1:   j1 (nevts=10)
	// t2:   j2 (nevts=10)
	// t3:   j3 (nevts=10)
	// join: join_j1_j2_j3
	// entries: 10
	// join[000][b10]: 101
	// join[000][b11]: 101
	// join[000][b12]: j1-101
	// join[000][b20]: 201
	// join[000][b21]: 201
	// join[000][b22]: j2-201
	// join[000][b30]: 301
	// join[000][b31]: 301
	// join[000][b32]: j3-301
	// join[000][fun]: 603
	// join[001][b10]: 102
	// join[001][b11]: 102
	// join[001][b12]: j1-102
	// join[001][b20]: 202
	// join[001][b21]: 202
	// join[001][b22]: j2-202
	// join[001][b30]: 302
	// join[001][b31]: 302
	// join[001][b32]: j3-302
	// join[001][fun]: 606
	// join[002][b10]: 103
	// join[002][b11]: 103
	// join[002][b12]: j1-103
	// join[002][b20]: 203
	// join[002][b21]: 203
	// join[002][b22]: j2-203
	// join[002][b30]: 303
	// join[002][b31]: 303
	// join[002][b32]: j3-303
	// join[002][fun]: 609
	// join[003][b10]: 104
	// join[003][b11]: 104
	// join[003][b12]: j1-104
	// join[003][b20]: 204
	// join[003][b21]: 204
	// join[003][b22]: j2-204
	// join[003][b30]: 304
	// join[003][b31]: 304
	// join[003][b32]: j3-304
	// join[003][fun]: 612
	// join[004][b10]: 105
	// join[004][b11]: 105
	// join[004][b12]: j1-105
	// join[004][b20]: 205
	// join[004][b21]: 205
	// join[004][b22]: j2-205
	// join[004][b30]: 305
	// join[004][b31]: 305
	// join[004][b32]: j3-305
	// join[004][fun]: 615
	// join[005][b10]: 106
	// join[005][b11]: 106
	// join[005][b12]: j1-106
	// join[005][b20]: 206
	// join[005][b21]: 206
	// join[005][b22]: j2-206
	// join[005][b30]: 306
	// join[005][b31]: 306
	// join[005][b32]: j3-306
	// join[005][fun]: 618
	// join[006][b10]: 107
	// join[006][b11]: 107
	// join[006][b12]: j1-107
	// join[006][b20]: 207
	// join[006][b21]: 207
	// join[006][b22]: j2-207
	// join[006][b30]: 307
	// join[006][b31]: 307
	// join[006][b32]: j3-307
	// join[006][fun]: 621
	// join[007][b10]: 108
	// join[007][b11]: 108
	// join[007][b12]: j1-108
	// join[007][b20]: 208
	// join[007][b21]: 208
	// join[007][b22]: j2-208
	// join[007][b30]: 308
	// join[007][b31]: 308
	// join[007][b32]: j3-308
	// join[007][fun]: 624
	// join[008][b10]: 109
	// join[008][b11]: 109
	// join[008][b12]: j1-109
	// join[008][b20]: 209
	// join[008][b21]: 209
	// join[008][b22]: j2-209
	// join[008][b30]: 309
	// join[008][b31]: 309
	// join[008][b32]: j3-309
	// join[008][fun]: 627
	// join[009][b10]: 110
	// join[009][b11]: 110
	// join[009][b12]: j1-110
	// join[009][b20]: 210
	// join[009][b21]: 210
	// join[009][b22]: j2-210
	// join[009][b30]: 310
	// join[009][b31]: 310
	// join[009][b32]: j3-310
	// join[009][fun]: 630
}

func ExampleJoin_withReadVarSelection() {

	get := func(fname, tname string) (rtree.Tree, func() error) {
		f, err := groot.Open(fname)
		if err != nil {
			log.Fatal(err)
		}
		t, err := f.Get(tname)
		if err != nil {
			_ = f.Close()
			log.Fatal(err)
		}
		return t.(rtree.Tree), f.Close
	}
	chk := func(f func() error) {
		err := f()
		if err != nil {
			log.Fatal(err)
		}
	}

	t1, f1 := get("../testdata/join1.root", "j1")
	defer chk(f1)

	t2, f2 := get("../testdata/join2.root", "j2")
	defer chk(f2)

	t3, f3 := get("../testdata/join3.root", "j3")
	defer chk(f3)

	join, err := rtree.Join(t1, t2, t3)
	if err != nil {
		log.Fatalf("could not join trees: %+v", err)
	}

	fmt.Printf("t1:   %s (nevts=%d)\n", t1.Name(), t1.Entries())
	fmt.Printf("t2:   %s (nevts=%d)\n", t2.Name(), t2.Entries())
	fmt.Printf("t3:   %s (nevts=%d)\n", t3.Name(), t3.Entries())
	fmt.Printf("join: %s\n", join.Name())
	fmt.Printf("entries: %d\n", join.Entries())

	rvars := []rtree.ReadVar{
		{Name: "b10", Value: new(float64)},
		{Name: "b20", Value: new(float64)},
	}

	r, err := rtree.NewReader(join, rvars, rtree.WithRange(3, 8))
	if err != nil {
		log.Fatalf("could not create reader for joined trees: %+v", err)
	}
	defer r.Close()

	rf1, err := r.FormulaFunc(
		[]string{"b12", "b32", "b22"},
		func(b1, b3, b2 string) string {
			return strings.Join([]string{b1, b3, b2}, ", ")
		},
	)
	if err != nil {
		log.Fatalf("could not bind formula: %+v", err)
	}
	fct1 := rf1.Func().(func() string)

	err = r.Read(func(rctx rtree.RCtx) error {
		for _, rv := range rvars {
			fmt.Printf("join[%03d][%s]: %v\n", rctx.Entry, rv.Name, rv.Deref())
		}
		fmt.Printf("join[%03d][fun]: %v\n", rctx.Entry, fct1())
		return nil
	})

	if err != nil {
		log.Fatalf("could not process events: %+v", err)
	}

	// Output:
	// t1:   j1 (nevts=10)
	// t2:   j2 (nevts=10)
	// t3:   j3 (nevts=10)
	// join: join_j1_j2_j3
	// entries: 10
	// join[003][b10]: 104
	// join[003][b20]: 204
	// join[003][fun]: j1-104, j3-304, j2-204
	// join[004][b10]: 105
	// join[004][b20]: 205
	// join[004][fun]: j1-105, j3-305, j2-205
	// join[005][b10]: 106
	// join[005][b20]: 206
	// join[005][fun]: j1-106, j3-306, j2-206
	// join[006][b10]: 107
	// join[006][b20]: 207
	// join[006][fun]: j1-107, j3-307, j2-207
	// join[007][b10]: 108
	// join[007][b20]: 208
	// join[007][fun]: j1-108, j3-308, j2-208
}
