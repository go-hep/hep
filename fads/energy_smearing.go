// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"math"
	"reflect"
	"sync"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type EnergySmearing struct {
	fwk.TaskBase

	input  string
	output string

	smear func(eta, ene float64) float64
	seed  uint64
	src   *rand.Rand
	srcmu sync.Mutex
}

func (tsk *EnergySmearing) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *EnergySmearing) StartTask(ctx fwk.Context) error {
	var err error
	tsk.src = rand.New(rand.NewSource(tsk.seed))
	return err
}

func (tsk *EnergySmearing) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *EnergySmearing) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	msg.Debugf(">>> input: %v\n", len(input))

	output := make([]Candidate, 0, len(input))
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	for i := range input {
		cand := &input[i]
		eta := cand.Pos.Eta()
		phi := cand.Pos.Phi()
		ene := cand.Mom.E()

		// apply smearing
		tsk.srcmu.Lock()
		smearEne := distuv.Normal{Mu: ene, Sigma: tsk.smear(eta, ene), Source: tsk.src}
		ene = smearEne.Rand()
		tsk.srcmu.Unlock()

		if ene <= 0 {
			continue
		}

		mother := cand
		c := cand.Clone()
		eta = cand.Mom.Eta()
		phi = cand.Mom.Phi()
		pt := ene / math.Cosh(eta)

		pxs := pt * math.Cos(phi)
		pys := pt * math.Sin(phi)
		pzs := pt * math.Sinh(eta)

		c.Mom = fmom.NewPxPyPzE(pxs, pys, pzs, ene)
		c.Add(mother)

		output = append(output, *c)
	}

	msg.Debugf(">>> smeared: %v\n", len(output))

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(EnergySmearing{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &EnergySmearing{
				TaskBase: fwk.NewTask(typ, name, mgr),
				input:    "InputParticles",
				output:   "OutputParticles",
				smear:    func(x, y float64) float64 { return 0 },
				seed:     1234,
			}

			err = tsk.DeclProp("Input", &tsk.input)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Output", &tsk.output)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Resolution", &tsk.smear)
			if err != nil {
				return nil, err
			}

			err = tsk.DeclProp("Seed", &tsk.seed)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
