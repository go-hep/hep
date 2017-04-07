// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/fmom"
)

// history holds information about the clustering
type history struct {
	parent1 int // index of first parent of this jet were created
	parent2 int // index of second parent of this jet were created
	child   int // index where the current jet is recombined with another
	jet     int
	dij     float64
	maxdij  float64
}

const (
	invalidIndex     = -3
	inexistentParent = -2
	beamJetIndex     = -1
)

type ClusterSequence struct {
	def      JetDefinition
	alg      JetAlgorithm
	strategy Strategy
	r        float64
	r2       float64
	invR2    float64
	qtot     float64

	jets      []Jet
	history   []history
	structure JetStructure
}

func NewClusterSequence(jets []Jet, def JetDefinition) (*ClusterSequence, error) {
	var err error
	cs := &ClusterSequence{
		def:      def,
		alg:      def.Algorithm(),
		strategy: def.Strategy(),
		r:        def.R(),
		jets:     make([]Jet, len(jets), len(jets)*2),
	}

	cs.r2 = cs.r * cs.r
	cs.invR2 = 1.0 / cs.r2
	cs.structure = ClusterSequenceStructure{cs}

	copy(cs.jets, jets)
	err = cs.init()
	if err != nil {
		return nil, err
	}

	err = cs.run()
	if err != nil {
		return nil, err
	}

	return cs, err
}

func (cs *ClusterSequence) InclusiveJets(ptmin float64) ([]Jet, error) {
	var err error
	dcut := ptmin * ptmin
	jets := make([]Jet, 0)
	i := len(cs.history) - 1 // last jet

	switch cs.alg {
	case KtAlgorithm:
		for ; 0 <= i; i-- {
			// with our specific definition of dij and dib (ie: R appears only in
			// dij) then dij==dib is the same as the jet.Pt2() and we can exploit
			// this in selecting the jets...
			if cs.history[i].maxdij < dcut {
				break
			}
			if hh := cs.history[i]; hh.parent2 == beamJetIndex && hh.dij >= dcut {
				// for beam jets
				jets = append(jets, cs.jets[cs.history[hh.parent1].jet])
			}
		}
	case CambridgeAlgorithm:
		for ; 0 <= i; i-- {
			// inclusive jets are all at the end of clustering sequence in the
			// cambridge algorithm.
			// if we find a non-exclusive jet, exit.
			if cs.history[i].parent2 != beamJetIndex {
				break
			}
			parent1 := cs.history[i].parent1
			jet := &cs.jets[cs.history[parent1].jet]
			if jet.Pt2() >= dcut {
				jets = append(jets, *jet)
			}
		}

	case PluginAlgorithm, EeKtAlgorithm, AntiKtAlgorithm,
		GenKtAlgorithm, EeGenKtAlgorithm, CambridgeForPassiveAlgorithm:
		// for inclusive jets with a plugin algorithm, we make no
		// assumption about anything (relation of dij to momenta,
		// ordering of the dij, etc...)
		for ; 0 <= i; i-- {
			hh := cs.history[i]
			if hh.parent2 != beamJetIndex {
				continue
			}
			parent1 := hh.parent1
			jet := &cs.jets[cs.history[parent1].jet]
			if jet.Pt2() >= dcut {
				jets = append(jets, *jet)
			}
		}
	}
	return jets, err
}

func (cs *ClusterSequence) init() error {
	var err error
	cs.history = make([]history, 0, len(cs.jets)*2)
	cs.qtot = 0

	for i := range cs.jets {
		jet := &cs.jets[i]
		cs.history = append(cs.history,
			history{
				parent1: inexistentParent,
				parent2: inexistentParent,
				child:   invalidIndex,
				jet:     i,
				dij:     0.0,
				maxdij:  0.0,
			},
		)
		// perform any momentum pre-processing needed by the recombination scheme
		err = cs.def.Recombiner().Preprocess(jet)
		if err != nil {
			return err
		}

		jet.hidx = i
		jet.structure = cs.structure

		cs.qtot += jet.E()
	}
	return err
}

func (cs *ClusterSequence) run() error {
	var err error
	// nothing to run when event is empty
	if len(cs.jets) <= 0 {
		return err
	}

	// FIXME
	err = cs.runN3Dumb()
	if err != nil {
		return err
	}

	return err
}

// Constituents retrieves the list of constituents of a given jet
func (cs *ClusterSequence) Constituents(jet *Jet) ([]Jet, error) {
	return cs.addConstituents(jet)
}

func (cs *ClusterSequence) addConstituents(jet *Jet) ([]Jet, error) {
	var err error
	var subjets []Jet

	// find position in cluster history
	i := jet.hidx
	hh := &cs.history[i]
	parent1 := hh.parent1
	if parent1 == inexistentParent {
		// It is an original particle (labelled by its parent having value
		// inexistentParent), therefore add it on to the subjet vector
		// Note: we add the initial particle and not simply 'jet' so that
		//       calling addCconstituents with a subtracted jet containing
		//       only one particle will work.
		subjets = append(subjets, cs.jets[i])
		return subjets, err
	}

	// add parent 1
	sub1, err := cs.addConstituents(&cs.jets[cs.history[parent1].jet])
	if err != nil {
		return subjets, err
	}
	subjets = append(subjets, sub1...)

	// see if parent2 is a real jet, then add its constituents
	parent2 := hh.parent2
	if parent2 == beamJetIndex {
		return subjets, err
	}

	sub2, err := cs.addConstituents(&cs.jets[cs.history[parent2].jet])
	if err != nil {
		return subjets, err
	}
	subjets = append(subjets, sub2...)

	return subjets, err
}

