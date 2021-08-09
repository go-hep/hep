// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/font/liberation"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg/draw"
)

var (
	// DefaultStyle is the default style used for hplot plots.
	DefaultStyle Style
)

// Style stores a given plot style.
type Style struct {
	Fonts struct {
		Name    string    // font name of this style
		Default font.Font // font used for the plot
		Title   font.Font // font used for the plot title
		Label   font.Font // font used for the plot labels
		Legend  font.Font // font used for the plot legend
		Tick    font.Font // font used for the plot's axes' ticks

		Cache *font.Cache // cache of fonts for this plot.
	}
	TextHandler text.Handler
}

// Apply setups the plot p with the current style.
func (s *Style) Apply(p *Plot) {
	p.Plot.Title.TextStyle.Font = s.Fonts.Title
	p.Plot.X.Label.TextStyle.Font = s.Fonts.Label
	p.Plot.Y.Label.TextStyle.Font = s.Fonts.Label
	p.Plot.X.Tick.Label.Font = s.Fonts.Tick
	p.Plot.Y.Tick.Label.Font = s.Fonts.Tick
	p.Plot.Legend.TextStyle.Font = s.Fonts.Legend
	p.Plot.Legend.YPosition = draw.PosCenter
	p.Plot.TextHandler = s.TextHandler
}

func (s *Style) reset(fnt font.Font) {
	plot.DefaultFont = fnt
}

// NewStyle creates a new style with the given font cache.
func NewStyle(fnt font.Font, cache *font.Cache) (Style, error) {
	var sty Style

	sty.Fonts.Name = fnt.Name()
	sty.Fonts.Default = fnt
	sty.Fonts.Title = fnt
	sty.Fonts.Label = fnt
	sty.Fonts.Legend = fnt
	sty.Fonts.Tick = fnt
	sty.Fonts.Cache = cache

	err := sty.init(fnt.Name())
	if err != nil {
		return sty, err
	}

	sty.TextHandler = text.Plain{
		Fonts: sty.Fonts.Cache,
	}

	return sty, nil
}

func (sty *Style) init(name string) error {
	sty.Fonts.Name = name
	for _, t := range []struct {
		ft   *font.Font
		size font.Length
	}{
		{&sty.Fonts.Title, 12},
		{&sty.Fonts.Label, 12},
		{&sty.Fonts.Legend, 12},
		{&sty.Fonts.Tick, 10},
	} {
		if !sty.Fonts.Cache.Has(*t.ft) {
			return fmt.Errorf("hplot: no font %v in cache", *t.ft)
		}
		t.ft.Size = t.size
	}

	return nil
}

func init() {
	var (
		cache = font.NewCache(liberation.Collection())
		err   error
	)
	DefaultStyle, err = NewStyle(
		font.Font{
			Typeface: "Liberation",
			Variant:  "Sans",
		},
		cache,
	)
	if err != nil {
		panic(err)
	}
}
