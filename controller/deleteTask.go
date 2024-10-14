package controller

import (
	"net/http"
)

// DeleteTask обработчик для удаления задачи
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idString := r.FormValue("id")

	id, err := checkID(idString)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DeleteTask(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK)
}
