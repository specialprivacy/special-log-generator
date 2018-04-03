package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
	"time"

	"github.com/urfave/cli"
)

// makeLog creates a log statement from a random selection of the values in config.
func makeLog(config config) interface{} {
	return log{
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Process:   getRandomValue(config.Process),
		Purpose:   getRandomValue(config.Purpose),
		Storage:   getRandomValue(config.Storage),
		UserID:    getRandomValue(config.UserID),
		Data:      getRandomList(config.Data),
	}
}

// makeConsent creates a consent event from a random selection of the values in the config.
func makeConsent(config config) interface{} {
	return consent{
		ConsentID:  randomUUID(),
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Purpose:    getRandomValue(config.Purpose),
		Processing: getRandomValue(config.Processing),
		Recipient:  getRandomValue(config.Recipient),
		Storage:    getRandomValue(config.Storage),
		UserID:     getRandomValue(config.UserID),
		Data:       getRandomValue(config.Data),
	}
}

// generateLog sends a maximum of n random messages through channel c at the a particular rate.
// The function is meant to run a a go-routine
// In case n <= 0 the function will keep the channel running indefinitely
func generateLog(
	config config,
	n int,
	rate time.Duration,
	producer func(config) interface{},
	c chan interface{},
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
			Name:  "rate",
			Value: time.Duration(0),
			Usage: "The `rate` at which the generator outputs events. Understands golang time syntax eg: 1s",
		},
		cli.IntFlag{
			Name:  "num",
			Value: 10,
			Usage: "The `number` of events to create. Numbers <= 0 will create an infinite stream",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to config `file` containing alternative values for the events",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "The `file` to which the generated events should be written (default: stdout)",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "The serialization `format` used to write the events (json or ttl)",
		},
		cli.StringFlag{
			Name:  "type, t",
			Value: "log",
			Usage: "The `type` of event to be generated (log or consent)",
		},
	},
	Action: func(c *cli.Context) error {
		rate := c.Duration("rate")
		num := c.Int("num")

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

		// Parse the output flag
		output, err := getOutput(c.String("output"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer output.Close()

		// Parse out the type flag (log or consent)
		eventType := c.String("type")
		var producer func(config) interface{}
		var ttlTemplate *template.Template
		if eventType == "log" {
			producer = makeLog
			ttlTemplate = getLogTTLTemplate()
		} else if eventType == "consent" {
			producer = makeConsent
			ttlTemplate = getConsentTTLTemplate()
		} else {
			return cli.NewExitError(fmt.Sprintf("type should be oneOf ['log']"), 1)
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
		ch := make(chan interface{})
		go generateLog(conf, num, rate, producer, ch)

		// For each message call the serializer and write to the output
		for log := range ch {
			b, err := serializer(log)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			_, err = fmt.Fprintf(output, "%s\n", b)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}

		return nil
	},
}
