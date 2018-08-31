// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootcnv_test

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/rootio"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleH1D() {
	f, err := rootio.Open("testdata/gauss-h1.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("h1d")
	if err != nil {
		log.Fatal(err)
	}

	root := obj.(*rootio.H1D)
	h, err := rootcnv.H1D(root)
	if err != nil {
		log.Fatalf("error converting TH1D: %v\n", err)
	}

	fmt.Printf("name:    %q\n", root.Name())
	fmt.Printf("mean:    %v\n", h.XMean())
	fmt.Printf("std-dev: %v\n", h.XStdDev())
	fmt.Printf("std-err: %v\n", h.XStdErr())

	// Output:
	// name:    "h1d"
	// mean:    0.028120161729965475
	// std-dev: 2.5450388581847907
	// std-err: 0.025447022905060374
}

func ExampleH2D() {
	f, err := rootio.Open("testdata/gauss-h2.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("h2d")
	if err != nil {
		log.Fatal(err)
	}

	root := obj.(*rootio.H2D)
	h, err := rootcnv.H2D(root)
	if err != nil {
		log.Fatalf("error converting TH2D: %v\n", err)
	}

	fmt.Printf("name:      %q\n", root.Name())
	fmt.Printf("x-mean:    %v\n", h.XMean())
	fmt.Printf("x-std-dev: %v\n", h.XStdDev())
	fmt.Printf("x-std-err: %v\n", h.XStdErr())
	fmt.Printf("y-mean:    %v\n", h.YMean())
	fmt.Printf("y-std-dev: %v\n", h.YStdDev())
	fmt.Printf("y-std-err: %v\n", h.YStdErr())

	// Output:
	// name:      "h2d"
	// x-mean:    -0.005792200073827981
	// x-std-dev: 2.2708057036302542
	// x-std-err: 0.06540325697671126
	// y-mean:    0.894201827242525
	// y-std-dev: 1.8307942514179008
	// y-std-err: 0.05273014186354511
}

func ExampleS2D() {
	f, err := rootio.Open("../../rootio/testdata/graphs.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tgae")
	if err != nil {
		log.Fatal(err)
	}

	root := obj.(rootio.GraphErrors)
	g, err := rootcnv.S2D(root)
	if err != nil {
		log.Fatalf("error converting TGraphErrors: %v\n", err)
	}

	fmt.Printf("name:  %q\n", g.Annotation()["name"])
	fmt.Printf("title: %q\n", g.Annotation()["title"])
	fmt.Printf("#pts:  %v\n", g.Len())
	for i, pt := range g.Points() {
		x := pt.X
		y := pt.Y
		xlo := pt.ErrX.Min
		xhi := pt.ErrX.Max
		ylo := pt.ErrY.Min
		yhi := pt.ErrY.Max
		fmt.Printf("(x,y)[%d] = (%+e +/- [%+e, %+e], %+e +/- [%+e, %+e])\n", i, x, xlo, xhi, y, ylo, yhi)
	}

	// Output:
	// name:  "tgae"
	// title: "graph with asymmetric errors"
	// #pts:  4
	// (x,y)[0] = (+1.000000e+00 +/- [+1.000000e-01, +2.000000e-01], +2.000000e+00 +/- [+3.000000e-01, +4.000000e-01])
	// (x,y)[1] = (+2.000000e+00 +/- [+2.000000e-01, +4.000000e-01], +4.000000e+00 +/- [+6.000000e-01, +8.000000e-01])
	// (x,y)[2] = (+3.000000e+00 +/- [+3.000000e-01, +6.000000e-01], +6.000000e+00 +/- [+9.000000e-01, +1.200000e+00])
	// (x,y)[3] = (+4.000000e+00 +/- [+4.000000e-01, +8.000000e-01], +8.000000e+00 +/- [+1.200000e+00, +1.600000e+00])
}

func TestH1D(t *testing.T) {
	f, err := rootio.Open("testdata/gauss-h1.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, test := range []struct {
		name string
		want []byte
	}{
		{
			name: "h1d",
			want: []byte(`BEGIN YODA_HISTO1D /h1d
Path=/h1d
Title=h1d
Type=Histo1D
# Mean: 2.812016e-02
# Area: 1.100600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.100600e+04	1.211000e+04	3.094905e+02	7.128989e+04	10004
Underflow	Underflow	2.000000e+00	2.000000e+00	0.000000e+00	0.000000e+00	2
Overflow	Overflow	4.000000e+00	8.000000e+00	0.000000e+00	0.000000e+00	2
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
-4.000000e+00	-3.200000e+00	6.600000e+00	7.260000e+00	0.000000e+00	0.000000e+00	6
-3.200000e+00	-2.400000e+00	7.260000e+01	7.986000e+01	0.000000e+00	0.000000e+00	66
-2.400000e+00	-1.600000e+00	5.434000e+02	5.977400e+02	0.000000e+00	0.000000e+00	494
-1.600000e+00	-8.000000e-01	1.708300e+03	1.879130e+03	0.000000e+00	0.000000e+00	1553
-8.000000e-01	2.220446e-16	3.130600e+03	3.443660e+03	0.000000e+00	0.000000e+00	2846
0.000000e+00	8.000000e-01	3.136100e+03	3.449710e+03	0.000000e+00	0.000000e+00	2851
8.000000e-01	1.600000e+00	1.753400e+03	1.928740e+03	0.000000e+00	0.000000e+00	1594
1.600000e+00	2.400000e+00	5.401000e+02	5.941100e+02	0.000000e+00	0.000000e+00	491
2.400000e+00	3.200000e+00	1.012000e+02	1.113200e+02	0.000000e+00	0.000000e+00	92
3.200000e+00	4.000000e+00	7.700000e+00	8.470000e+00	0.000000e+00	0.000000e+00	7
END YODA_HISTO1D

`),
		},
		{
			name: "h1f",
			want: []byte(`BEGIN YODA_HISTO1D /h1f
Path=/h1f
Title=h1f
Type=Histo1D
# Mean: 2.812016e-02
# Area: 1.100600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.100600e+04	1.211000e+04	3.094905e+02	7.128989e+04	10004
Underflow	Underflow	2.000000e+00	2.000000e+00	0.000000e+00	0.000000e+00	2
Overflow	Overflow	4.000000e+00	8.000000e+00	0.000000e+00	0.000000e+00	2
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
-4.000000e+00	-3.200000e+00	6.600000e+00	7.260000e+00	0.000000e+00	0.000000e+00	6
-3.200000e+00	-2.400000e+00	7.259995e+01	7.986000e+01	0.000000e+00	0.000000e+00	66
-2.400000e+00	-1.600000e+00	5.434013e+02	5.977400e+02	0.000000e+00	0.000000e+00	494
-1.600000e+00	-8.000000e-01	1.708276e+03	1.879130e+03	0.000000e+00	0.000000e+00	1553
-8.000000e-01	2.220446e-16	3.130664e+03	3.443660e+03	0.000000e+00	0.000000e+00	2846
0.000000e+00	8.000000e-01	3.136165e+03	3.449710e+03	0.000000e+00	0.000000e+00	2851
8.000000e-01	1.600000e+00	1.753375e+03	1.928740e+03	0.000000e+00	0.000000e+00	1594
1.600000e+00	2.400000e+00	5.401014e+02	5.941100e+02	0.000000e+00	0.000000e+00	491
2.400000e+00	3.200000e+00	1.011999e+02	1.113200e+02	0.000000e+00	0.000000e+00	92
3.200000e+00	4.000000e+00	7.700000e+00	8.470000e+00	0.000000e+00	0.000000e+00	7
END YODA_HISTO1D

`),
		},
		{
			name: "h1d-var",
			want: []byte(`BEGIN YODA_HISTO1D /h1d-var
Path=/h1d-var
Title=h1d-var
Type=Histo1D
# Mean: 2.812016e-02
# Area: 1.100600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.100600e+04	1.211000e+04	3.094905e+02	7.128989e+04	10004
Underflow	Underflow	2.000000e+00	2.000000e+00	0.000000e+00	0.000000e+00	2
Overflow	Overflow	4.000000e+00	8.000000e+00	0.000000e+00	0.000000e+00	2
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
-4.000000e+00	-3.200000e+00	6.600000e+00	7.260000e+00	0.000000e+00	0.000000e+00	6
-3.200000e+00	-2.400000e+00	7.259995e+01	7.986000e+01	0.000000e+00	0.000000e+00	66
-2.400000e+00	-1.600000e+00	5.434013e+02	5.977400e+02	0.000000e+00	0.000000e+00	494
-1.600000e+00	-8.000000e-01	1.708276e+03	1.879130e+03	0.000000e+00	0.000000e+00	1553
-8.000000e-01	0.000000e+00	3.130664e+03	3.443660e+03	0.000000e+00	0.000000e+00	2846
0.000000e+00	8.000000e-01	3.136165e+03	3.449710e+03	0.000000e+00	0.000000e+00	2851
8.000000e-01	1.600000e+00	1.753375e+03	1.928740e+03	0.000000e+00	0.000000e+00	1594
1.600000e+00	2.400000e+00	5.401014e+02	5.941100e+02	0.000000e+00	0.000000e+00	491
2.400000e+00	3.200000e+00	1.011999e+02	1.113200e+02	0.000000e+00	0.000000e+00	92
3.200000e+00	4.000000e+00	7.700000e+00	8.470000e+00	0.000000e+00	0.000000e+00	7
END YODA_HISTO1D

`),
		},
		{
			name: "h1f-var",
			want: []byte(`BEGIN YODA_HISTO1D /h1f-var
Path=/h1f-var
Title=h1f-var
Type=Histo1D
# Mean: 2.812016e-02
# Area: 1.100600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.100600e+04	1.211000e+04	3.094905e+02	7.128989e+04	10004
Underflow	Underflow	2.000000e+00	2.000000e+00	0.000000e+00	0.000000e+00	2
Overflow	Overflow	4.000000e+00	8.000000e+00	0.000000e+00	0.000000e+00	2
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
-4.000000e+00	-3.200000e+00	6.600000e+00	7.260000e+00	0.000000e+00	0.000000e+00	6
-3.200000e+00	-2.400000e+00	7.259995e+01	7.986000e+01	0.000000e+00	0.000000e+00	66
-2.400000e+00	-1.600000e+00	5.434013e+02	5.977400e+02	0.000000e+00	0.000000e+00	494
-1.600000e+00	-8.000000e-01	1.708276e+03	1.879130e+03	0.000000e+00	0.000000e+00	1553
-8.000000e-01	0.000000e+00	3.130664e+03	3.443660e+03	0.000000e+00	0.000000e+00	2846
0.000000e+00	8.000000e-01	3.136165e+03	3.449710e+03	0.000000e+00	0.000000e+00	2851
8.000000e-01	1.600000e+00	1.753375e+03	1.928740e+03	0.000000e+00	0.000000e+00	1594
1.600000e+00	2.400000e+00	5.401014e+02	5.941100e+02	0.000000e+00	0.000000e+00	491
2.400000e+00	3.200000e+00	1.011999e+02	1.113200e+02	0.000000e+00	0.000000e+00	92
3.200000e+00	4.000000e+00	7.700000e+00	8.470000e+00	0.000000e+00	0.000000e+00	7
END YODA_HISTO1D

`),
		},
	} {
		obj, err := f.Get(test.name)
		if err != nil {
			t.Errorf("%s: error: %v", test.name, err)
			continue
		}
		rhisto := obj.(rootio.H1)

		h, err := rootcnv.H1D(rhisto)
		if err != nil {
			t.Errorf("%s: convertion error: %v", test.name, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = yodacnv.Write(buf, h)
		if err != nil {
			t.Errorf("%s: YODA error: %v", test.name, err)
			continue
		}

		if !reflect.DeepEqual(buf.Bytes(), test.want) {
			t.Errorf("error converting %s:\ngot:\n%s\nwant:\n%s\n",
				test.name,
				string(buf.Bytes()),
				string(test.want),
			)
			continue
		}
	}
}

func TestH2D(t *testing.T) {
	f, err := rootio.Open("testdata/gauss-h2.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	for _, test := range []struct {
		name string
		want []byte
	}{
		{
			name: "h2f",
			want: []byte(`BEGIN YODA_HISTO2D /h2f
Path=/h2f
Title=h2f
Type=Histo2D
# Mean: (-5.792200e-03, 8.942018e-01)
# Volume: 1.083600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 sumwy	 sumwy2	 sumwxy	 numEntries
Total   	Total   	1.083600e+04	9.740400e+04	-6.276428e+01	5.583048e+04	9.689571e+03	4.495449e+04	-1.878975e+02	10008
# 2D outflow persistency not currently supported until API is stable
# xlow	 xhigh	 ylow	 yhigh	 sumw	 sumw2	 sumwx	 sumwx2	 sumwy	 sumwy2	 sumwxy	 numEntries
0.000000e+00	1.000000e+00	0.000000e+00	1.000000e+00	5.010000e+02	5.010000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	501
0.000000e+00	1.000000e+00	1.000000e+00	2.000000e+00	4.880000e+02	4.880000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	488
0.000000e+00	1.000000e+00	2.000000e+00	3.000000e+00	3.140000e+02	3.140000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	314
1.000000e+00	2.000000e+00	0.000000e+00	1.000000e+00	3.850000e+02	3.850000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	385
1.000000e+00	2.000000e+00	1.000000e+00	2.000000e+00	3.790000e+02	3.790000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	379
1.000000e+00	2.000000e+00	2.000000e+00	3.000000e+00	2.210000e+02	2.210000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	221
2.000000e+00	3.000000e+00	0.000000e+00	1.000000e+00	2.280000e+02	2.280000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	228
2.000000e+00	3.000000e+00	1.000000e+00	2.000000e+00	2.320000e+02	2.320000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	232
2.000000e+00	3.000000e+00	2.000000e+00	3.000000e+00	1.640000e+02	1.640000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	164
END YODA_HISTO2D

`),
		},
		{
			name: "h2d",
			want: []byte(`BEGIN YODA_HISTO2D /h2d
Path=/h2d
Title=h2d
Type=Histo2D
# Mean: (-5.792200e-03, 8.942018e-01)
# Volume: 1.083600e+04
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 sumwy	 sumwy2	 sumwxy	 numEntries
Total   	Total   	1.083600e+04	9.740400e+04	-6.276428e+01	5.583048e+04	9.689571e+03	4.495449e+04	-1.878975e+02	10008
# 2D outflow persistency not currently supported until API is stable
# xlow	 xhigh	 ylow	 yhigh	 sumw	 sumw2	 sumwx	 sumwx2	 sumwy	 sumwy2	 sumwxy	 numEntries
0.000000e+00	1.000000e+00	0.000000e+00	1.000000e+00	5.010000e+02	5.010000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	501
0.000000e+00	1.000000e+00	1.000000e+00	2.000000e+00	4.880000e+02	4.880000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	488
0.000000e+00	1.000000e+00	2.000000e+00	3.000000e+00	3.140000e+02	3.140000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	314
1.000000e+00	2.000000e+00	0.000000e+00	1.000000e+00	3.850000e+02	3.850000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	385
1.000000e+00	2.000000e+00	1.000000e+00	2.000000e+00	3.790000e+02	3.790000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	379
1.000000e+00	2.000000e+00	2.000000e+00	3.000000e+00	2.210000e+02	2.210000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	221
2.000000e+00	3.000000e+00	0.000000e+00	1.000000e+00	2.280000e+02	2.280000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	228
2.000000e+00	3.000000e+00	1.000000e+00	2.000000e+00	2.320000e+02	2.320000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	232
2.000000e+00	3.000000e+00	2.000000e+00	3.000000e+00	1.640000e+02	1.640000e+02	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	164
END YODA_HISTO2D

`),
		},
	} {
		obj, err := f.Get(test.name)
		if err != nil {
			t.Errorf("%s: error: %v", test.name, err)
			continue
		}
		rhisto := obj.(rootio.H2)

		h, err := rootcnv.H2D(rhisto)
		if err != nil {
			t.Errorf("%s: convertion error: %v", test.name, err)
			continue
		}

		buf := new(bytes.Buffer)
		err = yodacnv.Write(buf, h)
		if err != nil {
			t.Errorf("%s: YODA error: %v", test.name, err)
			continue
		}

		if !reflect.DeepEqual(buf.Bytes(), test.want) {
			t.Errorf("error converting %s:\ngot:\n%s\nwant:\n%s\n",
				test.name,
				string(buf.Bytes()),
				string(test.want),
			)
			continue
		}
	}
}

func TestFromH1D(t *testing.T) {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	h.Fill(-10, 1) // fill underflow
	h.Fill(-20, 2)
	h.Fill(+10, 1) // fill overflow
	h.Fill(+10, 2)
	h.Annotation()["name"] = "my-name"
	h.Annotation()["title"] = "my-title"

	for _, tc := range []struct {
		name   string
		h1     rootio.H1
		sumw   float64
		sumw2  float64
		sumwx  float64
		sumwx2 float64
	}{
		{
			name:   "TH1D",
			h1:     rootcnv.FromH1D(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.h1.SumW(), h.SumW(); got != want {
				t.Fatalf("sumw: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumW2(), h.SumW2(); got != want {
				t.Fatalf("sumw2: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumWX(), h.SumWX(); got != want {
				t.Fatalf("sumwx: got=%v, want=%v", got, want)
			}
			if got, want := tc.h1.SumWX2(), h.SumWX2(); got != want {
				t.Fatalf("sumwx2: got=%v, want=%v", got, want)
			}

			rraw, err := tc.h1.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			hh, err := rootcnv.H1D(tc.h1)
			if err != nil {
				t.Fatal(err)
			}

			hraw, err := hh.MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			var hr = rootio.Factory.Get(tc.name)().Interface().(rootio.H1)
			if err := hr.(yodacnv.Unmarshaler).UnmarshalYODA(hraw); err != nil {
				t.Fatal(err)
			}

			rgot, err := hr.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(rgot, rraw) {
				t.Fatalf("round trip error:\nraw:\n%s\ngot:\n%s\n", rraw, rgot)
			}
		})
	}
}

func TestFromH2D(t *testing.T) {
	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH2D(5, -4, +4, 6, -4, +4)
	for i := 0; i < npoints; i++ {
		x := dist.Rand()
		y := dist.Rand()
		h.Fill(x, y, 1)
	}
	h.Fill(+0, +5, 1) // N
	h.Fill(-5, +5, 2) // N-W
	h.Fill(-5, +0, 3) // W
	h.Fill(-5, -5, 4) // S-W
	h.Fill(+0, -5, 5) // S
	h.Fill(+5, -5, 6) // S-E
	h.Fill(+5, +0, 7) // E
	h.Fill(+5, +5, 8) // N-E

	h.Annotation()["name"] = "my-name"
	h.Annotation()["title"] = "my-title"

	for _, tc := range []struct {
		name   string
		h2     rootio.H2
		sumw   float64
		sumw2  float64
		sumwx  float64
		sumwx2 float64
	}{
		{
			name:   "TH2D",
			h2:     rootcnv.FromH2D(h),
			sumw:   h.SumW(),
			sumw2:  h.SumW2(),
			sumwx:  h.SumWX(),
			sumwx2: h.SumWX2(),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got, want := tc.h2.SumW(), h.SumW(); got != want {
				t.Fatalf("sumw: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumW2(), h.SumW2(); got != want {
				t.Fatalf("sumw2: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumWX(), h.SumWX(); got != want {
				t.Fatalf("sumwx: got=%v, want=%v", got, want)
			}
			if got, want := tc.h2.SumWX2(), h.SumWX2(); got != want {
				t.Fatalf("sumwx2: got=%v, want=%v", got, want)
			}

			rraw, err := tc.h2.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			hh, err := rootcnv.H2D(tc.h2)
			if err != nil {
				t.Fatal(err)
			}

			hraw, err := hh.MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			var hr = rootio.Factory.Get(tc.name)().Interface().(rootio.H2)
			if err := hr.(yodacnv.Unmarshaler).UnmarshalYODA(hraw); err != nil {
				t.Fatal(err)
			}

			rgot, err := hr.(yodacnv.Marshaler).MarshalYODA()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(rgot, rraw) {
				t.Fatalf("round trip error:\nraw:\n%s\ngot:\n%s\n", rraw, rgot)
			}
		})
	}
}

func TestFromS2D(t *testing.T) {
	hg := hbook.NewS2D(
		hbook.Point2D{X: 1, Y: 1, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
		hbook.Point2D{X: 2, Y: 1.5, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
		hbook.Point2D{X: -1, Y: +2, ErrX: hbook.Range{Min: 1, Max: 2}, ErrY: hbook.Range{Min: 3, Max: 4}},
	)

	rg := rootcnv.FromS2D(hg)

	hr, err := rootcnv.S2D(rg)
	if err != nil {
		t.Fatal(err)
	}

	want, err := hg.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	got, err := hr.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("got:\n%s\nwant:\n%s\n", got, want)
	}
}
