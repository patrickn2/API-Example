package config

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-envconfig"
)

type Envs struct {
	PostgresUser     string `env:"POSTGRES_USER, required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD, required"`
	PostgresHostname string `env:"POSTGRES_HOSTNAME, required"`
	PostgresDatabase string `env:"POSTGRES_DATABASE, required"`
	PostgresPort     string `env:"POSTGRES_PORT, required"`
	ApiPort          string `env:"API_PORT, required"`
}

var envs Envs

func GetEnvs() *Envs {
	return &envs
}

func Init() *Envs {
	ctx := context.Background()
	if err := envconfig.Process(ctx, &envs); err != nil {
		log.Fatal(err)
	}
	return &envs
}
