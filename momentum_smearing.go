package fads

import (
	"math"
	"math/rand"
	"reflect"

	"github.com/go-hep/fmom"
	"github.com/go-hep/fwk"
	"github.com/go-hep/random"
)

type MomentumSmearing struct {
	fwk.TaskBase

	input  string
	output string

	smear func(x, y float64) float64
	seed  int64
	src   rand.Source
}

func (tsk *MomentumSmearing) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *MomentumSmearing) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	tsk.src = rand.NewSource(tsk.seed)
	return err
}

func (tsk *MomentumSmearing) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *MomentumSmearing) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.input)
	if err != nil {
		return err
	}

	input := v.([]Candidate)
	msg.Infof(">>> input: %v\n", len(input))

	output := make([]Candidate, 0, len(input))
	defer func() {
		err = store.Put(tsk.output, output)
	}()

	for i := range input {
		cand := &input[i]
		eta := cand.Pos.Eta()
		phi := cand.Pos.Phi()
		pt := cand.Mom.Pt()

		// apply smearing
		smearPt := random.Gauss(pt, tsk.smear(pt, eta)*pt, &tsk.src)
		pt = smearPt()

		if pt <= 0 {
			continue
		}

		mother := cand
		c := cand.Clone()
		eta = cand.Mom.Eta()
		phi = cand.Mom.Phi()

		pxs := pt * math.Cos(phi)
		pys := pt * math.Sin(phi)
		pzs := pt * math.Sinh(eta)
		es := pt * math.Cosh(eta)
		c.Mom = fmom.NewPxPyPzE(pxs, pys, pzs, es)
		c.Add(mother)

		output = append(output, *c)
	}

	msg.Infof(">>> smeared: %v\n", len(output))

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(MomentumSmearing{}),
		func(name string, mgr fwk.App) (fwk.Component, fwk.Error) {
			var err fwk.Error
			tsk := &MomentumSmearing{
				TaskBase: fwk.NewTask(name, mgr),
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

			err = tsk.DeclInPort(tsk.input, reflect.TypeOf([]Candidate{}))
			if err != nil {
				return nil, err
			}

			err = tsk.DeclOutPort(tsk.output, reflect.TypeOf([]Candidate{}))
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
