// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hifact // import "go-hep.org/x/hep/hifact"

import (
	"math"

	"gonum.org/v1/gonum/floats"
)

type Model interface {
	Alphas(pars []float64) []float64
	Expected(pars []float64) []float64
	PDF(a, alpha float64) float64
}

type Data struct {
	Binning []float64 `json:"binning"`
	BinData BinData   `json:"bindata"`
}

type BinData struct {
	Data       []float64 `json:"data"`
	Bkg        []float64 `json:"bkg"`
	BkgErr     []float64 `json:"bkgerr"`
	Sig        []float64 `json:"sig"`
	BkgSysUp   []float64 `json:"bkgsys_up"`
	BkgSysDown []float64 `json:"bkgsys_dn"`
}

type Sample struct {
	Data   []float64
	Models []ModelDescr
}

type ModelDescr struct {
	Name string
	Type ModelType
	Data map[string][]float64
}

type ModelType string

const (
	NormFactorModel ModelType = "normfactor"
	ShapeSysModel   ModelType = "shapesys"
	NormSysModel    ModelType = "normsys"
	HistoSysModel   ModelType = "histosys"
)

type ModelConfig struct {
	nuisance map[string]nuisance
	next     int
	poi      int
}

func newConfig() ModelConfig {
	return ModelConfig{
		nuisance: make(map[string]nuisance),
	}
}

func (mc *ModelConfig) AddModel(name string, npars int, m Model) {
	mc.nuisance[name] = nuisance{
		Beg:   mc.next,
		End:   mc.next + npars,
		Model: m,
	}
	mc.next += npars
}

func (mc *ModelConfig) Model(name string) Model {
	return mc.nuisance[name].Model
}

func (mc *ModelConfig) Slice(name string) (int, int) {
	m := mc.nuisance[name]
	return m.Beg, m.End
}

type nuisance struct {
	Beg   int
	End   int
	Model Model
}

type PDF struct {
	Samples  []Sample
	Config   ModelConfig
	AuxData  []float64
	AuxNames []string
}

func New(samples []Sample) *PDF {
	pdf := PDF{
		Samples: samples,
		Config:  newConfig(),
	}

	for _, sample := range samples {
		for _, descr := range sample.Models {
			switch descr.Type {
			case NormFactorModel:
				var m Model = nil // no object for factors
				pdf.Config.AddModel(descr.Name, 1, m)
			case ShapeSysModel:
				// we reserve one parameter for each bin
				m := NewShapeSys(sample.Data, descr.Data["???"])
				pdf.Config.AddModel(descr.Name, len(sample.Data), m)
				// it's a constraint, so this implies more data
				pdf.AuxData = append(pdf.AuxData, m.Aux...)
			case NormSysModel:
				m := NewNormSys(sample.Data, descr.Data["lo"], descr.Data["hi"])
				pdf.Config.AddModel(descr.Name, len(sample.Data), m)
				// it's a constraint, so this implies more data
				pdf.AuxData = append(pdf.AuxData, m.Aux...)
			case HistoSysModel:
				m := NewHistoSys(sample.Data, descr.Data["lo_hist"], descr.Data["hi_hist"])
				pdf.Config.AddModel(descr.Name, len(sample.Data), m)
				// it's a constraint, so this implies more data
				pdf.AuxData = append(pdf.AuxData, m.Aux...)
			default:
				panic("hifact: unknown model [" + descr.Type + "]")
			}
			pdf.AuxNames = append(pdf.AuxNames, descr.Name)
		}
	}
	return &pdf
}

func (pdf *PDF) histoSysDelta(sample Sample, pars []float64) []float64 {
	var models []string
	for _, m := range sample.Models {
		if m.Type != HistoSysModel {
			continue
		}
		models = append(models, m.Name)
	}

	deltas := make([]float64, len(sample.Data))
	for _, name := range models {
		m := pdf.Config.Model(name).(*HistoSys)
		beg, end := pdf.Config.Slice(name)
		val := pars[beg:end]
		// interpolate for each bin
		var mdelta []float64
		for i := range m.Zero {
			lo := m.Mone[i]
			nom := m.Zero[i]
			up := m.Pone[i]
			mdelta = append(mdelta, interp0(val[0], lo, nom, up))
		}
		floats.Add(deltas, mdelta)
	}
	return deltas
}

func (pdf *PDF) ExpectedSample(sample Sample, pars []float64) []float64 {
	nom := make([]float64, len(sample.Data))
	copy(nom, sample.Data)
	hdelta := pdf.histoSysDelta(sample, pars)
	interp := nom
	floats.Add(interp, hdelta)

	var factors []float64
	return factors
}

func (pdf *PDF) ExpectedAuxData(pars []float64) []float64 {
	var auxdata []float64
	for _, name := range pdf.AuxNames {
		m := pdf.Config.Model(name)
		beg, end := pdf.Config.Slice(name)
		aux := m.Expected(pars[beg:end])
		auxdata = append(auxdata, aux...)
	}
	return auxdata
}

func (pdf *PDF) ExpectedActualData(pars []float64) []float64 {
	counts := make([][]float64, len(pdf.Samples))
	for i, sample := range pdf.Samples {
		counts[i] = pdf.ExpectedSample(sample, pars)
	}
	out := make([]float64, len(counts))
	for i := range out {
		out[i] = floats.Sum(counts[i])
	}
	return out
}

func (pdf *PDF) ExpectedData(pars []float64) []float64 {
	actual := pdf.ExpectedActualData(pars)
	constraints := pdf.ExpectedAuxData(pars)
	return append(actual, constraints...)
}

type Result struct {
	QMu    float64
	QMuA   float64
	Sigma  float64
	CLsb   float64
	CLb    float64
	CLs    float64
	CLsExp []float64
}

func RunOnePoint(mu float64, data []float64, pdf *PDF, pars []float64, bounds [][]float64) Result {
	var res Result
	//asimovMu := 0.0
	//asimovData := generateAsimovData(asimovMu, data, pdf, pars, bounds)

	return res
}

func generateAsimovData(mu float64, data []float64, pdf *PDF, pars []float64, bounds [][]float64) []float64 {

	nuisance := constrainedBestFit(mu, data, pdf, pars, bounds)
	return pdf.ExpectedData(nuisance)
}

func constrainedBestFit(mu float64, data []float64, pdf *PDF, pars []float64, bounds [][]float64) []float64 {
	panic("not implemented")
}

func interp0(alpha, mone, zero, pone float64) float64 {
	if alpha > 0 {
		return (pone - zero) * alpha
	}
	return (zero - mone) * alpha
}

func interp1(alpha, mone, zero, pone float64) float64 {
	if alpha > 0 {
		return math.Pow(pone/zero, alpha)
	}
	return math.Pow(mone/zero, -alpha)
}
