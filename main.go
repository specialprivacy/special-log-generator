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
	"os"

	"github.com/urfave/cli"
)

// Schema of a SPECIAL log message.
type log struct {
	Timestamp  int64    `json:"timestamp"`
	Process    string   `json:"process"`
	Purpose    string   `json:"purpose"`
	Location   string   `json:"location"`
	UserID     string   `json:"userID"`
	Attributes []string `json:"attributes"`
}

// Schema of a SPECIAL consent event
type consent struct {
	Timestamp int64  `json:"timestamp"`
	Purpose   string `json:"purpose"`
	Location  string `json:"location"`
	UserID    string `json:"userID"`
	Attribute string `json:"attribute"`
}

// Schema of the configuration file of this application.
// The configure command will output a valid configuration file.
// It's code relies heavily on reflection, so that any changes to this struct
// are immediately reflected in the command.
type config struct {
	Process    []string `json:"process,omitempty"`
	Purpose    []string `json:"purpose,omitempty"`
	Location   []string `json:"location,omitempty"`
	UserID     []string `json:"userID,omitempty"`
	Attributes []string `json:"attributes,omitempty"`
}

// Some hardcoded default values to make life easier for the user.
var defaultConfig = config{
	Process:    []string{"mailinglist", "send-invoice"},
	Purpose:    []string{"marketing", "billing"},
	Location:   []string{"belgium", "germany", "austria", "france"},
	UserID:     makeUUIDList(5),
	Attributes: []string{"name", "age", "email", "address", "hartrate"},
}

func main() {
	app := cli.NewApp()
	app.Name = "Special Log Generator"
	app.HelpName = "slg"
	app.Usage = "Create a stream of pseudo random events in the Special format"
	app.ArgsUsage = " "
	app.EnableBashCompletion = true
	app.Version = "1.0.0"
	app.Authors = []cli.Author{
		{
			Name:  "Wouter Dullaert",
			Email: "wouter.dullaert@tenforce.com",
		},
	}
	app.Copyright = "(c) 2018 Tenforce"

	app.Commands = []cli.Command{
		generateCommand,
		configureCommand,
	}

	app.Action = func(c *cli.Context) error {
		cli.ShowAppHelp(c)
		return nil
	}

	app.Run(os.Args)
}
