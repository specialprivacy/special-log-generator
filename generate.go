package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"text/template"
	"time"

	"github.com/urfave/cli"
	sarama "gopkg.in/Shopify/sarama.v1"
)

type message struct {
	Key   string
	Value interface{}
}

// makeLog creates a log statement from a random selection of the values in config.
func makeLog(config config) message {
	log := log{
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Process:    getRandomValue(config.Process),
		Purpose:    getRandomValue(config.Purpose),
		Processing: getRandomValue(config.Processing),
		Recipient:  getRandomValue(config.Recipient),
		Storage:    getRandomValue(config.Storage),
		UserID:     getRandomValue(config.UserID),
		Data:       getRandomList(config.Data),
		EventID:    randomUUID(),
	}
	return message{
		Key:   log.EventID,
		Value: log,
	}
}

// makeConsent creates a consent event from a random selection of the values in the config.
func makeConsent(config config) message {
	simplePolicies := make([]simplepolicy, rand.Intn(config.MaxPolicySize))
	for i := range simplePolicies {
		simplePolicies[i] = simplepolicy{
			Purpose:    getRandomValue(config.Purpose),
			Processing: getRandomValue(config.Processing),
			Recipient:  getRandomValue(config.Recipient),
			Storage:    getRandomValue(config.Storage),
			Data:       getRandomValue(config.Data),
		}
	}
	policy := policy{
		ConsentID:      randomUUID(),
		Timestamp:      time.Now().UnixNano() / int64(time.Millisecond),
		UserID:         getRandomValue(config.UserID),
		SimplePolicies: simplePolicies,
	}
	return message{
		Key:   policy.ConsentID,
		Value: policy,
	}
}

// generateLog sends a maximum of n random messages through channel c at the a particular rate.
// The function is meant to run a a go-routine
// In case n <= 0 the function will keep the channel running indefinitely
func generateLog(
	config config,
	n int,
	rate time.Duration,
	producer func(config) message,
	c chan message,
) {
	if n <= 0 {
		for {
			payload := producer(config)
			c <- payload
			time.Sleep(rate)
		}
	} else {
		for i := 0; i < n; i++ {
			payload := producer(config)
			c <- payload
			time.Sleep(rate)
		}
		close(c)
	}
}

