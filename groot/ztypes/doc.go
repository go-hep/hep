// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package ztypes holds all the types registered with the rtypes factory.
package ztypes // import "go-hep.org/x/hep/groot/ztypes"

import (
	_ "go-hep.org/x/hep/groot/rbase"
	_ "go-hep.org/x/hep/groot/rcont"
	_ "go-hep.org/x/hep/groot/rdict"
	_ "go-hep.org/x/hep/groot/rhist"
	_ "go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/rpad"
	_ "go-hep.org/x/hep/groot/rphys"
	_ "go-hep.org/x/hep/groot/rtree"

	// ROOT::Experimental
	_ "go-hep.org/x/hep/groot/exp/rntup"
)
