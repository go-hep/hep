package fmom

type PxPyPzE [4]float64

func NewPxPyPzE(px, py, pz, e float64) PxPyPzE {
	return PxPyPzE([4]float64{px, py, pz, e})
}

func (vec *PxPyPzE) Px() float64 {
	return vec[0]
}

func (vec *PxPyPzE) Py() float64 {
	return vec[1]
}

func (vec *PxPyPzE) Pz() float64 {
	return vec[2]
}

func (vec *PxPyPzE) E() float64 {
	return vec[3]
}

// EOF
