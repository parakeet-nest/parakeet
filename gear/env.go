package gear

import (
	"os"
	"strconv"
)


func GetEnvFloat(key string, defaultValue float64) float64 {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	value, err := strconv.ParseFloat(envValue, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetEnvInt(key string, defaultValue int) int {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(envValue)
	if err != nil {
		return defaultValue
	}
	return value
}

func GetEnvString(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}