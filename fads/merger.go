package fads

import (
	"math"
	"reflect"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
)

type Merger struct {
	fwk.TaskBase

	inputs []string
	output string
	outene string
	outmom string
}

func (tsk *Merger) Configure(ctx fwk.Context) error {
	var err error

	for _, input := range tsk.inputs {
		err = tsk.DeclInPort(input, reflect.TypeOf([]Candidate{}))
		if err != nil {
			return err
		}
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.outene, reflect.TypeOf(Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.outmom, reflect.TypeOf(Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *Merger) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Merger) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Merger) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	msg := ctx.Msg()

	output := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	sumpt := 0.0
	sumene := 0.0
	p4 := fmom.NewPxPyPzE(0, 0, 0, 0)
	var mom fmom.P4 = &p4

	for _, k := range tsk.inputs {
		v, err := store.Get(k)
		if err != nil {
			return err
		}
		input := v.([]Candidate)
		msg.Debugf(">>> input[%s]: %v\n", k, len(input))

		for i := range input {
			cand := &input[i]
			cmom := cand.Mom
			mom = fmom.IAdd(mom, &cmom)
			sumpt += cmom.Pt()
			sumene += cmom.E()

			output = append(output, *cand)
		}
	}

	var cmom Candidate
	cmom.Mom.Set(mom)

	err = store.Put(tsk.outmom, cmom)
	if err != nil {
		return err
	}

	var cene Candidate

	eta := 0.0
	phi := 0.0
	px := sumpt * math.Cos(phi)
	py := sumpt * math.Sin(phi)
	pz := sumpt * math.Sinh(eta)

	cene.Mom = fmom.NewPxPyPzE(px, py, pz, sumene)
	err = store.Put(tsk.outene, cene)
	if err != nil {
		return err
	}

	msg.Debugf(">>> output: %v\n", len(output))
	return err
}

func newMerger(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &Merger{
		TaskBase: fwk.NewTask(typ, name, mgr),
		inputs:   []string{},
		output:   "candidates",
		outene:   "energy",
		outmom:   "momentum",
	}

	err = tsk.DeclProp("Inputs", &tsk.inputs)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("MomentumOutput", &tsk.outmom)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EnergyOutput", &tsk.outene)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(Merger{}), newMerger)
}
