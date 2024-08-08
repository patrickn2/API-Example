package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func InitPostgres() (gorm.Dialector, *gorm.Config) {
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	hostname := os.Getenv("POSTGRES_HOSTNAME")
	database := os.Getenv("POSTGRES_DATABASE")
	port := os.Getenv("POSTGRES_PORT")
	sch := os.Getenv("POSTGRES_SCHEMA")

	config := gorm.Config{}
	if sch != "" {
		config = gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: sch + ".",
			}}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", hostname, user, password, database, port)
	return postgres.Open(dsn), &config
}
