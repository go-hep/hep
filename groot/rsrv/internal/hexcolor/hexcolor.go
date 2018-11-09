// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// http://code.google.com/p/gorilla/source/browse/LICENSE

// Taken from
// http://code.google.com/p/gorilla/source/browse/color/hex.go?r=8029f2a42eab5f62c298e5fbe942ed472509a090
// (it disappeared after that revision)

// Enhanced by Peter Waller <peter@pwaller.net>

package hexcolor

import (
	"fmt"
	"image/color"
	"log"
	"strconv"
)

// HexModel converts any Color to an Hex color.
var HexModel = color.ModelFunc(hexModel)

// Hex represents an RGB color in hexadecimal format.
//
// The length must be 3 or 6 characters, preceded or not by a '#'.
type Hex string

// RGBA returns the alpha-premultiplied red, green, blue and alpha values
// for the Hex.
func (c Hex) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, a := HexToRGBA(c)
	return uint32(r) * 0x101, uint32(g) * 0x101, uint32(b) * 0x101, uint32(a) * 0x101
}

// hexModel converts a Color to Hex.
func hexModel(c color.Color) color.Color {
	if _, ok := c.(Hex); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBAToHex(uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

// RGBAToHex converts an RGBA to a Hex string.
// If a == 255, the A is not specified in the hex string
func RGBAToHex(r, g, b, a uint8) Hex {
	if a == 255 {
		return Hex(fmt.Sprintf("#%02X%02X%02X", r, g, b))
	}
	return Hex(fmt.Sprintf("#%02X%02X%02X%02X", r, g, b, a))
}

// Converts an Hex string to RGBA.
// If alpha is not specified, it defaults to 255
// If it is not a valid hexadecimal number of the right width, a horrible yellow
// color is returned
func HexToRGBA(h Hex) (r, g, b, a uint8) {
	if len(h) > 0 && h[0] == '#' {
		h = h[1:]
	}
	r, g, b, a = 0xBC, 0xB6, 0x04, 0x88
	rgb, err := strconv.ParseUint(string(h), 16, 32)
	if err != nil {
		log.Printf("Invalid color: %q err=%v", string(h), err)
		return
	}
	// Some of the 0xFF masks are presumably un-necessary but are left here
	// to remind what is happening.
	switch len(h) {
	case 3:
		r = uint8(rgb>>8) & 0xF * 0x11
		g = uint8(rgb>>4) & 0xF * 0x11
		b = uint8(rgb>>0) & 0xF * 0x11
		a = 0xFF
	case 4:
		r = uint8(rgb>>12) & 0xF * 0x11
		g = uint8(rgb>>8) & 0xF * 0x11
		b = uint8(rgb>>4) & 0xF * 0x11
		a = uint8(rgb>>0) & 0xF * 0x11
	case 6:
		r = uint8(rgb>>16) & 0xFF
		g = uint8(rgb>>8) & 0xFF
		b = uint8(rgb>>0) & 0xFF
		a = 0xFF
	case 8:
		r = uint8(rgb>>24) & 0xFF
		g = uint8(rgb>>16) & 0xFF
		b = uint8(rgb>>8) & 0xFF
		a = uint8(rgb>>0) & 0xFF
	default:
		log.Printf("Invalid color: %q err: invalid length %d %v",
			string(h), len(h), " != {3,4,6,8}")
	}
	return
}
