// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"testing"

	"go-hep.org/x/hep/groot/rbytes"
)

func TestRef(t *testing.T) {
	ref := Ref{pid: &gPID}
	if obj := ref.Object(); obj != nil {
		t.Fatalf("invalid referenced object")
	}
	if got, want := ref.UID(), uint32(0); got != want {
		t.Fatalf("invalid UID: got=%d, want=%d", got, want)
	}

	obj := NewObject()
	obj.ID = 42
	gPID.objs[obj.UID()] = obj

	ref.obj = *obj
	if ptr := ref.Object(); ptr != obj {
		t.Fatalf("invalid referenced object: got=%v, want=%v", ptr, obj)
	}
	if got, want := ref.UID(), obj.ID; got != want {
		t.Fatalf("invalid UID: got=%d, want=%d", got, want)
	}

	obj.SetBit(kIsReferenced)

	wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
	wbuf.WriteObject(obj)
	err := wbuf.Err()
	if err != nil {
		t.Fatalf("could not marshal ROOT: %+v", err)
	}

	rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, nil)
	var got Ref
	err = got.UnmarshalROOT(rbuf)
	if err != nil {
		t.Fatalf("could not unmarshal ROOT: %+v", err)
	}

	if got.obj.ID != obj.ID {
		t.Fatalf("invalid unmarshaled ref: got=%d, want=%d", got.obj.ID, obj.ID)
	}
}
