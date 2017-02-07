package fads

import (
	"reflect"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
)

type EnergyScale struct {
	fwk.TaskBase

	input  string
	output string

	scale func(pt, eta float64) float64
}

func (tsk *EnergyScale) Configure(ctx fwk.Context) error {
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

func (tsk *EnergyScale) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *EnergyScale) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *EnergyScale) Process(ctx fwk.Context) error {
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
		cand := input[i].Clone()

		eta := cand.Mom.Eta()
		pt := cand.Mom.Pt()

		// get new scale
		scale := tsk.scale(pt, eta)
		if scale > 0 {
			mom := fmom.Scale(scale, &cand.Mom)
			cand.Mom.Set(mom)
		}

		output = append(output, *cand)
	}

	msg.Debugf(">>> scaled: %v\n", len(output))

	return err
}

func newEnergyScale(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &EnergyScale{
		TaskBase: fwk.NewTask(typ, name, mgr),
		input:    "InputParticles",
		output:   "OutputParticles",
		scale:    func(pt, eta float64) float64 { return 0.0 },
	}

	err = tsk.DeclProp("Input", &tsk.input)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Output", &tsk.output)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Scale", &tsk.scale)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(EnergyScale{}), newEnergyScale)
}
