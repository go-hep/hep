// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package fastjet_test

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/gonum/floats"
	"go-hep.org/x/hep/fastjet"
)

func TestJetAlgorithms(t *testing.T) {
	const tol = 1e-6

	for _, test := range []struct {
		input string
		name  string
		def   fastjet.JetDefinition
		ptmin float64
		want  []fastjet.Jet
	}{
		{
			input: "testdata/single-pp-event.dat",
			name:  "antikt_r0.4_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.AntiKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "antikt_r0.7_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.AntiKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "antikt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.AntiKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "kt_r0.4_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.KtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "kt_r0.7_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.KtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "kt_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.KtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "cam_r0.4_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.CambridgeAachenAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "cam_r0.7_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.CambridgeAachenAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "cam_r1.0_escheme_best",
			def: fastjet.NewJetDefinition(
				fastjet.CambridgeAachenAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p-1.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+0.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+1.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p-1.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+0.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+1.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p-1.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+0.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-pp-event.dat",
			name:  "genkt_p+1.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.GenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},

		// e+e- algs
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p-1.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+0.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+1.0_r0.4_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.4, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p-1.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+0.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+1.0_r0.7_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 0.7, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p-1.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, -1,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+0.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, 0,
			),
			ptmin: 0,
		},
		{
			input: "testdata/single-ee-event.dat",
			name:  "eegenkt_p+1.0_r1.0_escheme_best",
			def: fastjet.NewJetDefinitionExtra(
				fastjet.EeGenKtAlgorithm, 1.0, fastjet.EScheme, fastjet.BestStrategy, 1,
			),
			ptmin: 0,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			test := test
			particles, err := loadParticles(test.input)
			if err != nil {
				t.Fatal(err)
			}

			cs, err := fastjet.NewClusterSequence(particles, test.def)
			if err != nil {
				t.Fatalf("error for jet definition: %v", err)
			}

			jets, err := cs.InclusiveJets(test.ptmin)
			if err != nil {
				t.Fatalf("incl-jets error: %v", err)
			}

			sort.Sort(fastjet.ByPt(jets))

			want, err := loadRef("testdata/" + test.name + ".ref")
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
				jet := jets[i]
				rap := jet.Rapidity()
				phi := angle0to2Pi(jet.Phi())
				pt := jet.Pt()

				got := []float64{rap, phi, pt}
				if !floats.EqualApprox(got, ref, tol) {
					t.Errorf("#%d\ngot= %v\nwant=%v", i, got, ref)
				}
			}
		})
	}
}

func loadParticles(name string) ([]fastjet.Jet, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var particles []fastjet.Jet
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		txt := scan.Text()
		toks := make([]string, 0, 4)
		for _, tok := range strings.Split(txt, " ") {
			if tok != "" {
				toks = append(toks, tok)
			}
		}
		px, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return nil, err
		}
		py, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return nil, err
		}
		pz, err := strconv.ParseFloat(toks[2], 64)
		if err != nil {
			return nil, err
		}
		e, err := strconv.ParseFloat(toks[3], 64)
		if err != nil {
			return nil, err
		}
		particles = append(particles, fastjet.NewJet(px, py, pz, e))
	}
	err = scan.Err()
	if err != nil {
		return nil, err
	}

	return particles, nil
}

func loadRef(name string) ([][3]float64, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var refs [][3]float64
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		var i int
		var ref [3]float64
		_, err = fmt.Sscanf(scan.Text(), "%5d %f %f %f", &i, &ref[0], &ref[1], &ref[2])
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

const twoPi = 2 * math.Pi

func angle0to2Pi(v float64) float64 {
	v = math.Mod(v, twoPi)
	if v == 0 {
		return 0
	}
	if v < 0 {
		v += twoPi
	}
	if v == twoPi {
		v = 0
	}
	return v
}
