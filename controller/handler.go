package controller

import "github.com/Ararat25/go_final_project/model"

type Handler struct {
	service *model.Service
}

func NewHandler(service *model.Service) *Handler {
	return &Handler{
		service: service,
	}
}
