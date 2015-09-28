package linux

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
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

func (p *Plugin) validateConfigSpecs(path string) bool {
	valid := validation.Validation{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false
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
			p.errorLog = append(p.errorLog, "Atleast one Process.Env is empty")
			break
		}
	}
	if result := valid.Required(p.config.Root.Path, "Root.Path"); !result.Ok {
		p.errorLog = append(p.errorLog, "Root.Path is empty")
	}
	//Iterate over Mount array
	for _, mount := range p.config.Mounts {
		//If Mount points defined, it must define these three.
		if result := valid.Required(mount.Name, "Mount.Name"); !result.Ok {
			p.errorLog = append(p.errorLog, "Atleast one Mount.Name is empty")
			break
		}
		if result := valid.Required(mount.Path, "Mount.Path"); !result.Ok {
			p.errorLog = append(p.errorLog, "Atleast one Mount.Path is empty")
			break
		}
	}

	return true
}

func (p *Plugin) validateRuntimeSpecs(path string) bool {
	valid := validation.Validation{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	json.Unmarshal(data, &p.runtime)

	//Iterate over Mount array
	for _, mount := range p.runtime.Mounts {
		//If Mount points defined, it must define these three.
		if result := valid.Required(mount.Type, "Mount.Type"); !result.Ok {
			p.errorLog = append(p.errorLog, "Atleast one Mount.Type is empty")
			break
		}
		if result := valid.Required(mount.Source, "Mount.Source"); !result.Ok {
			p.errorLog = append(p.errorLog, "Atleast one Mount.Source is empty")
			break
		}
	}
	// Hooks Prestart
	for _, hook := range p.runtime.Hooks.Prestart {
		if result := valid.Required(hook.Path, "Hooks.Path"); !result.Ok {
			p.errorLog = append(p.errorLog, "Prestart hook Path cannot be empty")
			break
		}
	}

	// Hooks Poststop
	for _, hook := range p.runtime.Hooks.Poststop {
		if result := valid.Required(hook.Path, "Hooks.Path"); !result.Ok {
			p.errorLog = append(p.errorLog, "Poststop hook Path cannot be empty")
			break
		}
	}

	// UIDMappings mapping check.
	for _, uid := range p.runtime.Linux.UIDMappings {
		if result := valid.Range(uid.HostID, 0, 2147483647, "IDMapping.HostID"); !result.Ok {
			p.errorLog = append(p.errorLog, "UIDMapping's HostID must be valid integer")
			break
		}
		if result := valid.Range(uid.ContainerID, 0, 2147483647, "IDMapping.ContainerID"); !result.Ok {
			p.errorLog = append(p.errorLog, "UIDMapping's ContainerID must be valid integer")
			break
		}
		if result := valid.Range(uid.Size, 0, 2147483647, "IDMapping.Size"); !result.Ok {
			p.errorLog = append(p.errorLog, "UIDMapping's Size must be valid integer")
			break
		}
	}

	// GIDMappings mapping check.
	for _, gid := range p.runtime.Linux.GIDMappings {
		if result := valid.Range(gid.HostID, 0, 2147483647, "IDMapping.HostID"); !result.Ok {
			p.errorLog = append(p.errorLog, "GIDMapping's HostID must be valid integer")
			break
		}
		if result := valid.Range(gid.ContainerID, 0, 2147483647, "IDMapping.ContainerID"); !result.Ok {
			p.errorLog = append(p.errorLog, "GIDMapping's ContainerID must be valid integer")
			break
		}
		if result := valid.Range(gid.Size, 0, 2147483647, "IDMapping.Size"); !result.Ok {
			p.errorLog = append(p.errorLog, "GIDMapping's Size must be valid integer")
			break
		}
	}

	return true
}

func (p Plugin) ValidatePluginSpecs(path string) ([]string, bool) {

	validOCI := true

	if !p.validateConfigSpecs(path + "/config.json") {
		validOCI = false
		log.Errorf("Unable to Validate config.json")
	}
	if !p.validateRuntimeSpecs(path + "/runtime.json") {
		validOCI = false
		log.Errorf("Unable to Validate runtime.json")
	}

	if len(p.errorLog) > 0 {
		validOCI = false
		//p.errorLog = append(p.errorLog, "NOTE: Some errors may appear due to invalid OCI format")
		log.Info("NOTE: Some errors may appear due to invalid OCI format")
	}

	return p.errorLog, validOCI
}

//FIXME: Still runc has not implemented the changes, so state.json
//	 file has diffrent structure, so cannot verify.
//	Implementtion incomplete.
func (p Plugin) ValidatePluginRuntimeSpecs(containerID string) ([]string, bool) {
	//path := "./runtime.json"
	//path := "/run/oci" + "/" + containerID + "/state.json"

	validOCIStatus := true

	if len(p.errorLog) > 0 {
		validOCIStatus = false
		p.errorLog = append(p.errorLog, "NOTE: Some errors may appear due to invalid OCI format")
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
