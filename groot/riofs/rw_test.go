// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

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
