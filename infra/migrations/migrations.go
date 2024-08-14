package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
)

func initMigration() *migrate.Migrate {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf(".env file not found %v\n", err)
	}
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	hostname := os.Getenv("POSTGRES_HOSTNAME")
	port := os.Getenv("POSTGRES_PORT")
	database := os.Getenv("POSTGRES_DATABASE")

	postgresStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, hostname, port, database)
	fmt.Println(postgresStr)
	db, err := sql.Open("postgres", postgresStr)
	if err != nil {
		log.Fatalln("error:", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("error:", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalln("error:", err)
	}
	return m
}

func Up() {
	m := initMigration()
	err := m.Up()
	if err != nil {
		log.Fatalln("Error Migration Up: ", err)
	}
	fmt.Println("Migration Up finished")
}

func Down() {
	m := initMigration()
	err := m.Down()
	if err != nil {
		log.Fatalln("Error Migration Down: ", err)
	}
	fmt.Println("Migration Down finished")
}

func Create(name string) {
	fmt.Println("Creating Migration ", name)
	parsedName := strings.ReplaceAll(name, " ", "_")
	parsedName = strings.ToLower(parsedName)
	data := time.Now().Unix()
	downName := fmt.Sprintf("./migrations/%d_%s.down.sql", data, parsedName)
	upName := fmt.Sprintf("./migrations/%d_%s.up.sql", data, parsedName)

	fileDown, err := os.Create(downName)
	if err != nil {
		log.Fatalln("Error creating Down file: ", err)
	}
	_, err = fileDown.WriteString(fmt.Sprintf("-- Migration Down: %s\n\n", name))
	if err != nil {
		log.Fatalln("Error writing on Down file: ", err)
	}

	defer fileDown.Close()

	fileUp, err := os.Create(upName)
	if err != nil {
		log.Fatalln("Error creating Up file: ", err)
	}
	_, err = fileUp.WriteString(fmt.Sprintf("-- Migration Up: %s\n\n", name))
	if err != nil {
		log.Fatalln("Error writing on Up file: ", err)
	}
	defer fileUp.Close()

}
