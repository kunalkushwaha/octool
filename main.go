package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/kunalkushwaha/octool/plugins"
	_ "github.com/kunalkushwaha/octool/plugins/linux"
)

func main() {
	log.SetLevel(log.InfoLevel)
	app := cli.NewApp()
	app.Name = "octool"
	app.Usage = "Toolchain for OpenContainer Format"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:  "lint",
			Usage: "validate container config file(s)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "image",
					Value: ".",
					Usage: "path of image to validate",
				},
				cli.StringFlag{
					Name:  "os",
					Value: "linux",
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
	imagePath := c.String("image")
	targetOS := c.String("os")

	_, err := os.Stat(imagePath)
	if os.IsNotExist(err) {
		cli.ShowCommandHelp(c, "lint")
		return
	}
	//FIXME: Instead of default as linux, detect os

	plugin, err := plugin.NewPlugin(targetOS)
	if err != nil {
		//fmt.Println(err)
		log.Error(err)
		return
	}
	errors, valid := plugin.ValidatePluginSpecs(imagePath)
	if !valid {
		fmt.Println("")
		for _, err := range errors {
			log.Warn(err)
			//fmt.Println(err)
		}
		fmt.Printf("\nInvalid OCI config format\n")
	} else {
		fmt.Printf("\nConfig is Valid OCI\n")
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
		log.Error(err)
		//fmt.Println(err)
		return
	}
	errors, valid := plugin.ValidatePluginRuntimeSpecs(containerID)
	if !valid {
		for _, err := range errors {
			//fmt.Println(err)
			log.Warn(err)
		}
		fmt.Printf("\nInvalid OCI runtime format\n")
	} else {
		fmt.Println("Container State Valid OCI")
	}

	return
}
