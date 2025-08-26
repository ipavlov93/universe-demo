package env

import (
	"os"
	"strconv"
)

// EnvironmentVariable returns env variable by given key.
// Otherwise, returns fallback.
func EnvironmentVariable(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// ParseIntEnv lookups and parse env variable by given key.
// Otherwise, returns fallback.
func ParseIntEnv(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	parsedValue, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return fallback
	}
	return int(parsedValue)
}
