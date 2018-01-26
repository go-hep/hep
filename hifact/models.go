// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hifact // import "go-hep.org/x/hep/hifact"

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

type NormSys struct {
	Zero []float64
	Mone []float64
	Pone []float64
	Aux  []float64
}

func NewNormSys(data, mone, pone []float64) *NormSys {
	return &NormSys{
		Zero: data,
		Mone: mone,
		Pone: pone,
		Aux:  []float64{0}, // observed data is always at a = 1
	}
}

func (sys *NormSys) Alphas(pars []float64) []float64 {
	// nuisance parameters correspond directly to the alpha
	return pars
}

func (sys *NormSys) Expected(pars []float64) []float64 {
	return sys.Alphas(pars)
}

func (sys *NormSys) PDF(a, alpha float64) float64 {
	return distuv.Normal{Mu: alpha, Sigma: 1}.Prob(a)
}

type HistoSys struct {
	Zero []float64
	Mone []float64
	Pone []float64
	Aux  []float64
}

func NewHistoSys(data, mone, pone []float64) *HistoSys {
	return &HistoSys{
		Zero: data,
		Mone: mone,
		Pone: pone,
		Aux:  []float64{0}, // observed data is always at a = 1
	}
}

func (sys *HistoSys) Alphas(pars []float64) []float64 {
	// nuisance parameters correspond directly to the alpha
	return pars
}

func (sys *HistoSys) Expected(pars []float64) []float64 {
	return sys.Alphas(pars)
}

func (sys *HistoSys) PDF(a, alpha float64) float64 {
	return distuv.Normal{Mu: alpha, Sigma: 1}.Prob(a)
}

type ShapeSys struct {
	Aux []float64
	db2 []float64
}

func NewShapeSys(data, deltab []float64) *ShapeSys {
	if len(data) != len(deltab) {
		panic("hifact: lengths do not match")
	}
	var sys ShapeSys
	sys.Aux = make([]float64, len(data))
	sys.db2 = make([]float64, len(data))
	for i := range data {
		b := data[i]
		delta := deltab[i]
		db2 := b * b / (delta * delta) // tau*b
		sys.Aux[i] = db2
		sys.db2[i] = db2
	}
	return &sys
}

func (sys *ShapeSys) Alphas(pars []float64) []float64 {
	o := make([]float64, len(pars))
	for i := range pars {
		o[i] = pars[i] * sys.db2[i]
	}
	return o
}

func (sys *ShapeSys) Expected(pars []float64) []float64 {
	return sys.Alphas(pars)
}

func (sys *ShapeSys) PDF(a, alpha float64) float64 {
	return distuv.Normal{Mu: alpha, Sigma: math.Sqrt(alpha)}.Prob(a)
}

var (
	_ Model = (*NormSys)(nil)
	_ Model = (*HistoSys)(nil)
	_ Model = (*ShapeSys)(nil)
)
