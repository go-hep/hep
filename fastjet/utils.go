// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

func imin(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func imax(i, j int) int {
	if i > j {
		return i
	}
	return j
}

// ByPt sorts jets by descending Pt
type ByPt []Jet

func (p ByPt) Len() int {
	return len(p)
}

func (p ByPt) Less(i, j int) bool {
	return p[j].Pt() < p[i].Pt()
}

func (p ByPt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
