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

	app := job.New(nil)

	app.Create(job.C{
		Type: "github.com/go-hep/fads.ParticlePropagator",
		Name: "pprop",
	})

	app.Create(job.C{
		Type: "github.com/go-hep/fads.HepMcReader",
		Name: "hepmcreader",
		Props: job.P{
			"Input": "testdata/hepmc.data",
		},
	})

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

	/*
		c, err = fwk.New("github.com/go-hep/fads.MomentumSmearing", "charged-hadron-smearer")
		if err != nil {
			panic(err)
		}
		mgr.AddTask(c.(fwk.Task))
	*/
	//mgr.DeclProp(c, "Input", )

	app.Run()

	fmt.Printf("::: fads-app... [done]\n")
}
