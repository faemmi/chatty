package utils

import (
	"os"
)

func GetEnvWithDefault(key string, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
