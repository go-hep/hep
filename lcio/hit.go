// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

// Hit is an abstract Hit in the LCIO event data model.
type Hit interface {
	GetCellID0() int32
	GetCellID1() int32
}
