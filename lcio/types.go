// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"fmt"
	"log"
)

func typeFrom(name string) interface{} {
	switch name {
	case "MCParticle":
		return new(McParticleContainer)
	case "SimTrackerHit":
		return new(SimTrackerHitContainer)
	case "SimCalorimeterHit":
		return new(SimCalorimeterHitContainer)
	case "LCFloatVec":
		return new(FloatVec)
	case "LCIntVec":
		return new(IntVec)
	case "LCStrVec":
		return new(StrVec)
	case "RawCalorimeterHit":
		return new(RawCalorimeterHitContainer)
	case "CalorimeterHit":
		return new(CalorimeterHitContainer)
	case "TrackerData":
		return new(TrackerDataContainer)
	case "TrackerHit":
		return new(TrackerHitContainer)
	case "TrackerHitPlane":
		return new(TrackerHitPlaneContainer)
	case "TrackerHitZCylinder":
		return new(TrackerHitZCylinderContainer)
	case "TrackerPulse":
		return new(TrackerPulseContainer)
	case "TrackerRawData":
		return new(TrackerRawDataContainer)
	case "Track":
		return new(TrackContainer)
	case "Cluster":
		return new(ClusterContainer)
	case "Vertex":
		return new(VertexContainer)
	case "ReconstructedParticle":
		return new(RecParticleContainer)

	case "LCGenericObject":
		return new(GenericObject)
	case "LCRelation":
		return new(RelationContainer)
	}
	log.Printf("unhandled type %q", name)
	return nil
}

func typeName(t interface{}) string {
	switch t.(type) {
	case *McParticleContainer:
		return "MCParticle"
	case *SimTrackerHitContainer:
		return "SimTrackerHit"
	case *SimCalorimeterHitContainer:
		return "SimCalorimeterHit"
	case *FloatVec:
		return "LCFloatVec"
	case *IntVec:
		return "LCIntVec"
	case *StrVec:
		return "LCStrVec"
	case *RawCalorimeterHitContainer:
		return "RawCalorimeterHit"
	case *CalorimeterHitContainer:
		return "CalorimeterHit"
	case *TrackerDataContainer:
		return "TrackerData"
	case *TrackerHitContainer:
		return "TrackerHit"
	case *TrackerHitPlaneContainer:
		return "TrackerHitPlane"
	case *TrackerHitZCylinderContainer:
		return "TrackerHitZCylinder"
	case *TrackerPulseContainer:
		return "TrackerPulse"
	case *TrackerRawDataContainer:
		return "TrackerRawData"
	case *TrackContainer:
		return "Track"
	case *ClusterContainer:
		return "Cluster"
	case *VertexContainer:
		return "Vertex"
	case *RecParticleContainer:
		return "ReconstructedParticle"

	case *GenericObject:
		return "LCGenericObject"
	case *RelationContainer:
		return "LCRelation"
	}
	panic(fmt.Errorf("lcio: unhandled type %T", t))
}
