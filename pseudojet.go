package fastjet

import (
	"github.com/go-hep/fmom"
)

// UserInfo holds extra user information in a PseudoJet
type UserInfo interface{}

// PseudoJet holds minimal information of use for jet-clustering routines
type PseudoJet struct {
	fmom.PxPyPzE

	UserInfo UserInfo
}
