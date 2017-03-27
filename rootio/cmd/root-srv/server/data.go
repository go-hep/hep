// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"bytes"
	"fmt"
	"math"
	"strings"
	"sync"

	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgsvg"

	"go-hep.org/x/hep/hplot"
	"go-hep.org/x/hep/rootio"
)

type dbFiles struct {
	sync.RWMutex
	files map[string]*rootio.File
}

func newDbFiles() *dbFiles {
	return &dbFiles{
		files: make(map[string]*rootio.File),
	}
}

func (db *dbFiles) close() {
	db.Lock()
	defer db.Unlock()
	for _, f := range db.files {
		f.Close()
	}
	db.files = nil
}

func (db *dbFiles) get(name string) *rootio.File {
	db.RLock()
	defer db.RUnlock()
	f, _ := db.files[name]
	return f
}

func (db *dbFiles) set(name string, f *rootio.File) {
	db.Lock()
	defer db.Unlock()
	if old, dup := db.files[name]; dup {
		old.Close()
	}
	db.files[name] = f
}

func (db *dbFiles) del(name string) {
	db.Lock()
	defer db.Unlock()
	f, ok := db.files[name]
	if !ok {
		return
	}
	f.Close()
	delete(db.files, name)
}

type jsNode struct {
	ID    string `json:"id,omitempty"`
	Text  string `json:"text,omitempty"`
	Icon  string `json:"icon,omitempty"`
	State struct {
		Opened   bool `json:"opened,omitempty"`
		Disabled bool `json:"disabled,omitempty"`
		Selected bool `json:"selected,omitempty"`
	} `json:"state,omitempty"`
	Children []jsNode `json:"children,omitempty"`
	LiAttr   jsAttr   `json:"li_attr,omitempty"`
	Attr     jsAttr   `json:"a_attr,omitempty"`
}

type jsAttr map[string]interface{}

type brancher interface {
	Branches() []rootio.Branch
}

type jsNodes []jsNode

func (p jsNodes) Len() int           { return len(p) }
func (p jsNodes) Less(i, j int) bool { return p[i].ID < p[j].ID }
func (p jsNodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func newJsNodes(bres brancher, id string) ([]jsNode, error) {
	var err error
	branches := bres.Branches()
	if len(branches) <= 0 {
		return nil, err
	}
	var nodes []jsNode
	for _, b := range branches {
		id := strings.Join([]string{id, b.Name()}, "/")
		node := jsNode{
			ID:   id,
			Text: b.Name(),
			Icon: "fa fa-leaf",
			Attr: attrFor(b.(rootio.Object), id),
		}
		node.Children, err = newJsNodes(b, node.ID)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func fileJsTree(f *rootio.File, fname string) ([]jsNode, error) {
	root := jsNode{
		ID:   fname,
		Text: fmt.Sprintf("%s (version=%v)", fname, f.Version()),
		Icon: "fa fa-file",
	}
	root.State.Opened = true
	return dirTree(f, fname, root)
}

func dirTree(dir rootio.Directory, path string, root jsNode) ([]jsNode, error) {
	var nodes []jsNode
	for _, k := range dir.Keys() {
		obj, err := k.Object()
		if err != nil {
			return nil, fmt.Errorf("failed to extract key %q: %v", k.Name(), err)
		}
		switch obj := obj.(type) {
		case rootio.Tree:
			tree := obj
			node := jsNode{
				ID:   strings.Join([]string{path, k.Name()}, "/"),
				Text: fmt.Sprintf("%s (entries=%d)", k.Name(), tree.Entries()),
				Icon: "fa fa-tree",
			}
			node.Children, err = newJsNodes(tree, node.ID)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		case rootio.Directory:
			dir := obj
			node := jsNode{
				ID:   strings.Join([]string{path, k.Name()}, "/"),
				Text: k.Name(),
				Icon: "fa fa-folder",
			}
			node.Children, err = dirTree(dir, path+"/"+k.Name(), node)
			if err != nil {
				return nil, err
			}
			node.Children = node.Children[0].Children
			nodes = append(nodes, node)
		default:
			id := strings.Join([]string{path, k.Name() + fmt.Sprintf(";%d", k.Cycle())}, "/")
			nodes = append(nodes, jsNode{
				ID:   id,
				Text: fmt.Sprintf("%s;%d", k.Name(), k.Cycle()),
				Icon: iconFor(obj),
				Attr: attrFor(obj, id),
			})

		}
	}
	root.Children = nodes
	return []jsNode{root}, nil
}

func iconFor(obj rootio.Object) string {
	cls := obj.Class()
	switch {
	case strings.HasPrefix(cls, "TH1"):
		return "fa fa-bar-chart-o"
	case strings.HasPrefix(cls, "TH2"):
		return "fa fa-bar-chart-o"
	}
	return "fa fa-cube"
}

func attrFor(obj rootio.Object, id string) jsAttr {
	cls := obj.Class()
	switch {
	case strings.HasPrefix(cls, "TH1"):
		return jsAttr{
			"plot": true,
			"href": "/plot-1d/" + id,
		}
	case strings.HasPrefix(cls, "TH2"):
		return jsAttr{
			"plot": true,
			"href": "/plot-2d/" + id,
		}
	case strings.HasPrefix(cls, "TBranch"):
		return jsAttr{
			"plot": true,
			"href": "/plot-branch/" + id,
		}
	}
	return nil
}

func renderSVG(p *hplot.Plot) ([]byte, error) {
	size := 20 * vg.Centimeter
	canvas := vgsvg.New(size, size/vg.Length(math.Phi))
	p.Draw(draw.New(canvas))
	out := new(bytes.Buffer)
	_, err := canvas.WriteTo(out)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
