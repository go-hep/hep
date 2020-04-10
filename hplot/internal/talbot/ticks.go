// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package talbot

import (
	"math"
	"strconv"

	"gonum.org/v1/plot"
)

// Ticks returns Ticks in the specified range.
func Ticks(min, max float64, n int) []plot.Tick {
	if max <= min {
		panic("illegal range")
	}

	suggestedTicks := n

	labels, step, q, mag := talbotLinHanrahan(min, max, suggestedTicks, withinData, nil, nil, nil)
	majorDelta := step * math.Pow10(mag)
	if q == 0 {
		// Simple fall back was chosen, so
		// majorDelta is the label distance.
		majorDelta = labels[1] - labels[0]
	}

	// Choose a reasonable, but ad
	// hoc formatting for labels.
	fc := byte('f')
	var off int
	if mag < -1 || 6 < mag {
		off = 1
		fc = 'g'
	}
	if math.Trunc(q) != q {
		off += 2
	}
	prec := minInt(6, maxInt(off, -mag))
	var ticks []plot.Tick
	for _, v := range labels {
		ticks = append(ticks, plot.Tick{Value: v, Label: strconv.FormatFloat(v, fc, prec, 64)})
	}
	lastMajor := ticks[len(ticks)-1]

	var minorDelta float64
	// See talbotLinHanrahan for the values used here.
	switch step {
	case 1, 2.5:
		minorDelta = majorDelta / 5
	case 2, 3, 4, 5:
		minorDelta = majorDelta / step
	default:
		if majorDelta/2 < dlamchP {
			return ticks
		}
		minorDelta = majorDelta / 2
	}

	// Find the first minor tick not greater
	// than the lowest data value.
	var i float64
	for labels[0]+(i-1)*minorDelta > min {
		i--
	}
	// Add ticks at minorDelta intervals when
	// they are not within minorDelta/2 of a
	// labelled tick.
	for {
		val := labels[0] + i*minorDelta
		if val > max {
			break
		}
		found := false
		for _, t := range ticks {
			if math.Abs(t.Value-val) < minorDelta/2 {
				found = true
			}
		}
		if !found {
			ticks = append(ticks, plot.Tick{Value: val})
		}
		i++
	}

	return ticks
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
