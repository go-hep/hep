package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type ParticlePropagator struct {
	fwk.TaskBase

	radius  float64
	radius2 float64
	halflen float64
	bz      float64

	input  string
	output string

	hadrons string
	eles    string
	muons   string
}

func (tsk *ParticlePropagator) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	tsk.radius = 1.0
	err = tsk.DeclProp("Radius", &tsk.radius)
	if err != nil {
		return err
	}
	tsk.radius2 = tsk.radius * tsk.radius

	tsk.halflen = 3.0
	err = tsk.DeclProp("HalfLength", &tsk.halflen)
	if err != nil {
		return err
	}

	tsk.bz = 0.0
	err = tsk.DeclProp("Bz", &tsk.bz)
	if err != nil {
		return err
	}

	if tsk.radius < 1.0e-2 {
		return fwk.Errorf("")
	}

	if tsk.halflen < 1.0e-2 {
		return fwk.Errorf("")
	}

	tsk.input = "/fads/StableParticles"
	err = tsk.DeclProp("InputArray", &tsk.input)
	if err != nil {
		return err
	}

	tsk.output = "StableParticles"
	err = tsk.DeclProp("OutputArray", &tsk.output)
	if err != nil {
		return err
	}

	tsk.hadrons = "ChargedHadrons"
	err = tsk.DeclProp("ChargedHadrons", &tsk.hadrons)
	if err != nil {
		return err
	}

	tsk.eles = "Electrons"
	err = tsk.DeclProp("Electrons", &tsk.eles)
	if err != nil {
		return err
	}

	tsk.muons = "Muons"
	err = tsk.DeclProp("Muons", &tsk.muons)
	if err != nil {
		return err
	}

	err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.hadrons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.eles, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.muons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *ParticlePropagator) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *ParticlePropagator) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *ParticlePropagator) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	msg.Infof(">>> candidates: %v\n", len(input))

	output := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	hadrons := make([]Candidate, 0)
	eles := make([]Candidate, 0)
	muons := make([]Candidate, 0)
	defer func() {
		err = store.Put(tsk.hadrons, hadrons)
		if err != nil {
			return
		}
		err = store.Put(tsk.eles, eles)
		if err != nil {
			return
		}
		err = store.Put(tsk.muons, muons)
		if err != nil {
			return
		}
	}()

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(ParticlePropagator{}))
}

// EOF
