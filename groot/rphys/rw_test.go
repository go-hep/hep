// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rphys

import (
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestWRBuffer(t *testing.T) {
	for _, tc := range []struct {
		name string
		want rtests.ROOTer
		cmp  func(a, b rtests.ROOTer) bool
	}{
		{
			name: "TLorentzVector",
			want: &LorentzVector{
				obj: rbase.Object{ID: 0x0, Bits: 0x3000000},
				p: Vector3{
					obj: rbase.Object{ID: 0x0, Bits: 0x3000000},
					x:   1, y: 2, z: 3,
				},
				e: 4,
			},
		},
		{
			name: "TVector2",
			want: &Vector2{
				obj: rbase.Object{ID: 0x0, Bits: 0x3000000},
				x:   1, y: 2,
			},
		},
		{
			name: "TVector3",
			want: &Vector3{
				obj: rbase.Object{ID: 0x0, Bits: 0x3000000},
				x:   1, y: 2, z: 3,
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
			}
			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := tc.want.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not marshal ROOT: %v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
			class := tc.want.Class()
			obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
			{
				rbuf.SetErr(io.EOF)
				err = obj.UnmarshalROOT(rbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}
				rbuf.SetErr(nil)
			}
			err = obj.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not unmarshal ROOT: %v", err)
			}

			switch tc.cmp {
			case nil:
				if !reflect.DeepEqual(obj, tc.want) {
					t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
				}
			default:
				obj := obj.(rtests.ROOTer)
				if !tc.cmp(obj, tc.want) {
					t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
				}
			}
		})
	}
}
