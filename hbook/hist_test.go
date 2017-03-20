// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"math/rand"
)

var seed = rand.Uint32()
var g_seed = &seed

// rnd is a simple minded non-thread-safe version of math/rand.Float64
func rnd() float64 {
	ss := *g_seed
	ss += ss
	ss ^= 1
	if int32(ss) < 0 {
		ss ^= 0x88888eef
	}
	*g_seed = ss
	return float64(*g_seed%95) / float64(95)
}

func panics(fn func()) (panicked bool, message string) {
	defer func() {
		r := recover()
		panicked = r != nil
		message = fmt.Sprint(r)
	}()
	fn()
	return
}
