package router

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/handler"
)

func InitRoutes(handler *handler.Handler, router *gin.Engine) {
	router.GET("/populate", timeout.New(
		timeout.WithTimeout(360*time.Second),
		timeout.WithHandler(handler.Populate),
	))
	router.GET("/clerks", timeout.New(
		timeout.WithTimeout(30*time.Second),
		timeout.WithHandler(handler.Clerks),
	))
}
