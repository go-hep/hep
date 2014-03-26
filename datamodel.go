package fads

import (
	"github.com/go-hep/fmom"
	"github.com/go-hep/hepmc"
)

type Particle interface {
	P4() fmom.P4
	Charge() int32
}

type MissingEt struct {
	MET float64 // missing transverse energy
	Phi float64 // missing energy azimuthal angle
}

// scalar sum of transverse momenta
type ScalarHt float64

// rho energy density
type Rho float64

type mcParticle struct {
	pid    int32 // pdg id number
	status int32 // particle status
	ispu   byte  // 0 or 1 for particles from pile-up interactions

	m1 int // particle 1st mother
	m2 int // particle 2nd mother
	d1 int // particle 1st daughter
	d2 int // particle 2nd daughter

	charge int     // particle charge
	mass   float64 // particle mass

	mom fmom.PxPyPzE // particle momentum (px,py,pz,e)
	pt  float64      // particle transverse momentum
	eta float64      // particle pseudo-rapidity
	phi float64      // particle azimuthal angle

	rapidity float64 // particle rapidity

	pos [4]float64 // particle vertex position (t,x,y,z)
}

type Photon struct {
	p4    fmom.PtEtaPhiM // photon momentum
	ehoem float64        // ratio of the hadronic over electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

type Electron struct {
	p4          fmom.PtEtaPhiM // electron momentum
	charge      int32          // electron charge
	EhadOverEem float64        // ratio of the hadronic versus electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

type Muon struct {
	p4     fmom.PtEtaPhiM // muon momentum
	charge int32          // muon charge

	McPart *hepmc.Particle // generated particle
}

type Jet struct {
	p4     fmom.PtEtaPhiM // jet momentum
	charge int32          // jet charge

	deta float64 // jet radius in pseudo-rapidity
	dphi float64 // jet radius in azimuthal angle

	btag   byte // 0 or 1 for a jet that has been tagged as containing a heavy quark
	tautag byte // 0 or 1 for a jet that has been tagged as a tau

	Constituents []Particle        // pointers to constituents
	McParts      []*hepmc.Particle // pointers to generated particles
}

type Track struct {
	pid    int32          // HEP ID number
	charge int32          // track charge
	p4     fmom.PtEtaPhiM // track momentum

	eta float64 // track pseudo-rapidity at the tracker edge
	phi float64 // track azimuthal angle at the tracker edge

	x float64 // track vertex position
	y float64 // track vertex position
	z float64 // track vertex position

	xout float64 // track vertex position at the tracker edge
	yout float64 // track vertex position at the tracker edge
	zout float64 // track vertex position at the tracker edge

	mc *hepmc.Particle // pointer to generated particle
}

type Tower struct {
	p4   fmom.EtEtaPhiM // calorimeter tower momentum
	ene  float64        // calorimeter tower energy
	eem  float64        // calorimeter tower electromagnetic energy
	ehad float64        // calorimter tower hadronic energy

	edges [4]float64 // calorimeter tower edges

	mcparts []*hepmc.Particle // pointers to generated particles
}

type Candidate struct {
	pid            int32   // HEP ID number
	status         int32   // particle status
	m1, m2, d1, d2 int32   // particle mothers and daughters
	charge         int32   // particle charge
	mass           float64 // particle mass

	ispu          byte // 0 or 1 for particles from pile-up interactions
	isconstituent byte // 0 or 1 for particles being constituents
	btag          byte // 0 or 1 for a candidate that has been tagged as containing a heavy quark
	tautag        byte // 0 or 1 for a candidate that has been tagged as a tau

	eem  float64 // electromagnetic energy
	ehad float64 // hadronic energy

	edges [4]float64
	deta  float64
	dphi  float64

	mom  fmom.P4
	pos  fmom.P4
	area fmom.P4

	arr []Candidate
}

// EOF
