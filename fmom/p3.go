// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fmom

// Vec3 is a 3-dim vector.
type Vec3 [3]float64

func (vec *Vec3) X() float64 {
	return vec[0]
}

func (vec *Vec3) Y() float64 {
	return vec[1]
}

func (vec *Vec3) Z() float64 {
	return vec[2]
}
