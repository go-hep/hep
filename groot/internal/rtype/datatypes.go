// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtype // import "go-hep.org/x/hep/groot/internal/rtype"

// FIXME(sbinet): implement a real ROOT float16
// FIXME(sbinet): implement a real ROOT double32

// Float16 is a float32 in memory, written with a truncated mantissa.
type Float16 float32

// Double32 is a float64 in memory, written as a float32 to disk.
type Double32 float64
