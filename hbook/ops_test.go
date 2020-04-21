// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"fmt"
	"reflect"
	"testing"
)

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
