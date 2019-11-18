// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

import (
	"math"

	"go-hep.org/x/hep/fmom"
	"golang.org/x/xerrors"
)

type Recombiner interface {
	Description() string
	Recombine(j1, j2 *Jet) (Jet, error)
	Preprocess(jet *Jet) error
	Scheme() RecombinationScheme
}

type DefaultRecombiner struct {
	scheme RecombinationScheme
}

func NewRecombiner(scheme RecombinationScheme) DefaultRecombiner {
	return DefaultRecombiner{
		scheme: scheme,
	}
}

func (rec DefaultRecombiner) Description() string {
	str := rec.scheme.String()
	return str + " scheme recombination"
}

func (rec DefaultRecombiner) Recombine(j1, j2 *Jet) (Jet, error) {
	w1 := 0.0
	w2 := 0.0

	switch rec.Scheme() {
	case EScheme:
		return NewJet(
			j1.Px()+j2.Px(),
			j1.Py()+j2.Py(),
			j1.Pz()+j2.Pz(),
			j1.E()+j2.E(),
		), nil

	case PtScheme, EtScheme, BIPtScheme:
		w1 = j1.Pt()
		w2 = j2.Pt()

	case Pt2Scheme, Et2Scheme, BIPt2Scheme:
		w1 = j1.Pt2()
		w2 = j2.Pt2()

	default:
		return Jet{}, xerrors.Errorf("fastjet.Recombine: invalid recombination scheme (%v)", rec.Scheme())
	}

	pt := j1.Pt() + j2.Pt()
	if pt != 0.0 {
		y := (w1*j1.Rapidity() + w2*j2.Rapidity()) / (w1 + w2)
		phi1 := j1.Phi()
		phi2 := j2.Phi()
		if phi1-phi2 > math.Pi {
			phi2 += 2 * math.Pi
		}
		if phi1-phi2 < -math.Pi {
			phi2 -= 2 * math.Pi
		}
		phi := (w1*phi1 + w2*phi2) / (w1 + w2)
		return NewJet(
			pt*math.Cos(phi),
			pt*math.Sin(phi),
			pt*math.Sinh(y),
			pt*math.Cosh(y),
		), nil
	}

	return NewJet(0, 0, 0, 0), nil
}

func (rec DefaultRecombiner) Preprocess(jet *Jet) error {

	switch rec.Scheme() {
	case EScheme, BIPtScheme, BIPt2Scheme:
		return nil

	case PtScheme, Pt2Scheme:
		// these schemes (as in the ktjet impl.) need massless
		// initial 4-vectors, with essentially E=|p|
		pz := jet.Pz()
		e := math.Sqrt(jet.Pt2() + pz*pz)
		jet.PxPyPzE = fmom.NewPxPyPzE(jet.Px(), jet.Py(), jet.Pz(), e)
		return nil

	case EtScheme, Et2Scheme:
		// these schemes (as in the ktjet impl.) need massless
		// initial 4-vectors, with essentially E=|p|
		pz := jet.Pz()
		rescale := jet.E() / math.Sqrt(jet.Pt2()+pz*pz)
		jet.PxPyPzE = fmom.NewPxPyPzE(
			rescale*jet.Px(),
			rescale*jet.Py(),
			rescale*jet.Pz(),
			jet.E(),
		)
		return nil

	default:
		return xerrors.Errorf("fastjet.Preprocess: invalid recombination scheme (%v)", rec.Scheme())
	}
}

func (rec DefaultRecombiner) Scheme() RecombinationScheme {
	return rec.scheme
}
