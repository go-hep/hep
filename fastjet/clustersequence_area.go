// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

type ClusterSequenceArea struct {
	cs   *ClusterSequence
	area AreaDefinition
}

func NewClusterSequenceArea(jets []Jet, def JetDefinition, area AreaDefinition) (*ClusterSequenceArea, error) {
	cs, err := NewClusterSequence(jets, def)
	if err != nil {
		return nil, err
	}

	csa := ClusterSequenceArea{
		cs:   cs,
		area: area,
	}
	return &csa, nil
}

func (csa *ClusterSequenceArea) Area(jet *Jet) float64 {
	panic("not implemented")
}

func (csa *ClusterSequenceArea) AreaErr(jet *Jet) float64 {
	panic("not implemented")
}

func (csa *ClusterSequenceArea) NumExclusiveJets(dcut float64) int {
	panic("not implemented")
}

func (cs *ClusterSequenceArea) ExclusiveJets(dcut float64) ([]Jet, error) {
	panic("not implemented")
}

func (cs *ClusterSequenceArea) ExclusiveJetsUpTo(njets int) ([]Jet, error) {
	panic("not implemented")
}

func (csa *ClusterSequenceArea) InclusiveJets(ptmin float64) ([]Jet, error) {
	panic("not implemented")
}
