package fwk

import (
	"fmt"
	"sort"
)

var g_app App = nil

type fsm int

const (
	fsm_ONLINE     fsm = 0
	fsm_CONFIGURED fsm = 1
	fsm_STARTED    fsm = 2
	fsm_RUNNING    fsm = 3
	fsm_STOPPED    fsm = 4
	fsm_OFFLINE    fsm = 5
)

type appmgr struct {
	state    fsm
	comptype string
	compname string

	props map[Component]map[string]interface{}
	dflow *dflowsvc
	store *datastore

	evtmax int64

	tsks []Task
	svcs []Svc
}

func NewApp() App {
	if g_app != nil {
		return g_app
	}

	var app *appmgr

	app = &appmgr{
		state:    fsm_ONLINE,
		comptype: "appmgr",
		compname: "app",
		props:    make(map[Component]map[string]interface{}),
		dflow:    newDflowSvc("dataflow"),
		store:    newDataStore("evtstore"),
		evtmax:   4,
		tsks:     make([]Task, 0),
		svcs:     make([]Svc, 0),
	}

	var err Error
	err = app.AddSvc(app.store)
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create evtstore: %v\n", err)
		return nil
	}

	err = app.AddSvc(app.dflow)
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create dataflow svc: %v\n", err)
		return nil
	}

	g_app = app
	return app
}

func (app *appmgr) CompType() string {
	return app.comptype
}

func (app *appmgr) CompName() string {
	return app.compname
}

func (app *appmgr) AddTask(tsk Task) Error {
	var err Error
	app.tsks = append(app.tsks, tsk)
	return err
}

func (app *appmgr) DelTask(tsk Task) Error {
	var err Error
	tsks := make([]Task, 0, len(app.tsks))
	for _, t := range app.tsks {
		if t.CompName() != tsk.CompName() {
			tsks = append(tsks, t)
		}
	}
	app.tsks = tsks
	return err
}

func (app *appmgr) HasTask(n string) bool {
	for _, t := range app.tsks {
		if t.CompName() == n {
			return true
		}
	}
	return false
}

func (app *appmgr) GetTask(n string) Task {
	for _, t := range app.tsks {
		if t.CompName() == n {
			return t
		}
	}
	return nil
}

func (app *appmgr) Tasks() []Task {
	return app.tsks
}

func (app *appmgr) AddSvc(svc Svc) Error {
	var err Error
	app.svcs = append(app.svcs, svc)
	return err
}

func (app *appmgr) DelSvc(svc Svc) Error {
	var err Error
	svcs := make([]Svc, 0, len(app.svcs))
	for _, s := range app.svcs {
		if s.CompName() != svc.CompName() {
			svcs = append(svcs, s)
		}
	}
	app.svcs = svcs
	return err
}

func (app *appmgr) HasSvc(n string) bool {
	for _, s := range app.svcs {
		if s.CompName() == n {
			return true
		}
	}
	return false
}

func (app *appmgr) GetSvc(n string) Svc {
	for _, s := range app.svcs {
		if s.CompName() == n {
			return s
		}
	}
	return nil
}

func (app *appmgr) Svcs() []Svc {
	return app.svcs
}

func (app *appmgr) Run() Error {
	var err Error
	var ctx Context

	if app.state == fsm_ONLINE {
		err = app.configure(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm_CONFIGURED {
		err = app.start(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm_STARTED {
		err = app.run(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm_RUNNING {
		err = app.stop(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm_STOPPED {
		err = app.shutdown(ctx)
		if err != nil {
			return err
		}
	}

	return err
}

func (app *appmgr) configure(ctx Context) Error {
	var err Error
	fmt.Printf(">>> app.configure...\n")
	for _, svc := range app.svcs {
		fmt.Printf(">>> configuring [%v:%v]...\n", svc.CompType(), svc.CompName())
		cfg, ok := svc.(Configurer)
		if !ok {
			continue
		}
		err = cfg.Configure(ctx)
		if err != nil {
			return err
		}
	}

	for _, tsk := range app.tsks {
		fmt.Printf(">>> configuring [%v:%v]...\n", tsk.CompType(), tsk.CompName())
		cfg, ok := tsk.(Configurer)
		if !ok {
			continue
		}
		err = cfg.Configure(ctx)
		if err != nil {
			return err
		}
	}

	fmt.Printf(">>> --- [data flow] --- nodes...\n")
	for tsk, node := range app.dflow.nodes {
		fmt.Printf(">>> ---[%s]---\n", tsk)
		fmt.Printf("    in:  %v\n", node.in)
		fmt.Printf("    out: %v\n", node.out)
	}

	fmt.Printf(">>> --- [data flow] --- edges...\n")
	edges := make([]string, 0, len(app.dflow.edges))
	for n := range app.dflow.edges {
		edges = append(edges, n)
	}
	sort.Strings(edges)
	fmt.Printf(" edges: %v\n", edges)

	app.state = fsm_CONFIGURED
	fmt.Printf(">>> app.configure... [done]\n")
	return err
}

func (app *appmgr) start(ctx Context) Error {
	var err Error
	for _, svc := range app.svcs {
		fmt.Printf(">>> starting [%v:%v]...\n", svc.CompType(), svc.CompName())
		err = svc.StartSvc(ctx)
		if err != nil {
			return err
		}
	}

	for _, tsk := range app.tsks {
		fmt.Printf(">>> starting [%v:%v]...\n", tsk.CompType(), tsk.CompName())
		err = tsk.StartTask(ctx)
		if err != nil {
			return err
		}
	}

	app.state = fsm_STARTED
	return err
}

func (app *appmgr) run(ctx Context) Error {
	var err Error
	app.state = fsm_RUNNING
	for ievt := int64(0); ievt < app.evtmax; ievt++ {
		fmt.Printf(">>> app.running evt=%d...\n", ievt)
		for k := range app.dflow.edges {
			app.store.store[k] = make(achan, 1)
		}
		errch := make(chan Error, len(app.tsks))
		for _, tsk := range app.tsks {
			go func(tsk Task) {
				//fmt.Printf(">>> running [%v:%v]...\n", tsk.CompType(), tsk.CompName())
				ctx := context{id: ievt, slot: 0, store: app.store}
				errch <- tsk.Process(ctx)
			}(tsk)
		}
		for i := 0; i < len(app.tsks); i++ {
			err := <-errch
			if err != nil {
				close(errch)
				return err
			}

		}
	}
	return err
}

func (app *appmgr) stop(ctx Context) Error {
	var err Error
	for _, tsk := range app.tsks {
		err = tsk.StopTask(ctx)
		if err != nil {
			return err
		}
	}

	for _, svc := range app.svcs {
		err = svc.StopSvc(ctx)
		if err != nil {
			return err
		}
	}

	app.state = fsm_STOPPED
	return err
}

func (app *appmgr) shutdown(ctx Context) Error {
	var err Error
	app.tsks = nil
	app.svcs = nil
	app.state = fsm_OFFLINE

	app.props = nil
	app.dflow = nil
	app.store = nil

	g_app = nil
	return err
}

// EOF
