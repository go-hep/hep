package fwk

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"runtime"
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
	nprocs int

	comps map[string]Component
	tsks  []Task
	svcs  []Svc
}

func NewApp() App {

	var err error
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
		nprocs: 0,
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

	err = app.DeclProp(app, "EvtMax", &app.evtmax)
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not declare property 'EvtMax': %v\n", err)
		return nil
	}

	err = app.DeclProp(app, "NProcs", &app.nprocs)
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not declare property 'NProcs': %v\n", err)
		return nil
	}

	err = app.DeclProp(app, "MsgLevel", &app.msg.lvl)
	if err != nil {
		app.msg.Errorf("fwk.NewApp: could not declare property 'MsgLevel': %v\n", err)
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

func (app *appmgr) addComponent(c Component) error {
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

func (app *appmgr) AddTask(tsk Task) error {
	var err error
	app.tsks = append(app.tsks, tsk)
	app.comps[tsk.Name()] = tsk
	return err
}

func (app *appmgr) DelTask(tsk Task) error {
	var err error
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

func (app *appmgr) AddSvc(svc Svc) error {
	var err error
	app.svcs = append(app.svcs, svc)
	app.comps[svc.Name()] = svc
	return err
}

func (app *appmgr) DelSvc(svc Svc) error {
	var err error
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

func (app *appmgr) DeclProp(c Component, name string, ptr interface{}) error {
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

func (app *appmgr) SetProp(c Component, name string, value interface{}) error {
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

func (app *appmgr) GetProp(c Component, name string) (interface{}, error) {
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

func (app *appmgr) HasProp(c Component, name string) bool {
	cname := c.Name()
	_, ok := app.props[cname]
	if !ok {
		return ok
	}
	_, ok = app.props[cname][name]
	return ok
}

func (app *appmgr) DeclInPort(c Component, name string, t reflect.Type) error {
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

func (app *appmgr) DeclOutPort(c Component, name string, t reflect.Type) error {
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

func (app *appmgr) Run() error {
	var err error
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

func (app *appmgr) configure(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.msg.Debugf("configure...\n")
	app.state = fsm_CONFIGURING

	if app.evtmax == -1 {
		app.evtmax = math.MaxInt64
	}

	if app.nprocs < 0 {
		app.nprocs = runtime.NumCPU()
	}
	if app.nprocs > 1 {
		runtime.GOMAXPROCS(app.nprocs)
	}

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

	err = app.printDataFlow()
	if err != nil {
		return err
	}

	app.state = fsm_CONFIGURED
	app.msg.Debugf("configure... [done]\n")
	return err
}

func (app *appmgr) start(ctx Context) error {
	var err error
	defer app.msg.flush()
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

func (app *appmgr) run(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.state = fsm_RUNNING

	switch app.nprocs {
	case 0:
		err = app.runSequential(ctx)
	default:
		err = app.runConcurrent(ctx)
	}

	return err
}

func (app *appmgr) runSequential(ctx Context) error {
	var err error
	keys := app.dflow.keys()
	ctxs := make([]context, len(app.tsks))
	for j, tsk := range app.tsks {
		ctxs[j] = context{
			id:    -1,
			slot:  0,
			store: app.store,
			msg:   NewMsgStream(tsk.Name(), app.msg.lvl, nil),
		}
	}

	for ievt := int64(0); ievt < app.evtmax; ievt++ {
		app.msg.Infof(">>> running evt=%d...\n", ievt)
		err = app.store.reset(keys)
		if err != nil {
			return err
		}
		errch := make(chan error, len(app.tsks))
		quit := make(chan struct{})
		for i, tsk := range app.tsks {
			go func(i int, tsk Task) {
				//app.msg.Infof(">>> running [%s]...\n", tsk.Name())
				ctx := ctxs[i]
				ctx.id = ievt
				select {
				case errch <- tsk.Process(ctx):
					// FIXME(sbinet) dont be so eager to flush...
					ctx.msg.flush()
				case <-quit:
					// FIXME(sbinet) dont be so eager to flush...
					ctx.msg.flush()
				}
			}(i, tsk)
		}
		ndone := 0
	errloop:
		for err = range errch {
			ndone += 1
			if err != nil {
				close(quit)
				app.msg.flush()
				return err
			}
			if ndone == len(app.tsks) {
				break errloop
			}
		}
		app.msg.flush()
	}
	return err
}

func (app *appmgr) runConcurrent(ctx Context) error {
	var err error

	evts := make(chan int64, 100*app.nprocs)
	quit := make(chan struct{})
	done := make(chan struct{})
	errch := make(chan error)

	workers := make([]worker, app.nprocs)
	for i := 0; i < app.nprocs; i++ {
		workers[i] = worker{
			slot:  i,
			keys:  app.dflow.keys(),
			store: *app.store,
			ctxs:  make([]context, len(app.tsks)),
			msg:   NewMsgStream(fmt.Sprintf("%s-worker-%03d", app.name, i), app.msg.lvl, nil),
			evts:  evts,
			quit:  quit,
			errch: errch,
		}
		wrk := &workers[i]
		wrk.store.store = make(map[string]achan, len(app.dflow.keys()))
		for j, tsk := range app.tsks {
			wrk.ctxs[j] = context{
				id:    -1,
				slot:  i,
				store: &wrk.store,
				msg:   NewMsgStream(tsk.Name(), app.msg.lvl, nil),
			}
		}
		go func(wrk *worker) {
			wrk.run(app.tsks)
			done <- struct{}{}
		}(wrk)
	}

	go func() {
		for ievt := int64(0); ievt < app.evtmax; ievt++ {
			evts <- ievt
		}
		close(evts)
	}()

	ndone := 0
ctrl:
	for {
		select {
		case err = <-errch:
			if err != nil {
				// FIXME(sbinet) gather status of drained workers
				close(quit)
				return err
			}

		case <-done:
			ndone += 1
			if ndone == len(workers) {
				break ctrl
			}
		}
	}

	return err
}

func (app *appmgr) stop(ctx Context) error {
	var err error
	defer app.msg.flush()
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

func (app *appmgr) shutdown(ctx Context) error {
	var err error
	defer app.msg.flush()
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

func (app *appmgr) printDataFlow() error {
	var err error

	app.msg.Debugf(">>> --- [data flow] --- nodes...\n")
	for tsk, node := range app.dflow.nodes {
		app.msg.Debugf(">>> ---[%s]---\n", tsk)
		app.msg.Debugf("    in:  %v\n", node.in)
		app.msg.Debugf("    out: %v\n", node.out)
	}

	app.msg.Debugf(">>> --- [data flow] --- edges...\n")
	edges := make([]string, 0, len(app.dflow.edges))
	for n := range app.dflow.edges {
		edges = append(edges, n)
	}
	sort.Strings(edges)
	app.msg.Debugf(" edges: %v\n", edges)

	return err
}

func init() {
	Register(
		reflect.TypeOf(appmgr{}),
		func(t, name string, mgr App) (Component, error) {
			app := NewApp().(*appmgr)
			app.name = name
			return app, nil
		},
	)
}

// EOF
