package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"

	"github.com/go-hep/fads"
	"github.com/go-hep/fwk/job"
)

var (
	g_lvl      = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
	g_evtmax   = flag.Int("evtmax", -1, "number of events to process")
	g_cpu_prof = flag.Bool("cpu-prof", false, "enable CPU profiling")

	abs  = math.Abs
	sqrt = math.Sqrt
	pow  = math.Pow
)

func main() {
	fmt.Printf("::: fads-app...\n")

	flag.Parse()

	if *g_cpu_prof {
		f, err := os.Create("cpu.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	app := job.New(nil, job.P{
		"EvtMax":   int64(*g_evtmax),
		"MsgLevel": job.MsgLevel(*g_lvl),
	})

	// propagate particles in cylinder
	app.Create(job.C{
		Type: "github.com/go-hep/fads.ParticlePropagator",
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

	// read HepMC data
	app.Create(job.C{
		Type: "github.com/go-hep/fads.HepMcReader",
		Name: "hepmcreader",
		Props: job.P{
			"Input": "testdata/hepmc.data",
			//"Input": "testdata/full.hepmc.data",
		},
	})

	// charged hadron tracking efficiency
	app.Create(job.C{
		Type: "github.com/go-hep/fads.efficiency",
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
		Type: "github.com/go-hep/fads.efficiency",
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
		Type: "github.com/go-hep/fads.efficiency",
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
		Type: "github.com/go-hep/fads.MomentumSmearing",
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
		Type: "github.com/go-hep/fads.energySmearing",
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
		Type: "github.com/go-hep/fads.MomentumSmearing",
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
		Type: "github.com/go-hep/fads.Merger",
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
		Type: "github.com/go-hep/fads.calorimeter",
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
		Type: "github.com/go-hep/fads.Merger",
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

	app.Run()

	fmt.Printf("::: fads-app... [done]\n")
}
