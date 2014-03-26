package fmom

type EtEtaPhiM [4]float64

func NewEtEtaPhiM(et, eta, phi, m float64) EtEtaPhiM {
	return EtEtaPhiM([4]float64{et, eta, phi, m})
}

func (vec *EtEtaPhiM) Et() float64 {
	return vec[0]
}

func (vec *EtEtaPhiM) Eta() float64 {
	return vec[1]
}

func (vec *EtEtaPhiM) Phi() float64 {
	return vec[2]
}

func (vec *EtEtaPhiM) M() float64 {
	return vec[3]
}