// createTTLMarshal creates a function that renders a value in tuttle syntax according to the ttlTemplate.
// The created function is meant to be API compatible with json.Marshal.
//
// In the future we should replace this with a generic RDF library that renders
// a struct based on meta data defined by field tags (comparable to how json
// and xml work right now)
func createTTLMarshal(ttlTemplate *template.Template) func(v interface{}) ([]byte, error) {
	return func(v interface{}) ([]byte, error) {
		var buf bytes.Buffer
		err := ttlTemplate.Execute(&buf, v)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
}

var generateCommand = cli.Command{
	Name:      "generate",
	Aliases:   []string{"g"},
	Usage:     "Generate events in the SPECIAL format",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		cli.DurationFlag{
			Name:   "rate",
			Value:  time.Duration(0),
			Usage:  "The `rate` at which the generator outputs events. Understands golang duration syntax eg: 1s",
			EnvVar: "RATE",
		},
		cli.IntFlag{
			Name:   "num",
			Value:  10,
			Usage:  "The `number` of events to create. Numbers <= 0 will create an infinite stream",
			EnvVar: "NUM",
		},
		cli.StringFlag{
			Name:   "config, c",
			Usage:  "Path to config `file` containing alternative values for the events",
			EnvVar: "CONFIG",
		},
		cli.StringFlag{
			Name:   "output, o",
			Usage:  "The `file` to which the generated events should be written. If the special value 'kafka' is used, logs will be produced on kafka. (default: stdout)",
			EnvVar: "OUTPUT",
		},
		cli.StringFlag{
			Name:   "format, f",
			Value:  "json",
			Usage:  "The serialization `format` used to write the events (json or ttl)",
			EnvVar: "FORMAT",
		},
		cli.StringFlag{
			Name:   "type, t",
			Value:  "log",
			Usage:  "The `type` of event to be generated (log or consent)",
			EnvVar: "TYPE",
		},
		cli.StringSliceFlag{
			Name:   "kafka-broker-list",
			Usage:  "A comma separated list of `brokers` used to bootstrap the connection to a kafka cluster. eg: 127.0.0.1,172.10.50.4",
			EnvVar: "KAFKA_BROKER_LIST",
		},
		cli.StringFlag{
			Name:   "kafka-topic",
			Value:  "application-logs",
			Usage:  "The name of the topic on which logs will be produced. (default: application-logs)",
			EnvVar: "KAFKA_TOPIC",
		},
		cli.StringFlag{
			Name:   "kafka-cert-file",
			Usage:  "The `path` to a certificate file used for client authentication to kafka.",
			EnvVar: "KAFKA_CERT_FILE",
		},
		cli.StringFlag{
			Name:   "kafka-key-file",
			Usage:  "The `path` to a key file used for client authentication to kafka.",
			EnvVar: "KAFKA_KEY_FILE",
		},
		cli.StringFlag{
			Name:   "kafka-ca-file",
			Usage:  "The `path` to a ca file used for client authentication to kafka.",
			EnvVar: "KAFKA_CA_FILE",
		},
		cli.BoolFlag{
			Name:   "kafka-verify-ssl",
			Usage:  "Set to verify the SSL chain when connecting to kafka",
			EnvVar: "KAFKA_VERIFY_SSL",
		},
	},
	Action: func(c *cli.Context) error {
		rate := c.Duration("rate")
		num := c.Int("num")

		// Ensure rate and num are using sane combinations
		if rate == 0 && num <= 0 {
			return cli.NewExitError("Streaming (num <= 0) must be used with a non-zero rate duration", 1)
		}

		// Parse out the configuration should there be any
		configFlag := c.String("config")
		// This only makes a shallow copy, but the defaultConfig is never reused anyway, so it's not causing any issues for now
		conf := defaultConfig
		if configFlag != "" {
			rawConfig, err := ioutil.ReadFile(configFlag)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			err = json.Unmarshal(rawConfig, &conf)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		// Parse the output flag and kafka options
		var kafkaProducer sarama.SyncProducer
		var output *os.File
		kafkaTopic := c.String("kafka-topic")
		if c.String("output") == "kafka" {
			fmt.Println("[INFO] Writing logs to kafka")
			var err error
			kafkaProducer, err = createKafkaProducer(kafkaConfig{
				BrokerList: c.StringSlice("kafka-broker-list"),
				CertFile:   c.String("kafka-cert-file"),
				KeyFile:    c.String("kafka-key-file"),
				CaFile:     c.String("kafka-ca-file"),
				VerifySsl:  c.Bool("kafka-verify-ssl"),
			})
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			defer kafkaProducer.Close()
			fmt.Printf("[INFO] Successfully connected to kafka cluster at %s\n", c.StringSlice("kafka-broker-list"))
		} else {
			var err error
			output, err = getOutput(c.String("output"))
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			defer output.Close()
		}

		// Parse out the type flag (log or consent)
		eventType := c.String("type")
		var producer func(config) message
		var ttlTemplate *template.Template
		if eventType == "log" {
			producer = makeLog
			ttlTemplate = getLogTTLTemplate()
		} else if eventType == "consent" {
			producer = makeConsent
			ttlTemplate = getConsentTTLTemplate()
		} else {
			return cli.NewExitError(fmt.Sprintf("type should be oneOf ['log', 'consent']"), 1)
		}

		// Parse out the format flag (json or ttl)
		format := c.String("format")
		var serializer func(interface{}) ([]byte, error)
		if format == "json" {
			serializer = json.Marshal
		} else if format == "ttl" {
			serializer = createTTLMarshal(ttlTemplate)
		} else {
			return cli.NewExitError(fmt.Sprintf("format should be oneOf ['json', 'ttl']. Recieved %s", format), 1)
		}

		// Create the channel and start emitting messages
		ch := make(chan message)
		go generateLog(conf, num, rate, producer, ch)

		// For each message call the serializer and write to the output
		if kafkaProducer != nil {
			for log := range ch {
				b, err := serializer(log.Value)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				_, _, err = kafkaProducer.SendMessage(&sarama.ProducerMessage{
					Topic: kafkaTopic,
					Key:   sarama.StringEncoder(log.Key),
					Value: sarama.StringEncoder(b),
				})
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
			}
			fmt.Printf("[INFO] Done writing %d messages to kafka\n", c.Int("num"))
		} else {
			for log := range ch {
				b, err := serializer(log.Value)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				_, err = fmt.Fprintf(output, "%s\n", b)
				if err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
			}
		}

		return nil
	},
}
