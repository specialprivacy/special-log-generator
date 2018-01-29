package main

import (
	"math/rand"

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
	length := rand.Intn(len(values))
	for i := len(values) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
	return values[0:length]
}
