package controller

import (
	"net/http"
)

// Task перенаправляет запрос на соответствующий обработчик в зависимости от метода
func (h *Handler) Task(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.AddTask(w, r)
	case http.MethodGet:
		h.GetTask(w, r)
	case http.MethodPut:
		h.EditTask(w, r)
	}
}
