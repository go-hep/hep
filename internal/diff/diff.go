// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package diff provides basic text comparison (like Unix's diff(1)).
package diff // import "go-hep.org/x/hep/internal/diff"

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/diff"
)

// Format returns a formatted diff of the two texts,
// showing the entire text and the minimum line-level
// additions and removals to turn got into want.
// (That is, lines only in got appear with a leading -,
// and lines only in want appear with a leading +.)
func Format(got, want string) string {
	o := new(strings.Builder)
	err := diff.Text("a/got", "b/want", got, want, o)
	if err != nil {
		panic(err)
	}
	return o.String()
}

// Files returns a formatted diff of the two texts from the provided
// two file names.
// Files returns nil if they compare equal.
func Files(got, want string) error {
	g, err := os.ReadFile(got)
	if err != nil {
		return fmt.Errorf("diff: could not read chk file %q: %w", got, err)
	}
	w, err := os.ReadFile(want)
	if err != nil {
		return fmt.Errorf("diff: could not read ref file %q: %w", want, err)
	}

	if got, want := string(g), string(w); got != want {
		return fmt.Errorf("diff: files differ:\n%s", Format(got, want))
	}

	return nil
}
