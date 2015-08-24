package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	"github.com/astaxie/beego/validation"
	"github.com/codegangsta/cli"
	"github.com/opencontainers/specs"
)

type Config struct {
	configLinux specs.LinuxSpec
	//	comfigWindows	spec.WindowsSpec
}

func validateOCImage(c *cli.Context) {
	configJson := c.String("json")

	if len(configJson) == 0 {
		cli.ShowCommandHelp(c, "validate")
		return
	}

	// Read json file and load into spec.LinuxSpec struct.
	// Validate as per rules and spec defination.
	config, err := NewConfig(configJson)
	if err != nil {
		fmt.Println("Error while opening File ", err)
		return
	}
	config.ValidateCommonSpecs()
	//fmt.Println(config)
	//	dumpJSON(config)
	return

}

func dumpJSON(config Config) {
	b, err := json.Marshal(config.configLinux)
	if err != nil {
		fmt.Println(err)
		return
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)
	//	fmt.Println(out)
}

func getOS() string {
	return runtime.GOOS
}

func NewConfig(path string) (Config, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	config := Config{}
	if getOS() == "linux" {
		json.Unmarshal(data, &config.configLinux)
	}
	//	if getOS() == "windows" {
	//		json.Unmarshal(data, &config.configWindows)
	//	}

	return config, nil
}

func (conf *Config) ValidateCommonSpecs() bool {
	valid := validation.Validation{}
	valid.Required(conf.configLinux.Version, "Version")
	valid.Required(conf.configLinux.Platform.OS, "OS")
	valid.Required(conf.configLinux.Platform.Arch, "Platform.Arch")
	valid.Required(conf.configLinux.Process.User.UID, "User.UID")
	valid.Required(conf.configLinux.Process.User.GID, "User.GID")
	valid.Required(conf.configLinux.Root.Path, "Root.Path")
	//Iterate over Mount array
	//valid.Required(conf.configLinux.Mounts.Type, "Mount.Type")
	//valid.Required(conf.configLinux.Mounts.Source, "Mount.Source")
	//valid.Required(conf.configLinux.Mounts.Destination, "Mount.Destination")

	if valid.HasErrors() {
		// validation does not pass
		// print invalid message
		for _, err := range valid.Errors {
			fmt.Println(err.Key, err.Message)
		}
		fmt.Println("\nNOTE: Try to fix errors from top\n",
			"     as intital errors case parsing issue, resulting too many erorrs")

	}

	return false
}

func (conf *Config) ValidateLinuxSpecs() bool {
	return false
}

func (conf *Config) Analyze() {

}

func testOContainer(c *cli.Context) {
	fmt.Println("NOT-IMPLEMENTED")
	return
}
