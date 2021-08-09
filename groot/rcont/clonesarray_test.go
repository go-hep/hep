// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcont_test

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
)

func TestTClonesArray(t *testing.T) {
	for _, fname := range []string{
		"../testdata/tclonesarray-no-streamerbypass.root",
		// "../testdata/tclonesarray-with-streamerbypass.root", // FIXME(sbinet): needs member-wise streaming.
	} {
		t.Run(fname, func(t *testing.T) {
			f, err := groot.Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get("clones")
			if err != nil {
				t.Fatal(err)
			}

			tca := o.(*rcont.ClonesArray)
			if got, want := tca.Len(), 3; got != want {
				t.Fatalf("invalid length: got=%d, want=%d", got, want)
			}
			if got, want := tca.Last(), 2; got != want {
				t.Fatalf("invalid last: got=%d, want=%d", got, want)
			}

			for i, want := range []root.Object{
				rbase.NewObjString("Elem-0"),
				rbase.NewObjString("elem-1"),
				rbase.NewObjString("Elem-20"),
			} {
				got := tca.At(i)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("invalid obj[%d]: got=%+v, want=%+v", i, got, want)
				}
			}
		})
	}
}

func TestTClonesArrayRW(t *testing.T) {
	dir, err := os.MkdirTemp("", "groot-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for i, tc := range []struct {
		name string
		want *rcont.ClonesArray
		cmp  func(a, b *rcont.ClonesArray) bool
	}{
		{
			name: "TClonesArray",
			want: func() *rcont.ClonesArray {
				o := rcont.NewClonesArray()
				o.SetElems([]root.Object{
					rbase.NewObjString("Elem-0"),
					rbase.NewObjString("elem-1"),
					rbase.NewObjString("Elem-20"),
				})
				return o
			}(),
			cmp: func(got, want *rcont.ClonesArray) bool {
				if g, w := got.Len(), want.Len(); g != w {
					return false
				}
				if g, w := got.Last(), want.Last(); g != w {
					return false
				}
				for i := 0; i < got.Len(); i++ {
					if g, w := got.At(i), want.At(i); !reflect.DeepEqual(g, w) {
						return false
					}
				}
				return true
			},
		},
	} {
		fname := filepath.Join(dir, fmt.Sprintf("tclonesarray-%d.root", i))
		t.Run(tc.name, func(t *testing.T) {
			const kname = "my-key"

			w, err := groot.Create(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()

			err = w.Put(kname, tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := len(w.Keys()), 1; got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			{
				wbuf := rbytes.NewWBuffer(nil, nil, 0, w)
				wbuf.SetErr(io.EOF)
				_, err := tc.want.MarshalROOT(wbuf)
				if err == nil {
					t.Fatalf("expected an error")
				}
				if err != io.EOF {
					t.Fatalf("got=%v, want=%v", err, io.EOF)
				}

				rbuf := rbytes.NewRBuffer(wbuf.Bytes(), nil, 0, w)
				class := tc.want.Class()
				obj := rtypes.Factory.Get(class)().Interface().(rbytes.Unmarshaler)
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

			err = w.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			r, err := groot.Open(fname)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			si := r.StreamerInfos()
			if len(si) == 0 {
				t.Fatalf("empty list of streamers")
			}

			if got, want := len(r.Keys()), 1; got != want {
				t.Fatalf("invalid number of keys. got=%d, want=%d", got, want)
			}

			rgot, err := r.Get(kname)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := rgot.(*rcont.ClonesArray), tc.want; !tc.cmp(got, want) {
				t.Fatalf("error reading back objstring.\ngot = %#v\nwant = %#v", got, want)
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("error closing file: %v", err)
			}

			if !rtests.HasROOT {
				t.Logf("skip test with ROOT/C++")
				return
			}

			const rootls = `#include <iostream>
#include "TFile.h"
#include "TClonesArray.h"

void rootls(const char *fname, const char *kname) {
	auto f = TFile::Open(fname);
	auto o = f->Get<TClonesArray>(kname);
	if (o == NULL) {
		std:cerr << "could not retrieve [" << kname << "]" << std::endl;
		o->ClassName();
	}
	std::cout << "retrieved TClonesArray: [" << kname << "]" << std::endl;
}
`
			out, err := rtests.RunCxxROOT("rootls", []byte(rootls), fname, kname)
			if err != nil {
				t.Fatalf("ROOT/C++ could not open file %q:\n%s", fname, string(out))
			}
		})
	}
}
