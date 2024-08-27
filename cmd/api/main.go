package main

import (
	"github.com/patrickn2/api-challenge/config"
	"github.com/patrickn2/api-challenge/handler"
	"github.com/patrickn2/api-challenge/httpserver"
	"github.com/patrickn2/api-challenge/infra/db"
	"github.com/patrickn2/api-challenge/infra/migrations"
	"github.com/patrickn2/api-challenge/repository"
	"github.com/patrickn2/api-challenge/service"
)

func main() {

	// applying migrations before init
	migrations.Up()

	env := config.Init()
	dbConnection := &db.DatabaseConnection{
		User:     env.PostgresUser,
		Password: env.PostgresPassword,
		Hostname: env.PostgresHostname,
		Database: env.PostgresDatabase,
		Port:     env.PostgresPort,
	}
	db := db.NewDatabase(dbConnection)
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	handler := handler.New(userService)
	httpHdlr := httpserver.New("", env.ApiPort, handler)
	httpHdlr.Start()
}
