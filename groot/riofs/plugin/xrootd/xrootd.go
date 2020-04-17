// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrootd is a plugin for riofs.Open to support opening ROOT files over xrootd.
package xrootd

import (
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/xrootd/xrdio"
)

func init() {
	riofs.Register("root", openFile)
	riofs.Register("xroot", openFile)
}

func openFile(path string) (riofs.Reader, error) {
	return xrdio.Open(path)
}

var (
	_ riofs.Reader = (*xrdio.File)(nil)
	_ riofs.Writer = (*xrdio.File)(nil)
)
