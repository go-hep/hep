package fwk

import (
	"io"
	"math"
	"reflect"
	"runtime"
	"sort"
	"time"

	"github.com/go-hep/fwk/fsm"
	nctx "golang.org/x/net/context"
)

type appmgr struct {
	state fsm.State
	name  string

	props map[string]map[string]interface{}
	dflow *dflowsvc
	store *datastore
	msg   msgstream

	evtmax int64
	nprocs int

	comps   map[string]Component
	tsks    []Task
	svcs    []Svc
	istream Task
	ctxs    [2][]context
}

// NewApp creates a (default) fwk application with (default and) sensible options.
func NewApp() App {

	var err error
	var app *appmgr

	const appname = "app"

	app = &appmgr{
		state: fsm.Undefined,
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
		nprocs: -1,
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

// Type returns the fully qualified type of this application
func (app *appmgr) Type() string {
	return "github.com/go-hep/fwk.appmgr"
}

// Name returns the name of this application
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
	if app.state < fsm.Configuring {
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
	if app.state < fsm.Configuring {
		return Errorf(
			"fwk.DeclOutPort: invalid App state (%s). put the DeclInPort in Configure() of %s:%s",
			app.state,
			c.Type(),
			c.Name(),
		)
	}
	return app.dflow.addOutNode(c.Name(), name, t)
}

func (app *appmgr) FSMState() fsm.State {
	return app.state
}

func (app *appmgr) Run() error {
	var err error
	ctx := context{
		id:    0,
		slot:  0,
		store: nil,
		msg:   NewMsgStream("<root>", app.msg.lvl, nil),
		mgr:   app,
	}

	start := time.Now()
	var mstart runtime.MemStats
	runtime.ReadMemStats(&mstart)

	if app.state == fsm.Undefined {
		err = app.configure(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm.Configured {
		err = app.start(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm.Started {
		err = app.run(ctx)
		if err != nil && err != io.EOF {
			return err
		}
	}

	if app.state == fsm.Running {
		err = app.stop(ctx)
		if err != nil {
			return err
		}
	}

	if app.state == fsm.Stopped {
		err = app.shutdown(ctx)
		if err != nil {
			return err
		}
	}

	app.msg.Infof("cpu: %v\n", time.Since(start))
	var mdone runtime.MemStats
	runtime.ReadMemStats(&mdone)

	diff := func(v1, v0 uint64) int64 {
		if v0 > v1 {
			return -int64(v0 - v1)
		}
		return int64(v1 - v0)
	}

	app.msg.Infof("mem: alloc:     %10d kB\n", diff(mdone.Alloc, mstart.Alloc)/1024)
	app.msg.Infof("mem: tot-alloc: %10d kB\n", diff(mdone.TotalAlloc, mstart.TotalAlloc)/1024)
	app.msg.Infof("mem: n-mallocs: %10d\n", diff(mdone.Mallocs, mstart.Mallocs))
	app.msg.Infof("mem: n-frees:   %10d\n", diff(mdone.Frees, mstart.Frees))
	app.msg.Infof("mem: gc-pauses: %10d ms\n", diff(mdone.PauseTotalNs, mstart.PauseTotalNs)/1000000)

	return err
}

func (app *appmgr) Scripter() Scripter {
	return &irunner{app}
}

func (app *appmgr) configure(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.msg.Debugf("configure...\n")
	app.state = fsm.Configuring

	if app.evtmax == -1 {
		app.evtmax = math.MaxInt64
	}

	if app.nprocs < 0 {
		app.nprocs = runtime.NumCPU()
	}

	tsks := make([]context, len(app.tsks))
	for j, tsk := range app.tsks {
		tsks[j] = context{
			id:    -1,
			slot:  0,
			store: app.store,
			msg:   NewMsgStream(tsk.Name(), app.msg.lvl, nil),
			mgr:   app,
		}
	}

	svcs := make([]context, len(app.svcs))
	for j, svc := range app.svcs {
		svcs[j] = context{
			id:    -1,
			slot:  0,
			store: app.store,
			msg:   NewMsgStream(svc.Name(), app.msg.lvl, nil),
			mgr:   app,
		}
	}

	for i, svc := range app.svcs {
		app.msg.Debugf("configuring [%s]...\n", svc.Name())
		cfg, ok := svc.(Configurer)
		if !ok {
			continue
		}
		err = cfg.Configure(svcs[i])
		if err != nil {
			return err
		}
	}

	for i, tsk := range app.tsks {
		app.msg.Debugf("configuring [%s]...\n", tsk.Name())
		cfg, ok := tsk.(Configurer)
		if !ok {
			continue
		}
		err = cfg.Configure(tsks[i])
		if err != nil {
			return err
		}
	}

	err = app.printDataFlow()
	if err != nil {
		return err
	}

	app.ctxs[0] = tsks
	app.ctxs[1] = svcs
	app.state = fsm.Configured
	app.msg.Debugf("configure... [done]\n")
	return err
}

func (app *appmgr) start(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.state = fsm.Starting
	for i, svc := range app.svcs {
		app.msg.Debugf("starting [%s]...\n", svc.Name())
		err = svc.StartSvc(app.ctxs[1][i])
		if err != nil {
			return err
		}
	}

	for i, tsk := range app.tsks {
		app.msg.Debugf("starting [%s]...\n", tsk.Name())
		err = tsk.StartTask(app.ctxs[0][i])
		if err != nil {
			return err
		}
	}

	app.state = fsm.Started
	return err
}

func (app *appmgr) run(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.state = fsm.Running

	maxprocs := runtime.GOMAXPROCS(app.nprocs)

	switch app.nprocs {
	case 0:
		err = app.runSequential(ctx)
	default:
		err = app.runConcurrent(ctx)
	}

	runtime.GOMAXPROCS(maxprocs)

	return err
}

func (app *appmgr) runSequential(ctx Context) error {
	var err error

	runctx, runCancel := nctx.WithCancel(nctx.Background())
	defer runCancel()

	keys := app.dflow.keys()
	ctxs := make([]context, len(app.tsks))
	store := *app.store
	for j, tsk := range app.tsks {
		ctxs[j] = context{
			id:    -1,
			slot:  0,
			store: &store,
			msg:   NewMsgStream(tsk.Name(), app.msg.lvl, nil),
			mgr:   app,
		}
	}

	ictrl, err := app.startInputStream()
	if err != nil {
		return err
	}
	defer close(ictrl.Quit)

	octrl, err := app.startOutputStreams()
	if err != nil {
		return err
	}

	defer close(octrl.Quit)

	for ievt := int64(0); ievt < app.evtmax; ievt++ {
		evtctx, evtCancel := nctx.WithCancel(runctx)

		app.msg.Infof(">>> running evt=%d...\n", ievt)
		err = store.reset(keys)
		if err != nil {
			return err
		}
		err = app.istream.Process(ctxs[0])
		if err != nil {
			store.close()
			app.msg.flush()
			return err
		}
		run := taskrunner{
			ievt:   ievt,
			errc:   make(chan error, len(app.tsks)),
			evtctx: evtctx,
		}
		for i, tsk := range app.tsks {
			go run.run(i, ctxs[i], tsk)
		}
		ndone := 0
	errloop:
		for err = range run.errc {
			ndone++
			if err != nil {
				evtCancel()
				store.close()
				app.msg.flush()
				return err
			}
			if ndone == len(app.tsks) {
				break errloop
			}
		}
		store.close()
		app.msg.flush()
	}

	return err
}

func (app *appmgr) runConcurrent(ctx Context) error {
	var err error

	runctx, runCancel := nctx.WithCancel(nctx.Background())
	defer runCancel()

	ctrl := workercontrol{
		evts:   make(chan context, 2*app.nprocs),
		done:   make(chan struct{}),
		errc:   make(chan error),
		runctx: runctx,
	}

	istream, err := app.startInputStream()
	if err != nil {
		return err
	}
	defer close(istream.Quit)

	ostream, err := app.startOutputStreams()
	if err != nil {
		return err
	}
	defer close(ostream.Quit)

	workers := make([]worker, app.nprocs)
	for i := 0; i < app.nprocs; i++ {
		workers[i] = *newWorker(i, app, &ctrl)
	}

	go func() {
		keys := app.dflow.keys()
		msg := NewMsgStream(app.istream.Name(), app.msg.lvl, nil)
		for ievt := int64(0); ievt < app.evtmax; ievt++ {
			evtctx, evtCancel := nctx.WithCancel(runctx)
			store := *app.store
			store.store = make(map[string]achan, len(keys))
			err := store.reset(keys)
			if err != nil {
				evtCancel()
				close(ctrl.evts)
				ctrl.errc <- err
				return
			}
			ctx := context{
				id:    ievt,
				slot:  0,
				store: &store,
				msg:   msg,
				mgr:   nil, // nobody's supposed to access mgr's state during event-loop
				ctx:   evtctx,
			}

			err = app.istream.Process(ctx)
			if err != nil {
				if err != io.EOF {
					evtCancel()
					ctrl.errc <- err
				}
				close(ctrl.evts)
				return
			}
			ctrl.evts <- ctx
		}
		close(ctrl.evts)
	}()

	ndone := 0
ctrl:
	for {
		select {
		case eworker, ok := <-ctrl.errc:
			if !ok {
				continue
			}
			if eworker != nil && err == nil {
				// only record first error.
				// FIXME(sbinet) record all of them (errstack)
				err = eworker
			}

		case <-runctx.Done():
			return runctx.Err()

		case <-ctrl.done:
			ndone++
			app.msg.Infof("workers done: %d/%d\n", ndone, app.nprocs)
			if ndone == len(workers) {
				break ctrl
			}
		}
	}

	return err
}

func (app *appmgr) startInputStream() (StreamControl, error) {
	var err error

	ctrl := StreamControl{
		Ctx:  make(chan Context),
		Err:  make(chan error), // FIXME: impl. back-pressure
		Quit: make(chan struct{}),
	}

	idx := -1
	inputs := make([]*InputStream, 0)

	// collect input streams
	for i, tsk := range app.tsks {
		in, ok := tsk.(*InputStream)
		if !ok {
			continue
		}
		inputs = append(inputs, in)
		idx = i
	}

	switch len(inputs) {
	case 0:
		// create an event "pumper"
		tsk := &inputStream{NewTask("fwk.inputStream", "app-evtloop", app)}
		app.istream = tsk
	case 1:
		app.istream = inputs[0]
		app.tsks = append(app.tsks[:idx], app.tsks[idx+1:]...)
		err := inputs[0].connect(ctrl)
		if err != nil {
			return ctrl, err
		}

	default:
		return ctrl, Errorf("found more than one InputStream! (n=%d)", len(inputs))
	}

	return ctrl, err
}

func (app *appmgr) startOutputStreams() (StreamControl, error) {
	var err error

	ctrl := StreamControl{
		Ctx:  make(chan Context),
		Err:  make(chan error), // FIXME: impl. back-pressure
		Quit: make(chan struct{}),
	}

	// start output streams
	for _, tsk := range app.tsks {
		in, ok := tsk.(*OutputStream)
		if !ok {
			continue
		}
		err = in.connect(ctrl)
		if err != nil {
			return ctrl, err
		}
	}

	return ctrl, err
}

func (app *appmgr) stop(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.state = fsm.Stopping
	for i, tsk := range app.tsks {
		err = tsk.StopTask(app.ctxs[0][i])
		if err != nil {
			return err
		}
	}

	for i, svc := range app.svcs {
		err = svc.StopSvc(app.ctxs[1][i])
		if err != nil {
			return err
		}
	}

	app.state = fsm.Stopped
	return err
}

func (app *appmgr) shutdown(ctx Context) error {
	var err error
	defer app.msg.flush()
	app.comps = nil
	app.tsks = nil
	app.svcs = nil
	app.state = fsm.Offline

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
