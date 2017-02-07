// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// Object is the general handle to any hbook data analysis object.
type Object interface {
	Annotation() Annotation
	Name() string
}
