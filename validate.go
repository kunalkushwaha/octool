package main

import (
	"fmt"
	"os"
	"encoding/json"

	"github.com/codegangsta/cli"
	"github.com/opencontainers/specs"
)

func validateOCImage(c *cli.Context) {
	var config	specs.LinuxSpec
	configJson := c.String("json")

	if len(configJson) > 0 {
		// Read json file and load into spec.LinuxSpec struct.
		// Validate as per rules and spec defination.
		configFile, err := os.Open(configJson)
		if err != nil {
			fmt.Println("opening config file", err.Error())
			return
		}

		jsonParser := json.NewDecoder(configFile)
		if err = jsonParser.Decode(&config); err != nil {
			fmt.Println(configJson, "is not valid json file")
			return
		}
		fmt.Println(config)
		return
	}
	return
}

func testOContainer(c *cli.Context) {
	fmt.Println("NOT-IMPLEMENTED")
	return
}
