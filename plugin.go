package fastjet

import (
	"fmt"
)

type Plugin interface {
	Description() string
	RunClustering(builder ClusterBuilder) error
	R() float64
}

var (
	g_plugins = make(map[string]Plugin)
)

func Register(name string, plugin Plugin) {
	_, dup := g_plugins[name]
	if dup {
		panic(fmt.Errorf("fastjet.Register: duplicate plugin [%s] (%s)", name, plugin.Description()))
	}
	g_plugins[name] = plugin
}

func GetPlugin(name string) (Plugin, error) {
	plugin, ok := g_plugins[name]
	if !ok {
		return nil, fmt.Errorf("fastjet.JetPlugin: no such plugin [%s]", name)
	}

	return plugin, nil
}
