# Special Log Generator
Commandline utility which will generate pseudo random log statements in the
special format.

## Download and Install
Binaries for a large number of platforms can be downloaded for every tagged release.
These can be found in the gitlab UI, by going to 'reposiotry -> tags' and then clicking the download button next to the version you are looking for.
The binaries are standalone and do not have any further dependencies

## Usage

### Basic usage
```bash
slg [global options]
```

### Options
- `--rate`: The rate at which the generator outputs log statements. This parameter understands golang time syntax eg: `1s` or `10ms` (default: `0s`)
- `--num`: The total number of statements that will be generated. When this parameters is <=0 it will create an infinite stream (default: `10`)
- `--config`: The path to a config file in json containing alternative values for the logs
- `--output`: The file to which the logs will be written. Be carefull, because this will overwrite the file should it already exist (default: `stdout`)

### Config file format
The config file format is json which takes the following keys:
- `process`: An array of strings with potential values for `process`
- `purpose`: An array of strings with potential values for `purpose`
- `location`: An array of strings with potential values for `location`
- `userId`: An array of strings with potential values for `userId`
- `attributes`: An array of strings with potential values for `attributes`

If the config file contains unknown keys, they will be ignored.
If the type of any of the defined keys does not match, an error with a (hopefully) useful description will be shown.
**All keys are optional.**

Example of a fully specified config file:

```json
{
  "process": ["foo", "bar"],
  "purpose": ["marketing", "delivery"],
  "location": ["greenland", "iceland", "mordor"],
  "userId": ["1", "2", "3"],
  "attributes": ["height", "gender"]
}
```

## Build
Dependencies are vendored into the codebase, and managed through dep (https://golang.github.io/dep/)
This makes building as simple as

```bash
git clone https://git.ai.wu.ac.at/specialprivacy/special-log-generator.git # Ensure this is somewhere on the gopath
go build
```

In case you do not have go installed, but do have docker. The application can be built as follows
```bash
docker run -it golang:1.9-alpine /bin/sh
apk --update add git
mkdir -p /go/src
cd /go/src
git clone https://git.ai.wu.ac.at/specialprivacy/special-log-generator.git
cd special-log-generator
go build
```

## TODO
* Support additional serialization formats besides json

## LICENSE
Apache-2.0 © Tenforce
