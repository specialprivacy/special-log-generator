package main

/**
 * Copyright 2018 Tenforce.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli"
)

type log struct {
	Timestamp  int64    `json:"timestamp"`
	Process    string   `json:"process"`
	Purpose    string   `json:"purpose"`
	Location   string   `json:"location"`
	UserId     string   `json:"userId"`
	Attributes []string `json:"attributes"`
}

type config struct {
	Process    []string
	Purpose    []string
	Location   []string
	UserId     []string
	Attributes []string
}

func makeUUIDList(length int) []string {
	output := make([]string, length)
	for i := 0; i < length; i++ {
		output[i] = uuid.New().String()
	}
	return output
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

func main() {
	app := cli.NewApp()
	app.Name = "special-log-generator"
	app.Usage = "Create a stream of pseudo random log statements in the Special format"
	app.EnableBashCompletion = true
	app.Version = "0.1.0"

	var rateFlag string
	var numFlag int
	config := config{
		Process:    []string{"mailinglist", "send-invoice"},
		Purpose:    []string{"marketing", "billing"},
		Location:   []string{"belgium", "germany", "austria", "france"},
		UserId:     makeUUIDList(5),
		Attributes: []string{"name", "age", "email", "address", "hartrate"},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rate",
			Value:       "0s",
			Usage:       "The rate at which the generator outputs log statements. Understands golang time syntax eg: 1s",
			Destination: &rateFlag,
		},
		cli.IntFlag{
			Name:        "num",
			Value:       10,
			Usage:       "The number of log statements to create. Numbers <= 0 will create an infinite stream",
			Destination: &numFlag,
		},
	}

	app.Action = func(c *cli.Context) error {
		ch := make(chan log)
		rate, err := time.ParseDuration(rateFlag)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		go generateLog(config, numFlag, rate, ch)

		for log := range ch {
			b, err := json.Marshal(log)
			if err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
			fmt.Printf("%s\n", b)
		}

		return nil
	}

	app.Run(os.Args)
}
