package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/urfave/cli"
)

// makeLog creates a log statement from a random selection of the values in config.
func makeLog(config config) log {
	return log{
		Timestamp:  time.Now().UnixNano() / int64(time.Millisecond),
		Process:    getRandomValue(config.Process),
		Purpose:    getRandomValue(config.Purpose),
		Location:   getRandomValue(config.Location),
		UserId:     getRandomValue(config.UserId),
		Attributes: getRandomList(config.Attributes),
	}
}

// generateLog sends a maximum of n random log messages through channel c at the a particular rate.
// The function is meant to run a a go-routine
// In case n <= 0 the function will keep the channel running indefinitely
func generateLog(config config, n int, rate time.Duration, c chan log) {
	if n <= 0 {
		for {
			payload := makeLog(config)
			c <- payload
			time.Sleep(rate)
		}
	} else {
		for i := 0; i < n; i++ {
			payload := makeLog(config)
			c <- payload
			time.Sleep(rate)
		}
		close(c)
	}
}

var ttlTemplate = getTtlTemplate()

// ttlMarshal renders a value in tuttle syntax according to the ttlTemplate.
// It is meant to be API compatible with json.Marshal.
//
// In the future we should replace this with a generic RDF library that renders
// a struct based on meta data defined by field tags (comparable to how json
// and xml work right now)
func ttlMarshal(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := ttlTemplate.Execute(&buf, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var generateCommand = cli.Command{
	Name:      "generate",
	Aliases:   []string{"g"},
	Usage:     "Generate log messages in the SPECIAL format",
	ArgsUsage: " ",
	Flags: []cli.Flag{
		cli.DurationFlag{
			Name:  "rate",
			Value: time.Duration(0),
			Usage: "The `rate` at which the generator outputs log statements. Understands golang time syntax eg: 1s",
		},
		cli.IntFlag{
			Name:  "num",
			Value: 10,
			Usage: "The `number` of log statements to create. Numbers <= 0 will create an infinite stream",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to config `file` containing alternative values for the logs",
		},
		cli.StringFlag{
			Name:  "output, o",
			Usage: "The `file` to which the generated log statements should be written (default: stdout)",
		},
		cli.StringFlag{
			Name:  "format, f",
			Value: "json",
			Usage: "The serialization `format` used to write the logs (json or ttl)",
		},
	},
	Action: func(c *cli.Context) error {
		rate := c.Duration("rate")
		num := c.Int("num")

		// Parse out the configuration should there be any
		configFlag := c.String("config")
		// This only makes a shallow copy, but the defaultConfig is never reused anyway, so it's not causing any issues for now
		config := defaultConfig
		if configFlag != "" {
			rawConfig, err := ioutil.ReadFile(configFlag)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			err = json.Unmarshal(rawConfig, &config)
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

		// Parse out the format flag (json or ttl)
		format := c.String("format")
		var serializer func(interface{}) ([]byte, error)
		if format == "json" {
			serializer = json.Marshal
		} else if format == "ttl" {
			serializer = ttlMarshal
		} else {
			return cli.NewExitError(fmt.Sprintf("format should be oneOf ['json', 'ttl']. Recieved %s", format), 1)
		}

		// Create the channel and start emitting messages
		ch := make(chan log)
		go generateLog(config, num, rate, ch)

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
