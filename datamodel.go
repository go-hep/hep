package fads

import (
	"github.com/go-hep/fmom"
	"github.com/go-hep/hepmc"
)

type Particle interface {
	fmom.P4
	Charge() int32
}

type MissingEt struct {
	MET float32 // missing transverse energy
	Phi float32 // missing energy azimuthal angle
}

// scalar sum of transverse momenta
type ScalarHt float32

// rho energy density
type Rho float32

type Photon struct {
	Pt          float32 // photon transverse momentum
	Eta         float32 // photon pseudo-rapidity
	Phi         float32 // photon azimuthal angle
	Ene         float32 // photon energy
	EhadOverEem float32 // ratio of the hadronic versus electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

type Electron struct {
	Pt          float32 // electron transverse momentum
	Eta         float32 // electron pseudo-rapidity
	Phi         float32 // electron azimuthal angle
	charge      int32   // electron charge
	EhadOverEem float32 // ratio of the hadronic versus electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

type Muon struct {
	Pt     float32 // muon transverse momentum
	Eta    float32 // muon pseudo-rapidity
	Phi    float32 // muon azimuthal angle
	charge int32   // muon charge

	McPart *hepmc.Particle // generated particle
}

type Jet struct {
	Pt  float32 // jet transverse momentum
	Eta float32 // jet pseudo-rapidity
	Phi float32 // jet azimuthal angle
	M   float32 // jet invariant mass

	DeltaEta float32 // jet radius in pseudo-rapidity
	DeltaPhi float32 // jet radius in azimuthal angle

	BTag   byte // 0 or 1 for a jet that has been tagged as containing a heavy quark
	TauTag byte // 0 or 1 for a jet that has been tagged as a tau

	charge int32 // jet charge

	Constituents []Particle        // references to constituents
	McParts      []*hepmc.Particle // references to generated particles
}

// EOF
