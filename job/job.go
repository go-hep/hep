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

func New(app fwk.App) *Job {
	if app == nil {
		app = fwk.NewApp()
	}
	return &Job{
		app: app,
	}
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

func (job *Job) Debugf(format string, a ...interface{}) (int, fwk.Error) {
	return job.app.Msg().Debugf(format, a...)
}

func (job *Job) Infof(format string, a ...interface{}) (int, fwk.Error) {
	return job.app.Msg().Infof(format, a...)
}

func (job *Job) Warnf(format string, a ...interface{}) (int, fwk.Error) {
	return job.app.Msg().Warnf(format, a...)
}

func (job *Job) Errorf(format string, a ...interface{}) (int, fwk.Error) {
	return job.app.Msg().Errorf(format, a...)
}

func (job *Job) Load(r io.Reader) fwk.Error {
	var err fwk.Error
	panic("not implemented")
	return err
}

func (job *Job) RunScripts(files ...string) fwk.Error {
	panic("not implemented")
}
