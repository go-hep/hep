// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"math"
	"reflect"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
)

type isoclassifier struct {
	PtMin float64
}

func (iso isoclassifier) Category(track *Candidate) int {
	if track.Mom.Pt() < iso.PtMin {
		return -1
	}
	return 0
}

type Isolation struct {
	fwk.TaskBase

	candidates string
	isolations string
	rhos       string
	output     string

	deltaRMax  float64
	ptRatioMax float64
	ptSumMax   float64
	usePtSum   bool

	classifier isoclassifier
}

func (tsk *Isolation) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.candidates, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclInPort(tsk.isolations, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	if tsk.rhos != "" {
		err = tsk.DeclInPort(tsk.rhos, reflect.TypeOf([]Candidate{}))
		if err != nil {
			return err
		}
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *Isolation) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Isolation) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Isolation) Process(ctx fwk.Context) error {
	var err error

	deltaRMax2 := tsk.deltaRMax * tsk.deltaRMax

	store := ctx.Store()

	v, err := store.Get(tsk.candidates)
	if err != nil {
		return err
	}

	candidates := v.([]Candidate)

	v, err = store.Get(tsk.isolations)
	if err != nil {
		return err
	}

	isolations := v.([]Candidate)

	var rhos []Candidate = nil
	if tsk.rhos != "" {
		v, err = store.Get(tsk.rhos)
		if err != nil {
			return err
		}
		rhos = v.([]Candidate)
	}

	output := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	input := make([]Candidate, 0, len(isolations))
	for i := range isolations {
		cand := &isolations[i]
		if tsk.classifier.Category(cand) < 0 {
			continue
		}
		input = append(input, *cand)
	}

	if len(input) <= 0 {
		return err
	}

	for i := range candidates {
		cand := &input[i]
		eta := math.Abs(cand.Mom.Eta())
		sum := 0.0

		for j := range input {
			iso := &input[j]
			if fmom.DeltaR(&cand.Mom, &iso.Mom) <= tsk.deltaRMax && !cand.Overlaps(iso) {
				sum += iso.Mom.Pt()
			}
		}

		// find rho
		rho := 0.0
		for j := range rhos {
			obj := &rhos[j]
			if eta >= obj.Edges[0] && eta < obj.Edges[1] {
				rho = obj.Mom.Pt()
			}
		}

		// correct sum for pile-up contamination
		sum = sum - rho*deltaRMax2*math.Pi

		ratio := sum / cand.Mom.Pt()
		if (tsk.usePtSum && sum > tsk.ptSumMax) || ratio > tsk.ptRatioMax {
			continue
		}
		output = append(output, *cand)
	}
	return err
}

func newIsolation(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &Isolation{
		TaskBase:   fwk.NewTask(typ, name, mgr),
		candidates: "InputCandidates",
		isolations: "InputIsolations",
		rhos:       "",
		output:     "OutputIsolations",

		deltaRMax:  0.5,
		ptRatioMax: 0.1,
		ptSumMax:   5.0,
		usePtSum:   false,

		classifier: isoclassifier{
			PtMin: 0.5,
		},
	}

	err = tsk.DeclProp("Candidates", &tsk.candidates)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Isolations", &tsk.isolations)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Rhos", &tsk.rhos)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("DeltaRMax", &tsk.deltaRMax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PtRatioMax", &tsk.ptRatioMax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PtSumMax", &tsk.ptSumMax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("UsePtSum", &tsk.usePtSum)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PtMin", &tsk.classifier.PtMin)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(Isolation{}), newIsolation)
}
