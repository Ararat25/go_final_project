package controller

import (
	"net/http"
)

// DoneTask обработчик для отметки о выполнении задачи
func (h *Handler) DoneTask(w http.ResponseWriter, r *http.Request) {
	idString := r.FormValue("id")

	id, err := checkID(idString)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DoneTask(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK)
}
