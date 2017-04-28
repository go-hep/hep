// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package fastjet_test

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/gonum/floats"
	"go-hep.org/x/hep/fastjet"
)

func TestClusterSequenceArea(t *testing.T) {
	const tol = 1e-6

	for _, test := range []struct {
		input string
		name  string
		def   fastjet.JetDefinition
		area  fastjet.AreaDefinition
		ptmin float64
	}{
		{
			input: "testdata/single-pp-event.dat",
			name:  "area_ghost_active_kt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.KtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			area:  fastjet.AreaDefinition{}, // ghost-area, active-area
			ptmin: 5.0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "area_ghost_passive_kt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.KtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			area:  fastjet.AreaDefinition{}, // ghost-area, passive-area
			ptmin: 5.0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "area_ghost_active_antikt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.AntiKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			area:  fastjet.AreaDefinition{}, // ghost-area, active-area
			ptmin: 5.0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "area_ghost_passive_antikt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.AntiKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			area:  fastjet.AreaDefinition{}, // ghost-area, passive-area
			ptmin: 5.0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			// TODO
			if strings.Contains(test.name, "passive") {
				t.Skipf("passive area: not implemented")
			}
			// TODO
			if strings.Contains(test.name, "active") {
				t.Skipf("active area: not implemented")
			}
			test := test
			particles, err := loadParticles(test.input)
			if err != nil {
				t.Fatal(err)
			}

			csa, err := fastjet.NewClusterSequenceArea(particles, test.def, test.area)
			if err != nil {
				t.Fatalf("error for jet definition: %v", err)
			}

			jets, err := csa.InclusiveJets(test.ptmin)
			if err != nil {
				t.Fatalf("incl-jets error: %v", err)
			}

			sort.Sort(fastjet.ByPt(jets))

			want, err := loadRefAreas("testdata/" + test.name + ".ref")
			if err != nil {
				t.Fatalf("error reading reference file: %v", err)
			}

			if len(want) != len(jets) {
				t.Fatalf("got %d jets, want %d", len(jets), len(want))
			}

			n := len(jets)
			if len(want) < n {
				n = len(want)
			}
			for i := 0; i < n; i++ {
				ref := want[i][:]
				jet := &jets[i]
				rap := jet.Rapidity()
				phi := angle0to2Pi(jet.Phi())
				pt := jet.Pt()

				area := csa.Area(jet)
				areaErr := csa.AreaErr(jet)

				got := []float64{rap, phi, pt, area, areaErr}
				if !floats.EqualApprox(got, ref, tol) {
					t.Errorf("#%d\ngot= %v\nwant=%v", i, got, ref)
				}
			}
		})
	}

}

func loadRefAreas(name string) ([][5]float64, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var refs [][5]float64
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var i int
		var ref [5]float64
		_, err = fmt.Sscanf(scan.Text(), "%5d %f %f %f %f +- %f", &i, &ref[0], &ref[1], &ref[2], &ref[3], &ref[4])
		if err != nil {
			return nil, err
		}
		refs = append(refs, ref)
	}
	err = scan.Err()
	if err != nil {
		return nil, err
	}
	return refs, nil
}
