package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/usecase"
)

type HandlerUseCases struct {
	UserUsecase *usecase.UserUsecase
}

type Handler struct {
	UserCase *usecase.UserUsecase
	Shutdown struct {
		RunningProcess int
		ShutdownSignal bool
	}
}

func NewHandler(h *HandlerUseCases) *Handler {
	return &Handler{
		UserCase: h.UserUsecase,
		Shutdown: struct {
			RunningProcess int
			ShutdownSignal bool
		}{0, false},
	}
}

func (h *Handler) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if h.Shutdown.ShutdownSignal {
			fmt.Println("Server is shutting down - not accepting new requests")
			c.AbortWithStatus(http.StatusGone)
			return
		}
		h.Shutdown.RunningProcess++

		c.Next()

		h.Shutdown.RunningProcess--
	}
}
