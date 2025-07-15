// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rvers contains the ROOT version and the classes' versions
// groot is supporting and currently reading.
package rvers // import "go-hep.org/x/hep/groot/rvers"

const (
	// Groot version for STL-based classes.
	// This used to be just StreamerInfo (v=9), but ROOT-6.36.xx bumped TStreamerInfo to v=10.
	// And, apparently, the "artificial" class version we see in front of std::xxx classes is
	// still v=9.
	// So: make up our own class version name for the streamer elements.
	StreamerBaseSTL = 9
)
