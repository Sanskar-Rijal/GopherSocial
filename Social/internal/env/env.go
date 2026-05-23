package env

import (
	"os"
	"strconv"
)

func GetString(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func GetInt(key string, fallback int) int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	// converting the string value to int
	valAsInt, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return valAsInt
}

func GetBool(key string, fallback bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	//converting the string value to bool
	valAsBool, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return valAsBool
}
