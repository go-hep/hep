package main

import (
	"fmt"

	"github.com/go-hep/fads"
	"github.com/go-hep/fwk"
)

func main() {
	fmt.Printf("::: fads-app...\n")
	app := fwk.NewApp()
	mgr := app.(fwk.TaskMgr)
	mgr.AddTask(&task1{
		TaskBase: fwk.TaskBase{
			Type: "main.task1",
			Name: "t1",
		},
	})
	mgr.AddTask(&task2{
		TaskBase: fwk.TaskBase{
			Type: "main.task2",
			Name: "t2",
		},
	})

	mgr.AddTask(&task3{
		TaskBase: fwk.TaskBase{
			Type: "main.task3",
			Name: "reader",
		},
	})

	mgr.AddTask(&fads.ParticlePropagator{
		TaskBase: fwk.TaskBase{
			Type: "fads.ParticlePropagator",
			Name: "pprop",
		},
	})
	err := app.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("::: fads-app... [done]\n")
}
