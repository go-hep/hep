package fwk

import (
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
	Base
	nodes map[Task]*node
	edges map[string]struct{}
}

func newDflowSvc(name string) *dflowsvc {
	return &dflowsvc{
		Base: Base{
			Name: name,
			Type: "fwk.dflowsvc",
		},
		nodes: make(map[Task]*node),
		edges: make(map[string]struct{}),
	}
}

func (svc *dflowsvc) StartSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) StopSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) addInNode(tsk Task, name string, value interface{}) Error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.in[name]
	if ok {
		return Errorf(
			"fwk.DeclInPort: component [%s] already declare in-port with name [%s]",
			tsk.CompName(),
			name,
		)
	}

	node.in[name] = reflect.TypeOf(value)
	svc.edges[name] = struct{}{}
	return nil
}

func (svc *dflowsvc) addOutNode(tsk Task, name string, value interface{}) Error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.out[name]
	if ok {
		return Errorf(
			"fwk.DeclInPort: component [%s] already declare out-port with name [%s]",
			tsk.CompName(),
			name,
		)
	}

	node.out[name] = reflect.TypeOf(value)
	svc.edges[name] = struct{}{}
	return nil
}

// EOF
