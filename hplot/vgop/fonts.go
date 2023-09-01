// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vgop // import "go-hep.org/x/hep/hplot/vgop"

import (
	xfnt "golang.org/x/image/font"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type fontCtx struct {
	fonts map[fontID]font.Face
}

type fontID struct {
	Typeface font.Typeface `json:"typeface,omitempty"`
	Variant  font.Variant  `json:"variant,omitempty"`
	Style    xfnt.Style    `json:"style,omitempty"`
	Weight   xfnt.Weight   `json:"weight,omitempty"`
	Size     vg.Length     `json:"size,omitempty"`
}
