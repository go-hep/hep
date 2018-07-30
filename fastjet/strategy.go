// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

import "fmt"

// Strategy defines the algorithmic strategy used while clustering.
type Strategy int

const (
	N2MinHeapTiledStrategy Strategy = -4
	N2TiledStrategy                 = -3
	N2PoorTiledStrategy             = -2
	N2PlainStrategy                 = -1
	N3DumbStrategy                  = 0
	BestStrategy                    = 1
	NlnNStrategy                    = 2
	NlnN3piStrategy                 = 3
	NlnN4piStrategy                 = 4
	NlnNCam4piStrategy              = 14
	NlnNCam2pi2RStrategy            = 13
	NlnNCamStrategy                 = 12

	PluginStrategy Strategy = 999
)

func (s Strategy) String() string {
	switch s {
	case N2MinHeapTiledStrategy:
		return "N2MinHeapTiled"
	case N2TiledStrategy:
		return "N2Tiled"
	case N2PoorTiledStrategy:
		return "N2PoorTiled"
	case N2PlainStrategy:
		return "N2Plain"
	case N3DumbStrategy:
		return "N3Dumb"
	case BestStrategy:
		return "Best"
	case NlnNStrategy:
		return "NlnN"
	case NlnN3piStrategy:
		return "NlnN3pi"
	case NlnN4piStrategy:
		return "NlnN4pi"
	case NlnNCam4piStrategy:
		return "NlnNCam4pi"
	case NlnNCam2pi2RStrategy:
		return "NlnNCam2pi2R"
	case NlnNCamStrategy:
		return "NlnNCam"

	case PluginStrategy:
		return "Plugin"

	default:
		panic(fmt.Errorf("fastjet: invalid Strategy (%d)", int(s)))
	}
}
