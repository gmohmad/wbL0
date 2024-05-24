package utils

import (
	"log"
	"log/slog"
	"os"
)

func GetEnvOrFatal(envVar string) string {
	env := os.Getenv(envVar)

	if env == "" {
		log.Fatalf("%s is not set", envVar)
	}

	return env
}

func LogFatal(log *slog.Logger, err error) {
	log.Error(err.Error())
	os.Exit(1)
}
