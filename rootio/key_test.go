// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestKeyNewKeyFrom(t *testing.T) {
	var (
		werr = fmt.Errorf("rootio: invalid")
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

			otyp := k.ObjectType()
			if otyp == nil {
				t.Fatalf("could not retrieve key's payload's type")
			}
			if otyp != nil {
				switch v := reflect.New(otyp).Elem().Interface(); v.(type) {
				case ObjString:
					// ok
				default:
					t.Fatalf("expected a rootio.ObjString (got %T)", v)
				}
			}
		})
	}
}

func TestKeyObjectType(t *testing.T) {
	for _, tc := range []struct {
		name string
		key  Key
		typ  reflect.Type
	}{
		{
			name: "TObjString",
			key:  Key{class: "TObjString"},
			typ:  reflect.TypeOf((*tobjstring)(nil)),
		},
		{
			name: "rootio.tobjstring",
			key:  Key{class: "*rootio.tobjstring"},
			typ:  reflect.TypeOf((*tobjstring)(nil)),
		},
		{
			name: "invalid",
			key:  Key{class: "invalid"},
			typ:  nil,
		},
		{
			name: "TH1D",
			key:  Key{class: "TH1D"},
			typ:  reflect.TypeOf((*H1D)(nil)),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			otyp := tc.key.ObjectType()
			if otyp != tc.typ {
				t.Fatalf("got: %+v, want=%+v", otyp, tc.typ)
			}
			if otyp == nil {
				return
			}

			otyp2 := tc.key.ObjectType()
			if otyp != otyp2 {
				t.Fatalf("got: %+v, want=%+v", otyp2, otyp)
			}
		})
	}
}
