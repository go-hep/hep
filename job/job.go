package job

import (
	"github.com/go-hep/fwk"
)

// C describes the configuration data of a fwk.Component
type C struct {
	Name  string // name of the fwk.Component to create (eg "my-propagator")
	Type  string // type of the fwk.Component to create (eg "github.com/go-hep/fads.Propagator")
	Props P      // properties of the fwk.Component to create
}

// P holds the configuration data (the properties) of a fwk.Component
// a map of key-value pairs.
type P map[string]interface{}

// Job is the go-based scripting interface to create, configure and run a fwk.App.
type Job struct {
	app fwk.App
}

func New(app fwk.App, props P) *Job {
	if app == nil {
		app = fwk.NewApp()
	}

	job := &Job{app: app}
	for k, v := range props {
		job.SetProp(app, k, v)

	}
	return job
}

func (job *Job) App() fwk.App {
	return job.app
}

// Create creates a fwk.Component according to the configuration
// data held by cfg.
// Create panics if no such component was registered with fwk.
func (job *Job) Create(cfg C) fwk.Component {
	c, err := job.app.New(cfg.Type, cfg.Name)
	if err != nil {
		job.Errorf("could not create [%s:%s]: %v\n", cfg.Type, cfg.Name, err)
		panic(err)
	}
	if cfg.Props == nil {
		return c
	}

	for k, v := range cfg.Props {
		job.SetProp(c, k, v)
	}
	return c
}

// SetProp sets the property name of the component c with the value v.
// SetProp panics if the component does not have such property or
// if the types do not match.
func (job *Job) SetProp(c fwk.Component, name string, value interface{}) {
	if !job.app.HasProp(c, name) {
		err := fwk.Errorf("component [%s:%s] has no property named %q\n",
			c.Type(),
			c.Name(),
			name,
		)
		job.Errorf(err.Error())
		panic(err)
	}

	err := job.app.SetProp(c, name, value)
	if err != nil {
		job.Errorf(
			"could not set property name=%q value=%#v on component [%s]: %v\n",
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
			"could not run job: %v\n",
			err,
		)
		panic(err)
	}
}

// Debugf displays a (formated) DBG message
func (job *Job) Debugf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Debugf(format, a...)
}

// Infof displays a (formated) INFO message
func (job *Job) Infof(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Infof(format, a...)
}

// Warnf displays a (formated) WARN message
func (job *Job) Warnf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Warnf(format, a...)
}

// Errorf displays a (formated) ERR message
func (job *Job) Errorf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Errorf(format, a...)
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
