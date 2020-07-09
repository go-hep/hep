// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdio

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want URL
	}{
		{
			name: "root://example.org/file1.root",
			want: URL{
				Addr: "example.org",
				User: "",
				Path: "/file1.root",
			},
		},
		{
			name: "xroot://example.org/file1.root",
			want: URL{
				Addr: "example.org",
				User: "",
				Path: "/file1.root",
			},
		},
		{
			name: "root://example.org//file1.root",
			want: URL{
				Addr: "example.org",
				User: "",
				Path: "/file1.root",
			},
		},
		{
			name: "root://bob@example.org/file1.root",
			want: URL{
				Addr: "example.org",
				User: "bob",
				Path: "/file1.root",
			},
		},
		{
			name: "root://bob:s3cr3t@example.org/file1.root",
			want: URL{
				Addr: "example.org",
				User: "bob",
				Path: "/file1.root",
			},
		},
		{
			name: "root://bob:s3cr3t@example.org:1024/file1.root",
			want: URL{
				Addr: "example.org:1024",
				User: "bob",
				Path: "/file1.root",
			},
		},
		{
			name: "root://bob:s3cr3t@example.org:1024/dir/file1.root",
			want: URL{
				Addr: "example.org:1024",
				User: "bob",
				Path: "/dir/file1.root",
			},
		},
		{
			name: "root://example.org/file1.%c.root",
			want: URL{
				Addr: "example.org",
				Path: "/file1.%c.root",
			},
		},
		{
			name: "root://localhost:1094/file1.root",
			want: URL{
				Addr: "localhost:1094",
				Path: "/file1.root",
			},
		},
		{
			name: "root://127.0.0.1:1094/file1.root",
			want: URL{
				Addr: "127.0.0.1:1094",
				Path: "/file1.root",
			},
		},
		{
			name: "root://[2001:db8:85a3:8d3:1319:8a2e:370:7348]/file1.root",
			want: URL{
				Addr: "[2001:db8:85a3:8d3:1319:8a2e:370:7348]",
				Path: "/file1.root",
			},
		},
		{
			name: "root://[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1094/file1.root",
			want: URL{
				Addr: "[2001:db8:85a3:8d3:1319:8a2e:370:7348]:1094",
				Path: "/file1.root",
			},
		},
		{
			name: "root://[::1]:1094/file1.root",
			want: URL{
				Addr: "[::1]:1094",
				Path: "/file1.root",
			},
		},
		{
			name: "root://[::1%lo0]:1094/file1.root",
			want: URL{
				Addr: "[::1%lo0]:1094",
				Path: "/file1.root",
			},
		},
		{
			name: "file:///dir/file1.root",
			want: URL{
				Addr: "",
				Path: "/dir/file1.root",
			},
		},
		{
			// this is an incorrectly written URI.
			// unfortunately, we can't distinguish it from other well-formed URIs
			name: "file://dir/file1.root",
			want: URL{
				Addr: "dir",
				Path: "/file1.root",
			},
		},
		{
			name: "dir/file1.root",
			want: URL{
				Addr: "",
				Path: "dir/file1.root",
			},
		},
		{
			name: "root://example.org:1:2/file1.root",
			err:  fmt.Errorf(`could not parse URI "root://example.org:1:2/file1.root": could not extract host+port from URI: address example.org:1:2: too many colons in address`),
		},
		{
			name: "root://user@example.org:1:2/file1.root",
			err:  fmt.Errorf(`could not parse URI "root://user@example.org:1:2/file1.root": could not extract host+port from URI: address example.org:1:2: too many colons in address`),
		},
		{
			name: "root://user:pass@example.org:1:2/file1.root",
			err:  fmt.Errorf(`could not parse URI "root://user:pass@example.org:1:2/file1.root": could not extract host+port from URI: address example.org:1:2: too many colons in address`),
		},
		{
			name: "root://user:pass@[::1]:1:2/file1.root",
			err:  fmt.Errorf(`could not parse URI "root://user:pass@[::1]:1:2/file1.root": could not extract host+port from URI: address [::1]:1:2: too many colons in address`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(tc.name)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v",
						got, want,
					)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("could not parse URI: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %+v", tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}

			if got, want := got, tc.want; got != want {
				t.Fatalf("invalid parse result:\ngot= %#v\nwant=%#v",
					got, want,
				)
			}
		})
	}
}
