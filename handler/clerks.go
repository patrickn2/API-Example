package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/model"
)

func (h *Handler) Clerks(ctx *gin.Context) {
	// default: 10 users sorted by created_at desc
	// Parameters:
	// limit: Default 10, between 1 and 100
	// starting_after: cursor for pagination (user id)
	// ending_before: cursor for pagination (user id)
	// email (case insenstive)
	limit := ctx.DefaultQuery("limit", "")
	startingAfter := ctx.DefaultQuery("starting_after", "")
	endingBefore := ctx.DefaultQuery("ending_before", "")
	email := ctx.DefaultQuery("email", "")

	l, limitOk := strconv.Atoi(limit)
	_, startingAfterOk := strconv.Atoi(startingAfter)
	_, endingBeforeOk := strconv.Atoi(endingBefore)
	if (limitOk != nil && limit != "") || (startingAfterOk != nil && startingAfter != "") || (endingBeforeOk != nil && endingBefore != "") {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if limitOk == nil && (l < 1 || l > 100) {
		ctx.AbortWithError(http.StatusBadRequest, errors.New("limit should be between 1 and 100"))
		return
	}

	users, err := model.GetUsers(h.db, limit, startingAfter, endingBefore, email)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx.JSON(http.StatusOK, users)
}
