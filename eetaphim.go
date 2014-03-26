package fmom

type EEtaPhiM [4]float64

func NewEEtaPhiM(e, eta, phi, m float64) EEtaPhiM {
	return EEtaPhiM([4]float64{e, eta, phi, m})
}

func (vec *EEtaPhiM) E() float64 {
	return vec[0]
}

func (vec *EEtaPhiM) Eta() float64 {
	return vec[1]
}

func (vec *EEtaPhiM) Phi() float64 {
	return vec[2]
}

func (vec *EEtaPhiM) M() float64 {
	return vec[3]
}
