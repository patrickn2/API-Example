package main

import (
	"github.com/patrickn2/Clerk-Challenge/config"
	"github.com/patrickn2/Clerk-Challenge/handler"
	"github.com/patrickn2/Clerk-Challenge/infra/database"
	"github.com/patrickn2/Clerk-Challenge/router"
)

func main() {
	config.Init()
	db := database.InitDB(database.InitPostgres())
	handler := handler.NewHandler(db)
	router.Init(handler)
}
