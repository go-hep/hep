// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import "testing"

func TestStreamKind(t *testing.T) {
	for _, tc := range []struct {
		kind StreamKind
		want string
	}{
		{
			kind: ObjectWise,
			want: "object-wise",
		},
		{
			kind: MemberWise,
			want: "member-wise",
		},
		{
			kind: 255,
			want: "0xff",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.kind.String()
			if got != tc.want {
				t.Fatalf("invalid kind: got=%q, want=%q", got, tc.want)
			}
		})
	}
}
