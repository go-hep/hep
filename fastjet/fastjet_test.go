// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet_test

import (
	"fmt"
	"sort"
	"testing"

	"go-hep.org/x/hep/fastjet"
	"go-hep.org/x/hep/fmom"
)

func TestSimple(t *testing.T) {
	particles := []fastjet.Jet{
		fastjet.NewJet(+99.0, +0.1, 0, 100.0),
		fastjet.NewJet(+04.0, -0.1, 0, 005.0),
		fastjet.NewJet(-99.0, +0.0, 0, 099.0),
		fastjet.NewJet(+99.0, +0.1, 0, 199.0),
		fastjet.NewJet(-99.0, +0.0, 0, 299.0),
		fastjet.NewJet(-99.0, +1.0, 0, 399.0),
		fastjet.NewJet(+50.0, +1.0, 100, 399.0),
	}

	// choose a jet definition
	r := 0.7
	def := fastjet.NewJetDefinition(fastjet.AntiKtAlgorithm, r, fastjet.EScheme, fastjet.N3DumbStrategy)

	if def.R() != r {
		t.Fatalf("expected r-param=%v. got=%v", r, def.R())
	}

	if def.RecombinationScheme() != fastjet.EScheme {
		t.Fatalf("got scheme=%v. want=%v", def.RecombinationScheme(), fastjet.EScheme)
	}

	if def.Strategy() != fastjet.N3DumbStrategy {
		t.Fatalf("got strategy=%v. want=%v", def.Strategy(), fastjet.N3DumbStrategy)
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

func TestStrategies(t *testing.T) {
	type result struct {
		jet fmom.PxPyPzE
		cts []fmom.PxPyPzE
	}

	for _, test := range []struct {
		name       string
		particles  []fastjet.Jet
		r          float64
		scheme     fastjet.RecombinationScheme
		ptmin      float64
		strategies []fastjet.Strategy
		want       []result
	}{
		{
			name: "simple-1",
			particles: []fastjet.Jet{
				fastjet.NewJet(+99.0, +0.1, 0, 100.0),
				fastjet.NewJet(+04.0, -0.1, 0, 005.0),
				fastjet.NewJet(-99.0, +0.0, 0, 099.0),
				fastjet.NewJet(+99.0, +0.1, 0, 199.0),
				fastjet.NewJet(-99.0, +0.0, 0, 299.0),
				fastjet.NewJet(-99.0, +1.0, 0, 399.0),
				fastjet.NewJet(+50.0, +1.0, 100, 399.0),
			},
			r:          0.7,
			scheme:     fastjet.EScheme,
			ptmin:      0,
			strategies: []fastjet.Strategy{fastjet.BestStrategy, fastjet.N3DumbStrategy /*FIXME(sbinet): add fastjet.NlnNStrategy*/},
			want: []result{
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
			},
		},
		{
			name: "simple-2",
			particles: []fastjet.Jet{
				fastjet.NewJet(+99.0, +0.1, 0, 100.0),
				fastjet.NewJet(+04.0, -0.1, 0, 005.0),
				fastjet.NewJet(-99.0, +0.0, 0, 099.0),
				fastjet.NewJet(+99.0, +0.2, 0, 199.0),
				fastjet.NewJet(-99.0, +2.0, 0, 299.0),
				fastjet.NewJet(-99.0, +1.0, 0, 399.0),
				fastjet.NewJet(+50.0, +1.0, 100, 399.0),
			},
			r:          0.7,
			scheme:     fastjet.EScheme,
			ptmin:      0,
			strategies: []fastjet.Strategy{fastjet.BestStrategy, fastjet.N3DumbStrategy, fastjet.NlnNStrategy},
			want: []result{
				{
					jet: fmom.NewPxPyPzE(-2.970000e+02, +3.000000e+00, +0.000000e+00, +7.970000e+02),
					cts: []fmom.PxPyPzE{
						fmom.NewPxPyPzE(-9.900000e+01, +0.000000e+00, +0.000000e+00, +0.990000e+02),
						fmom.NewPxPyPzE(-9.900000e+01, +2.000000e+00, +0.000000e+00, +2.990000e+02),
						fmom.NewPxPyPzE(-9.900000e+01, +1.000000e+00, +0.000000e+00, +3.990000e+02),
					},
				},
				{
					jet: fmom.NewPxPyPzE(+2.520000e+02, +1.200000e+00, +1.000000e+02, +7.030000e+02),
					cts: []fmom.PxPyPzE{
						fmom.NewPxPyPzE(+5.000000e+01, +1.000000e+00, +1.000000e+02, +3.990000e+02),
						fmom.NewPxPyPzE(+4.000000e+00, -1.000000e-01, +0.000000e+00, +5.000000e+00),
						fmom.NewPxPyPzE(+9.900000e+01, +1.000000e-01, +0.000000e+00, +1.000000e+02),
						fmom.NewPxPyPzE(+9.900000e+01, +2.000000e-01, +0.000000e+00, +1.990000e+02),
					},
				},
			},
		},
	} {
		test := test
		for _, strategy := range test.strategies {
			t.Run(fmt.Sprintf("%s-%v", test.name, strategy), func(t *testing.T) {

				// choose a jet definition
				def := fastjet.NewJetDefinition(fastjet.AntiKtAlgorithm, test.r, test.scheme, strategy)

				if def.R() != test.r {
					t.Fatalf("got r-param=%v. want=%v", def.R(), test.r)
				}

				if def.RecombinationScheme() != test.scheme {
					t.Fatalf("got scheme=%v. want=%v", def.RecombinationScheme(), test.scheme)
				}

				if def.Strategy() != strategy {
					t.Fatalf("got strategy=%v. want=%v", def.Strategy(), strategy)
				}

				// run the clustering, extract jets
				cs, err := fastjet.NewClusterSequence(test.particles, def)
				if err != nil {
					t.Fatalf("clustering failed: %v", err)
				}

				jets, err := cs.InclusiveJets(test.ptmin)
				if err != nil {
					t.Fatalf("could not retrieve inclusive jets: %v", err)
				}
				sort.Sort(fastjet.ByPt(jets))

				if len(jets) != len(test.want) {
					t.Fatalf("got %d jets. want=%d", len(jets), len(test.want))
				}
				// print the jets
				for i := range jets {
					ref := &test.want[i]
					jet := &jets[i]
					if !fmom.Equal(&ref.jet, jet) {
						t.Fatalf("jet[%d] differ:\ngot:  %v\nwant: %v\n", i,
							jet.PxPyPzE,
							ref.jet,
						)
					}
					constituents := jet.Constituents()
					if len(constituents) != len(ref.cts) {
						t.Fatalf("jet[%d]: got %d constituents. want=%d",
							i, len(constituents), len(ref.cts),
						)
					}
					for j := range constituents {
						jj := &constituents[j]
						if !fmom.Equal(&ref.cts[j], jj) {
							t.Fatalf("jet[%d].constituent[%d] differ:\ngot:  %v\nwant: %v\n",
								i, j,
								jj.PxPyPzE,
								ref.cts[j],
							)
						}
					}
				}
			})
		}
	}
}
