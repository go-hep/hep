// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fads-app is a command that runs a simple ATLAS-like detector simulation,
// modelled after the C++ Delphes ATLAS data-card.
//
// Example:
//
//  $> fads-app -help
//  Usage: fads-app [options] <hepmc-input-file>
//
//  ex:
//   $ fads-app -l=INFO -evtmax=-1 ./testdata/hepmc.data
//
//  options:
//    -cpu-prof
//      	enable CPU profiling
//    -evtmax int
//      	number of events to process (default -1)
//    -l string
//      	log level (DEBUG|INFO|WARN|ERROR) (default "INFO")
//    -nprocs int
//      	number of concurrent events to process (default -1)
//    -o string
//      	name of output events file (default "data.rio")
//    -trace string
//      	path to file where to store traces
//
//  $> fads-app ./testdata/hepmc.data
//  ::: fads-app...
//  app                  INFO workers done: 1/4
//  app                  INFO workers done: 2/4
//  app                  INFO workers done: 3/4
//  app                  INFO workers done: 4/4
//  app                  INFO cpu: 340.092148ms
//  app                  INFO mem: alloc:          17963 kB
//  app                  INFO mem: tot-alloc:      35590 kB
//  app                  INFO mem: n-mallocs:      55399
//  app                  INFO mem: n-frees:        54777
//  app                  INFO mem: gc-pauses:          2 ms
//  ::: fads-app... [done] (time=343.533436ms)
//
package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"runtime/pprof"
	"runtime/trace"
	"time"

	"go-hep.org/x/hep/fads"
	"go-hep.org/x/hep/fastjet"
	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/fwk/job"
	"go-hep.org/x/hep/fwk/rio"
	"go-hep.org/x/hep/hepmc"
)

