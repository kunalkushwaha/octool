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
	errorLog []string
}

func init() {
	plugin.Register("linux",
		&plugin.RegisteredPlugin{New: NewPlugin})
}

func NewPlugin(pluginName string, path string) (plugin.Plugin, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Erorr in reading file!!")
		return Plugin{}, err
	}

	plugin := Plugin{}
	json.Unmarshal(data, &plugin.config)
	return plugin, nil
}

func (p Plugin) ValidatePluginSpecs() ([]string, bool) {

	validOCI := true
	valid := validation.Validation{}

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
	///	valid.Required(p.config.Process.User.UID, "User.UID")
	//	valid.Required(p.config.Process.User.GID, "User.GID")
	if result := valid.Required(p.config.Root.Path, "Root.Path"); !result.Ok {
		p.errorLog = append(p.errorLog, "Root.Path is empty")
	}
	//Iterate over Mount array
	//	for _, mount := range p.config.Mounts {
	//If Mount points defined, it must define these three.
	//valid.Required(mount.Type, "Mount.Type")
	//		valid.Required(mount.Source, "Mount.Source")
	//		valid.Required(mount.Destination, "Mount.Destination")
	//	}
	//fmt.Println(errorLog)
	if len(p.errorLog) > 0 {
		validOCI = false
	}
	return p.errorLog, validOCI
}

func (p Plugin) Analyze() []string {
	fmt.Println("none: Analyze() ")
	return []string{"none", "two"}
}

func (p Plugin) TestExecution() []string {
	fmt.Println("none: TestExecution() ")
	return []string{"none", "three"}
}

func dumpJSON(config Plugin) {
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
