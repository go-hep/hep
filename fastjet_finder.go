package fads

import (
	"reflect"

	"github.com/go-hep/fastjet"
	"github.com/go-hep/fwk"
)

// fastjetFinder finds jets using the fastjet library
type fastjetFinder struct {
	fwk.TaskBase

	input  string
	output string
	rho    string

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

func (tsk *fastjetFinder) Configure(ctx fwk.Context) error {
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

	return err
}

func (tsk *fastjetFinder) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *fastjetFinder) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *fastjetFinder) Process(ctx fwk.Context) error {
	var err error

	return err
}

func newFastJetFinder(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &fastjetFinder{
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
	fwk.Register(reflect.TypeOf(fastjetFinder{}), newFastJetFinder)
}
