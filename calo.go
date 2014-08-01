package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type EtaPhiBin struct {
	EtaBins []float64
	PhiBins []float64
}

type EneFrac struct {
	ECal float64
	HCal float64
}

type Calorimeter struct {
	fwk.TaskBase

	fracmap map[int]EneFrac
	binmap  map[float64]map[float64]struct{} // std::map<float64, std::set<float64>>

	etaphibins []EtaPhiBin
	ecalres    func(eta, ene float64) float64
	hcalres    func(eta, ene float64) float64

	particles   string
	tracks      string
	towers      string
	photons     string
	eflowtracks string
	eflowtowers string
}

func (tsk *Calorimeter) Configure(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	// err = tsk.DeclInPort(tsk.input, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	// err = tsk.DeclOutPort(tsk.output, reflect.TypeOf(sometype{}))
	// if err != nil {
	//	return err
	// }

	tsk.binmap = make(map[float64]map[float64]struct{}, len(tsk.etaphibins))
	for i := range tsk.etaphibins {
		bin := tsk.etaphibins[i]
		for _, eta := range bin.EtaBins {
			for _, phi := range bin.PhiBins {
				if _, ok := tsk.binmap[eta]; !ok {
					tsk.binmap[eta] = make(map[float64]struct{}, len(bin.PhiBins))
				}
				tsk.binmap[eta][phi] = struct{}{}
			}
		}
	}
	return err
}

func (tsk *Calorimeter) StartTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *Calorimeter) StopTask(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func (tsk *Calorimeter) Process(ctx fwk.Context) fwk.Error {
	var err fwk.Error

	return err
}

func newCalorimeter(typ, name string, mgr fwk.App) (fwk.Component, fwk.Error) {
	var err fwk.Error

	tsk := &Calorimeter{
		TaskBase: fwk.NewTask(typ, name, mgr),

		fracmap:    make(map[int]EneFrac),
		etaphibins: make([]EtaPhiBin, 0, 2),
		ecalres:    func(eta, ene float64) float64 { return 0 },
		hcalres:    func(eta, ene float64) float64 { return 0 },

		particles:   "/fads/particles",
		tracks:      "/fads/tracks",
		towers:      "/fads/towers",
		photons:     "/fads/photons",
		eflowtracks: "/fads/eflowtracks",
		eflowtowers: "/fads/eflowtowers",
	}

	// --

	err = tsk.DeclProp("EtaPhiBins", &tsk.etaphibins)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EnergyFraction", &tsk.fracmap)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ECalResolution", &tsk.ecalres)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("HCalResolution", &tsk.hcalres)
	if err != nil {
		return nil, err
	}

	// --

	err = tsk.DeclProp("Particles", &tsk.particles)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Tracks", &tsk.tracks)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Towers", &tsk.towers)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Photons", &tsk.photons)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EFlowTracks", &tsk.eflowtracks)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EFlowTowers", &tsk.eflowtowers)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(Calorimeter{}), newCalorimeter)
}
