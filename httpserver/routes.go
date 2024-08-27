package httpserver

import (
	"time"

	"github.com/gin-contrib/timeout"
	"github.com/gin-gonic/gin"
)

func (r *httpserver) InitRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	{
		api.POST("/populate", timeout.New(
			timeout.WithTimeout(360*time.Second),
			timeout.WithHandler(r.handler.Populate),
		))
		api.GET("/clerks", timeout.New(
			timeout.WithTimeout(30*time.Second),
			timeout.WithHandler(r.handler.Clerks),
		))
	}
}
