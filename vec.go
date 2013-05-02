package hepmc

type FourVector [4]float64

func (vec *FourVector) Px() float64 {
	return vec[0]
}

func (vec *FourVector) Py() float64 {
	return vec[1]
}

func (vec *FourVector) Pz() float64 {
	return vec[2]
}

func (vec *FourVector) E() float64 {
	return vec[3]
}

func (vec *FourVector) X() float64 {
	return vec[0]
}

func (vec *FourVector) Y() float64 {
	return vec[1]
}

func (vec *FourVector) Z() float64 {
	return vec[2]
}

func (vec *FourVector) T() float64 {
	return vec[3]
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
