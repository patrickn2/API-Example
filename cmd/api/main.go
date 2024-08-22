package main

import (
	"github.com/patrickn2/api-challenge/config"
	"github.com/patrickn2/api-challenge/handler"
	"github.com/patrickn2/api-challenge/httpHandler"
	"github.com/patrickn2/api-challenge/infra/db"
	"github.com/patrickn2/api-challenge/infra/migrations"
	"github.com/patrickn2/api-challenge/repository"
	"github.com/patrickn2/api-challenge/service"
)

func main() {

	// applying migrations before init
	migrations.Up()

	envs := config.Init()
	dbConnection := &db.DatabaseConnection{
		User:     envs["POSTGRES_USER"],
		Password: envs["POSTGRES_PASSWORD"],
		Hostname: envs["POSTGRES_HOSTNAME"],
		Database: envs["POSTGRES_DATABASE"],
		Port:     envs["POSTGRES_PORT"],
	}
	db := db.NewDatabase(dbConnection)
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	handler := handler.NewHandler(userService)
	httpHdlr := httpHandler.NewHttpHandler("", envs["API_PORT"], handler)
	httpHdlr.Init()
}
