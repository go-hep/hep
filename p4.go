package fmom

type P4 interface {
	Px() float64
	Py() float64
	Pz() float64
	E() float64

	Pt() float64
	Eta() float64
	Phi() float64
}

type ThreeVector [3]float64

func (vec *ThreeVector) X() float64 {
	return vec[0]
}

func (vec *ThreeVector) Y() float64 {
	return vec[1]
}

func (vec *ThreeVector) Z() float64 {
	return vec[2]
}

// EOF
