package job

import (
	"io"

	"github.com/go-hep/fwk"
)

type C struct {
	Name  string
	Type  string
	Props P
}

type P map[string]interface{}

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

func (job *Job) Debugf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Debugf(format, a...)
}

func (job *Job) Infof(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Infof(format, a...)
}

func (job *Job) Warnf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Warnf(format, a...)
}

func (job *Job) Errorf(format string, a ...interface{}) (int, error) {
	return job.app.Msg().Errorf(format, a...)
}

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
