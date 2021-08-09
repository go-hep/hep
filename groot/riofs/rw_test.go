// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import (
	"io"
	"os"
	"path/filepath"
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
	}{
		{
			name: "TFree",
			want: &freeSegment{
				first: 21,
				last:  24,
			},
		},
		{
			name: "TFree",
			want: &freeSegment{
				first: 21,
				last:  kStartBigFile + 24,
			},
		},
		{
			name: "TKey",
			want: &Key{
				nbytes:   1024,
				rvers:    4, // small file
				objlen:   10,
				datetime: datime2time(1576331001),
				keylen:   12,
				cycle:    2,
				seekkey:  1024,
				seekpdir: 2048,
				class:    "MyClass",
				name:     "my-key",
				title:    "my key title",
			},
		},
		{
			name: "TKey",
			want: &Key{
				nbytes:   1024,
				rvers:    1004, // big file
				objlen:   10,
				datetime: datime2time(1576331001),
				keylen:   12,
				cycle:    2,
				seekkey:  1024,
				seekpdir: 2048,
				class:    "MyClass",
				name:     "my-key",
				title:    "my key title",
			},
		},
		{
			name: "TDirectory",
			want: &tdirectory{
				rvers: 4, // small file
				named: *rbase.NewNamed("my-name", "my-title"),
				uuid: rbase.UUID{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					10, 11, 12, 13, 14, 15,
				},
			},
		},
		{
			name: "TDirectory",
			want: &tdirectory{
				rvers: 1004, // big file
				named: *rbase.NewNamed("my-name", "my-title"),
				uuid: rbase.UUID{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					10, 11, 12, 13, 14, 15,
				},
			},
		},
		{
			name: "TDirectoryFile",
			want: &tdirectoryFile{
				dir: tdirectory{
					rvers: 4, // small file
					named: *rbase.NewNamed("", ""),
					uuid: rbase.UUID{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
						10, 11, 12, 13, 14, 15,
					},
				},
				ctime:      datime2time(1576331001),
				mtime:      datime2time(1576331010),
				nbyteskeys: 1,
				nbytesname: 2,
				seekdir:    3,
				seekparent: 4,
				seekkeys:   5,
			},
		},
		{
			name: "TDirectoryFile",
			want: &tdirectoryFile{
				dir: tdirectory{
					rvers: 1004, // big file
					named: *rbase.NewNamed("", ""),
					uuid: rbase.UUID{
						0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
						10, 11, 12, 13, 14, 15,
					},
				},
				ctime:      datime2time(1576331001),
				mtime:      datime2time(1576331010),
				nbyteskeys: 1,
				nbytesname: 2,
				seekdir:    3,
				seekparent: 4,
				seekkeys:   5,
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

			if !reflect.DeepEqual(obj, tc.want) {
				t.Fatalf("error\ngot= %+v\nwant=%+v\n", obj, tc.want)
			}
		})
	}
}

func TestWriteBigFile(t *testing.T) {
	tmp, err := os.MkdirTemp("", "groot-riofs-")
	if err != nil {
		t.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	fname := filepath.Join(tmp, "big-file.root")

	kvals := []struct {
		k string
		v string
	}{
		{k: "key1", v: "obj1"},
		{k: "key2", v: "obj2"},
	}

	func() {
		f, err := Create(fname)
		if err != nil {
			t.Fatalf("could not create output file: %+v", err)
		}
		defer f.Close()

		kv := kvals[0]
		err = f.Put(kv.k, rbase.NewObjString(kv.v))
		if err != nil {
			t.Fatalf("could not write %s: %+v", kv.k, err)
		}

		_, err = f.WriteAt([]byte{1}, kStartBigFile+1)
		if err != nil {
			t.Fatalf("could not write past big-file-mark: %+v", err)
		}
		f.end = kStartBigFile + 1

		kv = kvals[1]
		err = f.Put(kv.k, rbase.NewObjString(kv.v))
		if err != nil {
			t.Fatalf("could not write %s: %+v", kv.k, err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("could not close ROOT file: %+v", err)
		}

		if f.units != 8 {
			t.Fatalf("not a big file")
		}
	}()

	f, err := Open(fname)
	if err != nil {
		t.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	for _, kv := range kvals {
		obj, err := f.Get(kv.k)
		if err != nil {
			t.Fatalf("could not get %s: %+v", kv.k, err)
		}
		if got, want := obj.(*rbase.ObjString).String(), kv.v; got != want {
			t.Fatalf("invalid %s value: got=%q, want=%q", kv.k, got, want)
		}
	}

	if f.units != 8 {
		t.Fatalf("not a big file")
	}

	if !rtests.HasROOT {
		t.Logf("skip test with ROOT/C++")
		return
	}

	const rootls = `#include <iostream>
#include "TFile.h"
#include "TNamed.h"
#include "TObjString.h"

#include <string>

void rootls(const char *fname, const char *kname, const char *v) {
	auto f = TFile::Open(fname);
	auto o = f->Get<TObjString>(kname);
	if (o == NULL) {
		std:cerr << "could not retrieve [" << kname << "]" << std::endl;
		o->ClassName();
	}
	std::cout << "retrieved: [" << kname << "]: [" << o->GetString() << "]" << std::endl;

	auto got = std::string(o->GetString());
	auto want = std::string(v);
	if (got != want) {
		std::cerr << "invalid key value for [" << kname << "]: got=[" << got << "], want=[" << want << "]\n";
		exit(1);
	}
}
`
	for _, kv := range kvals {
		out, err := rtests.RunCxxROOT("rootls", []byte(rootls), fname, kv.k, kv.v)
		if err != nil {
			t.Fatalf("ROOT/C++ could not process file %q:\n%s", fname, string(out))
		}
	}
}
