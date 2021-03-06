package main

import (
	"math/rand"
	"os"

	"github.com/google/uuid"
)

// makeUUIDList creates a list of random UUID strings of length n.
func makeUUIDList(n int) []string {
	output := make([]string, n)
	for i := 0; i < n; i++ {
		output[i] = randomUUID()
	}
	return output
}

// randomUUID creates a random UUID and return its string representation
func randomUUID() string {
	return uuid.New().String()
}

// getRandomValue picks a random value from the values and returns it.
func getRandomValue(values []string) string {
	return values[rand.Intn(len(values))]
}

// getRandomList selects a random subset of values and returns it as an array.
func getRandomList(values []string) []string {
	length := rand.Intn(len(values)) + 1
	for i := len(values) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
	return values[0:length]
}

// getOutput will open a writable file or return stdout if file is empty.
func getOutput(file string) (*os.File, error) {
	if file == "" {
		return os.Stdout, nil
	}
	output, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	return output, nil
}
