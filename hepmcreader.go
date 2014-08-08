package fads

import (
	"bufio"
	"io"
	"os"
	"reflect"
	"sort"

	"github.com/go-hep/fmom"
	"github.com/go-hep/fwk"
	"github.com/go-hep/hepmc"
	"github.com/go-hep/heppdt"
)

type HepMcReader struct {
	fwk.TaskBase

	fname string // input filename
	r     io.ReadCloser
	dec   *hepmc.Decoder // hepmc decoder

	mcevt       string // original HepMC event
	allparts    string // all particles
	stableparts string // all stable particles
	partons     string // all partons
}

func (tsk *HepMcReader) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	err = tsk.DeclOutPort(tsk.mcevt, reflect.TypeOf(hepmc.Event{}))
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

func (tsk *HepMcReader) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	tsk.r, err = os.Open(tsk.fname)
	if err != nil {
		return err
	}
	tsk.dec = hepmc.NewDecoder(bufio.NewReader(tsk.r))
	return err
}

func (tsk *HepMcReader) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	err = tsk.r.Close()
	if err != nil {
		return err
	}
	return err
}

func (tsk *HepMcReader) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error
	store := ctx.Store()
	msg := ctx.Msg()

	var evt hepmc.Event
	err = tsk.dec.Decode(&evt)
	if err != nil {
		return err
	}

	err = store.Put(tsk.mcevt, &evt)
	if err != nil {
		return err
	}

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
			Pid:        int32(p.PdgId),
			Status:     int32(p.Status),
			M2:         1,
			D2:         1,
			CandCharge: -999,
			CandMass:   -999.9,
			Mom:        fmom.PxPyPzE(p.Momentum),
		})
		c := &allparts[len(allparts)-1]
		pdg := heppdt.ParticleByID(heppdt.PID(p.PdgId))
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

		pdgcode := p.PdgId
		if pdgcode < 0 {
			pdgcode = -pdgcode
		}

		switch {
		case p.Status == 1 && pdg.IsStable():
			stableparts = append(stableparts, *c)

		case pdgcode <= 5 || pdgcode == 21 || pdgcode == 15:
			partons = append(partons, *c)
		}

	}

	sort.Sort(sort.Reverse(ByPt(allparts)))
	err = store.Put(tsk.allparts, allparts)
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(ByPt(stableparts)))
	err = store.Put(tsk.stableparts, stableparts)
	if err != nil {
		return err
	}

	sort.Sort(sort.Reverse(ByPt(partons)))
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
		func(typ, name string, mgr fwk.App) (fwk.Component, fwk.Error) {
			var err fwk.Error

			tsk := &HepMcReader{
				TaskBase:    fwk.NewTask(typ, name, mgr),
				fname:       "hepmc.data",
				mcevt:       "/fads/McEvent",
				allparts:    "/fads/AllParticles",
				stableparts: "/fads/StableParticles",
				partons:     "/fads/Partons",
			}

			err = tsk.DeclProp("Input", &tsk.fname)
			if err != nil {
				return nil, err
			}

			return tsk, err
		},
	)
}
