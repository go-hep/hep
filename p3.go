package fmom

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
