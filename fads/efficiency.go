// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"reflect"
	"sync"

	"go-hep.org/x/hep/fwk"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Efficiency struct {
	fwk.TaskBase

	input  string
	output string

	eff  func(pt, eta float64) float64
	seed uint64
	dist distuv.Uniform
	dmu  sync.Mutex
}

func (tsk *Efficiency) Configure(ctx fwk.Context) error {
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

func (tsk *Efficiency) StartTask(ctx fwk.Context) error {
	var err error
	src := rand.New(rand.NewSource(tsk.seed))
	tsk.dist = distuv.Uniform{Min: 0, Max: 1, Src: src}
	return err
}

func (tsk *Efficiency) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Efficiency) Process(ctx fwk.Context) error {
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

		// apply efficiency
		tsk.dmu.Lock()
		eff := tsk.dist.Rand()
		tsk.dmu.Unlock()
		max := tsk.eff(pt, eta)
		if eff > max {
			continue
		}

		output = append(output, *cand)
	}

	msg.Debugf(">>> output: %v\n", len(output))

	return err
}

func newEfficiency(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	tsk := &Efficiency{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "InputParticles",
		output:   "OutputParticles",
		eff:      func(x, y float64) float64 { return 1 },
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
	fwk.Register(reflect.TypeOf(Efficiency{}), newEfficiency)
}
