package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf(".env file not found %v\n", err)
	}

	mandatoryEnvs := []string{
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_HOSTNAME",
		"POSTGRES_DATABASE",
		"POSTGRES_PORT",
		"API_PORT",
	}

	for _, val := range mandatoryEnvs {
		if os.Getenv(val) == "" {
			log.Fatalf("environment variable %s not found\n", val)
		}
	}
}
