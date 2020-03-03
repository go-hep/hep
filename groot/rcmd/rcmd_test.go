// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd // import "go-hep.org/x/hep/groot/rcmd"

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSplitArg(t *testing.T) {
	for _, tc := range []struct {
		cmd   string
		fname string
		sel   string
		err   error
	}{
		{
			cmd:   "file.root",
			fname: "file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root",
			fname: "dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "/dir/sub/file.root",
			fname: "/dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "../dir/sub/file.root",
			fname: "../dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:hist",
			fname: "dir/sub/file.root",
			sel:   "/hist",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:hist*",
			fname: "dir/sub/file.root",
			sel:   "/hist*",
			err:   nil,
		},
		{
			cmd:   "dir/sub/file.root:",
			fname: "dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "file://dir/sub/file.root:",
			fname: "file://dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "https://dir/sub/file.root",
			fname: "https://dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "http://dir/sub/file.root",
			fname: "http://dir/sub/file.root",
			sel:   "/.*",
			err:   nil,
		},
		{
			cmd:   "https://dir/sub/file.root:hist*",
			fname: "https://dir/sub/file.root",
			sel:   "/hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:hist*",
			fname: "root://dir/sub/file.root",
			sel:   "/hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:/hist*",
			fname: "root://dir/sub/file.root",
			sel:   "/hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:^/hist*",
			fname: "root://dir/sub/file.root",
			sel:   "^/hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:^hist*",
			fname: "root://dir/sub/file.root",
			sel:   "^/hist*",
			err:   nil,
		},
		{
			cmd:   "root://dir/sub/file.root:/^hist*",
			fname: "root://dir/sub/file.root",
			sel:   "/^hist*",
			err:   nil,
		},
		{
			cmd: "dir/sub/file.root:h:h",
			err: fmt.Errorf("root-cp: too many ':' in %q", "dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root:h:h",
			err: fmt.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root:h:h"),
		},
		{
			cmd: "root://dir/sub/file.root::h:",
			err: fmt.Errorf("root-cp: too many ':' in %q", "root://dir/sub/file.root::h:"),
		},
	} {
		t.Run(tc.cmd, func(t *testing.T) {
			fname, sel, err := splitArg(tc.cmd)
			switch {
			case err != nil && tc.err != nil:
				if !reflect.DeepEqual(err.Error(), tc.err.Error()) {
					t.Fatalf("got err=%v, want=%v", err, tc.err)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("got err=%v, want=%v", err, tc.err)
			case err == nil && tc.err != nil:
				t.Fatalf("got err=%v, want=%v", err, tc.err)
			}

			if got, want := fname, tc.fname; got != want {
				t.Fatalf("fname=%q, want=%q", got, want)
			}

			if got, want := sel, tc.sel; got != want {
				t.Fatalf("selection=%q, want=%q", got, want)
			}
		})
	}
}
