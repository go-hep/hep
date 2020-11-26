// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/exp/rand"
)

func TestRand1D(t *testing.T) {
	edges := []float64{0, 1, 2, 3, 4}
	h1 := NewH1DFromEdges(edges)
	h1.FillN([]float64{
		0,
		1, 1, 1,
		2,
		3, 3, 3, 3, 3,
	}, nil)

	hr := NewRand1D(h1, rand.NewSource(1234))
	h2 := NewH1DFromEdges(edges)

	if got, want := hr.cdf, []float64{0, 0.1, 0.4, 0.5, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("invalid cdf:\ngot= %g\nwant=%g", got, want)
	}

	for _, test := range []struct{ v, want float64 }{
		{-10, 0},
		{+10, 1},
	} {
		if got, want := hr.CDF(test.v), test.want; got != want {
			t.Errorf("CDF(%g): got=%g, want=%g", test.v, got, want)
		}
	}

	const N = 1000
	for i := 0; i < N; i++ {
		h2.Fill(hr.Rand(), 1)
	}

	h1.Scale(1. / h1.Integral(h1.XMin(), h1.XMax()))
	h2.Scale(1. / h2.Integral(h2.XMin(), h2.XMax()))

	txt1, err := h1.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h1: %+v", err)
	}

	const wantTxt1 = `BEGIN YODA_HISTO1D_V2 /
Path: /
Title: ""
Type: Histo1D
---
# Mean: 2.000000e+00
# Area: 1.000000e+00
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.000000e+00	1.000000e-01	2.000000e+00	5.200000e+00	1.000000e+01
Underflow	Underflow	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
Overflow	Overflow	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
0.000000e+00	1.000000e+00	1.000000e-01	1.000000e-02	0.000000e+00	0.000000e+00	1.000000e+00
1.000000e+00	2.000000e+00	3.000000e-01	3.000000e-02	3.000000e-01	3.000000e-01	3.000000e+00
2.000000e+00	3.000000e+00	1.000000e-01	1.000000e-02	2.000000e-01	4.000000e-01	1.000000e+00
3.000000e+00	4.000000e+00	5.000000e-01	5.000000e-02	1.500000e+00	4.500000e+00	5.000000e+00
END YODA_HISTO1D_V2

`
	if got, want := string(txt1), wantTxt1; got != want {
		t.Errorf(
			"invalid h1 distribution:\n%s",
			cmp.Diff(want, got),
		)
	}

	const wantTxt2 = `BEGIN YODA_HISTO1D_V2 /
Path: /
Title: ""
Type: Histo1D
---
# Mean: 2.502788e+00
# Area: 1.000000e+00
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.000000e+00	1.000000e-03	2.502788e+00	7.531549e+00	1.000000e+03
Underflow	Underflow	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
Overflow	Overflow	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
0.000000e+00	1.000000e+00	9.200000e-02	9.200000e-05	4.339693e-02	2.823985e-02	9.200000e+01
1.000000e+00	2.000000e+00	3.090000e-01	3.090000e-04	4.611849e-01	7.110522e-01	3.090000e+02
2.000000e+00	3.000000e+00	8.700000e-02	8.700000e-05	2.154025e-01	5.413794e-01	8.700000e+01
3.000000e+00	4.000000e+00	5.120000e-01	5.120000e-04	1.782804e+00	6.250877e+00	5.120000e+02
END YODA_HISTO1D_V2

`
	txt2, err := h2.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal h2: %+v", err)
	}

	if got, want := string(txt2), wantTxt2; got != want {
		t.Errorf(
			"invalid h2 distribution:\n%s",
			cmp.Diff(want, got),
		)
	}
}
