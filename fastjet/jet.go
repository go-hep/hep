// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

import (
	"math"

	"go-hep.org/x/hep/fmom"
)

const (
	// Used to protect against parton-level events where pt can be zero
	// for some partons, giving rapidity=infinity. KtJet fails in those cases.
	MaxRap = 1e5
)

// UserInfo holds extra user information in a Jet
type UserInfo interface{}

// Jet holds minimal information of use for jet-clustering routines
type Jet struct {
	fmom.PxPyPzE

	UserInfo  UserInfo // holds extra user information for this Jet
	hidx      int      // cluster sequence history index
	structure JetStructure

	// -- cache

	rap float64
	pt2 float64
	phi float64
}

func NewJet(px, py, pz, e float64) Jet {
	jet := Jet{
		PxPyPzE: fmom.NewPxPyPzE(px, py, pz, e),
		hidx:    -1,
	}
	jet.setupCache()
	return jet
}

func (jet *Jet) setupCache() {
	pt := jet.Pt()
	jet.pt2 = pt * pt

	var rap float64
	if jet.E() == math.Abs(jet.Pz()) && jet.Pt2() == 0 {
		rap = MaxRap + math.Abs(jet.Pz())
		if jet.Pz() < 0 {
			rap = -rap
		}
	} else {
		m := jet.M()
		m2 := math.Max(0, m*m) // effective mass - force non-tachyonic mass
		e := jet.E() + math.Abs(jet.Pz())
		rap = 0.5 * math.Log((jet.Pt2()+m2)/(e*e))
		if jet.Pz() > 0 {
			rap = -rap
		}
	}
	jet.rap = rap
	jet.phi = jet.PxPyPzE.Phi()
}

func (jet *Jet) Pt2() float64 {
	return jet.pt2
}

func (jet *Jet) Phi() float64 {
	return jet.phi
}

func (jet *Jet) Rapidity() float64 {
	return jet.rap
}

// Constituents returns the list of constituents for this jet.
func (jet *Jet) Constituents() []Jet {
	subjets, err := jet.structure.Constituents(jet)
	if err != nil {
		panic(err)
	}
	return subjets
}

// Distance returns the squared cylinder (rapidity-phi) distance between 2 jets
func Distance(j1, j2 *Jet) float64 {
	//dphi := deltaPhi(j1, j2)
	dphi := math.Abs(j1.Phi() - j2.Phi())
	if dphi > math.Pi {
		dphi = 2*math.Pi - dphi
	}
	drap := j1.Rapidity() - j2.Rapidity()
	return dphi*dphi + drap*drap
}

func deltaPhi(j1, j2 *Jet) float64 {
	dphi := math.Abs(j1.Phi() - j2.Phi())
	if dphi > math.Pi {
		dphi = 2*math.Pi - dphi
	}
	return dphi
}

func deltaRap(j1, j2 *Jet) float64 {
	return j1.Rapidity() - j2.Rapidity()
}
