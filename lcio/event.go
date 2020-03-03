// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"go-hep.org/x/hep/sio"
)

// RunHeader provides metadata about a Run.
type RunHeader struct {
	RunNumber    int32
	Detector     string
	Descr        string
	SubDetectors []string
	Params       Params
}

func (hdr *RunHeader) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, "        Run:   %d\n", hdr.RunNumber)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, " description: %s\n", hdr.Descr)
	fmt.Fprintf(o, " detector   : %s\n", hdr.Detector)
	fmt.Fprintf(o, " sub-dets   : %v\n", hdr.SubDetectors)
	fmt.Fprintf(o, " parameters :\n%v\n", hdr.Params)
	return string(o.Bytes())
}

func (*RunHeader) VersionSio() uint32 {
	return Version
}

func (hdr *RunHeader) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hdr.RunNumber)
	enc.Encode(&hdr.Detector)
	enc.Encode(&hdr.Descr)
	enc.Encode(&hdr.SubDetectors)
	enc.Encode(&hdr.Params)
	return enc.Err()
}

func (hdr *RunHeader) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hdr.RunNumber)
	dec.Decode(&hdr.Detector)
	dec.Decode(&hdr.Descr)
	dec.Decode(&hdr.SubDetectors)
	dec.Decode(&hdr.Params)
	return dec.Err()
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

func (hdr *EventHeader) String() string {
	o := new(bytes.Buffer)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, "        Event  : %d - run:   %d - timestamp %v - weight %v\n",
		hdr.EventNumber, hdr.RunNumber, hdr.TimeStamp, hdr.Weight(),
	)
	fmt.Fprintf(o, "%s\n", strings.Repeat("=", 80))
	fmt.Fprintf(o, " date       %v\n", time.Unix(0, hdr.TimeStamp).UTC().Format("02.01.2006 15:04:05.999999999"))
	fmt.Fprintf(o, " detector : %s\n", hdr.Detector)
	fmt.Fprintf(o, " event parameters:\n%v", hdr.Params)

	w := tabwriter.NewWriter(o, 8, 4, 1, ' ', 0)
	for _, blk := range hdr.Blocks {
		fmt.Fprintf(w, " collection name : %s\t(%s)\n", blk.Name, blk.Type)
	}
	w.Flush()

	return string(o.Bytes())
}

func (*EventHeader) VersionSio() uint32 {
	return Version
}

func (hdr *EventHeader) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	enc.Encode(&hdr.RunNumber)
	enc.Encode(&hdr.EventNumber)
	enc.Encode(&hdr.TimeStamp)
	enc.Encode(&hdr.Detector)
	enc.Encode(&hdr.Blocks)
	enc.Encode(&hdr.Params)
	return enc.Err()
}

func (hdr *EventHeader) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&hdr.RunNumber)
	dec.Decode(&hdr.EventNumber)
	dec.Decode(&hdr.TimeStamp)
	dec.Decode(&hdr.Detector)
	dec.Decode(&hdr.Blocks)
	dec.Decode(&hdr.Params)
	return dec.Err()
}

func (hdr *EventHeader) Weight() float64 {
	if v, ok := hdr.Params.Floats["_weight"]; ok {
		return float64(v[0])
	}
	return 1.0
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
	{
		keys := make([]string, 0, len(p.Ints))
		for k := range p.Ints {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vec := p.Ints[k]
			fmt.Fprintf(o, " parameter %s [int]: ", k)
			if len(vec) == 0 {
				fmt.Fprintf(o, " [empty] \n")
			}
			for _, v := range vec {
				fmt.Fprintf(o, "%v, ", v)
			}
			fmt.Fprintf(o, "\n")
		}
	}
	{
		keys := make([]string, 0, len(p.Floats))
		for k := range p.Floats {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vec := p.Strings[k]
			fmt.Fprintf(o, " parameter %s [float]: ", k)
			if len(vec) == 0 {
				fmt.Fprintf(o, " [empty] \n")
			}
			for _, v := range vec {
				fmt.Fprintf(o, "%v, ", v)
			}
			fmt.Fprintf(o, "\n")
		}
	}
	{
		keys := make([]string, 0, len(p.Strings))
		for k := range p.Floats {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			vec := p.Strings[k]
			fmt.Fprintf(o, " parameter %s [string]: ", k)
			if len(vec) == 0 {
				fmt.Fprintf(o, " [empty] \n")
			}
			for _, v := range vec {
				fmt.Fprintf(o, "%v, ", v)
			}
			fmt.Fprintf(o, "\n")
		}
	}
	return string(o.Bytes())
}

func (*Params) VersionSio() uint32 {
	return Version
}

func (p *Params) MarshalSio(w sio.Writer) error {
	enc := sio.NewEncoder(w)
	{
		data := p.Ints
		enc.Encode(int32(len(data)))
		if n := len(data); n > 0 {
			keys := make([]string, 0, n)
			for k := range data {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				enc.Encode(&k)
				v := data[k]
				enc.Encode(&v)
			}
		}
	}
	{
		data := p.Floats
		enc.Encode(int32(len(data)))
		if n := len(data); n > 0 {
			keys := make([]string, 0, n)
			for k := range data {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				enc.Encode(&k)
				v := data[k]
				enc.Encode(&v)
			}
		}
	}
	{
		data := p.Strings
		enc.Encode(int32(len(data)))
		if n := len(data); n > 0 {
			keys := make([]string, 0, n)
			for k := range data {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			for _, k := range keys {
				enc.Encode(&k)
				v := data[k]
				enc.Encode(&v)
			}
		}
	}
	return enc.Err()
}

func (p *Params) UnmarshalSio(r sio.Reader) error {
	dec := sio.NewDecoder(r)
	dec.Decode(&p.Ints)
	dec.Decode(&p.Floats)
	dec.Decode(&p.Strings)
	return dec.Err()
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

var (
	_ sio.Versioner = (*RunHeader)(nil)
	_ sio.Codec     = (*RunHeader)(nil)
	_ sio.Versioner = (*EventHeader)(nil)
	_ sio.Codec     = (*EventHeader)(nil)
	_ sio.Versioner = (*Params)(nil)
	_ sio.Codec     = (*Params)(nil)
)
