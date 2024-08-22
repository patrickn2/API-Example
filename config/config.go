package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func Init() map[string]string {
	envs := make(map[string]string)

	// Mandatory Envs
	mandatoryEnvs := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOSTNAME",
		"POSTGRES_DATABASE",
		"POSTGRES_PORT",
		"API_PORT",
	}

	for _, val := range mandatoryEnvs {
		envVal := os.Getenv(val)
		if envVal == "" {
			log.Fatalf("environment variable %s not found\n", val)
		}
		envs[val] = envVal
	}

	// Optional Envs
	optionalEnvs := []string{}

	for _, val := range optionalEnvs {
		envVal := os.Getenv(val)
		envs[val] = envVal
	}
	return envs
}
