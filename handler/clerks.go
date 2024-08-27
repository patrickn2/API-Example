package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/api-challenge/schema"
)

func (h *Handler) Clerks(ctx *gin.Context) {
	var params schema.GetClerksParams
	ctx.BindQuery(&params)
	response, err := h.Services.UserService.Clerks(ctx, &params)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
