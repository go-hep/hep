package fastjet

import (
	"github.com/go-hep/fmom"
)

// PseudoJet holds minimal information of use for jet-clustering routines
type PseudoJet struct {
	fmom.PxPyPzE
}
