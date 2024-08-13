package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Clerks(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "")
	startingAfter := ctx.DefaultQuery("starting_after", "")
	endingBefore := ctx.DefaultQuery("ending_before", "")
	email := ctx.DefaultQuery("email", "")

	_, limitOk := strconv.Atoi(limit)
	_, startingAfterOk := strconv.Atoi(startingAfter)
	_, endingBeforeOk := strconv.Atoi(endingBefore)
	if (limitOk != nil && limit != "") || (startingAfterOk != nil && startingAfter != "") || (endingBeforeOk != nil && endingBefore != "") {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	users, err := h.UserCase.GetUsers(limit, startingAfter, endingBefore, email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}
