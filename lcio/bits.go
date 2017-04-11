// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

// Flags are bit patterns describing detector and simulation states.
type Flags uint32

// Flags for SimCalorimeterHit (CH)
const (
	ChBitLong   Flags = 31 // long(1) - short(0), (position)
	ChBitBarrel Flags = 30 // barrel(1) - endcap(0)
	ChBitID1    Flags = 29 // cellid1 stored
	ChBitPDG    Flags = 28 // PDG(1) - no PDG(0) (detailed shower contributions) // DEPRECATED: use ChBitStep
	ChBitStep   Flags = 28 // detailed shower contributions
)

// Flags for the (raw) Calorimeter hits
const (
	RChBitLong        Flags = 31 // long(1) - short(0), incl./excl. position
	RChBitBarrel      Flags = 30 // barrel(1) - endcap(0)
	RChBitID1         Flags = 29 // cellid1 stored
	RChBitNoPtr       Flags = 28 // 1: pointer tag not added
	RChBitTime        Flags = 27 // 1: time information stored
	RChBitEnergyError Flags = 26 // 1: store energy error
)

// Flags for the (raw) tracker data (pulses)
const (
	TRawBitID1 Flags = 31 // cellid1 stored
	TRawBitCM  Flags = 30 // covariant matrix stored(1) - not stored(0)
)

// Flags for the raw tracker hit
const (
	RThBitID1 Flags = 31 // cellid1 stored
)

// Flags for the tracker hit plane
const (
	RThPBitID1 Flags = 31 // cellid1 stored
)

// Flags for the tracker hit z-cylinder
const (
	RThZBitID1 Flags = 31 // cellid1 stored
)

// Flags for the SimTrackerHit
const (
	ThBitBarrel   Flags = 31 // barrel(1) - endcap(0)
	ThBitMomentum Flags = 30 // momentum of particle stored(1) - not stored(0)
	ThBitID1      Flags = 29 // cellid1 stored
)

// Flags for the Tracks
const (
	TrBitHits Flags = 31 // hits stored(1) - not stored(0)
)

// Flags for the Cluster
const (
	ClBitHits Flags = 31 // hits stored(1) - not stored(0)
)

// Flags for the TPCHit
const (
	TPCBitRaw   Flags = 31 // raw data stored(1) - not stored(0)
	TPCBitNoPtr Flags = 30 // 1: pointer tag not added (needed for TrackerHit)
)

// Flags for Relation
const (
	RelWeighted Flags = 31 // relation has weights
)

// Flags for GenericObject
const (
	GOBitFixed Flags = 31 // is fixed size
)

// Test returns whether the given bit is != 0
func (flags Flags) Test(bit Flags) bool {
	return flags&(1<<bit) != 0
}