func (cs *ClusterSequence) jetScaleForAlgorithm(jet *Jet) float64 {
	switch cs.alg {

	case KtAlgorithm:
		return jet.Pt2()

	case CambridgeAlgorithm:
		return 1.0

	case AntiKtAlgorithm:
		kt2 := jet.Pt2()
		if kt2 > 1e-300 {
			return 1.0 / kt2
		}
		return 1e300

	case GenKtAlgorithm:
		kt2 := jet.Pt2()
		p := cs.def.ExtraParam()
		if p <= 0 && kt2 < 1e-300 {
			kt2 = 1e-300
		}
		return math.Pow(kt2, p)

	case CambridgeForPassiveAlgorithm:
		kt2 := jet.Pt2()
		lim := cs.def.ExtraParam()
		if kt2 < lim*lim && kt2 != 0 {
			return 1.0 / kt2
		}
		return 1.0

	case EeGenKtAlgorithm:
		kt2 := jet.E()
		p := cs.def.ExtraParam()
		if p <= 0 && kt2 < 1e-300 {
			kt2 = 1e-300
		}
		return math.Pow(kt2, 2*p)

	default:
		panic(fmt.Errorf("fastjet: unrecognised jet algorithm (%v)", cs.alg))
	}
}

func (cs *ClusterSequence) setStructure(j *Jet) {
	j.structure = cs.structure
}

// do_ij_recombination_step
func (cs *ClusterSequence) ijRecombinationStep(i, j int, dij float64) (int, error) {

	k := -1
	// create the new jet by recombining the first two
	ijet := &cs.jets[i]
	jjet := &cs.jets[j]
	kjet, err := cs.def.Recombiner().Recombine(ijet, jjet)
	if err != nil {
		return k, err
	}
	k = len(cs.jets)
	khist := len(cs.history)
	kjet.hidx = khist
	cs.jets = append(cs.jets, kjet)

	ihist := ijet.hidx
	jhist := jjet.hidx

	err = cs.addStepToHistory(khist, imin(ihist, jhist), imax(ihist, jhist), k, dij)
	return k, err
}

func (cs *ClusterSequence) ibRecombinationStep(i int, dib float64) error {
	k := len(cs.history)
	err := cs.addStepToHistory(k, cs.jets[i].hidx, beamJetIndex, invalidIndex, dib)
	return err
}

func (cs *ClusterSequence) addStepToHistory(istep, i1, i2, idx int, dij float64) error {
	var err error

	cs.history = append(cs.history,
		history{
			parent1: i1,
			parent2: i2,
			jet:     idx,
			child:   invalidIndex,
			dij:     dij,
			maxdij:  math.Max(dij, cs.history[len(cs.history)-1].maxdij),
		},
	)
	step := len(cs.history) - 1
	if step != istep {
		panic(fmt.Errorf("fastjet: internal logic error (step number dont match (%d != %d))",
			step, istep,
		))
	}

	cs.history[i1].child = step
	if i2 >= 0 {
		cs.history[i2].child = step
	}

	// get cross-referencing right
	if idx != invalidIndex {
		cs.jets[idx].hidx = step
		cs.setStructure(&cs.jets[idx])
	}

	return err
}

// runs the N3Dumb strategy
func (cs *ClusterSequence) runN3Dumb() error {
	var err error
	njets := len(cs.jets)
	type jetinfo struct {
		jet *Jet
		idx int
	}
	jets := make([]jetinfo, njets)
	indices := make([]int, njets)

	for i := range cs.jets {
		jets[i] = jetinfo{
			jet: &cs.jets[i],
			idx: i,
		}
		indices[i] = i
	}

	for n := njets; n > 0; n-- {
		ii := 0
		jj := -2
		// find smallest beam distance
		ymin := cs.jetScaleForAlgorithm(jets[0].jet)
		for i := 0; i < n; i++ {
			y := cs.jetScaleForAlgorithm(jets[i].jet)
			if y < ymin {
				ymin = y
				ii = i
				jj = -2
			}
		}

		// find smallest distance between pair of jets
		for i := 0; i < n-1; i++ {
			ijet := jets[i].jet
			for j := i + 1; j < n; j++ {
				jjet := jets[j].jet
				jetscale := math.Min(
					cs.jetScaleForAlgorithm(ijet),
					cs.jetScaleForAlgorithm(jjet),
				)
				y := math.MaxFloat64
				switch cs.alg {
				case EeGenKtAlgorithm:
					den := 1 - math.Cos(cs.r)
					if cs.r > math.Pi {
						den = 3 + math.Cos(cs.r)
					}
					if den != 0 {
						y = jetscale * (1 - fmom.CosTheta(&ijet.PxPyPzE, &jjet.PxPyPzE)) / den
					}
				default:
					y = jetscale * Distance(ijet, jjet) * cs.invR2
				}
				if y < ymin {
					ymin = y
					ii = i
					jj = j
				}
			}
		}

		// now recombine
		newn := 2*len(jets) - n
		if jj >= 0 {
			//combine pair
			nn, err := cs.ijRecombinationStep(jets[ii].idx, jets[jj].idx, ymin)
			if err != nil {
				return err
			}
			// internal bookkeeping
			jets[ii] = jetinfo{
				jet: &cs.jets[nn],
				idx: nn,
			}
			// have jj point to jet that was pointed at by n-1
			// since original jj is no longer current
			jets[jj] = jets[n-1]
			indices[ii] = newn
			indices[jj] = indices[n-1]
		} else {
			// combine ii with beam
			err = cs.ibRecombinationStep(jets[ii].idx, ymin)
			if err != nil {
				return err
			}
			// put last jet in place of ii which has disappeared
			jets[ii] = jets[n-1]
			indices[ii] = indices[n-1]
		}
	}
	return err
}
