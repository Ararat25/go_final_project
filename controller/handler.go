package controller

import (
	"github.com/Ararat25/go_final_project/task"
)

// Handler структура для обраотчиков запросов
type Handler struct {
	service *task.Service
}

// NewHandler создает новый объект Handler
func NewHandler(service *task.Service) *Handler {
	return &Handler{
		service: service,
	}
}
