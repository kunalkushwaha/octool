# octool
OpenContainer tool set for image validation, analysis and test.

#### Current Status
In Development.
##### Feature wise status
- import : Docker Registry supported
- spec   : Completed (Docker Manifest to OCI's confg.json)
- lint   : Completed.
- validate-state   : Pending for runc implementation.

#### Usage

``go get github.com/kunalkushwaha/octool``

```
$ octool
Toolchain for OpenContainer Format(OCI)

Usage:
  octool [command]

Available Commands:
  import         Imports container image from remote registery and convert it to runc's rootfs
  lint           validate container config file(s)
  spec           genrates runc compatible spec from manifest file
  validate-state Validates the Container state

Flags:
  -h, --help   help for octool

Use "octool [command] --help" for more information about a command.



```

##### example
```
octool import docker://kunalkushwaha/demoapp_image:v1 -t demoapp
rootfs is prepared at : demoapp/rootfs

$ cd demoapp

$ octool spec

Succesfully generated config.json

$ octool lint
INFO[0000] NOTE: Some errors may appear due to invalid OCI format

WARN[0000] Version must be in format of X.X.X (complient to Semver v2.0.0)

Invalid OCI config format

```
