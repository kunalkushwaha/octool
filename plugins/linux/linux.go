package linux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/astaxie/beego/validation"
	"github.com/kunalkushwaha/octool/plugins"
	"github.com/opencontainers/specs"
)

type Plugin struct {
	config   specs.LinuxSpec
	runtime  specs.LinuxRuntimeSpec
	errorLog []string
}

func init() {
	plugin.Register("linux",
		&plugin.RegisteredPlugin{New: NewPlugin})
}

func NewPlugin(pluginName string) (plugin.Plugin, error) {

	return Plugin{}, nil
}

func (p Plugin) ValidatePluginSpecs(path string) ([]string, bool) {

	validOCI := true
	valid := validation.Validation{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Erorr in reading file!!")
		return p.errorLog, false
	}

	json.Unmarshal(data, &p.config)

	//Validate mandatory fields.
	if result := valid.Required(p.config.Version, "Version"); !result.Ok {
		p.errorLog = append(p.errorLog, "Version cannot be empty")
	}
	//Version must complient with  SemVer v2.0.0
	if result := valid.Match(p.config.Version, regexp.MustCompile("^(\\d+\\.)?(\\d+\\.)?(\\*|\\d+)$"), "Version"); !result.Ok {
		p.errorLog = append(p.errorLog, "Version must be in format of X.X.X (complient to Semver v2.0.0)")
	}
	if result := valid.Required(p.config.Platform.OS, "OS"); !result.Ok {
		p.errorLog = append(p.errorLog, "OS can be not empty")
	}
	if result := valid.Required(p.config.Platform.Arch, "Platform.Arch"); !result.Ok {
		p.errorLog = append(p.errorLog, "Platform.Arch is empty")
	}

	for _, env := range p.config.Process.Env {
		//If Process defined, env cannot be empty
		if result := valid.Required(env, "Process.Env"); !result.Ok {
			p.errorLog = append(p.errorLog, "Process.Env is empty")
		}
	}
	if result := valid.Required(p.config.Root.Path, "Root.Path"); !result.Ok {
		p.errorLog = append(p.errorLog, "Root.Path is empty")
	}
	//Iterate over Mount array
	for _, mount := range p.config.Mounts {
		//If Mount points defined, it must define these three.
		if result := valid.Required(mount.Name, "Mount.Name"); !result.Ok {
			p.errorLog = append(p.errorLog, "Mount.Name is required")
		}
		if result := valid.Required(mount.Path, "Mount.Path"); !result.Ok {
			p.errorLog = append(p.errorLog, "Mount.Path is required")
		}
	}
	if len(p.errorLog) > 0 {
		validOCI = false
	}
	return p.errorLog, validOCI
}

//FIXME: Still runc has not implemented the changes, so state.json
//	 file has diffrent structure, so cannot verify.
//	Implementtion incomplete.
func (p Plugin) ValidatePluginRuntimeSpecs(containerID string) ([]string, bool) {
	path := specs.LinuxStateDirectory + "/" + containerID + "/state.json"
	//path := "/run/oci" + "/" + containerID + "/state.json"

	validOCIStatus := true
	valid := validation.Validation{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Erorr in reading file!!")
		return p.errorLog, false
	}

	json.Unmarshal(data, &p.runtime)

	//Iterate over Mount array
	for _, mount := range p.runtime.Mounts {
		//If Mount points defined, it must define these three.
		if result := valid.Required(mount.Type, "Mount.Type"); !result.Ok {
			p.errorLog = append(p.errorLog, "Mount.Type is empty")
		}
		if result := valid.Required(mount.Source, "Mount.Source"); !result.Ok {
			p.errorLog = append(p.errorLog, "Mount.Path is empty")
		}
	}

	if len(p.errorLog) > 0 {
		validOCIStatus = false
	}

	return p.errorLog, validOCIStatus
}

func (p Plugin) Analyze() []string {
	fmt.Println("none: Analyze() ")
	return p.errorLog
}

func (p Plugin) TestExecution() []string {
	fmt.Println("none: TestExecution() ")
	return p.errorLog
}

// Debugging functions.
func dumpConfig(config specs.LinuxSpec) {
	b, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Println("")
}

func dumpRuntimeConfig(config specs.LinuxRuntimeSpec) {
	b, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteTo(os.Stdout)
	fmt.Println("")
}