var (
	lvl     = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	evtmax  = flag.Int("evtmax", -1, "number of events to process")
	nprocs  = flag.Int("nprocs", -1, "number of concurrent events to process")
	cpuprof = flag.Bool("cpu-prof", false, "enable CPU profiling")
	ptrace  = flag.String("trace", "", "path to file where to store traces")
	output  = flag.String("o", "data.rio", "name of output events file")

	abs  = math.Abs
	sqrt = math.Sqrt
	pow  = math.Pow
	tanh = math.Tanh
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: fads-app [options] <hepmc-input-file>

ex:
 $ fads-app -l=INFO -evtmax=-1 ./testdata/hepmc.data

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	start := time.Now()

	fmt.Printf("::: fads-app...\n")
	if *cpuprof {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *ptrace != "" {
		f, err := os.Create(*ptrace)
		if err != nil {
			log.Fatalf(
				"error creating trace-file [%s]: %v\n",
				*ptrace,
				err,
			)
		}
		defer f.Close()
		err = trace.Start(f)
		if err != nil {
			log.Fatalf(
				"error starting runtime/trace: %v\n",
				err,
			)
		}
		defer trace.Stop()
	}

	app := job.New(job.P{
		"EvtMax":   int64(*evtmax),
		"NProcs":   *nprocs,
		"MsgLevel": job.MsgLevel(*lvl),
	})

	// propagate particles in cylinder
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Propagator",
		Name: "pprop",
		Props: job.P{
			"Input":          "/fads/StableParticles",
			"Output":         "/fads/pprop/StableParticles",
			"ChargedHadrons": "/fads/pprop/ChargedHadrons",
			"Electrons":      "/fads/pprop/Electrons",
			"Muons":          "/fads/pprop/Muons",

			// radius of the magnetic field coverage, in meters
			"Radius": 1.15,
			// half-length of the magnetic field coverage, in meters
			"HalfLength": 3.51,
			// magnetic field
			"Bz": 2.0,
		},
	})

	input := "testdata/hepmc.data"
	//input := "testdata/full.hepmc.data"
	if flag.NArg() > 0 {
		input = flag.Arg(0)
	}

	// read HepMC data
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk.InputStream",
		Name: "hepmc-streamer",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "/fads/McEvent",
					Type: reflect.TypeOf(hepmc.Event{}),
				},
			},
			"Streamer": &fads.HepMcStreamer{
				Name: input,
			},
		},
	})

	// transform HepMC data into fads collection
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.HepMcReader",
		Name: "hepmcreader",
		Props: job.P{
			"Input": "/fads/McEvent",
		},
	})

	// charged hadron tracking efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "charged-hadron-trk-eff",
		Props: job.P{
			"Input":  "/fads/pprop/ChargedHadrons",
			"Output": "/fads/charged-hadron-trk-eff/ChargedHadrons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case (pt <= 0.1):
					return (0.00)
				case (math.Abs(eta) <= 1.5) && (pt > 0.1 && pt <= 1.0):
					return (0.70)
				case (math.Abs(eta) <= 1.5) && (pt > 1.0):
					return (0.95)
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 0.1 && pt <= 1.0):
					return (0.60)
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 1.0):
					return (0.85)
				case (math.Abs(eta) > 2.5):
					return (0.00)
				}

				return 0
			},
		},
	})

	// electron tracking efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "electron-trk-eff",
		Props: job.P{
			"Input":  "/fads/pprop/Electrons",
			"Output": "/fads/electron-trk-eff/Electrons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case pt <= 0.1:
					return 0
				case math.Abs(eta) <= 1.5 && (pt > 0.1 && pt <= 1.0):
					return 0.70
				case math.Abs(eta) <= 1.5 && (pt > 1.0 && pt <= 100):
					return 0.95
				case (math.Abs(eta) <= 1.5) && (pt > 100):
					return 0.99
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 0.1 && pt <= 1.0):
					return 0.50
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 1.0 && pt <= 100):
					return 0.83
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 100):
					return 0.90
				case (math.Abs(eta) > 2.5):
					return 0
				}

				return 0
			},
		},
	})

	// muon tracking efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "muon-trk-eff",
		Props: job.P{
			"Input":  "/fads/pprop/Muons",
			"Output": "/fads/muon-trk-eff/Muons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case (pt <= 0.1):
					return 0.00
				case (math.Abs(eta) <= 1.5) && (pt > 0.1 && pt <= 1.0):
					return (0.75)
				case (math.Abs(eta) <= 1.5) && (pt > 1.0):
					return (0.99)
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 0.1 && pt <= 1.0):
					return (0.70)
				case (math.Abs(eta) > 1.5 && math.Abs(eta) <= 2.5) && (pt > 1.0):
					return (0.98)
				case (math.Abs(eta) > 2.5):
					return (0.00)
				}

				return 0
			},
		},
	})

	// momentum resolution for charged tracks
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.MomentumSmearing",
		Name: "charged-hadron-mom-smearing",
		Props: job.P{
			"Input":  "/fads/charged-hadron-trk-eff/ChargedHadrons",
			"Output": "/fads/charged-hadron-mom-smearing/ChargedHadrons",
			"Resolution": func(eta, pt float64) float64 {
				switch {
				case (abs(eta) <= 1.5) && (pt > 0.1 && pt <= 1.0):
					return 0.02
				case (abs(eta) <= 1.5) && (pt > 1.0 && pt <= 1.0e1):
					return (0.01)
				case (abs(eta) <= 1.5) && (pt > 1.0e1 && pt <= 2.0e2):
					return (0.03)
				case (abs(eta) <= 1.5) && (pt > 2.0e2):
					return (0.05)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 0.1 && pt <= 1.0):
					return (0.03)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 1.0 && pt <= 1.0e1):
					return (0.02)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 1.0e1 && pt <= 2.0e2):
					return (0.04)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 2.0e2):
					return (0.05)
				}
				return 0
			},
		},
	})

	// energy resolution for electrons
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.EnergySmearing",
		Name: "electron-ene-smearing",
		Props: job.P{
			"Input":  "/fads/electron-trk-eff/Electrons",
			"Output": "/fads/electron-ene-smearing/Electrons",
			"Resolution": func(eta, ene float64) float64 {
				switch {
				case (abs(eta) <= 2.5) && (ene > 0.1 && ene <= 2.5e1):
					return (ene * 0.015)
				case (abs(eta) <= 2.5) && (ene > 2.5e1):
					return sqrt(pow(ene*0.005, 2) + ene*pow(0.05, 2) + pow(0.25, 2))
				case (abs(eta) > 2.5 && abs(eta) <= 3.0):
					return sqrt(pow(ene*0.005, 2) + ene*pow(0.05, 2) + pow(0.25, 2))
				case (abs(eta) > 3.0 && abs(eta) <= 5.0):
					return sqrt(pow(ene*0.107, 2) + ene*pow(2.08, 2))
				}

				return 0
			},
		},
	})

	// momentum resolution for muons
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.MomentumSmearing",
		Name: "muon-mom-smearing",
		Props: job.P{
			"Input":  "/fads/muon-trk-eff/Muons",
			"Output": "/fads/muon-mom-smearing/Muons",
			"Resolution": func(eta, pt float64) float64 {
				switch {
				case (abs(eta) <= 1.5) && (pt > 0.1 && pt <= 1.0):
					return (0.03)
				case (abs(eta) <= 1.5) && (pt > 1.0 && pt <= 5.0e1):
					return (0.03)
				case (abs(eta) <= 1.5) && (pt > 5.0e1 && pt <= 1.0e2):
					return (0.04)
				case (abs(eta) <= 1.5) && (pt > 1.0e2):
					return (0.07)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 0.1 && pt <= 1.0):
					return (0.04)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 1.0 && pt <= 5.0e1):
					return (0.04)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 5.0e1 && pt <= 1.0e2):
					return (0.05)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 1.0e2):
					return (0.10)
				}
				return 0
			},
		},
	})

	// track merger
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Merger",
		Name: "track-merger",
		Props: job.P{
			"Inputs": []string{
				"/fads/charged-hadron-mom-smearing/ChargedHadrons",
				"/fads/electron-ene-smearing/Electrons",
				"/fads/muon-mom-smearing/Muons",
			},
			"Output":         "/fads/track-merger/tracks",
			"MomentumOutput": "/fads/track-merger/momentum",
			"EnergyOutput":   "/fads/track-merger/energy",
		},
	})

	phi10 := make([]float64, 0, 37)
	for i := -18; i <= 18; i++ {
		phi10 = append(phi10, float64(i)*math.Pi/18.0)
	}

	phi20 := make([]float64, 0, 19)
	for i := -9; i <= 9; i++ {
		phi20 = append(phi20, float64(i)*math.Pi/9.0)
	}

	// calorimeter
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Calorimeter",
		Name: "calo",
		Props: job.P{
			"Particles":   "/fads/pprop/StableParticles",
			"Tracks":      "/fads/track-merger/tracks",
			"Towers":      "/fads/calo/towers",
			"Photons":     "/fads/calo/photons",
			"EFlowTracks": "/fads/calo/eflowtracks",
			"EFlowTowers": "/fads/calo/eflowtowers",

			"EtaPhiBins": fads.NewEtaPhiGrid(
				[]fads.EtaPhiBin{
					// 10-degrees towers: 0 <= |eta| <= 3.2
					{
						EtaBins: []float64{
							-3.2,
							-2.5, -2.4, -2.3, -2.2, -2.1, -2.0,
							-1.9, -1.8, -1.7, -1.6, -1.5, -1.4, -1.3, -1.2, -1.1, -1.0,
							-0.9, -0.8, -0.7, -0.6, -0.5, -0.4, -0.3, -0.2, -0.1, +0.0,
							+0.1, +0.2, +0.3, +0.4, +0.5, +0.6, +0.7, +0.8, +0.9,
							+1.0, +1.1, +1.2, +1.3, +1.4, +1.5, +1.6, +1.7, +1.8, +1.9,
							+2.0, +2.1, +2.2, +2.3, +2.4, +2.5, +2.6,
							+3.3,
						},
						PhiBins: phi10,
					},

					// 20-degrees towers: 2.8 <= |eta| <= 4.9
					{
						EtaBins: []float64{
							-4.9, -4.7, -4.5, -4.3, -4.1,
							-3.9, -3.7, -3.5, -3.3, -3.0,
							-2.8, -2.6,
							+2.8, +3.0, +3.2, +3.5, +3.7, +3.9,
							+4.1, +4.3, +4.5, +4.7, +4.9,
						},
						PhiBins: phi20,
					},
				},
			),

			// default energy fractions: abs(pid) -> {ECal, HCal}
			"EnergyFraction": map[int]fads.EneFrac{
				0: {
					ECal: 0,
					HCal: 1,
				},
				// energy fractions for ele,gamma and pi0
				11:  {ECal: 1, HCal: 0},
				22:  {ECal: 1, HCal: 0},
				111: {ECal: 1, HCal: 0},
				// energy fractions for muons, neutrinos and neutralinos
				12:      {ECal: 0.0, HCal: 0.0},
				13:      {ECal: 0.0, HCal: 0.0},
				14:      {ECal: 0.0, HCal: 0.0},
				16:      {ECal: 0.0, HCal: 0.0},
				1000022: {ECal: 0.0, HCal: 0.0},
				1000023: {ECal: 0.0, HCal: 0.0},
				1000025: {ECal: 0.0, HCal: 0.0},
				1000035: {ECal: 0.0, HCal: 0.0},
				1000045: {ECal: 0.0, HCal: 0.0},
				// energy fractions for K0short and Lambda
				310:  {ECal: 0.3, HCal: 0.7},
				3122: {ECal: 0.3, HCal: 0.7},
			},

			// ecal resolution (eta and energy)
			// http://arxiv.org/pdf/physics/0608012v1 jinst8_08_s08003
			// http://villaolmo.mib.infn.it/ICATPP9th_2005/Calorimetry/Schram.p.pdf
			// http://www.physics.utoronto.ca/~krieger/procs/ComoProceedings.pdf
			"ECalResolution": func(eta, ene float64) float64 {
				switch {
				case (abs(eta) <= 3.2):
					return sqrt(pow(ene, 2)*pow(0.0017, 2) + ene*pow(0.101, 2))

				case (abs(eta) > 3.2 && abs(eta) <= 4.9):
					return sqrt(pow(ene, 2)*pow(0.0350, 2) + ene*pow(0.285, 2))
				}
				return 0
			},

			// hcal resolution (eta and energy)
			// http://arxiv.org/pdf/hep-ex/0004009v1
			// http://villaolmo.mib.infn.it/ICATPP9th_2005/Calorimetry/Schram.p.pdf
			"HCalResolution": func(eta, ene float64) float64 {
				switch {
				case (abs(eta) <= 1.7):
					return sqrt(pow(ene, 2)*pow(0.0302, 2) + ene*pow(0.5205, 2) + pow(1.59, 2))
				case (abs(eta) > 1.7 && abs(eta) <= 3.2):
					return sqrt(pow(ene, 2)*pow(0.0500, 2) + ene*pow(0.706, 2))
				case (abs(eta) > 3.2 && abs(eta) <= 4.9):
					return sqrt(pow(ene, 2)*pow(0.9420, 2) + ene*pow(0.075, 2))
				}
				return 0
			},
		},
	})

	// eflow merger
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Merger",
		Name: "eflow-merger",
		Props: job.P{
			"Inputs": []string{
				"/fads/calo/eflowtracks",
				"/fads/calo/eflowtowers",
			},
			"Output":         "/fads/eflow-merger/eflow",
			"MomentumOutput": "/fads/eflow-merger/momentum",
			"EnergyOutput":   "/fads/eflow-merger/energy",
		},
	})

	// photon efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "photon-eff",
		Props: job.P{
			"Input":  "/fads/calo/photons",
			"Output": "/fads/photon-eff/photons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case (pt <= 10.0):
					return 0.00
				case (abs(eta) <= 1.5) && (pt > 10.0):
					return 0.95
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 10.0):
					return 0.85
				case (abs(eta) > 2.5):
					return 0.00
				}

				return 0
			},
		},
	})

	// photon isolation
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Isolation",
		Name: "photon-iso",
		Props: job.P{
			"Candidates": "/fads/photon-eff/photons",
			"Isolations": "/fads/eflow-merger/eflow",
			"Output":     "/fads/photon-iso/photons",

			"DeltaRMax":  0.5,
			"PtMin":      0.5,
			"PtRatioMax": 0.1,
		},
	})

	// electron efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "electron-eff",
		Props: job.P{
			"Input":  "/fads/electron-ene-smearing/Electrons",
			"Output": "/fads/electron-eff/electrons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case (pt <= 10.0):
					return (0.00)
				case (abs(eta) <= 1.5) && (pt > 10.0):
					return (0.95)
				case (abs(eta) > 1.5 && abs(eta) <= 2.5) && (pt > 10.0):
					return (0.85)
				case (abs(eta) > 2.5):
					return (0.00)
				}
				return 0
			},
		},
	})

	// electron isolation
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Isolation",
		Name: "electron-iso",
		Props: job.P{
			"Candidates": "/fads/electron-eff/electrons",
			"Isolations": "/fads/eflow-merger/eflow",
			"Output":     "/fads/electron-iso/electrons",

			"DeltaRMax":  0.5,
			"PtMin":      0.5,
			"PtRatioMax": 0.1,
		},
	})

	// muon efficiency
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Efficiency",
		Name: "muon-eff",
		Props: job.P{
			"Input":  "/fads/muon-mom-smearing/Muons",
			"Output": "/fads/muon-eff/muons",
			"Eff": func(eta, pt float64) float64 {
				switch {
				case (pt <= 10.0):
					return (0.00)
				case (abs(eta) <= 1.5) && (pt > 10.0):
					return (0.95)
				case (abs(eta) > 1.5 && abs(eta) <= 2.7) && (pt > 10.0):
					return (0.85)
				case (abs(eta) > 2.7):
					return (0.00)
				}
				return 0
			},
		},
	})

	// muon isolation
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Isolation",
		Name: "muon-iso",
		Props: job.P{
			"Candidates": "/fads/muon-eff/muons",
			"Isolations": "/fads/eflow-merger/eflow",
			"Output":     "/fads/muon-iso/muons",

			"DeltaRMax":  0.5,
			"PtMin":      0.5,
			"PtRatioMax": 0.1,
		},
	})

	// missing-et merger
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Merger",
		Name: "missing-et",
		Props: job.P{
			"Inputs": []string{
				"/fads/calo/eflowtracks",
				"/fads/calo/eflowtowers",
			},
			"Output":         "/fads/missing-et",
			"MomentumOutput": "/fads/missing-et/momentum",
			"EnergyOutput":   "/fads/missing-et/energy",
		},
	})

	// mc truth jet finder
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.FastJetFinder",
		Name: "mc-jet-finder",
		Props: job.P{
			"Input":  "/fads/StableParticles",
			"Output": "/fads/mc-jet-finder/jets",
			"Rho":    "/fads/mc-jet-finder/rho",

			"JetAlgorithm": fastjet.AntiKtAlgorithm,
			"ParameterR":   0.6,

			"JetPtMin": 20.0,
		},
	})

	// jet finder
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.FastJetFinder",
		Name: "fastjet-finder",
		Props: job.P{
			"Input":  "/fads/calo/towers",
			"Output": "/fads/fastjet-finder/jets",
			"Rho":    "/fads/fastjet-finder/rho",

			"JetAlgorithm": fastjet.AntiKtAlgorithm,
			"ParameterR":   0.6,

			"JetPtMin": 20.0,
		},
	})

	// jet energy scale
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.EnergyScale",
		Name: "jet-ene-scale",
		Props: job.P{
			"Input":  "/fads/fastjet-finder/jets",
			"Output": "/fads/jet-ene-scale/jets",
			"Scale":  func(eta, pt float64) float64 { return 1.08 },
		},
	})

	// b-tagging
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.BTagging",
		Name: "btag",
		Props: job.P{
			"Partons": "/fads/Partons",
			"Jets":    "/fads/jet-ene-scale/jets",
			"Output":  "/fads/btag/jets",

			"BitNumber":    uint(0),
			"DeltaR":       0.5,
			"PartonPtMin":  1.0,
			"PartonEtaMax": 2.5,

			// efficiency formula: [pdg-code] -> (pt,eta)
			// pdg-code: the highest PDG code of a quark or gluon inside a DeltaR-cone
			//           around the jet axis
			//           gluon's pdg-code has the lowest priority.
			"Eff": map[int]func(pt, eta float64) float64{
				// default efficiency (mis-identification rate)
				0: func(pt, eta float64) float64 { return 0.001 },

				// efficiency for c-jets (mis-identification rate)
				4: func(pt, eta float64) float64 {
					switch {
					case pt <= 15.0:
						return (0.000)

					case (abs(eta) <= 1.2) && (pt > 15.0):
						return (0.2 * tanh(pt*0.03-0.4))

					case (abs(eta) > 1.2 && abs(eta) <= 2.5) && (pt > 15.0):
						return (0.1 * tanh(pt*0.03-0.4))

					case (abs(eta) > 2.5):
						return (0.000)
					}
					return 0
				},

				// efficiency for b-jets
				5: func(pt, eta float64) float64 {
					switch {
					case (pt <= 15.0):
						return (0.000)
					case (abs(eta) <= 1.2) && (pt > 15.0):
						return (0.5 * tanh(pt*0.03-0.4))
					case (abs(eta) > 1.2 && abs(eta) <= 2.5) && (pt > 15.0):
						return (0.4 * tanh(pt*0.03-0.4))
					case (abs(eta) > 2.5):
						return (0.000)
					}
					return 0
				},
			},
		},
	})

	// tau-tagging
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.TauTagging",
		Name: "tau-tag",
		Props: job.P{
			"Particles": "/fads/AllParticles",
			"Partons":   "/fads/Partons",
			"Jets":      "/fads/btag/jets",
			"Output":    "/fads/tau-tag/jets",

			"DeltaR":    0.5,
			"TauPtMin":  1.0,
			"TauEtaMax": 2.5,

			// efficiency formula: [pdg-code] -> (pt,eta)
			"Eff": map[int]func(pt, eta float64) float64{
				// default efficiency (mis-identification rate)
				0: func(pt, eta float64) float64 { return 0.001 },

				// efficiency for tau-jets
				15: func(pt, eta float64) float64 { return 0.4 },
			},
		},
	})

	// find uniquely identified photons/electrons/taus/jets
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.UniqueObjectFinder",
		Name: "uobj-finder",
		Props: job.P{
			"Keys": []fads.ObjPair{
				{
					In:  "/fads/photon-iso/photons",
					Out: "/fads/uobj-finder/photons",
				},
				{
					In:  "/fads/electron-iso/electrons",
					Out: "/fads/uobj-finder/electrons",
				},
				{
					In:  "/fads/muon-iso/muons",
					Out: "/fads/uobj-finder/muons",
				},
				{
					In:  "/fads/tau-tag/jets",
					Out: "/fads/uobj-finder/jets",
				},
			},
		},
	})

	// scalar HT merger
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fads.Merger",
		Name: "scalar-ht",
		Props: job.P{
			"Inputs": []string{
				"/fads/uobj-finder/jets",
				"/fads/uobj-finder/electrons",
				"/fads/uobj-finder/photons",
				"/fads/uobj-finder/muons",
			},
			"Output":         "/fads/scalar-ht",
			"MomentumOutput": "/fads/scalar-ht/momentum",
			"EnergyOutput":   "/fads/scalar-ht/energy",
		},
	})

	// output
	app.Create(job.C{
		Type: "go-hep.org/x/hep/fwk.OutputStream",
		Name: "rio-output",
		Props: job.P{
			"Ports": []fwk.Port{
				{
					Name: "/fads/McEvent",
					Type: reflect.TypeOf(hepmc.Event{}),
				},
				{
					Name: "/fads/uobj-finder/jets",
					Type: reflect.TypeOf([]fads.Candidate{}),
				},
				{
					Name: "/fads/uobj-finder/electrons",
					Type: reflect.TypeOf([]fads.Candidate{}),
				},
				{
					Name: "/fads/uobj-finder/photons",
					Type: reflect.TypeOf([]fads.Candidate{}),
				},
				{
					Name: "/fads/uobj-finder/muons",
					Type: reflect.TypeOf([]fads.Candidate{}),
				},
			},
			"Streamer": &rio.OutputStreamer{
				Name: *output,
			},
		},
	})

	app.Run()
	fmt.Printf("::: fads-app... [done] (time=%v)\n", time.Since(start))
}
