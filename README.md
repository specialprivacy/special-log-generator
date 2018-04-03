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
slg [global options] command [command options]
```

The application has the following commands:
- **generate** Generate messages in the special formats

### Generate Options
- `--rate`: The rate at which the generator outputs events. This parameter understands golang time syntax eg: `1s` or `10ms` (default: `0s`)
- `--num`: The total number of events that will be generated. When this parameters is <=0 it will create an infinite stream (default: `10`)
- `--config`: The path to a config file in json containing alternative values for the events
- `--output`: The file to which the events will be written. Be carefull, because this will overwrite the file should it already exist (default: `stdout`)
- `--format`: The serialization format used to write the events (json or ttl) (default: `json`)
- `--type`: The type of event to be generated (log or consent) (default: "log")

### Config file format
The config file format is json which takes the following keys:
- `process`: An array of strings with potential values for `process`
- `purpose`: An array of strings with potential values for `purpose`
- `location`: An array of strings with potential values for `location`
- `userID`: An array of strings with potential values for `userId`
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
  "userID": ["1", "2", "3"],
  "attributes": ["height", "gender"]
}
```

### Configure Options
- `--output`: The file to which the generated configuration should be written (default: `stdout`)
- `--processNum`: The number of Process attribute values to generate (default: `0`)
- `--processPrefix`: The prefix string to be used for the generated Process attributes (default: `Process`)
- `--purposeNum`: The number of Purpose attribute values to generate (default: `0`)
- `--purposePrefix`: The prefix string to be used for the generated Purpose attributes (default: `Purpose`)
- `--locationNum`: The number of Location attribute values to generate (default: `0`)
- `--locationPrefix`: The prefix string to be used for the generated Location attributes (default: `Location`)
- `--userIDNum`: The number of UserId attribute values to generate (default: `0`)
- `--userIDPrefix`: The prefix string to be used for the generated UserId attributes (default: `UserId`)
- `--attributesNum`: The number of Attributes attribute values to generate (default: `0`)
- `--attributesPrefix`: The prefix string to be used for the generated Attributes attributes (default: `Attributes`)

### Examples
- Print 10 random logs to `stdout`
```bash
slg generate
```
- Print 10 random consents to `stdout`
```bash
sgl generate -t consent
```
- Print 100 random logs to `logs.json`
```bash
slg generate --num 100 --output logs.json
```
- Print 10 logs at a rate of 1 every second to `stdout`
```bash
slg generate --rate 1s
```
- Pipe an infinite stream of logs every 10ms to apache kafka
```bash
slg generate --rate 10ms --num -1 | kafka-cli-producer --broker-list http://kafka:9300 --zookeeper http://zookeeper:2181 --topic special-logs
```

## Build
Dependencies are vendored into the codebase, and managed through dep (https://golang.github.io/dep/)
This makes building as simple as

```bash
git clone https://git.ai.wu.ac.at/specialprivacy/special-log-generator.git # Ensure this is somewhere on the $GOPATH
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
* Bring log format in ttl in line with deliverable
* Investigate an option to group generated policies by userID (might have memory usage consequences at high rates / volumes, will most likely be mutually exlcusive streaming)
* Get a decision whether policies are linked with a datasubject through `#hasPolicy` or `#hasDataSubject` and add these properties to the vocabulary

## LICENSE
Apache-2.0 © Tenforce
