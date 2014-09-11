package fwk

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"

	"github.com/go-hep/fwk/utils/tarjan"
)

type node struct {
	in  map[string]reflect.Type
	out map[string]reflect.Type
}

func newNode() *node {
	return &node{
		in:  make(map[string]reflect.Type),
		out: make(map[string]reflect.Type),
	}
}

// dflowsvc models and describes the runtime data-flow and (data) dependencies between
// components as declared during configuration.
type dflowsvc struct {
	SvcBase
	nodes map[string]*node
	edges map[string]reflect.Type
}

func (svc *dflowsvc) Configure(ctx Context) error {
	return nil
}

func (svc *dflowsvc) StartSvc(ctx Context) error {
	var err error

	// sort node-names for reproducibility
	nodenames := make([]string, 0, len(svc.nodes))
	for n := range svc.nodes {
		nodenames = append(nodenames, n)
	}
	sort.Strings(nodenames)

	// - make sure all input keys of components are available
	//   as output keys of a task
	// - also detect whether a key is labeled as an out-port
	//   by 2 different components
	out := make(map[string]string) // outport-name -> producer-name
	for _, tsk := range nodenames {
		node := svc.nodes[tsk]
		for k := range node.out {
			n, dup := out[k]
			if dup {
				return Errorf("%s: component [%s] already declared port [%s] as its output (current=%s)",
					svc.Name(), n, k, tsk,
				)
			}
			out[k] = tsk
		}
	}

	for _, tsk := range nodenames {
		node := svc.nodes[tsk]
		for k := range node.in {
			_, ok := out[k]
			if !ok {
				return Errorf("%s: component [%s] declared port [%s] as input but NO KNOWN producer",
					svc.Name(), tsk, k,
				)
			}
		}
	}

	// detect cycles.
	graph := make(map[interface{}][]interface{})
	for _, n := range nodenames {
		node := svc.nodes[n]
		graph[n] = []interface{}{}
		for in := range node.in {
			for _, o := range nodenames {
				if o == n {
					continue
				}
				onode := svc.nodes[o]
				connected := false
				for out := range onode.out {
					if in == out {
						connected = true
						break
					}
				}
				if connected {
					graph[n] = append(graph[n], o)
				}
			}
		}
	}

	cycles := tarjan.Connections(graph)
	if len(cycles) > 0 {
		msg := ctx.Msg()
		ncycles := 0
		for _, cycle := range cycles {
			if len(cycle) > 1 {
				ncycles++
				msg.Errorf("cycle detected: %v\n", cycle)
			}
		}
		s := ""
		if ncycles > 1 {
			s = "s"
		}
		if ncycles > 0 {
			return Errorf("%s: cycle%s detected: %d", svc.Name(), s, ncycles)
		}
	}

	return err
}

func (svc *dflowsvc) StopSvc(ctx Context) error {
	return nil
}

func (svc *dflowsvc) keys() []string {
	keys := make([]string, 0, len(svc.edges))
	for k := range svc.edges {
		keys = append(keys, k)
	}
	return keys
}

func (svc *dflowsvc) addInNode(tsk string, name string, t reflect.Type) error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.in[name]
	if ok {
		return Errorf(
			"fwk.DeclInPort: component [%s] already declared in-port with name [%s]",
			tsk,
			name,
		)
	}

	node.in[name] = t
	edgetyp, dup := svc.edges[name]
	if dup {
		// make sure types match
		if edgetyp != t {
			type elemT struct {
				port string // in/out
				task string // task which defined the port
				typ  reflect.Type
			}
			cont := []elemT{}
			nodenames := make([]string, 0, len(svc.nodes))
			for tskname := range svc.nodes {
				nodenames = append(nodenames, tskname)
			}
			sort.Strings(nodenames)
			for _, tskname := range nodenames {
				node := svc.nodes[tskname]
				for k, in := range node.in {
					if k != name {
						continue
					}
					cont = append(cont,
						elemT{
							port: "in ",
							task: tskname,
							typ:  in,
						},
					)
				}
				for k, out := range node.out {
					if k != name {
						continue
					}
					cont = append(cont,
						elemT{
							port: "out",
							task: tskname,
							typ:  out,
						},
					)
				}
			}
			var o bytes.Buffer
			fmt.Fprintf(&o, "fwk.DeclInPort: detected type inconsistency for port [%s]:\n", name)
			for _, c := range cont {
				fmt.Fprintf(&o, " component=%q port=%s type=%v\n", c.task, c.port, c.typ)
			}
			return Errorf(string(o.Bytes()))
		}
	}

	svc.edges[name] = t
	return nil
}

func (svc *dflowsvc) addOutNode(tsk string, name string, t reflect.Type) error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.out[name]
	if ok {
		return Errorf(
			"fwk.DeclOutPort: component [%s] already declared out-port with name [%s]",
			tsk,
			name,
		)
	}

	node.out[name] = t

	edgetyp, dup := svc.edges[name]
	if dup {
		// edge already exists
		// loop over nodes, find out who already defined that edge
		nodenames := make([]string, 0, len(svc.nodes))
		for tskname := range svc.nodes {
			nodenames = append(nodenames, tskname)
		}
		sort.Strings(nodenames)
		for _, duptsk := range nodenames {
			dupnode := svc.nodes[duptsk]
			if duptsk == tsk {
				continue
			}
			for out := range dupnode.out {
				if out == name {
					return Errorf(
						"fwk.DeclOutPort: component [%s] already declared out-port with name [%s (type=%v)].\nfwk.DeclOutPort: component [%s] is trying to add a duplicate out-port [%s (type=%v)]",
						duptsk,
						name,
						edgetyp,
						tsk,
						name,
						t,
					)
				}
			}
		}
	}
	svc.edges[name] = t
	return nil
}

func init() {
	Register(reflect.TypeOf(dflowsvc{}),
		func(t, name string, mgr App) (Component, error) {
			svc := &dflowsvc{
				SvcBase: NewSvc(t, name, mgr),
				nodes:   make(map[string]*node),
				edges:   make(map[string]reflect.Type),
			}
			return svc, nil
		},
	)
}

// EOF
