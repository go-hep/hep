package fwk

import (
	"bytes"
	"fmt"
	"reflect"
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

type dflowsvc struct {
	SvcBase
	nodes map[string]*node
	edges map[string]reflect.Type
}

func (svc *dflowsvc) Configure(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) StartSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) StopSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) addInNode(tsk string, name string, t reflect.Type) Error {
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
			type elem_t struct {
				port string // in/out
				task string // task which defined the port
				typ  reflect.Type
			}
			cont := []elem_t{}
			for tskname, node := range svc.nodes {
				for k, in := range node.in {
					if k != name {
						continue
					}
					cont = append(cont,
						elem_t{
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
						elem_t{
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
			return fmt.Errorf(string(o.Bytes()))
		}
	}

	svc.edges[name] = t
	return nil
}

func (svc *dflowsvc) addOutNode(tsk string, name string, t reflect.Type) Error {
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
		for duptsk, dupnode := range svc.nodes {
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
		func(name string, mgr App) (Component, Error) {
			svc := &dflowsvc{
				SvcBase: NewSvc(name, mgr),
				nodes:   make(map[string]*node),
				edges:   make(map[string]reflect.Type),
			}
			return svc, nil
		},
	)
}

// EOF
