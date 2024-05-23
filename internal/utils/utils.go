package utils

import (
	"log"
	"os"
)

func GetEnvOrFatal(envVar string) string {
	env := os.Getenv(envVar)

	if env == "" {
		log.Fatalf("%s is not set", envVar)
	}

	return env
}
