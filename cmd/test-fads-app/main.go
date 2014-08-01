package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"runtime/pprof"

	"github.com/go-hep/fwk/job"
)

var (
	g_lvl      = flag.String("l", "INFO", "log level (DEBUG|INFO|WARN|ERROR)")
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
		"EvtMax":   int64(-1),
		"MsgLevel": job.MsgLevel(*g_lvl),
	})

	// propagate particles in cylinder
	app.Create(job.C{
		Type: "github.com/go-hep/fads.ParticlePropagator",
		Name: "pprop",
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
		Type: "github.com/go-hep/fads.Efficiency",
		Name: "charged-hadron-trk-eff",
		Props: job.P{
			"Input":  "ChargedHadrons",
			"Output": "EffChargedHadrons",
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
		Type: "github.com/go-hep/fads.Efficiency",
		Name: "electron-trk-eff",
		Props: job.P{
			"Input":  "Electrons",
			"Output": "EffElectrons",
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
		Type: "github.com/go-hep/fads.Efficiency",
		Name: "muon-trk-eff",
		Props: job.P{
			"Input":  "Muons",
			"Output": "EffMuons",
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
			"Input":  "EffChargedHadrons",
			"Output": "SmearChargedHadrons",
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
		Type: "github.com/go-hep/fads.EnergySmearing",
		Name: "electron-ene-smearing",
		Props: job.P{
			"Input":  "EffElectrons",
			"Output": "SmearElectrons",
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
			"Input":  "EffMuons",
			"Output": "SmearMuons",
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

	app.Run()

	fmt.Printf("::: fads-app... [done]\n")
}
