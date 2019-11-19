// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job // import "go-hep.org/x/hep/fwk/job"

import (
	"go-hep.org/x/hep/fwk"
	"golang.org/x/xerrors"
)

// C describes the configuration data of a fwk.Component
type C struct {
	Name  string // name of the fwk.Component to create (eg "my-propagator")
	Type  string // type of the fwk.Component to create (eg "go-hep.org/x/hep/fads.Propagator")
	Props P      // properties of the fwk.Component to create
}

// P holds the configuration data (the properties) of a fwk.Component
// a map of key-value pairs.
type P map[string]interface{}

// Job is the go-based scripting interface to create, configure and run a fwk.App.
type Job struct {
	stmts []Stmt
	app   fwk.App
}

func newJob(app fwk.App, props P) *Job {
	if app == nil {
		app = fwk.NewApp()
	}

	job := &Job{
		stmts: []Stmt{
			{
				Type: StmtNewApp,
				Data: C{
					Name:  app.Name(),
					Type:  app.Type(),
					Props: props,
				},
			},
		},
		app: app,
	}
	for k, v := range props {
		job.setProp(app, k, v)
	}
	return job
}

// NewJob create a new Job from the given fwk.App value
// and configures it with the given properties P.
func NewJob(app fwk.App, props P) *Job {
	return newJob(app, props)
}

// New create a new Job with the default fwk.App implementation
// and configures it with the given properties P.
func New(props P) *Job {
	return newJob(nil, props)
}

// App returns the underlying fwk.App value of this Job.
func (job *Job) App() fwk.App {
	return job.app
}

// UI returns a fwk Scripter
func (job *Job) UI() UI {
	return UI{job.app.Scripter()}
}

// Create creates a fwk.Component according to the configuration
// data held by cfg.
// Create panics if no such component was registered with fwk.
func (job *Job) Create(cfg C) fwk.Component {
	c, err := job.app.New(cfg.Type, cfg.Name)
	if err != nil {
		job.Errorf("could not create [%s:%s]: %w\n", cfg.Type, cfg.Name, err)
		panic(err)
	}
	if cfg.Props == nil {
		return c
	}

	for k, v := range cfg.Props {
		job.setProp(c, k, v)
	}

	job.stmts = append(job.stmts, Stmt{
		Type: StmtCreate,
		Data: cfg,
	})
	return c
}

// SetProp sets the property name of the component c with the value v.
// SetProp panics if the component does not have such property or
// if the types do not match.
func (job *Job) SetProp(c fwk.Component, name string, value interface{}) {
	job.setProp(c, name, value)
	job.stmts = append(job.stmts, Stmt{
		Type: StmtSetProp,
		Data: C{
			Type: c.Type(),
			Name: c.Name(),
			Props: P{
				name: value,
			},
		},
	})
}

func (job *Job) setProp(c fwk.Component, name string, value interface{}) {
	if !job.app.HasProp(c, name) {
		err := xerrors.Errorf("component [%s:%s] has no property named %q\n",
			c.Type(),
			c.Name(),
			name,
		)
		job.Errorf("%+v", err)
		panic(err)
	}

	err := job.app.SetProp(c, name, value)
	if err != nil {
		job.Errorf(
			"could not set property name=%q value=%#v on component [%s]: %+v\n",
			name, value,
			c.Name(),
			err,
		)
		panic(err)
	}
}

// Run runs the underlying fwk.App.
// Run panics if an error occurred during any of the execution
// stages of the application.
func (job *Job) Run() {
	err := job.app.Run()
	if err != nil {
		job.Errorf(
			"could not run job: %+v\n",
			err,
		)
		panic(err)
	}
}

// Stmts returns the list of statements this Job has seen so far.
func (job *Job) Stmts() []Stmt {
	stmts := make([]Stmt, len(job.stmts))
	copy(stmts, job.stmts)
	return stmts
}

// Debugf displays a (formated) DBG message
func (job *Job) Debugf(format string, a ...interface{}) {
	job.app.Msg().Debugf(format, a...)
}

// Infof displays a (formated) INFO message
func (job *Job) Infof(format string, a ...interface{}) {
	job.app.Msg().Infof(format, a...)
}

// Warnf displays a (formated) WARN message
func (job *Job) Warnf(format string, a ...interface{}) {
	job.app.Msg().Warnf(format, a...)
}

// Errorf displays a (formated) ERR message
func (job *Job) Errorf(format string, a ...interface{}) {
	job.app.Msg().Errorf(format, a...)
}

/*
func (job *Job) Load(r io.Reader) error {
	var err error
	panic("not implemented")
	return err
}

func (job *Job) RunScripts(files ...string) error {
	var err error
	panic("not implemented")
	return err
}
*/
