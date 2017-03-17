// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

// SpinState contains the total spin, spin and orbital angular momentum
type SpinState struct {
	TotalSpin float64 // total spin
	Spin      float64
	OrbAngMom float64 // orbital angular momentum
}
