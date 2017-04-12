// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import "log"

func typeFrom(name string) interface{} {
	switch name {
	case "MCParticle":
		return new(McParticles)
	case "SimTrackerHit":
		return new(SimTrackerHits)
	case "SimCalorimeterHit":
		return new(SimCalorimeterHits)
	case "LCFloatVec":
		return new(FloatVec)
	case "LCIntVec":
		return new(IntVec)
	case "LCStrVec":
		return new(StrVec)
	case "RawCalorimeterHit":
		return new(RawCalorimeterHits)
	case "CalorimeterHit":
		return new(CalorimeterHits)

	case "LCGenericObject":
		return new(GenericObject)
	}
	log.Printf("unhandled type %q", name)
	return nil
}
