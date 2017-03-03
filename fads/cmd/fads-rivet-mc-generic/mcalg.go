// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"
	"reflect"

	"go-hep.org/x/hep/fads"
	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/heppdt"
)

// McGeneric is a simple task mimicking Rivet's MC_GENERIC analysis looking at various distributions of final state particles
type McGeneric struct {
	fwk.TaskBase

	mcevt  string  // original HepMC event
	etamax float64 // maximum absolute eta acceptance
	ptmin  float64 // minimum pt acceptance in GeV

	hsvc    fwk.HistSvc // handle to the histogram service
	hstream string      // histogram output stream

	hmult fwk.H1D
	hpt   fwk.H1D
	hene  fwk.H1D
	heta  fwk.H1D
	hrap  fwk.H1D
	hphi  fwk.H1D

	hEtaPlus  fwk.H1D
	hEtaMinus fwk.H1D
	hRapPlus  fwk.H1D
	hRapMinus fwk.H1D

	hetaSumEt fwk.P1D

	hetaPMRatio   fwk.S2D
	hrapPMRatio   fwk.S2D
	hetaChPMRatio fwk.S2D
	hrapChPMRatio fwk.S2D

	hmultCh fwk.H1D
	hptCh   fwk.H1D
	heneCh  fwk.H1D
	hetaCh  fwk.H1D
	hrapCh  fwk.H1D
	hphiCh  fwk.H1D

	hEtaPlusCh  fwk.H1D
	hEtaMinusCh fwk.H1D
	hRapPlusCh  fwk.H1D
	hRapMinusCh fwk.H1D
}

