// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsql_test

import (
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rsql"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hbook"
)

func ExampleScan() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	var data []float64

	err = rsql.Scan(tree, "SELECT two FROM tree", func(x float64) error {
		data = append(data, x)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tree[%q]: %v\n", "two", data)

	// Output:
	// tree["two"]: [1.1 2.2 3.3 4.4]
}

func ExampleScan_nVars() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	var (
		v1s []int32
		v2s []float64
		v3s []string
	)

	err = rsql.Scan(tree, "SELECT (one, two, three) FROM tree", func(x int32, y float64, z string) error {
		v1s = append(v1s, x)
		v2s = append(v2s, y)
		v3s = append(v3s, z)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tree[%q]: %v\n", "one", v1s)
	fmt.Printf("tree[%q]: %v\n", "two", v2s)
	fmt.Printf("tree[%q]: %q\n", "three", v3s)

	// Output:
	// tree["one"]: [1 2 3 4]
	// tree["two"]: [1.1 2.2 3.3 4.4]
	// tree["three"]: ["uno" "dos" "tres" "quatro"]
}

func ExampleScanH1D() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	h, err := rsql.ScanH1D(tree, "SELECT two FROM tree", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("entries: %v\n", h.Entries())
	fmt.Printf("x-axis: (min=%v, max=%v)\n", h.XMin(), h.XMax())
	fmt.Printf("x-mean: %v\n", h.XMean())
	fmt.Printf("x-std-dev: %v\nx-std-err: %v\n", h.XStdDev(), h.XStdErr())

	// Output:
	// entries: 4
	// x-axis: (min=1.1, max=4.400000000000001)
	// x-mean: 2.75
	// x-std-dev: 1.4200938936093859
	// x-std-err: 0.7100469468046929
}

func ExampleScanH1D_withH1D() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	h, err := rsql.ScanH1D(tree, "SELECT two FROM tree", hbook.NewH1D(100, 0, 10))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("entries: %v\n", h.Entries())
	fmt.Printf("x-axis: (min=%v, max=%v)\n", h.XMin(), h.XMax())
	fmt.Printf("x-mean: %v\n", h.XMean())
	fmt.Printf("x-std-dev: %v\nx-std-err: %v\n", h.XStdDev(), h.XStdErr())

	// Output:
	// entries: 4
	// x-axis: (min=0, max=10)
	// x-mean: 2.75
	// x-std-dev: 1.4200938936093859
	// x-std-err: 0.7100469468046929
}

func ExampleScanH2D() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	h, err := rsql.ScanH2D(tree, "SELECT (one, two) FROM tree", nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("entries: %v\n", h.Entries())
	fmt.Printf("x-axis: (min=%v, max=%v)\n", h.XMin(), h.XMax())
	fmt.Printf("x-mean: %v\n", h.XMean())
	fmt.Printf("x-std-dev: %v\nx-std-err: %v\n", h.XStdDev(), h.XStdErr())
	fmt.Printf("y-axis: (min=%v, max=%v)\n", h.YMin(), h.YMax())
	fmt.Printf("y-mean: %v\n", h.YMean())
	fmt.Printf("y-std-dev: %v\ny-std-err: %v\n", h.YStdDev(), h.YStdErr())

	// Output:
	// entries: 4
	// x-axis: (min=1, max=4.000000000000001)
	// x-mean: 2.5
	// x-std-dev: 1.2909944487358056
	// x-std-err: 0.6454972243679028
	// y-axis: (min=1.1, max=4.400000000000001)
	// y-mean: 2.75
	// y-std-dev: 1.4200938936093859
	// y-std-err: 0.7100469468046929
}

func ExampleScanH2D_withH2D() {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	h, err := rsql.ScanH2D(tree, "SELECT (one, two) FROM tree", hbook.NewH2D(100, 0, 10, 100, 0, 10))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("entries: %v\n", h.Entries())
	fmt.Printf("x-axis: (min=%v, max=%v)\n", h.XMin(), h.XMax())
	fmt.Printf("x-mean: %v\n", h.XMean())
	fmt.Printf("x-std-dev: %v\nx-std-err: %v\n", h.XStdDev(), h.XStdErr())
	fmt.Printf("y-axis: (min=%v, max=%v)\n", h.YMin(), h.YMax())
	fmt.Printf("y-mean: %v\n", h.YMean())
	fmt.Printf("y-std-dev: %v\ny-std-err: %v\n", h.YStdDev(), h.YStdErr())

	// Output:
	// entries: 4
	// x-axis: (min=0, max=10)
	// x-mean: 2.5
	// x-std-dev: 1.2909944487358056
	// x-std-err: 0.6454972243679028
	// y-axis: (min=0, max=10)
	// y-mean: 2.75
	// y-std-dev: 1.4200938936093859
	// y-std-err: 0.7100469468046929
}
