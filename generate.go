package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli"
)

func makeUUIDList(length int) []string {
	output := make([]string, length)
	for i := 0; i < length; i++ {
		output[i] = randomUUID()
	}
	return output
}

func randomUUID() string {
	return uuid.New().String()
}

func getRandomValue(values []string) string {
	return values[rand.Intn(len(values))]
}

func getRandomList(values []string) []string {
	length := rand.Intn(len(values))
	for i := len(values) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
	return values[0:length]
}

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

		outputFlag := c.String("output")
		var output io.Writer
		if outputFlag == "" {
			output = os.Stdout
		} else {
			file, err := os.Create(outputFlag)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			output = file
			defer file.Close()
		}

		format := c.String("format")
		var serializer func(interface{}) ([]byte, error)
		if format == "json" {
			serializer = json.Marshal
		} else if format == "ttl" {
			serializer = ttlMarshal
		} else {
			return cli.NewExitError(fmt.Sprintf("format should be oneOf ['json', 'ttl']. Recieved %s", format), 1)
		}

		ch := make(chan log)
		go generateLog(config, num, rate, ch)

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
