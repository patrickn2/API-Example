package handler

import (
	"github.com/patrickn2/api-challenge/service"
)

type HandlerService struct {
	UserService *service.UserService
}

type Handler struct {
	Services *HandlerService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{
		Services: &HandlerService{
			UserService: userService,
		},
	}
}
