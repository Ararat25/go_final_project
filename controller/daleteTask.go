package controller

import (
	"net/http"
)

// DeleteTask обработчик для удаления задачи
func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")

	idInt, err := checkID(id)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.DB.DeleteTask(idInt)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	sendSuccessResponse(w, http.StatusOK)
}
