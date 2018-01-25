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
	"github.com/urfave/cli"
	"os"
	"time"
)

type log struct {
	Timestamp  int64    `json:"timestamp"`
	Process    string   `json:"process"`
	Purpose    string   `json:"purpose"`
	Location   string   `json:"location"`
	UserId     string   `json:"userId"`
	Attributes []string `json:"attributes"`
}

func generateLog(n int, rate time.Duration, c chan log) {
	for i := 0; i < n; i++ {
		payload := log{
			Timestamp:  1516714505234,
			Process:    "test-process",
			Purpose:    "marketing",
			Location:   "belgium",
			UserId:     "150e2356-eaf0-4ea3-b5e0-188fb5548ccb",
			Attributes: []string{"email", "age"},
		}
		c <- payload
		time.Sleep(rate)
	}
	close(c)
}

func main() {
	app := cli.NewApp()
	app.Name = "special-log-generator"
	app.Usage = "Create a stream of pseudo random log statements in the Special format"
	app.EnableBashCompletion = true
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		ch := make(chan log)
		rate, err := time.ParseDuration("2s")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		go generateLog(10, rate, ch)

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
