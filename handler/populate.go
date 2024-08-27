package handler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Populate(ginCtx *gin.Context) {
	ctx, cancel := context.WithCancel(ginCtx)
	done := make(chan bool)
	go checkConnectionStatus(cancel, ginCtx, done)

	rows, err := h.Services.UserService.NewUsers(ctx)
	if err != nil {
		log.Printf("Error inserting users %v\n", err)
		done <- true
		ginCtx.Status(http.StatusInternalServerError)
		return
	}
	done <- true
	ginCtx.JSON(http.StatusOK, map[string]int{"rows_affected": rows})
}

func checkConnectionStatus(cancel context.CancelFunc, ctx *gin.Context, done chan bool) {
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			continue
		case <-done:
			return
		case <-ctx.Request.Context().Done():
			log.Println("Connection Aborted")
			cancel() // Cancel all http requests
			return
		}
	}
}
