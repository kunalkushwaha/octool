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
			Name:  "validate",
			Usage: "validate container image / Json",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "json config file to validate",
				},
				cli.StringFlag{
					Name:  "os",
					Usage: "Target OS",
				},
			},
			Action: validateOCImage,
		},
		{
			Name:  "test",
			Usage: "Test the Container",
			//Action: testOContainer,
		},
	}

	app.Run(os.Args)
}

func validateOCImage(c *cli.Context) {
	configJson := c.String("json")
	//os := c.String("os")

	if len(configJson) == 0 {
		cli.ShowCommandHelp(c, "validate")
		return
	}
	plugin, err := plugin.NewPlugin("linux", "test.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	errors, valid := plugin.ValidatePluginSpecs()
	if !valid {
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Config is Valid OCI")
	}
	return

}
