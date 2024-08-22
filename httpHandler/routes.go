package httpHandler

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
	"github.com/patrickn2/api-challenge/handler"
)

func InitRoutes(handler *handler.Handler, router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/populate", timeout.New(
			timeout.WithTimeout(360*time.Second),
			timeout.WithHandler(handler.Populate),
		))
		api.GET("/clerks", timeout.New(
			timeout.WithTimeout(30*time.Second),
			timeout.WithHandler(handler.Clerks),
		))
	}
}
