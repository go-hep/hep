package fmom

type PtEtaPhiM [4]float64

func NewPtEtaPhiM(pt, eta, phi, m float64) PtEtaPhiM {
	return PtEtaPhiM([4]float64{pt, eta, phi, m})
}

func (vec *PtEtaPhiM) Pt() float64 {
	return vec[0]
}

func (vec *PtEtaPhiM) Eta() float64 {
	return vec[1]
}

func (vec *PtEtaPhiM) Phi() float64 {
	return vec[2]
}

func (vec *PtEtaPhiM) M() float64 {
	return vec[3]
}
