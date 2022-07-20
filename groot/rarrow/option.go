// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"git.sr.ht/~sbinet/go-arrow/memory"
)

type config struct {
	mem    memory.Allocator
	chunks int64
	beg    int64
	end    int64
}

func newConfig(opts []Option) *config {
	cfg := &config{
		mem: memory.NewGoAllocator(),
		end: -1,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// Option allows to configure how Records and Tables are constructed
// from input ROOT Trees.
type Option func(*config)

// WithAllocator configures an Arrow value to use the specified memory allocator
// instead of the default Go one.
func WithAllocator(mem memory.Allocator) Option {
	return func(cfg *config) {
		cfg.mem = mem
	}
}

// WithChunk specifies the number of entries to populate Records with.
//
// The default is to populate Records with the whole set of entries the input
// ROOT Tree contains.
func WithChunk(nentries int64) Option {
	return func(cfg *config) {
		cfg.chunks = nentries
	}
}

// WithStart specifies the first entry to read from the input ROOT Tree.
func WithStart(entry int64) Option {
	return func(cfg *config) {
		cfg.beg = entry
	}
}

// WithEnd specifies the last entry (excluded) to read from the input ROOT Tree.
//
// The default (-1) is to read all the entries of the input ROOT Tree.
func WithEnd(entry int64) Option {
	return func(cfg *config) {
		cfg.end = entry
	}
}
