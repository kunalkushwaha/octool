# octool
OpenContainer tool set for image validation, analysis and test.

#### Current Status
In Development.
##### Feature wise status
- Validate  : Almost completed.
- Analyse   : Yet to be done.
- Test      : Yet to be done

#### Usage 

```
$ ./octool
NAME:
   octool - Toolchain for OpenContainer Format

USAGE:
   octool [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   validate     validate container image / Json
   test         Test the Container
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h           show help
   --version, -v        print the version

```

##### example 
```
$ ./octool  validate  --json ./bad.json
0 Platform.Arch Can not be empty
1 Mount.Source Can not be empty

Linux Specific config errors

0 Namespace.Type Can not be empty

NOTE: One or more errors found in ./bad.json

$ ./octool  validate  --json ./test.json

 ./test.json has Valid OC Format !!

```
