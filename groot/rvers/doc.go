// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rvers contains the ROOT version and the classes' versions
// groot is supporting and currently reading.
package rvers // import "go-hep.org/x/hep/groot/rvers"

const (
	// Groot version for STL-based classes.
	// This used to be just StreamerInfo (v=9), but ROOT-6.36.xx bumped TStreamerInfo to v=10
	// and this demonstrated our handling of STL-based classes (the reading part) was subpar.
	//
	// So now we still use the latest version of StreamerInfo, but under a new name to ease
	// later (if any) debugging.
	StreamerBaseSTL = StreamerInfo
)
