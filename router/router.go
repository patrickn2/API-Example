package router

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/handler"
)

func Init(h *handler.Handler) {
	router := gin.Default()
	InitRoutes(h, router)
	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("api port not defined")
	}
	router.Run(":" + port)
}
