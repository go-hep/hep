// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package endsess_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/endsess"
)

func TestRequest(t *testing.T) {
	for _, want := range []endsess.Request{
		{},
		{SessionID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got endsess.Request
			)

			if want.ReqID() != endsess.RequestID {
				t.Fatalf("invalid request ID: got=%d want=%d", want.ReqID(), endsess.RequestID)
			}

			if want.ShouldSign() {
				t.Fatalf("invalid")
			}

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal request: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal request: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed:\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

var (
	_ xrdproto.Request = (*endsess.Request)(nil)
)
