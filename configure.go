package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/urfave/cli"
)

// Create an array of the names of all config struct keys
var configAttributes = func() []string {
	configType := reflect.TypeOf(config{})
	numAttributes := configType.NumField()
	attributes := make([]string, numAttributes)
	for i := 0; i < numAttributes; i++ {
		attributes[i] = configType.Field(i).Name
	}
	return attributes
}()

// getNumFlag creates the name of the flag for the number input for a config property v.
func getNumFlag(v string) string {
	camelV := strings.Replace(v, v[:1], strings.ToLower(v[:1]), 1)
	return fmt.Sprintf("%sNum", camelV)
}

// getPrefixFlag creates the name of the flag for the prefix input for a config property v.
func getPrefixFlag(v string) string {
	camelV := strings.Replace(v, v[:1], strings.ToLower(v[:1]), 1)
	return fmt.Sprintf("%sPrefix", camelV)
}

// createCommandFlags returns a list of cli.Flag configurations.
// It will append a few hardcoded flags to a list of flags for each config property.
func createCommandFlags(attributes []string) []cli.Flag {
	numRegularFlags := 1
	output := make([]cli.Flag, len(attributes)*2+numRegularFlags)
	output[0] = cli.StringFlag{
		Name:  "output, o",
		Usage: "The `file` to which the generated configuration should be written (default: stdout)",
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

// generateValues creates a list of length n with values based on a particular prefix.
func generateValues(num int, prefix string) []string {
	output := make([]string, num)
	for i := 0; i < num; i++ {
		output[i] = fmt.Sprintf("%s%d", prefix, i)
	}
	return output
}

var configureCommand = cli.Command{
	Name:      "configure",
	Aliases:   []string{"c"},
	Usage:     "Generate a configuration file with certain properties",
	ArgsUsage: " ",
	Flags:     createCommandFlags(configAttributes),
	Action: func(c *cli.Context) error {
		// Parse the output flag
		output, err := getOutput(c.String("output"))
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		defer output.Close()

		// Generate values for all fields of the config struct
		// Since the default number of values to be created is 0, we can just
		// naively iterate over all of them
		result := &config{}
		resultValue := reflect.ValueOf(result)
		for _, attr := range configAttributes {
			resultValue.Elem().FieldByName(attr).Set(reflect.ValueOf(generateValues(
				c.Int(getNumFlag(attr)),
				c.String(getPrefixFlag(attr)),
			)))
		}

		// Marshal the result into json and write to the specified output
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
