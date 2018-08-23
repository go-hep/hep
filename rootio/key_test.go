// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
)

func TestKeyNewKeyFrom(t *testing.T) {
	var (
		werr = errors.Errorf("rootio: invalid")
	)
	for _, tc := range []struct {
		want *tobjstring
		wbuf *WBuffer
		err  error
	}{
		{
			want: NewObjString("hello"),
			wbuf: nil,
		},
		{
			want: NewObjString("hello"),
			wbuf: func() *WBuffer {
				wbuf := NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				return wbuf
			}(),
		},
		{
			want: NewObjString("hello"),
			wbuf: func() *WBuffer {
				wbuf := NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				wbuf.err = werr
				return wbuf
			}(),
			err: werr,
		},
		{
			want: NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: nil,
		},
		{
			want: NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: func() *WBuffer {
				wbuf := NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				return wbuf
			}(),
		},
		{
			want: NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: func() *WBuffer {
				wbuf := NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				wbuf.err = werr
				return wbuf
			}(),
			err: werr,
		},
	} {
		t.Run("", func(t *testing.T) {
			k, err := newKeyFrom(tc.want, tc.wbuf)
			switch {
			case err == nil && tc.err == nil:
				// ok
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error (%v)", tc.err)
			case err != nil && tc.err == nil:
				t.Fatalf("could not generate key from tobjstring: %v", err)
			case err != nil && tc.err != nil:
				if !reflect.DeepEqual(err, tc.err) {
					t.Fatalf("error: got=%#v, want=%#v", err, tc.err)
				}
				return
			}

			v, err := k.Object()
			if err != nil {
				t.Fatal(err)
			}

			if got := v.(*tobjstring); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v\n", got, tc.want)
			}
		})
	}
}
