package fads

import (
	"io"
	"os"
	"reflect"

	"github.com/go-hep/fwk"
	"github.com/go-hep/hepmc"
)

type HepMcReader struct {
	fwk.TaskBase

	fname string // input filename
	r     io.ReadCloser
	dec   *hepmc.Decoder // hepmc decoder

	allparts    string // all particles
	stableparts string // all stable particles
	hadrons     string // all hadrons
}

func (tsk *HepMcReader) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	tsk.fname = "hepmc.data"
	tsk.fname = "/home/binet/dev/hepsw/atlasdelphes/data/test.hepmc"
	err = tsk.DeclProp("Input", &tsk.fname)
	if err != nil {
		return err
	}

	tsk.allparts = "/fads/AllParticles"
	tsk.stableparts = "/fads/StableParticles"
	tsk.hadrons = "/fads/Hadrons"

	err = tsk.DeclOutPort(tsk.allparts, reflect.TypeOf(hepmc.Event{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.stableparts, reflect.TypeOf(hepmc.Event{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.hadrons, reflect.TypeOf(hepmc.Event{}))
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
	tsk.dec = hepmc.NewDecoder(tsk.r)
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

	err = store.Put(tsk.allparts, &evt)
	if err != nil {
		return err
	}

	msg.Infof(
		"event number: %d, #parts=%d #vtx=%d\n",
		evt.EventNumber,
		len(evt.Particles), len(evt.Vertices),
	)
	return err
}

func init() {
	fwk.Register(reflect.TypeOf(HepMcReader{}))
}
