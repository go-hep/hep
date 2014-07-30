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
	Pid    int32 // pdg id number
	Status int32 // particle status
	IsPU   byte  // 0 or 1 for particles from pile-up interactions

	M1 int // particle 1st mother
	M2 int // particle 2nd mother
	D1 int // particle 1st daughter
	D2 int // particle 2nd daughter

	McCharge int     // particle charge
	Mass     float64 // particle mass

	Mom fmom.PxPyPzE // particle momentum (px,py,pz,e)

	Pt  float64 // particle transverse momentum
	Eta float64 // particle pseudo-rapidity
	Phi float64 // particle azimuthal angle

	Rapidity float64 // particle rapidity

	Pos [4]float64 // particle vertex position (t,x,y,z)
}

func (mc *mcParticle) P4() fmom.P4 {
	return &mc.Mom
}

func (mc *mcParticle) Charge() int32 {
	return int32(mc.McCharge)
}

type Photon struct {
	Mom    fmom.PtEtaPhiM // photon momentum (mass=0.0)
	EhoEem float64        // ratio of the hadronic over electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

func (pho *Photon) P4() fmom.P4 {
	return &pho.Mom
}

func (pho *Photon) Charge() int32 {
	return 0
}

type Electron struct {
	Mom       fmom.PtEtaPhiM // electron momentum (mass=0.0)
	EleCharge int32          // electron charge
	EhoEem    float64        // ratio of the hadronic versus electromagnetic energy deposited in the calorimeter

	McPart *hepmc.Particle // generated particle
}

func (ele *Electron) P4() fmom.P4 {
	return &ele.Mom
}

func (ele *Electron) Charge() int32 {
	return ele.EleCharge
}

type Muon struct {
	Mom      fmom.PtEtaPhiM // muon momentum (mass=0.0)
	MuCharge int32          // muon charge

	McPart *hepmc.Particle // generated particle
}

func (muon *Muon) P4() fmom.P4 {
	return &muon.Mom
}

func (muon *Muon) Charge() int32 {
	return muon.MuCharge
}

type Jet struct {
	Mom       fmom.PtEtaPhiM // jet momentum
	JetCharge int32          // jet charge

	DEta float64 // jet radius in pseudo-rapidity
	DPhi float64 // jet radius in azimuthal angle

	BTag   byte // 0 or 1 for a jet that has been tagged as containing a heavy quark
	TauTag byte // 0 or 1 for a jet that has been tagged as a tau

	Constituents []Particle        // pointers to constituents
	McParts      []*hepmc.Particle // pointers to generated particles
}

func (jet *Jet) P4() fmom.P4 {
	return &jet.Mom
}

func (jet *Jet) Charge() int32 {
	return jet.JetCharge
}

type Track struct {
	Pid       int32          // HEP ID number
	Mom       fmom.PtEtaPhiM // track momentum (mass=0.0)
	TrkCharge int32          // track charge

	Eta float64 // track pseudo-rapidity at the tracker edge
	Phi float64 // track azimuthal angle at the tracker edge

	X float64 // track vertex position
	Y float64 // track vertex position
	Z float64 // track vertex position

	Xout float64 // track vertex position at the tracker edge
	Yout float64 // track vertex position at the tracker edge
	Zout float64 // track vertex position at the tracker edge

	McPart *hepmc.Particle // pointer to generated particle
}

func (trk *Track) P4() fmom.P4 {
	return &trk.Mom
}

func (trk *Track) Charge() int32 {
	return trk.TrkCharge
}

type Tower struct {
	Mom  fmom.EtEtaPhiM // calorimeter tower momentum
	Ene  float64        // calorimeter tower energy
	Eem  float64        // calorimeter tower electromagnetic energy
	Ehad float64        // calorimter tower hadronic energy

	Edges [4]float64 // calorimeter tower edges

	McParts []*hepmc.Particle // pointers to generated particles
}

func (tower *Tower) P4() fmom.P4 {
	return &tower.Mom
}

type Candidate struct {
	Pid            int32   // HEP ID number
	Status         int32   // particle status
	M1, M2, D1, D2 int32   // particle mothers and daughters
	CandCharge     int32   // particle charge
	CandMass       float64 // particle mass

	IsPU          byte // 0 or 1 for particles from pile-up interactions
	IsConstituent byte // 0 or 1 for particles being constituents
	BTag          byte // 0 or 1 for a candidate that has been tagged as containing a heavy quark
	TauTag        byte // 0 or 1 for a candidate that has been tagged as a tau

	Eem  float64 // electromagnetic energy
	Ehad float64 // hadronic energy

	Edges [4]float64
	DEta  float64
	DPhi  float64

	Mom  fmom.PxPyPzE
	Pos  fmom.PxPyPzE
	Area fmom.PxPyPzE

	Candidates []Candidate
}

func (cand *Candidate) Clone() *Candidate {
	c := *cand
	c.Candidates = make([]Candidate, 0, len(cand.Candidates))
	for i := range cand.Candidates {
		cc := &cand.Candidates[i]
		c.Add(cc)
	}

	return &c
}

func (cand *Candidate) P4() fmom.P4 {
	return &cand.Mom
}

func (cand *Candidate) Charge() int32 {
	return cand.CandCharge
}

func (cand *Candidate) Add(c *Candidate) {
	cand.Candidates = append(cand.Candidates, *c)
}

func (cand *Candidate) Overlaps(o *Candidate) bool {
	if cand == o {
		return true
	}

	for i := range cand.Candidates {
		cc := &cand.Candidates[i]
		if cc.Overlaps(o) {
			return true
		}
	}

	for i := range o.Candidates {
		cc := &o.Candidates[i]
		if cc.Overlaps(cand) {
			return true
		}
	}

	return false
}

// EOF
