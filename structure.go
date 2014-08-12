package fastjet

// JetStructure allows to retrieve information related to the clustering.
type JetStructure interface {
	Constituents(jet *Jet) ([]Jet, error)
}
