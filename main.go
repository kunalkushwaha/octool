package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/kunalkushwaha/octool/plugins"
	_ "github.com/kunalkushwaha/octool/plugins/linux"
)

func main() {

	app := cli.NewApp()
	app.Name = "octool"
	app.Usage = "Toolchain for OpenContainer Format"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:  "lint",
			Usage: "validate container config file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config",
					Usage: "config file to validate [If not present, expects config.json in current folder]",
				},
				cli.StringFlag{
					Name:  "os",
					Usage: "Target OS",
				},
			},
			Action: validateContainerConfig,
		},
		{
			Name:  "validate-state",
			Usage: "Validates the Container state",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "id",
					Usage: "Container ID",
				},
			},
			Action: validateContainerState,
		},
	}

	app.Run(os.Args)
}

func validateContainerConfig(c *cli.Context) {
	configJson := c.String("config")
	targetOS := c.String("os")

	if len(configJson) == 0 {
		configJson = "config.json"
		_, err := os.Stat(configJson)
		if os.IsNotExist(err) {
			cli.ShowCommandHelp(c, "lint")
			return
		}
	}

	if len(targetOS) == 0 {
		//FIXME: Instead of default as linux, detect os
		targetOS = "linux"
	}

	plugin, err := plugin.NewPlugin(targetOS)
	if err != nil {
		fmt.Println(err)
		return
	}
	errors, valid := plugin.ValidatePluginSpecs(configJson)
	if !valid {
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Config is Valid OCI")
	}
	return

}

func validateContainerState(c *cli.Context) {
	containerID := c.String("id")

	if len(containerID) == 0 {
		cli.ShowCommandHelp(c, "validate-state")
		return
	}

	//FIXME: Instead of default as linux, detect os
	targetOS := "linux"

	plugin, err := plugin.NewPlugin(targetOS)
	if err != nil {
		fmt.Println(err)
		return
	}
	errors, valid := plugin.ValidatePluginRuntimeSpecs(containerID)
	if !valid {
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Container State Valid OCI")
	}

	return
}
