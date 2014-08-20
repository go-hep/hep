package main

import (
	"fmt"

	"github.com/go-hep/fwk/job"
)

func handle_err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("::: fwk-app...\n")

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
			"Input": "ChargedHadrons",
		},
	})

	app.Run()

	fmt.Printf("::: fwk-app... [done]\n")
}
