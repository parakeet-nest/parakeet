package gear

import (
	"os"
	"strconv"
)

/*
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
*/

func GetEnvFloat(key string, defaultValue float64) float64 {
	if value, err := strconv.ParseFloat(os.Getenv(key), 64); err == nil {
		return value
	}
	return defaultValue
}

/*
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
*/

func GetEnvInt(key string, defaultValue int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}


func GetEnvString(key string, defaultValue string) string {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}
	return envValue
}