package plugin

import (
	"fmt"
)

type Plugin interface {
	ValidatePluginSpecs() ([]string, bool)
	Analyze() []string
	TestExecution() []string
}

type RegisteredPlugin struct {
	New func(pluginName string, path string) (Plugin, error)
}

var (
	plugins map[string]*RegisteredPlugin
)

func init() {
	plugins = make(map[string]*RegisteredPlugin)
}

// Register a plugin
func Register(name string, registeredPlugin *RegisteredPlugin) error {

	if _, exists := plugins[name]; exists {
		return fmt.Errorf("Name already registered %s", name)
	}

	plugins[name] = registeredPlugin
	return nil
}

func NewPlugin(name string, path string) (Plugin, error) {
	plugin, exists := plugins[name]
	if !exists {
		return nil, fmt.Errorf("Plugin: Unknown plugin %q", name)
	}
	return plugin.New(name, path)
}
