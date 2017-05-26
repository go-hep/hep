// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"math"
	"math/rand"
	"reflect"
	"sync"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
	"gonum.org/v1/gonum/stat/distuv"
)

type tauclassifier struct {
	PtMin  float64
	EtaMax float64
}

func (tag tauclassifier) Category(tau *Candidate, particles []Candidate) int {

	pdg := tau.Pid
	if pdg < 0 {
		pdg = -pdg
	}
	if pdg != 15 {
		return -1
	}

	if tau.Mom.Pt() <= tag.PtMin || math.Abs(tau.Mom.Eta()) > tag.EtaMax {
		return -1
	}

	if tau.D1 < 0 {
		return -1
	}

	for i := tau.D1; i <= tau.D2; i++ {
		daughter := &particles[i]
		pdg := daughter.Pid
		if pdg < 0 {
			pdg = -pdg
		}
		switch pdg {
		case 11, 13, 15, 24:
			return -1
		}
	}
	return 0
}

type TauTagging struct {
	fwk.TaskBase

	particles string
	partons   string
	jets      string
	output    string

	dR float64

	tag tauclassifier
	eff map[int]func(pt, eta float64) float64

	seed int64
	src  *rand.Rand

	flatmu sync.Mutex
	flat   distuv.Uniform
}

func (tsk *TauTagging) Configure(ctx fwk.Context) error {
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

	tsk.src = rand.New(rand.NewSource(tsk.seed))
	tsk.flat = distuv.Uniform{Min: 0, Max: 1, Source: tsk.src}
	return err
}

func (tsk *TauTagging) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *TauTagging) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *TauTagging) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.particles)
	if err != nil {
		return err
	}

	particles := v.([]Candidate)

	v, err = store.Get(tsk.partons)
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
		err = store.Put(tsk.output, output)
	}()

	msg.Debugf("particles: %d\n", len(particles))
	msg.Debugf("partons:   %d\n", len(allpartons))
	msg.Debugf("jets:      %d\n", len(jets))

	taus := make([]Candidate, 0, len(allpartons))
	for i := range allpartons {
		cand := &allpartons[i]
		if tsk.tag.Category(cand, particles) < 0 {
			continue
		}
		taus = append(taus, *cand)
	}

	for i := range jets {
		jet := jets[i].Clone()
		pdg := 0
		eta := jet.Mom.Eta()
		pt := jet.Mom.Pt()

		charge := int32(-1)
		tsk.flatmu.Lock()
		if tsk.flat.Rand() > 0.5 {
			charge = 1
		}
		tsk.flatmu.Unlock()

		for j := range taus {
			mc := &taus[j]
			if mc.D1 < 0 {
				continue
			}

			var p4 fmom.PxPyPzE
			for ii := mc.D1; ii < mc.D2; ii++ {
				daughter := &particles[ii]
				pdg := daughter.Pid
				if pdg == -16 || pdg == 16 {
					continue
				}
				fmom.IAdd(&p4, &daughter.Mom)
			}

			if fmom.DeltaR(&jet.Mom, &p4) < tsk.dR {
				pdg = 15
				charge = mc.CandCharge
			}
		}

		eff, ok := tsk.eff[pdg]
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
		jet.TauTag = tag
		jet.CandCharge = charge

		output = append(output, *jet)
	}

	msg.Debugf("output:  %d\n", len(output))
	return err
}

func newTauTagging(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &TauTagging{
		TaskBase:  fwk.NewTask(typ, name, mgr),
		particles: "InputParticles",
		partons:   "InputPartons",
		jets:      "InputJets",
		output:    "OutputJets",

		dR: 0.5,
		tag: tauclassifier{
			PtMin:  1.0,
			EtaMax: 2.5,
		},
		eff: map[int]func(pt, eta float64) float64{
			0: func(pt, eta float64) float64 { return 0 },
		},

		seed: 1234,
	}

	err = tsk.DeclProp("Particles", &tsk.particles)
	if err != nil {
		return nil, err
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

	err = tsk.DeclProp("DeltaR", &tsk.dR)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("TauPtMin", &tsk.tag.PtMin)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("TauEtaMax", &tsk.tag.EtaMax)
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
	fwk.Register(reflect.TypeOf(TauTagging{}), newTauTagging)
}
