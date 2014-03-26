package fads

import (
	"github.com/go-hep/fwk"
)

type Calorimeter struct {
	fwk.TaskBase

	fracmap map[int64][2]float64
	binmap  map[float64]map[float64]struct{}

	ecalres func(float64) float64
	hcalres func(float64) float64

	tower *Candidate
}

func (proc *Calorimeter) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	return err
}

func (proc *Calorimeter) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	return err
}

func (proc *Calorimeter) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	return err
}

// EOF
