// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootcnv_test

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/rootio"
)

func ExampleH1D() {
	f, err := rootio.Open("testdata/gauss-h1.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	obj, ok := f.Get("h1d")
	if !ok {
		log.Fatalf("no such histo %q\n", "h1d")
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
		obj, ok := f.Get(test.name)
		if !ok {
			t.Errorf("%s: no key %q", test.name, test.name)
			continue
		}
		rhisto := obj.(yodacnv.Marshaler)

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
