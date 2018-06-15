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
	Processing string   `json:"processing"`
	Recipient  string   `json:"recipient"`
	Storage    string   `json:"storage"`
	UserID     string   `json:"userID"`
	Data       []string `json:"data"`
	EventID    string   `json:"eventID"`
}

// Schema of a SPECIAL simplepolicy event
type simplepolicy struct {
	Purpose    string `json:"purposeCollection"`
	Processing string `json:"processingCollection"`
	Recipient  string `json:"recipientCollection"`
	Storage    string `json:"storageCollection"`
	Data       string `json:"dataCollection"`
}

type policy struct {
	ConsentID      string         `json:"-"`
	Timestamp      int64          `json:"timestamp"`
	UserID         string         `json:"userID"`
	SimplePolicies []simplepolicy `json:"simplePolicies"`
}

// Schema of the configuration file of this application.
// The configure command will output a valid configuration file.
// It's code relies heavily on reflection, so that any changes to this struct
// are immediately reflected in the command.
type config struct {
	Process       []string `json:"process,omitempty"`
	Purpose       []string `json:"purpose,omitempty"`
	Processing    []string `json:"processing,omitempty"`
	Recipient     []string `json:"recipient,omitempty"`
	Storage       []string `json:"storage,omitempty"`
	UserID        []string `json:"userID,omitempty"`
	Data          []string `json:"data,omitempty"`
	MaxPolicySize int      `json:"maxPolicySize,omitempty"`
}

func makeDefaultConfig() config {
	// Some hardcoded default values to make life easier for the user.
	defaultProcess := []string{"mailinglist", "send-invoice"}
	defaultPurpose := []string{"spl:AnyPurpose", "svpu:Account", "svpu:Admin", "svpu:AnyContact", "svpu:Arts", "svpu:AuxPurpose", "svpu:Browsing", "svpu:Charity", "svpu:Communicate", "svpu:Current", "svpu:Custom", "svpu:Delivery", "svpu:Develop", "svpu:Downloads", "svpu:Education", "svpu:Feedback", "svpu:Finmgt", "svpu:Gambling", "svpu:Gaming", "svpu:Government", "svpu:Health", "svpu:Historical", "svpu:Login", "svpu:Marketing", "svpu:News", "svpu:OtherContact", "svpu:Payment", "svpu:Sales", "svpu:Search", "svpu:State", "svpu:Tailoring", "svpu:Telemarketing"}
	for index, purpose := range defaultPurpose {
		defaultPurpose[index] = expandPrefix(purpose)
	}
	defaultProcessing := []string{"spl:AnyProcessing", "svpr:Aggregate", "svpr:Analyze", "svpr:Anonymize", "svpr:Collect", "svpr:Copy", "svpr:Derive", "svpr:Move", "svpr:Query", "svpr:Transfer"}
	for index, processing := range defaultProcessing {
		defaultProcessing[index] = expandPrefix(processing)
	}
	defaultRecipient := []string{"spl:AnyRecipient", "svr:Delivery", "svr:OtherRecipient", "svr:Ours", "svr:Public", "svr:Same", "svr:Unrelated"}
	for index, recipient := range defaultRecipient {
		defaultRecipient[index] = expandPrefix(recipient)
	}
	defaultStorage := []string{"spl:AnyLocation", "svl:ControllerServers", "svl:EU", "svl:EULike", "svl:ThirdCountries", "svl:OurServers", "svl:ProcessorServers", "svl:ThirdParty"}
	for index, storage := range defaultStorage {
		defaultStorage[index] = expandPrefix(storage)
	}
	defaultData := []string{"spl:AnyData", "svd:Activity", "svd:Anonymized", "svd:AudiovisualActivity", "svd:Computer", "svd:Content", "svd:Demographic", "svd:Derived", "svd:Financial", "svd:Government", "svd:Health", "svd:Interactive", "svd:Judicial", "svd:Location", "svd:Navigation", "svd:Online", "svd:OnlineActivity", "svd:Physical", "svd:PhysicalActivity", "svd:Political", "svd:Preference", "svd:Profile", "svd:Purchase", "svd:Social", "svd:State", "svd:Statistical", "svd:TelecomActivity", "svd:UniqueId"}
	for index, data := range defaultData {
		defaultData[index] = expandPrefix(data)
	}

	return config{
		Process:       defaultProcess,
		Purpose:       defaultPurpose,
		Processing:    defaultProcessing,
		Storage:       defaultStorage,
		Recipient:     defaultRecipient,
		UserID:        makeUUIDList(5),
		Data:          defaultData,
		MaxPolicySize: 5,
	}

}

var defaultConfig = makeDefaultConfig()

func main() {
	app := cli.NewApp()
	app.Name = "Special Log Generator"
	app.HelpName = "slg"
	app.Usage = "Create a stream of pseudo random events in the Special format"
	app.ArgsUsage = " "
	app.EnableBashCompletion = true
	app.Version = "2.0.0"
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
