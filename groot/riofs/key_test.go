// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"reflect"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rhist"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rvers"
)

func TestKeyNewKeyFrom(t *testing.T) {
	var (
		werr = errors.Errorf("riofs: invalid")
	)
	for _, tc := range []struct {
		want *rbase.ObjString
		wbuf *rbytes.WBuffer
		err  error
	}{
		{
			want: rbase.NewObjString("hello"),
			wbuf: nil,
		},
		{
			want: rbase.NewObjString("hello"),
			wbuf: func() *rbytes.WBuffer {
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				return wbuf
			}(),
		},
		{
			want: rbase.NewObjString("hello"),
			wbuf: func() *rbytes.WBuffer {
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				wbuf.SetErr(werr)
				return wbuf
			}(),
			err: werr,
		},
		{
			want: rbase.NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: nil,
		},
		{
			want: rbase.NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: func() *rbytes.WBuffer {
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				return wbuf
			}(),
		},
		{
			want: rbase.NewObjString(strings.Repeat("+", 512) + "hello"),
			wbuf: func() *rbytes.WBuffer {
				wbuf := rbytes.NewWBuffer(nil, nil, 0, nil)
				wbuf.WriteString(strings.Repeat("=+=", 80))
				wbuf.SetErr(werr)
				return wbuf
			}(),
			err: werr,
		},
	} {
		t.Run("", func(t *testing.T) {
			var parent Directory
			k, err := newTestKeyFrom(parent, tc.want, tc.wbuf)
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

			if got := v.(*rbase.ObjString); !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v\n", got, tc.want)
			}

			otyp := k.ObjectType()
			if otyp == nil {
				t.Fatalf("could not retrieve key's payload's type")
			}
			if otyp != nil {
				switch v := reflect.New(otyp).Elem().Interface(); v.(type) {
				case root.ObjString:
					// ok
				default:
					t.Fatalf("expected a root.ObjString (got %T)", v)
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
			typ:  reflect.TypeOf((*rbase.ObjString)(nil)),
		},
		{
			name: "invalid",
			key:  Key{class: "invalid"},
			typ:  nil,
		},
		{
			name: "TH1D",
			key:  Key{class: "TH1D"},
			typ:  reflect.TypeOf((*rhist.H1D)(nil)),
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

func newTestKeyFrom(dir Directory, obj root.Object, wbuf *rbytes.WBuffer) (Key, error) {
	if wbuf == nil {
		wbuf = rbytes.NewWBuffer(nil, nil, 0, nil)
	}
	beg := int(wbuf.Pos())
	n, err := obj.(rbytes.Marshaler).MarshalROOT(wbuf)
	if err != nil {
		return Key{}, err
	}
	end := beg + n
	data := wbuf.Bytes()[beg:end]

	name := ""
	title := ""
	if obj, ok := obj.(root.Named); ok {
		name = obj.Name()
		title = obj.Title()
	}

	k := Key{
		rvers:    rvers.Key,
		objlen:   int32(n),
		datetime: nowUTC(),
		class:    obj.Class(),
		name:     name,
		title:    title,
		buf:      data,
		obj:      obj,
		otyp:     reflect.TypeOf(obj),
	}
	return k, nil
}
