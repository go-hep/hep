// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package huntex // import "go-hep.org/x/hep/hplot/huntex"

import (
	"gonum.org/v1/plot/vg"
)

// Canvas implements the vg.Canvas interface, rendering LaTeX strings
// with their UTF-8 equivalent.
type Canvas struct {
	vg.Canvas

	rep replace
}

// FillString fills in text at the specified
// location using the given font.
// If the font size is zero, the text is not drawn.
func (c Canvas) FillString(f vg.Font, pt vg.Point, text string) {
	txt := c.rep.replace(text)
	c.Canvas.FillString(f, pt, txt)
}

var (
	_ vg.Canvas = (*Canvas)(nil)
)
