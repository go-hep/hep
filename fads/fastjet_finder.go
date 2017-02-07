package fads

import (
	"math"
	"reflect"
	"sort"

	"go-hep.org/x/hep/fastjet"
	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
)

// FastJetFinder finds jets using the fastjet library
type FastJetFinder struct {
	fwk.TaskBase

	input  string
	output string
	rho    string

	jetDef           fastjet.JetDefinition
	jetAlg           fastjet.JetAlgorithm
	paramR           float64
	jetPtMin         float64
	coneRadius       float64
	seedThreshold    float64
	coneAreaFraction float64
	maxIters         int
	maxPairSize      int
	iratch           int
	adjacencyCut     int
	overlapThreshold float64

	// fastjet area method ---
	areaDef    interface{}
	areaAlg    int
	computeRho bool

	// ghost based areas ---
	ghostEtaMax float64
	repeat      int
	ghostArea   float64
	gridScatter float64
	ptScatter   float64
	meanGhostPt float64

	// voronoi areas ---
	effectiveRfact float64
	etaRangeMap    map[float64]float64
}

func (tsk *FastJetFinder) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.rho, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	if tsk.jetAlg != fastjet.AntiKtAlgorithm {
		return fwk.Errorf("fastjet-finder: only implemented for AntiKt")
	}

	if tsk.areaAlg != 0 {
		return fwk.Errorf("fastjet-finder: only implemented with *NO* area-definition")
	}

	tsk.jetDef = fastjet.NewJetDefinition(tsk.jetAlg, tsk.paramR, fastjet.EScheme, fastjet.BestStrategy)

	return err
}

func (tsk *FastJetFinder) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *FastJetFinder) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *FastJetFinder) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}
	input := v.([]Candidate)

	output := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	injets := make([]fastjet.Jet, 0, len(input))
	for i := range input {
		cand := &input[i]
		jet := fastjet.NewJet(cand.Mom.Px(), cand.Mom.Py(), cand.Mom.Pz(), cand.Mom.E())
		jet.UserInfo = i
		injets = append(injets, jet)
	}

	// construct jets
	var bldr fastjet.Builder
	if tsk.areaDef != nil {
		// FIXME
		panic("not implemented")
	} else {
		bldr, err = fastjet.NewClusterSequence(injets, tsk.jetDef)
		if err != nil {
			return err
		}
	}

	// compute rho and store it
	if tsk.computeRho {
		// FIXME
		panic("not implemented")
	}

	outjets, err := bldr.InclusiveJets(tsk.jetPtMin)
	if err != nil {
		return err
	}
	sort.Sort(fastjet.ByPt(outjets))

	detaMax := 0.0
	dphiMax := 0.0
	output = make([]Candidate, 0, len(outjets))
	for i := range outjets {
		jet := &outjets[i]
		area := fmom.PxPyPzE{0, 0, 0, 0}
		if tsk.areaDef != nil {
			// FIXME
			panic("not implemented")
			// area = jet.Area()
		}

		cand := Candidate{
			Mom: jet.PxPyPzE,
		}

		time := 0.0
		wtime := 0.0
		csts, err := bldr.Constituents(jet)
		if err != nil {
			return err
		}

		for j := range csts {
			idx := csts[j].UserInfo.(int)
			cst := &input[idx]
			deta := math.Abs(cand.Mom.Eta() - cst.Mom.Eta())
			dphi := math.Abs(fmom.DeltaPhi(&cand.Mom, &cst.Mom))
			if deta > detaMax {
				detaMax = deta
			}
			if dphi > dphiMax {
				dphiMax = dphi
			}

			esqrt := math.Sqrt(cst.Mom.E())
			time += esqrt * cst.Pos.T()
			wtime += esqrt

			cand.Add(cst)
		}

		cand.Pos[3] = time / wtime
		cand.Area = area
		cand.DEta = detaMax
		cand.DPhi = dphiMax

		output = append(output, cand)
	}

	// fmt.Printf("%s: input=%02d outjets=%02d\n", tsk.Name(), len(input), len(output))
	return err
}

func newFastJetFinder(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &FastJetFinder{
		TaskBase: fwk.NewTask(typ, name, mgr),

		input:  "/fads/fastjet/input",
		output: "/fads/fastjet/output",
		rho:    "/fads/fastjet/rho",

		jetAlg:           fastjet.AntiKtAlgorithm,
		paramR:           0.5,
		jetPtMin:         10.0,
		coneRadius:       0.5,
		seedThreshold:    1.0,
		coneAreaFraction: 1.0,
		maxIters:         100,
		maxPairSize:      2,
		iratch:           1,
		adjacencyCut:     2,
		overlapThreshold: 0.75,

		// fastjet area method ---
		areaDef:    nil,
		areaAlg:    0,
		computeRho: false,

		// ghost based areas ---
		ghostEtaMax: 5.0,
		repeat:      1,
		ghostArea:   0.01,
		gridScatter: 1.0,
		ptScatter:   0.1,
		meanGhostPt: 1e-100,

		// voronoi areas ---
		effectiveRfact: 1.0,
		etaRangeMap:    make(map[float64]float64),
	}

	err = tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Rho", &tsk.rho)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("JetAlgorithm", &tsk.jetAlg)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ParameterR", &tsk.paramR)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("JetPtMin", &tsk.jetPtMin)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ConeRadius", &tsk.coneRadius)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("SeedThreshold", &tsk.seedThreshold)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ConeAreaFraction", &tsk.coneAreaFraction)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("MaxIterations", &tsk.maxIters)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("MaxPairSize", &tsk.maxPairSize)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Iratch", &tsk.iratch)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("AdjacencyCut", &tsk.adjacencyCut)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("OverlapThreshold", &tsk.overlapThreshold)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("AreaAlgorithm", &tsk.areaAlg)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ComputeRho", &tsk.computeRho)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("GhostEtaMax", &tsk.ghostEtaMax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Repeat", &tsk.repeat)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("GhostArea", &tsk.ghostArea)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("GridScatter", &tsk.gridScatter)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PtScatter", &tsk.ptScatter)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("MeanGhostPt", &tsk.meanGhostPt)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EffectiveRfact", &tsk.effectiveRfact)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("RhoEtaRange", &tsk.etaRangeMap)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(FastJetFinder{}), newFastJetFinder)
}
