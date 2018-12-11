// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
	_ "go-hep.org/x/hep/groot/ztypes"
)

type VisitT1 struct {
	Name   string
	F64    float64
	ArrF64 [2]float64
	// SliF64 []float64 // FIXME(sbinet)
}

type VisitT2 struct {
	VisitT1
	ArrF64 [2]float64
}

func TestVisit(t *testing.T) {
	rdict.Streamers.Add(rdict.StreamerOf(rdict.Streamers, reflect.TypeOf([2]float64{})))
	rdict.Streamers.Add(rdict.StreamerOf(rdict.Streamers, reflect.TypeOf([]float64{})))

	rdict.Streamers.Add(rdict.StreamerOf(rdict.Streamers, reflect.TypeOf((*VisitT1)(nil)).Elem()))
	rdict.Streamers.Add(rdict.StreamerOf(rdict.Streamers, reflect.TypeOf((*VisitT2)(nil)).Elem()))

	for _, tc := range []struct {
		si   rbytes.StreamerInfo
		want []string
	}{
		{
			si:   loadSI(t, "TObject"),
			want: []string{"fUniqueID", "fBits"},
		},
		{
			si:   loadSI(t, "TNamed"),
			want: []string{"TObject", "fUniqueID", "fBits", "fName", "fTitle"},
		},
		{
			si:   loadSI(t, "TObjString"),
			want: []string{"TObject", "fUniqueID", "fBits", "fString"},
		},
		{
			si:   loadSI(t, "VisitT1"),
			want: []string{"Name", "F64", "ArrF64"},
		},
		{
			si:   loadSI(t, "VisitT2"),
			want: []string{"VisitT1", "Name", "F64", "ArrF64", "ArrF64"},
		},
	} {
		t.Run(tc.si.Name(), func(t *testing.T) {
			var got []string
			err := rdict.Visit(nil, tc.si, func(depth int, se rbytes.StreamerElement) error {
				got = append(got, se.Name())
				return nil
			})
			if err != nil {
				t.Fatalf("could not visit %q: %v", tc.si.Name(), err)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid element list.\ngot= %v\nwant=%v\n", got, tc.want)
			}
		})
	}
}

func loadSI(t *testing.T, name string) rbytes.StreamerInfo {
	t.Helper()

	si, err := rdict.Streamers.StreamerInfo(name, -1)
	if err != nil {
		t.Fatal(err)
	}
	return si
}