func (tsk *McGeneric) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.mcevt, reflect.TypeOf(hepmc.Event{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *McGeneric) StartTask(ctx fwk.Context) error {
	var err error

	svc, err := ctx.Svc("histsvc")
	if err != nil {
		return err
	}

	tsk.hsvc = svc.(fwk.HistSvc)

	tsk.hetaSumEt, err = tsk.hsvc.BookP1D(tsk.hstream+"/EtaSumEt", 25, 0, 5)
	if err != nil {
		return err
	}

	bookH1D := func(name string, nbins int, xmin, xmax float64) fwk.H1D {
		h, e := tsk.hsvc.BookH1D(tsk.hstream+"/"+name, nbins, xmin, xmax)
		if e != nil && err == nil {
			err = e
		}
		return h
	}

	tsk.hmult = bookH1D("Mult", 100, -0.5, 199.5)
	tsk.hpt = bookH1D("Pt", 300, 0, 30)
	tsk.hene = bookH1D("E", 100, 0, 200)
	tsk.heta = bookH1D("Eta", 50, -5, 5)
	tsk.hrap = bookH1D("Rapidity", 50, -5, 5)
	tsk.hphi = bookH1D("Phi", 50, 0, 2*math.Pi)

	tsk.hmultCh = bookH1D("MultCh ", 100, -0.5, 199.5)
	tsk.hptCh = bookH1D("PtCh ", 300, 0, 30)
	tsk.heneCh = bookH1D("ECh ", 100, 0, 200)
	tsk.hetaCh = bookH1D("EtaCh ", 50, -5, 5)
	tsk.hrapCh = bookH1D("RapidityCh ", 50, -5, 5)
	tsk.hphiCh = bookH1D("PhiCh ", 50, 0, 2*math.Pi)

	// temp H1D
	bookH1D = func(name string, nbins int, xmin, xmax float64) fwk.H1D {
		h, e := tsk.hsvc.BookH1D("/"+name, nbins, xmin, xmax)
		if e != nil && err == nil {
			err = e
		}
		return h
	}

	tsk.hEtaPlus = bookH1D("EtaPlus", 25, 0, 5)
	tsk.hEtaMinus = bookH1D("EtaMinus", 25, 0, 5)
	tsk.hRapPlus = bookH1D("RapPlus", 25, 0, 5)
	tsk.hRapMinus = bookH1D("RapMinus", 25, 0, 5)

	tsk.hEtaPlusCh = bookH1D("EtaPlusCh", 25, 0, 5)
	tsk.hEtaMinusCh = bookH1D("EtaMinusCh", 25, 0, 5)
	tsk.hRapPlusCh = bookH1D("RapPlusCh", 25, 0, 5)
	tsk.hRapMinusCh = bookH1D("RapMinusCh", 25, 0, 5)

	bookS2D := func(name string) fwk.S2D {
		s, e := tsk.hsvc.BookS2D(tsk.hstream + "/" + name)
		if e != nil && err == nil {
			err = e
		}
		return s
	}

	tsk.hetaPMRatio = bookS2D("EtaPMRatio")
	tsk.hrapPMRatio = bookS2D("RapidityPMRatio")
	tsk.hetaChPMRatio = bookS2D("EtaChPMRatio")
	tsk.hrapChPMRatio = bookS2D("RapidityChPMRatio")

	return err
}

func (tsk *McGeneric) StopTask(ctx fwk.Context) error {
	var err error
	for _, h := range []*hbook.H1D{
		tsk.hmult.Hist, tsk.heta.Hist, tsk.hrap.Hist, tsk.hpt.Hist, tsk.hene.Hist, tsk.hphi.Hist,
		tsk.hmultCh.Hist, tsk.hetaCh.Hist, tsk.hrapCh.Hist, tsk.hptCh.Hist, tsk.heneCh.Hist, tsk.hphiCh.Hist,
	} {
		area := h.Integral()
		h.Scale(1 / area)
	}

	for _, v := range []struct {
		num *hbook.H1D
		den *hbook.H1D
		res *hbook.S2D
	}{
		{tsk.hEtaPlus.Hist, tsk.hEtaMinus.Hist, tsk.hetaPMRatio.Scatter},
		{tsk.hRapPlus.Hist, tsk.hRapMinus.Hist, tsk.hrapPMRatio.Scatter},
		{tsk.hEtaPlusCh.Hist, tsk.hEtaMinusCh.Hist, tsk.hetaChPMRatio.Scatter},
		{tsk.hRapPlusCh.Hist, tsk.hRapMinusCh.Hist, tsk.hrapChPMRatio.Scatter},
	} {
		res, err := hbook.DivideH1D(v.num, v.den)
		if err != nil {
			return err
		}
		v.res.Fill(res.Points()...)
	}
	return err
}

func (tsk *McGeneric) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.mcevt)
	if err != nil {
		return err
	}
	evt := v.(hepmc.Event)
	weight := 1.0
	if len(evt.Weights.Slice) >= 1 {
		weight = evt.Weights.Slice[0]
	}

	msg.Debugf(
		"event number: %d, #parts=%d #vtx=%d weight=%v\n",
		evt.EventNumber,
		len(evt.Particles), len(evt.Vertices), weight,
	)

	nfsparts := 0
	ncfsparts := 0
	for _, p := range evt.Particles {
		pdg := heppdt.ParticleByID(heppdt.PID(p.PdgID))
		if pdg == nil {
			continue
		}

		mc := fads.Candidate{
			Pid:        int32(p.PdgID),
			Status:     int32(p.Status),
			M2:         1,
			D2:         1,
			CandCharge: int32(pdg.Charge),
			CandMass:   pdg.Mass,
			Mom:        fmom.PxPyPzE(p.Momentum),
		}
		if vtx := p.ProdVertex; vtx != nil {
			mc.M1 = 1
			mc.Pos = fmom.PxPyPzE(vtx.Position)
		}
		pdgcode := p.PdgID
		if pdgcode < 0 {
			pdgcode = -pdgcode
		}

		if p.Status != 1 {
			continue
		}

		eta := mc.Mom.Eta()
		abseta := math.Abs(eta)
		if abseta >= tsk.etamax {
			continue
		}

		if mc.Mom.Pt() <= tsk.ptmin {
			continue
		}

		width := pdg.Resonance.Width.Value
		if width > 1e-10 {
			continue
		}
		rap := mc.Mom.Rapidity()
		absrap := math.Abs(rap)

		pt := mc.Mom.Pt()
		ene := mc.Mom.E()
		// hep/fmom.P4 returns phi in [-pi,pi) range.
		// convert to [0,2pi) to match Rivet convention.
		phi := angle0to2Pi(mc.Mom.Phi())

		tsk.hsvc.FillP1D(tsk.hetaSumEt.ID, abseta, mc.Mom.Et(), weight)
		nfsparts++
		tsk.hsvc.FillH1D(tsk.heta.ID, eta, weight)
		tsk.hsvc.FillH1D(tsk.hrap.ID, rap, weight)
		tsk.hsvc.FillH1D(tsk.hpt.ID, pt, weight)
		tsk.hsvc.FillH1D(tsk.hene.ID, ene, weight)
		tsk.hsvc.FillH1D(tsk.hphi.ID, phi, weight)

		switch eta > 0 {
		case true:
			tsk.hsvc.FillH1D(tsk.hEtaPlus.ID, abseta, weight)
		case false:
			tsk.hsvc.FillH1D(tsk.hEtaMinus.ID, abseta, weight)
		}

		switch rap > 0 {
		case true:
			tsk.hsvc.FillH1D(tsk.hRapPlus.ID, absrap, weight)
		case false:
			tsk.hsvc.FillH1D(tsk.hRapMinus.ID, absrap, weight)
		}

		if mc.CandCharge == 0 {
			continue
		}

		ncfsparts++
		tsk.hsvc.FillH1D(tsk.hetaCh.ID, eta, weight)
		tsk.hsvc.FillH1D(tsk.hrapCh.ID, rap, weight)
		tsk.hsvc.FillH1D(tsk.hptCh.ID, pt, weight)
		tsk.hsvc.FillH1D(tsk.heneCh.ID, ene, weight)
		tsk.hsvc.FillH1D(tsk.hphiCh.ID, phi, weight)

		switch eta > 0 {
		case true:
			tsk.hsvc.FillH1D(tsk.hEtaPlusCh.ID, abseta, weight)
		case false:
			tsk.hsvc.FillH1D(tsk.hEtaMinusCh.ID, abseta, weight)
		}

		switch rap > 0 {
		case true:
			tsk.hsvc.FillH1D(tsk.hRapPlusCh.ID, absrap, weight)
		case false:
			tsk.hsvc.FillH1D(tsk.hRapMinusCh.ID, absrap, weight)
		}

	}

	msg.Debugf("total multiplicity = %d\n", nfsparts)
	tsk.hsvc.FillH1D(tsk.hmult.ID, float64(nfsparts), weight)

	tsk.hsvc.FillH1D(tsk.hmultCh.ID, float64(ncfsparts), weight)
	msg.Debugf("total charged multiplicity = %d\n", ncfsparts)
	return err
}

func newMcGeneric(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &McGeneric{
		TaskBase: fwk.NewTask(typ, name, mgr),
		mcevt:    "/fads/McEvent",
		etamax:   5,
		ptmin:    0.5, // GeV
		hstream:  "/MC_GENERIC",
	}

	err = tsk.DeclProp("Input", &tsk.mcevt)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EtaMax", &tsk.etamax)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("PtMin", &tsk.ptmin)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Stream", &tsk.hstream)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(McGeneric{}), newMcGeneric)
}

const twoPi = 2 * math.Pi

func angle0to2Pi(v float64) float64 {
	v = math.Mod(v, twoPi)
	if v == 0 {
		return 0
	}
	if v < 0 {
		v += twoPi
	}
	if v == twoPi {
		v = 0
	}
	return v
}
