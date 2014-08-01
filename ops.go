package fmom

import "fmt"

// Add returns the sum p1+p2.
func Add(p1, p2 P4) P4 {
	// FIXME(sbinet):
	// dispatch most efficient/less-lossy addition
	// based on type(dst) (and, optionally, type(src))
	var sum P4
	switch p1 := p1.(type) {

	case *PxPyPzE:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		sum = &p

	case *EEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp EEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *EtEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp EtEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *PtEtaPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp PtEtaPhiM
		pp.Set(&p)
		sum = &pp

	case *IPtCotThPhiM:
		p := NewPxPyPzE(p1.Px()+p2.Px(), p1.Py()+p2.Py(), p1.Pz()+p2.Pz(), p1.E()+p2.E())
		var pp IPtCotThPhiM
		pp.Set(&p)
		sum = &pp

	default:
		panic(fmt.Errorf("fmom: invalid P4 concrete value: %#v", p1))
	}
	return sum
}

// IAdd adds src into dst, and returns dst
func IAdd(dst, src P4) P4 {
	// FIXME(sbinet):
	// dispatch most efficient/less-lossy addition
	// based on type(dst) (and, optionally, type(src))
	var sum P4
	var p4 *PxPyPzE = nil
	switch p1 := dst.(type) {

	case *PxPyPzE:
		p4 = p1
		sum = dst

	case *EEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *EtEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *PtEtaPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	case *IPtCotThPhiM:
		p := NewPxPyPzE(p1.Px(), p1.Py(), p1.Pz(), p1.E())
		p4 = &p
		sum = dst

	default:
		panic(fmt.Errorf("fmom: invalid P4 concrete value: %#v", dst))
	}
	p4[0] += src.Px()
	p4[1] += src.Py()
	p4[2] += src.Pz()
	p4[3] += src.E()
	sum.Set(p4)
	return sum
}
