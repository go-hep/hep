package fastjet

import (
	"math"

	"github.com/go-hep/fmom"
)

// // Used to protect against parton-level events where pt can be zero
// // for some partons, giving rapidity=infinity. KtJet fails in those cases.
// const (
// 	MaxRap = 1e5
// )

// UserInfo holds extra user information in a Jet
type UserInfo interface{}

// Jet holds minimal information of use for jet-clustering routines
type Jet struct {
	fmom.PxPyPzE

	UserInfo  UserInfo // holds extra user information for this Jet
	hidx      int      // cluster sequence history index
	structure JetStructure
}

func NewJet(px, py, pz, e float64) Jet {
	return Jet{
		PxPyPzE: fmom.NewPxPyPzE(px, py, pz, e),
		hidx:    -1,
	}
}

func (jet *Jet) Pt2() float64 {
	pt := jet.Pt()
	return pt * pt
}

func (jet *Jet) Rapidity() float64 {
	m := jet.M()
	m2 := math.Max(0, m*m) // effective mass - force non-tachyonic mass
	e := jet.E() + math.Abs(jet.Pz())
	rap := 0.5 * math.Log((jet.Pt2()+m2)/(e*e))
	if jet.Pz() > 0 {
		rap = -rap
	}
	return rap
}

func (jet *Jet) Constituents() []Jet {
	// FIXME
	return nil
}

// Distance returns the squared cylinder (rapidity-phi) distance between 2 jets
func Distance(j1, j2 *Jet) float64 {
	dphi := math.Abs(j1.Phi() - j2.Phi())
	if dphi > math.Pi {
		dphi = 2*math.Phi - dphi
	}
	drap := j1.Rapidity() - j2.Rapidity()
	return dphi*dphi + drap*drap
}
