package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func LoadFromEnvBool(key string, defaultValue bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	b, err := strconv.ParseBool(val)
	if err != nil {
		return defaultValue
	}

	return b
}

func LoadFromEnvInt(key string, defaultValue int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultValue
	}

	return i
}

func LoadFromEnvInt64(key string, defaultValue int64) int64 {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return defaultValue
	}

	return i
}

func LoadFromEnvString(key string, defaultValue string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	return val
}

func LoadFromEnvStringSlice(key string, defaultValue []string) []string {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	slice := strings.Split(val, ",")
	if len(slice) == 0 {
		return defaultValue
	}

	return slice
}

func LoadFromEnvTimeDuration(key string, defaultValue time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultValue
	}

	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultValue
	}

	return d
}
