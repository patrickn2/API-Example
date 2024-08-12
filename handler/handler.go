package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db       *gorm.DB
	Shutdown struct {
		RunningProcess int
		ShutdownSignal bool
	}
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		db: db,
		Shutdown: struct {
			RunningProcess int
			ShutdownSignal bool
		}{0, false},
	}
}

func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.Shutdown.ShutdownSignal {
			fmt.Println("Server is Down")
			c.AbortWithStatus(http.StatusGone)
			return
		}
		h.Shutdown.RunningProcess++

		c.Next()

		h.Shutdown.RunningProcess--
	}
}
