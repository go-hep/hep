// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package htex

import (
	"runtime"

	"golang.org/x/sync/errgroup"
)

// GoHandler is a Latex handler that compiles
// Latex document in background goroutines.
type GoHandler struct {
	ch   chan int // throttling channel
	grp  *errgroup.Group
	hdlr Handler
}

// NewGoHandler creates a new Latex handler that compiles Latex
// document in the background with the cmd executable.
//
// The handler allows for up to n concurrent compilations.
// If n<=0, the concurrency will be set to the number of cores+1.
func NewGoHandler(n int, cmd string) *GoHandler {
	if n <= 0 {
		n = runtime.NumCPU() + 1
	}

	h := &GoHandler{
		ch:   make(chan int, n),
		grp:  new(errgroup.Group),
		hdlr: NewHandler(cmd),
	}

	return h
}

// CompileLatex compiles the provided .tex document.
func (gh *GoHandler) CompileLatex(fname string) error {
	gh.grp.Go(func() error {
		gh.ch <- 1
		defer func() { <-gh.ch }()
		return gh.hdlr.CompileLatex(fname)
	})
	return nil
}

func (gh *GoHandler) Wait() error {
	return gh.grp.Wait()
}
