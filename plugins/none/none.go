package none

import "github.com/kunalkushwaha/octool/plugins"

type Plugin struct {
}

func init() {
	plugin.Register("none",
		&plugin.RegisteredPlugin{New: NewPlugin})
}

func NewPlugin(pluginName string, path string) (plugin.Plugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) ValidatePluginSpecs() ([]string, bool) {
	return []string{"none", "one"}, true
}

func (p *Plugin) Analyze() []string {
	return []string{"none", "two"}
}

func (p *Plugin) TestExecution() []string {
	return []string{"none", "three"}
}
