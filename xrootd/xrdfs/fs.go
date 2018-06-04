// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrdfs contains structures representing the XRootD-based filesystem.
package xrdfs

import (
	"context"
)

// FileSystem implements access to a collection of named files over XRootD.
type FileSystem interface {
	Dirlist(ctx context.Context, path string) ([]EntryStat, error)
}
