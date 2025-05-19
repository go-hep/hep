// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"gonum.org/v1/plot/vg/draw"
)

// Options encodes various options to pass to a plot.
type Options func(cfg *config)

// Step kind
type StepsKind byte

const (
	NoSteps StepsKind = iota

	// HiSteps connects two points by following lines: horizontal, vertical, horizontal.
	// Vertical line is placed following the histogram/error-bins informations.
	HiSteps

	// PreSteps connects two points by following lines: vertical, horizontal.
	PreSteps

	// MidSteps connects two points by following lines: horizontal, vertical, horizontal.
	// Vertical line is placed in the middle of the interval.
	MidSteps

	// PostSteps connects two points by following lines: horizontal, vertical.
	PostSteps
)

type config struct {
	bars struct {
		xerrs bool
		yerrs bool
	}
	band   bool
	hinfos HInfos
	log    struct {
		y bool
	}
	glyph draw.GlyphStyle
	steps StepsKind
}

func newConfig(opts []Options) *config {
	cfg := new(config)
	cfg.steps = NoSteps
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

// WithLogY sets whether the plotter in Y should handle log-scale.
func WithLogY(v bool) Options {
	return func(c *config) {
		c.log.y = v
	}
}

// WithXErrBars enables or disables the display of X-error bars.
func WithXErrBars(v bool) Options {
	return func(c *config) {
		c.bars.xerrs = v
	}
}

// WithYErrBars enables or disables the display of Y-error bars.
func WithYErrBars(v bool) Options {
	return func(c *config) {
		c.bars.yerrs = v
	}
}

// WithBand enables or disables the display of a colored band between Y-error bars.
func WithBand(v bool) Options {
	return func(c *config) {
		c.band = v
	}
}

// WithStepsKind sets the style of the connecting line (NoSteps, HiSteps, etc...)
func WithStepsKind(s StepsKind) Options {
	return func(c *config) {
		c.steps = s
	}
}

// WithGlyphStyle sets the glyph style of a plotter.
func WithGlyphStyle(sty draw.GlyphStyle) Options {
	return func(c *config) {
		c.glyph = sty
	}
}

// WithHInfo sets a given histogram info style.
func WithHInfo(v HInfoStyle) Options {
	return func(c *config) {
		c.hinfos.Style = v
	}
}
