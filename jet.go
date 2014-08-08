package fastjet

import (
	"github.com/go-hep/fmom"
)

// UserInfo holds extra user information in a Jet
type UserInfo interface{}

// Jet holds minimal information of use for jet-clustering routines
type Jet struct {
	fmom.PxPyPzE

	UserInfo UserInfo // holds extra user information for this Jet
}
