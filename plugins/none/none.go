package none

import "github.com/kunalkushwaha/octool/plugins"

type Plugin struct {
}

func init() {
	plugin.Register("none",
		&plugin.RegisteredPlugin{New: NewPlugin})
}

func NewPlugin(pluginName string) (plugin.Plugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) ValidatePluginSpecs(path string) ([]string, bool) {
	return []string{}, true
}

func (p Plugin) ValidatePluginState(containerID string) ([]string, bool) {
	return []string{}, true
}

func (p *Plugin) Analyze() []string {
	return []string{}
}

func (p *Plugin) TestExecution() []string {
	return []string{}
}
