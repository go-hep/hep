// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"
	"testing"
)

func TestSave(t *testing.T) {
	p := New()
	p.Title.Text = "my title"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	for _, tc := range []struct {
		name  string
		files []string
		want  error
	}{
		{
			name:  "empty-fnames",
			files: []string{},
			want:  fmt.Errorf(`hplot: need at least 1 file name`),
		},
		{
			name:  "invalid-format",
			files: []string{"invalid-format"},
			want:  fmt.Errorf(`hplot: could not save plot: hplot: could not create canvas: unsupported format: ""`),
		},
		{
			name:  "unknown-format",
			files: []string{"file.txt"},
			want:  fmt.Errorf(`hplot: could not save plot: hplot: could not create canvas: unsupported format: "txt"`),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := Save(p, -1, -1, tc.files...)
			if err == nil {
				t.Fatalf("expected an error")
			}

			if got, want := err.Error(), tc.want.Error(); got != want {
				t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
			}
		})
	}
}
