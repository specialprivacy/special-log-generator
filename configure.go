package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/urfave/cli"
)

var configAttributes = func() []string {
	configType := reflect.TypeOf(config{})
	numAttributes := configType.NumField()
	attributes := make([]string, numAttributes)
	for i := 0; i < numAttributes; i++ {
		attributes[i] = configType.Field(i).Name
	}
	return attributes
}()

func getNumFlag(v string) string {
	camelV := strings.Replace(v, v[:1], strings.ToLower(v[:1]), 1)
	return fmt.Sprintf("%sNum", camelV)
}

func getPrefixFlag(v string) string {
	camelV := strings.Replace(v, v[:1], strings.ToLower(v[:1]), 1)
	return fmt.Sprintf("%sPrefix", camelV)
}

func createCommandFlags(attributes []string) []cli.Flag {
	numRegularFlags := 1
	output := make([]cli.Flag, len(attributes)*2+numRegularFlags)
	output[0] = cli.StringFlag{
		Name:  "output, o",
		Usage: "The `file` to which the generated log statements should be written (default: stdout)",
	}
	for i, v := range attributes {
		index := i*2 + numRegularFlags
		output[index] = cli.IntFlag{
			Name:  getNumFlag(v),
			Usage: fmt.Sprintf("The `number` of %s attribute values to generate", v),
		}
		output[index+1] = cli.StringFlag{
			Name:  getPrefixFlag(v),
			Usage: fmt.Sprintf("The prefix `string` to be used for the generated %s attributes", v),
			Value: v,
		}
	}
	return output
}

func generateValues(num int, prefix string) []string {
	output := make([]string, num)
	for i := 0; i < num; i++ {
		output[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return output
}

func generateUserIds(num int, prefix string) []string {
	if prefix != "UserId" {
		return generateValues(num, prefix)
	}
	return makeUUIDList(num)
}

var configureCommand = cli.Command{
	Name:      "configure",
	Aliases:   []string{"c"},
	Usage:     "Generate a configuration file with certain properties",
	ArgsUsage: " ",
	Flags:     createCommandFlags(configAttributes),
	Action: func(c *cli.Context) error {
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

		result := config{}
		resultValue := reflect.ValueOf(result)
		for _, attr := range configAttributes {
			resultValue.FieldByName(attr).Set(reflect.ValueOf(generateValues(
				c.Int(getNumFlag(attr)),
				c.String(getPrefixFlag(attr)),
			)))
		}

		b, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		_, err = fmt.Fprintf(output, "%s", b)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		return nil
	},
}
