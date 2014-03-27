package fmom

// P4 models a Lorentz 4-vector.
type P4 interface {
	Px() float64       // x component of 4-momentum
	Py() float64       // y component of 4-momentum
	Pz() float64       // z component of 4-momentum
	M() float64        // mass
	M2() float64       // mass squared
	P() float64        // momentum magnitude
	P2() float64       // square of momentum magnitude
	Eta() float64      // pseudo-rapidity
	Rapidity() float64 // rapidity
	Phi() float64      // azimuthal angle in [-pi,pi)
	E() float64        // energy of 4-momentum
	Et() float64       // transverse energy defined to be E*sin(Theta)
	Pt() float64       // transverse momentum
	IPt() float64      // inverse of transverse momentum
	CosPhi() float64   // cosine(Phi)
	SinPhi() float64   // sine(Phi)
	CosTh() float64    // cosine(Theta)
	SinTh() float64    // sine(Theta)
	CotTh() float64    // cottan(Theta)
	TanTh() float64    // tan(Theta)

	Set(p4 P4)
}

// EOF
