// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func ExampleDivideH1D(){

	h1 := NewH1D(5, 0, 5)
	h1.Fill(0, 1)
	h1.Fill(1, 2)
	h1.Fill(2, 3)
	h1.Fill(3, 3)
	h1.Fill(4, 4)

	h2 := NewH1D(5, 0, 5)
	h2.Fill(0, 11)
	h2.Fill(1, 22)
	h2.Fill(2, 0)
	h2.Fill(3, 33)
	h2.Fill(4, 44)

	s0, err := DivideH1D(h1, h2)
	if err != nil {
		panic(err)
	}
	s1, err := DivideH1D(h1, h2, DivIgnoreNaNs())
	if err != nil {
		panic(err)
	}
	s2, err := DivideH1D(h1, h2, DivReplaceNaNs(1.0))
	if err != nil {
		panic(err)
	}

	fmt.Println("Default:")
	for i, pt := range s0.Points() {
		fmt.Printf("Point %v: %.2f  + %.2f  - %.2f\n", i, pt.Y, pt.ErrY.Min, pt.ErrY.Min)
	}

	fmt.Println("\nDivIgnoreNaNs:")
	for i, pt := range s1.Points() {
		fmt.Printf("Point %v: %.2f  + %.2f  - %.2f\n", i, pt.Y, pt.ErrY.Min, pt.ErrY.Min)
	}

	fmt.Println("\nDivReplaceNaNs with v=1.0:")
	for i, pt := range s2.Points() {
		fmt.Printf("Point %v: %.2f  + %.2f  - %.2f\n", i, pt.Y, pt.ErrY.Min, pt.ErrY.Min)
	}
	
	// Output:
	// Default:
	// Point 0: 0.09  + 0.13  - 0.13
	// Point 1: 0.09  + 0.13  - 0.13
	// Point 2: NaN  + 0.00  - 0.00
	// Point 3: 0.09  + 0.13  - 0.13
	// Point 4: 0.09  + 0.13  - 0.13
	//
	// DivIgnoreNaNs:
	// Point 0: 0.09  + 0.13  - 0.13
	// Point 1: 0.09  + 0.13  - 0.13
	// Point 2: 0.09  + 0.13  - 0.13
	// Point 3: 0.09  + 0.13  - 0.13
	// 
	// DivReplaceNaNs with v=1.0:
	// Point 0: 0.09  + 0.13  - 0.13
	// Point 1: 0.09  + 0.13  - 0.13
	// Point 2: 1.00  + 0.00  - 0.00
	// Point 3: 0.09  + 0.13  - 0.13
	// Point 4: 0.09  + 0.13  - 0.13
}


func TestDivideH1D(t *testing.T) {
	h1 := NewH1D(5, 0, 5)
	h2 := NewH1D(5, 0, 5)
	for i := 0; i < 5; i++ {
		h1.Fill(float64(i), float64(i+1))
		h2.Fill(float64(i), float64(i+3))
	}
	s, err := DivideH1D(h1, h2)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := s.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	want := []byte(`BEGIN YODA_SCATTER2D /
Path=/
Title=
Type=Scatter2D
# xval	 xerr-	 xerr+	 yval	 yerr-	 yerr+
5.000000e-01	5.000000e-01	5.000000e-01	3.333333e-01	4.714045e-01	4.714045e-01
1.500000e+00	5.000000e-01	5.000000e-01	5.000000e-01	7.071068e-01	7.071068e-01
2.500000e+00	5.000000e-01	5.000000e-01	6.000000e-01	8.485281e-01	8.485281e-01
3.500000e+00	5.000000e-01	5.000000e-01	6.666667e-01	9.428090e-01	9.428090e-01
4.500000e+00	5.000000e-01	5.000000e-01	7.142857e-01	1.010153e+00	1.010153e+00
END YODA_SCATTER2D

`)

	if !reflect.DeepEqual(chk, want) {
		t.Fatalf("divide(num,den) differ:\ngot:\n%s\nwant:\n%s\n", string(chk), string(want))
	}
}

func TestAddH1DPanics(t *testing.T) {
	for _, tc := range []struct {
		h1, h2 *H1D
		panics error
	}{
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(5, 0, 10),
			panics: fmt.Errorf("hbook: h1 and h2 have different number of bins"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 1, 10),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 0, 11),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 1, 11),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
	} {
		t.Run("", func(t *testing.T) {
			if tc.panics != nil {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(error).Error(), tc.panics.Error(); got != want {
						t.Fatalf("invalid panic message.\ngot= %v\nwant=%v", got, want)
					}
				}()
			}
			_ = AddH1D(tc.h1, tc.h2)
		})
	}
}

