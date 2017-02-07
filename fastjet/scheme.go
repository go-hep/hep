package fastjet

import "fmt"

// RecombinationScheme defines the recombination choice for the 4-momenta of
// pseudo-jets during the clustering procedure
type RecombinationScheme int

const (
	EScheme   RecombinationScheme = iota // summing the 4-momenta
	PtScheme                             // pt-weighted recombination of y,phi
	Pt2Scheme                            // pt^2 weighted recombination of y,phi
	EtScheme
	Et2Scheme
	BIPtScheme
	BIPt2Scheme

	ExternalScheme RecombinationScheme = 99
)

func (s RecombinationScheme) String() string {
	switch s {
	case EScheme:
		return "E"
	case PtScheme:
		return "Pt"
	case Pt2Scheme:
		return "Pt2"
	case EtScheme:
		return "Et"
	case Et2Scheme:
		return "Et2"
	case BIPtScheme:
		return "BIPt"
	case BIPt2Scheme:
		return "BIPt2"

	case ExternalScheme:
		return "External"

	default:
		panic(fmt.Errorf("fastjet: invalid RecombinationScheme (%d)", int(s)))
	}
}
