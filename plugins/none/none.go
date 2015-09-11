package none

import (
	"fmt"
	"github.com/kunalkushwaha/octool/plugins"
)

type Plugin struct {
}


func init() {
	fmt.Println("init().. none invoked")
	plugin.Register("none",
	&plugin.RegisteredPlugin{New: NewPlugin})
}

func NewPlugin(pluginName string) (plugin.Plugin, error) {
	return &Plugin{}, nil
}

func (p *Plugin) ValidatePluginSpecs() []string {
	fmt.Println("none: ValidatePluginSpecs() ")
	return []string{"none","one"}
}

func (p *Plugin) Analyze() []string {
	fmt.Println("none: Analyze() ")
	return []string{"none","two"}
}

func (p *Plugin) TestExecution() []string {
	fmt.Println("none: TestExecution() ")
	return []string{"none","three"}
}

