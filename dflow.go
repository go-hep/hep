package fwk

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
	Base
	nodes map[string]*node
	edges map[string]struct{}
}

func newDflowSvc(name string) *dflowsvc {
	return &dflowsvc{
		Base: Base{
			Name: name,
			Type: "fwk.dflowsvc",
		},
		nodes: make(map[string]*node),
		edges: make(map[string]struct{}),
	}
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

// EOF
