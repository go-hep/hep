package fads

import (
	"reflect"

	"github.com/go-hep/fwk"
)

type ObjPair struct {
	In  string
	Out string
}

type UniqueObjectFinder struct {
	fwk.TaskBase

	colls []ObjPair
}

func (tsk *UniqueObjectFinder) Configure(ctx fwk.Context) error {
	var err error

	for _, pair := range tsk.colls {
		err = tsk.DeclInPort(pair.In, reflect.TypeOf([]Candidate{}))
		if err != nil {
			return err
		}

		err = tsk.DeclOutPort(pair.Out, reflect.TypeOf([]Candidate{}))
		if err != nil {
			return err
		}
	}

	return err
}

func (tsk *UniqueObjectFinder) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *UniqueObjectFinder) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *UniqueObjectFinder) Process(ctx fwk.Context) error {
	var err error
	store := ctx.Store()

	type Pair struct {
		In  []Candidate
		Out []Candidate
	}

	colls := make([]Pair, 0, len(tsk.colls))
	for i, pair := range tsk.colls {
		v, err2 := store.Get(pair.In)
		if err2 != nil {
			return err2
		}

		input := v.([]Candidate)
		output := make([]Candidate, 0, len(input))
		colls = append(colls, Pair{
			In:  input,
			Out: output,
		})

		defer func(i int) {
			err2 := store.Put(tsk.colls[i].Out, colls[i].Out)
			if err2 != nil {
				err = err2
			}
		}(i)
	}

	for icol := range colls {
		pair := &colls[icol]
		input := pair.In
		output := pair.Out
		for i := range input {
			cand := &input[i]
			unique := false
		uniqueloop:
			for jcol := 0; jcol < icol; jcol++ {
				jcands := colls[jcol].In
				for j := range jcands {
					jcand := &jcands[j]
					if cand.Overlaps(jcand) {
						unique = true
						break uniqueloop
					}
				}
			}

			if !unique {
				continue
			}
			output = append(output, *cand)
		}

		pair.Out = output
	}

	return err
}

func newUniqueObjectFinder(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &UniqueObjectFinder{
		TaskBase: fwk.NewTask(typ, name, mgr),
		colls:    make([]ObjPair, 0),
	}

	err = tsk.DeclProp("Keys", &tsk.colls)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(UniqueObjectFinder{}), newUniqueObjectFinder)
}
