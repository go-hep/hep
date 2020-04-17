// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fads

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/heppdt"
)

type HepMcStreamer struct {
	Name string // input filename
	r    io.ReadCloser
	dec  *hepmc.Decoder

	mcevt string // hepmc event key
}

func (s *HepMcStreamer) Connect(ports []fwk.Port) error {
	var err error
	s.r, err = os.Open(s.Name)
	if err != nil {
		return err
	}

	s.dec = hepmc.NewDecoder(bufio.NewReader(s.r))

	port := ports[0]
	if port.Type != reflect.TypeOf(hepmc.Event{}) {
		err = fmt.Errorf("fads: invalid port. expected type=hepmc.Event. got=%v", port.Type)
		return err
	}

	s.mcevt = port.Name

	return err
}

func (s *HepMcStreamer) Read(ctx fwk.Context) error {
	var err error

	var evt hepmc.Event
	err = s.dec.Decode(&evt)
	if err != nil {
		return err
	}

	store := ctx.Store()
	err = store.Put(s.mcevt, evt)
	return err
}

func (s *HepMcStreamer) Disconnect() error {
	return s.r.Close()
}

type HepMcReader struct {
	fwk.TaskBase

	mcevt       string // original HepMC event
	allparts    string // all particles
	stableparts string // all stable particles
	partons     string // all partons
}

func (tsk *HepMcReader) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.mcevt, reflect.TypeOf(hepmc.Event{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.allparts, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.stableparts, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.partons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	return err
}

func (tsk *HepMcReader) StartTask(ctx fwk.Context) error {
	var err error
	return err
}

func (tsk *HepMcReader) StopTask(ctx fwk.Context) error {
	var err error
	return err
}

func (tsk *HepMcReader) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.mcevt)
	if err != nil {
		return err
	}
	evt := v.(hepmc.Event)

	msg.Debugf(
		"event number: %d, #parts=%d #vtx=%d\n",
		evt.EventNumber,
		len(evt.Particles), len(evt.Vertices),
	)

	allparts := make([]Candidate, 0, len(evt.Particles)/2)
	stableparts := make([]Candidate, 0)
	partons := make([]Candidate, 0)

	for _, p := range evt.Particles {
		allparts = append(allparts, Candidate{
			Pid:        int32(p.PdgID),
			Status:     int32(p.Status),
			M2:         1,
			D2:         1,
			CandCharge: -999,
			CandMass:   -999.9,
			Mom:        fmom.PxPyPzE(p.Momentum),
		})
		c := &allparts[len(allparts)-1]
		pdg := heppdt.ParticleByID(heppdt.PID(p.PdgID))
		if pdg != nil {
			c.CandCharge = int32(pdg.Charge)
			c.CandMass = pdg.Mass
		}

		// FIXME(sbinet)
		if vtx := p.ProdVertex; vtx != nil {
			c.M1 = 1
			c.Pos = fmom.PxPyPzE(vtx.Position)
		}

		if pdg == nil {
			continue
		}

		pdgcode := p.PdgID
		if pdgcode < 0 {
			pdgcode = -pdgcode
		}

		switch {
		case p.Status == 1: // && pdg.IsStable():
			width := pdg.Resonance.Width.Value
			if width <= 1e-10 {
				stableparts = append(stableparts, *c)
			}

		case pdgcode <= 5 || pdgcode == 21 || pdgcode == 15:
			partons = append(partons, *c)
		}

	}

	sort.Sort(ByPt(allparts))
	err = store.Put(tsk.allparts, allparts)
	if err != nil {
		return err
	}

	sort.Sort(ByPt(stableparts))
	err = store.Put(tsk.stableparts, stableparts)
	if err != nil {
		return err
	}

	sort.Sort(ByPt(partons))
	err = store.Put(tsk.partons, partons)
	if err != nil {
		return err
	}

	msg.Debugf("allparts: %d\nstables: %d\npartons: %d\n",
		len(allparts),
		len(stableparts),
		len(partons),
	)

	return err
}

func init() {
	fwk.Register(reflect.TypeOf(HepMcReader{}),
		func(typ, name string, mgr fwk.App) (fwk.Component, error) {
			var err error

			tsk := &HepMcReader{
				TaskBase:    fwk.NewTask(typ, name, mgr),
				mcevt:       "/fads/McEvent",
				allparts:    "/fads/AllParticles",
				stableparts: "/fads/StableParticles",
				partons:     "/fads/Partons",
			}

			err = tsk.DeclProp("Input", &tsk.mcevt)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
