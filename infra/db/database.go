package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	sch "gorm.io/gorm/schema"
)

type DatabaseConnection struct {
	User     string
	Password string
	Hostname string
	Database string
	Port     string
	Sch      string
}

type Database struct {
	Conn *gorm.DB
}

func NewDatabase(c *DatabaseConnection) *Database {
	config := &gorm.Config{}
	if c.Sch != "" {
		config = &gorm.Config{
			NamingStrategy: sch.NamingStrategy{
				TablePrefix: c.Sch + ".",
			}}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", c.Hostname, c.User, c.Password, c.Database, c.Port)
	dial := postgres.Open(dsn)
	conn, err := gorm.Open(dial, config)
	if err != nil {
		log.Fatalln("error connecting to database ", err)
	}
	return &Database{conn}
}
