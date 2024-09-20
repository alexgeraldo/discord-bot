package config

import (
	"log"
	"os"
	"strconv"
)

// GetEnv retrieves an environment variable or returns a default value if not set
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// GetEnvAsBool retrieves an environment variable as a boolean or returns a default value
func GetEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Error parsing bool env %s: %v", key, err)
			return defaultValue
		}
		return boolValue
	}
	return defaultValue
}

// GetEnvAsInt retrieves an environment variable as an integer or returns a default value
func GetEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("Error parsing int env %s: %v", key, err)
			return defaultValue
		}
		return intValue
	}
	return defaultValue
}
