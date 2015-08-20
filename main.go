package main

import (
	"github.com/codegangsta/cli"
	"os"
)

func main() {

	app := cli.NewApp()
	app.Name = "octool"
	app.Usage = "Toolchain for OpenContainer Specs"
	app.Version = "0.1.0"
	app.Commands = []cli.Command{
		{
			Name:  "validate",
			Usage: "validate container image / Json",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "json",
					Usage: "[Optional] dry_run the command",
				},
			},
			Action: validateOCImage,
		},
		{
			Name:   "test",
			Usage:  "Test the Container",
			Action: testOContainer,
		},
	}

	app.Run(os.Args)
}

