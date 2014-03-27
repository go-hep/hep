package main

import (
	"fmt"

	"github.com/go-hep/fwk"
)

func main() {
	fmt.Printf("::: fads-app...\n")
	app := fwk.NewApp()
	mgr := app.(fwk.TaskMgr)

	c, err := fwk.New("main.task1", "t1")
	if err != nil {
		panic(err)
	}
	mgr.AddTask(c.(fwk.Task))

	c, err = fwk.New("main.task2", "t2")
	if err != nil {
		panic(err)
	}
	mgr.AddTask(c.(fwk.Task))

	c, err = fwk.New("main.task3", "reader")
	if err != nil {
		panic(err)
	}
	mgr.AddTask(c.(fwk.Task))

	c, err = fwk.New("github.com/go-hep/fads.ParticlePropagator", "pprop")
	if err != nil {
		panic(err)
	}
	mgr.AddTask(c.(fwk.Task))

	err = app.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("::: fads-app... [done]\n")
}