func TestAddH1D(t *testing.T) {
	t.Skipf("missing some dist") // FIXME(sbinet)

	h1 := NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	h3 := AddH1D(h1, h2)

	got, err := h3.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal to yoda: %+v", err)
	}

	want := []byte(`BEGIN YODA_HISTO1D_V2 /
Path: /
Title: 
Type: Histo1D
---
# Mean: 2.526554e+00
# Area: 1.770000e+01
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	1.770000e+01	2.225000e+01	4.472000e+01	2.115580e+02	1.600000e+01
Underflow	Underflow	1.700000e+00	1.490000e+00	-8.500000e-01	4.250000e-01	2.000000e+00
Overflow	Overflow	1.500000e+00	1.250000e+00	9.800000e+00	6.424000e+01	2.000000e+00
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
0.000000e+00	1.000000e+00	4.700000e+00	5.690000e+00	1.540000e+00	8.780000e-01	4.000000e+00
1.000000e+00	2.000000e+00	1.800000e+00	1.640000e+00	2.400000e+00	3.240000e+00	2.000000e+00
2.000000e+00	3.000000e+00	2.700000e+00	4.490000e+00	5.740000e+00	1.220800e+01	2.000000e+00
3.000000e+00	4.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
4.000000e+00	5.000000e+00	2.300000e+00	2.690000e+00	9.790000e+00	4.167700e+01	2.000000e+00
5.000000e+00	6.000000e+00	3.000000e+00	5.000000e+00	1.630000e+01	8.889000e+01	2.000000e+00
END YODA_HISTO1D_V2

`)

	if !bytes.Equal(got, want) {
		t.Fatalf(
			"invalid yoda marshal response:\ngot:\n%s\nwant:\n%s\n",
			got, want,
		)
	}
}

func TestSubH1DPanics(t *testing.T) {
	for _, tc := range []struct {
		h1, h2 *H1D
		panics error
	}{
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(5, 0, 10),
			panics: fmt.Errorf("hbook: h1 and h2 have different number of bins"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 1, 10),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 0, 11),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
		{
			h1:     NewH1D(10, 0, 10),
			h2:     NewH1D(10, 1, 11),
			panics: fmt.Errorf("hbook: h1 and h2 have different range"),
		},
	} {
		t.Run("", func(t *testing.T) {
			if tc.panics != nil {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic")
					}
					if got, want := err.(error).Error(), tc.panics.Error(); got != want {
						t.Fatalf("invalid panic message.\ngot= %v\nwant=%v", got, want)
					}
				}()
			}
			_ = SubH1D(tc.h1, tc.h2)
		})
	}
}

func TestSubH1D(t *testing.T) {
	t.Skipf("missing some dist") // FIXME(sbinet)

	h1 := NewH1D(6, 0, 6)
	h1.Fill(-0.5, 1)
	h1.Fill(0, 1.5)
	h1.Fill(0.5, 1)
	h1.Fill(1.2, 1)
	h1.Fill(2.1, 2)
	h1.Fill(4.2, 1)
	h1.Fill(5.9, 1)
	h1.Fill(6, 0.5)

	h2 := NewH1D(6, 0, 6)
	h2.Fill(-0.5, 0.7)
	h2.Fill(0.2, 1)
	h2.Fill(0.7, 1.2)
	h2.Fill(1.5, 0.8)
	h2.Fill(2.2, 0.7)
	h2.Fill(4.3, 1.3)
	h2.Fill(5.2, 2)
	h2.Fill(6.8, 1)

	h3 := SubH1D(h1, h2)

	got, err := h3.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal to yoda: %+v", err)
	}

	want := []byte(`BEGIN YODA_HISTO1D_V2 /
Path: /
Title: 
Type: Histo1D
---
# Mean: -2.573333e+01
# Area: 3.000000e-01
# ID	 ID	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
Total   	Total   	3.000000e-01	2.225000e+01	-7.720000e+00	-4.913800e+01	1.600000e+01
Underflow	Underflow	3.000000e-01	1.490000e+00	-1.500000e-01	7.500000e-02	2.000000e+00
Overflow	Overflow	-5.000000e-01	1.250000e+00	-3.800000e+00	-2.824000e+01	2.000000e+00
# xlow	 xhigh	 sumw	 sumw2	 sumwx	 sumwx2	 numEntries
0.000000e+00	1.000000e+00	3.000000e-01	5.690000e+00	-5.400000e-01	-3.780000e-01	4.000000e+00
1.000000e+00	2.000000e+00	2.000000e-01	1.640000e+00	-2.220446e-16	-3.600000e-01	2.000000e+00
2.000000e+00	3.000000e+00	1.300000e+00	4.490000e+00	2.660000e+00	5.432000e+00	2.000000e+00
3.000000e+00	4.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00	0.000000e+00
4.000000e+00	5.000000e+00	-3.000000e-01	2.690000e+00	-1.390000e+00	-6.397000e+00	2.000000e+00
5.000000e+00	6.000000e+00	-1.000000e+00	5.000000e+00	-4.500000e+00	-1.927000e+01	2.000000e+00
END YODA_HISTO1D_V2

`)

	if !bytes.Equal(got, want) {
		t.Fatalf(
			"invalid yoda marshal response:\ngot:\n%s\nwant:\n%s\n",
			got, want,
		)
	}
}


