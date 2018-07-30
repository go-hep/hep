// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

// Flags are bit patterns describing detector and simulation states.
type Flags uint32

// Flags for SimCalorimeterHit (CH)
const (
	BitsChLong   Flags = 1 << 31 // long(1) - short(0), (position)
	BitsChBarrel Flags = 1 << 30 // barrel(1) - endcap(0)
	BitsChID1    Flags = 1 << 29 // cellid1 stored
	BitsChPDG    Flags = 1 << 28 // PDG(1) - no PDG(0) (detailed shower contributions) // DEPRECATED: use ChBitStep
	BitsChStep   Flags = 1 << 28 // detailed shower contributions
)

// Flags for the (raw) Calorimeter hits
const (
	BitsRChLong        Flags = 1 << 31 // long(1) - short(0), incl./excl. position
	BitsRChBarrel      Flags = 1 << 30 // barrel(1) - endcap(0)
	BitsRChID1         Flags = 1 << 29 // cellid1 stored
	BitsRChNoPtr       Flags = 1 << 28 // 1: pointer tag not added
	BitsRChTime        Flags = 1 << 27 // 1: time information stored
	BitsRChEnergyError Flags = 1 << 26 // 1: store energy error
)

// Flags for the (raw) tracker data (pulses)
const (
	BitsTRawID1 Flags = 1 << 31 // cellid1 stored
	BitsTRawCM  Flags = 1 << 30 // covariant matrix stored(1) - not stored(0)
)

// Flags for the raw tracker hit
const (
	BitsRThID1 Flags = 1 << 31 // cellid1 stored
)

// Flags for the tracker hit plane
const (
	BitsThPID1 Flags = 1 << 31 // cellid1 stored
)

// Flags for the tracker hit z-cylinder
const (
	BitsThZID1 Flags = 1 << 31 // cellid1 stored
)

// Flags for the SimTrackerHit
const (
	BitsThBarrel   Flags = 1 << 31 // barrel(1) - endcap(0)
	BitsThMomentum Flags = 1 << 30 // momentum of particle stored(1) - not stored(0)
	BitsThID1      Flags = 1 << 29 // cellid1 stored
)

// Flags for the Tracks
const (
	BitsTrHits Flags = 1 << 31 // hits stored(1) - not stored(0)
)

// Flags for the Cluster
const (
	BitsClHits Flags = 1 << 31 // hits stored(1) - not stored(0)
)

// Flags for the TPCHit
const (
	BitsTPCRaw   Flags = 1 << 31 // raw data stored(1) - not stored(0)
	BitsTPCNoPtr Flags = 1 << 30 // 1: pointer tag not added (needed for TrackerHit)
)

// Flags for Relation
const (
	BitsRelWeighted Flags = 1 << 31 // relation has weights
)

// Flags for GenericObject
const (
	BitsGOFixed Flags = 1 << 31 // is fixed size
)

// Test returns whether the given bit is != 0
func (flags Flags) Test(bit Flags) bool {
	return flags&bit != 0
}
