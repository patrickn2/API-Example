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
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/joho/godotenv/autoload"
	"github.com/patrickn2/api-challenge/config"
)

func initMigration() *migrate.Migrate {
	envs := config.Init()

	postgresStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", envs["POSTGRES_USER"], envs["POSTGRES_PASSWORD"], envs["POSTGRES_HOSTNAME"], envs["POSTGRES_PORT"], envs["POSTGRES_DATABASE"])
	db, err := sql.Open("postgres", postgresStr)
	if err != nil {
		log.Fatalln("sql Open error:", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalln("postgres connection error :", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalln("migrate error:", err)
	}
	return m
}

func Up() {
	m := initMigration()
	err := m.Up()
	if err != nil {
		log.Println("Migration Up: ", err)
		return
	}
	log.Println("Migration Up finished")
}

func Down() {
	m := initMigration()
	err := m.Down()
	if err != nil {
		log.Fatalln("Error Migration Down: ", err)
	}
	log.Println("Migration Down finished")
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
