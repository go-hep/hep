// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	stdpath "path"
	"strings"

	"gonum.org/v1/plot/vg"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rsrv"
	"go-hep.org/x/hep/groot/rtree"
)

const (
	plotH1     = "h1"
	plotH2     = "h2"
	plotS2     = "s2"
	plotBranch = "branch"
)

type plot struct {
	Type string   `json:"type"`
	URI  string   `json:"uri"`
	Dir  string   `json:"dir"`
	Obj  string   `json:"obj"`
	Vars []string `json:"vars"`

	Options rsrv.PlotOptions `json:"options"`
}

type plotRequest struct {
	cookie *http.Cookie
	req    plot
	resp   chan plotResponse
}

type plotResponse struct {
	err    error
	ctype  string
	status int
	body   []byte
}

type jsNode struct {
	ID    string `json:"id,omitempty"`
	URI   string `json:"uri,omitempty"`
	Dir   string `json:"dir,omitempty"`
	Obj   string `json:"obj,omitempty"`
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

type jsAttr map[string]any

type brancher interface {
	Branches() []rtree.Branch
}

type jsNodes []jsNode

func (p jsNodes) Len() int           { return len(p) }
func (p jsNodes) Less(i, j int) bool { return p[i].ID < p[j].ID }
func (p jsNodes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func newJsNodes(bres brancher, parent jsNode) ([]jsNode, error) {
	var err error
	branches := bres.Branches()
	if len(branches) <= 0 {
		return nil, err
	}
	var nodes []jsNode
	for _, b := range branches {
		id := parent.ID
		bid := strings.Join([]string{id, b.Name()}, "/")
		node := jsNode{
			ID:   bid,
			URI:  parent.URI,
			Dir:  stdpath.Join(parent.Dir, parent.Obj),
			Obj:  b.Name(),
			Text: b.Name(),
			Icon: "fa fa-leaf",
		}
		node.Attr, err = attrFor(b.(root.Object), node)
		if err != nil {
			return nil, err
		}
		node.Children, err = newJsNodes(b, node)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func fileJsTree(f *riofs.File, fname string) ([]jsNode, error) {
	root := jsNode{
		ID:   f.Name(),
		URI:  fname,
		Dir:  "/",
		Text: fmt.Sprintf("%s (version=%v)", fname, f.Version()),
		Icon: "fa fa-file",
	}
	root.State.Opened = true
	return dirTree(f, fname, root)
}

func dirTree(dir riofs.Directory, path string, root jsNode) ([]jsNode, error) {
	var nodes []jsNode
	for _, k := range dir.Keys() {
		obj, err := k.Object()
		if err != nil {
			return nil, fmt.Errorf("failed to extract key %q: %w", k.Name(), err)
		}
		switch obj := obj.(type) {
		case rtree.Tree:
			tree := obj
			node := jsNode{
				ID:   strings.Join([]string{path, k.Name()}, "/"),
				URI:  root.URI,
				Dir:  stdpath.Join(root.Dir, root.Obj),
				Obj:  k.Name(),
				Text: fmt.Sprintf("%s (entries=%d)", k.Name(), tree.Entries()),
				Icon: "fa fa-tree",
			}
			node.Children, err = newJsNodes(tree, node)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		case riofs.Directory:
			dir := obj
			node := jsNode{
				ID:   strings.Join([]string{path, k.Name()}, "/"),
				URI:  root.URI,
				Dir:  stdpath.Join(root.Dir, root.Obj),
				Obj:  k.Name(),
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
			node := jsNode{
				ID:   id,
				URI:  root.URI,
				Dir:  stdpath.Join(root.Dir, root.Obj),
				Obj:  k.Name(),
				Text: fmt.Sprintf("%s;%d", k.Name(), k.Cycle()),
				Icon: iconFor(obj),
			}
			node.Attr, err = attrFor(obj, node)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, node)
		}
	}
	root.Children = nodes
	return []jsNode{root}, nil
}

func iconFor(obj root.Object) string {
	cls := obj.Class()
	switch {
	case strings.HasPrefix(cls, "TH1"):
		return "fa fa-bar-chart-o"
	case strings.HasPrefix(cls, "TH2"):
		return "fa fa-bar-chart-o"
	case strings.HasPrefix(cls, "TGraph"):
		return "fa fa-bar-chart-o"
	}
	return "fa fa-cube"
}

func attrFor(obj root.Object, node jsNode) (jsAttr, error) {
	cmd := new(bytes.Buffer)
	cls := obj.Class()
	switch {
	case strings.HasPrefix(cls, "TH1"):
		req := plot{
			Type: plotH1,
			URI:  node.URI,
			Dir:  node.Dir,
			Obj:  node.Obj,
			Options: rsrv.PlotOptions{
				Title:  node.Obj,
				Type:   "svg",
				Height: -1,
				Width:  20 * vg.Centimeter,
			},
		}
		err := json.NewEncoder(cmd).Encode(req)
		if err != nil {
			return nil, err
		}
		return jsAttr{
			"plot": true,
			"href": "/plot",
			"cmd":  cmd.String(),
		}, nil
	case strings.HasPrefix(cls, "TH2"):
		req := plot{
			Type: plotH2,
			URI:  node.URI,
			Dir:  node.Dir,
			Obj:  node.Obj,
			Options: rsrv.PlotOptions{
				Title:  node.Obj,
				Type:   "svg",
				Height: -1,
				Width:  20 * vg.Centimeter,
			},
		}
		err := json.NewEncoder(cmd).Encode(req)
		if err != nil {
			return nil, err
		}
		return jsAttr{
			"plot": true,
			"href": "/plot",
			"cmd":  cmd.String(),
		}, nil
	case strings.HasPrefix(cls, "TGraph"):
		req := plot{
			Type: plotS2,
			URI:  node.URI,
			Dir:  node.Dir,
			Obj:  node.Obj,
			Options: rsrv.PlotOptions{
				Title:  node.Obj,
				Type:   "svg",
				Height: -1,
				Width:  20 * vg.Centimeter,
			},
		}
		err := json.NewEncoder(cmd).Encode(req)
		if err != nil {
			return nil, err
		}
		return jsAttr{
			"plot": true,
			"href": "/plot",
			"cmd":  cmd.String(),
		}, nil
	case strings.HasPrefix(cls, "TBranch"):
		req := plot{
			Type: plotBranch,
			URI:  node.URI,
			Dir:  stdpath.Dir(node.Dir),
			Obj:  stdpath.Base(node.Dir),
			Vars: []string{node.Obj},
			Options: rsrv.PlotOptions{
				Title:  node.Obj,
				Type:   "svg",
				Height: -1,
				Width:  20 * vg.Centimeter,
			},
		}
		err := json.NewEncoder(cmd).Encode(req)
		if err != nil {
			return nil, err
		}
		return jsAttr{
			"plot": true,
			"href": "/plot",
			"cmd":  cmd.String(),
		}, nil
	}
	// TODO(sbinet) do something clever with things we don't know how to handle?
	return nil, nil
}
