package controller

import "github.com/Ararat25/go_final_project/model"

// Handler структура для обраотчиков запросов
type Handler struct {
	service *model.Service
}

// NewHandler создает новый объект Handler
func NewHandler(service *model.Service) *Handler {
	return &Handler{
		service: service,
	}
}
