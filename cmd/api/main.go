package main

import (
	"github.com/patrickn2/Clerk-Challenge/config"
	"github.com/patrickn2/Clerk-Challenge/handler"
	"github.com/patrickn2/Clerk-Challenge/infra/database"
	"github.com/patrickn2/Clerk-Challenge/router"
	"github.com/patrickn2/Clerk-Challenge/usecase"
)

func main() {
	config.Init()
	db := database.InitDB(database.InitPostgres())
	userUsecase := usecase.NewUserUsecase(db)

	handlerUseCases := &handler.HandlerUseCases{
		UserUsecase: userUsecase,
	}

	handler := handler.NewHandler(handlerUseCases)
	router.Init(handler)
}
