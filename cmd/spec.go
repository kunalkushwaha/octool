package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/user"
	"runtime"
	"strconv"

	"github.com/docker/distribution/manifest/schema1"
	ctr "github.com/docker/docker/api/types/container"
	"github.com/opencontainers/runtime-spec/specs-go"
	"github.com/spf13/cobra"
)

// specCmd represents the spec command
var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "genrates runc compatible spec from manifest file",
	Run:   generateSpec,
}

type ImageInspect struct {
	Architecture    string
	Author          string
	Config          *ctr.Config
	Container       string
	ContainerConfig *ctr.Config
	DockerVersion   string
	Created         string
	ID              string
	Os              string
	Parent          string
	Throwway        bool
}

var (
	spec = specs.Spec{
		Version: specs.Version,
		Platform: specs.Platform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Root: specs.Root{
			Path:     "rootfs",
			Readonly: true,
		},
		Process: specs.Process{
			Terminal: true,
			User:     specs.User{},
			Args: []string{
				"sh",
			},
			Env: []string{
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
				"TERM=xterm",
			},
			Cwd:             "/",
			NoNewPrivileges: true,
			Capabilities: []string{
				"CAP_NET_RAW",
				"CAP_NET_BIND_SERVICE",
				"CAP_AUDIT_READ",
				"CAP_AUDIT_WRITE",
				"CAP_DAC_OVERRIDE",
				"CAP_SETFCAP",
				"CAP_SETPCAP",
				"CAP_SETGID",
				"CAP_SETUID",
				"CAP_MKNOD",
				"CAP_CHOWN",
				"CAP_FOWNER",
				"CAP_FSETID",
				"CAP_KILL",
				"CAP_SYS_CHROOT",
			},
			Rlimits: []specs.Rlimit{
				{
					Type: "RLIMIT_NOFILE",
					Hard: uint64(1024),
					Soft: uint64(1024),
				},
			},
		},
		Hostname: "runc",
		Mounts: []specs.Mount{
			{
				Destination: "/proc",
				Type:        "proc",
				Source:      "proc",
				Options:     []string{"nosuid", "noexec", "nodev"},
			},
			{
				Destination: "/dev",
				Type:        "tmpfs",
				Source:      "tmpfs",
				Options:     []string{"nosuid", "strictatime", "mode=755"},
			},
			{
				Destination: "/dev/pts",
				Type:        "devpts",
				Source:      "devpts",
				Options:     []string{"nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5"},
			},
			{
				Destination: "/dev/shm",
				Type:        "tmpfs",
				Source:      "shm",
				Options:     []string{"nosuid", "noexec", "nodev", "mode=1777", "size=65536k"},
			},
			{
				Destination: "/dev/mqueue",
				Type:        "mqueue",
				Source:      "mqueue",
				Options:     []string{"nosuid", "noexec", "nodev"},
			},
			{
				Destination: "/sys",
				Type:        "sysfs",
				Source:      "sysfs",
				Options:     []string{"nosuid", "noexec", "nodev", "ro"},
			},
			{
				Destination: "/sys/fs/cgroup",
				Type:        "cgroup",
				Source:      "cgroup",
				Options:     []string{"nosuid", "noexec", "nodev", "relatime", "ro"},
			},
		},
		Linux: &specs.Linux{
			MaskedPaths: []string{
				"/proc/kcore",
				"/proc/latency_stats",
				"/proc/timer_list",
				"/proc/timer_stats",
				"/proc/sched_debug",
				//FIXME: Fails to start the container saying "not a directory"
				//"/sys/firmware",
			},
			ReadonlyPaths: []string{
				"/proc/asound",
				"/proc/bus",
				"/proc/fs",
				"/proc/irq",
				"/proc/sys",
				"/proc/sysrq-trigger",
			},
			Resources: &specs.Resources{
				Devices: []specs.DeviceCgroup{
					{
						Allow:  false,
						Access: sPtr("rwm"),
					},
				},
			},
			Namespaces: []specs.Namespace{
				{
					Type: "pid",
				},
				{
					Type: "network",
				},
				{
					Type: "ipc",
				},
				{
					Type: "uts",
				},
				{
					Type: "mount",
				},
			},
		},
	}
)

func init() {
	RootCmd.AddCommand(specCmd)

	specCmd.Flags().StringP("runas", "", "", "user as which container should run ")
	specCmd.Flags().StringP("security-profile", "", "", "Apparmor security profile for container")
	//specCmd.Flags().StringP("pre-start", "", "", "pre-start hook script/binary to execute before container start")
	//specCmd.Flags().StringP("post-start", "", "", "post-start hook script/binary to execute after container start")
	//specCmd.Flags().StringP("post-stop", "", "", "post-stop hook script/binary to execute after container stops")

}

func generateSpec(cmd *cobra.Command, args []string) {

	//fmt.Println(spec)
	manifestFileData, err := ioutil.ReadFile("manifest.json")
	if err != nil {
		fmt.Println("Error while reading manifest file.")
		return
	}
	testManifest := schema1.SignedManifest{}
	json.Unmarshal(manifestFileData, &testManifest)
	//fmt.Println(testManifest.History[0].V1Compatibility)
	configJSON := ImageInspect{}
	err = json.Unmarshal([]byte(testManifest.History[0].V1Compatibility), &configJSON)
	if err != nil {
		fmt.Println(err)
	}

	spec.Platform.Arch = configJSON.Architecture
	spec.Platform.OS = configJSON.Os
	//Append xterm
	spec.Process.Env = configJSON.Config.Env
	if configJSON.Config.WorkingDir != "" {
		spec.Process.Cwd = configJSON.Config.WorkingDir
	}

	//spec.Process.Args = configJSON.Config.Cmd

	//TODO: take user input or assin current user and group.
	// NOTE: This will run the container with current user in container
	runasUser, _ := cmd.Flags().GetString("runas")
	if runasUser != "" {
		u, err := user.Lookup(runasUser)
		if err != nil {
			fmt.Printf("%s user not found in system\n", runasUser)
			return
		}
		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)
		//	fmt.Println(uid, gid)
		spec.Process.User.UID = uint32(uid)
		spec.Process.User.GID = uint32(gid)
	}

	securityProfile, _ := cmd.Flags().GetString("security-profile")
	spec.Process.ApparmorProfile = securityProfile

	//TODO: Capablities?
	//TODO: Parse /dev/ram* (tmpfs) and setup devices in cgroup list
	//TODO: Can we use the jessfraz/netns for setting up network?
	jsonSpec, _ := json.MarshalIndent(spec, "", "\t")
	err = ioutil.WriteFile("config.json", jsonSpec, 0666)
	if err != nil {
		fmt.Println("Error while writing config file : ", err)
	}
	fmt.Printf("\nSuccesfully generated config.json\n")
}

func sPtr(s string) *string { return &s }
