package fastjet

type ClusterSequenceStructure struct {
	cs *ClusterSequence
}

func (css ClusterSequenceStructure) Constituents(jet *Jet) ([]Jet, error) {
	return css.cs.Constituents(jet)
}
