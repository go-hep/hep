// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet_test

import (
	"sort"
	"testing"

	"go-hep.org/x/hep/fastjet"
	"go-hep.org/x/hep/fmom"
)

func TestSimple(t *testing.T) {
	t.Parallel()

	particles := []fastjet.Jet{
		fastjet.NewJet(+99.0, +0.1, 0, 100.0),
		fastjet.NewJet(+04.0, -0.1, 0, 005.0),
		fastjet.NewJet(-99.0, +0.0, 0, 099.0),
		fastjet.NewJet(+99.0, +0.1, 0, 199.0),
		fastjet.NewJet(-99.0, +0.0, 0, 299.0),
		fastjet.NewJet(-99.0, +1.0, 0, 399.0),
		fastjet.NewJet(+50.0, +1.0, 100, 399.0),
	}

	// for i, jet := range particles {
	// 	fmt.Printf("part[%d]: pt=%+e eta=%+e rap=%+e phi=%+e\n",
	// 		i, jet.Pt(), jet.Eta(), jet.Rapidity(), jet.Phi(),
	// 	)
	// }

	// choose a jet definition
	r := 0.7
	def := fastjet.NewJetDefinition(fastjet.AntiKtAlgorithm, r, fastjet.EScheme, fastjet.BestStrategy)

	if def.R() != r {
		t.Fatalf("expected r-param=%v. got=%v", r, def.R())
	}

	if def.RecombinationScheme() != fastjet.EScheme {
		t.Fatalf("expected scheme=%v. got=%v", def.RecombinationScheme(), fastjet.EScheme)
	}

	if def.Strategy() != fastjet.BestStrategy {
		t.Fatalf("expected strategy=%v. got=%v", def.Strategy(), fastjet.BestStrategy)
	}

	// run the clustering, extract jets
	cs, err := fastjet.NewClusterSequence(particles, def)
	if err != nil {
		t.Fatalf("clustering failed: %v", err)
	}

	expected := []struct {
		jet fmom.PxPyPzE
		cts []fmom.PxPyPzE
	}{
		{
			jet: fmom.NewPxPyPzE(-2.970000e+02, +1.000000e+00, +0.000000e+00, +7.970000e+02),
			cts: []fmom.PxPyPzE{

				fmom.NewPxPyPzE(-9.900000e+01, +1.000000e+00, +0.000000e+00, +3.990000e+02),
				fmom.NewPxPyPzE(-9.900000e+01, +0.000000e+00, +0.000000e+00, +9.900000e+01),
				fmom.NewPxPyPzE(-9.900000e+01, +0.000000e+00, +0.000000e+00, +2.990000e+02),
			},
		},
		{
			jet: fmom.NewPxPyPzE(+2.520000e+02, +1.100000e+00, +1.000000e+02, +7.030000e+02),
			cts: []fmom.PxPyPzE{
				fmom.NewPxPyPzE(+5.000000e+01, +1.000000e+00, +1.000000e+02, +3.990000e+02),
				fmom.NewPxPyPzE(+4.000000e+00, -1.000000e-01, +0.000000e+00, +5.000000e+00),
				fmom.NewPxPyPzE(+9.900000e+01, +1.000000e-01, +0.000000e+00, +1.000000e+02),
				fmom.NewPxPyPzE(+9.900000e+01, +1.000000e-01, +0.000000e+00, +1.990000e+02),
			},
		},
	}

	const ptmin = 0
	jets, err := cs.InclusiveJets(ptmin)
	if err != nil {
		t.Fatalf("could not retrieve inclusive jets: %v", err)
	}
	sort.Sort(fastjet.ByPt(jets))

	if len(jets) != len(expected) {
		t.Fatalf("expected %d jets. got=%d", len(expected), len(jets))
	}
	// print the jets
	for i := range jets {
		ref := &expected[i]
		jet := &jets[i]
		if !fmom.Equal(&ref.jet, jet) {
			t.Fatalf("jet[%d] differ:\nexp: %v\ngot: %v\n", i,
				ref.jet,
				jet.PxPyPzE,
			)
		}
		constituents := jet.Constituents()
		if len(constituents) != len(ref.cts) {
			t.Fatalf("jet[%d]: expected %d constituents. got=%d",
				i, len(ref.cts), len(constituents),
			)
		}
		for j := range constituents {
			jj := &constituents[j]
			if !fmom.Equal(&ref.cts[j], jj) {
				t.Fatalf("jet[%d].constituent[%d] differ:\nexp: %v\ngot: %v\n",
					i, j,
					ref.cts[j],
					jj.PxPyPzE,
				)
			}
		}
	}
}
