// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ntroot_test

import (
	"fmt"
	"testing"

	"go-hep.org/x/hep/hbook/ntup/ntroot"
)

func TestOpen(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree string
		err  error
	}{
		{
			name: "../../../groot/testdata/simple.root",
			tree: "tree",
		},
		{
			name: "../../../groot/testdata/graphs.root",
			tree: "tg",
			err:  fmt.Errorf(`ROOT object "tg" is not a tree`),
		},
		{
			name: "../../../groot/testdata/simple.root",
			tree: "treeXXX",
			err:  fmt.Errorf(`could not find ROOT tree "treeXXX": riofs: simple.root: could not find key "treeXXX;9999"`),
		},
		{
			name: "../../../groot/testdata/simple.rootXXX",
			tree: "tree",
			err:  fmt.Errorf(`could not open ROOT file: riofs: unable to open "../../../groot/testdata/simple.rootXXX": riofs: no ROOT plugin to open [../../../groot/testdata/simple.rootXXX] (scheme=)`),
		},
	} {
		t.Run(tc.name+":"+tc.tree, func(t *testing.T) {
			nt, err := ntroot.Open(tc.name, tc.tree)
			if err == nil {
				_ = nt.DB().Close()
			}

			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
				}
			case err != nil && tc.err == nil:
				t.Fatalf("unexpected error: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error (got=%v, want=%v)", err, tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}
		})
	}
}
