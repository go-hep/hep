// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbytes"
)

func TestObjectFrom(t *testing.T) {
	sictx := StreamerInfos
	loadSI := func(name string) rbytes.StreamerInfo {
		t.Helper()
		si, err := sictx.StreamerInfo(name, -1)
		if err != nil {
			t.Fatalf("could not load streamer %q: %+v", name, err)
		}
		return si
	}

	for _, tc := range []struct {
		name string
		si   rbytes.StreamerInfo
	}{
		{
			name: "TObject",
			si:   loadSI("TObject"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			obj := ObjectFrom(tc.si, sictx)

			if got, want := obj.Class(), tc.si.Name(); got != want {
				t.Fatalf("invalid class name: got=%v, want=%v", got, want)
			}

			if got, want := obj.RVersion(), int16(tc.si.ClassVersion()); got != want {
				t.Fatalf("invalid class version: got=%v, want=%v", got, want)
			}

			wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
			_, err := obj.MarshalROOT(wbuf)
			if err != nil {
				t.Fatalf("could not write object: %+v", err)
			}

			rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)

			got := ObjectFrom(tc.si, sictx)
			err = got.UnmarshalROOT(rbuf)
			if err != nil {
				t.Fatalf("could not read object: %+v", err)
			}

			if !reflect.DeepEqual(got, obj) {
				t.Fatalf("round-trip failed:\ngot= %#v\nwant=%#v", got, obj)
			}
		})
	}
}
