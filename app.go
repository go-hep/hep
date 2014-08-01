package fwk

import (
	"io"
	"math"
	"reflect"
	"sort"
)

type fsm int

const (
	fsm_UNDEFINED fsm = iota
	fsm_CONFIGURING
	fsm_CONFIGURED
	fsm_STARTING
	fsm_STARTED
	fsm_RUNNING
	fsm_STOPPING
	fsm_STOPPED
	fsm_OFFLINE
)

func (state fsm) String() string {
	switch state {
	case fsm_UNDEFINED:
		return "UNDEFINED"
	case fsm_CONFIGURING:
		return "CONFIGURING"
	case fsm_CONFIGURED:
		return "CONFIGURED"
	case fsm_STARTING:
		return "STARTING"
	case fsm_STARTED:
		return "STARTED"
	case fsm_RUNNING:
		return "RUNNING"
	case fsm_STOPPING:
		return "STOPPING"
	case fsm_STOPPED:
		return "STOPPED"
	case fsm_OFFLINE:
		return "OFFLINE"

	default:
		panic(Errorf("invalid FSM value %d", int(state)))
	}
}

type appmgr struct {
	state fsm
	name  string

	props map[string]map[string]interface{}
	dflow *dflowsvc
	store *datastore
	msg   msgstream

	evtmax int64

	comps map[string]Component
	tsks  []Task
	svcs  []Svc
}

func NewApp() App {

	var err Error
	var app *appmgr

	const appname = "app"

	app = &appmgr{
		state: fsm_UNDEFINED,
		name:  appname,
		props: make(map[string]map[string]interface{}),
		dflow: nil,
		store: nil,
		msg: NewMsgStream(
			appname,
			LvlInfo,
			//LvlDebug,
			//LvlError,
			nil,
		),
		evtmax: -1,
		comps:  make(map[string]Component),
		tsks:   make([]Task, 0),
		svcs:   make([]Svc, 0),
	}

	svc, err := app.New("github.com/go-hep/fwk.datastore", "evtstore")
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not create evtstore: %v\n", err)
		return nil
	}
	app.store = svc.(*datastore)

	//app.store = newDataStore("evtstore")
	err = app.AddSvc(app.store)
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not create evtstore: %v\n", err)
		return nil
	}

	svc, err = app.New("github.com/go-hep/fwk.dflowsvc", "dataflow")
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not create dataflow svc: %v\n", err)
		return nil
	}
	app.dflow = svc.(*dflowsvc)

	//app.dflow = newDflowSvc("dataflow")
	err = app.AddSvc(app.dflow)
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not create dataflow svc: %v\n", err)
		return nil
	}

	return app
}

func (app *appmgr) Type() string {
	return "github.com/go-hep/fwk.appmgr"
}

func (app *appmgr) Name() string {
	return app.name
}

func (app *appmgr) Component(n string) Component {
	c, ok := app.comps[n]
	if !ok {
		return nil
	}
	return c
}

func (app *appmgr) addComponent(c Component) Error {
	app.comps[c.Name()] = c
	return nil
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
	cname := c.Name()
	_, ok := app.props[cname]
	if !ok {
		app.props[cname] = make(map[string]interface{})
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
	app.props[cname][name] = ptr
	return nil
}

func (app *appmgr) SetProp(c Component, name string, value interface{}) Error {
	cname := c.Name()
	m, ok := app.props[cname]
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
	cname := c.Name()
	m, ok := app.props[cname]
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

func (app *appmgr) DeclInPort(c Component, name string, t reflect.Type) Error {
	if app.state < fsm_CONFIGURING {
		return Errorf(
			"fwk.DeclInPort: invalid App state (%s). put the DeclInPort in Configure() of %s:%s",
			app.state,
			c.Type(),
			c.Name(),
		)
	}
	return app.dflow.addInNode(c.Name(), name, t)
}

func (app *appmgr) DeclOutPort(c Component, name string, t reflect.Type) Error {
	if app.state < fsm_CONFIGURING {
		return Errorf(
			"fwk.DeclOutPort: invalid App state (%s). put the DeclInPort in Configure() of %s:%s",
			app.state,
			c.Type(),
			c.Name(),
		)
	}
	return app.dflow.addOutNode(c.Name(), name, t)
}

func (app *appmgr) Run() Error {
	var err Error
	var ctx Context = context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", app.msg.lvl, nil),
	}

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
		if err != nil && err != io.EOF {
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
	app.msg.Debugf("configure...\n")
	app.state = fsm_CONFIGURING
	for _, svc := range app.svcs {
		app.msg.Debugf("configuring [%s]...\n", svc.Name())
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
		app.msg.Debugf("configuring [%s]...\n", tsk.Name())
		cfg, ok := tsk.(Configurer)
		if !ok {
			continue
		}
		err = cfg.Configure(ctx)
		if err != nil {
			return err
		}
	}

	app.msg.Infof(">>> --- [data flow] --- nodes...\n")
	for tsk, node := range app.dflow.nodes {
		app.msg.Infof(">>> ---[%s]---\n", tsk)
		app.msg.Infof("    in:  %v\n", node.in)
		app.msg.Infof("    out: %v\n", node.out)
	}

	app.msg.Infof(">>> --- [data flow] --- edges...\n")
	edges := make([]string, 0, len(app.dflow.edges))
	for n := range app.dflow.edges {
		edges = append(edges, n)
	}
	sort.Strings(edges)
	app.msg.Infof(" edges: %v\n", edges)

	app.state = fsm_CONFIGURED
	app.msg.Debugf("configure... [done]\n")
	return err
}

func (app *appmgr) start(ctx Context) Error {
	var err Error
	app.state = fsm_STARTING
	for _, svc := range app.svcs {
		app.msg.Debugf("starting [%s]...\n", svc.Name())
		err = svc.StartSvc(ctx)
		if err != nil {
			return err
		}
	}

	for _, tsk := range app.tsks {
		app.msg.Debugf("starting [%s]...\n", tsk.Name())
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
	if app.evtmax == -1 {
		app.evtmax = math.MaxInt64
	}
	for ievt := int64(0); ievt < app.evtmax; ievt++ {
		app.msg.Infof(">>> running evt=%d...\n", ievt)
		for k := range app.dflow.edges {
			//app.msg.Infof("--- edge [%s]... (%v)\n", k, rt)
			ch, ok := app.store.store[k]
			if ok {
				select {
				case vv := <-ch:
					if vv, ok := vv.(Deleter); ok {
						err = vv.Delete()
						if err != nil {
							return err
						}
					}
				default:
				}
			}
			app.store.store[k] = make(achan, 1)
		}
		errch := make(chan Error, len(app.tsks))
		for _, tsk := range app.tsks {
			go func(tsk Task) {
				//app.msg.Infof(">>> running [%s]...\n", tsk.Name())
				ctx := context{
					id:    ievt,
					slot:  0,
					store: app.store,
					msg:   NewMsgStream(tsk.Name(), app.msg.lvl, nil),
				}
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
	app.state = fsm_STOPPING
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

	return err
}

func (app *appmgr) Msg() MsgStream {
	return app.msg
}

func init() {
	Register(
		reflect.TypeOf(appmgr{}),
		func(t, name string, mgr App) (Component, Error) {
			app := NewApp().(*appmgr)
			app.name = name
			return app, nil
		},
	)
}

// EOF
