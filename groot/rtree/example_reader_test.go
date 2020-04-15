// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleReader() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatalf("could not retrieve ROOT tree: %+v", err)
	}
	t := o.(rtree.Tree)

	var (
		v1 int32
		v2 float32
		v3 string

		rvars = []rtree.ReadVar{
			{Name: "one", Value: &v1},
			{Name: "two", Value: &v2},
			{Name: "three", Value: &v3},
		}
	)

	r, err := rtree.NewReader(t, rvars)
	if err != nil {
		log.Fatalf("could not create tree reader: %+v", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf("evt[%d]: %v, %v, %v\n", ctx.Entry, v1, v2, v3)
		return nil
	})
	if err != nil {
		log.Fatalf("could not process tree: %+v", err)
	}

	// Output:
	// evt[0]: 1, 1.1, uno
	// evt[1]: 2, 2.2, dos
	// evt[2]: 3, 3.3, tres
	// evt[3]: 4, 4.4, quatro
}

func ExampleReader_withRange() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatalf("could not retrieve ROOT tree: %+v", err)
	}
	t := o.(rtree.Tree)

	var (
		v1 int32
		v2 float32
		v3 string

		rvars = []rtree.ReadVar{
			{Name: "one", Value: &v1},
			{Name: "two", Value: &v2},
			{Name: "three", Value: &v3},
		}
	)

	r, err := rtree.NewReader(t, rvars, rtree.WithRange(1, 3))
	if err != nil {
		log.Fatalf("could not create tree reader: %+v", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf("evt[%d]: %v, %v, %v\n", ctx.Entry, v1, v2, v3)
		return nil
	})
	if err != nil {
		log.Fatalf("could not process tree: %+v", err)
	}

	// Output:
	// evt[1]: 2, 2.2, dos
	// evt[2]: 3, 3.3, tres
}

func ExampleReader_withChain() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatalf("could not retrieve ROOT tree: %+v", err)
	}
	t := o.(rtree.Tree)

	t = rtree.Chain(t, t, t, t)

	var (
		v1 int32
		v2 float32
		v3 string

		rvars = []rtree.ReadVar{
			{Name: "one", Value: &v1},
			{Name: "two", Value: &v2},
			{Name: "three", Value: &v3},
		}
	)

	r, err := rtree.NewReader(t, rvars, rtree.WithRange(0, -1))
	if err != nil {
		log.Fatalf("could not create tree reader: %+v", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf("evt[%d]: %v, %v, %v\n", ctx.Entry, v1, v2, v3)
		return nil
	})
	if err != nil {
		log.Fatalf("could not process tree: %+v", err)
	}

	// Output:
	// evt[0]: 1, 1.1, uno
	// evt[1]: 2, 2.2, dos
	// evt[2]: 3, 3.3, tres
	// evt[3]: 4, 4.4, quatro
	// evt[4]: 1, 1.1, uno
	// evt[5]: 2, 2.2, dos
	// evt[6]: 3, 3.3, tres
	// evt[7]: 4, 4.4, quatro
	// evt[8]: 1, 1.1, uno
	// evt[9]: 2, 2.2, dos
	// evt[10]: 3, 3.3, tres
	// evt[11]: 4, 4.4, quatro
	// evt[12]: 1, 1.1, uno
	// evt[13]: 2, 2.2, dos
	// evt[14]: 3, 3.3, tres
	// evt[15]: 4, 4.4, quatro
}

func ExampleReader_withReadVarsFromStruct() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatalf("could not retrieve ROOT tree: %+v", err)
	}
	t := o.(rtree.Tree)

	var (
		data struct {
			V1 int32   `groot:"one"`
			V2 float32 `groot:"two"`
			V3 string  `groot:"three"`
		}
		rvars = rtree.ReadVarsFromStruct(&data)
	)

	r, err := rtree.NewReader(t, rvars)
	if err != nil {
		log.Fatalf("could not create tree reader: %+v", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		fmt.Printf("evt[%d]: %v, %v, %v\n", ctx.Entry, data.V1, data.V2, data.V3)
		return nil
	})
	if err != nil {
		log.Fatalf("could not process tree: %+v", err)
	}

	// Output:
	// evt[0]: 1, 1.1, uno
	// evt[1]: 2, 2.2, dos
	// evt[2]: 3, 3.3, tres
	// evt[3]: 4, 4.4, quatro
}
