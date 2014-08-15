package fastjet

// Builder builds jets out of 4-vectors
type Builder interface {
	// InclusiveJets returns all jets (in the sense of
	// the inclusive algorithm) with pt >= ptmin
	InclusiveJets(ptmin float64) ([]Jet, error)

	// ExclusiveJets

	// Constituents retrieves the constituents of a jet
	Constituents(jet *Jet) ([]Jet, error)
}
