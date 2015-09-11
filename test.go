package main

import (
	"fmt"
	"github.com/kunalkushwaha/octool/plugins"
	_ "github.com/kunalkushwaha/octool/plugins/none"

)

func main() {
	plugin, err := plugin.NewPlugin("none")
	if err!= nil {
		fmt.Println(err)
		return
	}
	plugin.ValidatePluginSpecs()
	plugin.Analyze()
	plugin.TestExecution()

}
