package fwk

import (
	"fmt"
	"reflect"
	"sort"
)

var g_mgr *appmgr = nil

type fsm int

const (
	fsm_UNDEFINED  fsm = 0
	fsm_CONFIGURED fsm = 1
	fsm_STARTED    fsm = 2
	fsm_RUNNING    fsm = 3
	fsm_STOPPED    fsm = 4
	fsm_OFFLINE    fsm = 5
)

type appmgr struct {
	state fsm
	name  string

	props map[Component]map[string]interface{}
	dflow *dflowsvc
	store *datastore

	evtmax int64

	comps map[string]Component
	tsks  []Task
	svcs  []Svc
}

func NewApp() App {
	if g_mgr != nil {
		return g_mgr
	}

	var err Error
	var app *appmgr

	app = &appmgr{
		state:  fsm_UNDEFINED,
		name:   "app",
		props:  make(map[Component]map[string]interface{}),
		dflow:  nil,
		store:  nil,
		evtmax: 4,
		comps:  make(map[string]Component),
		tsks:   make([]Task, 0),
		svcs:   make([]Svc, 0),
	}

	svc, err := New("github.com/go-hep/fwk.datastore", "evtstore")
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create evtstore: %v\n", err)
		return nil
	}
	app.store = svc.(*datastore)

	//app.store = newDataStore("evtstore")
	err = app.AddSvc(app.store)
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create evtstore: %v\n", err)
		return nil
	}

	svc, err = New("github.com/go-hep/fwk.dflowsvc", "dataflow")
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create dataflow svc: %v\n", err)
		return nil
	}
	app.dflow = svc.(*dflowsvc)

	//app.dflow = newDflowSvc("dataflow")
	err = app.AddSvc(app.dflow)
	if err != nil {
		fmt.Printf("**error** fwk.NewApp: could not create dataflow svc: %v\n", err)
		return nil
	}

	g_mgr = app
	return app
}

func (app *appmgr) Name() string {
	return app.name
}

func (app *appmgr) SetName(n string) {
	app.name = n
}

func (app *appmgr) Component(n string) Component {
	c, ok := app.comps[n]
	if !ok {
		return nil
	}
	return c
}

func (app *appmgr) HasComponent(n string) bool {
	_, ok := app.comps[n]
	return ok
}

func (app *appmgr) Components() []Component {
	comps := make([]Component, 0, len(app.comps))
	for _, c := range app.comps {
		comps = append(comps, c)
	}
	return comps
}

func (app *appmgr) AddTask(tsk Task) Error {
	var err Error
	app.tsks = append(app.tsks, tsk)
	app.comps[tsk.Name()] = tsk
	return err
}

func (app *appmgr) DelTask(tsk Task) Error {
	var err Error
	tsks := make([]Task, 0, len(app.tsks))
	for _, t := range app.tsks {
		if t.Name() != tsk.Name() {
			tsks = append(tsks, t)
		}
	}
	app.tsks = tsks
	return err
}

func (app *appmgr) HasTask(n string) bool {
	for _, t := range app.tsks {
		if t.Name() == n {
			return true
		}
	}
	return false
}

func (app *appmgr) GetTask(n string) Task {
	for _, t := range app.tsks {
		if t.Name() == n {
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
	app.comps[svc.Name()] = svc
	return err
}

func (app *appmgr) DelSvc(svc Svc) Error {
	var err Error
	svcs := make([]Svc, 0, len(app.svcs))
	for _, s := range app.svcs {
		if s.Name() != svc.Name() {
			svcs = append(svcs, s)
		}
	}
	app.svcs = svcs
	return err
}

func (app *appmgr) HasSvc(n string) bool {
	for _, s := range app.svcs {
		if s.Name() == n {
			return true
		}
	}
	return false
}

func (app *appmgr) GetSvc(n string) Svc {
	for _, s := range app.svcs {
		if s.Name() == n {
			return s
		}
	}
	return nil
}

func (app *appmgr) Svcs() []Svc {
	return app.svcs
}

func (app *appmgr) DeclProp(c Component, name string, ptr interface{}) Error {
	_, ok := app.props[c]
	if !ok {
		app.props[c] = make(map[string]interface{})
	}
	switch reflect.TypeOf(ptr).Kind() {
	case reflect.Ptr:
		// ok
	default:
		return Errorf(
			"fwk.DeclProp: component [%s] didn't pass a pointer for the property [%s] (type=%T)",
			c.Name(),
			name,
			ptr,
		)
	}
	app.props[c][name] = ptr
	return nil
}

func (app *appmgr) SetProp(c Component, name string, value interface{}) Error {
	m, ok := app.props[c]
	if !ok {
		return Errorf(
			"fwk.SetProp: component [%s] didn't declare any property",
			c.Name(),
		)
	}
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	ptr := reflect.ValueOf(m[name])
	dst := ptr.Elem().Type()
	if !rt.AssignableTo(dst) {
		return Errorf(
			"fwk.SetProp: component [%s] has property [%s] with type [%s]. got value=%v (type=%s)",
			c.Name(),
			name,
			dst.Name(),
			value,
			rt.Name(),
		)
	}
	ptr.Elem().Set(rv)
	return nil
}

func (app *appmgr) GetProp(c Component, name string) (interface{}, Error) {
	m, ok := app.props[c]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property",
			c.Name(),
		)
	}

	ptr, ok := m[name]
	if !ok {
		return nil, Errorf(
			"fwk.GetProp: component [%s] didn't declare any property with name [%s]",
			c.Name(),
			name,
		)
	}

	v := reflect.Indirect(reflect.ValueOf(ptr)).Interface()
	return v, nil
}

func (app *appmgr) Run() Error {
	var err Error
	var ctx Context

	if app.state == fsm_UNDEFINED {
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
		fmt.Printf(">>> configuring [%s]...\n", svc.Name())
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
		fmt.Printf(">>> configuring [%s]...\n", tsk.Name())
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
		fmt.Printf(">>> starting [%s]...\n", svc.Name())
		err = svc.StartSvc(ctx)
		if err != nil {
			return err
		}
	}

	for _, tsk := range app.tsks {
		fmt.Printf(">>> starting [%s]...\n", tsk.Name())
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
				//fmt.Printf(">>> running [%s]...\n", tsk.Name())
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
	app.comps = nil
	app.tsks = nil
	app.svcs = nil
	app.state = fsm_OFFLINE

	app.props = nil
	app.dflow = nil
	app.store = nil

	g_mgr = nil
	return err
}

func init() {
	Register(reflect.TypeOf(appmgr{}))
}

// EOF
