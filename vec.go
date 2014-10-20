package hepmc

// FourVector is a simple quadri-vector representation.
type FourVector [4]float64

// Px returns the x-component of the 4-momentum
func (vec *FourVector) Px() float64 {
	return vec[0]
}

// Py returns the y-component of the 4-momentum
func (vec *FourVector) Py() float64 {
	return vec[1]
}

// Pz returns the z-component of the 4-momentum
func (vec *FourVector) Pz() float64 {
	return vec[2]
}

// E returns the energy of the 4-momentum
func (vec *FourVector) E() float64 {
	return vec[3]
}

// X returns the x-component of the 4-momentum
func (vec *FourVector) X() float64 {
	return vec[0]
}

// Y returns the y-component of the 4-momentum
func (vec *FourVector) Y() float64 {
	return vec[1]
}

// Z returns the z-component of the 4-momentum
func (vec *FourVector) Z() float64 {
	return vec[2]
}

// T returns the t-component of the 4-momentum
func (vec *FourVector) T() float64 {
	return vec[3]
}

// ThreeVector is a simple 3d-vector representation.
type ThreeVector [3]float64

// X returns the x-component of the 3d-vector
func (vec *ThreeVector) X() float64 {
	return vec[0]
}

// Y returns the y-component of the 3d-vector
func (vec *ThreeVector) Y() float64 {
	return vec[1]
}

// Z returns the z-component of the 3d-vector
func (vec *ThreeVector) Z() float64 {
	return vec[2]
}

// EOF
