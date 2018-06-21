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
special-log-generator [global options] command [command options]
```

The application has the following commands:
- **generate** Generate messages in the special formats

### Generate Options
- `--rate`: The rate at which the generator outputs events. This parameter understands golang duration syntax eg: `1s` or `10ms` (default: `0s`) [$RATE]
- `--num`: The total number of events that will be generated. When this parameters is <=0 it will create an infinite stream (default: `10`) [$NUM]
- `--config`: The path to a config file in json containing alternative values for the events [$CONFIG]
- `--output`: The file to which the generated events should be written. If the special value 'kafka' is used, logs will be produced on kafka. (default: `stdout`) [$OUTPUT]
- `--format`: The serialization format used to write the events (json or ttl) (default: `json`) [$FORMAT]
- `--type`: The type of event to be generated (log or consent) (default: `log`) [$TYPE]
- `--max-policy-size number`: The maximum number of policies to be used in a single consent (only applicable for type consent) (default: `5`) [$MAX_POLICY_SIZE]
- `--kafka-broker-list`: A comma separated list of brokers used to bootstrap the connection to a kafka cluster. eg: `127.0.0.1,172.10.50.4` [$KAFKA_BROKER_LIST]
- `--kafka-topic`: The name of the topic on which logs will be produced. (default: `application-logs`) [$KAFKA_TOPIC]
- `--kafka-cert-file`: The path to a certificate file used for client authentication to kafka. [$KAFKA_CERT_FILE]
- `--kafka-key-file`: The path to a key file used for client authentication to kafka. [$KAFKA_KEY_FILE]
- `--kafka-ca-file`: The path to a ca file used for client authentication to kafka. [$KAFKA_CA_FILE]
- `--kafka-verify-ssl`: Set to verify the SSL chain when connecting to kafka [$KAFKA_VERIFY_SSL]

### Config file format
The config file format is json which takes the following keys:
- `process`: An array of strings with potential values for `process`
- `purpose`: An array of strings with potential values for `purpose`
- `processing`: An array of strings with potential values for `processing`
- `recipient`: An array of strings with potential values for `recipient`
- `storage`: An array of strings with potential values for `storage`
- `userID`: An array of strings with potential values for `userID`
- `data`: An array of strings with potential values for `data`

If the config file contains unknown keys, they will be ignored.
If the type of any of the defined keys does not match, an error with a (hopefully) useful description will be shown.
**All keys are optional.**

Example of a fully specified config file:

```json
{
  "process": ["foo", "bar"],
  "purpose": ["marketing", "delivery"],
  "processing": ["categorization", "archiving"],
  "recipient": ["affiliates", "google"],
  "storage": ["greenland", "iceland", "mordor"],
  "userID": ["1", "2", "3"],
  "data": ["height", "gender"]
}
```

### Configure Options
- `--output file, -o file`: The file to which the generated configuration should be written (default: `stdout`)
- `--processNum number`: The number of Process attribute values to generate (default: `0`)
- `--processPrefix string`: The prefix string to be used for the generated Process attributes (default: `Process`)
- `--purposeNum number`: The number of Purpose attribute values to generate (default: `0`)
- `--purposePrefix string`: The prefix string to be used for the generated Purpose attributes (default: `Purpose`)
- `--processingNum number`: The number of Processing attribute values to generate (default: `0`)
- `--processingPrefix string`: The prefix string to be used for the generated Processing attributes (default: `Processing`)
- `--recipientNum number`: The number of Recipient attribute values to generate (default: `0`)
- `--recipientPrefix string`: The prefix string to be used for the generated Recipient attributes (default: `Recipient`)
- `--storageNum number`: The number of Storage attribute values to generate (default: `0`)
- `--storagePrefix string`: The prefix string to be used for the generated Storage attributes (default: `Storage`)
- `--userIDNum number`: The number of UserID attribute values to generate (default: `0`)
- `--userIDPrefix string`: The prefix string to be used for the generated UserID attributes (default: `UserID`)
- `--dataNum number`: The number of Data attribute values to generate (default: `0`)
- `--dataPrefix string`: The prefix string to be used for the generated Data attributes (default: `Data`)

### Examples
- Print 10 random logs to `stdout`
```bash
special-log-generator generate
```
- Print 10 random consents to `stdout`
```bash
special-log-generator generate -t consent
```
- Print 100 random logs to `logs.json`
```bash
special-log-generator generate --num 100 --output logs.json
```
- Print 10 logs at a rate of 1 every second to `stdout`
```bash
special-log-generator generate --rate 1s
```
- Pipe an infinite stream of logs every 10ms to apache kafka
```bash
special-log-generator generate --rate 10ms --num -1 --output kafka --kafka-broker-list kafka:9092 --kafka-topic special-logs
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
* Investigate an option to group generated policies by userID (might have memory usage consequences at high rates / volumes, will most likely be mutually exclusive streaming)
* Get a decision whether policies are linked with a datasubject through `#hasPolicy` or `#hasDataSubject` and add these properties to the vocabulary
* Change json outputs so that we can easily add an `@context` which results in the ttl represenation

## LICENSE
Apache-2.0 © Tenforce
