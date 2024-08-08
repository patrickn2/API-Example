package router

import (
	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/handler"
)

func InitRoutes(handler *handler.Handler, router *gin.Engine) {
	router.GET("/populate", handler.Populate)
}
