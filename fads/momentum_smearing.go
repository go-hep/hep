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

type MomentumSmearing struct {
	fwk.TaskBase

	input  string
	output string

	smear func(x, y float64) float64
	seed  uint64
	src   *rand.Rand
	srcmu sync.Mutex
}

func (tsk *MomentumSmearing) Configure(ctx fwk.Context) error {
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

func (tsk *MomentumSmearing) StartTask(ctx fwk.Context) error {
	var err error
	tsk.src = rand.New(rand.NewPCG(tsk.seed, tsk.seed))
	return err
}

func (tsk *MomentumSmearing) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *MomentumSmearing) Process(ctx fwk.Context) error {
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
		pt := cand.Mom.Pt()

		// apply smearing
		tsk.srcmu.Lock()
		smearPt := distuv.Normal{Mu: pt, Sigma: tsk.smear(pt, eta) * pt, Src: tsk.src}
		pt = smearPt.Rand()
		tsk.srcmu.Unlock()

		if pt <= 0 {
			continue
		}

		mother := cand
		c := cand.Clone()
		eta = cand.Mom.Eta()
		phi := cand.Mom.Phi()

		pxs := pt * math.Cos(phi)
		pys := pt * math.Sin(phi)
		pzs := pt * math.Sinh(eta)
		es := pt * math.Cosh(eta)
		c.Mom = fmom.NewPxPyPzE(pxs, pys, pzs, es)
		c.Add(mother)

		output = append(output, *c)
	}

	msg.Debugf(">>> smeared: %v\n", len(output))

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(MomentumSmearing{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error
			tsk := &MomentumSmearing{
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
