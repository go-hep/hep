// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot

import (
	"github.com/golang/freetype/truetype"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/fonts"
)

var (
	// DefaultStyle is the default style used for hplot plots.
	DefaultStyle Style
)

// Style stores a given plot style.
type Style struct {
	Fonts struct {
		Name   string  // font name of this style
		Title  vg.Font // font used for the plot title
		Label  vg.Font // font used for the plot labels
		Legend vg.Font // font used for the plot legend
		Tick   vg.Font // font used for the plot's axes' ticks
	}
}

// Apply setups the plot p with the current style.
func (s *Style) Apply(p *Plot) {
	p.Plot.Title.TextStyle.Font = s.Fonts.Title
	p.Plot.X.Label.TextStyle.Font = s.Fonts.Label
	p.Plot.Y.Label.TextStyle.Font = s.Fonts.Label
	p.Plot.X.Tick.Label.Font = s.Fonts.Tick
	p.Plot.Y.Tick.Label.Font = s.Fonts.Tick
	p.Plot.Legend.TextStyle.Font = s.Fonts.Legend
}

func (s *Style) reset(name string) {
	plot.DefaultFont = name
}

// NewStyle creates a new style with the given truetype font.
func NewStyle(name string, ft *truetype.Font) (Style, error) {
	var sty Style
	err := sty.init(name, ft)
	if err != nil {
		return sty, err
	}

	return sty, nil
}

func (sty *Style) init(name string, ft *truetype.Font) error {
	sty.Fonts.Name = name
	vg.AddFont(name, ft)
	for _, t := range []struct {
		ft   *vg.Font
		size vg.Length
	}{
		{&sty.Fonts.Title, 12},
		{&sty.Fonts.Label, 12},
		{&sty.Fonts.Legend, 12},
		{&sty.Fonts.Tick, 10},
	} {
		ft, err := vg.MakeFont(sty.Fonts.Name, t.size)
		if err != nil {
			return err
		}
		*t.ft = ft
	}

	return nil
}

func init() {
	vgfonts, err := fonts.Asset("LiberationSans-Regular.ttf")
	if err != nil {
		panic(err)
	}

	ft, err := truetype.Parse(vgfonts)
	if err != nil {
		panic(err)
	}

	err = DefaultStyle.init("Helvetica", ft)
	if err != nil {
		panic(err)
	}
}
