package fwk

import (
	"reflect"
)

type node struct {
	in  map[string]struct{}
	out map[string]struct{}
}

func newNode() *node {
	return &node{
		in:  make(map[string]struct{}),
		out: make(map[string]struct{}),
	}
}

type dflowsvc struct {
	SvcBase
	nodes map[string]*node
	edges map[string]struct{}
}

func (svc *dflowsvc) Configure(ctx Context) Error {
	svc.nodes = make(map[string]*node)
	svc.edges = make(map[string]struct{})
	return nil
}

func (svc *dflowsvc) StartSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) StopSvc(ctx Context) Error {
	return nil
}

func (svc *dflowsvc) addInNode(tsk string, name string) Error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.in[name]
	if ok {
		return Errorf(
			"fwk.DeclInPort: component [%s] already declare in-port with name [%s]",
			tsk,
			name,
		)
	}

	node.in[name] = struct{}{}
	svc.edges[name] = struct{}{}
	return nil
}

func (svc *dflowsvc) addOutNode(tsk string, name string) Error {
	node, ok := svc.nodes[tsk]
	if !ok {
		node = newNode()
		svc.nodes[tsk] = node
	}
	_, ok = node.out[name]
	if ok {
		return Errorf(
			"fwk.DeclInPort: component [%s] already declare out-port with name [%s]",
			tsk,
			name,
		)
	}

	node.out[name] = struct{}{}
	svc.edges[name] = struct{}{}
	return nil
}

func init() {
	Register(reflect.TypeOf(dflowsvc{}))
}

// EOF
