// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"math"
	"math/rand/v2"
	"reflect"
	"sync"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
	"gonum.org/v1/gonum/stat/distuv"
)

type btagclassifier struct {
	PtMin  float64
	EtaMax float64
}

func (btag btagclassifier) Category(parton *Candidate) int {
	if parton.Mom.Pt() <= btag.PtMin || math.Abs(parton.Mom.Eta()) > btag.EtaMax {
		return -1
	}
	pdg := parton.Pid
	if pdg < 0 {
		pdg = -pdg
	}

	if pdg != 21 && pdg > 5 {
		return -1
	}

	return 0
}

type BTagging struct {
	fwk.TaskBase

	partons string
	jets    string
	output  string

	dR  float64
	bit uint

	btag btagclassifier
	eff  map[int]func(pt, eta float64) float64

	seed uint64
	src  *rand.Rand

	flat   distuv.Uniform
	flatmu sync.Mutex
}

func (tsk *BTagging) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.partons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclInPort(tsk.jets, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	tsk.src = rand.New(rand.NewPCG(tsk.seed, tsk.seed))
	tsk.flat = distuv.Uniform{Min: 0, Max: 1, Src: tsk.src}
	return err
}

func (tsk *BTagging) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *BTagging) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *BTagging) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.partons)
	if err != nil {
		return err
	}

	allpartons := v.([]Candidate)

	v, err = store.Get(tsk.jets)
	if err != nil {
		return err
	}
	jets := v.([]Candidate)

	output := make([]Candidate, 0, len(jets))
	defer func() {
		err = store.Put(tsk.output, jets)
	}()

	msg.Debugf("partons: %d\n", len(allpartons))
	msg.Debugf("jets:    %d\n", len(jets))

	partons := make([]Candidate, 0, len(allpartons))
	for i := range allpartons {
		cand := &allpartons[i]
		if tsk.btag.Category(cand) < 0 {
			continue
		}
		partons = append(partons, *cand)
	}

	for i := range jets {
		jet := jets[i].Clone()
		pdgmax := -1
		eta := jet.Mom.Eta()
		pt := jet.Mom.Pt()

		for j := range partons {
			p := &partons[j]
			pdg := int(p.Pid)
			if pdg < 0 {
				pdg = -pdg
			}
			if pdg == 21 {
				pdg = 0
			}
			if fmom.DeltaR(&jet.Mom, &p.Mom) < tsk.dR {
				if pdgmax < pdg {
					pdgmax = pdg
				}
			}
		}

		switch pdgmax {
		case 0:
			pdgmax = 21
		case -1:
			pdgmax = 0
		}

		eff, ok := tsk.eff[pdgmax]
		if !ok {
			eff = tsk.eff[0]
		}

		// apply efficiency
		tag := uint32(0)
		tsk.flatmu.Lock()
		if tsk.flat.Rand() <= eff(pt, eta) {
			tag = 1
		}
		tsk.flatmu.Unlock()
		jet.BTag |= tag << tsk.bit

		output = append(output, *jet)
	}

	msg.Debugf("output:  %d\n", len(output))
	return err
}

func newBTagging(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &BTagging{
		TaskBase: fwk.NewTask(typ, name, mgr),
		partons:  "InputPartons",
		jets:     "InputJets",
		output:   "OutputJets",

		bit: 0,
		dR:  0.5,
		btag: btagclassifier{
			PtMin:  1.0,
			EtaMax: 2.5,
		},
		eff: map[int]func(pt, eta float64) float64{
			0: func(pt, eta float64) float64 { return 0 },
		},

		seed: 1234,
	}

	err = tsk.DeclProp("Partons", &tsk.partons)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Jets", &tsk.jets)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("BitNumber", &tsk.bit)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("DeltaR", &tsk.dR)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PartonPtMin", &tsk.btag.PtMin)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PartonEtaMax", &tsk.btag.EtaMax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Eff", &tsk.eff)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Seed", &tsk.seed)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(BTagging{}), newBTagging)
}
