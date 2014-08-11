package fastjet

type ClusterBuilder interface {
	// InclusiveJets returns all jets (in the sense of
	// the inclusive algorithm) with pt >= ptmin
	InclusiveJets(ptmin float64) ([]Jet, error)

	// ExclusiveJets
}
