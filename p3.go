package fmom

type Vec3 [3]float64

func (vec *Vec3) X() float64 {
	return vec[0]
}

func (vec *Vec3) Y() float64 {
	return vec[1]
}

func (vec *Vec3) Z() float64 {
	return vec[2]
}

// EOF
