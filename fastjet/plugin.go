// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fastjet

import (
	"golang.org/x/xerrors"
)

type Plugin interface {
	Description() string
	RunClustering(builder Builder) error
	R() float64
}

var (
	g_plugins = make(map[string]Plugin)
)

func Register(name string, plugin Plugin) {
	_, dup := g_plugins[name]
	if dup {
		panic(xerrors.Errorf("fastjet.Register: duplicate plugin [%s] (%s)", name, plugin.Description()))
	}
	g_plugins[name] = plugin
}

func GetPlugin(name string) (Plugin, error) {
	plugin, ok := g_plugins[name]
	if !ok {
		return nil, xerrors.Errorf("fastjet.JetPlugin: no such plugin [%s]", name)
	}

	return plugin, nil
}
