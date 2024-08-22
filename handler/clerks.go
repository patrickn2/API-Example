package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/api-challenge/interfaces"
)

func (h *Handler) Clerks(ctx *gin.Context) {
	var params interfaces.GetClerksParams
	ctx.BindQuery(&params)
	response, err := h.Services.UserService.Clerks(&params)
	if err != nil {
		log.Printf("Error: %v", err)
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
