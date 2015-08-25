package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"regexp"

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
	common := config.ValidateCommonSpecs()
	platform := config.ValidateLinuxSpecs()
	if !common || !platform {
		fmt.Println("\nNOTE: One or more errors found in", configJson)
	} else {
		fmt.Println("\n",configJson, "has Valid OC Format !!")
	}
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

	//Validate mandatory fields.
	valid.Required(conf.configLinux.Version, "Version")
	//Version must complient with  SemVer v2.0.0
	valid.Match(conf.configLinux.Version, regexp.MustCompile("^(\\d+\\.)?(\\d+\\.)?(\\*|\\d+)$"),"Version")
	valid.Required(conf.configLinux.Platform.OS, "OS")
	valid.Required(conf.configLinux.Platform.Arch, "Platform.Arch")

	for _, env := range conf.configLinux.Process.Env {
		//If Process defined, env cannot be empty
		valid.Required(env, "Process.Env")
	}
	valid.Required(conf.configLinux.Process.User.UID, "User.UID")
	valid.Required(conf.configLinux.Process.User.GID, "User.GID")
	valid.Required(conf.configLinux.Root.Path, "Root.Path")
	//Iterate over Mount array
	for _, mount := range conf.configLinux.Mounts {
		//If Mount points defined, it must define these three.
		valid.Required(mount.Type, "Mount.Type")
		valid.Required(mount.Source, "Mount.Source")
		valid.Required(mount.Destination, "Mount.Destination")
	}

	if valid.HasErrors() {
		// validation does not pass
		for i, err := range valid.Errors {
			fmt.Println(i, err.Key, err.Message)
		}
		return false
	}

	return true
}

func (conf *Config) ValidateLinuxSpecs() bool {
	valid := validation.Validation{}

	for _, namespace := range conf.configLinux.Linux.Namespaces  {
		valid.Required(namespace.Type, "Namespace.Type")
	}


	if valid.HasErrors() {
		// validation does not pass
		fmt.Println("\nLinux Specific config errors\n")

		for i, err := range valid.Errors {
			fmt.Println(i, err.Key, err.Message)
		}
		return false
	}

	return true
}

func (conf *Config) Analyze() {
	fmt.Println("NOT-IMPLEMETED")
	return
}

func testOContainer(c *cli.Context) {
	fmt.Println("NOT-IMPLEMENTED")
	return
}
