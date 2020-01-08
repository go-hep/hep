// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rcmd provides helper functions containing the logic of various root-xyz commands.
package rcmd // import "go-hep.org/x/hep/groot/rcmd"

import (
	"path/filepath"
	"strings"

	"golang.org/x/xerrors"
)

func splitArg(cmd string) (fname, sel string, err error) {
	fname = cmd
	prefix := ""
	for _, p := range []string{"https://", "http://", "root://", "file://"} {
		if strings.HasPrefix(cmd, p) {
			prefix = p
			break
		}
	}
	fname = fname[len(prefix):]

	vol := filepath.VolumeName(fname)
	if vol != fname {
		fname = fname[len(vol):]
	}

	if strings.Count(fname, ":") > 1 {
		return "", "", xerrors.Errorf("root-cp: too many ':' in %q", cmd)
	}

	i := strings.LastIndex(fname, ":")
	switch {
	case i > 0:
		sel = fname[i+1:]
		fname = fname[:i]
	default:
		sel = ".*"
	}
	if sel == "" {
		sel = ".*"
	}
	fname = prefix + vol + fname
	switch {
	case strings.HasPrefix(sel, "/"):
	case strings.HasPrefix(sel, "^/"):
	case strings.HasPrefix(sel, "^"):
		sel = "^/" + sel[1:]
	default:
		sel = "/" + sel
	}
	return fname, sel, err
}
