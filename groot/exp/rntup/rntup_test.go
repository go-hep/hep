// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rntup

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/riofs"
)

func TestNTuple(t *testing.T) {
	for _, tc := range []struct {
		want rtests.ROOTer
	}{
		{
			want: &NTuple{1, 2, span{1, 2, 3}, span{4, 5, 6}, 7},
		},
	} {
		t.Run("", func(t *testing.T) {
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal: %+v", err)
			}

			rt := reflect.Indirect(reflect.ValueOf(tc.want)).Type()
			got := reflect.New(rt).Interface().(rtests.ROOTer)
			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)

			err = got.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal: %+v", err)
			}

			if got, want := got, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid r/w round-trip:\ngot= %#v\nwant=%#v", got, want)
			}
		})
	}
}

func TestReadNTuple(t *testing.T) {
	f, err := riofs.Open("../../testdata/ntpl001_staff.root")
	if err != nil {
		t.Fatalf("could not open file: +%v", err)
	}
	defer f.Close()

	obj, err := f.Get("Staff")
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	nt, ok := obj.(*NTuple)
	if !ok {
		t.Fatalf("%q not an NTuple: %T", "Staff", obj)
	}

	want := NTuple{
		rvers: 0x0,
		size:  0x30,
		header: span{
			seek:   854,
			nbytes: 537,
			length: 2495,
		},
		footer: span{
			seek:   72369,
			nbytes: 285,
			length: 804,
		},
		reserved: 0,
	}

	if got, want := *nt, want; got != want {
		t.Fatalf("error:\ngot= %#v\nwant=%#v", got, want)
	}

	if got, want := nt.String(), want.String(); got != want {
		t.Fatalf("error:\ngot= %v\nwant=%v", got, want)
	}
}
