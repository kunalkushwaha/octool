# octool
OpenContainer tool set for image validation, analysis and test.

#### Current Status
In Development.
##### Feature wise status
- lint  : Completed.
- validate-state   : Pending for runc implementation.

#### Usage 

``go get github.com/kunalkushwaha/octool``

```
$ octool
NAME:
   octool - Toolchain for OpenContainer Format

USAGE:
   octool [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   lint                 validate container config file
   validate-state       Validates the Container state
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version


```

##### example 
```
$  octool lint

OR

$ octool lint --image /home/test/container/  --os linux

```
