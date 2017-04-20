// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// RunHeader provides metadata about a Run.
type RunHeader struct {
	RunNbr       int32
	Detector     string
	Descr        string
	SubDetectors []string
	Params       Params
}

func (*RunHeader) VersionSio() uint32 {
	return Version
}

// EventHeader provides metadata about an Event.
type EventHeader struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Blocks      []BlockDescr
	Params      Params
}

func (*EventHeader) VersionSio() uint32 {
	return Version
}

// BlockDescr describes a SIO block.
// BlockDescr provides the name of the SIO block and the type name of the
// data stored in that block.
type BlockDescr struct {
	Name string
	Type string
}

type Params struct {
	Ints    map[string][]int32
	Floats  map[string][]float32
	Strings map[string][]string
}

func (p Params) String() string {
	o := new(bytes.Buffer)
	for k, vec := range p.Ints {
		fmt.Fprintf(o, " parameter %s [int]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	for k, vec := range p.Floats {
		fmt.Fprintf(o, " parameter %s [float]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	for k, vec := range p.Strings {
		fmt.Fprintf(o, " parameter %s [string]: ", k)
		if len(vec) == 0 {
			fmt.Fprintf(o, " [empty] \n")
		}
		for _, v := range vec {
			fmt.Fprintf(o, "%v, ", v)
		}
		fmt.Fprintf(o, "\n")
	}
	return string(o.Bytes())
}

// Event holds informations about an LCIO event.
// Event also holds collection data for that event.
type Event struct {
	RunNumber   int32
	EventNumber int32
	TimeStamp   int64
	Detector    string
	Params      Params
	colls       map[string]interface{}
	names       []string
}

// Names returns the event data labels that define this event.
func (evt *Event) Names() []string {
	return evt.names
}

// Get returns the event data labelled name.
func (evt *Event) Get(name string) interface{} {
	return evt.colls[name]
}

// Has returns whether this event has data named name.
func (evt *Event) Has(name string) bool {
	_, ok := evt.colls[name]
	return ok
}

// Add attaches the (pointer to the) data ptr to this event,
// with the given name.
// Add panics if there is already some data labelled with the same name.
// Add panics if ptr is not a pointer to some data.
func (evt *Event) Add(name string, ptr interface{}) {
	if _, dup := evt.colls[name]; dup {
		panic(fmt.Errorf("lcio: duplicate key %q", name))
	}
	evt.names = append(evt.names, name)
	if evt.colls == nil {
		evt.colls = make(map[string]interface{})
	}
	if rv := reflect.ValueOf(ptr); rv.Type().Kind() != reflect.Ptr {
		panic("lcio: expects a pointer to a value")
	}
	evt.colls[name] = ptr
}

func (evt *Event) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, "        Event  : %d - run:   %d - timestamp %v - weight %v\n",
		evt.EventNumber, evt.RunNumber, evt.TimeStamp, evt.Weight(),
	)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, " date       %v\n", time.Unix(0, evt.TimeStamp).UTC().Format("02.01.2006 15:04:05.999999999"))
	fmt.Fprintf(o, " detector : %s\n", evt.Detector)
	fmt.Fprintf(o, " event parameters:\n%v\n", evt.Params)

	for _, name := range evt.names {
		coll := evt.colls[name]
		fmt.Fprintf(o, " collection name : %s\n parameters: \n%v\n", name, coll)
	}
	return string(o.Bytes())
}

func (evt *Event) Weight() float64 {
	if v, ok := evt.Params.Floats["_weight"]; ok {
		return float64(v[0])
	}
	return 1.0
}
