// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

// Measurement holds a value and its associated error
type Measurement struct {
	Value float64
	Sigma float64
}
